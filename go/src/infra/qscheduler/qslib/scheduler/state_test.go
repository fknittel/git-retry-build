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

package scheduler

import (
	"fmt"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"

	"infra/qscheduler/qslib/tutils"
	"infra/qscheduler/qslib/types/vector"

	. "github.com/smartystreets/goconvey/convey"
)

func shouldResemblePretty(actual interface{}, expected ...interface{}) string {
	diff := pretty.Compare(actual, expected[0])
	if diff != "" {
		return fmt.Sprintf("Unexpected diff (-got +want): %s", diff)
	}
	return ""
}

// TestMarkIdle tests that a new worker is marked idle by MarkIdle, and that
// a subsequent updates to it are only respected if their timestamp is newer
// that the previous one.
func TestMarkIdle(t *testing.T) {
	tm0 := time.Unix(0, 0)
	tm1 := time.Unix(1, 0)
	tm2 := time.Unix(2, 0)
	Convey("Given an empty state", t, func() {
		state := NewState()
		Convey("when a worker marked idle at t=1", func() {
			state.markIdle("w1", []string{"old_label"}, tm1)
			Convey("then the worker is added to the state.", func() {
				want := NewState()
				want.Workers["w1"] = &Worker{Labels: []string{"old_label"}, ConfirmedTime: tutils.TimestampProto(tm1)}
				So(state, shouldResemblePretty, want)
			})
			Convey("when marking idle again with newer time t=2", func() {
				state.markIdle("w1", []string{"new_label"}, tm2)
				Convey("then the update is applied.", func() {
					So(tutils.Timestamp(state.Workers["w1"].ConfirmedTime), ShouldEqual, tm2)
					So(state.Workers["w1"].Labels, ShouldResemble, []string{"new_label"})
				})
			})
			Convey("when marking idle again with older time t=0", func() {
				state.markIdle("w1", []string{"new_label"}, tm0)
				Convey("then the update is ignored.", func() {
					So(tutils.Timestamp(state.Workers["w1"].ConfirmedTime), ShouldEqual, tm1)
					So(state.Workers["w1"].Labels, ShouldResemble, []string{"old_label"})
				})
			})
		})

		Convey("given a worker running a task at t=1", func() {
			state.markIdle("w1", []string{}, tm1)
			state.addRequest("r1", &TaskRequest{}, tm1)
			state.applyAssignment(&Assignment{Type: Assignment_IDLE_WORKER, RequestId: "r1", WorkerId: "w1"})
			Convey("when marking idle again with newer time t=2", func() {
				state.markIdle("w1", []string{}, tm2)
				Convey("then the update is applied.", func() {
					So(state.Workers["w1"].isIdle(), ShouldBeTrue)
					So(tutils.Timestamp(state.Workers["w1"].ConfirmedTime), ShouldEqual, tm2)
				})
			})

			Convey("when marking idle again with older time t=0", func() {
				state.markIdle("w1", []string{}, tm0)
				Convey("then the update is ignored.", func() {
					So(state.Workers["w1"].isIdle(), ShouldBeFalse)
					So(tutils.Timestamp(state.Workers["w1"].ConfirmedTime), ShouldEqual, tm1)
				})
			})
		})
	})
}

func TestNotifyRequest(t *testing.T) {
	tm0 := time.Unix(0, 0)
	tm1 := time.Unix(1, 0)
	tm2 := time.Unix(2, 0)
	tm3 := time.Unix(3, 0)
	tm4 := time.Unix(4, 0)
	Convey("Given a state with a request(t=1) and idle worker(t=3) and a match between them", t, func() {
		state := NewState()
		state.addRequest("r1", &TaskRequest{}, tm1)
		state.markIdle("w1", []string{}, tm3)
		a := &Assignment{
			Type:      Assignment_IDLE_WORKER,
			WorkerId:  "w1",
			RequestId: "r1",
		}
		state.applyAssignment(a)
		stateBeforeNotification := state.Clone()
		Convey("when notifying (idle request) with an older time t=0", func() {
			state.notifyRequest("r1", "", tm0)
			Convey("then the update is ignored.", func() {
				So(state, shouldResemblePretty, stateBeforeNotification)
			})
		})
		Convey("when notifying (idle request) with an intermediate time (between current request and worker time) t=1", func() {
			Convey("then the update is ignored.", func() {
				state.notifyRequest("r1", "", tm2)
				So(state, shouldResemblePretty, stateBeforeNotification)
			})
		})
		Convey("when notifying (idle request) with newer time t=4", func() {
			state.notifyRequest("r1", "", tm4)
			Convey("then the worker is deleted.", func() {
				So(state.Workers, ShouldNotContainKey, "w1")
			})
			Convey("then the request is deleted.", func() {
				So(state.QueuedRequests, ShouldContainKey, "r1")
			})
			Convey("then the request time is updated.", func() {
				So(tutils.Timestamp(state.QueuedRequests["r1"].ConfirmedTime), ShouldEqual, tm4)
			})
		})

		Convey("when notifying (correct match) with older time t=0", func() {
			state.notifyRequest("r1", "w1", tm0)
			Convey("then the update is ignored.", func() {
				So(state, shouldResemblePretty, stateBeforeNotification)
			})
		})
		Convey("when notifying (correct match) with intermediate time t=2", func() {
			state.notifyRequest("r1", "w1", tm2)
			Convey("then the request time is updated.", func() {
				So(tutils.Timestamp(state.Workers["w1"].RunningTask.Request.ConfirmedTime), ShouldEqual, tm2)
			})
		})
		Convey("when notifying (correct match) with newer time t=4", func() {
			state.notifyRequest("r1", "w1", tm4)
			Convey("then the request time is updated.", func() {
				So(tutils.Timestamp(state.Workers["w1"].RunningTask.Request.ConfirmedTime), ShouldEqual, tm4)
			})
			Convey("then the worker time is updated.", func() {
				So(tutils.Timestamp(state.Workers["w1"].ConfirmedTime), ShouldEqual, tm4)
			})
		})
	})

	Convey("Given a state with a matched request and worker both at t=1 and a separate idle worker at t=3", t, func() {
		state := NewState()
		state.addRequest("r1", &TaskRequest{}, tm1)
		state.markIdle("w1", []string{}, tm1)
		state.markIdle("w2", []string{}, tm3)
		a := &Assignment{
			Type:      Assignment_IDLE_WORKER,
			WorkerId:  "w1",
			RequestId: "r1",
		}
		state.applyAssignment(a)
		stateBefore := state.Clone()
		Convey("when notifying (contradictory match) with an older time t=0", func() {
			state.notifyRequest("r1", "w2", tm0)
			Convey("then the update is ignored.", func() {
				So(state, shouldResemblePretty, stateBefore)
			})
		})
		Convey("when notifying (contradictory match) with a time newer than match but older than idle worker t=2", func() {
			state.notifyRequest("r1", "w2", tm2)
			Convey("then the matching worker and request are deleted.", func() {
				So(state.QueuedRequests, ShouldNotContainKey, "r1")
				So(state.Workers, ShouldNotContainKey, "w1")
			})
		})

		Convey("when notifying (contradictory match) with a newer time t=4", func() {
			state.notifyRequest("r1", "w2", tm4)
			Convey("then the request and both workers are deleted.", func() {
				So(state, shouldResemblePretty, NewState())
			})
		})

	})

	Convey("Given a state with an idle worker(t=1), and a notify call with a match to an unknown request for that worker", t, func() {
		state := NewState()
		state.markIdle("w1", []string{}, tm1)
		stateBefore := state.Clone()
		Convey("when notifying (unknown request for worker) with older time t=0", func() {
			state.notifyRequest("r1", "w1", tm0)
			Convey("then the update is ignored.", func() {
				So(state, shouldResemblePretty, stateBefore)
			})
		})
		Convey("when notifying (unknown request for worker) with newer time t=2", func() {
			state.notifyRequest("r1", "w1", tm2)
			Convey("then the worker is deleted.", func() {
				So(state, shouldResemblePretty, NewState())
			})
		})
	})
}

// stateAtTime is a testing helper that creates an initialized but empty
//  State instance with the given time as its LastAccountUpdate time.
func stateAtTime(t time.Time) *State {
	s := NewState()
	s.LastUpdateTime = tutils.TimestampProto(t)
	return s
}

// Helper method to assert that two State instances are deeply equal.
func assertStateEqual(t *testing.T, desc string, got *State, want *State) {
	t.Helper()
	// TODO(akeshet): eliminate this call to got.regenCache() and update tests
	// accordingly, so that we can test that cache is updated correctly by
	// the various state mutations.
	got.regenCache()
	want.regenCache()
	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf(fmt.Sprintf("Case [%s] unexpected state diff (-got +want): %s", desc, diff))
	}
}

// TestApplyIdleAssignment tests that Apply for IDLE_WORKER behaves correctly.
func TestApplyIdleAssignment(t *testing.T) {
	t.Parallel()
	state := &State{
		QueuedRequests: map[string]*TaskRequest{"t1": &TaskRequest{}},
		Workers:        map[string]*Worker{"w1": NewWorker()},
	}

	expect := &State{
		QueuedRequests: map[string]*TaskRequest{},
		Workers: map[string]*Worker{
			"w1": &Worker{RunningTask: &TaskRun{
				RequestId: "t1",
				Priority:  1,
				Request:   &TaskRequest{},
				Cost:      vector.New()}},
		},
		RunningRequestsCache: map[string]string{"t1": "w1"},
	}

	mut := &Assignment{Type: Assignment_IDLE_WORKER, Priority: 1, RequestId: "t1", WorkerId: "w1"}
	state.applyAssignment(mut)

	if diff := pretty.Compare(state, expect); diff != "" {
		t.Errorf(fmt.Sprintf("Unexpected state diff (-got +want): %s", diff))
	}
}

// TestApplyPreempt tests that Apply for PREEMPT_WORKER behaves correctly.
func TestApplyPreempt(t *testing.T) {
	t.Parallel()
	state := &State{
		Balances: map[string]*vector.Vector{
			"a1": vector.New(),
			"a2": vector.New(2),
		},
		QueuedRequests: map[string]*TaskRequest{
			"t2": &TaskRequest{AccountId: "a2"},
		},
		Workers: map[string]*Worker{
			"w1": &Worker{RunningTask: &TaskRun{
				Cost:      vector.New(1),
				Priority:  2,
				Request:   &TaskRequest{AccountId: "a1"},
				RequestId: "t1",
			}},
		},
	}

	expect := &State{
		Balances: map[string]*vector.Vector{
			"a1": vector.New(1),
			"a2": vector.New(1),
		},
		QueuedRequests: map[string]*TaskRequest{
			"t1": &TaskRequest{AccountId: "a1"},
		},
		Workers: map[string]*Worker{
			"w1": &Worker{RunningTask: &TaskRun{
				Cost:      vector.New(1),
				Priority:  1,
				Request:   &TaskRequest{AccountId: "a2"},
				RequestId: "t2",
			},
			}},
		RunningRequestsCache: map[string]string{"t2": "w1"},
	}

	mut := &Assignment{Type: Assignment_PREEMPT_WORKER, Priority: 1, RequestId: "t2", WorkerId: "w1", TaskToAbort: "t1"}
	state.applyAssignment(mut)

	if diff := pretty.Compare(state, expect); diff != "" {
		t.Errorf(fmt.Sprintf("Unexpected state diff (-got +want): %s", diff))
	}
}
