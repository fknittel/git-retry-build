# Copyright 2015 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import base64
import contextlib
import copy
import datetime
import json
import logging
import uuid

from parameterized import parameterized

from components import utils
utils.fix_protobuf_package()

from google import protobuf
from google.appengine.ext import ndb
from google.protobuf import timestamp_pb2

from components import auth
from components import net
from components import utils
from testing_utils import testing
from webob import exc
import mock
import webapp2

from legacy import api_common
from go.chromium.org.luci.buildbucket.proto import build_pb2
from go.chromium.org.luci.buildbucket.proto import builder_common_pb2
from go.chromium.org.luci.buildbucket.proto import common_pb2
from go.chromium.org.luci.buildbucket.proto import launcher_pb2
from go.chromium.org.luci.buildbucket.proto import project_config_pb2
from go.chromium.org.luci.buildbucket.proto import service_config_pb2
from test import test_util
from test.test_util import future, future_exception
import bbutil
import errors
import experiments
import model
import swarming
import tq
import user

linux_CACHE_NAME = (
    'builder_ccadafffd20293e0378d1f94d214c63a0f8342d1161454ef0acfa0405178106b'
)

NOW = datetime.datetime(2015, 11, 30)


def tspb(seconds, nanos=0):
  return timestamp_pb2.Timestamp(seconds=seconds, nanos=nanos)


class BaseTest(testing.AppengineTestCase):
  maxDiff = None

  def setUp(self):
    super(BaseTest, self).setUp()
    user.clear_request_cache()

    self.patch('tq.enqueue_async', autospec=True, return_value=future(None))

    self.now = NOW
    self.patch(
        'components.utils.utcnow', autospec=True, side_effect=lambda: self.now
    )
    self.patch(
        'resultdb.create_invocations_async',
        autospec=True,
        return_value=future(None),
    )

    self.settings = service_config_pb2.SettingsCfg(
        swarming=dict(
            milo_hostname='milo.example.com',
            bbagent_package=dict(
                package_name='infra/tools/bbagent',
                version='luci-runner-version',
                version_canary='luci-runner-version-canary',
                builders=service_config_pb2.BuilderPredicate(
                    regex=['chromium/try/linux'],
                ),
            ),
            kitchen_package=dict(
                package_name='infra/tools/kitchen',
                version='kitchen-version',
                version_canary='kitchen-version-canary',
            ),
            user_packages=[
                dict(
                    package_name='infra/tools/git',
                    version='git-version',
                    version_canary='git-version-canary',
                ),
                dict(
                    package_name='infra/cpython/python',
                    version='py-version',
                    subdir='python',
                ),
                dict(
                    package_name='infra/excluded',
                    version='excluded-version',
                    builders=service_config_pb2.BuilderPredicate(
                        regex_exclude=['.*'],
                    ),
                ),
            ],
        ),
        known_public_gerrit_hosts=['chromium-review.googlesource.com'],
    )
    self.patch(
        'config.get_settings_async',
        autospec=True,
        return_value=future(self.settings)
    )


class TaskDefTest(BaseTest):

  def setUp(self):
    super(TaskDefTest, self).setUp()

    self.task_template = {
        'name':
            'bb-${build_id}-${project}-${builder}',
        'priority':
            100,
        'tags': ['luci_project:${project}'],
        'task_slices': [{
            'properties': {
                'extra_args': [
                    'cook',
                    '-recipe',
                    '${recipe}',
                    '-properties',
                    '${properties_json}',
                    '-logdog-project',
                    '${project}',
                ],
            },
            'wait_for_capacity': False,
        },],
    }
    self.task_template_canary = self.task_template.copy()
    self.task_template_canary['name'] += '-canary'

    self.patch(
        'google.appengine.api.app_identity.get_default_version_hostname',
        return_value='cr-buildbucket.appspot.com'
    )

  def _test_build(self, **build_proto_fields):
    build = test_util.build(for_creation=True, **build_proto_fields)
    # Ensure "recipe" property is set
    if 'recipe' not in build.proto.input.properties.fields:  # pragma: no branch
      build.proto.input.properties['recipe'] = 'presubmit'
    return build

  def compute_task_def(self, build):
    return swarming.compute_task_def(build, self.settings, fake_build=False)

  def test_shared_cache(self):
    build = self._test_build(
        infra=dict(
            swarming=dict(
                caches=[
                    dict(path='builder', name='shared_builder_cache'),
                ],
            ),
        ),
    )

    slices = self.compute_task_def(build)['task_slices']
    self.assertEqual(
        slices[0]['properties']['caches'], [
            {'path': 'cache/builder', 'name': 'shared_builder_cache'},
        ]
    )

  def test_dimensions_and_cache_fallback(self):
    # Creates 4 task_slices by modifying the buildercfg in 2 ways:
    # - Add two named caches, one expiring at 60 seconds, one at 360 seconds.
    # - Add an optional builder dimension, expiring at 120 seconds.
    #
    # This ensures the combination of these features works correctly, and that
    # multiple 'caches' dimensions can be injected.
    build = self._test_build(
        scheduling_timeout=dict(seconds=3600),
        infra=dict(
            swarming=dict(
                caches=[
                    dict(
                        path='builder',
                        name='shared_builder_cache',
                        wait_for_warm_cache=dict(seconds=60),
                    ),
                    dict(
                        path='second',
                        name='second_cache',
                        wait_for_warm_cache=dict(seconds=360),
                    ),
                ],
                task_dimensions=[
                    dict(key='a', value='1', expiration=dict(seconds=120)),
                    dict(key='a', value='2', expiration=dict(seconds=120)),
                    dict(key='pool', value='Chrome'),
                ]
            )
        )
    )

    slices = self.compute_task_def(build)['task_slices']

    self.assertEqual(4, len(slices))
    for t in slices:
      # They all use the same cache definitions.
      self.assertEqual(
          t['properties']['caches'], [
              {'path': u'cache/builder', 'name': u'shared_builder_cache'},
              {'path': u'cache/second', 'name': u'second_cache'},
          ]
      )

    # But the dimensions are different. 'a' and 'caches' are injected.
    self.assertEqual(
        slices[0]['properties']['dimensions'], [
            {u'key': u'a', u'value': u'1'},
            {u'key': u'a', u'value': u'2'},
            {u'key': u'caches', u'value': u'second_cache'},
            {u'key': u'caches', u'value': u'shared_builder_cache'},
            {u'key': u'pool', u'value': u'Chrome'},
        ]
    )
    self.assertEqual(slices[0]['expiration_secs'], '60')

    # One 'caches' expired. 'a' and one 'caches' are still injected.
    self.assertEqual(
        slices[1]['properties']['dimensions'], [
            {u'key': u'a', u'value': u'1'},
            {u'key': u'a', u'value': u'2'},
            {u'key': u'caches', u'value': u'second_cache'},
            {u'key': u'pool', u'value': u'Chrome'},
        ]
    )
    # 120-60
    self.assertEqual(slices[1]['expiration_secs'], '60')

    # 'a' expired, one 'caches' remains.
    self.assertEqual(
        slices[2]['properties']['dimensions'], [
            {u'key': u'caches', u'value': u'second_cache'},
            {u'key': u'pool', u'value': u'Chrome'},
        ]
    )
    # 360-120
    self.assertEqual(slices[2]['expiration_secs'], '240')

    # The cold fallback; the last 'caches' expired.
    self.assertEqual(
        slices[3]['properties']['dimensions'], [
            {u'key': u'pool', u'value': u'Chrome'},
        ]
    )
    # 3600-360
    self.assertEqual(slices[3]['expiration_secs'], '3240')

  def test_execution_timeout(self):
    build = self._test_build(execution_timeout=dict(seconds=120))
    slices = self.compute_task_def(build)['task_slices']

    self.assertEqual(slices[0]['properties']['execution_timeout_secs'], '120')

  def test_scheduling_timeout(self):
    build = self._test_build(scheduling_timeout=dict(seconds=120))
    slices = self.compute_task_def(build)['task_slices']

    self.assertEqual(1, len(slices))
    self.assertEqual(slices[0]['expiration_secs'], '120')

  def test_compute_bbagent(self):
    build = self._test_build()
    proto = copy.deepcopy(build.proto)
    build.tags_to_protos(proto.tags)
    bbagentargs = swarming._cli_encode_proto(
        launcher_pb2.BBAgentArgs(
            payload_path=swarming._KITCHEN_CHECKOUT,
            cache_dir=swarming._CACHE_DIR,
            known_public_gerrit_hosts=self.settings.known_public_gerrit_hosts,
            build=proto,
        )
    )
    self.assertEqual(
        swarming._compute_bbagent(build, self.settings, fake_build=False),
        [u'bbagent${EXECUTABLE_SUFFIX}', bbagentargs],
    )

  def test_compute_bbagent_get_build(self):
    build = self._test_build(
        infra=dict(buildbucket=dict(hostname='cr-buildbucket.example.com'),)
    )
    build.experiments.append('+%s' % (experiments.BBAGENT_GET_BUILD,))
    self.assertEqual(
        swarming._compute_bbagent(build, self.settings, fake_build=False),
        [
            u'bbagent${EXECUTABLE_SUFFIX}', u'-host',
            u'cr-buildbucket.example.com', u'-build-id', u'9027773186396127232'
        ],
    )

  def test_compute_bbagent_get_build_fake(self):
    build = self._test_build(
        infra=dict(buildbucket=dict(hostname='cr-buildbucket.example.com'),)
    )
    build.experiments.append('+%s' % (experiments.BBAGENT_GET_BUILD,))
    proto = copy.deepcopy(build.proto)
    build.tags_to_protos(proto.tags)
    bbagentargs = swarming._cli_encode_proto(
        launcher_pb2.BBAgentArgs(
            payload_path=swarming._KITCHEN_CHECKOUT,
            cache_dir=swarming._CACHE_DIR,
            known_public_gerrit_hosts=self.settings.known_public_gerrit_hosts,
            build=proto,
        )
    )
    self.assertEqual(
        swarming._compute_bbagent(build, self.settings, fake_build=True),
        [u'bbagent${EXECUTABLE_SUFFIX}', bbagentargs],
    )

  def test_compute_cipd_input_exclusion(self):
    build = self._test_build()
    cipd_input = swarming._compute_cipd_input(build, self.settings)
    packages = {p['package_name']: p for p in cipd_input['packages']}
    self.assertIn('infra/tools/git', packages)
    self.assertIn('infra/cpython/python', packages)
    self.assertNotIn('infra/excluded', packages)

  def test_compute_cipd_input_path(self):
    build = self._test_build()
    cipd_input = swarming._compute_cipd_input(build, self.settings)
    packages = {p['package_name']: p for p in cipd_input['packages']}
    self.assertEqual(
        packages['infra/tools/git']['path'],
        swarming.USER_PACKAGE_DIR,
    )
    self.assertEqual(
        packages['infra/cpython/python']['path'],
        '%s/python' % swarming.USER_PACKAGE_DIR,
    )

  def test_compute_cipd_input_canary(self):
    build = self._test_build(canary=True)
    cipd_input = swarming._compute_cipd_input(build, self.settings)
    packages = {p['package_name']: p for p in cipd_input['packages']}
    self.assertEqual(
        packages['infra/tools/bbagent']['version'],
        'luci-runner-version-canary',
    )
    self.assertEqual(
        packages['infra/tools/git']['version'],
        'git-version-canary',
    )

  def test_compute_cipd_input_bbagent_cipd_handling(self):
    build = self._test_build(
        input=dict(
            experiments=[
                experiments.BBAGENT_DOWNLOAD_CIPD, experiments.USE_BBAGENT
            ],
        ),
    )
    cipd_input = swarming._compute_cipd_input(build, self.settings)
    self.assertEqual(
        test_util.ununicode(cipd_input), {
            'packages': [{
                'package_name': 'infra/tools/bbagent',
                'path': '.',
                'version': 'luci-runner-version',
            }]
        }
    )

  def test_compute_env_prefixes_bbagent_cipd_handling(self):
    build = self._test_build(
        input=dict(
            experiments=[
                experiments.BBAGENT_DOWNLOAD_CIPD, experiments.USE_BBAGENT
            ],
        ),
        infra=dict(
            swarming=dict(
                caches=[
                    dict(
                        path='vpython',
                        name='vpython',
                        env_var='VPYTHON_VIRTUALENV_ROOT'
                    ),
                ],
            ),
        ),
    )
    env_prefixes = swarming._compute_env_prefixes(build, self.settings)
    self.assertEqual(
        test_util.ununicode(env_prefixes),
        [{'key': 'VPYTHON_VIRTUALENV_ROOT', 'value': ['cache/vpython']}]
    )

  def test_properties(self):
    self.patch(
        'components.auth.get_current_identity',
        autospec=True,
        return_value=auth.Identity('user', 'john@example.com')
    )

    now_ts = timestamp_pb2.Timestamp()
    now_ts.FromDatetime(utils.utcnow())
    build = model.Build(
        tags=['t:1'],
        created_by=auth.Anonymous,
        proto=build_pb2.Build(
            id=1,
            builder=dict(project='chromium', bucket='try', builder='linux-rel'),
            number=1,
            status=common_pb2.SCHEDULED,
            created_by='anonymous:anonymous',
            create_time=now_ts,
            update_time=now_ts,
            input=dict(
                properties=bbutil.dict_to_struct({
                    'recipe': 'recipe',
                    'a': 'b',
                }),
                gerrit_changes=[
                    dict(
                        host='chromium-review.googlesource.com',
                        project='chromium/src',
                        change=1234,
                        patchset=5,
                    )
                ],
            ),
            output=dict(),
            infra=dict(
                buildbucket=dict(
                    requested_properties=bbutil.dict_to_struct({'a': 'b'}),
                ),
                recipe=dict(),
            ),
        ),
    )

    actual = bbutil.struct_to_dict(swarming._compute_legacy_properties(build))

    expected = {
        'a': 'b',
        'buildername': 'linux-rel',
        'buildnumber': 1,
        'recipe': 'recipe',
        'repository': 'https://chromium.googlesource.com/chromium/src',
        '$recipe_engine/buildbucket': {
            'hostname': 'cr-buildbucket.appspot.com',
            'build': {
                'id': '1',
                'builder': {
                    'project': 'chromium',
                    'bucket': 'try',
                    'builder': 'linux-rel',
                },
                'number': 1,
                'tags': [{'key': 't', 'value': '1'}],
                'input': {
                    'gerritChanges': [{
                        'host': 'chromium-review.googlesource.com',
                        'project': 'chromium/src',
                        'change': '1234',
                        'patchset': '5',
                    }],
                },
                'infra': {'buildbucket': {},},
                'createdBy': 'anonymous:anonymous',
                'createTime': '2015-11-30T00:00:00Z',
            },
        },
        '$recipe_engine/runtime': {
            'is_experimental': False,
            'is_luci': True,
        },
    }
    self.assertEqual(test_util.ununicode(actual), expected)

  def test_overall(self):
    self.patch(
        'components.auth.get_current_identity',
        autospec=True,
        return_value=auth.Identity('user', 'john@example.com')
    )

    build = self._test_build(
        id=1,
        number=1,
        scheduling_timeout=dict(seconds=3600),
        execution_timeout=dict(seconds=3600),
        grace_period=dict(seconds=45),
        builder=builder_common_pb2.BuilderID(
            project='chromium', bucket='try', builder='linux'
        ),
        exe=dict(
            cipd_package='infra/recipe_bundle',
            cipd_version='refs/heads/master',
        ),
        tags=[
            common_pb2.StringPair(key='custom', value='tag'),
        ],
        input=dict(
            properties=bbutil.dict_to_struct({
                'a': 'b',
                'recipe': 'recipe',
            }),
            gerrit_changes=[
                dict(
                    host='chromium-review.googlesource.com',
                    project='chromium/src',
                    change=1234,
                    patchset=5,
                ),
            ],
        ),
        infra=dict(
            logdog=dict(
                hostname='logs.example.com',
                project='chromium',
                prefix='bb',
            ),
            swarming=dict(
                task_service_account='robot@example.com',
                priority=108,
                task_dimensions=[
                    dict(key='cores', value='8'),
                    dict(key='os', value='Ubuntu'),
                    dict(key='pool', value='Chrome'),
                ],
                caches=[
                    dict(
                        path='vpython',
                        name='vpython',
                        env_var='VPYTHON_VIRTUALENV_ROOT'
                    ),
                ],
            ),
        ),
    )

    actual = self.compute_task_def(build)

    expected_args = launcher_pb2.BBAgentArgs(
        payload_path=swarming._KITCHEN_CHECKOUT,
        cache_dir=swarming._CACHE_DIR,
        known_public_gerrit_hosts=['chromium-review.googlesource.com'],
        build=copy.deepcopy(build.proto),
    )
    # build.proto doesn't have tags in storage.
    build.tags_to_protos(expected_args.build.tags)

    expected_swarming_props_def = {
        'env': [{
            'key': 'BUILDBUCKET_EXPERIMENTAL',
            'value': 'FALSE',
        }],
        'env_prefixes': [
            {
                'key':
                    'PATH',
                'value': [
                    'cipd_bin_packages',
                    'cipd_bin_packages/bin',
                    'cipd_bin_packages/python',
                    'cipd_bin_packages/python/bin',
                ],
            },
            {
                'key': 'VPYTHON_VIRTUALENV_ROOT',
                'value': ['cache/vpython'],
            },
        ],
        'execution_timeout_secs': '3600',
        'grace_period_secs': '225',
        'command': [
            'bbagent${EXECUTABLE_SUFFIX}',
            swarming._cli_encode_proto(expected_args),
        ],
        'dimensions': [
            {'key': 'cores', 'value': '8'},
            {'key': 'os', 'value': 'Ubuntu'},
            {'key': 'pool', 'value': 'Chrome'},
        ],
        'caches': [{'path': 'cache/vpython', 'name': 'vpython'},],
        'cipd_input': {
            'packages': [
                {
                    'package_name': 'infra/tools/bbagent',
                    'path': '.',
                    'version': 'luci-runner-version',
                },
                {
                    'package_name': 'infra/tools/kitchen',
                    'path': '.',
                    'version': 'kitchen-version',
                },
                {
                    'package_name': 'infra/recipe_bundle',
                    'path': 'kitchen-checkout',
                    'version': 'refs/heads/master',
                },
                {
                    'package_name': 'infra/tools/git',
                    'path': swarming.USER_PACKAGE_DIR,
                    'version': 'git-version',
                },
                {
                    'package_name': 'infra/cpython/python',
                    'path': '%s/python' % swarming.USER_PACKAGE_DIR,
                    'version': 'py-version',
                },
            ],
        },
    }
    expected = {
        'name':
            'bb-1-chromium/try/linux-1',
        'realm':
            'chromium:try',
        'priority':
            '108',
        'tags': [
            'build_address:luci.chromium.try/linux/1',
            'buildbucket_bucket:chromium/try',
            'buildbucket_build_id:1',
            'buildbucket_hostname:cr-buildbucket.appspot.com',
            'buildbucket_template_canary:0',
            'builder:linux',
            'buildset:1',
            'custom:tag',
            'luci_project:chromium',
        ],
        'task_slices': [{
            'expiration_secs': '3600',
            'properties': expected_swarming_props_def,
            'wait_for_capacity': False,
        }],
        'pubsub_topic':
            'projects/testbed-test/topics/swarming',
        'pubsub_userdata':
            json.dumps({
                'created_ts': 1448841600000000,
                'swarming_hostname': 'swarming.example.com',
                'build_id': 1L,
            },
                       sort_keys=True),
        'service_account':
            'robot@example.com',
    }
    self.assertEqual(test_util.ununicode(actual), expected)

    self.assertEqual(
        build.proto.infra.swarming.task_service_account, 'robot@example.com'
    )

    # Now check that the blob on the cli is actually reasonable.
    cli_blob = actual['task_slices'][0]['properties']['command'][1]
    # No newlines, no padding
    self.assertNotIn('\n', cli_blob)
    self.assertNotIn('=', cli_blob)
    # Restore padding, so python can decode it.
    #   l % 4 == 0 -> no padding       (ex "aGkh")
    #   l % 4 == 1 -> =                (ex "aGk=")
    #   l % 4 == 2 -> ==               (ex "aA==")
    #   l % 4 == 3 -> <invalid state>  (cannot happen in well-formed base64)
    remainder = len(cli_blob) % 4
    if remainder:  # pragma: no cover
      padding = '=' * (4 - remainder)
      self.assertLessEqual(len(padding), 2)  # should be '', '=', or '=='
    else:  # pragma: no cover
      padding = ''

    args = launcher_pb2.BBAgentArgs()
    args.ParseFromString((cli_blob + padding).decode('base64').decode('zlib'))
    self.assertEqual(args, expected_args)

    self.assertNotIn('buildbucket', build.proto.input.properties)
    self.assertNotIn('$recipe_engine/buildbucket', build.proto.input.properties)

  def test_legacy_kitchen(self):
    build = self._test_build(
        builder=builder_common_pb2.BuilderID(
            project='chromium', bucket='try', builder='linux_kitchen'
        ),
    )
    build.proto.exe.cmd[0] = 'recipes'
    actual = self.compute_task_def(build)

    self.assertEqual([
        "kitchen${EXECUTABLE_SUFFIX}", 'cook', '-buildbucket-hostname',
        'cr-buildbucket.appspot.com', '-buildbucket-build-id',
        '9027773186396127232', '-call-update-build', '-build-url',
        'https://milo.example.com/b/9027773186396127232',
        '-luci-system-account', 'system', '-recipe', 'presubmit', '-cache-dir',
        'cache', '-checkout-dir', 'kitchen-checkout', '-temp-dir', 'tmp',
        '-properties',
        api_common.properties_to_json(
            swarming._compute_legacy_properties(build)
        ), '-logdog-annotation-url',
        'logdog://logdog.example.com/chromium/bb/+/annotations',
        '-known-gerrit-host', 'chromium-review.googlesource.com'
    ], test_util.ununicode(actual['task_slices'][0]['properties']['command']))

  def test_experimental(self):
    build = self._test_build(input=dict(experimental=True))
    actual = self.compute_task_def(build)

    env = actual['task_slices'][0]['properties']['env']
    self.assertIn({
        'key': 'BUILDBUCKET_EXPERIMENTAL',
        'value': 'TRUE',
    }, env)

  def test_parent_run_id(self):
    build = self._test_build(
        infra=dict(swarming=dict(parent_run_id='deadbeef'))
    )
    actual = self.compute_task_def(build)
    self.assertEqual(actual['parent_task_id'], 'deadbeef')

  def test_parent(self):
    build = self._test_build(
        infra=dict(swarming=dict(parent_run_id='deadbeef')),
        ancestor_ids=[123],
        input=dict(experiments=['luci.buildbucket.parent_tracking']),
    )
    actual = self.compute_task_def(build)
    self.assertIsNone(actual.get('parent_task_id'))

  def test_parent_no_exp(self):
    # Even though the build has a parent,
    build = self._test_build(
        infra=dict(swarming=dict(parent_run_id='deadbeef')),
        ancestor_ids=[123],
    )
    actual = self.compute_task_def(build)
    self.assertEqual(actual['parent_task_id'], 'deadbeef')

  def test_generate_build_url(self):
    build = self._test_build(id=1)
    self.assertEqual(
        swarming._generate_build_url('milo.example.com', build),
        'https://milo.example.com/b/1',
    )

    self.assertEqual(
        swarming._generate_build_url(None, build),
        ('https://swarming.example.com/task?id=deadbeef')
    )


class SyncBuildTest(BaseTest):

  def setUp(self):
    super(SyncBuildTest, self).setUp()
    self.patch('components.net.json_request_async', autospec=True)

    self.build_token = 'beeff00d'
    self.gen_build_token_mock = self.patch(
        'tokens.generate_build_token',
        autospec=True,
        return_value=self.build_token,
    )

    self.task_def = {'is_task_def': True, 'task_slices': [{
        'properties': {},
    }]}
    self.patch(
        'swarming.compute_task_def', autospec=True, return_value=self.task_def
    )
    self.patch(
        'google.appengine.api.app_identity.get_default_version_hostname',
        return_value='cr-buildbucket.appspot.com'
    )

    self.build_bundle = test_util.build_bundle(
        id=1, created_by='user:john@example.com'
    )
    self.build.resultdb_update_token = self.resultdb_update_token = 'abc01d2e'

    self.build_bundle.build.swarming_task_key = None
    with self.build_bundle.infra.mutate() as infra:
      infra.swarming.task_id = ''
    self.build_bundle.put()

  @property
  def build(self):
    return self.build_bundle.build

  def _create_task(self):
    self.build_bundle.build.proto.infra.ParseFromString(
        self.build_bundle.infra.infra
    )
    self.build_bundle.build.proto.input.properties.ParseFromString(
        self.build_bundle.input_properties.properties
    )
    swarming._create_swarming_task(self.build_bundle.build)

  def _expected_task_def(self, build_token=None, **extra):
    expected_task_def = copy.deepcopy(self.task_def)
    expected_secrets = launcher_pb2.BuildSecrets(
        build_token=(build_token or self.build_token),
        resultdb_invocation_update_token=self.resultdb_update_token,
    )
    expected_task_def[u'task_slices'][0][u'properties'][u'secret_bytes'] = (
        base64.b64encode(expected_secrets.SerializeToString())
    )
    expected_task_def['request_uuid'] = str(uuid.UUID(int=self.build.proto.id))
    expected_task_def.update(**extra)
    return expected_task_def

  def test_build_entity_generates_update_token(self):
    net.json_request_async.return_value = future({'task_id': 'x'})
    self.assertFalse(self.build_bundle.build.update_token)
    self.build_bundle.put()
    swarming._sync_build(self.build.proto.id, 0)
    self.gen_build_token_mock.assert_called_once_with(1)

    net.json_request_async.assert_called_with(
        'https://swarming.example.com/_ah/api/swarming/v1/tasks/new',
        method='POST',
        scopes=net.EMAIL_SCOPE,
        payload=self._expected_task_def(),
        project_id=test_util.BUILD_DEFAULTS.builder.project,
        deadline=30,
        max_attempts=1,
    )

    # Assert that the build entity has the new token.
    bundle = model.BuildBundle.get(1)
    self.assertEqual(bundle.build.update_token, self.build_token)

  def test_create_task(self):
    net.json_request_async.return_value = future({'task_id': 'x'})
    swarming._sync_build(self.build.proto.id, 0)

    net.json_request_async.assert_called_with(
        'https://swarming.example.com/_ah/api/swarming/v1/tasks/new',
        method='POST',
        scopes=net.EMAIL_SCOPE,
        payload=self._expected_task_def(),
        project_id=test_util.BUILD_DEFAULTS.builder.project,
        deadline=30,
        max_attempts=1,
    )

    # Assert that we've persisted information about the new task.
    bundle = model.BuildBundle.get(1, infra=True)
    self.assertIsNotNone(bundle)
    self.assertTrue(bundle.infra.parse().swarming.task_id)

    expected_continuation_payload = {
        'id': 1,
        'generation': 1,
    }
    expected_continuation = {
        'name': 'sync-task-1-1',
        'url': '/internal/task/swarming/sync-build/1',
        'payload': json.dumps(expected_continuation_payload, sort_keys=True),
        'retry_options': {
            'task_age_limit': model.BUILD_TIMEOUT.total_seconds()
        },
        'countdown': 300,
    }
    tq.enqueue_async.assert_called_with(
        swarming.SYNC_QUEUE_NAME, [expected_continuation], transactional=False
    )

  @mock.patch('swarming.cancel_task', autospec=True)
  def test_already_exists_after_creation(self, cancel_task):

    @ndb.tasklet
    def json_request_async(*_args, **_kwargs):
      with self.build_bundle.infra.mutate() as infra:
        infra.swarming.task_id = 'deadbeef'
      yield self.build_bundle.infra.put_async()

      raise ndb.Return({'task_id': 'new task'})

    net.json_request_async.side_effect = json_request_async

    self._create_task()
    cancel_task.assert_called_with(
        'swarming.example.com', 'new task', 'chromium:try'
    )

  def test_http_400(self):
    net.json_request_async.return_value = future_exception(
        net.Error('HTTP 401', 400, 'invalid request')
    )

    self._create_task()

    build = self.build.key.get()
    self.assertEqual(build.status, common_pb2.INFRA_FAILURE)
    self.assertEqual(
        build.proto.summary_markdown,
        r'Swarming task creation API responded with HTTP 400: `invalid request`'
    )

  def test_http_500(self):
    net.json_request_async.return_value = future_exception(
        net.Error('internal', 500, 'Internal server error')
    )

    with self.assertRaises(net.Error):
      self._create_task()

  def test_http_500_give_up(self):
    net.json_request_async.return_value = future_exception(
        net.Error('internal', 500, 'Internal server error')
    )

    self.now += swarming._SWARMING_CREATE_TASK_GIVE_UP_TIMEOUT
    self.now += datetime.timedelta(seconds=1)

    self._create_task()

    build = self.build.key.get()
    self.assertEqual(build.status, common_pb2.INFRA_FAILURE)
    self.assertEqual(
        build.proto.summary_markdown,
        'Swarming task creation API responded with HTTP 500 after '
        'several attempts: `Internal server error`'
    )

  def test_http_timeout_give_up(self):
    net.json_request_async.return_value = future_exception(
        net.Error(None, None, None)
    )

    self.now += swarming._SWARMING_CREATE_TASK_GIVE_UP_TIMEOUT
    self.now += datetime.timedelta(seconds=1)

    self._create_task()

    build = self.build.key.get()
    self.assertEqual(build.status, common_pb2.INFRA_FAILURE)
    self.assertEqual(
        build.proto.summary_markdown,
        'Swarming task creation API timed-out after several '
        'attempts. (timeout=30 sec)',
    )

  def test_validate(self):
    build = test_util.build()
    swarming.validate_build(build)

  def test_validate_lease_key(self):
    build = test_util.build()
    build.lease_key = 123
    with self.assertRaises(errors.InvalidInputError):
      swarming.validate_build(build)

  @parameterized.expand([
      (
          dict(
              infra=dict(
                  swarming=dict(
                      task_dimensions=[
                          dict(
                              key='a',
                              value='b',
                              expiration=dict(seconds=60 * i)
                          ) for i in xrange(7)
                      ],
                  ),
              ),
          ),
      ),
  ])
  def test_validate_fails(self, build_params):
    build = test_util.build(for_creation=True, **build_params)
    with self.assertRaises(errors.InvalidInputError):
      swarming.validate_build(build)

  @parameterized.expand([
      ({
          'task_result': None,
          'status': common_pb2.INFRA_FAILURE,
          'end_time': test_util.dt2ts(NOW),
      },),
      ({
          'task_result': {'state': 'PENDING'},
          'status': common_pb2.SCHEDULED,
      },),
      ({
          'task_result': {
              'state': 'RUNNING',
              'started_ts': '2018-01-29T21:15:02.649750',
          },
          'status': common_pb2.STARTED,
          'start_time': tspb(seconds=1517260502, nanos=649750000),
      },),
      ({
          'task_result': {
              'state': 'COMPLETED',
              'started_ts': '2018-01-29T21:15:02.649750',
              'completed_ts': '2018-01-30T00:15:18.162860',
          },
          'status': common_pb2.SUCCESS,
          'start_time': tspb(seconds=1517260502, nanos=649750000),
          'end_time': tspb(seconds=1517271318, nanos=162860000),
      },),
      ({
          'task_result': {
              'state':
                  'COMPLETED',
              'started_ts':
                  '2018-01-29T21:15:02.649750',
              'completed_ts':
                  '2018-01-30T00:15:18.162860',
              'bot_dimensions': [
                  {'key': 'os', 'value': ['Ubuntu', 'Trusty']},
                  {'key': 'pool', 'value': ['luci.chromium.try']},
                  {'key': 'id', 'value': ['bot1']},
                  {'key': 'empty'},
              ],
          },
          'status': common_pb2.SUCCESS,
          'bot_dimensions': [
              common_pb2.StringPair(key='id', value='bot1'),
              common_pb2.StringPair(key='os', value='Trusty'),
              common_pb2.StringPair(key='os', value='Ubuntu'),
              common_pb2.StringPair(key='pool', value='luci.chromium.try'),
          ],
          'start_time': tspb(seconds=1517260502, nanos=649750000),
          'end_time': tspb(seconds=1517271318, nanos=162860000),
      },),
      ({
          'task_result': {
              'state': 'COMPLETED',
              'failure': True,
              'started_ts': '2018-01-29T21:15:02.649750',
              'completed_ts': '2018-01-30T00:15:18.162860',
          },
          'status': common_pb2.INFRA_FAILURE,
          'start_time': tspb(seconds=1517260502, nanos=649750000),
          'end_time': tspb(seconds=1517271318, nanos=162860000),
      },),
      ({
          'task_result': {
              'state': 'BOT_DIED',
              'started_ts': '2018-01-29T21:15:02.649750',
              'abandoned_ts': '2018-01-30T00:15:18.162860',
          },
          'status': common_pb2.INFRA_FAILURE,
          'start_time': tspb(seconds=1517260502, nanos=649750000),
          'end_time': tspb(seconds=1517271318, nanos=162860000),
      },),
      ({
          'task_result': {
              'state': 'TIMED_OUT',
              'started_ts': '2018-01-29T21:15:02.649750',
              'completed_ts': '2018-01-30T00:15:18.162860',
          },
          'status': common_pb2.INFRA_FAILURE,
          'is_timeout': True,
          'start_time': tspb(seconds=1517260502, nanos=649750000),
          'end_time': tspb(seconds=1517271318, nanos=162860000),
      },),
      ({
          'task_result': {
              'state': 'EXPIRED',
              'abandoned_ts': '2018-01-30T00:15:18.162860',
          },
          'status': common_pb2.INFRA_FAILURE,
          'is_resource_exhaustion': True,
          'is_timeout': True,
          'end_time': tspb(seconds=1517271318, nanos=162860000),
      },),
      ({
          'task_result': {
              'state': 'KILLED',
              'abandoned_ts': '2018-01-30T00:15:18.162860',
          },
          'status': common_pb2.CANCELED,
          'end_time': tspb(seconds=1517271318, nanos=162860000),
      },),
      ({
          'task_result': {
              'state': 'CANCELED',
              'abandoned_ts': '2018-01-30T00:15:18.162860',
          },
          'status': common_pb2.CANCELED,
          'end_time': tspb(seconds=1517271318, nanos=162860000),
      },),
      ({
          'task_result': {
              'state': 'NO_RESOURCE',
              'abandoned_ts': '2018-01-30T00:15:18.162860',
          },
          'status': common_pb2.INFRA_FAILURE,
          'is_resource_exhaustion': True,
          'end_time': tspb(seconds=1517271318, nanos=162860000),
      },),
      # NO_RESOURCE with abandoned_ts before creation time.
      (
          {
              'task_result': {
                  'state': 'NO_RESOURCE',
                  'abandoned_ts': '2015-11-29T00:15:18.162860',
              },
              'status': common_pb2.INFRA_FAILURE,
              'is_resource_exhaustion': True,
              'end_time': test_util.dt2ts(NOW),
          },
      ),
  ])
  def test_sync_with_task_result(self, case):
    logging.info('test case: %s', case)
    bundle = test_util.build_bundle(id=1)
    bundle.put()

    net.json_request_async.return_value = future(case['task_result'])

    swarming._sync_build(1, 1)

    net.json_request_async.assert_called_with(
        (
            'https://swarming.example.com/'
            '_ah/api/swarming/v1/task/deadbeef/result'
        ),
        method='GET',
        scopes=net.EMAIL_SCOPE,
        payload=None,
        project_id=test_util.BUILD_DEFAULTS.builder.project,
        deadline=None,
        max_attempts=None,
    )

    build = bundle.build.key.get()
    build_infra = bundle.infra.key.get()
    bp = build.proto
    self.assertEqual(bp.status, case['status'])
    self.assertEqual(
        bp.status_details.HasField('timeout'),
        case.get('is_timeout', False),
    )
    self.assertEqual(
        bp.status_details.HasField('resource_exhaustion'),
        case.get('is_resource_exhaustion', False)
    )

    self.assertEqual(bp.start_time, case.get('start_time', tspb(0)))
    self.assertEqual(bp.end_time, case.get('end_time', tspb(0)))

    self.assertEqual(
        list(build_infra.parse().swarming.bot_dimensions),
        case.get('bot_dimensions', [])
    )

    expected_continuation_payload = {
        'id': 1,
        'generation': 2,
    }
    expected_continuation = {
        'name': 'sync-task-1-2',
        'url': '/internal/task/swarming/sync-build/1',
        'payload': json.dumps(expected_continuation_payload, sort_keys=True),
        'retry_options': {
            'task_age_limit': model.BUILD_TIMEOUT.total_seconds()
        },
        'countdown': 300,
    }
    tq.enqueue_async.assert_called_with(
        swarming.SYNC_QUEUE_NAME, [expected_continuation], transactional=False
    )

  def test_termination(self):
    self.build.proto.status = common_pb2.SUCCESS
    self.build.proto.end_time.FromDatetime(utils.utcnow())
    self.build.put()

    swarming._sync_build(1, 1)
    self.assertFalse(tq.enqueue_async.called)


class CancelTest(BaseTest):

  def setUp(self):
    super(CancelTest, self).setUp()

    self.json_response = None

    def json_request_async(*_, **__):
      if self.json_response is not None:
        return future(self.json_response)
      self.fail('unexpected outbound request')  # pragma: no cover

    self.patch(
        'components.net.json_request_async',
        autospec=True,
        side_effect=json_request_async
    )

  def test_cancel_task(self):
    self.json_response = {'ok': True}
    swarming.cancel_task('swarming.example.com', 'deadbeef', 'chromium:try')
    net.json_request_async.assert_called_with(
        (
            'https://swarming.example.com/'
            '_ah/api/swarming/v1/task/deadbeef/cancel'
        ),
        method='POST',
        scopes=net.EMAIL_SCOPE,
        payload={'kill_running': True},
        project_id=test_util.BUILD_DEFAULTS.builder.project,
        deadline=None,
        max_attempts=None,
    )

  def test_cancel_running_task(self):
    self.json_response = {
        'was_running': True,
        'ok': False,
    }
    swarming.cancel_task('swarming.example.com', 'deadbeef', 'chromium:try')


class CancelTQTaskTest(BaseTest):

  @mock.patch('tq.enqueue_async')
  def test_tq_task(self, enqueue_async):
    enqueue_async.return_value = future(None)

    build = test_util.build(id=1, for_creation=True)

    @ndb.transactional
    def txn():
      swarming.cancel_task_transactionally_async(
          build, build.proto.infra.swarming
      ).get_result()

    txn()

    enqueue_async.assert_called_with(
        'backend-default', [{
            'url':
                '/internal/task/buildbucket/cancel_swarming_task/' +
                'swarming.example.com/deadbeef',
            'payload': {
                'hostname': 'swarming.example.com',
                'task_id': 'deadbeef',
                'realm': 'chromium:try',
            },
        }]
    )


class SubNotifyTest(BaseTest):

  def setUp(self):
    super(SubNotifyTest, self).setUp()
    self.handler = swarming.SubNotify(response=webapp2.Response())
    self.build_bundle = test_util.build_bundle(id=1)

  def test_unpack_msg(self):
    self.assertEqual(
        self.handler.unpack_msg({
            'messageId':
                '1', 'data':
                    b64json({
                        'task_id':
                            'deadbeef', 'userdata':
                                json.dumps({
                                    'created_ts': 1448841600000000,
                                    'swarming_hostname': 'swarming.example.com',
                                    'build_id': 1L,
                                })
                    })
        }), (
            'swarming.example.com', datetime.datetime(2015, 11,
                                                      30), 'deadbeef', 1
        )
    )

  def test_unpack_msg_with_err(self):
    with self.assert_bad_message():
      self.handler.unpack_msg({})
    with self.assert_bad_message():
      self.handler.unpack_msg({'data': b64json([])})

    bad_data = [
        # Bad task id.
        {
            'userdata':
                json.dumps({
                    'created_ts': 1448841600000,
                    'swarming_hostname': 'swarming.example.com',
                })
        },

        # Bad swarming hostname.
        {
            'task_id': 'deadbeef',
        },
        {
            'task_id': 'deadbeef',
            'userdata': '{}',
        },
        {
            'task_id': 'deadbeef',
            'userdata': json.dumps({
                'swarming_hostname': 1,
            }),
        },

        # Bad creation time
        {
            'task_id':
                'deadbeef',
            'userdata':
                json.dumps({
                    'swarming_hostname': 'swarming.example.com',
                }),
        },
        {
            'task_id':
                'deadbeef',
            'userdata':
                json.dumps({
                    'created_ts': 'foo',
                    'swarming_hostname': 'swarming.example.com',
                }),
        },
    ]

    for data in bad_data:
      with self.assert_bad_message():
        self.handler.unpack_msg({'data': b64json(data)})

  def mock_request(self, user_data, task_id='deadbeef'):
    msg_data = b64json({
        'task_id': task_id,
        'userdata': json.dumps(user_data),
    })
    self.handler.request = mock.Mock(
        json={
            'message': {
                'messageId': '1',
                'data': msg_data,
            },
        }
    )

  @mock.patch('swarming._load_task_result', autospec=True)
  def test_post(self, load_task_result):
    self.build_bundle.put()
    self.mock_request({
        'build_id': 1L,
        'created_ts': 1448841600000000,
        'swarming_hostname': 'swarming.example.com',
    })

    load_task_result.return_value = {
        'task_id': 'deadbeef',
        'state': 'COMPLETED',
    }

    self.handler.post()

    build = self.build_bundle.build.key.get()
    self.assertEqual(build.proto.status, common_pb2.SUCCESS)

  def test_post_with_different_swarming_hostname(self):
    self.build_bundle.put()

    self.mock_request({
        'build_id': 1L,
        'created_ts': 1448841600000000,
        'swarming_hostname': 'different-chromium.example.com',
    })

    with self.assert_bad_message(expect_redelivery=False):
      self.handler.post()

  def test_post_with_different_task_id(self):
    self.build_bundle.put()

    self.mock_request(
        {
            'build_id': 1L,
            'created_ts': 1448841600000000,
            'swarming_hostname': 'swarming.example.com',
        },
        task_id='deadbeefffffffffff',
    )

    with self.assert_bad_message(expect_redelivery=False):
      self.handler.post()

  def test_post_too_soon(self):
    with self.build_bundle.infra.mutate() as infra:
      infra.swarming.task_id = ''
    self.build_bundle.put()

    self.mock_request({
        'build_id': 1L,
        'created_ts': 1448841600000000,
        'swarming_hostname': 'swarming.example.com',
    })

    with self.assert_bad_message(expect_redelivery=True):
      self.handler.post()

  def test_post_without_task_id(self):
    self.mock_request(
        {
            'build_id': 1L,
            'created_ts': 1448841600000000,
            'swarming_hostname': 'swarming.example.com',
        },
        task_id=None,
    )

    with self.assert_bad_message(expect_redelivery=False):
      self.handler.post()

  def test_post_without_build_id(self):
    self.mock_request({
        'created_ts': 1448841600000000,
        'swarming_hostname': 'swarming.example.com',
    })
    with self.assert_bad_message(expect_redelivery=False):
      self.handler.post()

  def test_post_without_build(self):
    self.mock_request({
        'created_ts': 1438841600000000,
        'swarming_hostname': 'swarming.example.com',
        'build_id': 1L,
    })

    with self.assert_bad_message(expect_redelivery=False):
      self.handler.post()

  @contextlib.contextmanager
  def assert_bad_message(self, expect_redelivery=False):
    self.handler.bad_message = False
    err = exc.HTTPClientError if expect_redelivery else exc.HTTPOk
    with self.assertRaises(err):
      yield
    self.assertTrue(self.handler.bad_message)

  @mock.patch('swarming.SubNotify._process_msg', autospec=True)
  def test_dedup_messages(self, _process_msg):
    self.handler.request = mock.Mock(
        json={'message': {
            'messageId': '1',
            'data': b64json({}),
        }}
    )

    self.handler.post()
    self.handler.post()

    self.assertEquals(_process_msg.call_count, 1)


def b64json(data):
  return base64.b64encode(json.dumps(data))
