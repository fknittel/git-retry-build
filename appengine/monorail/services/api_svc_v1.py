# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""API service.

To manually test this API locally, use the following steps:
1. Start the development server via 'make serve'.
2. Start a new Chrome session via the command-line:
  PATH_TO_CHROME --user-data-dir=/tmp/test \
  --unsafely-treat-insecure-origin-as-secure=http://localhost:8080
3. Visit http://localhost:8080/_ah/api/explorer
4. Click shield icon in the omnibar and allow unsafe scripts.
5. Click on the "Services" menu item in the API Explorer.
"""
from __future__ import print_function
from __future__ import division
from __future__ import absolute_import

import calendar
import datetime
import endpoints
import functools
import logging
import re
import time
from google.appengine.api import oauth
from protorpc import message_types
from protorpc import protojson
from protorpc import remote

import settings
from businesslogic import work_env
from features import filterrules_helpers
from features import send_notifications
from framework import authdata
from framework import exceptions
from framework import framework_constants
from framework import framework_helpers
from framework import framework_views
from framework import monitoring
from framework import monorailrequest
from framework import permissions
from framework import ratelimiter
from framework import sql
from project import project_helpers
from proto import api_pb2_v1
from proto import project_pb2
from proto import tracker_pb2
from search import frontendsearchpipeline
from services import api_pb2_v1_helpers
from services import client_config_svc
from services import service_manager
from services import tracker_fulltext
from sitewide import sitewide_helpers
from tracker import field_helpers
from tracker import issuedetailezt
from tracker import tracker_bizobj
from tracker import tracker_constants
from tracker import tracker_helpers

from infra_libs import ts_mon


ENDPOINTS_API_NAME = 'monorail'
DOC_URL = (
    'https://chromium.googlesource.com/infra/infra/+/main/'
    'appengine/monorail/doc/api.md')


def monorail_api_method(
    request_message, response_message, **kwargs):
  """Extends endpoints.method by performing base checks."""
  time_fn = kwargs.pop('time_fn', time.time)
  method_name = kwargs.get('name', '')
  method_path = kwargs.get('path', '')
  http_method = kwargs.get('http_method', '')
  def new_decorator(func):
    @endpoints.method(request_message, response_message, **kwargs)
    @functools.wraps(func)
    def wrapper(self, *args, **kwargs):
      start_time = time_fn()
      approximate_http_status = 200
      request = args[0]
      ret = None
      c_id = None
      c_email = None
      mar = None
      try:
        if settings.read_only and http_method.lower() != 'get':
          raise permissions.PermissionException(
              'This request is not allowed in read-only mode')
        requester = endpoints.get_current_user()
        logging.info('requester is %r', requester)
        logging.info('args is %r', args)
        logging.info('kwargs is %r', kwargs)
        auth_client_ids, auth_emails = (
            client_config_svc.GetClientConfigSvc().GetClientIDEmails())
        if settings.local_mode:
          auth_client_ids.append(endpoints.API_EXPLORER_CLIENT_ID)
        if self._services is None:
          self._set_services(service_manager.set_up_services())
        cnxn = sql.MonorailConnection()
        c_id, c_email = api_base_checks(
            request, requester, self._services, cnxn,
            auth_client_ids, auth_emails)
        mar = self.mar_factory(request, cnxn)
        self.ratelimiter.CheckStart(c_id, c_email, start_time)
        monitoring.IncrementAPIRequestsCount(
            'endpoints', c_id, client_email=c_email, handler=method_name)
        ret = func(self, mar, *args, **kwargs)
      except exceptions.NoSuchUserException as e:
        approximate_http_status = 404
        raise endpoints.NotFoundException(
            'The user does not exist: %s' % str(e))
      except (exceptions.NoSuchProjectException,
              exceptions.NoSuchIssueException,
              exceptions.NoSuchComponentException) as e:
        approximate_http_status = 404
        raise endpoints.NotFoundException(str(e))
      except (permissions.BannedUserException,
              permissions.PermissionException) as e:
        approximate_http_status = 403
        logging.info('Allowlist ID %r email %r', auth_client_ids, auth_emails)
        raise endpoints.ForbiddenException(str(e))
      except endpoints.BadRequestException:
        approximate_http_status = 400
        raise
      except endpoints.UnauthorizedException:
        approximate_http_status = 401
        # Client will refresh token and retry.
        raise
      except oauth.InvalidOAuthTokenError:
        approximate_http_status = 401
        # Client will refresh token and retry.
        raise endpoints.UnauthorizedException(
            'Auth error: InvalidOAuthTokenError')
      except (exceptions.GroupExistsException,
              exceptions.InvalidComponentNameException,
              ratelimiter.ApiRateLimitExceeded) as e:
        approximate_http_status = 400
        raise endpoints.BadRequestException(str(e))
      except Exception as e:
        approximate_http_status = 500
        logging.exception('Unexpected error in monorail API')
        raise
      finally:
        if mar:
          mar.CleanUp()
        now = time_fn()
        if c_id and c_email:
          self.ratelimiter.CheckEnd(c_id, c_email, now, start_time)
        _RecordMonitoringStats(
            start_time, request, ret, (method_name or func.__name__),
            (method_path or func.__name__), approximate_http_status, now)

      return ret

    return wrapper
  return new_decorator


def _RecordMonitoringStats(
    start_time,
    request,
    response,
    method_name,
    method_path,
    approximate_http_status,
    now=None):
  now = now or time.time()
  elapsed_ms = int((now - start_time) * 1000)
  # Use the api name, not the request path, to prevent an explosion in
  # possible field values.
  method_identifier = (
      ENDPOINTS_API_NAME + '.endpoints.' + method_name + '/' + method_path)

  # Endpoints APIs don't return the full set of http status values.
  fields = monitoring.GetCommonFields(
      approximate_http_status, method_identifier)

  monitoring.AddServerDurations(elapsed_ms, fields)
  monitoring.IncrementServerResponseStatusCount(fields)
  request_length = len(protojson.encode_message(request))
  monitoring.AddServerRequesteBytes(request_length, fields)
  response_length = 0
  if response:
    response_length = len(protojson.encode_message(response))
  monitoring.AddServerResponseBytes(response_length, fields)


def _is_requester_in_allowed_domains(requester):
  if requester.email().endswith(settings.api_allowed_email_domains):
    if framework_constants.MONORAIL_SCOPE in oauth.get_authorized_scopes(
        framework_constants.MONORAIL_SCOPE):
      return True
    else:
      logging.info("User is not authenticated with monorail scope")
  return False

def api_base_checks(request, requester, services, cnxn,
                    auth_client_ids, auth_emails):
  """Base checks for API users.

  Args:
    request: The HTTP request from Cloud Endpoints.
    requester: The user who sends the request.
    services: Services object.
    cnxn: connection to the SQL database.
    auth_client_ids: authorized client ids.
    auth_emails: authorized emails when client is anonymous.

  Returns:
    Client ID and client email.

  Raises:
    endpoints.UnauthorizedException: If the requester is anonymous.
    exceptions.NoSuchUserException: If the requester does not exist in Monorail.
    NoSuchProjectException: If the project does not exist in Monorail.
    permissions.BannedUserException: If the requester is banned.
    permissions.PermissionException: If the requester does not have
        permisssion to view.
  """
  valid_user = False
  auth_err = ''
  client_id = None

  try:
    client_id = oauth.get_client_id(framework_constants.OAUTH_SCOPE)
    logging.info('Oauth client ID %s', client_id)
  except oauth.Error as ex:
    auth_err = 'oauth.Error: %s' % ex

  if not requester:
    try:
      requester = oauth.get_current_user(framework_constants.OAUTH_SCOPE)
      logging.info('Oauth requester %s', requester.email())
    except oauth.Error as ex:
      logging.info('Got oauth error: %r', ex)
      auth_err = 'oauth.Error: %s' % ex

  if client_id and requester:
    if client_id in auth_client_ids:
      # A allowlisted client app can make requests for any user or anon.
      logging.info('Client ID %r is allowlisted', client_id)
      valid_user = True
    elif requester.email() in auth_emails:
      # A allowlisted user account can make requests via any client app.
      logging.info('Client email %r is allowlisted', requester.email())
      valid_user = True
    elif _is_requester_in_allowed_domains(requester):
      # A user with an allowed-domain email and authenticated with the
      # monorail scope is allowed to make requests via any client app.
      logging.info(
          'User email %r is within the allowed domains', requester.email())
      valid_user = True
    else:
      auth_err = (
          'Neither client ID %r nor email %r is allowlisted' %
          (client_id, requester.email()))

  if not valid_user:
    raise endpoints.UnauthorizedException('Auth error: %s' % auth_err)
  else:
    logging.info('API request from user %s:%s', client_id, requester.email())

  project_name = None
  if hasattr(request, 'projectId'):
    project_name = request.projectId
  issue_local_id = None
  if hasattr(request, 'issueId'):
    issue_local_id = request.issueId
  # This could raise exceptions.NoSuchUserException
  requester_id = services.user.LookupUserID(cnxn, requester.email())
  auth = authdata.AuthData.FromUserID(cnxn, requester_id, services)
  if permissions.IsBanned(auth.user_pb, auth.user_view):
    raise permissions.BannedUserException(
        'The user %s has been banned from using Monorail' %
        requester.email())
  if project_name:
    project = services.project.GetProjectByName(
        cnxn, project_name)
    if not project:
      raise exceptions.NoSuchProjectException(
          'Project %s does not exist' % project_name)
    if project.state != project_pb2.ProjectState.LIVE:
      raise permissions.PermissionException(
          'API may not access project %s because it is not live'
          % project_name)
    if not permissions.UserCanViewProject(
        auth.user_pb, auth.effective_ids, project):
      raise permissions.PermissionException(
          'The user %s has no permission for project %s' %
          (requester.email(), project_name))
    if issue_local_id:
      # This may raise a NoSuchIssueException.
      issue = services.issue.GetIssueByLocalID(
          cnxn, project.project_id, issue_local_id)
      perms = permissions.GetPermissions(
          auth.user_pb, auth.effective_ids, project)
      config = services.config.GetProjectConfig(cnxn, project.project_id)
      granted_perms = tracker_bizobj.GetGrantedPerms(
          issue, auth.effective_ids, config)
      if not permissions.CanViewIssue(
          auth.effective_ids, perms, project, issue,
          granted_perms=granted_perms):
        raise permissions.PermissionException(
            'User is not allowed to view this issue %s:%d' %
            (project_name, issue_local_id))

  return client_id, requester.email()


@endpoints.api(name=ENDPOINTS_API_NAME, version='v1',
               description='Monorail API to manage issues.',
               auth_level=endpoints.AUTH_LEVEL.NONE,
               allowed_client_ids=endpoints.SKIP_CLIENT_ID_CHECK,
               documentation=DOC_URL)
class MonorailApi(remote.Service):

  # Class variables. Handy to mock.
  _services = None
  _mar = None

  ratelimiter = ratelimiter.ApiRateLimiter()

  @classmethod
  def _set_services(cls, services):
    cls._services = services

  def mar_factory(self, request, cnxn):
    if not self._mar:
      self._mar = monorailrequest.MonorailApiRequest(
          request, self._services, cnxn=cnxn)
    return self._mar

  def aux_delete_comment(self, mar, request, delete=True):
    action_name = 'delete' if delete else 'undelete'

    with work_env.WorkEnv(mar, self._services) as we:
      issue = we.GetIssueByLocalID(
          mar.project_id, request.issueId, use_cache=False)
      all_comments = we.ListIssueComments(issue)
      try:
        issue_comment = all_comments[request.commentId]
      except IndexError:
        raise exceptions.NoSuchIssueException(
              'The issue %s:%d does not have comment %d.' %
              (mar.project_name, request.issueId, request.commentId))

      issue_perms = permissions.UpdateIssuePermissions(
          mar.perms, mar.project, issue, mar.auth.effective_ids,
          granted_perms=mar.granted_perms)
      commenter = we.GetUser(issue_comment.user_id)

      if not permissions.CanDeleteComment(
          issue_comment, commenter, mar.auth.user_id, issue_perms):
        raise permissions.PermissionException(
              'User is not allowed to %s the comment %d of issue %s:%d' %
              (action_name, request.commentId, mar.project_name,
               request.issueId))

      we.DeleteComment(issue, issue_comment, delete=delete)
    return api_pb2_v1.IssuesCommentsDeleteResponse()

  @monorail_api_method(
      api_pb2_v1.ISSUES_COMMENTS_DELETE_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.IssuesCommentsDeleteResponse,
      path='projects/{projectId}/issues/{issueId}/comments/{commentId}',
      http_method='DELETE',
      name='issues.comments.delete')
  def issues_comments_delete(self, mar, request):
    """Delete a comment."""
    return self.aux_delete_comment(mar, request, True)

  def parse_imported_reporter(self, mar, request):
    """Handle the case where an API client is importing issues for users.

    Args:
      mar: monorail API request object including auth and perms.
      request: A request PB that defines author and published fields.

    Returns:
      A pair (reporter_id, timestamp) with the user ID of the user to
      attribute the comment to and timestamp of the original comment.
      If the author field is not set, this is not an import request
      and the comment is attributed to the API client as per normal.
      An API client that is attempting to post on behalf of other
      users must have the ImportComment permission in the current
      project.
    """
    reporter_id = mar.auth.user_id
    timestamp = None
    if (request.author and request.author.name and
        request.author.name != mar.auth.email):
      if not mar.perms.HasPerm(
          permissions.IMPORT_COMMENT, mar.auth.user_id, mar.project):
        logging.info('name is %r', request.author.name)
        raise permissions.PermissionException(
            'User is not allowed to attribue comments to others')
      reporter_id = self._services.user.LookupUserID(
              mar.cnxn, request.author.name, autocreate=True)
      logging.info('Importing issue or comment.')
      if request.published:
        timestamp = calendar.timegm(request.published.utctimetuple())

    return reporter_id, timestamp

  @monorail_api_method(
      api_pb2_v1.ISSUES_COMMENTS_INSERT_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.IssuesCommentsInsertResponse,
      path='projects/{projectId}/issues/{issueId}/comments',
      http_method='POST',
      name='issues.comments.insert')
  def issues_comments_insert(self, mar, request):
    # type (...) -> proto.api_pb2_v1.IssuesCommentsInsertResponse
    """Add a comment."""
    # Because we will modify issues, load from DB rather than cache.
    issue = self._services.issue.GetIssueByLocalID(
        mar.cnxn, mar.project_id, request.issueId, use_cache=False)
    old_owner_id = tracker_bizobj.GetOwnerId(issue)
    if not permissions.CanCommentIssue(
        mar.auth.effective_ids, mar.perms, mar.project, issue,
        mar.granted_perms):
      raise permissions.PermissionException(
          'User is not allowed to comment this issue (%s, %d)' %
          (request.projectId, request.issueId))

    # Temporary block on updating approval subfields.
    if request.updates and request.updates.fieldValues:
      fds_by_name = {fd.field_name.lower():fd for fd in mar.config.field_defs}
      for fv in request.updates.fieldValues:
        # Checking for fv.approvalName is unreliable since it can be removed.
        fd = fds_by_name.get(fv.fieldName.lower())
        if fd and fd.approval_id:
          raise exceptions.ActionNotSupported(
              'No API support for approval field changes: (approval %s owns %s)'
              % (fd.approval_id, fd.field_name))
        # if fd was None, that gets dealt with later.

    if request.content and len(
        request.content) > tracker_constants.MAX_COMMENT_CHARS:
      raise endpoints.BadRequestException(
          'Comment is too long on this issue (%s, %d' %
          (request.projectId, request.issueId))

    updates_dict = {}
    move_to_project = None
    if request.updates:
      if not permissions.CanEditIssue(
          mar.auth.effective_ids, mar.perms, mar.project, issue,
          mar.granted_perms):
        raise permissions.PermissionException(
            'User is not allowed to edit this issue (%s, %d)' %
            (request.projectId, request.issueId))
      if request.updates.moveToProject:
        move_to = request.updates.moveToProject.lower()
        move_to_project = issuedetailezt.CheckMoveIssueRequest(
            self._services, mar, issue, True, move_to, mar.errors)
        if mar.errors.AnyErrors():
          raise endpoints.BadRequestException(mar.errors.move_to)

      updates_dict['summary'] = request.updates.summary
      updates_dict['status'] = request.updates.status
      updates_dict['is_description'] = request.updates.is_description
      if request.updates.owner:
        # A current issue owner can be removed via the API with a
        # NO_USER_NAME('----') input.
        if request.updates.owner == framework_constants.NO_USER_NAME:
          updates_dict['owner'] = framework_constants.NO_USER_SPECIFIED
        else:
          new_owner_id = self._services.user.LookupUserID(
              mar.cnxn, request.updates.owner)
          valid, msg = tracker_helpers.IsValidIssueOwner(
              mar.cnxn, mar.project, new_owner_id, self._services)
          if not valid:
            raise endpoints.BadRequestException(msg)
          updates_dict['owner'] = new_owner_id
      updates_dict['cc_add'], updates_dict['cc_remove'] = (
          api_pb2_v1_helpers.split_remove_add(request.updates.cc))
      updates_dict['cc_add'] = list(self._services.user.LookupUserIDs(
          mar.cnxn, updates_dict['cc_add'], autocreate=True).values())
      updates_dict['cc_remove'] = list(self._services.user.LookupUserIDs(
          mar.cnxn, updates_dict['cc_remove']).values())
      updates_dict['labels_add'], updates_dict['labels_remove'] = (
          api_pb2_v1_helpers.split_remove_add(request.updates.labels))
      blocked_on_add_strs, blocked_on_remove_strs = (
          api_pb2_v1_helpers.split_remove_add(request.updates.blockedOn))
      updates_dict['blocked_on_add'] = api_pb2_v1_helpers.issue_global_ids(
          blocked_on_add_strs, issue.project_id, mar,
          self._services)
      updates_dict['blocked_on_remove'] = api_pb2_v1_helpers.issue_global_ids(
          blocked_on_remove_strs, issue.project_id, mar,
          self._services)
      blocking_add_strs, blocking_remove_strs = (
          api_pb2_v1_helpers.split_remove_add(request.updates.blocking))
      updates_dict['blocking_add'] = api_pb2_v1_helpers.issue_global_ids(
          blocking_add_strs, issue.project_id, mar,
          self._services)
      updates_dict['blocking_remove'] = api_pb2_v1_helpers.issue_global_ids(
          blocking_remove_strs, issue.project_id, mar,
          self._services)
      components_add_strs, components_remove_strs = (
          api_pb2_v1_helpers.split_remove_add(request.updates.components))
      updates_dict['components_add'] = (
          api_pb2_v1_helpers.convert_component_ids(
              mar.config, components_add_strs))
      updates_dict['components_remove'] = (
          api_pb2_v1_helpers.convert_component_ids(
              mar.config, components_remove_strs))
      if request.updates.mergedInto:
        merge_project_name, merge_local_id = tracker_bizobj.ParseIssueRef(
            request.updates.mergedInto)
        merge_into_project = self._services.project.GetProjectByName(
            mar.cnxn, merge_project_name or issue.project_name)
        # Because we will modify issues, load from DB rather than cache.
        merge_into_issue = self._services.issue.GetIssueByLocalID(
            mar.cnxn, merge_into_project.project_id, merge_local_id,
            use_cache=False)
        merge_allowed = tracker_helpers.IsMergeAllowed(
            merge_into_issue, mar, self._services)
        if not merge_allowed:
          raise permissions.PermissionException(
            'User is not allowed to merge into issue %s:%s' %
            (merge_into_issue.project_name, merge_into_issue.local_id))
        updates_dict['merged_into'] = merge_into_issue.issue_id
      (updates_dict['field_vals_add'], updates_dict['field_vals_remove'],
       updates_dict['fields_clear'], updates_dict['fields_labels_add'],
       updates_dict['fields_labels_remove']) = (
          api_pb2_v1_helpers.convert_field_values(
              request.updates.fieldValues, mar, self._services))

    field_helpers.ValidateCustomFields(
        mar.cnxn, self._services,
        (updates_dict.get('field_vals_add', []) +
         updates_dict.get('field_vals_remove', [])),
        mar.config, mar.project, ezt_errors=mar.errors)
    if mar.errors.AnyErrors():
      raise endpoints.BadRequestException(
          'Invalid field values: %s' % mar.errors.custom_fields)

    updates_dict['labels_add'] = (
        updates_dict.get('labels_add', []) +
        updates_dict.get('fields_labels_add', []))
    updates_dict['labels_remove'] = (
        updates_dict.get('labels_remove', []) +
        updates_dict.get('fields_labels_remove', []))

    # TODO(jrobbins): Stop using updates_dict in the first place.
    delta = tracker_bizobj.MakeIssueDelta(
        updates_dict.get('status'),
        updates_dict.get('owner'),
        updates_dict.get('cc_add', []),
        updates_dict.get('cc_remove', []),
        updates_dict.get('components_add', []),
        updates_dict.get('components_remove', []),
        (updates_dict.get('labels_add', []) +
         updates_dict.get('fields_labels_add', [])),
        (updates_dict.get('labels_remove', []) +
         updates_dict.get('fields_labels_remove', [])),
        updates_dict.get('field_vals_add', []),
        updates_dict.get('field_vals_remove', []),
        updates_dict.get('fields_clear', []),
        updates_dict.get('blocked_on_add', []),
        updates_dict.get('blocked_on_remove', []),
        updates_dict.get('blocking_add', []),
        updates_dict.get('blocking_remove', []),
        updates_dict.get('merged_into'),
        updates_dict.get('summary'))

    importer_id = None
    reporter_id, timestamp = self.parse_imported_reporter(mar, request)
    if reporter_id != mar.auth.user_id:
      importer_id = mar.auth.user_id

    # TODO(jrobbins): Finish refactoring to make everything go through work_env.
    _, comment = self._services.issue.DeltaUpdateIssue(
        cnxn=mar.cnxn, services=self._services,
        reporter_id=reporter_id, project_id=mar.project_id, config=mar.config,
        issue=issue, delta=delta, index_now=False, comment=request.content,
        is_description=updates_dict.get('is_description'),
        timestamp=timestamp, importer_id=importer_id)

    move_comment = None
    if move_to_project:
      old_text_ref = 'issue %s:%s' % (issue.project_name, issue.local_id)
      tracker_fulltext.UnindexIssues([issue.issue_id])
      moved_back_iids = self._services.issue.MoveIssues(
          mar.cnxn, move_to_project, [issue], self._services.user)
      new_text_ref = 'issue %s:%s' % (issue.project_name, issue.local_id)
      if issue.issue_id in moved_back_iids:
        content = 'Moved %s back to %s again.' % (old_text_ref, new_text_ref)
      else:
        content = 'Moved %s to now be %s.' % (old_text_ref, new_text_ref)
      move_comment = self._services.issue.CreateIssueComment(
        mar.cnxn, issue, mar.auth.user_id, content, amendments=[
            tracker_bizobj.MakeProjectAmendment(move_to_project.project_name)])

    if 'merged_into' in updates_dict:
      new_starrers = tracker_helpers.GetNewIssueStarrers(
          mar.cnxn, self._services, [issue.issue_id], merge_into_issue.issue_id)
      tracker_helpers.AddIssueStarrers(
          mar.cnxn, self._services, mar,
          merge_into_issue.issue_id, merge_into_project, new_starrers)
      # Load target issue again to get the updated star count.
      merge_into_issue = self._services.issue.GetIssue(
        mar.cnxn, merge_into_issue.issue_id, use_cache=False)
      merge_comment_pb = tracker_helpers.MergeCCsAndAddComment(
        self._services, mar, issue, merge_into_issue)
      hostport = framework_helpers.GetHostPort(
          project_name=merge_into_issue.project_name)
      send_notifications.PrepareAndSendIssueChangeNotification(
          merge_into_issue.issue_id, hostport,
          mar.auth.user_id, send_email=True, comment_id=merge_comment_pb.id)

    tracker_fulltext.IndexIssues(
        mar.cnxn, [issue], self._services.user, self._services.issue,
        self._services.config)

    comment = comment or move_comment
    if comment is None:
      return api_pb2_v1.IssuesCommentsInsertResponse()

    cmnts = self._services.issue.GetCommentsForIssue(mar.cnxn, issue.issue_id)
    seq = len(cmnts) - 1

    if request.sendEmail:
      hostport = framework_helpers.GetHostPort(project_name=issue.project_name)
      send_notifications.PrepareAndSendIssueChangeNotification(
          issue.issue_id, hostport, comment.user_id, send_email=True,
          old_owner_id=old_owner_id, comment_id=comment.id)

    issue_perms = permissions.UpdateIssuePermissions(
        mar.perms, mar.project, issue, mar.auth.effective_ids,
        granted_perms=mar.granted_perms)
    commenter = self._services.user.GetUser(mar.cnxn, comment.user_id)
    can_delete = permissions.CanDeleteComment(
        comment, commenter, mar.auth.user_id, issue_perms)
    return api_pb2_v1.IssuesCommentsInsertResponse(
        id=seq,
        kind='monorail#issueComment',
        author=api_pb2_v1_helpers.convert_person(
            comment.user_id, mar.cnxn, self._services),
        content=comment.content,
        published=datetime.datetime.fromtimestamp(comment.timestamp),
        updates=api_pb2_v1_helpers.convert_amendments(
            issue, comment.amendments, mar, self._services),
        canDelete=can_delete)

  @monorail_api_method(
      api_pb2_v1.ISSUES_COMMENTS_LIST_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.IssuesCommentsListResponse,
      path='projects/{projectId}/issues/{issueId}/comments',
      http_method='GET',
      name='issues.comments.list')
  def issues_comments_list(self, mar, request):
    """List all comments for an issue."""
    issue = self._services.issue.GetIssueByLocalID(
        mar.cnxn, mar.project_id, request.issueId)
    comments = self._services.issue.GetCommentsForIssue(
        mar.cnxn, issue.issue_id)
    comments = [comment for comment in comments if not comment.approval_id]
    visible_comments = []
    for comment in comments[
        request.startIndex:(request.startIndex + request.maxResults)]:
      visible_comments.append(
          api_pb2_v1_helpers.convert_comment(
              issue, comment, mar, self._services, mar.granted_perms))

    return api_pb2_v1.IssuesCommentsListResponse(
        kind='monorail#issueCommentList',
        totalResults=len(comments),
        items=visible_comments)

  @monorail_api_method(
      api_pb2_v1.ISSUES_COMMENTS_DELETE_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.IssuesCommentsDeleteResponse,
      path='projects/{projectId}/issues/{issueId}/comments/{commentId}',
      http_method='POST',
      name='issues.comments.undelete')
  def issues_comments_undelete(self, mar, request):
    """Restore a deleted comment."""
    return self.aux_delete_comment(mar, request, False)

  @monorail_api_method(
      api_pb2_v1.APPROVALS_COMMENTS_LIST_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.ApprovalsCommentsListResponse,
      path='projects/{projectId}/issues/{issueId}/'
            'approvals/{approvalName}/comments',
      http_method='GET',
      name='approvals.comments.list')
  def approvals_comments_list(self, mar, request):
    """List all comments for an issue approval."""
    issue = self._services.issue.GetIssueByLocalID(
        mar.cnxn, mar.project_id, request.issueId)
    if not permissions.CanViewIssue(
        mar.auth.effective_ids, mar.perms, mar.project, issue,
        mar.granted_perms):
      raise permissions.PermissionException(
          'User is not allowed to view this issue (%s, %d)' %
          (request.projectId, request.issueId))
    config = self._services.config.GetProjectConfig(mar.cnxn, issue.project_id)
    approval_fd = tracker_bizobj.FindFieldDef(request.approvalName, config)
    if not approval_fd:
      raise endpoints.BadRequestException(
          'Field definition for %s not found in project config' %
          request.approvalName)
    comments = self._services.issue.GetCommentsForIssue(
        mar.cnxn, issue.issue_id)
    comments = [comment for comment in comments
                if comment.approval_id == approval_fd.field_id]
    visible_comments = []
    for comment in comments[
        request.startIndex:(request.startIndex + request.maxResults)]:
      visible_comments.append(
          api_pb2_v1_helpers.convert_approval_comment(
              issue, comment, mar, self._services, mar.granted_perms))

    return api_pb2_v1.ApprovalsCommentsListResponse(
        kind='monorail#approvalCommentList',
        totalResults=len(comments),
        items=visible_comments)

  @monorail_api_method(
      api_pb2_v1.APPROVALS_COMMENTS_INSERT_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.ApprovalsCommentsInsertResponse,
      path=("projects/{projectId}/issues/{issueId}/"
            "approvals/{approvalName}/comments"),
      http_method='POST',
      name='approvals.comments.insert')
  def approvals_comments_insert(self, mar, request):
    # type (...) -> proto.api_pb2_v1.ApprovalsCommentsInsertResponse
    """Add an approval comment."""
    approval_fd = tracker_bizobj.FindFieldDef(
        request.approvalName, mar.config)
    if not approval_fd or (
        approval_fd.field_type != tracker_pb2.FieldTypes.APPROVAL_TYPE):
      raise endpoints.BadRequestException(
          'Field definition for %s not found in project config' %
          request.approvalName)
    try:
      issue = self._services.issue.GetIssueByLocalID(
          mar.cnxn, mar.project_id, request.issueId)
    except exceptions.NoSuchIssueException:
      raise endpoints.BadRequestException(
          'Issue %s:%s not found' % (request.projectId, request.issueId))
    approval = tracker_bizobj.FindApprovalValueByID(
        approval_fd.field_id, issue.approval_values)
    if not approval:
      raise endpoints.BadRequestException(
          'Approval %s not found in issue.' % request.approvalName)

    if not permissions.CanCommentIssue(
        mar.auth.effective_ids, mar.perms, mar.project, issue,
        mar.granted_perms):
      raise permissions.PermissionException(
          'User is not allowed to comment on this issue (%s, %d)' %
          (request.projectId, request.issueId))

    if request.content and len(
        request.content) > tracker_constants.MAX_COMMENT_CHARS:
      raise endpoints.BadRequestException(
          'Comment is too long on this issue (%s, %d' %
          (request.projectId, request.issueId))

    updates_dict = {}
    if request.approvalUpdates:
      if request.approvalUpdates.fieldValues:
        # Block updating field values that don't belong to the approval.
        approvals_fds_by_name = {
            fd.field_name.lower():fd for fd in mar.config.field_defs
            if fd.approval_id == approval_fd.field_id}
        for fv in request.approvalUpdates.fieldValues:
          if approvals_fds_by_name.get(fv.fieldName.lower()) is None:
            raise endpoints.BadRequestException(
              'Field defition for %s not found in %s subfields.' %
              (fv.fieldName, request.approvalName))
        (updates_dict['field_vals_add'], updates_dict['field_vals_remove'],
         updates_dict['fields_clear'], updates_dict['fields_labels_add'],
         updates_dict['fields_labels_remove']) = (
             api_pb2_v1_helpers.convert_field_values(
                 request.approvalUpdates.fieldValues, mar, self._services))
      if request.approvalUpdates.approvers:
        if not permissions.CanUpdateApprovers(
            mar.auth.effective_ids, mar.perms, mar.project,
            approval.approver_ids):
          raise permissions.PermissionException(
              'User is not allowed to update approvers')
        approvers_add, approvers_remove = api_pb2_v1_helpers.split_remove_add(
            request.approvalUpdates.approvers)
        updates_dict['approver_ids_add'] = list(
            self._services.user.LookupUserIDs(mar.cnxn, approvers_add,
              autocreate=True).values())
        updates_dict['approver_ids_remove'] = list(
            self._services.user.LookupUserIDs(mar.cnxn, approvers_remove,
              autocreate=True).values())
      if request.approvalUpdates.status:
        status = tracker_pb2.ApprovalStatus(
            api_pb2_v1.ApprovalStatus(request.approvalUpdates.status).number)
        if not permissions.CanUpdateApprovalStatus(
            mar.auth.effective_ids, mar.perms, mar.project,
            approval.approver_ids, status):
          raise permissions.PermissionException(
              'User is not allowed to make this status change')
        updates_dict['status'] = status
    logging.info(time.time)
    approval_delta = tracker_bizobj.MakeApprovalDelta(
        updates_dict.get('status'), mar.auth.user_id,
        updates_dict.get('approver_ids_add', []),
        updates_dict.get('approver_ids_remove', []),
        updates_dict.get('field_vals_add', []),
        updates_dict.get('field_vals_remove', []),
        updates_dict.get('fields_clear', []),
        updates_dict.get('fields_labels_add', []),
        updates_dict.get('fields_labels_remove', []))
    comment = self._services.issue.DeltaUpdateIssueApproval(
        mar.cnxn, mar.auth.user_id, mar.config, issue, approval, approval_delta,
        comment_content=request.content,
        is_description=request.is_description)

    cmnts = self._services.issue.GetCommentsForIssue(mar.cnxn, issue.issue_id)
    seq = len(cmnts) - 1

    if request.sendEmail:
      hostport = framework_helpers.GetHostPort(project_name=issue.project_name)
      send_notifications.PrepareAndSendApprovalChangeNotification(
          issue.issue_id, approval.approval_id,
          hostport, comment.id, send_email=True)

    issue_perms = permissions.UpdateIssuePermissions(
        mar.perms, mar.project, issue, mar.auth.effective_ids,
        granted_perms=mar.granted_perms)
    commenter = self._services.user.GetUser(mar.cnxn, comment.user_id)
    can_delete = permissions.CanDeleteComment(
        comment, commenter, mar.auth.user_id, issue_perms)
    return api_pb2_v1.ApprovalsCommentsInsertResponse(
        id=seq,
        kind='monorail#approvalComment',
        author=api_pb2_v1_helpers.convert_person(
            comment.user_id, mar.cnxn, self._services),
        content=comment.content,
        published=datetime.datetime.fromtimestamp(comment.timestamp),
        approvalUpdates=api_pb2_v1_helpers.convert_approval_amendments(
            comment.amendments, mar, self._services),
        canDelete=can_delete)

  @monorail_api_method(
      api_pb2_v1.USERS_GET_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.UsersGetResponse,
      path='users/{userId}',
      http_method='GET',
      name='users.get')
  def users_get(self, mar, request):
    """Get a user."""
    owner_project_only = request.ownerProjectsOnly
    with work_env.WorkEnv(mar, self._services) as we:
      (visible_ownership, visible_deleted, visible_membership,
       visible_contrib) = we.GetUserProjects(
           mar.viewed_user_auth.effective_ids)

    project_list = []
    for proj in (visible_ownership + visible_deleted):
      config = self._services.config.GetProjectConfig(
          mar.cnxn, proj.project_id)
      templates = self._services.template.GetProjectTemplates(
          mar.cnxn, config.project_id)
      proj_result = api_pb2_v1_helpers.convert_project(
          proj, config, api_pb2_v1.Role.owner, templates)
      project_list.append(proj_result)
    if not owner_project_only:
      for proj in visible_membership:
        config = self._services.config.GetProjectConfig(
            mar.cnxn, proj.project_id)
        templates = self._services.template.GetProjectTemplates(
            mar.cnxn, config.project_id)
        proj_result = api_pb2_v1_helpers.convert_project(
            proj, config, api_pb2_v1.Role.member, templates)
        project_list.append(proj_result)
      for proj in visible_contrib:
        config = self._services.config.GetProjectConfig(
            mar.cnxn, proj.project_id)
        templates = self._services.template.GetProjectTemplates(
            mar.cnxn, config.project_id)
        proj_result = api_pb2_v1_helpers.convert_project(
            proj, config, api_pb2_v1.Role.contributor, templates)
        project_list.append(proj_result)

    return api_pb2_v1.UsersGetResponse(
        id=str(mar.viewed_user_auth.user_id),
        kind='monorail#user',
        projects=project_list,
    )

  @monorail_api_method(
      api_pb2_v1.ISSUES_GET_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.IssuesGetInsertResponse,
      path='projects/{projectId}/issues/{issueId}',
      http_method='GET',
      name='issues.get')
  def issues_get(self, mar, request):
    """Get an issue."""
    issue = self._services.issue.GetIssueByLocalID(
        mar.cnxn, mar.project_id, request.issueId)

    return api_pb2_v1_helpers.convert_issue(
        api_pb2_v1.IssuesGetInsertResponse, issue, mar, self._services)

  @monorail_api_method(
      api_pb2_v1.ISSUES_INSERT_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.IssuesGetInsertResponse,
      path='projects/{projectId}/issues',
      http_method='POST',
      name='issues.insert')
  def issues_insert(self, mar, request):
    """Add a new issue."""
    if not mar.perms.CanUsePerm(
        permissions.CREATE_ISSUE, mar.auth.effective_ids, mar.project, []):
      raise permissions.PermissionException(
          'The requester %s is not allowed to create issues for project %s.' %
          (mar.auth.email, mar.project_name))

    with work_env.WorkEnv(mar, self._services) as we:
      owner_id = framework_constants.NO_USER_SPECIFIED
      if request.owner and request.owner.name:
        try:
          owner_id = self._services.user.LookupUserID(
              mar.cnxn, request.owner.name)
        except exceptions.NoSuchUserException:
          raise endpoints.BadRequestException(
              'The specified owner %s does not exist.' % request.owner.name)

      cc_ids = []
      request.cc = [cc for cc in request.cc if cc]
      if request.cc:
        cc_ids = list(self._services.user.LookupUserIDs(
            mar.cnxn, [ap.name for ap in request.cc],
            autocreate=True).values())
      comp_ids = api_pb2_v1_helpers.convert_component_ids(
          mar.config, request.components)
      fields_add, _, _, fields_labels, _ = (
          api_pb2_v1_helpers.convert_field_values(
              request.fieldValues, mar, self._services))
      field_helpers.ValidateCustomFields(
          mar.cnxn, self._services, fields_add, mar.config, mar.project,
          ezt_errors=mar.errors)
      if mar.errors.AnyErrors():
        raise endpoints.BadRequestException(
            'Invalid field values: %s' % mar.errors.custom_fields)

      logging.info('request.author is %r', request.author)
      reporter_id, timestamp = self.parse_imported_reporter(mar, request)
      # To preserve previous behavior, do not raise filter rule errors.
      try:
        new_issue, _ = we.CreateIssue(
            mar.project_id,
            request.summary,
            request.status,
            owner_id,
            cc_ids,
            request.labels + fields_labels,
            fields_add,
            comp_ids,
            request.description,
            blocked_on=api_pb2_v1_helpers.convert_issueref_pbs(
                request.blockedOn, mar, self._services),
            blocking=api_pb2_v1_helpers.convert_issueref_pbs(
                request.blocking, mar, self._services),
            reporter_id=reporter_id,
            timestamp=timestamp,
            send_email=request.sendEmail,
            raise_filter_errors=False)
        we.StarIssue(new_issue, True)
      except exceptions.InputException as e:
        raise endpoints.BadRequestException(str(e))

    return api_pb2_v1_helpers.convert_issue(
        api_pb2_v1.IssuesGetInsertResponse, new_issue, mar, self._services)

  @monorail_api_method(
      api_pb2_v1.ISSUES_LIST_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.IssuesListResponse,
      path='projects/{projectId}/issues',
      http_method='GET',
      name='issues.list')
  def issues_list(self, mar, request):
    """List issues for projects."""
    if request.additionalProject:
      for project_name in request.additionalProject:
        project = self._services.project.GetProjectByName(
            mar.cnxn, project_name)
        if project and not permissions.UserCanViewProject(
            mar.auth.user_pb, mar.auth.effective_ids, project):
          raise permissions.PermissionException(
              'The user %s has no permission for project %s' %
              (mar.auth.email, project_name))
    # TODO(jrobbins): This should go through work_env.
    pipeline = frontendsearchpipeline.FrontendSearchPipeline(
        mar.cnxn,
        self._services,
        mar.auth, [mar.me_user_id],
        mar.query,
        mar.query_project_names,
        mar.num,
        mar.start,
        mar.can,
        mar.group_by_spec,
        mar.sort_spec,
        mar.warnings,
        mar.errors,
        mar.use_cached_searches,
        mar.profiler,
        project=mar.project)
    if not mar.errors.AnyErrors():
      pipeline.SearchForIIDs()
      pipeline.MergeAndSortIssues()
      pipeline.Paginate()
    else:
      raise endpoints.BadRequestException(mar.errors.query)

    issue_list = [
        api_pb2_v1_helpers.convert_issue(
            api_pb2_v1.IssueWrapper, r, mar, self._services)
        for r in pipeline.visible_results]
    return api_pb2_v1.IssuesListResponse(
        kind='monorail#issueList',
        totalResults=pipeline.total_count,
        items=issue_list)

  @monorail_api_method(
      api_pb2_v1.GROUPS_SETTINGS_LIST_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.GroupsSettingsListResponse,
      path='groupsettings',
      http_method='GET',
      name='groups.settings.list')
  def groups_settings_list(self, mar, request):
    """List all group settings."""
    all_groups = self._services.usergroup.GetAllUserGroupsInfo(mar.cnxn)
    group_settings = []
    for g in all_groups:
      setting = g[2]
      wrapper = api_pb2_v1_helpers.convert_group_settings(g[0], setting)
      if not request.importedGroupsOnly or wrapper.ext_group_type:
        group_settings.append(wrapper)
    return api_pb2_v1.GroupsSettingsListResponse(
        groupSettings=group_settings)

  @monorail_api_method(
      api_pb2_v1.GROUPS_CREATE_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.GroupsCreateResponse,
      path='groups',
      http_method='POST',
      name='groups.create')
  def groups_create(self, mar, request):
    """Create a new user group."""
    if not permissions.CanCreateGroup(mar.perms):
      raise permissions.PermissionException(
          'The user is not allowed to create groups.')

    user_dict = self._services.user.LookupExistingUserIDs(
        mar.cnxn, [request.groupName])
    if request.groupName.lower() in user_dict:
      raise exceptions.GroupExistsException(
          'group %s already exists' % request.groupName)

    if request.ext_group_type:
      ext_group_type = str(request.ext_group_type).lower()
    else:
      ext_group_type = None
    group_id = self._services.usergroup.CreateGroup(
        mar.cnxn, self._services, request.groupName,
        str(request.who_can_view_members).lower(),
        ext_group_type)

    return api_pb2_v1.GroupsCreateResponse(
        groupID=group_id)

  @monorail_api_method(
      api_pb2_v1.GROUPS_GET_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.GroupsGetResponse,
      path='groups/{groupName}',
      http_method='GET',
      name='groups.get')
  def groups_get(self, mar, request):
    """Get a group's settings and users."""
    if not mar.viewed_user_auth:
      raise exceptions.NoSuchUserException(request.groupName)
    group_id = mar.viewed_user_auth.user_id
    group_settings = self._services.usergroup.GetGroupSettings(
        mar.cnxn, group_id)
    member_ids, owner_ids = self._services.usergroup.LookupAllMembers(
          mar.cnxn, [group_id])
    (owned_project_ids, membered_project_ids,
     contrib_project_ids) = self._services.project.GetUserRolesInAllProjects(
         mar.cnxn, mar.auth.effective_ids)
    project_ids = owned_project_ids.union(
        membered_project_ids).union(contrib_project_ids)
    if not permissions.CanViewGroupMembers(
        mar.perms, mar.auth.effective_ids, group_settings, member_ids[group_id],
        owner_ids[group_id], project_ids):
      raise permissions.PermissionException(
          'The user is not allowed to view this group.')

    member_ids, owner_ids = self._services.usergroup.LookupMembers(
        mar.cnxn, [group_id])

    member_emails = list(self._services.user.LookupUserEmails(
        mar.cnxn, member_ids[group_id]).values())
    owner_emails = list(self._services.user.LookupUserEmails(
        mar.cnxn, owner_ids[group_id]).values())

    return api_pb2_v1.GroupsGetResponse(
      groupID=group_id,
      groupSettings=api_pb2_v1_helpers.convert_group_settings(
          request.groupName, group_settings),
      groupOwners=owner_emails,
      groupMembers=member_emails)

  @monorail_api_method(
      api_pb2_v1.GROUPS_UPDATE_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.GroupsUpdateResponse,
      path='groups/{groupName}',
      http_method='POST',
      name='groups.update')
  def groups_update(self, mar, request):
    """Update a group's settings and users."""
    group_id = mar.viewed_user_auth.user_id
    member_ids_dict, owner_ids_dict = self._services.usergroup.LookupMembers(
        mar.cnxn, [group_id])
    owner_ids = owner_ids_dict.get(group_id, [])
    member_ids = member_ids_dict.get(group_id, [])
    if not permissions.CanEditGroup(
        mar.perms, mar.auth.effective_ids, owner_ids):
      raise permissions.PermissionException(
          'The user is not allowed to edit this group.')

    group_settings = self._services.usergroup.GetGroupSettings(
        mar.cnxn, group_id)
    if (request.who_can_view_members or request.ext_group_type
        or request.last_sync_time or request.friend_projects):
      group_settings.who_can_view_members = (
          request.who_can_view_members or group_settings.who_can_view_members)
      group_settings.ext_group_type = (
          request.ext_group_type or group_settings.ext_group_type)
      group_settings.last_sync_time = (
          request.last_sync_time or group_settings.last_sync_time)
      if framework_constants.NO_VALUES in request.friend_projects:
        group_settings.friend_projects = []
      else:
        id_dict = self._services.project.LookupProjectIDs(
            mar.cnxn, request.friend_projects)
        group_settings.friend_projects = (
            list(id_dict.values()) or group_settings.friend_projects)
      self._services.usergroup.UpdateSettings(
          mar.cnxn, group_id, group_settings)

    if request.groupOwners or request.groupMembers:
      self._services.usergroup.RemoveMembers(
          mar.cnxn, group_id, owner_ids + member_ids)
      owners_dict = self._services.user.LookupUserIDs(
          mar.cnxn, request.groupOwners, autocreate=True)
      self._services.usergroup.UpdateMembers(
          mar.cnxn, group_id, list(owners_dict.values()), 'owner')
      members_dict = self._services.user.LookupUserIDs(
          mar.cnxn, request.groupMembers, autocreate=True)
      self._services.usergroup.UpdateMembers(
          mar.cnxn, group_id, list(members_dict.values()), 'member')

    return api_pb2_v1.GroupsUpdateResponse()

  @monorail_api_method(
      api_pb2_v1.COMPONENTS_LIST_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.ComponentsListResponse,
      path='projects/{projectId}/components',
      http_method='GET',
      name='components.list')
  def components_list(self, mar, _request):
    """List all components of a given project."""
    config = self._services.config.GetProjectConfig(mar.cnxn, mar.project_id)
    components = [api_pb2_v1_helpers.convert_component_def(
        cd, mar, self._services) for cd in config.component_defs]
    return api_pb2_v1.ComponentsListResponse(
        components=components)

  @monorail_api_method(
      api_pb2_v1.COMPONENTS_CREATE_REQUEST_RESOURCE_CONTAINER,
      api_pb2_v1.Component,
      path='projects/{projectId}/components',
      http_method='POST',
      name='components.create')
  def components_create(self, mar, request):
    """Create a component."""
    if not mar.perms.CanUsePerm(
        permissions.EDIT_PROJECT, mar.auth.effective_ids, mar.project, []):
      raise permissions.PermissionException(
          'User is not allowed to create components for this project')

    config = self._services.config.GetProjectConfig(mar.cnxn, mar.project_id)
    leaf_name = request.componentName
    if not tracker_constants.COMPONENT_NAME_RE.match(leaf_name):
      raise exceptions.InvalidComponentNameException(
          'The component name %s is invalid.' % leaf_name)

    parent_path = request.parentPath
    if parent_path:
      parent_def = tracker_bizobj.FindComponentDef(parent_path, config)
      if not parent_def:
        raise exceptions.NoSuchComponentException(
            'Parent component %s does not exist.' % parent_path)
      if not permissions.CanEditComponentDef(
          mar.auth.effective_ids, mar.perms, mar.project, parent_def, config):
        raise permissions.PermissionException(
            'User is not allowed to add a subcomponent to component %s' %
            parent_path)

      path = '%s>%s' % (parent_path, leaf_name)
    else:
      path = leaf_name

    if tracker_bizobj.FindComponentDef(path, config):
      raise exceptions.InvalidComponentNameException(
          'The name %s is already in use.' % path)

    created = int(time.time())
    user_emails = set()
    user_emails.update([mar.auth.email] + request.admin + request.cc)
    user_ids_dict = self._services.user.LookupUserIDs(
        mar.cnxn, list(user_emails), autocreate=False)
    request.admin = [admin for admin in request.admin if admin]
    admin_ids = [user_ids_dict[uname] for uname in request.admin]
    request.cc = [cc for cc in request.cc if cc]
    cc_ids = [user_ids_dict[uname] for uname in request.cc]
    label_ids = []  # TODO(jrobbins): allow API clients to specify this too.

    component_id = self._services.config.CreateComponentDef(
        mar.cnxn, mar.project_id, path, request.description, request.deprecated,
        admin_ids, cc_ids, created, user_ids_dict[mar.auth.email], label_ids)

    return api_pb2_v1.Component(
        componentId=component_id,
        projectName=request.projectId,
        componentPath=path,
        description=request.description,
        admin=request.admin,
        cc=request.cc,
        deprecated=request.deprecated,
        created=datetime.datetime.fromtimestamp(created),
        creator=mar.auth.email)

  @monorail_api_method(
      api_pb2_v1.COMPONENTS_DELETE_REQUEST_RESOURCE_CONTAINER,
      message_types.VoidMessage,
      path='projects/{projectId}/components/{componentPath}',
      http_method='DELETE',
      name='components.delete')
  def components_delete(self, mar, request):
    """Delete a component."""
    config = self._services.config.GetProjectConfig(mar.cnxn, mar.project_id)
    component_path = request.componentPath
    component_def = tracker_bizobj.FindComponentDef(
        component_path, config)
    if not component_def:
      raise exceptions.NoSuchComponentException(
          'The component %s does not exist.' % component_path)
    if not permissions.CanViewComponentDef(
        mar.auth.effective_ids, mar.perms, mar.project, component_def):
      raise permissions.PermissionException(
          'User is not allowed to view this component %s' % component_path)
    if not permissions.CanEditComponentDef(
        mar.auth.effective_ids, mar.perms, mar.project, component_def, config):
      raise permissions.PermissionException(
          'User is not allowed to delete this component %s' % component_path)

    allow_delete = not tracker_bizobj.FindDescendantComponents(
        config, component_def)
    if not allow_delete:
      raise permissions.PermissionException(
          'User tried to delete component that had subcomponents')

    self._services.issue.DeleteComponentReferences(
        mar.cnxn, component_def.component_id)
    self._services.config.DeleteComponentDef(
        mar.cnxn, mar.project_id, component_def.component_id)
    return message_types.VoidMessage()

  @monorail_api_method(
      api_pb2_v1.COMPONENTS_UPDATE_REQUEST_RESOURCE_CONTAINER,
      message_types.VoidMessage,
      path='projects/{projectId}/components/{componentPath}',
      http_method='POST',
      name='components.update')
  def components_update(self, mar, request):
    """Update a component."""
    config = self._services.config.GetProjectConfig(mar.cnxn, mar.project_id)
    component_path = request.componentPath
    component_def = tracker_bizobj.FindComponentDef(
        component_path, config)
    if not component_def:
      raise exceptions.NoSuchComponentException(
          'The component %s does not exist.' % component_path)
    if not permissions.CanViewComponentDef(
        mar.auth.effective_ids, mar.perms, mar.project, component_def):
      raise permissions.PermissionException(
          'User is not allowed to view this component %s' % component_path)
    if not permissions.CanEditComponentDef(
        mar.auth.effective_ids, mar.perms, mar.project, component_def, config):
      raise permissions.PermissionException(
          'User is not allowed to edit this component %s' % component_path)

    original_path = component_def.path
    new_path = component_def.path
    new_docstring = component_def.docstring
    new_deprecated = component_def.deprecated
    new_admin_ids = component_def.admin_ids
    new_cc_ids = component_def.cc_ids
    update_filterrule = False
    for update in request.updates:
      if update.field == api_pb2_v1.ComponentUpdateFieldID.LEAF_NAME:
        leaf_name = update.leafName
        if not tracker_constants.COMPONENT_NAME_RE.match(leaf_name):
          raise exceptions.InvalidComponentNameException(
              'The component name %s is invalid.' % leaf_name)

        if '>' in original_path:
          parent_path = original_path[:original_path.rindex('>')]
          new_path = '%s>%s' % (parent_path, leaf_name)
        else:
          new_path = leaf_name

        conflict = tracker_bizobj.FindComponentDef(new_path, config)
        if conflict and conflict.component_id != component_def.component_id:
          raise exceptions.InvalidComponentNameException(
              'The name %s is already in use.' % new_path)
        update_filterrule = True
      elif update.field == api_pb2_v1.ComponentUpdateFieldID.DESCRIPTION:
        new_docstring = update.description
      elif update.field == api_pb2_v1.ComponentUpdateFieldID.ADMIN:
        user_ids_dict = self._services.user.LookupUserIDs(
            mar.cnxn, list(update.admin), autocreate=True)
        new_admin_ids = list(set(user_ids_dict.values()))
      elif update.field == api_pb2_v1.ComponentUpdateFieldID.CC:
        user_ids_dict = self._services.user.LookupUserIDs(
            mar.cnxn, list(update.cc), autocreate=True)
        new_cc_ids = list(set(user_ids_dict.values()))
        update_filterrule = True
      elif update.field == api_pb2_v1.ComponentUpdateFieldID.DEPRECATED:
        new_deprecated = update.deprecated
      else:
        logging.error('Unknown component field %r', update.field)

    new_modified = int(time.time())
    new_modifier_id = self._services.user.LookupUserID(
        mar.cnxn, mar.auth.email, autocreate=False)
    logging.info(
        'Updating component id %d: path-%s, docstring-%s, deprecated-%s,'
        ' admin_ids-%s, cc_ids-%s modified by %s', component_def.component_id,
        new_path, new_docstring, new_deprecated, new_admin_ids, new_cc_ids,
        new_modifier_id)
    self._services.config.UpdateComponentDef(
        mar.cnxn, mar.project_id, component_def.component_id,
        path=new_path, docstring=new_docstring, deprecated=new_deprecated,
        admin_ids=new_admin_ids, cc_ids=new_cc_ids, modified=new_modified,
        modifier_id=new_modifier_id)

    # TODO(sheyang): reuse the code in componentdetails
    if original_path != new_path:
      # If the name changed then update all of its subcomponents as well.
      subcomponent_ids = tracker_bizobj.FindMatchingComponentIDs(
          original_path, config, exact=False)
      for subcomponent_id in subcomponent_ids:
        if subcomponent_id == component_def.component_id:
          continue
        subcomponent_def = tracker_bizobj.FindComponentDefByID(
            subcomponent_id, config)
        subcomponent_new_path = subcomponent_def.path.replace(
            original_path, new_path, 1)
        self._services.config.UpdateComponentDef(
            mar.cnxn, mar.project_id, subcomponent_def.component_id,
            path=subcomponent_new_path)

    if update_filterrule:
      filterrules_helpers.RecomputeAllDerivedFields(
          mar.cnxn, self._services, mar.project, config)

    return message_types.VoidMessage()


@endpoints.api(name='monorail_client_configs', version='v1',
               description='Monorail API client configs.')
class ClientConfigApi(remote.Service):

  # Class variables. Handy to mock.
  _services = None
  _mar = None

  @classmethod
  def _set_services(cls, services):
    cls._services = services

  def mar_factory(self, request, cnxn):
    if not self._mar:
      self._mar = monorailrequest.MonorailApiRequest(
          request, self._services, cnxn=cnxn)
    return self._mar

  @endpoints.method(
      message_types.VoidMessage,
      message_types.VoidMessage,
      path='client_configs',
      http_method='POST',
      name='client_configs.update')
  def client_configs_update(self, request):
    if self._services is None:
      self._set_services(service_manager.set_up_services())
    mar = self.mar_factory(request, sql.MonorailConnection())
    if not mar.perms.HasPerm(permissions.ADMINISTER_SITE, None, None):
      raise permissions.PermissionException(
          'The requester %s is not allowed to update client configs.' %
           mar.auth.email)

    ROLE_DICT = {
        1: permissions.COMMITTER_ROLE,
        2: permissions.CONTRIBUTOR_ROLE,
    }

    client_config = client_config_svc.GetClientConfigSvc()

    cfg = client_config.GetConfigs()
    if not cfg:
      msg = 'Failed to fetch client configs.'
      logging.error(msg)
      raise endpoints.InternalServerErrorException(msg)

    for client in cfg.clients:
      if not client.client_email:
        continue
      # 1: create the user if non-existent
      user_id = self._services.user.LookupUserID(
          mar.cnxn, client.client_email, autocreate=True)
      user_pb = self._services.user.GetUser(mar.cnxn, user_id)

      logging.info('User ID %d for email %s', user_id, client.client_email)

      # 2: set period and lifetime limit
      # new_soft_limit, new_hard_limit, new_lifetime_limit
      new_limit_tuple = (
          client.period_limit, client.period_limit, client.lifetime_limit)
      action_limit_updates = {'api_request': new_limit_tuple}
      self._services.user.UpdateUserSettings(
          mar.cnxn, user_id, user_pb, action_limit_updates=action_limit_updates)

      logging.info('Updated api request limit %r', new_limit_tuple)

      # 3: Update project role and extra perms
      projects_dict = self._services.project.GetAllProjects(mar.cnxn)
      project_name_to_ids = {
          p.project_name: p.project_id for p in projects_dict.values()}

      # Set project role and extra perms
      for perm in client.project_permissions:
        project_ids = self._GetProjectIDs(perm.project, project_name_to_ids)
        logging.info('Matching projects %r for name %s',
                     project_ids, perm.project)

        role = ROLE_DICT[perm.role]
        for p_id in project_ids:
          project = projects_dict[p_id]
          people_list = []
          if role == 'owner':
            people_list = project.owner_ids
          elif role == 'committer':
            people_list = project.committer_ids
          elif role == 'contributor':
            people_list = project.contributor_ids
          # Onlu update role/extra perms iff changed
          if not user_id in people_list:
            logging.info('Update project %s role %s for user %s',
                         project.project_name, role, client.client_email)
            owner_ids, committer_ids, contributor_ids = (
                project_helpers.MembersWithGivenIDs(project, {user_id}, role))
            self._services.project.UpdateProjectRoles(
                mar.cnxn, p_id, owner_ids, committer_ids,
                contributor_ids)
          if perm.extra_permissions:
            logging.info('Update project %s extra perm %s for user %s',
                         project.project_name, perm.extra_permissions,
                         client.client_email)
            self._services.project.UpdateExtraPerms(
                mar.cnxn, p_id, user_id, list(perm.extra_permissions))

    mar.CleanUp()
    return message_types.VoidMessage()

  def _GetProjectIDs(self, project_str, project_name_to_ids):
    result = []
    if any(ch in project_str for ch in ['*', '+', '?', '.']):
      pattern = re.compile(project_str)
      for p_name in project_name_to_ids.keys():
        if pattern.match(p_name):
          project_id = project_name_to_ids.get(p_name)
          if project_id:
            result.append(project_id)
    else:
      project_id = project_name_to_ids.get(project_str)
      if project_id:
        result.append(project_id)

    if not result:
      logging.warning('Cannot find projects for specified name %s',
                      project_str)
    return result
