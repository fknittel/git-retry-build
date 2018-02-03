# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
import datetime
import mock
import logging

from google.protobuf import timestamp_pb2

from libs import analysis_status
from model import suspected_cl_status
from model.flake.flake_culprit import FlakeCulprit
from model.flake.master_flake_analysis import MasterFlakeAnalysis
from model.flake.master_flake_analysis import DataPoint
from model.proto.gen import findit_pb2
from model.proto.gen.compile_analysis_pb2 import CompileAnalysisCompletionEvent
from model.proto.gen.test_analysis_pb2 import TestAnalysisCompletionEvent
from model.wf_try_job import WfTryJob
from model.wf_suspected_cl import WfSuspectedCL
from model.wf_analysis import WfAnalysis
from services import bigquery_helper
from services import event_reporting
from waterfall.flake import triggering_sources

from waterfall.test import wf_testcase


class EventReportingTest(wf_testcase.WaterfallTestCase):
  @mock.patch.object(bigquery_helper, 'ReportEventsToBigquery')
  @mock.patch.object(event_reporting, 'CreateTestFlakeAnalysisCompletionEvent')
  def testReportTestFlakeAnalysisCompletionEvent(self, mock_create_proto_fn,
      bq_helper_fn):
    proto = TestAnalysisCompletionEvent()
    mock_create_proto_fn.return_value = proto

    analysis = MasterFlakeAnalysis.Create('m', 'b', 123, 's', 't')

    event_reporting.ReportTestFlakeAnalysisCompletionEvent(analysis)
    self.assertTrue(mock_create_proto_fn.called)
    self.assertTrue(bq_helper_fn.called)

  @mock.patch.object(bigquery_helper, 'ReportEventsToBigquery')
  @mock.patch.object(event_reporting,
                     'CreateCompileFailureAnalysisCompletionEvent')
  def testReportCompileFailureAnalysisCompletionEvent(
      self, mock_create_proto_fn, bq_helper_fn):
    proto = CompileAnalysisCompletionEvent()
    mock_create_proto_fn.return_value = proto

    analysis = WfAnalysis.Create('m', 'b', 0)

    event_reporting.ReportCompileFailureAnalysisCompletionEvent(analysis)
    self.assertTrue(mock_create_proto_fn.called)
    self.assertTrue(bq_helper_fn.called)

  def testCreateTestFlakeAnalysisCompletionEventWithNoSuspectedBuild(self):
    """Test reporting event where the Cr has been notified."""
    master = 'master'
    builder = 'builder'
    step = 'step'
    test = 'test'
    build_number = 100

    analysis = MasterFlakeAnalysis.Create(master, builder, build_number, step,
                                          test)
    analysis.data_points = [
      DataPoint.Create(),
      DataPoint.Create(),
      DataPoint.Create()
    ]
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.put()

    event = event_reporting.CreateTestFlakeAnalysisCompletionEvent(analysis)
    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)

    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)

    self.assertEqual(event.analysis_info.outcomes, [findit_pb2.REPRODUCIBLE])
    self.assertTrue(event.flake)

  def testCreateTestFlakeAnalysisCompletionEventWithCrNotification(self):
    """Test reporting event where the Cr has been notified."""
    master = 'master'
    builder = 'builder'
    step = 'step'
    test = 'test'

    build_number = 10
    suspected_build_number = 5

    repo = 'chromium'
    revision = 'revision'
    culprit_commit_position = 2
    culprit = FlakeCulprit.Create(repo, revision, culprit_commit_position)

    suspect_1 = FlakeCulprit.Create(repo, revision, culprit_commit_position - 1)
    suspect_1.put()
    suspect_2 = FlakeCulprit.Create(repo, revision, culprit_commit_position)
    suspect_2.put()
    suspect_3 = FlakeCulprit.Create(repo, revision, culprit_commit_position + 1)
    suspect_3.put()

    analysis = MasterFlakeAnalysis.Create(master, builder, build_number, step,
                                          test)
    analysis.data_points = [
      DataPoint.Create(),
      DataPoint.Create(),
      DataPoint.Create()
    ]
    analysis.suspected_flake_build_number = 5
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.culprit_urlsafe_key = culprit.key.urlsafe()
    culprit.cr_notification_status = analysis_status.COMPLETED
    culprit.put()
    analysis.confidence_in_culprit = 1.0
    analysis.suspect_urlsafe_keys = [
      suspect_1.key.urlsafe(),
      suspect_2.key.urlsafe(),
      suspect_3.key.urlsafe()
    ]
    analysis.put()

    event = event_reporting.CreateTestFlakeAnalysisCompletionEvent(analysis)
    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)

    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number,
                     suspected_build_number)

    self.assertEqual(
        str(event.analysis_info.suspects),
        '[host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n, host: '
        '"chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n, host: '
        '"chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n]')
    self.assertEqual(
        str(event.analysis_info.culprit),
        'host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\nconfidence: 1.0\n')

    self.assertEqual(event.analysis_info.outcomes, [
      findit_pb2.CULPRIT, findit_pb2.SUSPECT, findit_pb2.REGRESSION_IDENTIFIED
    ])
    self.assertEqual(event.analysis_info.actions, [findit_pb2.CL_COMMENTED])
    self.assertTrue(event.flake)

  def testCreateTestFlakeAnalysisCompletionEventWithBugCreation(self):
    """Test reporting event where the Cr has been notified."""
    master = 'master'
    builder = 'builder'
    step = 'step'
    test = 'test'

    build_number = 10
    suspected_build_number = 5

    repo = 'chromium'
    revision = 'revision'
    culprit_commit_position = 2
    culprit = FlakeCulprit.Create(repo, revision, culprit_commit_position)

    suspect_1 = FlakeCulprit.Create(repo, revision, culprit_commit_position - 1)
    suspect_1.put()
    suspect_2 = FlakeCulprit.Create(repo, revision, culprit_commit_position)
    suspect_2.put()
    suspect_3 = FlakeCulprit.Create(repo, revision, culprit_commit_position + 1)
    suspect_3.put()

    analysis = MasterFlakeAnalysis.Create(master, builder, build_number, step,
                                          test)
    analysis.data_points = [
      DataPoint.Create(),
      DataPoint.Create(),
      DataPoint.Create()
    ]
    analysis.suspected_flake_build_number = 5
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.culprit_urlsafe_key = culprit.key.urlsafe()
    culprit.put()
    analysis.confidence_in_culprit = 1.0
    analysis.suspect_urlsafe_keys = [
      suspect_1.key.urlsafe(),
      suspect_2.key.urlsafe(),
      suspect_3.key.urlsafe()
    ]
    analysis.bug_id = 1
    analysis.has_attempted_filing = True
    analysis.put()

    event = event_reporting.CreateTestFlakeAnalysisCompletionEvent(analysis)
    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)
    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number,
                     suspected_build_number)

    self.assertEqual(
        str(event.analysis_info.suspects),
        '[host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n, host: '
        '"chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n, host: '
        '"chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n]')
    self.assertEqual(
        str(event.analysis_info.culprit),
        'host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\nconfidence: 1.0\n')

    self.assertEqual(event.analysis_info.outcomes, [
      findit_pb2.CULPRIT, findit_pb2.SUSPECT, findit_pb2.REGRESSION_IDENTIFIED
    ])
    self.assertEqual(event.analysis_info.actions, [findit_pb2.BUG_CREATED])

  def testCreateTestFlakeAnalysisCompletionEventWithBugComment(self):
    """Test reporting event where the Cr has been notified."""
    master = 'master'
    builder = 'builder'
    step = 'step'
    test = 'test'

    build_number = 10
    suspected_build_number = 5

    repo = 'chromium'
    revision = 'revision'
    culprit_commit_position = 2
    culprit = FlakeCulprit.Create(repo, revision, culprit_commit_position)

    suspect_1 = FlakeCulprit.Create(repo, revision, culprit_commit_position - 1)
    suspect_1.put()
    suspect_2 = FlakeCulprit.Create(repo, revision, culprit_commit_position)
    suspect_2.put()
    suspect_3 = FlakeCulprit.Create(repo, revision, culprit_commit_position + 1)
    suspect_3.put()

    analysis = MasterFlakeAnalysis.Create(master, builder, build_number, step,
                                          test)
    analysis.data_points = [
      DataPoint.Create(),
      DataPoint.Create(),
      DataPoint.Create()
    ]
    analysis.suspected_flake_build_number = 5
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.culprit_urlsafe_key = culprit.key.urlsafe()
    culprit.put()
    analysis.confidence_in_culprit = 1.0
    analysis.confidence_in_suspected_build = 1.0
    analysis.suspect_urlsafe_keys = [
      suspect_1.key.urlsafe(),
      suspect_2.key.urlsafe(),
      suspect_3.key.urlsafe()
    ]
    analysis.bug_id = 1
    analysis.bug_reported_by = triggering_sources.FINDIT_API
    analysis.put()

    event = event_reporting.CreateTestFlakeAnalysisCompletionEvent(analysis)
    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)
    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number,
                     suspected_build_number)

    self.assertEqual(
        str(event.analysis_info.suspects),
        '[host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n, host: '
        '"chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n, host: '
        '"chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n]')
    self.assertEqual(
        str(event.analysis_info.culprit),
        'host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\nconfidence: 1.0\n')

    self.assertEqual(event.analysis_info.outcomes, [
      findit_pb2.CULPRIT, findit_pb2.SUSPECT, findit_pb2.REGRESSION_IDENTIFIED
    ])
    self.assertEqual(event.analysis_info.actions, [findit_pb2.BUG_COMMENTED])

  def testCreateTestFlakeAnalysisCompletionEventWithSuspects(self):
    """Test reporting event where the Cr has been notified."""
    master = 'master'
    builder = 'builder'
    step = 'step'
    test = 'test'

    build_number = 10
    suspected_build_number = 5

    repo = 'chromium'
    revision = 'revision'
    culprit_commit_position = 2
    suspect_1 = FlakeCulprit.Create(repo, revision, culprit_commit_position - 1)
    suspect_1.put()
    suspect_2 = FlakeCulprit.Create(repo, revision, culprit_commit_position)
    suspect_2.put()
    suspect_3 = FlakeCulprit.Create(repo, revision, culprit_commit_position + 1)
    suspect_3.put()

    analysis = MasterFlakeAnalysis.Create(master, builder, build_number, step,
                                          test)
    analysis.data_points = [
      DataPoint.Create(),
      DataPoint.Create(),
      DataPoint.Create()
    ]
    analysis.suspected_flake_build_number = 5
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.confidence_in_culprit = 1.0
    analysis.suspect_urlsafe_keys = [
      suspect_1.key.urlsafe(),
      suspect_2.key.urlsafe(),
      suspect_3.key.urlsafe()
    ]
    analysis.put()

    event = event_reporting.CreateTestFlakeAnalysisCompletionEvent(analysis)
    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)
    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number,
                     suspected_build_number)

    self.assertEqual(
        str(event.analysis_info.suspects),
        '[host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n, host: '
        '"chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n, host: '
        '"chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "revision"\n]')

    self.assertEqual(event.analysis_info.outcomes,
                     [findit_pb2.SUSPECT, findit_pb2.REGRESSION_IDENTIFIED])

  def testCreateTestFlakeAnalysisCompletionEventWithNoSuspects(self):
    """Test reporting event where the Cr has been notified."""
    master = 'master'
    builder = 'builder'
    step = 'step'
    test = 'test'

    build_number = 10
    suspected_build_number = 5

    analysis = MasterFlakeAnalysis.Create(master, builder, build_number, step,
                                          test)
    # Need two or more points for a valid regression range.
    analysis.data_points = [DataPoint.Create(), DataPoint.Create()]
    analysis.suspected_flake_build_number = 5
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.put()

    event = event_reporting.CreateTestFlakeAnalysisCompletionEvent(analysis)
    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)
    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number,
                     suspected_build_number)

    self.assertEqual(event.analysis_info.outcomes,
                     [findit_pb2.REGRESSION_IDENTIFIED])

  def testCreateTestFlakeAnalysisCompletionEventWithNoRegressionRange(self):
    """Test reporting event where the Cr has been notified."""
    master = 'master'
    builder = 'builder'
    step = 'step'
    test = 'test'

    build_number = 10

    analysis = MasterFlakeAnalysis.Create(master, builder, build_number, step,
                                          test)
    analysis.data_points = [DataPoint.Create(pass_rate=1)]
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.put()

    event = event_reporting.CreateTestFlakeAnalysisCompletionEvent(analysis)
    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)
    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)

    self.assertEqual(event.analysis_info.outcomes,
                     [findit_pb2.NOT_REPRODUCIBLE])

  def testCreateTestFlakeAnalysisCompletionEventWithNoDataPoints(self):
    """Test reporting event where the Cr has been notified."""
    master = 'master'
    builder = 'builder'
    step = 'step'
    test = 'test'

    build_number = 10

    analysis = MasterFlakeAnalysis.Create(master, builder, build_number, step,
                                          test)
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.put()

    event = event_reporting.CreateTestFlakeAnalysisCompletionEvent(analysis)
    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)
    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)

    self.assertEqual(event.analysis_info.outcomes, [])

  def testCreateCompileFailureAnalysisCompletionEvent(self):
    master = 'm'
    builder = 'b'
    build_number = 10
    step = 'compile'

    suspected_build_number = 1

    failure_info = {
      'failed_steps': {
        'compile': {
          'current_failure': 2,
          'first_failure': suspected_build_number,
        }
      }
    }
    signals_json = {
      'compile': {
        'files': {
          'a/b/c.cc': [307],
          'a/b/d.cc': [123],
        },
        'keywords': {},
        'failed_output_nodes': [
          'obj/a/b/test.c.o',
          'obj/a/b/test.d.o',
        ],
        'failed_edges': [{
          'rule': 'CXX',
          'output_nodes': ['obj/a/b/test.c.o'],
          'dependencies': ["a/b/c.cc", "new/c.cc"]
        }, {
          'rule': 'LINK',
          'output_nodes': ['obj/a/b/test.d.o'],
          'dependencies': []
        }]
      }
    }
    repo_name = 'chromium'
    revision = 'rev1'
    commit_position = 1

    try_job = WfTryJob.Create(master, builder, suspected_build_number)
    try_job.compile_results = [{
      "culprit": {
        "compile": {
          "url": "https://chromium-review.googlesource.com/q/asdf",
          "author": "fdsa@chromium.com",
          "commit_position": commit_position,
          "revision": revision,
          "repo_name": repo_name
        }
      }
    }]
    try_job.put()

    culprit_cl = WfSuspectedCL.Create(repo_name, revision, commit_position)
    culprit_cl.revert_submission_status = analysis_status.COMPLETED
    culprit_cl.put()

    analysis = WfAnalysis.Create(master, builder, build_number)
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.failure_info = failure_info
    analysis.signals = signals_json
    analysis.suspected_cls = [{
      'repo_name': repo_name,
      'revision': revision,
      'commit_position': commit_position,
      'url': 'https://codereview.chromium.org/123',
      'status': suspected_cl_status.CORRECT
    },
      {
        'repo_name': repo_name,
        'revision': revision,
        'commit_position': commit_position,
        'url': 'https://codereview.chromium.org/123',
        'status': suspected_cl_status.CORRECT,
        'top_score': None
      }]
    analysis.put()

    event = event_reporting.CreateCompileFailureAnalysisCompletionEvent(
        analysis)

    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)

    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number,
                     suspected_build_number)

    self.assertEqual(
        str(event.analysis_info.suspects),
        '[host: "codereview.chromium.org"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "rev1"\n]')
    self.assertEqual(
        str(event.analysis_info.culprit),
        'host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "rev1"\n')

    self.assertEqual(event.analysis_info.outcomes,
                     [findit_pb2.CULPRIT, findit_pb2.SUSPECT])
    self.assertEqual(event.analysis_info.actions, [findit_pb2.REVERT_SUBMITTED])

    self.assertEqual(event.failed_build_rules, ['CXX', 'LINK'])

  def testCreateCompileFailureAnalysisCompletionEventNoFailureInfo(self):
    master = 'm'
    builder = 'b'
    build_number = 10
    step = 'compile'

    suspected_build_number = 1

    signals_json = {
      'compile': {
        'files': {
          'a/b/c.cc': [307],
          'a/b/d.cc': [123],
        },
        'keywords': {},
        'failed_output_nodes': [
          'obj/a/b/test.c.o',
          'obj/a/b/test.d.o',
        ],
        'failed_edges': [{
          'rule': 'CXX',
          'output_nodes': ['obj/a/b/test.c.o'],
          'dependencies': ["a/b/c.cc", "new/c.cc"]
        }, {
          'rule': 'LINK',
          'output_nodes': ['obj/a/b/test.d.o'],
          'dependencies': []
        }]
      }
    }
    repo_name = 'chromium'
    revision = 'rev1'
    commit_position = 1

    try_job = WfTryJob.Create(master, builder, suspected_build_number)
    try_job.compile_results = [{
      "culprit": {
        "compile": {
          "url": "https://chromium-review.googlesource.com/q/asdf",
          "author": "fdsa@chromium.com",
          "commit_position": commit_position,
          "revision": revision,
          "repo_name": repo_name
        }
      }
    }]
    try_job.put()

    culprit_cl = WfSuspectedCL.Create(repo_name, revision, commit_position)
    culprit_cl.revert_submission_status = analysis_status.COMPLETED
    culprit_cl.put()

    analysis = WfAnalysis.Create(master, builder, build_number)
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.failure_info = {}
    analysis.signals = signals_json
    analysis.suspected_cls = [{
      'repo_name': repo_name,
      'revision': revision,
      'commit_position': commit_position,
      'url': 'https://codereview.chromium.org/123',
      'status': suspected_cl_status.CORRECT
    }]
    analysis.put()

    event = event_reporting.CreateCompileFailureAnalysisCompletionEvent(
        analysis)

    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)

    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number, 0)

    self.assertEqual(
        str(event.analysis_info.suspects),
        '[host: "codereview.chromium.org"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "rev1"\n]')
    self.assertEqual(str(event.analysis_info.culprit), '')

    self.assertEqual(event.analysis_info.outcomes, [findit_pb2.SUSPECT])
    self.assertEqual(event.analysis_info.actions, [])

    self.assertEqual(event.failed_build_rules, ['CXX', 'LINK'])

  def testCreateCompileFailureAnalysisCompletionEventNoCulprit(self):
    master = 'm'
    builder = 'b'
    build_number = 10
    step = 'compile'

    suspected_build_number = 1

    failure_info = {
      'failed_steps': {
        'compile': {
          'current_failure': 2,
          'first_failure': suspected_build_number,
        }
      }
    }
    signals_json = {
      'compile': {
        'files': {
          'a/b/c.cc': [307],
          'a/b/d.cc': [123],
        },
        'keywords': {},
        'failed_output_nodes': [
          'obj/a/b/test.c.o',
          'obj/a/b/test.d.o',
        ],
        'failed_edges': [{
          'rule': 'CXX',
          'output_nodes': ['obj/a/b/test.c.o'],
          'dependencies': ["a/b/c.cc", "new/c.cc"]
        }, {
          'rule': 'LINK',
          'output_nodes': ['obj/a/b/test.d.o'],
          'dependencies': []
        }]
      }
    }
    repo_name = 'chromium'
    revision = 'rev1'
    commit_position = 1

    try_job = WfTryJob.Create(master, builder, suspected_build_number)
    try_job.compile_results = []
    try_job.put()

    culprit_cl = WfSuspectedCL.Create(repo_name, revision, commit_position)
    culprit_cl.revert_submission_status = analysis_status.COMPLETED
    culprit_cl.put()

    analysis = WfAnalysis.Create(master, builder, build_number)
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.failure_info = failure_info
    analysis.signals = signals_json
    analysis.suspected_cls = [{
      'repo_name': repo_name,
      'revision': revision,
      'commit_position': commit_position,
      'url': 'https://codereview.chromium.org/123',
      'status': suspected_cl_status.CORRECT
    }]
    analysis.put()

    event = event_reporting.CreateCompileFailureAnalysisCompletionEvent(
        analysis)

    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)

    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number,
                     suspected_build_number)

    self.assertEqual(
        str(event.analysis_info.suspects),
        '[host: "codereview.chromium.org"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "rev1"\n]')

    self.assertEqual(event.analysis_info.outcomes, [findit_pb2.SUSPECT])
    self.assertEqual(event.analysis_info.actions, [])

    self.assertEqual(event.failed_build_rules, ['CXX', 'LINK'])

  def testCreateCompileFailureAnalysisCompletionEventNoSuspect(self):
    master = 'm'
    builder = 'b'
    build_number = 10
    step = 'compile'

    failure_info = {
      'failed_steps': {
        'compile': {
          'current_failure': 2,
          'first_failure': build_number,
        }
      }
    }
    signals_json = {
      'compile': {
        'files': {
          'a/b/c.cc': [307],
          'a/b/d.cc': [123],
        },
        'keywords': {},
        'failed_output_nodes': [
          'obj/a/b/test.c.o',
          'obj/a/b/test.d.o',
        ],
        'failed_edges': [{
          'rule': 'CXX',
          'output_nodes': ['obj/a/b/test.c.o'],
          'dependencies': ["a/b/c.cc", "new/c.cc"]
        }, {
          'rule': 'LINK',
          'output_nodes': ['obj/a/b/test.d.o'],
          'dependencies': []
        }]
      }
    }
    repo_name = 'chromium'
    revision = 'rev1'
    commit_position = 1

    try_job = WfTryJob.Create(master, builder, build_number)
    try_job.compile_results = []
    try_job.put()

    culprit_cl = WfSuspectedCL.Create(repo_name, revision, commit_position)
    culprit_cl.revert_submission_status = analysis_status.COMPLETED
    culprit_cl.put()

    analysis = WfAnalysis.Create(master, builder, build_number)
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.failure_info = failure_info
    analysis.signals = signals_json
    analysis.suspected_cls = []
    analysis.put()

    event = event_reporting.CreateCompileFailureAnalysisCompletionEvent(
        analysis)

    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)

    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number, build_number)

    self.assertEqual(list(event.analysis_info.suspects), [])

    self.assertEqual(event.analysis_info.outcomes, [])
    self.assertEqual(event.analysis_info.actions, [])

    self.assertEqual(event.failed_build_rules, ['CXX', 'LINK'])

  def testCreateCompileFailureAnalysisCompletionEventRevertCreated(self):
    master = 'm'
    builder = 'b'
    build_number = 10
    step = 'compile'

    failure_info = {
      'failed_steps': {
        'compile': {
          'current_failure': 2,
          'first_failure': build_number,
        }
      }
    }
    signals_json = {
      'compile': {
        'files': {
          'a/b/c.cc': [307],
          'a/b/d.cc': [123],
        },
        'keywords': {},
        'failed_output_nodes': [
          'obj/a/b/test.c.o',
          'obj/a/b/test.d.o',
        ],
        'failed_edges': [{
          'rule': 'CXX',
          'output_nodes': ['obj/a/b/test.c.o'],
          'dependencies': ["a/b/c.cc", "new/c.cc"]
        }, {
          'rule': 'LINK',
          'output_nodes': ['obj/a/b/test.d.o'],
          'dependencies': []
        }]
      }
    }
    repo_name = 'chromium'
    revision = 'rev1'
    commit_position = 1

    try_job = WfTryJob.Create(master, builder, build_number)
    try_job.compile_results = [{
      "culprit": {
        "compile": {
          "url": "https://chromium-review.googlesource.com/q/asdf",
          "author": "fdsa@chromium.com",
          "commit_position": commit_position,
          "revision": revision,
          "repo_name": repo_name
        }
      }
    }]
    try_job.put()

    culprit_cl = WfSuspectedCL.Create(repo_name, revision, commit_position)
    culprit_cl.revert_status = analysis_status.COMPLETED
    culprit_cl.put()

    analysis = WfAnalysis.Create(master, builder, build_number)
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.failure_info = failure_info
    analysis.signals = signals_json
    analysis.suspected_cls = [{
      'repo_name': repo_name,
      'revision': revision,
      'commit_position': commit_position,
      'url': 'https://codereview.chromium.org/123',
      'status': suspected_cl_status.CORRECT
    }]
    analysis.put()

    event = event_reporting.CreateCompileFailureAnalysisCompletionEvent(
        analysis)

    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)

    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number, build_number)

    self.assertEqual(
        str(event.analysis_info.suspects),
        '[host: "codereview.chromium.org"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "rev1"\n]')
    self.assertEqual(
        str(event.analysis_info.culprit),
        'host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "rev1"\n')

    self.assertEqual(event.analysis_info.outcomes,
                     [findit_pb2.CULPRIT, findit_pb2.SUSPECT])
    self.assertEqual(event.analysis_info.actions, [findit_pb2.REVERT_CREATED])

    self.assertEqual(event.failed_build_rules, ['CXX', 'LINK'])

  def testCreateCompileFailureAnalysisCompletionEventClComment(self):
    master = 'm'
    builder = 'b'
    build_number = 10
    step = 'compile'

    failure_info = {
      'failed_steps': {
        'compile': {
          'current_failure': 2,
          'first_failure': build_number,
        }
      }
    }
    signals_json = {
      'compile': {
        'files': {
          'a/b/c.cc': [307],
          'a/b/d.cc': [123],
        },
        'keywords': {},
        'failed_output_nodes': [
          'obj/a/b/test.c.o',
          'obj/a/b/test.d.o',
        ],
        'failed_edges': [{
          'rule': 'CXX',
          'output_nodes': ['obj/a/b/test.c.o'],
          'dependencies': ["a/b/c.cc", "new/c.cc"]
        }, {
          'rule': 'LINK',
          'output_nodes': ['obj/a/b/test.d.o'],
          'dependencies': []
        }]
      }
    }
    repo_name = 'chromium'
    revision = 'rev1'
    commit_position = 1

    try_job = WfTryJob.Create(master, builder, build_number)
    try_job.compile_results = [{
      "culprit": {
        "compile": {
          "url": "https://chromium-review.googlesource.com/q/asdf",
          "author": "fdsa@chromium.com",
          "commit_position": commit_position,
          "revision": revision,
          "repo_name": repo_name
        }
      }
    }]
    try_job.put()

    culprit_cl = WfSuspectedCL.Create(repo_name, revision, commit_position)
    culprit_cl.cr_notification_status = analysis_status.COMPLETED
    culprit_cl.put()

    analysis = WfAnalysis.Create(master, builder, build_number)
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.failure_info = failure_info
    analysis.signals = signals_json
    analysis.suspected_cls = [{
      'repo_name': repo_name,
      'revision': revision,
      'commit_position': commit_position,
      'url': 'https://codereview.chromium.org/123',
      'status': suspected_cl_status.CORRECT
    }]
    analysis.put()

    event = event_reporting.CreateCompileFailureAnalysisCompletionEvent(
        analysis)

    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)

    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number, build_number)

    self.assertEqual(
        str(event.analysis_info.suspects),
        '[host: "codereview.chromium.org"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "rev1"\n]')
    self.assertEqual(
        str(event.analysis_info.culprit),
        'host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "rev1"\n')

    self.assertEqual(event.analysis_info.outcomes,
                     [findit_pb2.CULPRIT, findit_pb2.SUSPECT])
    self.assertEqual(event.analysis_info.actions, [findit_pb2.CL_COMMENTED])

    self.assertEqual(event.failed_build_rules, ['CXX', 'LINK'])

  def testCreateCompileFailureAnalysisCompletionEventNoSignals(self):
    master = 'm'
    builder = 'b'
    build_number = 10
    step = 'compile'

    failure_info = {
      'failed_steps': {
        'compile': {
          'current_failure': 2,
          'first_failure': build_number,
        }
      }
    }
    repo_name = 'chromium'
    revision = 'rev1'
    commit_position = 1

    try_job = WfTryJob.Create(master, builder, build_number)
    try_job.compile_results = [{
      "culprit": {
        "compile": {
          "url": "https://chromium-review.googlesource.com",
          "author": "fdsa@chromium.com",
          "commit_position": commit_position,
          "revision": revision,
          "repo_name": repo_name
        }
      }
    }]
    try_job.put()

    culprit_cl = WfSuspectedCL.Create(repo_name, revision, commit_position)
    culprit_cl.cr_notification_status = analysis_status.COMPLETED
    culprit_cl.put()

    analysis = WfAnalysis.Create(master, builder, build_number)
    analysis.start_time = datetime.datetime(2017, 1, 1)
    analysis.end_time = datetime.datetime(2017, 1, 2)
    analysis.failure_info = failure_info
    analysis.suspected_cls = [{
      'repo_name': repo_name,
      'revision': revision,
      'commit_position': commit_position,
      'url': 'https://codereview.chromium.org/123',
      'status': suspected_cl_status.CORRECT
    }]
    analysis.put()

    event = event_reporting.CreateCompileFailureAnalysisCompletionEvent(
        analysis)

    self.assertEqual(event.analysis_info.master_name, master)
    self.assertEqual(event.analysis_info.builder_name, builder)
    self.assertEqual(event.analysis_info.step_name, step)

    start = timestamp_pb2.Timestamp()
    start.FromDatetime(analysis.start_time)
    self.assertEqual(event.analysis_info.timestamp.started, start)

    complete = timestamp_pb2.Timestamp()
    complete.FromDatetime(analysis.end_time)
    self.assertEqual(event.analysis_info.timestamp.completed, complete)

    self.assertEqual(event.analysis_info.detected_build_number, build_number)
    self.assertEqual(event.analysis_info.culprit_build_number, build_number)

    self.assertEqual(
        str(event.analysis_info.suspects),
        '[host: "codereview.chromium.org"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "rev1"\n]')
    self.assertEqual(
        str(event.analysis_info.culprit),
        'host: "chromium-review.googlesource.com"\nproject: "chromium"\nref: '
        '"refs/heads/master"\nrevision: "rev1"\n')

    self.assertEqual(event.analysis_info.outcomes,
                     [findit_pb2.CULPRIT, findit_pb2.SUSPECT])
    self.assertEqual(event.analysis_info.actions, [findit_pb2.CL_COMMENTED])

    self.assertEqual(event.failed_build_rules, [])
