// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package autotest implements logic necessary for Autotest execution of an
// ExecuteRequest.
package autotest

import (
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"

	"go.chromium.org/chromiumos/infra/proto/go/test_platform"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/config"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/steps"
	swarming_api "go.chromium.org/luci/common/api/swarming/swarming/v1"
	"go.chromium.org/luci/common/clock"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"

	"infra/cmd/cros_test_platform/internal/execution/internal/autotest/parse"
	"infra/cmd/cros_test_platform/internal/execution/internal/common"
	"infra/cmd/cros_test_platform/internal/execution/isolate"
	"infra/cmd/cros_test_platform/internal/execution/swarming"
	"infra/libs/skylab/autotest/dynamicsuite"
)

const suiteName = "cros_test_platform"

// Runner runs a set of tests in autotest.
type Runner struct {
	invocations   []*steps.EnumerationResponse_AutotestInvocation
	requestParams *test_platform.Request_Params
	config        *config.Config_AutotestBackend

	response *steps.ExecuteResponse
}

// New returns a new autotest runner.
func New(tests []*steps.EnumerationResponse_AutotestInvocation, params *test_platform.Request_Params, config *config.Config_AutotestBackend) *Runner {
	return &Runner{invocations: tests, requestParams: params, config: config}
}

// LaunchAndWait launches an autotest execution and waits for it to complete.
func (r *Runner) LaunchAndWait(ctx context.Context, client swarming.Client, _ isolate.GetterFactory) error {
	if p := r.requestParams.GetScheduling().GetPriority(); p != 0 && p != 140 {
		// TODO(crbug.com/996301): Support priority in autotest backend.
		logging.Warningf(ctx, "request specifed a nondefault priority %d; this is unsupported in autotest backend, and ignored", p)
	}
	if len(r.requestParams.GetDecorations().GetTags()) != 0 {
		logging.Warningf(ctx, "request specified tags %s; this is unsupported in autotest backend, and ignored", r.requestParams.GetDecorations().GetTags())
	}

	taskID, err := r.launch(ctx, client)
	if err != nil {
		return err
	}
	logging.Infof(ctx, "launched task with ID %s", taskID)

	r.response = &steps.ExecuteResponse{State: &test_platform.TaskState{LifeCycle: test_platform.TaskState_LIFE_CYCLE_RUNNING}}

	if err = r.wait(ctx, client, taskID); err != nil {
		return err
	}

	r.response = &steps.ExecuteResponse{State: &test_platform.TaskState{LifeCycle: test_platform.TaskState_LIFE_CYCLE_COMPLETED}}

	resp, err := r.collect(ctx, client, taskID)
	if err != nil {
		return err
	}
	r.response = resp

	return nil
}

// Response constructs a response based on the current state of the
// Runner.
func (r *Runner) Response(swarming swarming.URLer) *steps.ExecuteResponse {
	return r.response
}

func (r *Runner) launch(ctx context.Context, client swarming.Client) (string, error) {
	req, err := r.proxyRequest()
	if err != nil {
		return "", errors.Annotate(err, "launch").Err()
	}

	logging.Debugf(ctx, "Launching proxy request %+v", req)

	resp, err := client.CreateTask(ctx, req)
	if err != nil {
		return "", errors.Annotate(err, "launch").Err()
	}

	logging.Debugf(ctx, "Launched proxy task at %s", client.GetTaskURL(resp.TaskId))

	return resp.TaskId, nil
}

func (r *Runner) proxyRequest() (*swarming_api.SwarmingRpcsNewTaskRequest, error) {
	if err := r.validate(); err != nil {
		return nil, errors.Annotate(err, "create proxy request").Err()
	}

	builds, err := common.ExtractBuilds(r.requestParams.SoftwareDependencies)
	if err != nil {
		return nil, errors.Annotate(err, "create proxy request").Err()
	}

	timeout, err := toTimeout(r.requestParams)
	if err != nil {
		return nil, errors.Annotate(err, "create proxy request").Err()
	}

	pool, err := toPool(r.requestParams.GetScheduling())
	if err != nil {
		return nil, errors.Annotate(err, "create proxy request").Err()
	}

	rArgs, err := r.reimageAndRunArgs()
	if err != nil {
		return nil, errors.Annotate(err, "create proxy request").Err()
	}

	afeHost := r.config.GetAfeHost()
	if afeHost == "" {
		return nil, errors.Reason("create proxy request: config specified no afe_host").Err()
	}

	if f := r.requestParams.GetFreeformAttributes().GetSwarmingDimensions(); len(f) != 0 {
		return nil, errors.Reason("swarming dimensions are not supported in autotest execution, but were specified as %s", f).Err()
	}

	dsArgs := dynamicsuite.Args{
		Board:             r.requestParams.SoftwareAttributes.BuildTarget.Name,
		Build:             builds.ChromeOS,
		FirmwareRWBuild:   builds.FirmwareRW,
		FirmwareROBuild:   builds.FirmwareRO,
		Model:             r.requestParams.HardwareAttributes.GetModel(),
		Timeout:           timeout,
		Pool:              pool,
		AfeHost:           afeHost,
		ReimageAndRunArgs: rArgs,
	}

	req, err := dynamicsuite.NewRequest(dsArgs)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func toPool(params *test_platform.Request_Params_Scheduling) (string, error) {
	switch v := params.GetPool().(type) {
	case *test_platform.Request_Params_Scheduling_ManagedPool_:
		if v.ManagedPool == test_platform.Request_Params_Scheduling_MANAGED_POOL_UNSPECIFIED {
			return "", errors.Reason("unspecified managed pool").Err()
		}
		if v.ManagedPool == test_platform.Request_Params_Scheduling_MANAGED_POOL_QUOTA {
			return "", errors.Reason("quota pool is not supported for autotest execution").Err()
		}
		longName := v.ManagedPool.String()
		return strings.ToLower(strings.TrimPrefix(longName, "MANAGED_POOL_")), nil
	case *test_platform.Request_Params_Scheduling_UnmanagedPool:
		return v.UnmanagedPool, nil
	case *test_platform.Request_Params_Scheduling_QuotaAccount:
		return "", errors.Reason("quota accounts are not valid for autotest execution").Err()
	default:
		return "", errors.Reason("no pool").Err()
	}
}

func toTimeout(params *test_platform.Request_Params) (time.Duration, error) {
	if params.Time == nil {
		return 0, errors.Reason("get timeout: nil params.time").Err()
	}
	duration, err := ptypes.Duration(params.Time.MaximumDuration)
	if err != nil {
		return 0, errors.Annotate(err, "get timeout").Err()
	}
	return duration, nil
}

func (r *Runner) validate() error {
	if r.requestParams == nil {
		return errors.Reason("nil request_params").Err()
	}
	if r.requestParams.SoftwareAttributes == nil {
		return errors.Reason("nil request_params.software_attributes").Err()
	}
	if r.requestParams.SoftwareAttributes.BuildTarget == nil {
		return errors.Reason("nil request_params.software_attributes.build_target").Err()
	}
	if r.requestParams.Time == nil {
		return errors.Reason("nil requests_params.time").Err()
	}
	return nil
}

func (r *Runner) wait(ctx context.Context, client swarming.Client, taskID string) error {
	for {
		complete, err := r.tick(ctx, client, taskID)
		if complete {
			return nil
		}
		if err != nil {
			return errors.Annotate(err, "wait for task %s completion", taskID).Err()
		}
		select {
		case <-ctx.Done():
			return errors.Annotate(ctx.Err(), "wait for task %s completion", taskID).Err()
		case <-clock.After(ctx, 15*time.Second):
		}
	}
}

func (r *Runner) tick(ctx context.Context, client swarming.Client, taskID string) (complete bool, err error) {
	results, err := client.GetResults(ctx, []string{taskID})
	if err != nil {
		return false, err
	}

	if len(results) != 1 {
		return false, errors.Reason("expected 1 result, found %d", len(results)).Err()
	}

	taskState, err := swarming.AsTaskState(results[0].State)
	if err != nil {
		return false, err
	}

	return !swarming.UnfinishedTaskStates[taskState], nil
}

func (r Runner) collect(ctx context.Context, client swarming.Client, taskID string) (*steps.ExecuteResponse, error) {
	resps, err := client.GetTaskOutputs(ctx, []string{taskID})
	if err != nil {
		return nil, errors.Annotate(err, "collect results").Err()
	}

	if len(resps) != 1 {
		return nil, errors.Reason("collect results: expected 1 result, got %d", len(resps)).Err()
	}

	output := resps[0].Output
	response, err := parse.RunSuite(output)
	if err != nil {
		return nil, errors.Annotate(err, "collect results").Err()
	}

	return response, nil
}

func (r *Runner) reimageAndRunArgs() (interface{}, error) {
	testNames := make([]string, len(r.invocations))
	for i, v := range r.invocations {
		switch {
		case v.DisplayName != "":
			// crbug.com/993998: This hack promotes the display name to the test
			// name for Autotest backend, ignoring any test arguments.
			// This allows the paygen stage in CI builders to request autoupdate
			// tests using dynamically generated control files in a uniform way
			// across the autotest and skylab backends.
			testNames[i] = v.DisplayName
		case v.TestArgs != "":
			return nil, errors.Reason(
				"test args %s were specified for test %s; test args are not supported in autotest backend",
				v.TestArgs, v.Test.Name).Err()
		default:
			testNames[i] = v.Test.Name
		}
	}
	return map[string]interface{}{
		// test_names is in argument to reimage_and_run which, if provided, short
		// cuts the normal test enumeration code and replaces it with this list
		// of tests.
		"test_names":  testNames,
		"name":        suiteName,
		"job_keyvals": r.requestParams.GetDecorations().GetAutotestKeyvals(),
	}, nil
}
