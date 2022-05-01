# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import inspect
import logging
import re
import urllib

from go.chromium.org.luci.buildbucket.proto import common_pb2
from go.chromium.org.luci.buildbucket.proto.builder_common_pb2 import BuilderID
from google.protobuf import json_format
from google.protobuf.field_mask_pb2 import FieldMask

from common.waterfall import buildbucket_client
from gae_libs.caches import PickledMemCache
from libs.cache_decorator import Cached
from model.wf_build import WfBuild
from services import git
from waterfall.build_info import BuildInfo

# TODO(crbug.com/787676): Use an api rather than parse urls to get the relevant
# data out of a build/tryjob url.
_HOST_NAME_PATTERN = (
    r'https?://(?:(?:build|ci)\.chromium\.org/p|\w+\.\w+\.google\.com/i)')

_MASTER_URL_PATTERN = re.compile(r'^%s/([^/]+)(?:/.*)?$' % _HOST_NAME_PATTERN)

_MILO_MASTER_URL_PATTERN = re.compile(
    r'^https?://luci-milo\.appspot\.com/buildbot/([^/]+)(?:/.*)?$')

_CI_MASTER_URL_PATTERN = re.compile(
    r'^https?://ci\.chromium\.org/buildbot/([^/]+)(?:/.*)?$')

_CI_LONG_MASTER_URL_PATTERN = re.compile(
    r'^https?://ci\.chromium\.org/p/chromium/([^/]+)(?:/.*)?$')

_MASTER_URL_PATTERNS = [  # yapf: disable
    _MASTER_URL_PATTERN,
    _MILO_MASTER_URL_PATTERN,
    _CI_MASTER_URL_PATTERN,
    _CI_LONG_MASTER_URL_PATTERN,
]

_MILO_SWARMING_TASK_URL_PATTERN = re.compile(
    r'^https?://luci-milo\.appspot\.com/swarming/task/([^/]+)(?:/.*)?$')

_BUILD_URL_PATTERN = re.compile(
    r'^%s/([^/]+)/builders/([^/]+)/builds/(\d+)(?:/.*)?$' % _HOST_NAME_PATTERN)

_MILO_BUILD_URL_PATTERN = re.compile(
    r'^https?://luci-milo\.appspot\.com/buildbot/([^/]+)/([^/]+)/(\d+)'
    '(?:/.*)?$')

_MILO_BUILD_LONG_URL_PATTERN = re.compile(
    r'^https?://luci-milo\.appspot\.com/p/([^/]+)/builders/([^/]+)/([^/]+)'
    '/(\d+)(?:/.*)?$')  # pylint: disable=anomalous-backslash-in-string

_CI_BUILD_URL_PATTERN = re.compile(
    r'^https?://ci\.chromium\.org/buildbot/([^/]+)/([^/]+)/(\d+)'
    '(?:/.*)?$')

_CI_BUILD_LONG_URL_PATTERN = re.compile(
    r'^https?://ci\.chromium\.org/p/([^/]+)/builders/([^/]+)/([^/]+)/(\d+)')

_BUILD_URL_PATTERNS = [  # yapf: disable
    _BUILD_URL_PATTERN,
    _MILO_BUILD_URL_PATTERN,
    _CI_BUILD_URL_PATTERN,
]

_STEP_URL_PATTERN = re.compile(
    r'^%s/([^/]+)/builders/([^/]+)/builds/(\d+)/steps/([^/]+)(/.*)?$' %
    _HOST_NAME_PATTERN)

_COMMIT_POSITION_PATTERN = re.compile(r'refs/heads/master@{#(\d+)}$',
                                      re.IGNORECASE)

_BUILDBUCKET_BUILD_URL_PATTERN = re.compile(
    r'^https?://ci\.chromium\.org/b/(\d+)$')


def GetRecentCompletedBuilds(master_name, builder_name, page_size=None):
  """Returns a sorted list of recent completed builds for the given builder.

  Sorted by completed time, newer builds at beginning of the returned list.
  """
  luci_project, luci_bucket = GetLuciProjectAndBucketForMaster(master_name)
  search_builds_response = buildbucket_client.SearchV2BuildsOnBuilder(
      BuilderID(project=luci_project, bucket=luci_bucket, builder=builder_name),
      status=common_pb2.ENDED_MASK,
      page_size=page_size)

  return [build.number for build in search_builds_response.builds
         ] if search_builds_response else None


def GetMasterNameFromUrl(url):
  """Parses the given url and returns the master name."""
  if not url:
    return None

  match = None
  for pattern in _MASTER_URL_PATTERNS:
    match = pattern.match(url)
    if match:
      break
  if not match:
    return None
  return match.group(1)


def _ComputeCacheKeyForLuciBuilder(func, args, kwargs, namespace):
  """Returns a key for the Luci builder passed over to _GetBuildbotMasterName"""
  params = inspect.getcallargs(func, *args, **kwargs)
  return '%s-%s::%s::%s' % (namespace, params['project'], params['bucket_name'],
                            params['builder_name'])


# TODO(crbug/802940): Remove this when the API of getting LUCI build is ready.
@Cached(
    PickledMemCache(),
    namespace='luci-builder-to-master',
    key_generator=_ComputeCacheKeyForLuciBuilder)
def _GetBuildbotMasterName(project, bucket_name, builder_name, build_number):
  """Gets buildbot master name using builder_info and build_number."""
  # TODO(https://crbug.com/1109276) Once builds with the mastername property are
  # beyond horizon that we care about, don't check mastername
  build = buildbucket_client.GetV2BuildByBuilderAndBuildNumber(
      project, bucket_name, builder_name, build_number,
      FieldMask(paths=[
          'input.properties.fields.builder_group',
          'input.properties.fields.mastername'
      ]))
  if not build:
    logging.error('Failed to get mastername for %s::%s::%s::%d', project,
                  bucket_name, builder_name, build_number)
    return None

  return (build.input.properties['builder_group']
          if 'builder_group' in build.input.properties else
          build.input.properties['mastername'])


# TODO(crbug/802940): Remove this when the API of getting LUCI build is ready.
def _ParseBuildLongUrl(url):
  """Parses _CI_BUILD_LONG_URL_PATTERN and _MILO_BUILD_LONG_URL_PATTERN urls."""
  match = None
  for pattern in (_CI_BUILD_LONG_URL_PATTERN, _MILO_BUILD_LONG_URL_PATTERN):
    match = pattern.match(url)
    if match:
      break
  if not match:
    return None

  project, bucket_name, builder_name, build_number = match.groups()
  # If the bucket_name is in the format of 'luci.chromium.ci', only uses 'ci'.
  bucket_name = bucket_name.split('.')[-1]
  builder_name = urllib.unquote(builder_name)

  master_name = _GetBuildbotMasterName(project, bucket_name, builder_name,
                                       int(build_number))
  if not master_name:
    return None

  return master_name, builder_name, int(build_number)


def _ParseBuildbotBuildUrl(url):
  """Parses the given build url using buildbot url format.

  Return:
    (master_name, builder_name, build_number)
  """
  match = None
  for pattern in _BUILD_URL_PATTERNS:
    match = pattern.match(url)
    if match:
      break
  if not match:
    return _ParseBuildLongUrl(url)

  master_name, builder_name, build_number = match.groups()
  builder_name = urllib.unquote(builder_name)
  return master_name, builder_name, int(build_number)


def ParseBuildUrl(url):
  """Parses the given build url.

  Return:
    (master_name, builder_name, build_number)
  """
  if not url:
    return None

  match = _BUILDBUCKET_BUILD_URL_PATTERN.match(url)
  if not match:
    return _ParseBuildbotBuildUrl(url)

  build_id = match.groups()[0]
  build = buildbucket_client.GetV2Build(
      build_id,
      FieldMask(paths=['id', 'number', 'builder', 'input.properties']))
  if not build:
    return None

  input_properties = json_format.MessageToDict(build.input.properties)
  # TODO(https://crbug.com/1109276) Once builds with the mastername property are
  # beyond horizon that we care about, don't check mastername
  master_name = (
      input_properties.get('builder_group') or
      input_properties.get('mastername'))
  if not master_name:
    logging.error('No master name for build %d', build_id)
    return None

  return master_name, build.builder.builder, build.number


def ParseStepUrl(url):
  """Parses the given step url.

  Return:
    (master_name, builder_name, build_number, step_name)
  """
  if not url:
    return None

  match = _STEP_URL_PATTERN.match(url)
  if not match:
    return None

  master_name, builder_name, build_number, step_name, _ = match.groups()
  builder_name = urllib.unquote(builder_name)
  return master_name, builder_name, int(build_number), step_name


def CreateBuildUrl(master_name, builder_name, build_number):
  """Creates the url for the given build."""
  builder_name = urllib.quote(builder_name)
  return 'https://ci.chromium.org/buildbot/%s/%s/%s' % (
      master_name, builder_name, build_number)


def GetBuildIdThroughBuildbucket(master_name, builder_name, build_number):
  """Gets build_id through buildbucket based on legacy buildbot info."""
  luci_project, bucket = GetLuciProjectAndBucketForMaster(master_name)
  build = buildbucket_client.GetV2BuildByBuilderAndBuildNumber(
      luci_project, bucket, builder_name, build_number, FieldMask(paths=['id']))
  if not build:
    logging.error('Failed to get build info for %s/%s/%d', master_name,
                  builder_name, build_number)
    return None
  return build.id


def GetBuildId(master_name, builder_name, build_number):
  """Gets build_id based on legacy buildbot info.

  Returns:
    build_id (str): Id of the build in string format.
  """
  if master_name is None or builder_name is None or build_number is None:
    return None
  build = WfBuild.Get(master_name, builder_name, build_number)
  build_id = build.build_id if build and build.build_id else None
  if not build_id:
    build_id = GetBuildIdThroughBuildbucket(master_name, builder_name,
                                            build_number)

  return str(build_id) if build_id else None


def CreateBuildbucketUrl(master_name, builder_name, build_number):
  """Creates the buildbucket build url for the given build."""
  build_id = GetBuildId(master_name, builder_name, build_number)
  return 'https://ci.chromium.org/b/{}'.format(build_id) if build_id else None


def GetCommitPosition(commit_position_line):
  if commit_position_line:
    match = _COMMIT_POSITION_PATTERN.match(commit_position_line)
    if match:
      return int(match.group(1))
  return None


def GetBlameListForV2Build(build):
  """ Uses gitiles_commit from the previous build and current build to get
     blame_list.

  Args:
    build (build_pb2.Build): All info about the build.

  Returns:
    (list of str): Blame_list of the build.
  """

  search_builds_response = buildbucket_client.SearchV2BuildsOnBuilder(
      build.builder, build_range=(None, build.id), page_size=2)

  previous_build = None
  for search_build in search_builds_response.builds:
    # TODO(crbug.com/969124): remove the loop when SearchBuilds RPC works as
    # expected.
    if search_build.id != build.id:
      previous_build = search_build
      break

  if not previous_build:
    logging.error(
        'No previous build found for build %d, cannot get blame list.',
        build.id)
    return []

  repo_url = git.GetRepoUrlFromV2Build(build)

  return git.GetCommitsBetweenRevisionsInOrder(
      previous_build.input.gitiles_commit.id, build.input.gitiles_commit.id,
      repo_url)


def ExtractBuildInfoFromV2Build(master_name, builder_name, build_number, build):
  """Generates BuildInfo using bb v2 build info.

  This conversion is needed to keep Findit v1 running, will be deprecated in
  v2 (TODO: crbug.com/966982).

  Args:
    master_name (str): The name of the master.
    builder_name (str): The name of the builder.
    build_number (int): The build number.
    build (build_pb2.Build): All info about the build.

  Returns:
    (BuildInfo)
  """
  build_info = BuildInfo(master_name, builder_name, build_number)

  input_properties = json_format.MessageToDict(build.input.properties)

  chromium_revision = build.input.gitiles_commit.id
  runtime = input_properties.get('$recipe_engine/runtime') or {}

  build_info.chromium_revision = chromium_revision

  repo_url = git.GetRepoUrlFromV2Build(build)
  build_info.commit_position = git.GetCommitPositionFromRevision(
      build.input.gitiles_commit.id, repo_url=repo_url)

  build_info.build_start_time = build.create_time.ToDatetime()
  build_info.build_end_time = build.end_time.ToDatetime()
  build_info.completed = bool(build_info.build_end_time)
  build_info.result = build.status
  build_info.parent_buildername = input_properties.get('parent_buildername')
  build_info.parent_mastername = (
      input_properties.get('parent_builder_group') or
      input_properties.get('parent_mastername'))
  build_info.buildbucket_id = str(build.id)
  build_info.buildbucket_bucket = build.builder.bucket
  build_info.is_luci = runtime.get('is_luci')

  build_info.blame_list = GetBlameListForV2Build(build)

  # Step categories:
  # 1. A step is passed if it is in SUCCESS status.
  # 2. A step is failed if it is in FAILURE status.
  # 3. A step is not passed if it is not in SUCCESS status. This category
  #   includes steps in statuses: FAILURE, INFRA_FAILURE, CANCELED, etc.
  for step in build.steps:
    step_name = step.name
    step_status = step.status

    if step_status in [
        common_pb2.STATUS_UNSPECIFIED, common_pb2.SCHEDULED, common_pb2.STARTED
    ]:
      continue

    if step_status != common_pb2.SUCCESS:
      build_info.not_passed_steps.append(step_name)

    if step_name == 'Failure reason':
      # 'Failure reason' is always red when the build breaks or has exception,
      # but it is not a failed step.
      continue

    if not step.logs:
      # Skip wrapping steps.
      continue

    if step_status == common_pb2.SUCCESS:
      build_info.passed_steps.append(step_name)
    elif step_status == common_pb2.FAILURE:
      build_info.failed_steps.append(step_name)

  return build_info


def ValidateBuildUrl(url):
  return bool(
      _MILO_MASTER_URL_PATTERN.match(url) or
      _MILO_SWARMING_TASK_URL_PATTERN.match(url) or
      _BUILD_URL_PATTERN.match(url))


def GetLuciProjectAndBucketForMaster(master_name):
  """Matches master_name to Luci project and bucket.

  This is a temporary and hacky solution to remove Findit's dependencies on
  milo API.

  (TODO: crbug.com/965557): deprecate this when a long solution is in place.

  Args:
    master_name (str): Name of the master for this test.

  Returns:
    (str, str): Luci project and bucket
  """

  bucket = 'ci'
  if master_name.startswith('tryserver'):
    bucket = 'try'

  return 'chromium', bucket
