// Copyright 2018 The LUCI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package frontend

import (
	"context"
	"sort"

	qscheduler "infra/appengine/qscheduler-swarming/api/qscheduler/v1"
	"infra/swarming"

	"github.com/pkg/errors"

	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/luci/common/data/stringset"
	"go.chromium.org/luci/common/data/strpair"
	"go.chromium.org/luci/common/logging"

	"infra/qscheduler/qslib/reconciler"
	"infra/qscheduler/qslib/scheduler"
	"infra/qscheduler/qslib/tutils"
)

// AccountIDTagKey is the key used in Task tags to specify which quotascheduler
// account the task should be charged to.
const AccountIDTagKey = "qs_account"

// QSchedulerState encapsulates the state of a scheduler.
type QSchedulerState struct {
	schedulerID string
	scheduler   *scheduler.Scheduler
	reconciler  *reconciler.State
	config      *qscheduler.SchedulerPoolConfig
}

// QSchedulerServerImpl implements the QSchedulerServer interface.
type QSchedulerServerImpl struct {
	// TODO(akeshet): Implement in-memory cache of SchedulerPool struct, so that
	// we don't need to load and re-persist its state on every call.
	// TODO(akeshet): Implement request batching for AssignTasks and NotifyTasks.
	// TODO(akeshet): Determine if go.chromium.org/luci/server/caching has a
	// solution for in-memory caching like this.
}

// AssignTasks implements QSchedulerServer.
func (s *QSchedulerServerImpl) AssignTasks(ctx context.Context, r *swarming.AssignTasksRequest) (*swarming.AssignTasksResponse, error) {
	var response *swarming.AssignTasksResponse

	doAssign := func(ctx context.Context) error {
		sp, err := load(ctx, r.SchedulerId)
		if err != nil {
			return err
		}

		idles := make([]*reconciler.IdleWorker, len(r.IdleBots))
		for i, v := range r.IdleBots {
			idles[i] = &reconciler.IdleWorker{
				ID: v.BotId,
				// TODO(akeshet): Compute provisionable labels properly. This should actually
				// be the workers label set minus the scheduler pool's label set.
				ProvisionableLabels: v.Dimensions,
			}
		}

		a, err := sp.reconciler.AssignTasks(ctx, sp.scheduler, tutils.Timestamp(r.Time), idles...)
		if err != nil {
			return nil
		}

		assignments := make([]*swarming.TaskAssignment, len(a))
		for i, v := range a {
			assignments[i] = &swarming.TaskAssignment{
				BotId:  v.WorkerID,
				TaskId: v.RequestID,
			}
		}
		if err := save(ctx, sp); err != nil {
			return err
		}
		response = &swarming.AssignTasksResponse{Assignments: assignments}
		return nil
	}

	if err := datastore.RunInTransaction(ctx, doAssign, nil); err != nil {
		return nil, err
	}

	return response, nil
}

// GetCancellations implements QSchedulerServer.
func (s *QSchedulerServerImpl) GetCancellations(ctx context.Context, r *swarming.GetCancellationsRequest) (*swarming.GetCancellationsResponse, error) {
	sp, err := load(ctx, r.SchedulerId)
	if err != nil {
		return nil, err
	}

	c := sp.reconciler.Cancellations(ctx)
	rc := make([]*swarming.GetCancellationsResponse_Cancellation, len(c))
	for i, v := range c {
		rc[i] = &swarming.GetCancellationsResponse_Cancellation{BotId: v.WorkerID, TaskId: v.RequestID}
	}
	return &swarming.GetCancellationsResponse{Cancellations: rc}, nil
}

// NotifyTasks implements QSchedulerServer.
func (s *QSchedulerServerImpl) NotifyTasks(ctx context.Context, r *swarming.NotifyTasksRequest) (*swarming.NotifyTasksResponse, error) {

	doNotify := func(ctx context.Context) error {
		sp, err := load(ctx, r.SchedulerId)
		if err != nil {
			return err
		}

		if sp.config == nil {
			return errors.Errorf("Scheduler with id %s has nil config.", r.SchedulerId)
		}

		for _, n := range r.Notifications {
			var t reconciler.TaskUpdate_Type
			var ok bool
			if t, ok = toTaskState(n.Task.State); !ok {
				// TODO(akeshet): Return an error about unknown notification state.
				logging.Warningf(ctx, "Received notification with unhandled state %s.", n.Task.State)
				continue
			}

			var provisionableLabels []string
			// ProvisionableLabels attribute only matters for TaskUpdate_NEW type
			// updates, because these are tasks that are in the queue (scheduler
			// pays no attention to labels of already-running tasks).
			if t == reconciler.TaskUpdate_NEW {
				if provisionableLabels, err = getProvisionableLabels(n); err != nil {
					return err
				}
			}

			// TODO(akeshet): Validate that new tasks have dimensions that match the
			// worker pool dimensions for this scheduler pool.
			update := &reconciler.TaskUpdate{
				// TODO(akeshet): implement me. This will be based upon task tags.
				AccountId: "",
				// TODO(akeshet): implement me properly. This should be a separate field
				// of the task state, not the notification time.
				EnqueueTime:         n.Time,
				ProvisionableLabels: provisionableLabels,
				RequestId:           n.Task.Id,
				Time:                n.Time,
				Type:                t,
				WorkerId:            n.Task.BotId,
			}
			if err := sp.reconciler.Notify(ctx, sp.scheduler, update); err != nil {
				return err
			}
			logging.Debugf(ctx, "To scheduler with id %s, applied task update %+v", r.SchedulerId, update)
		}
		logState(ctx, sp.scheduler.State)
		return save(ctx, sp)
	}

	if err := datastore.RunInTransaction(ctx, doNotify, nil); err != nil {
		return nil, err
	}
	return &swarming.NotifyTasksResponse{}, nil
}

// getProvisionableLabels determines the provisionable labels for a given task,
// based on the dimensions of its slices.
func getProvisionableLabels(n *swarming.NotifyTasksItem) ([]string, error) {
	switch len(n.Task.Slices) {
	case 1:
		return []string{}, nil
	case 2:
		s1 := stringset.NewFromSlice(n.Task.Slices[0].Dimensions...)
		s2 := stringset.NewFromSlice(n.Task.Slices[1].Dimensions...)
		// s2 must be a subset of s1 (i.e. the first slice must be more specific about dimensions than the second one)
		// otherwise this is an error.
		if flaws := s2.Difference(s1); flaws.Len() != 0 {
			return nil, errors.Errorf("Invalid slice dimensions; task's 2nd slice dimensions are not a subset of 1st slice dimensions.")
		}

		var provisionable sort.StringSlice
		provisionable = s1.Difference(s2).ToSlice()
		provisionable.Sort()
		return provisionable, nil
	default:
		return nil, errors.Errorf("Invalid slice count %d; quotascheduler only supports 1-slice or 2-slice tasks.", len(n.Task.Slices))
	}
}

// getAccountID determines the account id for a given task, based on its tags.
func getAccountID(n *swarming.NotifyTasksItem) (string, error) {
	m := strpair.ParseMap(n.Task.Tags)
	accounts := m[AccountIDTagKey]
	switch len(accounts) {
	case 0:
		return "", nil
	case 1:
		return accounts[0], nil
	default:
		return "", errors.Errorf("Too many account tags.")
	}
}

func toTaskState(s swarming.TaskState) (reconciler.TaskUpdate_Type, bool) {
	// These cases appear in the same order as they are defined in swarming/proto/tasks.proto
	// If you add any cases here, please preserve their in-order appearance.
	switch s {
	case swarming.TaskState_RUNNING:
		return reconciler.TaskUpdate_ASSIGNED, true
	case swarming.TaskState_PENDING:
		return reconciler.TaskUpdate_NEW, true
	// The following states all translate to "ABORTED", because they are all equivalent
	// to the task being neither running nor enqueued.
	case swarming.TaskState_EXPIRED:
		fallthrough
	case swarming.TaskState_TIMED_OUT:
		fallthrough
	case swarming.TaskState_BOT_DIED:
		fallthrough
	case swarming.TaskState_CANCELED:
		fallthrough
	case swarming.TaskState_COMPLETED:
		fallthrough
	case swarming.TaskState_KILLED:
		fallthrough
	case swarming.TaskState_NO_RESOURCE:
		return reconciler.TaskUpdate_ABORTED, true

	// Invalid state.
	default:
		return reconciler.TaskUpdate_NULL, false
	}
}

func logState(ctx context.Context, s *scheduler.State) {
	logging.Debugf(ctx, "Scheduler has %d queued tasks, %d workers.", len(s.QueuedRequests), len(s.Workers))
}
