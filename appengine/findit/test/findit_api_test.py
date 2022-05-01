# Copyright 2015 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import datetime
import json
import logging
import mock
import pickle
import re

from parameterized import parameterized

from go.chromium.org.luci.buildbucket.proto.build_pb2 import Build
from go.chromium.org.luci.buildbucket.proto.builder_common_pb2 import BuilderID
from google.appengine.api import taskqueue
import webtest

from testing_utils import testing

from common import exceptions
from common.waterfall import buildbucket_client
from common.waterfall import failure_type
import endpoint_api
from findit_v2.model.messages import findit_result
from gae_libs import appengine_util
from libs import analysis_status
from libs import time_util
from model import analysis_approach_type
from model.base_build_model import BaseBuildModel
from model.base_suspected_cl import RevertCL
from model.flake.detection.flake_occurrence import FlakeOccurrence
from model.flake.flake import Flake
from model.flake.flake_issue import FlakeIssue
from model.flake.flake_type import FlakeType
from model.test_inventory import LuciTest
from model.wf_analysis import WfAnalysis
from model.wf_suspected_cl import WfSuspectedCL
from model.wf_swarming_task import WfSwarmingTask
from model.wf_try_job import WfTryJob
from waterfall import suspected_cl_util
from waterfall import waterfall_config

# pylint:disable=unused-argument, unused-variable
# https://crbug.com/947753


# Create a sample flake, and the default properties correspond to that of the
# sample flake occurrence, override the properties as necessary.
def _CreateFlake(**kwargs):
  return Flake.Create(
      luci_project=kwargs.get('luci_project', 'chromium'),
      normalized_step_name=kwargs.get('normalized_step_name', 'browser_tests'),
      normalized_test_name=kwargs.get('test_name', 'foo.bar'),
      test_label_name=kwargs.get('test_label_name', 'foo.bar'))


def _CreateFlakeIssue(**kwargs):
  return FlakeIssue.Create(
      monorail_project=kwargs.get('monorail_project', 'chromium'),
      issue_id=kwargs.get('issue_id', 99999))


# Create a sample flake occurrence, override the properties as necessary.
def _CreateFlakeOccurrence(**kwargs):
  return FlakeOccurrence.Create(
      flake_type=kwargs.get('flake_type', FlakeType.RETRY_WITH_PATCH),
      build_id=kwargs.get('build_id', 1000000),
      step_ui_name=kwargs.get('step_ui_name', 'browser_tests (with patch)'),
      test_name=kwargs.get('test_name', 'foo.bar'),
      luci_project=kwargs.get('luci_project', 'chromium'),
      luci_bucket=kwargs.get('luci_bucket', 'try'),
      luci_builder=kwargs.get('luci_builder', 'linux-rel'),
      legacy_master_name='b',
      legacy_build_number=1,
      time_happened=kwargs.get('time_happened', datetime.datetime(2019, 12,
                                                                  12)),
      gerrit_cl_id=kwargs.get('gerrit_cl_id', 1234),
      parent_flake_key=kwargs.get('parent_flake_key', None))


class FinditApiTest(testing.EndpointsTestCase):
  api_service_cls = endpoint_api.FindItApi

  def setUp(self):
    super(FinditApiTest, self).setUp()
    self.taskqueue_requests = []

    def Mocked_taskqueue_add(**kwargs):
      self.taskqueue_requests.append(kwargs)

    self.mock(taskqueue, 'add', Mocked_taskqueue_add)

  def _MockMasterIsSupported(self, supported):

    def MockMasterIsSupported(*_):
      return supported

    self.mock(waterfall_config, 'MasterIsSupported', MockMasterIsSupported)

  @mock.patch.object(
      endpoint_api.acl,
      'ValidateOauthUserForNewAnalysis',
      side_effect=exceptions.UnauthorizedException)
  def testValidateOauthUserForAuthorizedUser(self, _):
    with self.assertRaises(endpoint_api.endpoints.UnauthorizedException):
      endpoint_api._ValidateOauthUser()

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testUnrecognizedMasterUrl(self, _):
    builds = {
        'builds': [{
            'master_url': 'https://not a master url',
            'builder_name': 'a',
            'build_number': 1
        }]
    }
    expected_results = []

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_results, response.json_body.get('results', []))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testMasterIsNotSupported(self, _):
    builds = {
        'builds': [{
            'master_url': 'https://build.chromium.org/p/a',
            'builder_name': 'a',
            'build_number': 1
        }]
    }
    expected_results = []

    self._MockMasterIsSupported(supported=False)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_results, response.json_body.get('results', []))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testNothingIsReturnedWhenNoAnalysisWasRun(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 1

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number
        }]
    }

    expected_result = []

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_result, response.json_body.get('results', []))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  @mock.patch.object(appengine_util, 'IsStaging', return_Value=True)
  @mock.patch.object(endpoint_api, '_AsyncProcessFailureAnalysisRequests')
  def testNoAnalysisTriggeredOnStaging(self, mock_trigger_analysis, *_):
    master_name = 'm'
    builder_name = 'b'
    build_number = 2

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number
        }]
    }

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.COMPLETED
    analysis.result = {
        'failures': [{
            'step_name':
                'test',
            'first_failure':
                3,
            'last_pass':
                1,
            'supported':
                True,
            'suspected_cls': [{
                'repo_name': 'chromium',
                'revision': 'git_hash',
                'commit_position': 123,
            }]
        }]
    }
    analysis.put()

    expected_result = []

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_result, response.json_body.get('results', []))
    self.assertFalse(mock_trigger_analysis.called)

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testFailedAnalysisIsNotReturnedEvenWhenItHasResults(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 5

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number
        }]
    }

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.ERROR
    analysis.result = {
        'failures': [{
            'step_name':
                'test',
            'first_failure':
                3,
            'last_pass':
                1,
            'supported':
                True,
            'suspected_cls': [{
                'repo_name': 'chromium',
                'revision': 'git_hash',
                'commit_position': 123,
            }]
        }]
    }
    analysis.put()

    expected_result = []

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_result, response.json_body.get('results', []))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testResultIsReturnedWhenNoAnalysisIsCompleted(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 3

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number
        }]
    }

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.RUNNING
    analysis.result = None
    analysis.put()

    expected_result = []

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_result, response.json_body.get('results', []))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testPreviousAnalysisResultIsReturnedWhileANewAnalysisIsRunning(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 4

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number,
            'failed_steps': ['a', 'b']
        }]
    }

    self._MockMasterIsSupported(supported=True)

    analysis_result = {
        'failures': [{
            'step_name':
                'a',
            'first_failure':
                23,
            'last_pass':
                22,
            'supported':
                True,
            'suspected_cls': [{
                'repo_name': 'chromium',
                'revision': 'git_hash',
                'commit_position': 123,
            }]
        }]
    }
    expected_results = [
        {
            'master_url':
                master_url,
            'builder_name':
                builder_name,
            'build_number':
                build_number,
            'step_name':
                'a',
            'is_sub_test':
                False,
            'first_known_failed_build_number':
                23,
            'suspected_cls': [{
                'repo_name': 'chromium',
                'revision': 'git_hash',
                'commit_position': 123,
                'analysis_approach': 'HEURISTIC'
            }],
            'analysis_approach':
                'HEURISTIC',
            'try_job_status':
                'FINISHED',
            'is_flaky_test':
                False,
            'has_findings':
                True,
            'is_finished':
                True,
            'is_supported':
                True,
        },
        {
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number,
            'step_name': 'b',
            'is_sub_test': False,
            'analysis_approach': 'HEURISTIC',
            'is_flaky_test': False,
            'has_findings': False,
            'is_finished': False,
            'is_supported': True,
        },
    ]

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.RUNNING
    analysis.result = analysis_result
    analysis.put()

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(
        sorted(expected_results), sorted(response.json_body['results']))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testAnalysisFindingNoSuspectedCLsIsNotReturned(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 6

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number,
            'failed_steps': ['test']
        }]
    }

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.COMPLETED
    analysis.result = {
        'failures': [{
            'step_name': 'test',
            'first_failure': 3,
            'last_pass': 1,
            'supported': True,
            'suspected_cls': []
        }]
    }
    analysis.put()

    expected_result = [{
        'master_url': master_url,
        'builder_name': builder_name,
        'build_number': build_number,
        'step_name': 'test',
        'is_sub_test': False,
        'first_known_failed_build_number': 3,
        'analysis_approach': 'HEURISTIC',
        'try_job_status': 'FINISHED',
        'is_flaky_test': False,
        'has_findings': False,
        'is_finished': True,
        'is_supported': True,
    }]

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_result, response.json_body.get('results', []))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testAnalysisFindingSuspectedCLsIsReturned(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 7

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number
        }]
    }

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.COMPLETED
    analysis.result = {
        'failures': [{
            'step_name':
                'test',
            'first_failure':
                3,
            'last_pass':
                1,
            'supported':
                True,
            'suspected_cls': [{
                'build_number': 2,
                'repo_name': 'chromium',
                'revision': 'git_hash1',
                'commit_position': 234,
                'score': 11,
                'hints': {
                    'add a/b/x.cc': 5,
                    'delete a/b/y.cc': 5,
                    'modify e/f/z.cc': 1,
                }
            },
                              {
                                  'build_number': 3,
                                  'repo_name': 'chromium',
                                  'revision': 'git_hash2',
                                  'commit_position': 288,
                                  'score': 1,
                                  'hints': {
                                      'modify d/e/f.cc': 1,
                                  }
                              }]
        }]
    }
    analysis.put()

    expected_results = [{
        'master_url':
            master_url,
        'builder_name':
            builder_name,
        'build_number':
            build_number,
        'step_name':
            'test',
        'is_sub_test':
            False,
        'first_known_failed_build_number':
            3,
        'suspected_cls': [{
            'repo_name': 'chromium',
            'revision': 'git_hash1',
            'commit_position': 234,
            'analysis_approach': 'HEURISTIC'
        },
                          {
                              'repo_name': 'chromium',
                              'revision': 'git_hash2',
                              'commit_position': 288,
                              'analysis_approach': 'HEURISTIC'
                          }],
        'analysis_approach':
            'HEURISTIC',
        'is_flaky_test':
            False,
        'try_job_status':
            'FINISHED',
        'has_findings':
            True,
        'is_finished':
            True,
        'is_supported':
            True,
    }]

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_results, response.json_body.get('results'))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testTryJobResultReturnedForCompileFailure(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 8

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number,
            'failed_steps': ['compile']
        }]
    }

    try_job = WfTryJob.Create(master_name, builder_name, 3)
    try_job.status = analysis_status.COMPLETED
    try_job.compile_results = [{
        'culprit': {
            'compile': {
                'repo_name': 'chromium',
                'revision': 'r3',
                'commit_position': 3,
                'url': None,
            },
        },
    }]
    try_job.put()

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.COMPLETED
    analysis.build_failure_type = failure_type.COMPILE
    analysis.failure_result_map = {
        'compile': '/'.join([master_name, builder_name, '3']),
    }
    analysis.result = {
        'failures': [{
            'step_name':
                'compile',
            'first_failure':
                3,
            'last_pass':
                1,
            'supported':
                True,
            'suspected_cls': [{
                'build_number': 3,
                'repo_name': 'chromium',
                'revision': 'git_hash2',
                'commit_position': 288,
                'score': 1,
                'hints': {
                    'modify d/e/f.cc': 1,
                }
            }]
        }]
    }
    analysis.put()

    culprit = WfSuspectedCL.Create('chromium', 'r3', 3)
    culprit.revert_submission_status = analysis_status.COMPLETED
    revert = RevertCL()
    revert.revert_cl_url = 'revert_cl_url'
    culprit.revert_cl = revert
    culprit.put()

    expected_results = [{
        'master_url':
            master_url,
        'builder_name':
            builder_name,
        'build_number':
            build_number,
        'step_name':
            'compile',
        'is_sub_test':
            False,
        'first_known_failed_build_number':
            3,
        'suspected_cls': [{
            'repo_name': 'chromium',
            'revision': 'r3',
            'commit_position': 3,
            'analysis_approach': 'TRY_JOB',
            'revert_cl_url': 'revert_cl_url',
            'revert_committed': True
        },],
        'analysis_approach':
            'TRY_JOB',
        'is_flaky_test':
            False,
        'try_job_status':
            'FINISHED',
        'has_findings':
            True,
        'is_finished':
            True,
        'is_supported':
            True,
    }]

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_results, response.json_body.get('results'))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testTryJobIsRunning(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 9

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number,
            'failed_steps': ['compile']
        }]
    }

    try_job = WfTryJob.Create(master_name, builder_name, 3)
    try_job.status = analysis_status.RUNNING
    try_job.put()

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.COMPLETED
    analysis.build_failure_type = failure_type.COMPILE
    analysis.failure_result_map = {
        'compile': '/'.join([master_name, builder_name, '3']),
    }
    analysis.result = {
        'failures': [{
            'step_name':
                'compile',
            'first_failure':
                3,
            'last_pass':
                1,
            'supported':
                True,
            'suspected_cls': [{
                'build_number': 3,
                'repo_name': 'chromium',
                'revision': 'git_hash2',
                'commit_position': 288,
                'score': 1,
                'hints': {
                    'modify d/e/f.cc': 1,
                }
            }]
        }]
    }
    analysis.put()

    expected_results = [{
        'master_url':
            master_url,
        'builder_name':
            builder_name,
        'build_number':
            build_number,
        'step_name':
            'compile',
        'is_sub_test':
            False,
        'first_known_failed_build_number':
            3,
        'suspected_cls': [{
            'repo_name': 'chromium',
            'revision': 'git_hash2',
            'commit_position': 288,
            'analysis_approach': 'HEURISTIC'
        },],
        'analysis_approach':
            'HEURISTIC',
        'is_flaky_test':
            False,
        'try_job_status':
            'RUNNING',
        'has_findings':
            True,
        'is_finished':
            False,
        'is_supported':
            True,
    }]

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_results, response.json_body.get('results'))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testTestIsFlaky(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 10

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number,
            'failed_steps': ['b on platform']
        }]
    }

    task = WfSwarmingTask.Create(master_name, builder_name, 3, 'b on platform')
    task.tests_statuses = {
        'Unittest3.Subtest1': {
            'total_run': 4,
            'SUCCESS': 2,
            'FAILURE': 2
        },
        'Unittest3.Subtest2': {
            'total_run': 4,
            'SUCCESS': 2,
            'FAILURE': 2
        }
    }
    task.put()

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.COMPLETED
    analysis.failure_result_map = {
        'b on platform': {
            'Unittest3.Subtest1': '/'.join([master_name, builder_name, '3']),
            'Unittest3.Subtest2': '/'.join([master_name, builder_name, '3']),
        },
    }
    analysis.result = {
        'failures': [{
            'step_name':
                'b on platform',
            'first_failure':
                3,
            'last_pass':
                2,
            'supported':
                True,
            'suspected_cls': [],
            'tests': [{
                'test_name': 'Unittest3.Subtest1',
                'first_failure': 3,
                'last_pass': 2,
                'suspected_cls': []
            },
                      {
                          'test_name': 'Unittest3.Subtest2',
                          'first_failure': 3,
                          'last_pass': 2,
                          'suspected_cls': []
                      }]
        }]
    }
    analysis.put()

    expected_results = [{
        'master_url': master_url,
        'builder_name': builder_name,
        'build_number': build_number,
        'step_name': 'b on platform',
        'is_sub_test': True,
        'test_name': 'Unittest3.Subtest1',
        'first_known_failed_build_number': 3,
        'analysis_approach': 'HEURISTIC',
        'is_flaky_test': True,
        'try_job_status': 'FINISHED',
        'has_findings': True,
        'is_finished': True,
        'is_supported': True,
    },
                        {
                            'master_url': master_url,
                            'builder_name': builder_name,
                            'build_number': build_number,
                            'step_name': 'b on platform',
                            'is_sub_test': True,
                            'test_name': 'Unittest3.Subtest2',
                            'first_known_failed_build_number': 3,
                            'analysis_approach': 'HEURISTIC',
                            'is_flaky_test': True,
                            'try_job_status': 'FINISHED',
                            'has_findings': True,
                            'is_finished': True,
                            'is_supported': True,
                        }]

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_results, response.json_body.get('results'))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  @mock.patch.object(suspected_cl_util,
                     'GetSuspectedCLConfidenceScoreAndApproach')
  def testTestLevelResultIsReturned(self, mock_fn, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 11

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number,
            'failed_steps': ['a', 'b on platform']
        }]
    }

    task = WfSwarmingTask.Create(master_name, builder_name, 4, 'b on platform')
    task.parameters['ref_name'] = 'b'
    task.status = analysis_status.COMPLETED
    task.put()

    try_job = WfTryJob.Create(master_name, builder_name, 4)
    try_job.status = analysis_status.COMPLETED
    try_job.test_results = [{
        'culprit': {
            'a': {
                'repo_name': 'chromium',
                'revision': 'r4_2',
                'commit_position': 42,
                'url': None,
            },
            'b': {
                'tests': {
                    'Unittest3.Subtest1': {
                        'repo_name': 'chromium',
                        'revision': 'r4_10',
                        'commit_position': 410,
                        'url': None,
                    },
                }
            }
        },
    }]
    try_job.put()

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.COMPLETED
    analysis.failure_result_map = {
        'a': '/'.join([master_name, builder_name, '4']),
        'b on platform': {
            'Unittest1.Subtest1': '/'.join([master_name, builder_name, '3']),
            'Unittest2.Subtest1': '/'.join([master_name, builder_name, '4']),
            'Unittest3.Subtest1': '/'.join([master_name, builder_name, '4']),
        },
    }
    analysis.result = {
        'failures': [
            {
                'step_name':
                    'a',
                'first_failure':
                    4,
                'last_pass':
                    3,
                'supported':
                    True,
                'suspected_cls': [{
                    'build_number': 4,
                    'repo_name': 'chromium',
                    'revision': 'r4_2_failed',
                    'commit_position': None,
                    'url': None,
                    'score': 2,
                    'hints': {
                        'modified f4_2.cc (and it was in log)': 2,
                    },
                }],
            },
            {
                'step_name':
                    'b on platform',
                'first_failure':
                    3,
                'last_pass':
                    2,
                'supported':
                    True,
                'suspected_cls': [
                    {
                        'build_number': 3,
                        'repo_name': 'chromium',
                        'revision': 'r3_1',
                        'commit_position': None,
                        'url': None,
                        'score': 5,
                        'hints': {
                            'added x/y/f3_1.cc (and it was in log)': 5,
                        },
                    },
                    {
                        'build_number': 4,
                        'repo_name': 'chromium',
                        'revision': 'r4_1',
                        'commit_position': None,
                        'url': None,
                        'score': 2,
                        'hints': {
                            'modified f4.cc (and it was in log)': 2,
                        },
                    }
                ],
                'tests': [
                    {
                        'test_name':
                            'Unittest1.Subtest1',
                        'first_failure':
                            3,
                        'last_pass':
                            2,
                        'suspected_cls': [{
                            'build_number': 2,
                            'repo_name': 'chromium',
                            'revision': 'r2_1',
                            'commit_position': None,
                            'url': None,
                            'score': 5,
                            'hints': {
                                'added x/y/f99_1.cc (and it was in log)': 5,
                            },
                        }]
                    },
                    {
                        'test_name':
                            'Unittest2.Subtest1',
                        'first_failure':
                            4,
                        'last_pass':
                            2,
                        'suspected_cls': [{
                            'build_number': 2,
                            'repo_name': 'chromium',
                            'revision': 'r2_1',
                            'commit_position': None,
                            'url': None,
                            'score': 5,
                            'hints': {
                                'added x/y/f99_1.cc (and it was in log)': 5,
                            },
                        }]
                    },
                    {
                        'test_name': 'Unittest3.Subtest1',
                        'first_failure': 4,
                        'last_pass': 2,
                        'suspected_cls': []
                    }
                ]
            },
            {
                'step_name': 'c',
                'first_failure': 4,
                'last_pass': 3,
                'supported': False,
                'suspected_cls': [],
            }
        ]
    }
    analysis.put()

    suspected_cl_42 = WfSuspectedCL.Create('chromium', 'r4_2', 42)
    suspected_cl_42.builds = {
        BaseBuildModel.CreateBuildKey(master_name, builder_name, 5): {
            'approaches': [analysis_approach_type.TRY_JOB]
        }
    }
    suspected_cl_42.put()

    suspected_cl_21 = WfSuspectedCL.Create('chromium', 'r2_1', None)
    suspected_cl_21.builds = {
        BaseBuildModel.CreateBuildKey(master_name, builder_name, 3): {
            'approaches': [analysis_approach_type.HEURISTIC],
            'top_score': 5
        },
        BaseBuildModel.CreateBuildKey(master_name, builder_name, 4): {
            'approaches': [analysis_approach_type.HEURISTIC],
            'top_score': 5
        },
        BaseBuildModel.CreateBuildKey(master_name, builder_name,
                                      build_number): {
            'approaches': [analysis_approach_type.HEURISTIC],
            'top_score': 5
        }
    }
    suspected_cl_21.put()

    suspected_cl_410 = WfSuspectedCL.Create('chromium', 'r4_10', None)
    suspected_cl_410.builds = {
        BaseBuildModel.CreateBuildKey(master_name, builder_name, 4): {
            'approaches': [
                analysis_approach_type.HEURISTIC, analysis_approach_type.TRY_JOB
            ],
            'top_score':
                5
        },
        BaseBuildModel.CreateBuildKey(master_name, builder_name,
                                      build_number): {
            'approaches': [analysis_approach_type.HEURISTIC],
            'top_score': 5
        }
    }
    revert_cl = RevertCL()
    revert_cl.revert_cl_url = 'revert_cl_url'
    suspected_cl_410.revert_cl = revert_cl
    suspected_cl_410.put()

    def confidence_side_effect(_, build_info, first_build_info):
      if (first_build_info and first_build_info.get('approaches') == [
          analysis_approach_type.HEURISTIC, analysis_approach_type.TRY_JOB
      ]):
        return 100, analysis_approach_type.TRY_JOB
      if build_info and build_info.get('top_score'):
        return 90, analysis_approach_type.HEURISTIC
      return 98, analysis_approach_type.TRY_JOB

    mock_fn.side_effect = confidence_side_effect

    expected_results = [{
        'master_url':
            master_url,
        'builder_name':
            builder_name,
        'build_number':
            build_number,
        'step_name':
            'a',
        'is_sub_test':
            False,
        'first_known_failed_build_number':
            4,
        'suspected_cls': [{
            'repo_name': 'chromium',
            'revision': 'r4_2',
            'commit_position': 42,
            'confidence': 98,
            'analysis_approach': 'TRY_JOB',
            'revert_committed': False
        }],
        'analysis_approach':
            'TRY_JOB',
        'is_flaky_test':
            False,
        'try_job_status':
            'FINISHED',
        'has_findings':
            True,
        'is_finished':
            True,
        'is_supported':
            True,
    },
                        {
                            'master_url':
                                master_url,
                            'builder_name':
                                builder_name,
                            'build_number':
                                build_number,
                            'step_name':
                                'b on platform',
                            'is_sub_test':
                                True,
                            'test_name':
                                'Unittest1.Subtest1',
                            'first_known_failed_build_number':
                                3,
                            'suspected_cls': [{
                                'repo_name': 'chromium',
                                'revision': 'r2_1',
                                'confidence': 90,
                                'analysis_approach': 'HEURISTIC',
                                'revert_committed': False
                            }],
                            'analysis_approach':
                                'HEURISTIC',
                            'is_flaky_test':
                                False,
                            'try_job_status':
                                'FINISHED',
                            'has_findings':
                                True,
                            'is_finished':
                                True,
                            'is_supported':
                                True,
                        },
                        {
                            'master_url':
                                master_url,
                            'builder_name':
                                builder_name,
                            'build_number':
                                build_number,
                            'step_name':
                                'b on platform',
                            'is_sub_test':
                                True,
                            'test_name':
                                'Unittest2.Subtest1',
                            'first_known_failed_build_number':
                                4,
                            'suspected_cls': [{
                                'repo_name': 'chromium',
                                'revision': 'r2_1',
                                'confidence': 90,
                                'analysis_approach': 'HEURISTIC',
                                'revert_committed': False
                            }],
                            'analysis_approach':
                                'HEURISTIC',
                            'is_flaky_test':
                                False,
                            'try_job_status':
                                'FINISHED',
                            'has_findings':
                                True,
                            'is_finished':
                                True,
                            'is_supported':
                                True,
                        },
                        {
                            'master_url':
                                master_url,
                            'builder_name':
                                builder_name,
                            'build_number':
                                build_number,
                            'step_name':
                                'b on platform',
                            'is_sub_test':
                                True,
                            'test_name':
                                'Unittest3.Subtest1',
                            'first_known_failed_build_number':
                                4,
                            'suspected_cls': [{
                                'repo_name': 'chromium',
                                'revision': 'r4_10',
                                'commit_position': 410,
                                'analysis_approach': 'TRY_JOB',
                                'confidence': 100,
                                'revert_cl_url': 'revert_cl_url',
                                'revert_committed': False
                            }],
                            'analysis_approach':
                                'TRY_JOB',
                            'is_flaky_test':
                                False,
                            'try_job_status':
                                'FINISHED',
                            'has_findings':
                                True,
                            'is_finished':
                                True,
                            'is_supported':
                                True,
                        },
                        {
                            'master_url': master_url,
                            'builder_name': builder_name,
                            'build_number': build_number,
                            'step_name': 'c',
                            'is_sub_test': False,
                            'analysis_approach': 'HEURISTIC',
                            'is_flaky_test': False,
                            'has_findings': False,
                            'is_finished': True,
                            'is_supported': False,
                        }]

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertItemsEqual(expected_results, response.json_body.get('results'))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testAnalysisRequestQueuedAsExpected(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 12

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number
        }]
    }

    expected_result = []

    self._MockMasterIsSupported(supported=True)

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_result, response.json_body.get('results', []))
    self.assertEqual(1, len(self.taskqueue_requests))

    expected_payload_json = {
        'builds': [{
            'master_name': master_name,
            'builder_name': builder_name,
            'build_number': build_number,
            'failed_steps': [],
        },]
    }
    self.assertEqual(expected_payload_json,
                     json.loads(self.taskqueue_requests[0].get('payload')))

  @mock.patch.object(
      endpoint_api,
      '_ValidateOauthUser',
      side_effect=endpoint_api.endpoints.UnauthorizedException())
  @mock.patch.object(endpoint_api, 'AsyncProcessFlakeReport', return_value=None)
  def testUnauthorizedRequestToAnalyzeFlake(self, mocked_func, _):
    flake = {
        'name':
            'suite.test',
        'is_step':
            False,
        'bug_id':
            123,
        'build_steps': [{
            'master_name': 'm',
            'builder_name': 'b',
            'build_number': 456,
            'step_name': 'name (with patch) on Windows-7-SP1',
        }]
    }

    self.assertRaisesRegexp(
        webtest.app.AppError,
        re.compile('.*401 Unauthorized.*', re.MULTILINE | re.DOTALL),
        self.call_api,
        'AnalyzeFlake',
        body=flake)
    self.assertFalse(mocked_func.called)

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testFlakeAnalysisRequestWithoutBugId(self, _):
    flake = {
        'name':
            'suite.test',
        'is_step':
            False,
        'bug_id':
            None,
        'build_steps': [{
            'master_name': 'm',
            'builder_name': 'b',
            'build_number': 456,
            'step_name': 'name (with patch) on Windows-7-SP1',
        }]
    }

    response = self.call_api('AnalyzeFlake', body=flake)
    self.assertEqual(200, response.status_int)
    self.assertTrue(response.json_body.get('queued'))
    self.assertEqual(1, len(self.taskqueue_requests))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  @mock.patch.object(
      endpoint_api, 'AsyncProcessFlakeReport', side_effect=Exception())
  def testAuthorizedRequestToAnalyzeFlakeNotQueued(self, mocked_func, _):
    flake = {
        'name':
            'suite.test',
        'is_step':
            False,
        'bug_id':
            123,
        'build_steps': [{
            'master_name': 'm',
            'builder_name': 'b',
            'build_number': 456,
            'step_name': 'name (with patch) on Windows-7-SP1',
        }]
    }

    response = self.call_api('AnalyzeFlake', body=flake)
    self.assertEqual(200, response.status_int)
    self.assertFalse(response.json_body.get('queued'))
    self.assertEqual(1, mocked_func.call_count)
    self.assertEqual(0, len(self.taskqueue_requests))

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testAuthorizedRequestToAnalyzeFlakeQueued(self, _):
    flake = {
        'name':
            'suite.test',
        'is_step':
            False,
        'bug_id':
            123,
        'build_steps': [{
            'master_name': 'm',
            'builder_name': 'b',
            'build_number': 456,
            'step_name': 'name (with patch) on Windows-7-SP1',
        }]
    }

    response = self.call_api('AnalyzeFlake', body=flake)
    self.assertEqual(200, response.status_int)
    self.assertTrue(response.json_body.get('queued'))
    self.assertEqual(1, len(self.taskqueue_requests))
    request, user_email, is_admin = pickle.loads(
        self.taskqueue_requests[0]['payload'])
    self.assertEqual('suite.test', request.name)
    self.assertFalse(request.is_step)
    self.assertEqual(123, request.bug_id)
    self.assertEqual(1, len(request.build_steps))
    self.assertEqual('m', request.build_steps[0].master_name)
    self.assertEqual('b', request.build_steps[0].builder_name)
    self.assertEqual(456, request.build_steps[0].build_number)
    self.assertEqual('name (with patch) on Windows-7-SP1',
                     request.build_steps[0].step_name)
    self.assertEqual('email', user_email)
    self.assertFalse(is_admin)

  def testGetStatusAndCulpritFromTryJobSwarmingTaskIsRunning(self):
    swarming_task = WfSwarmingTask.Create('m', 'b', 123, 'step')
    swarming_task.put()
    status, culprit = endpoint_api.FindItApi()._GetStatusAndCulpritFromTryJob(
        None, swarming_task, None, 'step', None)
    self.assertEqual(status, endpoint_api._TryJobStatus.RUNNING)
    self.assertIsNone(culprit)

  def testGetStatusAndCulpritFromTryJobTryJobFailed(self):
    try_job = WfTryJob.Create('m', 'b', 123)
    try_job.status = analysis_status.ERROR
    try_job.put()
    status, culprit = endpoint_api.FindItApi()._GetStatusAndCulpritFromTryJob(
        try_job, None, None, None, None)
    self.assertEqual(status, endpoint_api._TryJobStatus.FINISHED)
    self.assertIsNone(culprit)

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  def testAnalysisIsStillRunning(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 13

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number,
            'failed_steps': ['a']
        }]
    }

    self._MockMasterIsSupported(supported=True)

    expected_results = [{
        'master_url': master_url,
        'builder_name': builder_name,
        'build_number': build_number,
        'step_name': 'a',
        'analysis_approach': 'HEURISTIC',
        'is_sub_test': False,
        'is_flaky_test': False,
        'has_findings': False,
        'is_finished': False,
        'is_supported': True,
    }]

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.RUNNING
    analysis.result = None
    analysis.put()

    response = self.call_api('AnalyzeBuildFailures', body=builds)
    self.assertEqual(200, response.status_int)
    self.assertEqual(expected_results, response.json_body['results'])

  @mock.patch.object(
      endpoint_api,
      '_ValidateOauthUser',
      side_effect=endpoint_api.endpoints.UnauthorizedException('Unauthorized.'))
  @mock.patch.object(endpoint_api, '_AsyncProcessFailureAnalysisRequests')
  def testUserNotAuthorized(self, mocked_func, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 14

    master_url = 'https://build.chromium.org/p/%s' % master_name
    builds = {
        'builds': [{
            'master_url': master_url,
            'builder_name': builder_name,
            'build_number': build_number,
            'failed_steps': ['a']
        }]
    }

    self.assertRaisesRegexp(
        webtest.app.AppError,
        re.compile('.*401 Unauthorized.*', re.MULTILINE | re.DOTALL),
        self.call_api,
        'AnalyzeBuildFailures',
        body=builds)
    self.assertFalse(mocked_func.called)

  @mock.patch.object(logging, 'error')
  def testGetSwarmingTaskAndTryJobForFailureInfoMismatch(self, mock_log):
    failure_result_map = {'s': {'t': 'm/b/1'}}
    self.assertEqual(
        (None, None, None),
        endpoint_api.FindItApi()._GetSwarmingTaskAndTryJobForFailure(
            's', None, failure_result_map, None, None))
    # mock_log.assert_called_once_with(
    #     'Try_job_key in wrong format - failure_result_map: %s; step_name: %s;'
    #     ' test_name: %s.', json.dumps(failure_result_map, default=str), 's',
    #     None)

  @mock.patch.object(
      endpoint_api.FindItApi, '_GetV2AnalysisResultFromV1', return_value=None)
  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  @mock.patch.object(logging, 'info')
  @mock.patch(
      'endpoint_api.findit_v2_api.OnBuildFailureAnalysisResultRequested')
  def testAnalyzeLuciBuildFailures(self, mock_api, mock_logging, *_):
    api_input = {
        'requests': [
            {
                'build_id': 8000000000123,
                'failed_steps': ['a']
            },
            {
                'build_alternative_id': {
                    'project': 'chromium',
                    'bucket': 'ci',
                    'builder': 'Luci Tests',
                    'number': 124
                },
                'failed_steps': ['compile']
            },
        ]
    }

    mock_api.side_effect = [
        [],
        [
            findit_result.BuildFailureAnalysisResponse(
                build_alternative_id=findit_result.BuildIdentifierByNumber(
                    project='chromium',
                    bucket='ci',
                    builder='Luci Tests',
                    number=124),
                is_finished=False,
            )
        ]
    ]

    response = self.call_api('AnalyzeLuciBuildFailures', body=api_input)
    self.assertEqual(200, response.status_int)
    mock_logging.assert_called_once_with(
        '%d build failure(s), while findit_v2 can provide results for%d, and'
        ' findit_v1 can provide results for %d.', 2, 1, 0)

  @mock.patch.object(
      endpoint_api, '_ValidateOauthUser', return_value=('email', False))
  @mock.patch.object(buildbucket_client, 'GetV2BuildByBuilderAndBuildNumber')
  @mock.patch.object(buildbucket_client, 'GetV2Build')
  @mock.patch.object(logging, 'info')
  @mock.patch(
      'endpoint_api.findit_v2_api.OnBuildFailureAnalysisResultRequested')
  def testGetV1AnalysesResults(self, mock_api, mock_logging,
                               mock_get_build_by_id, mock_get_build_by_number,
                               _):
    api_input = {
        'requests': [
            {
                'build_id': 8000000000123,
                'failed_steps': ['a']
            },
            {
                'build_alternative_id': {
                    'project': 'chromium',
                    'bucket': 'ci',
                    'builder': 'Luci Tests',
                    'number': 124
                },
                'failed_steps': ['compile']
            },
        ]
    }

    mock_api.return_value = []

    mock_build1 = Build(
        id=8000000000123,
        builder=BuilderID(project='chromeos', bucket='ci', builder='builder'))
    mock_get_build_by_id.return_value = mock_build1

    master_name = 'chromium.linux'
    builder_name = 'Luci Tests'
    build_number = 124
    mock_build2 = Build(
        id=8000000000124,
        builder=BuilderID(
            project='chromium', bucket='ci', builder=builder_name),
        number=build_number)
    mock_build2.output.properties['builder_group'] = master_name
    mock_get_build_by_number.return_value = mock_build2

    analysis = WfAnalysis.Create(master_name, builder_name, build_number)
    analysis.status = analysis_status.COMPLETED
    analysis.result = {
        'failures': [{
            'step_name':
                'test',
            'first_failure':
                3,
            'last_pass':
                1,
            'supported':
                True,
            'suspected_cls': [{
                'build_number': 2,
                'repo_name': 'chromium',
                'revision': 'git_hash1',
                'commit_position': 234,
                'score': 11,
                'hints': {
                    'add a/b/x.cc': 5,
                    'delete a/b/y.cc': 5,
                    'modify e/f/z.cc': 1,
                }
            },
                              {
                                  'build_number': 3,
                                  'repo_name': 'chromium',
                                  'revision': 'git_hash2',
                                  'commit_position': 288,
                                  'score': 1,
                                  'hints': {
                                      'modify d/e/f.cc': 1,
                                  }
                              }]
        }]
    }
    analysis.put()

    response = self.call_api('AnalyzeLuciBuildFailures', body=api_input)
    self.assertEqual(200, response.status_int)
    mock_logging.assert_called_once_with(
        '%d build failure(s), while findit_v2 can provide results for%d, and'
        ' findit_v1 can provide results for %d.', 2, 0, 1)

  @parameterized.expand([
      ({
          'request_body': {},
          'expected_test_data': [
              {
                  'luci_project': 'chromium',
                  'normalized_step_name': 'normal_step1',
                  'normalized_test_name': 'test_1',
              },
              {
                  'luci_project': 'chromium',
                  'normalized_step_name': 'normal_step2',
                  'normalized_test_name': 'test_2',
              },
              {
                  'luci_project': 'chromium',
                  'normalized_step_name': 'normal_step3',
                  'normalized_test_name': 'test_3',
              },
          ],
      },),
      ({
          'request_body': {
              'include_tags': [
                  'component::mock_component',
                  'test_type::mock_step',
              ],
              'request_type':
                  'ALL',
          },
          'expected_test_data': [
              {
                  'luci_project': 'chromium',
                  'normalized_step_name': 'normal_step1',
                  'normalized_test_name': 'test_1',
                  'disabled_test_variants': [{
                      'variant': ['os:Mac'],
                  },],
              },
              {
                  'luci_project': 'chromium',
                  'normalized_step_name': 'normal_step2',
                  'normalized_test_name': 'test_2',
                  'disabled_test_variants': [{
                      'variant': ['os:Mac'],
                  },],
              },
          ],
      },),
      ({
          'request_body': {
              'include_tags': [
                  'component::mock_component',
                  'test_type::mock_step',
              ],
              'exclude_tags': ['step::mock_step (with patch)',],
              'request_type':
                  'ALL',
          },
          'expected_test_data': [{
              'luci_project': 'chromium',
              'normalized_step_name': 'normal_step2',
              'normalized_test_name': 'test_2',
              'disabled_test_variants': [{
                  'variant': ['os:Mac'],
              },],
          },],
      },),
      ({
          'request_body': {
              'include_tags': [
                  'component::mock_component',
                  'test_type::mock_step',
              ],
              'request_type':
                  'NAME_ONLY',
          },
          'expected_test_data': [
              {
                  'luci_project': 'chromium',
                  'normalized_step_name': 'normal_step1',
                  'normalized_test_name': 'test_1'
              },
              {
                  'luci_project': 'chromium',
                  'normalized_step_name': 'normal_step2',
                  'normalized_test_name': 'test_2'
              },
          ],
      },),
      ({
          'request_body': {
              'include_tags': [
                  'component::mock_component',
                  'test_type::mock_step',
              ],
              'exclude_tags': ['step::mock_step (with patch)',],
              'request_type':
                  'NAME_ONLY',
          },
          'expected_test_data': [{
              'luci_project': 'chromium',
              'normalized_step_name': 'normal_step2',
              'normalized_test_name': 'test_2'
          },],
      },),
  ])
  @mock.patch.object(endpoint_api, '_ValidateOauthUser')
  def testFilterDisabledTestsTestData(self, cases, _):
    test_1 = LuciTest(
        key=LuciTest.CreateKey('chromium', 'normal_step1', 'test_1'),
        disabled_test_variants={('os:Mac123',), ('os:Mac124',)},
        tags=[
            'component::mock_component',
            'step::mock_step (with patch)',
            'test_type::mock_step',
        ])
    test_1.put()
    test_2 = LuciTest(
        key=LuciTest.CreateKey('chromium', 'normal_step2', 'test_2'),
        disabled_test_variants={('os:Mac123',), ('os:Mac124',)},
        tags=[
            'component::mock_component',
            'step::mock_step (without patch)',
            'test_type::mock_step',
        ])
    test_2.put()
    test_3 = LuciTest(
        key=LuciTest.CreateKey('chromium', 'normal_step3', 'test_3'),
        disabled_test_variants={('os:Mac123',), ('os:Mac124',)},
        tags=[])
    test_3.put()
    test_4 = LuciTest(
        key=LuciTest.CreateKey('chromium', 'normal_step4', 'test_4'),
        disabled_test_variants=set(),
        tags=[
            'component::mock_component',
            'test_type::mock_step',
        ])
    test_4.put()
    test_3.put()

    response = self.call_api('FilterDisabledTests', body=cases['request_body'])
    self.assertEqual(200, response.status_int)

    actual_test_data = response.json_body.get('test_data', [])
    for expected_test in cases['expected_test_data']:
      self.assertIn(expected_test, actual_test_data)

  @parameterized.expand([
      ({
          'request_body': {
              'include_tags': [
                  'component::mock_component',
                  'test_type::mock_step',
              ],
              'request_type':
                  'COUNT',
          },
          'expected_test_count': 2
      },),
      ({
          'request_body': {
              'include_tags': [
                  'component::mock_component',
                  'test_type::mock_step',
              ],
              'exclude_tags': ['step::mock_step (with patch)',],
              'request_type':
                  'COUNT',
          },
          'expected_test_count': 1
      },),
      ({
          'request_body': {
              'include_tags': [
                  'component::mock_component',
                  'test_type::mock_step',
              ],
              'request_type':
                  'COUNT',
          },
          'expected_test_count': 2
      },),
      ({
          'request_body': {
              'include_tags': [
                  'component::mock_component',
                  'test_type::mock_step',
              ],
              'exclude_tags': ['step::mock_step (with patch)',],
              'request_type':
                  'COUNT',
          },
          'expected_test_count': 1
      },),
      ({
          'request_body': {
              'request_type': 'COUNT',
          },
          'expected_test_count': 3
      },),
  ])
  @mock.patch.object(endpoint_api, '_ValidateOauthUser')
  def testFilterDisabledTestsTestCount(self, cases, _):
    test_1 = LuciTest(
        key=LuciTest.CreateKey('chromium', 'normal_step1', 'test_1'),
        disabled_test_variants={('os:Mac123',), ('os:Mac124',)},
        tags=[
            'component::mock_component',
            'step::mock_step (with patch)',
            'test_type::mock_step',
        ])
    test_1.put()
    test_2 = LuciTest(
        key=LuciTest.CreateKey('chromium', 'normal_step2', 'test_2'),
        disabled_test_variants={('os:Mac123',), ('os:Mac124',)},
        tags=[
            'component::mock_component',
            'step::mock_step (without patch)',
            'test_type::mock_step',
        ])
    test_2.put()
    test_3 = LuciTest(
        key=LuciTest.CreateKey('chromium', 'normal_step3', 'test_3'),
        disabled_test_variants={('os:Mac123',), ('os:Mac124',)},
        tags=[])
    test_3.put()
    test_4 = LuciTest(
        key=LuciTest.CreateKey('chromium', 'normal_step4', 'test_4'),
        disabled_test_variants=set(),
        tags=[
            'component::mock_component',
            'test_type::mock_step',
        ])
    test_4.put()

    response = self.call_api('FilterDisabledTests', body=cases['request_body'])
    self.assertEqual(200, response.status_int)

    self.assertEqual(cases['expected_test_count'],
                     response.json_body.get('test_count'))

  def testGetCQFlakesInvalidRequestsMissingBuilder(self):
    request = {
        'project': 'chromium',
        'bucket': 'try',
        'tests': [],
    }

    response = self.call_api('GetCQFlakes', body=request, status=400)

  def testGetCQFlakesInvalidRequestsMissingTest(self):
    request = {
        'project': 'chromium',
        'bucket': 'try',
        'tests': [{
            'step_ui_name': 'browser_tests (with patch)',
        }],
    }

    response = self.call_api('GetCQFlakes', body=request, status=400)

  def testGetCQFlakesNoFlakes(self):
    request = {
        'project': 'chromium',
        'bucket': 'try',
        'builder': 'linux-rel',
        'tests': [],
    }

    response = self.call_api('GetCQFlakes', body=request)
    self.assertEqual(200, response.status_int)
    self.assertDictEqual({}, response.json_body)

  @mock.patch.object(
      time_util, 'GetUTCNow', return_value=datetime.datetime(2019, 12, 10))
  def testGetCQFlakes(self, _):
    destination_flake_issue = _CreateFlakeIssue(issue_id=99995)
    destination_flake_issue.put()
    flake_issue = _CreateFlakeIssue()
    flake_issue.merge_destination_key = destination_flake_issue.key
    flake_issue.put()
    flake = _CreateFlake()
    flake.flake_issue_key = flake_issue.key
    flake.put()

    _CreateFlakeOccurrence(
        build_id=987,
        time_happened=time_util.GetDatetimeBeforeNow(hours=1),
        gerrit_cl_id=1234,
        parent_flake_key=flake.key).put()
    _CreateFlakeOccurrence(
        build_id=986,
        time_happened=time_util.GetDatetimeBeforeNow(hours=14),
        gerrit_cl_id=1235,
        parent_flake_key=flake.key).put()
    _CreateFlakeOccurrence(
        build_id=985,
        time_happened=time_util.GetDatetimeBeforeNow(hours=15),
        gerrit_cl_id=1236,
        parent_flake_key=flake.key).put()

    request = {
        'project':
            'chromium',
        'bucket':
            'try',
        'builder':
            'linux-rel',
        'tests': [{
            'step_ui_name': 'browser_tests (with patch)',
            'test_name': 'foo.bar',
        },],
    }

    response = self.call_api('GetCQFlakes', body=request)

    self.assertEqual(200, response.status_int)
    self.assertDictEqual({
        'flakes': [{
            'test': {
                'step_ui_name': 'browser_tests (with patch)',
                'test_name': 'foo.bar',
            },
            'affected_gerrit_changes': ['1234', '1235', '1236'],
            'monorail_issue': '99995',
        }]
    }, response.json_body)

  # This test tests that enough occurrences are required in order for a test to
  # be used to skip retrying known flakes.
  @mock.patch.object(
      time_util, 'GetUTCNow', return_value=datetime.datetime(2019, 12, 10))
  def testGetCQFlakesNotEnoughOccurrences(self, _):
    flake_issue = _CreateFlakeIssue()
    flake = _CreateFlake()
    flake.flake_issue_key = flake_issue.key
    flake_issue.put()
    flake.put()

    _CreateFlakeOccurrence(
        build_id=987,
        time_happened=time_util.GetDatetimeBeforeNow(hours=1),
        gerrit_cl_id=1234).put()
    _CreateFlakeOccurrence(
        build_id=986,
        time_happened=time_util.GetDatetimeBeforeNow(hours=14),
        gerrit_cl_id=1235).put()

    request = {
        'project':
            'chromium',
        'bucket':
            'try',
        'builder':
            'linux-rel',
        'tests': [{
            'step_ui_name': 'browser_tests (with patch)',
            'test_name': 'foo.bar',
        },],
    }

    response = self.call_api('GetCQFlakes', body=request)
    self.assertEqual(200, response.status_int)
    self.assertEqual({}, response.json_body)

  # This test tests that enough uniquely affected CLs are required in order for
  # a test to be used to skip retrying known flakes.
  @mock.patch.object(
      time_util, 'GetUTCNow', return_value=datetime.datetime(2019, 12, 10))
  def testGetCQFlakesNotEnoughGerritChanges(self, _):
    flake_issue = _CreateFlakeIssue()
    flake = _CreateFlake()
    flake.flake_issue_key = flake_issue.key
    flake_issue.put()
    flake.put()

    # Two of the occurrences share the same gerrit change.
    _CreateFlakeOccurrence(
        build_id=987,
        time_happened=time_util.GetDatetimeBeforeNow(hours=1),
        gerrit_cl_id=1234).put()
    _CreateFlakeOccurrence(
        build_id=986,
        time_happened=time_util.GetDatetimeBeforeNow(hours=14),
        gerrit_cl_id=1235).put()
    _CreateFlakeOccurrence(
        build_id=985,
        time_happened=time_util.GetDatetimeBeforeNow(hours=15),
        gerrit_cl_id=1235).put()

    request = {
        'project':
            'chromium',
        'bucket':
            'try',
        'builder':
            'linux-rel',
        'tests': [{
            'step_ui_name': 'browser_tests (with patch)',
            'test_name': 'foo.bar',
        },],
    }

    response = self.call_api('GetCQFlakes', body=request)
    self.assertEqual(200, response.status_int)
    self.assertDictEqual({}, response.json_body)

  # This test tests that at least one recent flake occurrence is required in
  # order for a test to be determined as flaky.
  @mock.patch.object(
      time_util, 'GetUTCNow', return_value=datetime.datetime(2019, 12, 10))
  def testGetCQFlakesNoRecentActivity(self, _):
    flake_issue = _CreateFlakeIssue()
    flake = _CreateFlake()
    flake.flake_issue_key = flake_issue.key
    flake_issue.put()
    flake.put()

    _CreateFlakeOccurrence(
        build_id=987,
        time_happened=time_util.GetDatetimeBeforeNow(hours=13),
        gerrit_cl_id=1234).put()
    _CreateFlakeOccurrence(
        build_id=986,
        time_happened=time_util.GetDatetimeBeforeNow(hours=14),
        gerrit_cl_id=1235).put()
    _CreateFlakeOccurrence(
        build_id=985,
        time_happened=time_util.GetDatetimeBeforeNow(hours=15),
        gerrit_cl_id=1236).put()

    request = {
        'project':
            'chromium',
        'bucket':
            'try',
        'builder':
            'linux-rel',
        'tests': [{
            'step_ui_name': 'browser_tests (with patch)',
            'test_name': 'foo.bar',
        },],
    }

    response = self.call_api('GetCQFlakes', body=request)
    self.assertEqual(200, response.status_int)
    self.assertDictEqual({}, response.json_body)

  # This test tests that a bug must be filed in order for a test to be used to
  # skip retrying known flakes.
  @mock.patch.object(
      time_util, 'GetUTCNow', return_value=datetime.datetime(2019, 12, 10))
  def testGetCQFlakesNoBugFiled(self, _):
    flake = _CreateFlake()
    flake.put()

    _CreateFlakeOccurrence(
        build_id=987,
        time_happened=time_util.GetDatetimeBeforeNow(hours=1),
        gerrit_cl_id=1234,
        parent_flake_key=flake.key).put()
    _CreateFlakeOccurrence(
        build_id=986,
        time_happened=time_util.GetDatetimeBeforeNow(hours=14),
        gerrit_cl_id=1235,
        parent_flake_key=flake.key).put()
    _CreateFlakeOccurrence(
        build_id=985,
        time_happened=time_util.GetDatetimeBeforeNow(hours=15),
        gerrit_cl_id=1236,
        parent_flake_key=flake.key).put()

    request = {
        'project':
            'chromium',
        'bucket':
            'try',
        'builder':
            'linux-rel',
        'tests': [{
            'step_ui_name': 'browser_tests (with patch)',
            'test_name': 'foo.bar',
        },],
    }

    response = self.call_api('GetCQFlakes', body=request)
    self.assertEqual(200, response.status_int)
    self.assertDictEqual({}, response.json_body)

  # This test tests that the same step/test names on different builders will be
  # treated as two different tests.
  @mock.patch.object(
      time_util, 'GetUTCNow', return_value=datetime.datetime(2019, 12, 10))
  def testGetCQFlakesMultipleBuilders(self, _):
    flake_issue = _CreateFlakeIssue()
    flake = _CreateFlake()
    flake.flake_issue_key = flake_issue.key
    flake_issue.put()
    flake.put()

    _CreateFlakeOccurrence(
        build_id=987,
        luci_builder='linux-rel',
        time_happened=time_util.GetDatetimeBeforeNow(hours=1),
        gerrit_cl_id=1234).put()
    _CreateFlakeOccurrence(
        build_id=986,
        luci_builder='linux-rel',
        time_happened=time_util.GetDatetimeBeforeNow(hours=14),
        gerrit_cl_id=1235).put()
    _CreateFlakeOccurrence(
        build_id=985,
        luci_builder='linux-chromeos-rel',
        time_happened=time_util.GetDatetimeBeforeNow(hours=15),
        gerrit_cl_id=1235).put()

    request = {
        'project':
            'chromium',
        'bucket':
            'try',
        'builder':
            'linux-rel',
        'tests': [{
            'step_ui_name': 'browser_tests (with patch)',
            'test_name': 'foo.bar',
        },],
    }

    response = self.call_api('GetCQFlakes', body=request)
    self.assertEqual(200, response.status_int)
    self.assertDictEqual({}, response.json_body)

    request = {
        'project':
            'chromium',
        'bucket':
            'try',
        'builder':
            'linux-chromeos-rel',
        'tests': [{
            'step_ui_name': 'browser_tests (with patch)',
            'test_name': 'foo.bar',
        },],
    }

    response = self.call_api('GetCQFlakes', body=request)
    self.assertEqual(200, response.status_int)
    self.assertDictEqual({}, response.json_body)

  # This test tests that the same test names of different step names will be
  # treated as two different tests.
  @mock.patch.object(
      time_util, 'GetUTCNow', return_value=datetime.datetime(2019, 12, 10))
  def testGetCQFlakesMultipleSteps(self, _):
    flake_issue = _CreateFlakeIssue()
    flake = _CreateFlake()
    flake.flake_issue_key = flake_issue.key
    flake_issue.put()
    flake.put()

    _CreateFlakeOccurrence(
        build_id=987,
        time_happened=time_util.GetDatetimeBeforeNow(hours=1),
        gerrit_cl_id=1234,
        step_ui_name='browser_tests (with patch)').put()
    _CreateFlakeOccurrence(
        build_id=986,
        time_happened=time_util.GetDatetimeBeforeNow(hours=14),
        gerrit_cl_id=1235,
        step_ui_name='browser_tests (with patch)').put()
    _CreateFlakeOccurrence(
        build_id=985,
        time_happened=time_util.GetDatetimeBeforeNow(hours=15),
        gerrit_cl_id=1236,
        step_ui_name='non_viz_browser_tests (with patch)').put()

    request = {
        'project':
            'chromium',
        'bucket':
            'try',
        'builder':
            'linux-rel',
        'tests': [
            {
                'step_ui_name': 'browser_tests (with patch)',
                'test_name': 'foo.bar',
            },
            {
                'step_ui_name': 'non_viz_browser_tests (with patch)',
                'test_name': 'foo.bar',
            },
        ],
    }

    response = self.call_api('GetCQFlakes', body=request)
    self.assertEqual(200, response.status_int)
    self.assertDictEqual({}, response.json_body)
