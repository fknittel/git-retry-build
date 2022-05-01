# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import contextlib
import copy
import functools
import json
import logging

from google.appengine.ext import ndb
from google.protobuf import json_format
from protorpc import messages
from protorpc import message_types
from protorpc import remote
import endpoints

from components import auth
from components.config import validation as config_validation
from components import protoutil
from components import utils
import gae_ts_mon

from legacy import api_common
from go.chromium.org.luci.buildbucket.proto import builder_common_pb2
from go.chromium.org.luci.buildbucket.proto import builds_service_pb2 as rpc_pb2
from go.chromium.org.luci.buildbucket.proto import common_pb2
from go.chromium.org.luci.buildbucket.proto import project_config_pb2
import bbutil
import buildtags
import config
import creation
import errors
import model
import search
import service
import swarmingcfg
import user
import validation

_PARAM_SWARMING = 'swarming'
_PARAM_CHANGES = 'changes'


class ErrorMessage(messages.Message):
  reason = messages.EnumField(errors.LegacyReason, 1, required=True)
  message = messages.StringField(2, required=True)


def exception_to_error_message(ex):
  assert isinstance(ex, errors.Error)
  assert ex.legacy_reason is not None
  return ErrorMessage(reason=ex.legacy_reason, message=ex.message)


class PubSubCallbackMessage(messages.Message):
  topic = messages.StringField(1, required=True)
  user_data = messages.StringField(2)
  auth_token = messages.StringField(3)


def pubsub_callback_to_notification_config(pubsub_callback, notify):
  """Converts PubSubCallbackMessage to NotificationConfig.

  Ignores auth_token.
  """
  notify.pubsub_topic = pubsub_callback.topic
  notify.user_data = (pubsub_callback.user_data or '').encode('utf-8')


class PutRequestMessage(messages.Message):
  client_operation_id = messages.StringField(1)
  bucket = messages.StringField(2, required=True)
  tags = messages.StringField(3, repeated=True)
  parameters_json = messages.StringField(4)
  lease_expiration_ts = messages.IntegerField(5)
  pubsub_callback = messages.MessageField(PubSubCallbackMessage, 6)
  canary_preference = messages.EnumField(api_common.CanaryPreference, 7)
  experimental = messages.BooleanField(8)


class BuildResponseMessage(messages.Message):
  build = messages.MessageField(api_common.BuildMessage, 1)
  error = messages.MessageField(ErrorMessage, 2)


def parse_v1_tags(v1_tags):
  """Parses V1 tags.

  Returns a tuple of:
    v2_tags: list of StringPair
    gitiles_commit: common_pb2.GitilesCommit or None
    gerrit_changes: list of common_pb2.GerritChange.
  """
  v2_tags = []
  gitiles_commit = None
  gitiles_ref = None
  gerrit_changes = []

  for t in v1_tags:
    key, value = buildtags.parse(t)

    if key == buildtags.GITILES_REF_KEY:
      gitiles_ref = value
      continue

    if key == buildtags.BUILDSET_KEY:
      commit = buildtags.parse_gitiles_commit_buildset(value)
      if commit:
        if gitiles_commit:  # pragma: no cover
          raise errors.InvalidInputError('multiple gitiles commit')
        gitiles_commit = commit
        continue

      cl = buildtags.parse_gerrit_change_buildset(value)
      if cl:
        gerrit_changes.append(cl)
        continue

    v2_tags.append(common_pb2.StringPair(key=key, value=value))

  if gitiles_commit and not gitiles_commit.ref:
    gitiles_commit.ref = gitiles_ref or 'refs/heads/master'

  return v2_tags, gitiles_commit, gerrit_changes


def validate_known_build_parameters(params):
  """Raises errors.InvalidInputError if LUCI build parameters are invalid."""
  params = copy.deepcopy(params)

  ctx = config_validation.Context.raise_on_error(
      exc_type=errors.InvalidInputError
  )

  def bad(fmt, *args):
    raise errors.InvalidInputError(fmt % args)

  def assert_object(name, value):
    if not isinstance(value, dict):
      bad('%s parameter must be an object' % name)

  changes = params.get(_PARAM_CHANGES)
  if changes is not None:
    if not isinstance(changes, list):
      bad('changes param must be an array')
    for c in changes:  # pragma: no branch
      if not isinstance(c, dict):
        bad('changes param must contain only objects')
      repo_url = c.get('repo_url')
      if repo_url is not None and not isinstance(repo_url, basestring):
        bad('change repo_url must be a string')
      author = c.get('author')
      if not isinstance(author, dict):
        bad('change author must be an object')
      email = author.get('email')
      if not isinstance(email, basestring):
        bad('change author email must be a string')
      if not email:
        bad('change author email not specified')

  swarming = params.get(_PARAM_SWARMING)
  if swarming is not None:
    assert_object('swarming', swarming)
    swarming = copy.deepcopy(swarming)

    if swarming:  # pragma: no cover
      bad('unrecognized keys in swarming param: %r', swarming.keys())

  properties = params.get('properties')
  if properties:
    for k, v in sorted(properties.iteritems()):
      with ctx.prefix('property %r:', k):
        swarmingcfg.validate_recipe_property(k, v, ctx)


def put_request_message_to_build_request(put_request, well_known_experiments):
  """Converts PutRequest to BuildRequest.

  Raises errors.InvalidInputError if the put_request is invalid.
  """
  lease_expiration_date = parse_datetime(put_request.lease_expiration_ts)
  errors.validate_lease_expiration_date(lease_expiration_date)

  # Read parameters.
  parameters = parse_json_object(put_request.parameters_json, 'parameters_json')
  parameters = parameters or {}
  validate_known_build_parameters(parameters)

  builder = parameters.pop(api_common.BUILDER_PARAMETER, '') or ''

  # Validate tags.
  buildtags.validate_tags(put_request.tags, 'new', builder=builder)

  # Read properties. Remove them from parameters.
  props = parameters.pop(api_common.PROPERTIES_PARAMETER, None)
  if props is not None and not isinstance(props, dict):
    raise errors.InvalidInputError(
        '"properties" parameter must be a JSON object or null'
    )
  props = props or {}

  changes = parameters.get(_PARAM_CHANGES)
  if changes:  # pragma: no branch
    # Buildbucket-Buildbot integration passes repo_url of the first change in
    # build parameter "changes" as "repository" attribute of SourceStamp.
    # https://chromium.googlesource.com/chromium/tools/build/+/2c6023d
    # /scripts/master/buildbucket/changestore.py#140
    # Buildbot passes repository of the build source stamp as "repository"
    # build property. Recipes, in partiular bot_update recipe module, rely on
    # "repository" property and it is an almost sane property to support in
    # swarmbucket.
    repo_url = changes[0].get('repo_url')
    if repo_url:  # pragma: no branch
      props['repository'] = repo_url

    # Buildbot-Buildbucket integration converts emails in changes to blamelist
    # property.
    emails = [c.get('author', {}).get('email') for c in changes]
    props['blamelist'] = filter(None, emails)

  # Create a v2 request.
  sbr = rpc_pb2.ScheduleBuildRequest(
      builder=builder_common_pb2.BuilderID(builder=builder),
      properties=bbutil.dict_to_struct(props),
      request_id=put_request.client_operation_id,
      experimental=bbutil.BOOLISH_TO_TRINARY[put_request.experimental],
      canary=api_common.CANARY_PREFERENCE_TO_TRINARY.get(
          put_request.canary_preference, common_pb2.UNSET
      ),
  )
  sbr.builder.project, sbr.builder.bucket = config.parse_bucket_id(
      put_request.bucket
  )

  # Parse tags. Extract gitiles commit and gerrit changes.
  tags, gitiles_commit, gerrit_changes = parse_v1_tags(put_request.tags)
  sbr.tags.extend(tags)
  if gitiles_commit:
    sbr.gitiles_commit.CopyFrom(gitiles_commit)

  # Gerrit changes explicitly passed via "gerrit_changes" parameter win.
  gerrit_change_list = parameters.pop('gerrit_changes', None)
  if gerrit_change_list is not None:
    if not isinstance(gerrit_change_list, list):  # pragma: no cover
      raise errors.InvalidInputError('gerrit_changes must be a list')
    try:
      gerrit_changes = [
          json_format.ParseDict(c, common_pb2.GerritChange())
          for c in gerrit_change_list
      ]
    except json_format.ParseError as ex:  # pragma: no cover
      raise errors.InvalidInputError('Invalid gerrit_changes: %s' % ex)

  sbr.gerrit_changes.extend(gerrit_changes)

  if (not gerrit_changes and
      not sbr.builder.bucket.startswith('master.')):  # pragma: no cover
    changes = parameters.get('changes')
    if isinstance(changes, list) and changes and not gitiles_commit:
      legacy_revision = changes[0].get('revision')
      if legacy_revision:
        raise errors.InvalidInputError(
            'legacy revision without gitiles buildset tag'
        )

  # Populate Gerrit project from patch_project property.
  # V2 API users will have to provide this.
  patch_project = props.get('patch_project')
  if len(sbr.gerrit_changes) == 1 and isinstance(patch_project, basestring):
    sbr.gerrit_changes[0].project = patch_project

  # Read PubSub callback.
  pubsub_callback_auth_token = None
  if put_request.pubsub_callback:
    pubsub_callback_auth_token = put_request.pubsub_callback.auth_token
    pubsub_callback_to_notification_config(
        put_request.pubsub_callback, sbr.notify
    )

  # Validate the resulting v2 request before continuing.
  with _wrap_validation_error():
    validation.validate_schedule_build_request(sbr, well_known_experiments)

  return creation.BuildRequest(
      schedule_build_request=sbr,
      parameters=parameters,
      lease_expiration_date=lease_expiration_date,
      pubsub_callback_auth_token=pubsub_callback_auth_token,
  )


def builds_to_messages(builds, include_lease_key=False):
  """Converts model.Build objects to BuildMessage objects.

  Fetches children of the model.Build entities.
  """
  bundle_futs = [
      model.BuildBundle.get_async(
          b, infra=True, input_properties=True, output_properties=True
      ) for b in builds
  ]
  return [
      api_common.build_to_message(
          f.get_result(), include_lease_key=include_lease_key
      ) for f in bundle_futs
  ]


def build_to_message(build, include_lease_key=False):
  return builds_to_messages([build], include_lease_key=include_lease_key)[0]


def build_to_response_message(build, include_lease_key=False):
  msg = build_to_message(build, include_lease_key=include_lease_key)
  return BuildResponseMessage(build=msg)


def id_resource_container(body_message_class=message_types.VoidMessage):
  return endpoints.ResourceContainer(
      body_message_class,
      id=messages.IntegerField(1, required=True),
  )


def catch_errors(fn, response_message_class):

  @functools.wraps(fn)
  def decorated(svc, *args, **kwargs):
    try:
      return fn(svc, *args, **kwargs)
    except errors.Error as ex:
      assert hasattr(response_message_class, 'error')
      return response_message_class(error=exception_to_error_message(ex))
    except auth.AuthorizationError as ex:  # pragma: no cover
      logging.warning(
          'Authorization error.\n%s\nPeer: %s\nIP: %s', ex.message,
          auth.get_peer_identity().to_bytes(), svc.request_state.remote_address
      )
      raise endpoints.ForbiddenException(ex.message)

  return decorated


def convert_bucket(bucket):
  """Converts a bucket string to a bucket id and checks access.

  A synchronous wrapper for api_common.to_bucket_id_async that also checks
  access.

  Raises:
    auth.AuthorizationError if the requester doesn't have access to the bucket.
    errors.InvalidInputError if bucket is invalid or ambiguous.
  """
  bucket_id = api_common.to_bucket_id_async(bucket).get_result()

  # Check access here to return user-supplied bucket name,
  # as opposed to computed bucket id to prevent sniffing bucket names.
  if not bucket_id or not user.has_perm(user.PERM_BUCKETS_GET, bucket_id):
    raise user.current_identity_cannot('access bucket %r', bucket)

  return bucket_id


def convert_bucket_list(buckets):
  # This could be made concurrent, but in practice we search by at most one
  # bucket.
  return map(convert_bucket, buckets)


def buildbucket_api_method(
    request_message_class, response_message_class, **kwargs
):
  """Defines a buildbucket API method."""

  init_auth = auth.endpoints_method(
      request_message_class, response_message_class, **kwargs
  )

  def decorator(fn):
    fn = catch_errors(fn, response_message_class)
    fn = init_auth(fn)

    ts_mon_time = lambda: utils.datetime_to_timestamp(utils.utcnow()) / 1e6
    fn = gae_ts_mon.instrument_endpoint(time_fn=ts_mon_time)(fn)

    # ndb.toplevel must be the last one.
    # We use it because codebase uses the following pattern:
    #   results = [f.get_result() for f in futures]
    # without ndb.Future.wait_all.
    # If a future has an exception, get_result won't be called successive
    # futures, and thus may be left running.
    return ndb.toplevel(fn)

  return decorator


def parse_json_object(json_data, param_name):
  if not json_data:
    return None
  try:
    rv = json.loads(json_data)
  except ValueError as ex:
    raise errors.InvalidInputError('Could not parse %s: %s' % (param_name, ex))
  if rv is not None and not isinstance(rv, dict):
    raise errors.InvalidInputError(
        'Invalid %s: not a JSON object or null' % param_name
    )
  return rv


def parse_datetime(timestamp):
  if timestamp is None:
    return None
  try:
    return utils.timestamp_to_datetime(timestamp)
  except OverflowError:
    raise errors.InvalidInputError('Could not parse timestamp: %s' % timestamp)


def check_scheduling_permissions(bucket_ids):
  """Checks if the requester can schedule builds in any of the buckets.

  Raises auth.AuthorizationError on insufficient permissions.
  """
  bucket_ids = set(bucket_ids)
  can_add = user.filter_buckets_by_perm(user.PERM_BUILDS_ADD, bucket_ids)
  forbidden = sorted(bucket_ids - can_add)
  if forbidden:  # pragma: no cover
    raise user.current_identity_cannot('add builds to buckets %s', forbidden)


@auth.endpoints_api(
    name='buildbucket', version='v1', title='Build Bucket Service'
)
class BuildBucketApi(remote.Service):
  """API for scheduling builds."""

  ####### GET ##################################################################

  @buildbucket_api_method(
      id_resource_container(),
      BuildResponseMessage,
      path='builds/{id}',
      http_method='GET'
  )
  @auth.public
  def get(self, request):
    """Returns a build by id."""
    try:
      build = service.get_async(request.id).get_result()
    except auth.AuthorizationError:
      build = None
    if build is None:
      raise errors.BuildNotFoundError()
    return build_to_response_message(build)

  ####### PUT ##################################################################

  @buildbucket_api_method(
      PutRequestMessage, BuildResponseMessage, path='builds', http_method='PUT'
  )
  @auth.public
  def put(self, request):
    """Creates a new build."""
    request.bucket = convert_bucket(request.bucket)
    check_scheduling_permissions([request.bucket])
    settings = config.get_settings_async().get_result()
    build_req = put_request_message_to_build_request(
        request, set(exp.name for exp in settings.experiment.experiments)
    )
    build = creation.add_async(build_req).get_result()
    return build_to_response_message(build, include_lease_key=True)

  ####### SEARCH ###############################################################

  SEARCH_REQUEST_RESOURCE_CONTAINER = endpoints.ResourceContainer(
      message_types.VoidMessage,
      start_cursor=messages.StringField(1),
      bucket=messages.StringField(2, repeated=True),
      # All specified tags must be present in a build.
      tag=messages.StringField(3, repeated=True),
      status=messages.EnumField(search.StatusFilter, 4),
      result=messages.EnumField(model.BuildResult, 5),
      cancelation_reason=messages.EnumField(model.CancelationReason, 6),
      failure_reason=messages.EnumField(model.FailureReason, 7),
      created_by=messages.StringField(8),
      max_builds=messages.IntegerField(9, variant=messages.Variant.INT32),
      retry_of=messages.IntegerField(10),
      canary=messages.BooleanField(11),
      # search by canary_preference is not supported
      creation_ts_low=messages.IntegerField(12),  # inclusive
      creation_ts_high=messages.IntegerField(13),  # exclusive
      include_experimental=messages.BooleanField(14),
  )

  class SearchResponseMessage(messages.Message):
    builds = messages.MessageField(api_common.BuildMessage, 1, repeated=True)
    next_cursor = messages.StringField(2)
    error = messages.MessageField(ErrorMessage, 3)

  @buildbucket_api_method(
      SEARCH_REQUEST_RESOURCE_CONTAINER,
      SearchResponseMessage,
      path='search',
      http_method='GET'
  )
  @auth.public
  def search(self, request):
    """Searches for builds."""
    assert isinstance(request.tag, list)
    if (request.status is not None or request.result is not None or
        request.failure_reason is not None or
        request.cancelation_reason is not None):  # pragma: no cover
      logging.warning(
          '%s is searching by legacy property: %s',
          auth.get_peer_identity().to_bytes(), request
      )
    builds, next_cursor = search.search_async(
        search.Query(
            bucket_ids=convert_bucket_list(request.bucket),
            tags=request.tag,
            status=request.status,
            result=request.result,
            failure_reason=request.failure_reason,
            cancelation_reason=request.cancelation_reason,
            max_builds=request.max_builds,
            created_by=request.created_by,
            start_cursor=request.start_cursor,
            retry_of=request.retry_of,
            canary=request.canary,
            create_time_low=parse_datetime(request.creation_ts_low),
            create_time_high=parse_datetime(request.creation_ts_high),
            include_experimental=request.include_experimental,
        )
    ).get_result()
    return self.SearchResponseMessage(
        builds=builds_to_messages(builds),
        next_cursor=next_cursor,
    )

  ####### PEEK #################################################################

  PEEK_REQUEST_RESOURCE_CONTAINER = endpoints.ResourceContainer(
      message_types.VoidMessage,
      bucket=messages.StringField(1, repeated=True),
      max_builds=messages.IntegerField(2, variant=messages.Variant.INT32),
      start_cursor=messages.StringField(3),
  )

  @buildbucket_api_method(
      PEEK_REQUEST_RESOURCE_CONTAINER,
      SearchResponseMessage,
      path='peek',
      http_method='GET'
  )
  @auth.public
  def peek(self, request):
    """Returns available builds."""
    assert isinstance(request.bucket, list)
    builds, next_cursor = service.peek(
        convert_bucket_list(request.bucket),
        max_builds=request.max_builds,
        start_cursor=request.start_cursor,
    )
    return self.SearchResponseMessage(
        builds=builds_to_messages(builds),
        next_cursor=next_cursor,
    )

  ####### LEASE ################################################################

  class LeaseRequestBodyMessage(messages.Message):
    lease_expiration_ts = messages.IntegerField(1)

  @buildbucket_api_method(
      id_resource_container(LeaseRequestBodyMessage),
      BuildResponseMessage,
      path='builds/{id}/lease',
      http_method='POST'
  )
  @auth.public
  def lease(self, request):
    """Leases a build.

    Response may contain an error.
    """
    success, build = service.lease(
        request.id,
        lease_expiration_date=parse_datetime(request.lease_expiration_ts),
    )
    if not success:
      return BuildResponseMessage(
          error=ErrorMessage(
              message='Could not lease build',
              reason=errors.LegacyReason.CANNOT_LEASE_BUILD,
          )
      )

    assert build.lease_key is not None
    return build_to_response_message(build, include_lease_key=True)

  ####### START ################################################################

  class StartRequestBodyMessage(messages.Message):
    lease_key = messages.IntegerField(1)
    url = messages.StringField(2)

  @buildbucket_api_method(
      id_resource_container(StartRequestBodyMessage),
      BuildResponseMessage,
      path='builds/{id}/start',
      http_method='POST'
  )
  @auth.public
  def start(self, request):
    """Marks a build as started."""
    build = service.start(request.id, request.lease_key, request.url)
    return build_to_response_message(build)

  ####### HEARTBEAT ############################################################

  class HeartbeatRequestBodyMessage(messages.Message):
    lease_key = messages.IntegerField(1, required=True)
    lease_expiration_ts = messages.IntegerField(2, required=True)

  @buildbucket_api_method(
      id_resource_container(HeartbeatRequestBodyMessage),
      BuildResponseMessage,
      path='builds/{id}/heartbeat',
      http_method='POST'
  )
  @auth.public
  def heartbeat(self, request):
    """Updates build lease."""
    build = service.heartbeat(
        request.id, request.lease_key,
        parse_datetime(request.lease_expiration_ts)
    )
    return build_to_response_message(build)

  class HeartbeatBatchRequestMessage(messages.Message):

    class OneHeartbeat(messages.Message):
      build_id = messages.IntegerField(1, required=True)
      lease_key = messages.IntegerField(2, required=True)
      lease_expiration_ts = messages.IntegerField(3, required=True)

    heartbeats = messages.MessageField(OneHeartbeat, 1, repeated=True)

  class HeartbeatBatchResponseMessage(messages.Message):

    class OneHeartbeatResult(messages.Message):
      build_id = messages.IntegerField(1, required=True)
      lease_expiration_ts = messages.IntegerField(2)
      error = messages.MessageField(ErrorMessage, 3)

    results = messages.MessageField(OneHeartbeatResult, 1, repeated=True)
    error = messages.MessageField(ErrorMessage, 2)

  @buildbucket_api_method(
      HeartbeatBatchRequestMessage,
      HeartbeatBatchResponseMessage,
      path='heartbeat',
      http_method='POST'
  )
  @auth.public
  def heartbeat_batch(self, request):
    """Updates multiple build leases."""
    heartbeats = [{
        'build_id': h.build_id,
        'lease_key': h.lease_key,
        'lease_expiration_date': parse_datetime(h.lease_expiration_ts),
    } for h in request.heartbeats]

    def to_message((build_id, build, ex)):
      msg = self.HeartbeatBatchResponseMessage.OneHeartbeatResult(
          build_id=build_id
      )
      if build:
        msg.lease_expiration_ts = utils.datetime_to_timestamp(
            build.lease_expiration_date
        )
      elif isinstance(ex, errors.Error):
        msg.error = exception_to_error_message(ex)
      else:
        raise ex
      return msg

    results = service.heartbeat_batch(heartbeats)
    return self.HeartbeatBatchResponseMessage(results=map(to_message, results))

  ####### SUCCEED ##############################################################

  class SucceedRequestBodyMessage(messages.Message):
    lease_key = messages.IntegerField(1)
    result_details_json = messages.StringField(2)
    url = messages.StringField(3)
    new_tags = messages.StringField(4, repeated=True)

  @buildbucket_api_method(
      id_resource_container(SucceedRequestBodyMessage),
      BuildResponseMessage,
      path='builds/{id}/succeed',
      http_method='POST'
  )
  @auth.public
  def succeed(self, request):
    """Marks a build as succeeded."""
    build = service.succeed(
        request.id,
        request.lease_key,
        result_details=parse_json_object(
            request.result_details_json, 'result_details_json'
        ),
        url=request.url,
        new_tags=request.new_tags
    )
    return build_to_response_message(build)

  ####### FAIL #################################################################

  class FailRequestBodyMessage(messages.Message):
    lease_key = messages.IntegerField(1)
    result_details_json = messages.StringField(2)
    failure_reason = messages.EnumField(model.FailureReason, 3)
    url = messages.StringField(4)
    new_tags = messages.StringField(5, repeated=True)

  @buildbucket_api_method(
      id_resource_container(FailRequestBodyMessage),
      BuildResponseMessage,
      path='builds/{id}/fail',
      http_method='POST'
  )
  @auth.public
  def fail(self, request):
    """Marks a build as failed."""
    build = service.fail(
        request.id,
        request.lease_key,
        result_details=parse_json_object(
            request.result_details_json, 'result_details_json'
        ),
        failure_reason=request.failure_reason,
        url=request.url,
        new_tags=request.new_tags,
    )
    return build_to_response_message(build)

  ####### CANCEL ###############################################################

  class CancelRequestBodyMessage(messages.Message):
    result_details_json = messages.StringField(1)

  @buildbucket_api_method(
      id_resource_container(CancelRequestBodyMessage),
      BuildResponseMessage,
      path='builds/{id}/cancel',
      http_method='POST'
  )
  @auth.public
  def cancel(self, request):
    """Cancels a build."""
    build = service.cancel_async(
        request.id,
        result_details=parse_json_object(
            request.result_details_json, 'result_details_json'
        ),
    ).get_result()
    return build_to_response_message(build)


@contextlib.contextmanager
def _wrap_validation_error():
  """Converts validation.Error to errors.InvalidInputError."""
  try:
    yield
  except validation.Error as ex:
    raise errors.InvalidInputError(ex.message)
