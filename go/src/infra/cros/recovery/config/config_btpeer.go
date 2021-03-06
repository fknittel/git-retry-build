// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package config

import (
	"google.golang.org/protobuf/types/known/durationpb"
)

func btpeerRepairPlan() *Plan {
	return &Plan{
		CriticalActions: []string{
			"btpeer_state_broken",
			"Device is pingable",
			"cros_ssh",
			"check_server",
			"btpeer_state_working",
		},
		Actions: map[string]*Action{
			"check_server": {
				Docs:     []string{"To check if devices is responsive we request not empty list of detected statuses."},
				ExecName: "btpeer_get_detected_statuses",
			},
			"Device is pingable": {
				ExecTimeout: &durationpb.Duration{Seconds: 15},
				ExecName:    "cros_ping",
			},
		},
	}
}
