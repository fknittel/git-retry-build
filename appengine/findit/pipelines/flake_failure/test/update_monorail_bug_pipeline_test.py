# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import mock

from gae_libs import appengine_util
from gae_libs import pipelines
from gae_libs.pipeline_wrapper import pipeline_handlers
from model.flake.analysis.master_flake_analysis import DataPoint
from model.flake.analysis.master_flake_analysis import MasterFlakeAnalysis
from monorail_api import Issue
from pipelines.flake_failure.update_monorail_bug_pipeline import (
    UpdateMonorailBugInput)
from pipelines.flake_failure.update_monorail_bug_pipeline import (
    UpdateMonorailBugPipeline)
from services import issue_tracking_service
from services.flake_failure import flake_report_util
from waterfall.test.wf_testcase import WaterfallTestCase


class UpdateMonorailPipelineTestShouldNotUpdate(WaterfallTestCase):
  app_module = pipeline_handlers._APP

  @mock.patch.object(
      flake_report_util, 'ShouldUpdateBugForAnalysis', return_value=False)
  def testUpdateMonorailBugPipelineShouldNotUpdateBug(self, _):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 123, 's', 't')
    analysis.Save()

    update_monorail_bug_input = UpdateMonorailBugInput(
        analysis_urlsafe_key=analysis.key.urlsafe())

    pipeline_job = UpdateMonorailBugPipeline(update_monorail_bug_input)
    pipeline_job.start()
    self.execute_queued_tasks()

    pipeline_job = pipelines.pipeline.Pipeline.from_id(pipeline_job.pipeline_id)
    self.assertFalse(pipeline_job.outputs.default.value)

  @mock.patch.object(
      flake_report_util, 'ShouldUpdateBugForAnalysis', return_value=True)
  @mock.patch.object(
      issue_tracking_service,
      'GetMergedDestinationIssueForId',
      return_value=None)
  @mock.patch(
      'pipelines.flake_failure.update_monorail_bug_pipeline.IssueTrackerAPI')
  def testUpdateMonorailBugPipelineWithCulpritBugNotFound(self, mocked_api, *_):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 123, 's', 't')
    analysis.culprit_urlsafe_key = 'c'
    analysis.Save()

    update_monorail_bug_input = UpdateMonorailBugInput(
        analysis_urlsafe_key=analysis.key.urlsafe())

    pipeline_job = UpdateMonorailBugPipeline(update_monorail_bug_input)
    pipeline_job.start()
    self.execute_queued_tasks()

    pipeline_job = pipelines.pipeline.Pipeline.from_id(pipeline_job.pipeline_id)
    self.assertFalse(pipeline_job.outputs.default.value)
    mocked_api.assert_called_once()

  @mock.patch.object(
      flake_report_util, 'ShouldUpdateBugForAnalysis', return_value=True)
  @mock.patch.object(
      issue_tracking_service,
      'GetMergedDestinationIssueForId',
      return_value=None)
  @mock.patch(
      'pipelines.flake_failure.update_monorail_bug_pipeline.IssueTrackerAPI')
  def testUpdateMonorailBugPipelineWithoutCulpritBugNotFound(
      self, mocked_api, *_):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 123, 's', 't')
    analysis.Save()

    update_monorail_bug_input = UpdateMonorailBugInput(
        analysis_urlsafe_key=analysis.key.urlsafe())

    pipeline_job = UpdateMonorailBugPipeline(update_monorail_bug_input)
    pipeline_job.start()
    self.execute_queued_tasks()

    pipeline_job = pipelines.pipeline.Pipeline.from_id(pipeline_job.pipeline_id)
    self.assertFalse(pipeline_job.outputs.default.value)
    mocked_api.assert_called_once()

  @mock.patch.object(appengine_util, 'IsStaging', return_value=False)
  @mock.patch.object(
      flake_report_util, 'ShouldUpdateBugForAnalysis', return_value=True)
  @mock.patch.object(
      flake_report_util.FlakeAnalysisIssueGenerator,
      'GetComment',
      return_value='comment')
  @mock.patch.object(issue_tracking_service, 'UpdateBug')
  @mock.patch(
      'pipelines.flake_failure.update_monorail_bug_pipeline.IssueTrackerAPI')
  @mock.patch.object(issue_tracking_service, 'GetMergedDestinationIssueForId')
  def testUpdateMonorailBugPipelineWithCulprit(
      self, mock_get_merged_issue, issue_tracker, mock_update_bug, *_):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 1, 's', 't')
    analysis.bug_id = 123
    analysis.culprit_urlsafe_key = 'c'
    analysis.data_points = [DataPoint(), DataPoint(), DataPoint()]
    analysis.Save()

    mock_issue = Issue({})
    mocked_instance = mock.Mock()
    mocked_instance.getIssue.return_value = mock_issue
    issue_tracker.return_value = mocked_instance

    mock_get_merged_issue.return_value = mock_issue

    update_monorail_bug_input = UpdateMonorailBugInput(
        analysis_urlsafe_key=analysis.key.urlsafe())
    pipeline_job = UpdateMonorailBugPipeline(update_monorail_bug_input)
    pipeline_job.start()
    self.execute_queued_tasks()

    pipeline_job = pipelines.pipeline.Pipeline.from_id(pipeline_job.pipeline_id)
    self.assertTrue(pipeline_job.outputs.default.value)
    self.assertTrue(analysis.has_commented_on_bug)
    self.assertTrue(mock_update_bug.called)

  @mock.patch.object(appengine_util, 'IsStaging', return_value=False)
  @mock.patch.object(
      flake_report_util, 'ShouldUpdateBugForAnalysis', return_value=True)
  @mock.patch.object(issue_tracking_service, 'UpdateBug')
  @mock.patch(
      'pipelines.flake_failure.update_monorail_bug_pipeline.IssueTrackerAPI')
  @mock.patch.object(issue_tracking_service, 'GetMergedDestinationIssueForId')
  def testUpdateMonorailBugPipelineNoCulprit(
      self, mock_get_merged_issue, issue_tracker, mock_update_bug, *_):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 1, 's', 't')
    analysis.bug_id = 123
    analysis.data_points = [DataPoint(), DataPoint(), DataPoint()]
    analysis.Save()

    mock_issue = Issue({})
    mocked_instance = mock.Mock()
    mocked_instance.getIssue.return_value = mock_issue
    issue_tracker.return_value = mocked_instance

    mock_get_merged_issue.return_value = mock_issue

    update_monorail_bug_input = UpdateMonorailBugInput(
        analysis_urlsafe_key=analysis.key.urlsafe())
    pipeline_job = UpdateMonorailBugPipeline(update_monorail_bug_input)
    pipeline_job.start()
    self.execute_queued_tasks()

    pipeline_job = pipelines.pipeline.Pipeline.from_id(pipeline_job.pipeline_id)
    self.assertTrue(pipeline_job.outputs.default.value)
    self.assertTrue(mock_update_bug.called)
