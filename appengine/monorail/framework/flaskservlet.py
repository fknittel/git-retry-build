# Copyright 2022 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd
"""Base classes for Monorail Flask servlets.

This is derived from  servlet.py
This base class provides handler methods that conveniently drive
the process of parsing the request, checking base permissions,
gathering common page information, gathering page-specific information,
and adding on-page debugging information (when appropriate).
Subclasses can simply implement the page-specific logic.

Summary of page classes:
  FlaskServlet: abstract base class for all Monorail flask servlets.
"""

import gc
import os
import httplib
import logging
import time
from businesslogic import work_env

import ezt
from features import features_bizobj, hotlist_views
import flask
from proto import project_pb2
from search import query2ast

import settings
from framework import alerts, exceptions, framework_helpers
from framework import framework_views, servlet_helpers
from framework import framework_constants
from framework import monorailrequest
from framework import permissions
from framework import ratelimiter
from framework import template_helpers
from framework import xsrf
from third_party import httpagentparser

from google.appengine.api import app_identity
from google.appengine.api import modules
from google.appengine.api import users
from tracker import tracker_views

NONCE_LENGTH = 32

if not settings.unit_test_mode:
  import MySQLdb


class MethodNotSupportedError(NotImplementedError):
  """An exception class for indicating that the method is not supported.

  Used by GatherPageData and ProcessFormData to indicate that GET and POST,
  respectively, are not supported methods on the given Servlet.
  """
  pass


class FlaskServlet(object):
  """Base class for all Monorail flask servlets.

  Defines a framework of methods that build up parts of the EZT page data.

  Subclasses should override GatherPageData and/or ProcessFormData to
  handle requests.
  """
  _MAIN_TAB_MODE = None  # Normally overridden in subclasses to be one of these:

  # This value should not typically be overridden.
  _TEMPLATE_PATH = framework_constants.TEMPLATE_PATH

  _PAGE_TEMPLATE = None  # Normally overridden in subclasses.
  _ELIMINATE_BLANK_LINES = False

  _MISSING_PERMISSIONS_TEMPLATE = 'sitewide/403-page.ezt'

  def __init__(self, services=None, content_type='text/html; charset=UTF-8'):
    """Load and parse the template, saving it for later use."""
    if self._PAGE_TEMPLATE:  # specified in subclasses
      template_path = self._TEMPLATE_PATH + self._PAGE_TEMPLATE
      self.template = template_helpers.GetTemplate(
          template_path, eliminate_blank_lines=self._ELIMINATE_BLANK_LINES)
    else:
      self.template = None

    self._missing_permissions_template = template_helpers.MonorailTemplate(
        self._TEMPLATE_PATH + self._MISSING_PERMISSIONS_TEMPLATE)
    self.services = services or flask.current_app.config['services']
    self.content_type = content_type
    self.mr = None
    self.request = flask.request
    self.response = None
    self.ratelimiter = ratelimiter.RateLimiter()

  def handler(self):
    """Do common stuff then dispatch the request to get() or put() methods."""
    self.response = flask.make_response()
    handler_start_time = time.time()
    logging.info('\n\n\nRequest handler: %r', self)

    #TODO: add the ts_mon.NonCumulativeDistributionMetric
    # count0, count1, count2 = gc.get_count()
    # logging.info('gc counts: %d %d %d', count0, count1, count2)
    # GC_COUNT.add(count0, {'generation': 0})
    # GC_COUNT.add(count1, {'generation': 1})
    # GC_COUNT.add(count2, {'generation': 2})

    self.mr = monorailrequest.MonorailRequest(self.services)
    self.response.headers.add(
        'Strict-Transport-Security', 'max-age=31536000; includeSubDomains')

    if 'X-Cloud-Trace-Context' in self.request.headers:
      self.mr.profiler.trace_context = (
          self.request.headers.get('X-Cloud-Trace-Context'))

    if self.services.cache_manager:
      try:
        with self.mr.profiler.Phase('distributed invalidation'):
          self.services.cache_manager.DoDistributedInvalidation(self.mr.cnxn)

      except MySQLdb.OperationalError as e:
        logging.exception(e)
        page_data = {
            'http_response_code': httplib.SERVICE_UNAVAILABLE,
            'requested_url': self.request.url,
        }
        self.template = template_helpers.GetTemplate(
            'templates/framework/database-maintenance.ezt',
            eliminate_blank_lines=self._ELIMINATE_BLANK_LINES)
        self.template.WriteFlaskResponse(
            self.response, page_data, content_type='text/html')
        return self.response

    try:
      self.ratelimiter.CheckStart(self.request)

      with self.mr.profiler.Phase('parsing request and doing lookups'):
        self.mr.ParseFlaskRequest(self.request, self.services)

      self.response.headers['X-Frame-Options'] = 'SAMEORIGIN'

      if self.request.method == 'POST':
        self.post()
      elif self.request.method == 'GET':
        self.get()
    except exceptions.NoSuchUserException as e:
      logging.info('Trapped NoSuchUserException %s', e)
      flask.abort(404, 'user not found')

    except exceptions.NoSuchGroupException as e:
      logging.warning('Trapped NoSuchGroupException %s', e)
      flask.abort(404, 'user group not found')

    except exceptions.InputException as e:
      logging.info('Rejecting invalid input: %r', e)
      self.response.status_code = httplib.BAD_REQUEST

    except exceptions.NoSuchProjectException as e:
      logging.info('Rejecting invalid request: %r', e)
      self.response.status_code = httplib.NOT_FOUND

    except xsrf.TokenIncorrect as e:
      logging.info('Bad XSRF token: %r', e.message)
      self.response.status_code = httplib.BAD_REQUEST

    except permissions.BannedUserException as e:
      logging.warning('The user has been banned')
      #TODO: reidrect to userbanned page

    except ratelimiter.RateLimitExceeded as e:
      logging.info('RateLimitExceeded Exception %s', e)
      self.response.status_code = httplib.BAD_REQUEST
      self.response.response = 'Slow your roll.'

    finally:
      self.mr.CleanUp()
      self.ratelimiter.CheckEnd(self.request, time.time(), handler_start_time)

    total_processing_time = time.time() - handler_start_time
    logging.info(
        'Processed request in %d ms', int(total_processing_time * 1000))

    end_count0, end_count1, end_count2 = gc.get_count()
    logging.info('gc counts: %d %d %d', end_count0, end_count1, end_count2)
    # TODO: get the GC event back
    # if (end_count0 < count0) or (end_count1 < count1) or(end_count2 < count2):
    #   GC_EVENT_REQUEST.increment()

    if settings.enable_profiler_logging:
      self.mr.profiler.LogStats()

    return self.response

  def get(self):
    """Collect page-specific and generic info, then render the page.

    Args:
      Any path components parsed by webapp2 will be in kwargs, but we do
        our own parsing later anyway, so ignore them for now.
    """
    page_data = {}
    nonce = framework_helpers.MakeRandomKey(length=NONCE_LENGTH)
    try:
      csp_header = 'Content-Security-Policy'
      csp_scheme = 'https:'
      if settings.local_mode:
        csp_header = 'Content-Security-Policy-Report-Only'
        csp_scheme = 'http:'
      user_agent_str = self.mr.request.headers.get('User-Agent', '')
      ua = httpagentparser.detect(user_agent_str)
      browser, browser_major_version = 'Unknown browser', 0
      if ua.has_key('browser'):
        browser = ua['browser']['name']
        try:
          browser_major_version = int(ua['browser']['version'].split('.')[0])
        except ValueError:
          logging.warn('Could not parse version: %r', ua['browser']['version'])
      csp_supports_report_sample = (
          (browser == 'Chrome' and browser_major_version >= 59) or
          (browser == 'Opera' and browser_major_version >= 46))
      version_base = _VersionBaseURL(self.mr.request)
      self.response.headers.add(
          csp_header,
          (
              "default-src %(scheme)s ; "
              "script-src"
              " %(rep_samp)s"  # Report 40 chars of any inline violation.
              " 'unsafe-inline'"  # Only counts in browsers that lack CSP2.
              " 'strict-dynamic'"  # Allows <script nonce> to load more.
              " %(version_base)s/static/dist/"
              " 'self' 'nonce-%(nonce)s'; "
              "child-src 'none'; "
              "frame-src accounts.google.com"  # All used by gapi.js auth.
              " content-issuetracker.corp.googleapis.com"
              " login.corp.google.com up.corp.googleapis.com"
              # Used by Google Feedback
              " feedback.googleusercontent.com"
              " www.google.com; "
              "img-src %(scheme)s data: blob: ; "
              "style-src %(scheme)s 'unsafe-inline'; "
              "object-src 'none'; "
              "base-uri 'self'; "  # Used by Google Feedback
              "report-uri /csp.do" % {
                  'nonce':
                      nonce,
                  'scheme':
                      csp_scheme,
                  'rep_samp':
                      "'report-sample'" if csp_supports_report_sample else '',
                  'version_base':
                      version_base,
              }))

      # add the function to get data and render page
      page_data.update(self._GatherFlagData(self.mr))

      # # Page-specific work happens in this call.
      page_data.update(self._DoPageProcessing(self.mr, nonce))

      #TODO: (crbug.com/monorail/10863)
      # self._AddHelpDebugPageData(page_data)

      #TODO: (crbug.com/monorail/10868)
      # with self.mr.profiler.Phase('rendering template'):
      #   self._RenderResponse(page_data)

    except (MethodNotSupportedError, NotImplementedError) as e:
      # Instead of these pages throwing 500s display the 404 message and log.
      # The motivation of this is to minimize 500s on the site to keep alerts
      # meaningful during fuzzing. For more context see
      # https://bugs.chromium.org/p/monorail/issues/detail?id=659
      logging.warning('Trapped NotImplementedError %s', e)
      flask.abort(404, 'invalid page')
    except query2ast.InvalidQueryError as e:
      logging.warning('Trapped InvalidQueryError: %s', e)
      logging.exception(e)
      msg = e.message if e.message else 'invalid query'
      flask.abort(400, msg)
    except permissions.PermissionException as e:
      logging.warning('Trapped PermissionException %s', e)
      logging.warning('mr.auth.user_id is %s', self.mr.auth.user_id)
      logging.warning('mr.auth.effective_ids is %s', self.mr.auth.effective_ids)
      logging.warning('mr.perms is %s', self.mr.perms)
      if not self.mr.auth.user_id:
        # If not logged in, let them log in
        url = _SafeCreateLoginURL(self.mr)
        flask.redirect(url, code=307)
      else:
        # Display the missing permissions template.
        page_data = {
            'reason': e.message,
            'http_response_code': httplib.FORBIDDEN,
        }
        with self.mr.profiler.Phase('gather base data'):
          page_data.update(self.GatherBaseData(self.mr, nonce))
        # TODO: (crbug.com/monorail/10863)
        # self._AddHelpDebugPageData(page_data)
        self._missing_permissions_template.WriteFlaskResponse(
            self.response, page_data, content_type=self.content_type)

  def post(self):
    #TODO: implement basic data processing
    logging.info('process post request')

  def _GatherFlagData(self, mr):
    page_data = {
        'project_stars_enabled':
            ezt.boolean(settings.enable_project_stars),
        'user_stars_enabled':
            ezt.boolean(settings.enable_user_stars),
        'can_create_project':
            ezt.boolean(permissions.CanCreateProject(mr.perms)),
        'can_create_group':
            ezt.boolean(permissions.CanCreateGroup(mr.perms)),
    }

    return page_data

  # pylint: disable=unused-argument
  def _DoPageProcessing(self, mr, nonce):
    """Do user lookups and gather page-specific ezt data."""
    with mr.profiler.Phase('common request data'):

      # TODO: (crbug.com/monorail/10861)
      # self._DoCommonRequestProcessing(self.request, mr)

      # TODO: (crbug.com/monorail/10860)
      # self._MaybeRedirectToBrandedDomain(self.request, mr.project_name)

      page_data = self.GatherBaseData(mr, nonce)
      page_data = {}

    with mr.profiler.Phase('page processing'):
      page_data.update(self.GatherPageData(mr))
      page_data.update(mr.form_overrides)
      template_helpers.ExpandLabels(page_data)
      # self._RecordVisitTime(mr)

    return page_data

  # pylint: disable=unused-argument
  def GatherPageData(self, mr):
    """Return a dict of page-specific ezt data."""
    return {}
    # raise MethodNotSupportedError()

  def GatherBaseData(self, mr, nonce):
    """Return a dict of info used on almost all pages."""
    project = mr.project

    project_summary = ''
    project_alert = None
    project_read_only = False
    project_home_page = ''
    project_thumbnail_url = ''
    if project:
      project_summary = project.summary
      # TODO: (crbug.com/monorail/10870)
      # project_alert = _CalcProjectAlert(project)
      project_read_only = project.read_only_reason
      project_home_page = project.home_page
      project_thumbnail_url = tracker_views.LogoView(project).thumbnail_url

    with work_env.WorkEnv(mr, self.services) as we:
      is_project_starred = False
      project_view = None
      if mr.project:
        if permissions.UserCanViewProject(mr.auth.user_pb,
                                          mr.auth.effective_ids, mr.project):
          is_project_starred = we.IsProjectStarred(mr.project_id)
          project_view = template_helpers.PBProxy(mr.project)

    grid_x_attr = None
    grid_y_attr = None
    hotlist_view = None
    if mr.hotlist:
      users_by_id = framework_views.MakeAllUserViews(
          mr.cnxn, self.services.user,
          features_bizobj.UsersInvolvedInHotlists([mr.hotlist]))
      hotlist_view = hotlist_views.HotlistView(
          mr.hotlist, mr.perms, mr.auth, mr.viewed_user_auth.user_id,
          users_by_id,
          self.services.hotlist_star.IsItemStarredBy(
              mr.cnxn, mr.hotlist.hotlist_id, mr.auth.user_id))
      grid_x_attr = mr.x.lower()
      grid_y_attr = mr.y.lower()

    app_version = os.environ.get('CURRENT_VERSION_ID')

    viewed_username = None
    if mr.viewed_user_auth.user_view:
      viewed_username = mr.viewed_user_auth.user_view.username

    config = None
    if mr.project_id and self.services.config:
      with mr.profiler.Phase('getting config'):
        config = self.services.config.GetProjectConfig(mr.cnxn, mr.project_id)
      grid_x_attr = (mr.x or config.default_x_attr).lower()
      grid_y_attr = (mr.y or config.default_y_attr).lower()

    viewing_self = mr.auth.user_id == mr.viewed_user_auth.user_id
    offer_saved_queries_subtab = (
        viewing_self or mr.auth.user_pb and mr.auth.user_pb.is_site_admin)

    # TODO: (crbug.com/monorail/10871)
    # login_url = _SafeCreateLoginURL(mr)
    # logout_url = _SafeCreateLogoutURL(mr)
    login_url = '/'
    logout_url = '/'
    logout_url_goto_home = users.create_logout_url('/')
    version_base = _VersionBaseURL(mr.request)

    base_data = {
        # EZT does not have constants for True and False, so we pass them in.
        'True':
            ezt.boolean(True),
        'False':
            ezt.boolean(False),
        'local_mode':
            ezt.boolean(settings.local_mode),
        'site_name':
            settings.site_name,
        'show_search_metadata':
            ezt.boolean(False),
        'page_template':
            self._PAGE_TEMPLATE,
        'main_tab_mode':
            self._MAIN_TAB_MODE,
        'project_summary':
            project_summary,
        'project_home_page':
            project_home_page,
        'project_thumbnail_url':
            project_thumbnail_url,
        'hotlist_id':
            mr.hotlist_id,
        'hotlist':
            hotlist_view,
        'hostport':
            mr.request.host,
        'absolute_base_url':
            '%s://%s' % (mr.request.scheme, mr.request.host),
        'project_home_url':
            None,
        'link_rel_canonical':
            None,  # For specifying <link rel="canonical">
        'projectname':
            mr.project_name,
        'project':
            project_view,
        'project_is_restricted':
            ezt.boolean(_ProjectIsRestricted(mr)),
        'offer_contributor_list':
            ezt.boolean(permissions.CanViewContributorList(mr, mr.project)),
        'logged_in_user':
            mr.auth.user_view,
        'form_token':
            None,  # Set to a value below iff the user is logged in.
        'form_token_path':
            None,
        'token_expires_sec':
            None,
        'xhr_token':
            None,  # Set to a value below iff the user is logged in.
        'flag_spam_token':
            None,
        'nonce':
            nonce,
        'perms':
            mr.perms,
        'warnings':
            mr.warnings,
        'errors':
            mr.errors,
        'viewed_username':
            viewed_username,
        'viewed_user':
            mr.viewed_user_auth.user_view,
        'viewed_user_pb':
            template_helpers.PBProxy(mr.viewed_user_auth.user_pb),
        'viewing_self':
            ezt.boolean(viewing_self),
        'viewed_user_id':
            mr.viewed_user_auth.user_id,
        'offer_saved_queries_subtab':
            ezt.boolean(offer_saved_queries_subtab),
        'currentPageURL':
            mr.current_page_url,
        'currentPageURLEncoded':
            mr.current_page_url_encoded,
        'login_url':
            login_url,
        'logout_url':
            logout_url,
        'logout_url_goto_home':
            logout_url_goto_home,
        'continue_issue_id':
            mr.continue_issue_id,
        'feedback_email':
            settings.feedback_email,
        'category_css':
            None,  # Used to specify a category of stylesheet
        'category2_css':
            None,  # specify a 2nd category of stylesheet if needed.
        'page_css':
            None,  # Used to add a stylesheet to a specific page.
        'can':
            mr.can,
        'query':
            mr.query,
        'colspec':
            None,
        'sortspec':
            mr.sort_spec,

        # Options for issuelist display
        'grid_x_attr':
            grid_x_attr,
        'grid_y_attr':
            grid_y_attr,
        'grid_cell_mode':
            mr.cells,
        'grid_mode':
            None,
        'list_mode':
            None,
        'chart_mode':
            None,
        'is_cross_project':
            ezt.boolean(False),

        # for project search (some also used in issue search)
        'start':
            mr.start,
        'num':
            mr.num,
        'groupby':
            mr.group_by_spec,
        'q_field_size':
            (
                min(
                    framework_constants.MAX_ARTIFACT_SEARCH_FIELD_SIZE,
                    max(
                        framework_constants.MIN_ARTIFACT_SEARCH_FIELD_SIZE,
                        len(mr.query) + framework_constants.AUTOSIZE_STEP))),
        'mode':
            None,  # Display mode, e.g., grid mode.
        'ajah':
            mr.ajah,
        'table_title':
            mr.table_title,
        'alerts':
            alerts.AlertsView(mr),  # For alert.ezt
        'project_alert':
            project_alert,
        'title':
            None,  # First part of page title
        'title_summary':
            None,  # Appended to title on artifact detail pages
        'project_read_only':
            ezt.boolean(project_read_only),
        'site_read_only':
            ezt.boolean(settings.read_only),
        'banner_time':
            servlet_helpers.GetBannerTime(settings.banner_time),
        'read_only':
            ezt.boolean(settings.read_only or project_read_only),
        'site_banner_message':
            settings.banner_message,
        'robots_no_index':
            None,
        'analytics_id':
            settings.analytics_id,
        'is_project_starred':
            ezt.boolean(is_project_starred),
        'version_base':
            version_base,
        'app_version':
            app_version,
        'gapi_client_id':
            settings.gapi_client_id,
        'viewing_user_page':
            ezt.boolean(False),
        'old_ui_url':
            None,
        'new_ui_url':
            None,
        'is_member':
            ezt.boolean(False),
    }

    if mr.project:
      base_data['project_home_url'] = '/p/%s' % mr.project_name

    # Always add xhr-xsrf token because even anon users need some
    # pRPC methods, e.g., autocomplete, flipper, and charts.
    base_data['token_expires_sec'] = xsrf.TokenExpiresSec()
    base_data['xhr_token'] = xsrf.GenerateToken(
        mr.auth.user_id, xsrf.XHR_SERVLET_PATH)
    # Always add other anti-xsrf tokens when the user is logged in.
    if mr.auth.user_id:
      # TODO: (crbug.com/monorail/10871)
      # form_token_path = self._FormHandlerURL(mr.request.path)
      form_token_path = '/'
      base_data['form_token'] = xsrf.GenerateToken(
          mr.auth.user_id, form_token_path)
      base_data['form_token_path'] = form_token_path

    return base_data


def _ProjectIsRestricted(mr):
  """Return True if the mr has a 'private' project."""
  return (mr.project and mr.project.access != project_pb2.ProjectAccess.ANYONE)


def _VersionBaseURL(request):
  """Return a version-specific URL that we use to load static assets."""
  if settings.local_mode:
    version_base = '%s://%s' % (request.scheme, request.host)
  else:
    version_base = '%s://%s-dot-%s' % (
        request.scheme, modules.get_current_version_name(),
        app_identity.get_default_version_hostname())

  return version_base


def _SafeCreateLoginURL(mr, continue_url=None):
  """Make a login URL w/ a detailed continue URL, otherwise use a short one."""
  continue_url = continue_url or mr.current_page_url
  try:
    url = users.create_login_url(continue_url)
  except users.RedirectTooLongError:
    if mr.project_name:
      url = users.create_login_url('/p/%s' % mr.project_name)
    else:
      url = users.create_login_url('/')

  # Give the user a choice of existing accounts in their session
  # or the option to add an account, even if they are currently
  # signed in to exactly one account.
  if mr.auth.user_id:
    # Notice: this makes assumptions about the output of users.create_login_url,
    # which can change at any time. See https://crbug.com/monorail/3352.
    url = url.replace('/ServiceLogin', '/AccountChooser', 1)
  return url
