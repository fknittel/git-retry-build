// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package utilization provides functions to report DUT utilization metrics.
package utilization

import (
	"context"
	"fmt"
	invV1 "infra/libs/skylab/inventory"

	"go.chromium.org/luci/common/api/swarming/swarming/v1"
	"go.chromium.org/luci/common/tsmon/field"
	"go.chromium.org/luci/common/tsmon/metric"
)

var dutmonMetric = metric.NewInt(
	"chromeos/skylab/dut_mon/swarming_dut_count",
	"The number of DUTs in a given bucket and status",
	nil,
	field.String("board"),
	field.String("model"),
	field.String("pool"),
	field.String("status"),
	field.Bool("is_locked"),
)

// ReportMetrics reports DUT utilization metrics akin to dutmon in Autotest
//
// The reported fields closely match those reported by dutmon, but the metrics
// path is different.
func ReportMetrics(ctx context.Context, bis []*swarming.SwarmingRpcsBotInfo) {
	c := make(counter)
	for _, bi := range bis {
		b := getBucketForBotInfo(bi)
		s := getStatusForBotInfo(bi)
		c.Increment(b, s)
	}
	c.Report(ctx)
}

// bucket contains static DUT dimensions.
//
// These dimensions do not change often. If all DUTs with a given set of
// dimensions are removed, the related metric is not automatically reset. The
// metric will get reset eventually.
type bucket struct {
	board       string
	model       string
	pool        string
	environment string
}

func (b bucket) String() string {
	return fmt.Sprintf("board: %s, model: %s, pool: %s, env: %s", b.board, b.model, b.pool, b.environment)
}

// status is a dynamic DUT dimension.
//
// This dimension changes often. If no DUTs have a particular status value,
// the corresponding metric is immediately reset.
type status string

var allStatuses = []status{"[None]", "Ready", "RepairFailed", "NeedsRepair", "NeedsReset", "Running"}

// counter collects number of DUTs per bucket and status.
type counter map[bucket]map[status]int

func (c counter) Increment(b bucket, s status) {
	sc := c[b]
	if sc == nil {
		sc = make(map[status]int)
		c[b] = sc
	}
	sc[s]++
}

func (c counter) Report(ctx context.Context) {
	for b, counts := range c {
		for _, s := range allStatuses {
			// TODO(crbug/929872) Report locked status once DUT leasing is
			// implemented in Skylab.
			dutmonMetric.Set(ctx, int64(counts[s]), b.board, b.model, b.pool, string(s), false)
		}
	}
}

func getBucketForBotInfo(bi *swarming.SwarmingRpcsBotInfo) bucket {
	b := bucket{
		board: "[None]",
		model: "[None]",
		pool:  "[None]",
	}
	for _, d := range bi.Dimensions {
		switch d.Key {
		case "label-board":
			b.board = summarizeValues(d.Value)
		case "label-model":
			b.model = summarizeValues(d.Value)
		case "label-pool":
			b.pool = getReportPool(d.Value)
		default:
			// Ignore other dimensions.
		}
	}
	return b
}

func getStatusForBotInfo(bi *swarming.SwarmingRpcsBotInfo) status {
	// dutState values are defined at
	// https://chromium.googlesource.com/infra/infra/+/e70c5ed1f9dddec833fad7e87567c0ded19fd565/go/src/infra/cmd/skylab_swarming_worker/internal/botinfo/botinfo.go#32
	dutState := ""
	for _, d := range bi.Dimensions {
		switch d.Key {
		case "dut_state":
			dutState = summarizeValues(d.Value)
			break
		default:
			// Ignore other dimensions.
		}
	}

	// Order matters: a bot may be dead and still have a task associated with it.
	if !isBotHealthy(bi) {
		return "[None]"
	}

	botBusy := bi.TaskId != ""

	switch dutState {
	case "ready":
		if botBusy {
			return "Running"
		}
		return "Ready"
	case "running":
		return "Running"
	case "needs_reset":
		// We count time spent waiting for a reset task to be assigned as time
		// spent Resetting.
		return "NeedsReset"
	case "needs_repair":
		// We count time spent waiting for a repair task to be assigned as time
		// spent Repairing.
		return "NeedsRepair"
	case "repair_failed":
		return "RepairFailed"

	default:
		return "[None]"
		// We should never see this state
	}
}

func isBotHealthy(bi *swarming.SwarmingRpcsBotInfo) bool {
	return !(bi.Deleted || bi.IsDead || bi.Quarantined)
}

func summarizeValues(vs []string) string {
	switch len(vs) {
	case 0:
		return "[None]"
	case 1:
		return vs[0]
	default:
		return "[Multiple]"
	}
}

func isManagedPool(p string) bool {
	_, ok := invV1.SchedulableLabels_DUTPool_value[p]
	return ok
}

func getReportPool(pools []string) string {
	p := summarizeValues(pools)
	if isManagedPool(p) {
		return fmt.Sprintf("managed:%s", p)
	}
	return p
}
