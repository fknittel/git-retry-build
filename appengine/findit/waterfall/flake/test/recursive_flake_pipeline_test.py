# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from datetime import datetime
import mock

from common import constants
from gae_libs.pipeline_wrapper import pipeline_handlers
from libs import analysis_status
from libs import time_util
from model.flake.flake_swarming_task import FlakeSwarmingTask
from model.flake.master_flake_analysis import DataPoint
from model.flake.master_flake_analysis import MasterFlakeAnalysis
from model.wf_swarming_task import WfSwarmingTask

from waterfall.flake import confidence
from waterfall import swarming_util
from waterfall.flake import flake_constants
from waterfall.flake import lookback_algorithm
from waterfall.flake import recursive_flake_pipeline

from waterfall.flake.initialize_flake_try_job_pipeline import (
    InitializeFlakeTryJobPipeline)
from waterfall.flake.next_build_number_pipeline import NextBuildNumberPipeline
from waterfall.flake.recursive_flake_pipeline import RecursiveFlakePipeline
from waterfall.flake.save_last_attempted_swarming_task_id_pipeline import (
    SaveLastAttemptedSwarmingTaskIdPipeline)
from waterfall.flake.update_flake_analysis_data_points_pipeline import (
    UpdateFlakeAnalysisDataPointsPipeline)
from waterfall.flake.update_flake_bug_pipeline import UpdateFlakeBugPipeline
from waterfall.test import wf_testcase
from waterfall.test.wf_testcase import DEFAULT_CONFIG_DATA

_DEFAULT_CACHE_NAME = swarming_util.GetCacheName('pm', 'pb')


class MOCK_INFO(object):
  parent_buildername = 'pb'
  parent_mastername = 'pm'


class RecursiveFlakePipelineTest(wf_testcase.WaterfallTestCase):
  app_module = pipeline_handlers._APP

  @mock.patch.object(swarming_util, 'BotsAvailableForTask', return_value=True)
  @mock.patch.object(
      recursive_flake_pipeline, '_IsFinished', side_effect=[False, True])
  def testRecursiveFlakePipeline(self, *_):
    master_name = 'm'
    builder_name = 'b'
    master_build_number = 100
    build_number = 100
    run_build_number = 100
    step_name = 's'
    test_name = 't'
    queue_name = constants.DEFAULT_QUEUE
    task_id = 'task_id'

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.status = analysis_status.PENDING
    analysis.algorithm_parameters = DEFAULT_CONFIG_DATA['check_flake_settings']
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.status = analysis_status.COMPLETED
    task.put()

    self.MockPipeline(
        recursive_flake_pipeline.TriggerFlakeSwarmingTaskPipeline,
        'task_id',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, [test_name],
            100, 3 * 60 * 60
        ],
        expected_kwargs={'force': False})

    self.MockPipeline(
        SaveLastAttemptedSwarmingTaskIdPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), task_id, run_build_number],
        expected_kwargs={})

    self.MockPipeline(
        recursive_flake_pipeline.ProcessFlakeSwarmingTaskResultPipeline,
        'test_result_future',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, task_id,
            master_build_number, test_name, analysis.version_number
        ],
        expected_kwargs={})

    self.MockPipeline(
        UpdateFlakeAnalysisDataPointsPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), run_build_number],
        expected_kwargs={})

    self.MockPipeline(
        NextBuildNumberPipeline,
        100,
        expected_args=[analysis.key.urlsafe(), build_number, None, None, None])

    pipeline_job = RecursiveFlakePipeline(
        analysis.key.urlsafe(),
        build_number,
        None,
        None,
        None,
        use_nearby_neighbor=False)
    pipeline_job.start(queue_name=queue_name)
    self.execute_queued_tasks()

  @mock.patch.object(swarming_util, 'BotsAvailableForTask', return_value=True)
  @mock.patch.object(
      recursive_flake_pipeline, '_IsFinished', side_effect=[False, True])
  def testRecursiveFlakePipelineWithUserInput(self, *_):
    master_name = 'm'
    builder_name = 'b'
    master_build_number = 100
    run_build_number = 90
    lower_bound_build_number = 50
    upper_bound_build_number = 90
    iterations_to_rerun = 150
    step_name = 's'
    test_name = 't'
    queue_name = constants.DEFAULT_QUEUE
    task_id = 'task_id'

    analysis = MasterFlakeAnalysis.Create(
        master_name, builder_name, master_build_number, step_name, test_name)
    analysis.status = analysis_status.PENDING
    analysis.algorithm_parameters = DEFAULT_CONFIG_DATA['check_flake_settings']
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, run_build_number,
                                    step_name, test_name)
    task.status = analysis_status.ERROR
    task.put()

    self.MockPipeline(
        recursive_flake_pipeline.TriggerFlakeSwarmingTaskPipeline,
        'task_id',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, [test_name],
            iterations_to_rerun, 3 * 60 * 60
        ],
        expected_kwargs={'force': False})

    self.MockPipeline(
        SaveLastAttemptedSwarmingTaskIdPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), task_id, run_build_number],
        expected_kwargs={})

    self.MockPipeline(
        recursive_flake_pipeline.ProcessFlakeSwarmingTaskResultPipeline,
        'test_result_future',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, task_id,
            master_build_number, test_name, analysis.version_number
        ],
        expected_kwargs={})

    self.MockPipeline(
        UpdateFlakeAnalysisDataPointsPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), run_build_number],
        expected_kwargs={})

    self.MockPipeline(
        NextBuildNumberPipeline,
        100,
        expected_args=[
            analysis.key.urlsafe(), run_build_number, lower_bound_build_number,
            upper_bound_build_number, iterations_to_rerun
        ])

    pipeline_job = RecursiveFlakePipeline(
        analysis.key.urlsafe(),
        run_build_number,
        lower_bound_build_number,
        upper_bound_build_number,
        iterations_to_rerun,
        use_nearby_neighbor=False)
    pipeline_job.start(queue_name=queue_name)
    self.execute_queued_tasks()

  @mock.patch.object(swarming_util, 'BotsAvailableForTask', return_value=True)
  @mock.patch.object(
      recursive_flake_pipeline, '_IsFinished', side_effect=[False, True])
  def testRecursiveFlakePipelineWithUpperLowerBounds(self, *_):
    master_name = 'm'
    builder_name = 'b'
    master_build_number = 100
    lower_bound_build_number = 50
    upper_bound_build_number = 90
    run_build_number = 51
    step_name = 's'
    test_name = 't'
    queue_name = constants.DEFAULT_QUEUE
    task_id = 'task_id'

    analysis = MasterFlakeAnalysis.Create(
        master_name, builder_name, master_build_number, step_name, test_name)
    analysis.status = analysis_status.PENDING
    analysis.algorithm_parameters = DEFAULT_CONFIG_DATA['check_flake_settings']
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, run_build_number,
                                    step_name, test_name)
    task.status = analysis_status.ERROR
    task.put()

    self.MockPipeline(
        recursive_flake_pipeline.TriggerFlakeSwarmingTaskPipeline,
        'task_id',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, [test_name],
            master_build_number, 3 * 60 * 60
        ],
        expected_kwargs={'force': False})

    self.MockPipeline(
        SaveLastAttemptedSwarmingTaskIdPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), task_id, run_build_number],
        expected_kwargs={})

    self.MockPipeline(
        recursive_flake_pipeline.ProcessFlakeSwarmingTaskResultPipeline,
        'test_result_future',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, task_id,
            master_build_number, test_name, analysis.version_number
        ],
        expected_kwargs={})

    self.MockPipeline(
        UpdateFlakeAnalysisDataPointsPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), run_build_number])

    self.MockPipeline(
        NextBuildNumberPipeline,
        100,
        expected_args=[
            analysis.key.urlsafe(), run_build_number, lower_bound_build_number,
            upper_bound_build_number, None
        ])

    pipeline_job = RecursiveFlakePipeline(
        analysis.key.urlsafe(),
        run_build_number,
        50,
        90,
        None,
        use_nearby_neighbor=False)
    pipeline_job.start(queue_name=queue_name)
    self.execute_queued_tasks()

  @mock.patch.object(swarming_util, 'BotsAvailableForTask', return_value=True)
  @mock.patch.object(
      recursive_flake_pipeline, '_IsFinished', side_effect=[False, True])
  def testRecursiveFlakePipelineWithForceFlag(self, *_):
    master_name = 'm'
    builder_name = 'b'
    master_build_number = 100
    lower_bound_build_number = None
    upper_bound_build_number = None
    run_build_number = 51
    step_name = 's'
    test_name = 't'
    queue_name = constants.DEFAULT_QUEUE
    task_id = 'task_id'

    analysis = MasterFlakeAnalysis.Create(
        master_name, builder_name, master_build_number, step_name, test_name)
    analysis.status = analysis_status.PENDING
    analysis.algorithm_parameters = DEFAULT_CONFIG_DATA['check_flake_settings']
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, run_build_number,
                                    step_name, test_name)
    task.status = analysis_status.ERROR
    task.put()

    self.MockPipeline(
        recursive_flake_pipeline.TriggerFlakeSwarmingTaskPipeline,
        'task_id',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, [test_name],
            master_build_number, 3 * 60 * 60
        ],
        expected_kwargs={'force': True})

    self.MockPipeline(
        SaveLastAttemptedSwarmingTaskIdPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), task_id, run_build_number],
        expected_kwargs={})

    self.MockPipeline(
        recursive_flake_pipeline.ProcessFlakeSwarmingTaskResultPipeline,
        'test_result_future',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, task_id,
            master_build_number, test_name, analysis.version_number
        ],
        expected_kwargs={})

    self.MockPipeline(
        UpdateFlakeAnalysisDataPointsPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), run_build_number],
        expected_kwargs={})

    self.MockPipeline(
        NextBuildNumberPipeline,
        100,
        expected_args=[
            analysis.key.urlsafe(), run_build_number, lower_bound_build_number,
            upper_bound_build_number, None
        ])

    pipeline_job = RecursiveFlakePipeline(
        analysis.key.urlsafe(),
        run_build_number,
        lower_bound_build_number,
        upper_bound_build_number,
        None,
        force=True)
    pipeline_job.start(queue_name=queue_name)
    self.execute_queued_tasks()

  @mock.patch.object(flake_constants, 'BASE_COUNT_DOWN_SECONDS', 0)
  @mock.patch.object(
      recursive_flake_pipeline, '_IsFinished', side_effect=[False, True])
  @mock.patch.object(swarming_util, 'BotsAvailableForTask')
  def testTryLaterIfNoAvailableBots(self, mock_fn, *_):
    mock_fn.side_effect = [False, True]

    master_name = 'm'
    builder_name = 'b'
    master_build_number = 100
    build_number = 100
    run_build_number = 100
    step_name = 's'
    test_name = 't'
    queue_name = constants.DEFAULT_QUEUE
    task_id = 'task_id'

    analysis = MasterFlakeAnalysis.Create(
        master_name, builder_name, master_build_number, step_name, test_name)
    analysis.status = analysis_status.PENDING
    analysis.algorithm_parameters = DEFAULT_CONFIG_DATA['check_flake_settings']
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, run_build_number,
                                    step_name, test_name)
    task.status = analysis_status.ERROR
    task.put()

    self.MockPipeline(
        recursive_flake_pipeline.TriggerFlakeSwarmingTaskPipeline,
        'task_id',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, [test_name]
        ],
        expected_kwargs={})

    self.MockPipeline(
        recursive_flake_pipeline.ProcessFlakeSwarmingTaskResultPipeline,
        'test_result_future',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, task_id,
            master_build_number, test_name, analysis.version_number
        ],
        expected_kwargs={})

    self.MockPipeline(
        NextBuildNumberPipeline,
        100,
        expected_args=[analysis.key.urlsafe(), build_number, None, None, None])

    pipeline_job = RecursiveFlakePipeline(
        analysis.key.urlsafe(),
        build_number,
        None,
        None,
        None,
        use_nearby_neighbor=False)

    pipeline_job.start(queue_name=queue_name)
    self.execute_queued_tasks()

  @mock.patch.object(recursive_flake_pipeline, '_IsFinished', return_value=True)
  def testRecursiveFlakePipelineInitializeFlakeTryJobPipeline(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 123
    step_name = 's'
    test_name = 't'

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.status = analysis_status.PENDING
    analysis.algorithm_parameters = DEFAULT_CONFIG_DATA['check_flake_settings']
    analysis.data_points = [
        DataPoint.Create(pass_rate=0.9, build_number=build_number),
        DataPoint.Create(pass_rate=1.0, build_number=build_number - 1)
    ]
    analysis.suspected_flake_build_number = build_number
    analysis.confidence_in_suspected_build = 0.7
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.status = analysis_status.ERROR
    task.put()

    self.MockPipeline(
        InitializeFlakeTryJobPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), None, False],
        expected_kwargs={})
    self.MockPipeline(
        UpdateFlakeBugPipeline,
        '',
        expected_args=[analysis.key.urlsafe()],
        expected_kwargs={})

    pipeline_job = NextBuildNumberPipeline(analysis.key.urlsafe(), build_number,
                                           None, None, None)
    pipeline_job.start(queue_name=constants.DEFAULT_QUEUE)
    self.execute_queued_tasks()

  @mock.patch.object(RecursiveFlakePipeline, 'was_aborted', return_value=True)
  @mock.patch.object(
      recursive_flake_pipeline, '_IsFinished', side_effect=[False, True])
  def testRecursiveFlakePipelineAborted(self, *_):
    master_name = 'm'
    builder_name = 'b'
    master_build_number = 100
    build_number = 100
    step_name = 's'
    test_name = 't'

    analysis = MasterFlakeAnalysis.Create(
        master_name, builder_name, master_build_number, step_name, test_name)
    analysis.status = analysis_status.PENDING
    analysis.put()

    pipeline_job = RecursiveFlakePipeline(analysis.key.urlsafe(), build_number,
                                          None, None, None)
    pipeline_job._LogUnexpectedAbort()

    expected_error = {
        'error': 'RecursiveFlakePipeline was aborted unexpectedly',
        'message': 'RecursiveFlakePipeline was aborted unexpectedly'
    }

    self.assertEqual(analysis_status.ERROR, analysis.status)
    self.assertEqual(expected_error, analysis.error)

  @mock.patch.object(RecursiveFlakePipeline, 'was_aborted', return_value=True)
  @mock.patch.object(
      recursive_flake_pipeline, '_IsFinished', side_effect=[False, True])
  def testRecursiveFlakePipelineAbortedNotUpdateCompletedAnalysis(self, *_):
    master_name = 'm'
    builder_name = 'b'
    master_build_number = 100
    build_number = 100
    step_name = 's'
    test_name = 't'

    analysis = MasterFlakeAnalysis.Create(
        master_name, builder_name, master_build_number, step_name, test_name)
    analysis.status = analysis_status.COMPLETED
    analysis.Save()

    pipeline_job = RecursiveFlakePipeline(analysis.key.urlsafe(), build_number,
                                          None, None, None)
    pipeline_job._LogUnexpectedAbort()

    self.assertEqual(analysis_status.COMPLETED, analysis.status)

  @mock.patch.object(recursive_flake_pipeline, '_IsFinished', return_value=True)
  def testRecursiveFlakePipelineFinishes(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 100
    step_name = 's'
    test_name = 't'

    lower_bound = 1
    upper_bound = 10
    user_range = recursive_flake_pipeline._UserSpecifiedRange(
        lower_bound, upper_bound)
    iterations = 100

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.status = analysis_status.COMPLETED
    analysis.algorithm_parameters = DEFAULT_CONFIG_DATA['check_flake_settings']
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.status = analysis_status.COMPLETED
    task.put()

    self.MockPipeline(
        UpdateFlakeBugPipeline,
        '',
        expected_args=[analysis.key.urlsafe()],
        expected_kwargs={})

    self.MockPipeline(
        InitializeFlakeTryJobPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), iterations, user_range],
        expected_kwargs={})

    pipeline = RecursiveFlakePipeline(analysis.key.urlsafe(), build_number,
                                      lower_bound, upper_bound, iterations)
    pipeline.start(queue_name=constants.DEFAULT_QUEUE)
    self.execute_queued_tasks()

  @mock.patch.object(swarming_util, 'GetETAToStartAnalysis', return_value=None)
  @mock.patch.object(swarming_util, 'BotsAvailableForTask', return_value=False)
  @mock.patch.object(
      recursive_flake_pipeline, '_IsFinished', side_effect=[False, True])
  @mock.patch.object(FlakeSwarmingTask, 'Get')
  def testRetriesExceedMax(self, mock_flake_swarming_task, *_):
    master_name = 'm'
    builder_name = 'b'
    master_build_number = 100
    build_number = 100
    run_build_number = 100
    step_name = 's'
    test_name = 't'
    queue_name = constants.DEFAULT_QUEUE
    task_id = 'task_id'

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.status = analysis_status.PENDING
    analysis.algorithm_parameters = DEFAULT_CONFIG_DATA['check_flake_settings']
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.status = analysis_status.ERROR
    task.put()
    mock_flake_swarming_task.return_value = task

    self.MockPipeline(
        recursive_flake_pipeline.TriggerFlakeSwarmingTaskPipeline,
        'task_id',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, [test_name],
            100, 3 * 60 * 60
        ],
        expected_kwargs={'force': False})

    self.MockPipeline(
        SaveLastAttemptedSwarmingTaskIdPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), task_id, build_number],
        expected_kwargs={})

    self.MockPipeline(
        recursive_flake_pipeline.ProcessFlakeSwarmingTaskResultPipeline,
        'test_result_future',
        expected_args=[
            master_name, builder_name, run_build_number, step_name, task_id,
            master_build_number, test_name, analysis.version_number
        ],
        expected_kwargs={})

    self.MockPipeline(
        UpdateFlakeAnalysisDataPointsPipeline,
        '',
        expected_args=[analysis.key.urlsafe(), build_number],
        expected_kwargs={})

    self.MockPipeline(
        NextBuildNumberPipeline,
        100,
        expected_args=[analysis.key.urlsafe(), build_number, None, None, None])

    pipeline_job = RecursiveFlakePipeline(
        analysis.key.urlsafe(),
        build_number,
        None,
        None,
        None,
        step_metadata=None,
        manually_triggered=False,
        use_nearby_neighbor=False,
        retries=5)

    pipeline_job.start(queue_name=queue_name)
    self.execute_queued_tasks()

  #######################################
  #      Function unit tests.           #
  #######################################

  @mock.patch.object(
      recursive_flake_pipeline.confidence,
      'SteppinessForBuild',
      return_value=0.6)
  def testGetBuildConfidenceScore(self, _):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 124, 's', 't')
    analysis.data_points = [DataPoint.Create(pass_rate=0.7, build_number=123)]
    self.assertIsNone(
        recursive_flake_pipeline._GetBuildConfidenceScore(analysis, None, []))
    self.assertEqual(0.6,
                     recursive_flake_pipeline._GetBuildConfidenceScore(
                         analysis, 123, analysis.data_points))

  def testGetBuildConfidenceScoreIntroducedNewFlakyTest(self):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 124, 's', 't')
    analysis.data_points = [
        DataPoint.Create(pass_rate=0.7, build_number=123),
        DataPoint.Create(
            pass_rate=flake_constants.PASS_RATE_TEST_NOT_FOUND,
            build_number=122)
    ]
    self.assertEqual(1.0,
                     recursive_flake_pipeline._GetBuildConfidenceScore(
                         analysis, 123, analysis.data_points))

  def testGetBestBuildNumberToRunWithNearbyNeighborRunning(self):
    master_name = 'm'
    builder_name = 'b'
    preferred_run_build_number = 1000
    cached_build_number = 997
    step_name = 's'
    test_name = 't'
    number_of_iterations = 100
    step_size = 10

    task = FlakeSwarmingTask.Create(master_name, builder_name,
                                    cached_build_number, step_name, test_name)
    task.status = analysis_status.RUNNING
    task.tries = number_of_iterations
    task.put()

    self.assertEqual(
        recursive_flake_pipeline._GetBestBuildNumberToRun(
            master_name, builder_name, preferred_run_build_number, step_name,
            test_name, None, None, step_size, number_of_iterations),
        cached_build_number)

  def testGetBestBuildNumberToRunWithNearbyNeighborCompleted(self):
    # Completed build should take precedence over running build, even if it's
    # farther away.
    master_name = 'm'
    builder_name = 'b'
    preferred_run_build_number = 1000
    running_cached_build_number = 997
    completed_cached_build_number = 996
    step_name = 's'
    test_name = 't'
    number_of_iterations = 100
    step_size = 10

    task = FlakeSwarmingTask.Create(master_name, builder_name,
                                    completed_cached_build_number, step_name,
                                    test_name)
    task.status = analysis_status.COMPLETED
    task.tries = number_of_iterations
    task.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name,
                                    running_cached_build_number, step_name,
                                    test_name)
    task.status = analysis_status.RUNNING
    task.tries = number_of_iterations
    task.put()

    self.assertEqual(
        recursive_flake_pipeline._GetBestBuildNumberToRun(
            master_name, builder_name, preferred_run_build_number, step_name,
            test_name, None, None, step_size, number_of_iterations),
        completed_cached_build_number)

  def testGetBestBuildNumberToRunWithMultipleInProgress(self):
    # Completed builds should take precedence over running build, even if it's
    # farther away.
    master_name = 'm'
    builder_name = 'b'
    preferred_run_build_number = 1000
    running_cached_build_number_1 = 997
    running_cached_build_number_2 = 996
    step_name = 's'
    test_name = 't'
    number_of_iterations = 100
    step_size = 10

    task = FlakeSwarmingTask.Create(master_name, builder_name,
                                    running_cached_build_number_1, step_name,
                                    test_name)
    task.status = analysis_status.RUNNING
    task.tries = number_of_iterations
    task.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name,
                                    running_cached_build_number_2, step_name,
                                    test_name)
    task.status = analysis_status.RUNNING
    task.tries = number_of_iterations
    task.put()

    self.assertEqual(
        recursive_flake_pipeline._GetBestBuildNumberToRun(
            master_name, builder_name, preferred_run_build_number, step_name,
            test_name, None, None, step_size, number_of_iterations),
        running_cached_build_number_1)

  def testGetBestBuildNumberToRunPendingAndRunning(self):
    # Running builds should take precedence over pending builds.
    master_name = 'm'
    builder_name = 'b'
    preferred_run_build_number = 1000
    running_cached_build_number_1 = 997
    running_cached_build_number_2 = 996
    step_name = 's'
    test_name = 't'
    number_of_iterations = 100
    step_size = 10

    task = FlakeSwarmingTask.Create(master_name, builder_name,
                                    running_cached_build_number_1, step_name,
                                    test_name)
    task.status = analysis_status.PENDING
    task.tries = number_of_iterations
    task.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name,
                                    running_cached_build_number_2, step_name,
                                    test_name)
    task.status = analysis_status.RUNNING
    task.tries = number_of_iterations
    task.put()

    self.assertEqual(
        recursive_flake_pipeline._GetBestBuildNumberToRun(
            master_name, builder_name, preferred_run_build_number, step_name,
            test_name, None, None, step_size, number_of_iterations),
        running_cached_build_number_2)

  def testGetListOfNearbyBuildNumbers(self):
    self.assertEqual([1],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         1, None, None, 0))
    self.assertEqual([1],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         1, None, None, -1))
    self.assertEqual([1, 0, 2],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         1, None, None, 1))
    self.assertEqual([1, 0, 2, 3],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         1, None, None, 2))
    self.assertEqual([2, 1, 3],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         2, None, None, 1))
    self.assertEqual([100, 99, 101, 98, 102, 97, 103, 96, 104, 95, 105],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         100, None, None, 5))
    self.assertEqual([6, 5, 7, 8, 9],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         6, 5, None, 3))
    self.assertEqual([7, 6, 8, 5, 9, 10],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         7, 5, None, 3))
    self.assertEqual([8, 7, 9, 6, 5],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         8, None, 9, 3))
    self.assertEqual([8, 9],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         8, 8, 9, 3))
    self.assertEqual([9, 8],
                     recursive_flake_pipeline._GetListOfNearbyBuildNumbers(
                         9, 8, 9, 3))

  def testIsSwarmingTaskSufficientNoSwarmingTasks(self):
    self.assertFalse(
        recursive_flake_pipeline._IsSwarmingTaskSufficientForCacheHit(
            None, 100))

  def testIsSwarmingTaskSufficientForCacheHitInsufficientIterations(self):
    desired_iterations = 200
    flake_swarming_task = FlakeSwarmingTask.Create('m', 'b', 12345, 's', 't')
    flake_swarming_task.tries = 100
    flake_swarming_task.status = analysis_status.COMPLETED
    self.assertFalse(
        recursive_flake_pipeline._IsSwarmingTaskSufficientForCacheHit(
            flake_swarming_task, desired_iterations))

  def testIsSwarmingTaskSufficientForCacheHitError(self):
    desired_iterations = 100
    flake_swarming_task = FlakeSwarmingTask.Create('m', 'b', 12345, 's', 't')
    flake_swarming_task.tries = 200
    flake_swarming_task.status = analysis_status.ERROR
    self.assertFalse(
        recursive_flake_pipeline._IsSwarmingTaskSufficientForCacheHit(
            flake_swarming_task, desired_iterations))

  def testIsSwarmingTaskSufficientForCacheHitPending(self):
    desired_iterations = 100
    flake_swarming_task = FlakeSwarmingTask.Create('m', 'b', 12345, 's', 't')
    flake_swarming_task.tries = desired_iterations
    flake_swarming_task.status = analysis_status.PENDING
    self.assertTrue(
        recursive_flake_pipeline._IsSwarmingTaskSufficientForCacheHit(
            flake_swarming_task, desired_iterations))

  def testIsSwarmingTaskSufficientForCacheHitRunning(self):
    desired_iterations = 100
    flake_swarming_task = FlakeSwarmingTask.Create('m', 'b', 12345, 's', 't')
    flake_swarming_task.tries = desired_iterations
    flake_swarming_task.status = analysis_status.RUNNING
    self.assertTrue(
        recursive_flake_pipeline._IsSwarmingTaskSufficientForCacheHit(
            flake_swarming_task, desired_iterations))

  def testIsSwarmingTaskSufficientForCacheHitCompleted(self):
    desired_iterations = 100
    flake_swarming_task = FlakeSwarmingTask.Create('m', 'b', 12345, 's', 't')
    flake_swarming_task.tries = desired_iterations
    flake_swarming_task.status = analysis_status.COMPLETED
    self.assertTrue(
        recursive_flake_pipeline._IsSwarmingTaskSufficientForCacheHit(
            flake_swarming_task, desired_iterations))

  def testGetBestBuildNumberToRunWithStepSizeZero(self):
    self.assertEqual(12345,
                     recursive_flake_pipeline._GetBestBuildNumberToRun(
                         'm', 'b', 12345, 's', 't', None, None, 0, 100))

  def testGetBestBuildNumberToRunWithNoNearbyNeighbors(self):
    self.assertEqual(12345,
                     recursive_flake_pipeline._GetBestBuildNumberToRun(
                         'm', 'b', 12345, 's', 't', None, None, 10, 100))

  @mock.patch.object(
      recursive_flake_pipeline,
      '_CanEstimateExecutionTimeFromReferenceSwarmingTask',
      return_value=False)
  def testGetHardTimeoutSecondsDefault(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 123
    step_name = 's'

    self.UpdateUnitTestConfigSettings(
        config_property='check_flake_settings',
        override_data={'swarming_rerun': {
            'per_iteration_timeout_seconds': 60
        }})
    self.assertEqual(3 * 60 * 60,
                     recursive_flake_pipeline._GetHardTimeoutSeconds(
                         master_name, builder_name, build_number, step_name,
                         100))

  def testGetHardTimeoutSeconds(self):
    master_name = 'm'
    builder_name = 'b'
    build_number = 123
    step_name = 's'
    reference_swarming_task = WfSwarmingTask.Create(master_name, builder_name,
                                                    build_number, step_name)
    reference_swarming_task.completed_time = datetime(2017, 4, 16, 0, 0, 40, 12)
    reference_swarming_task.started_time = datetime(2017, 4, 16, 0, 0, 0, 10)
    reference_swarming_task.tests_statuses = {'1': 1, '2': 1}
    reference_swarming_task.parameters = {'iterations_to_rerun': 2}
    reference_swarming_task.put()
    self.UpdateUnitTestConfigSettings(
        config_property='check_flake_settings',
        override_data={'swarming_rerun': {
            'per_iteration_timeout_seconds': 1
        }})
    timeout = recursive_flake_pipeline._GetHardTimeoutSeconds(
        master_name, builder_name, build_number, step_name, 10)
    self.assertTrue(isinstance(timeout, int))
    self.assertEqual(60 * 60, timeout)

  def testGetEarliestBuildNumber(self):
    algorithm_settings = {'max_build_numbers_to_look_back': 10}

    self.assertEqual(5,
                     recursive_flake_pipeline._GetEarliestBuildNumber(
                         5, 6, algorithm_settings))
    self.assertEqual(0,
                     recursive_flake_pipeline._GetEarliestBuildNumber(
                         None, 5, algorithm_settings))
    self.assertEqual(15,
                     recursive_flake_pipeline._GetEarliestBuildNumber(
                         None, 25, algorithm_settings))

  def testUserSpecifiedRange(self):
    self.assertTrue(recursive_flake_pipeline._UserSpecifiedRange(123, 125))
    self.assertFalse(recursive_flake_pipeline._UserSpecifiedRange(None, 123))
    self.assertFalse(recursive_flake_pipeline._UserSpecifiedRange(123, None))
    self.assertFalse(recursive_flake_pipeline._UserSpecifiedRange(None, None))

  def testUpdateAnalysisResults(self):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 123, 's', 't')
    analysis.last_attempted_swarming_task_id = '12345'
    recursive_flake_pipeline._UpdateAnalysisResults(
        analysis, 100, analysis_status.COMPLETED, None)
    self.assertEqual(analysis.suspected_flake_build_number, 100)
    self.assertIsNone(analysis.last_attempted_swarming_task_id)

  @mock.patch.object(time_util, 'GetUTCNow', return_value=datetime(2017, 6, 27))
  def testUpdateAnalysisResultsWithError(self, _):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 123, 's', 't')
    analysis.last_attempted_swarming_task_id = '12345'
    recursive_flake_pipeline._UpdateAnalysisResults(
        analysis, None, analysis_status.ERROR, {'error': 'error'})
    self.assertEqual(analysis_status.ERROR, analysis.status)
    self.assertEqual(datetime(2017, 6, 27), analysis.end_time)

  def testIsFinished(self):
    self.assertTrue(recursive_flake_pipeline._IsFinished(1, 2, 3, 0))
    self.assertFalse(recursive_flake_pipeline._IsFinished(1, 2, 3, 10))
    self.assertTrue(recursive_flake_pipeline._IsFinished(4, 2, 3, 0))
    self.assertFalse(recursive_flake_pipeline._IsFinished(4, 2, 3, 10))
    self.assertTrue(recursive_flake_pipeline._IsFinished(1, 2, 3, 0))
