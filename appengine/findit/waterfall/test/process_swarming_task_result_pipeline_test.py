# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import datetime

from testing_utils import testing

from model import wf_analysis_status
from model.wf_swarming_task import WfSwarmingTask
from waterfall import process_swarming_task_result_pipeline
from waterfall import swarming_util
from waterfall.process_swarming_task_result_pipeline import (
    ProcessSwarmingTaskResultPipeline)


_ISOLATED_SERVER = 'https://isolateserver.appspot.com'
_ISOLATED_STORAGE_URL = 'isolateserver.storage.googleapis.com'
_SAMPLE_FAILURE_LOG = {
    'per_iteration_data': [
        {
            'TestSuite1.test1': [
                {
                    'status': 'SUCCESS',
                    'other_info': 'N/A'
                }
            ],
            'TestSuite1.test2': [
                {
                    'status': 'FAILURE',
                    'other_info': 'N/A'
                },
                {
                    'status': 'FAILURE',
                    'other_info': 'N/A'
                },
                {
                    'status': 'SUCCESS',
                    'other_info': 'N/A'
                }
            ],
            'TestSuite1.test3': [
                {
                    'status': 'FAILURE',
                    'other_info': 'N/A'
                },
                {
                    'status': 'FAILURE',
                    'other_info': 'N/A'
                },
                {
                    'status': 'FAILURE',
                    'other_info': 'N/A'
                }
            ]
        },
        {
            'TestSuite1.test1': [
                {
                    'status': 'SUCCESS',
                    'other_info': 'N/A'
                }
            ],
            'TestSuite1.test2': [
                {
                    'status': 'SUCCESS',
                    'other_info': 'N/A'
                }
            ],
            'TestSuite1.test3': [
                {
                    'status': 'SUCCESS',
                    'other_info': 'N/A'
                }
            ]
        }
    ]
}
_EXPECTED_TESTS_STATUESE = {
    'TestSuite1.test1': {
        'total_run': 2,
        'SUCCESS': 2
    },
    'TestSuite1.test2': {
        'total_run': 4,
        'SUCCESS': 2,
        'FAILURE': 2
    },
    'TestSuite1.test3': {
        'total_run': 4,
        'SUCCESS': 1,
        'FAILURE': 3
    }
}


class ProcessSwarmingTaskResultPipelineTest(testing.AppengineTestCase):

  def _MockedGetSwarmingTaskResultById(self, task_id, _):
    swarming_task_results = {
        'task_id1': {
            'state': 'COMPLETED',
            'outputs_ref': {
                'isolatedserver': _ISOLATED_SERVER,
                'namespace': 'default-gzip',
                'isolated': 'shard1_isolated'
            },
            'created_ts': '2016-02-10T18:32:06.538220',
            'started_ts': '2016-02-10T18:32:09.090550',
            'completed_ts': '2016-02-10T18:33:09'
        },
        'task_id2': {
            'state': 'TIMED_OUT',
            'outputs_ref': None
        }
    }

    mocked_result = swarming_task_results[task_id]
    return mocked_result

  def setUp(self):
    super(ProcessSwarmingTaskResultPipelineTest, self).setUp()
    self.master_name = 'm'
    self.builder_name = 'b'
    self.build_number = 121
    self.step_name = 'abc_tests'
    self.mock(swarming_util, 'GetSwarmingTaskResultById',
              self._MockedGetSwarmingTaskResultById)

  def testCheckTestsRunStatusesNoOutPutJson(self):
    tests_statuses = (
        process_swarming_task_result_pipeline._CheckTestsRunStatuses(None))
    self.assertEqual({}, tests_statuses)

  def testCheckTestsRunStatuses(self):
    tests_statuses = (
        process_swarming_task_result_pipeline._CheckTestsRunStatuses(
            _SAMPLE_FAILURE_LOG))
    self.assertEqual(_EXPECTED_TESTS_STATUESE, tests_statuses)

  def _MockedGetSwarmingTaskFailureLog(self, *_):
    return _SAMPLE_FAILURE_LOG

  def testProcessSwarmingTaskResultPipeline(self):
    task_id = 'task_id1'

    self.mock(swarming_util, 'GetSwarmingTaskFailureLog',
              self._MockedGetSwarmingTaskFailureLog)

    WfSwarmingTask.Create(
        self.master_name, self.builder_name,
        self.build_number, self.step_name).put()

    pipeline = ProcessSwarmingTaskResultPipeline()
    step_name, tests_statuses = pipeline.run(
        self.master_name, self.builder_name,
        self.build_number, self.step_name, task_id)

    self.assertEqual(_EXPECTED_TESTS_STATUESE, tests_statuses)

    task = WfSwarmingTask.Get(
        self.master_name, self.builder_name,self.build_number, self.step_name)

    self.assertEqual(self.step_name, step_name)
    self.assertEqual(wf_analysis_status.ANALYZED, task.status)
    self.assertEqual(_EXPECTED_TESTS_STATUESE, task.tests_statuses)
    self.assertEqual(datetime.datetime(2016, 2, 10, 18, 32, 6, 538220),
                     task.created_time)
    self.assertEqual(datetime.datetime(2016, 2, 10, 18, 32, 9, 90550),
                     task.started_time)
    self.assertEqual(datetime.datetime(2016, 2, 10, 18, 33, 9),
                     task.completed_time)


  def testProcessSwarmingTaskResultPipelineTaskNotRunning(self):
    task_id = 'task_id2'


    WfSwarmingTask.Create(
        self.master_name, self.builder_name,
        self.build_number, self.step_name).put()

    pipeline = ProcessSwarmingTaskResultPipeline()
    step_name, tests_statuses = pipeline.run(
        self.master_name, self.builder_name,
        self.build_number, self.step_name, task_id)

    self.assertEqual({}, tests_statuses)

    task = WfSwarmingTask.Get(
        self.master_name, self.builder_name,self.build_number, self.step_name)

    self.assertEqual(wf_analysis_status.ERROR, task.status)
    self.assertEqual({}, task.tests_statuses)
    self.assertEqual(self.step_name, step_name)

  def testProcessSwarmingTaskResultPipelineTaskTimeOut(self):
    task_id = 'task_id2'
    self.mock(ProcessSwarmingTaskResultPipeline, 'TIMEOUT_HOURS', -1)

    WfSwarmingTask.Create(
        self.master_name, self.builder_name,
        self.build_number, self.step_name).put()

    pipeline = ProcessSwarmingTaskResultPipeline()
    step_name, tests_statuses = pipeline.run(
        self.master_name, self.builder_name,
        self.build_number, self.step_name, task_id)

    self.assertEqual({}, tests_statuses)

    task = WfSwarmingTask.Get(
        self.master_name, self.builder_name,self.build_number, self.step_name)

    self.assertEqual(wf_analysis_status.ERROR, task.status)
    self.assertEqual({}, task.tests_statuses)
    self.assertEqual(self.step_name, step_name)
