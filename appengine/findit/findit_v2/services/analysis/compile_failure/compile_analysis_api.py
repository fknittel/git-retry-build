# Copyright 2019 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
"""Special logic of pre compile analysis.

Build with compile failures will be pre-processed to determine if a new compile
analysis is needed or not.
"""

from findit_v2.model import luci_build
from findit_v2.model.compile_failure import CompileFailure
from findit_v2.model.compile_failure import CompileFailureAnalysis
from findit_v2.model.compile_failure import CompileFailureGroup
from findit_v2.services.analysis.analysis_api import AnalysisAPI
from findit_v2.services.failure_type import StepTypeEnum


class CompileAnalysisAPI(AnalysisAPI):

  @property
  def step_type(self):
    return StepTypeEnum.COMPILE

  def _GetMergedFailureKey(self, failure_entities, referred_build_id,
                           step_ui_name, atomic_failure):
    return CompileFailure.GetMergedFailureKey(
        failure_entities, referred_build_id, step_ui_name, atomic_failure)

  def _GetFailuresInBuild(self, project_api, build, failed_steps):
    return project_api.GetCompileFailures(build, failed_steps)

  def _GetFailuresWithMatchingFailureGroups(self, project_api, context, build,
                                            first_failures_in_current_build):
    return project_api.GetFailuresWithMatchingCompileFailureGroups(
        context, build, first_failures_in_current_build)

  def _CreateFailure(self, failed_build_key, step_ui_name,
                     first_failed_build_id, last_passed_build_id,
                     merged_failure_key, atomic_failure, properties):
    """Creates a CompileFailure entity."""
    return CompileFailure.Create(
        failed_build_key=failed_build_key,
        step_ui_name=step_ui_name,
        output_targets=list(atomic_failure or []),
        rule=(properties or {}).get('rule'),
        first_failed_build_id=first_failed_build_id,
        last_passed_build_id=last_passed_build_id,
        # Default to first_failed_build_id, will be updated later if matching
        # group exists.
        failure_group_build_id=first_failed_build_id,
        merged_failure_key=merged_failure_key)

  def _GetFailureEntitiesForABuild(self, build):
    build_entity = luci_build.LuciFailedBuild.get_by_id(build.id)
    assert build_entity, 'No LuciFailedBuild entity for build {}'.format(
        build.id)

    compile_failure_entities = CompileFailure.query(
        ancestor=build_entity.key).fetch()
    assert compile_failure_entities, (
        'No compile failure saved in datastore for build {}'.format(build.id))
    return compile_failure_entities

  def _CreateFailureGroup(self, context, build, compile_failure_keys,
                          last_passed_gitiles_id, last_passed_commit_position,
                          first_failed_commit_position):
    group_entity = CompileFailureGroup.Create(
        luci_project=context.luci_project_name,
        luci_bucket=build.builder.bucket,
        build_id=build.id,
        gitiles_host=context.gitiles_host,
        gitiles_project=context.gitiles_project,
        gitiles_ref=context.gitiles_ref,
        last_passed_gitiles_id=last_passed_gitiles_id,
        last_passed_commit_position=last_passed_commit_position,
        first_failed_gitiles_id=context.gitiles_id,
        first_failed_commit_position=first_failed_commit_position,
        compile_failure_keys=compile_failure_keys)
    return group_entity

  def _CreateFailureAnalysis(
      self, luci_project, context, build, last_passed_gitiles_id,
      last_passed_commit_position, first_failed_commit_position,
      rerun_builder_id, compile_failure_keys):
    analysis = CompileFailureAnalysis.Create(
        luci_project=luci_project,
        luci_bucket=build.builder.bucket,
        luci_builder=build.builder.builder,
        build_id=build.id,
        gitiles_host=context.gitiles_host,
        gitiles_project=context.gitiles_project,
        gitiles_ref=context.gitiles_ref,
        last_passed_gitiles_id=last_passed_gitiles_id,
        last_passed_commit_position=last_passed_commit_position,
        first_failed_gitiles_id=context.gitiles_id,
        first_failed_commit_position=first_failed_commit_position,
        rerun_builder_id=rerun_builder_id,
        compile_failure_keys=compile_failure_keys)
    return analysis
