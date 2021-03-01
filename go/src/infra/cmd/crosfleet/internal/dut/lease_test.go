// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dut

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"infra/cmd/crosfleet/internal/common"
	"testing"
)

var testValidateData = []struct {
	leaseFlags
	wantValidationErrString string
}{
	{ // All flags raise errors
		leaseFlags{
			durationMins: 0,
			reason:       "this desc is barely too long!!!",
			host:         "",
			model:        "",
			board:        "",
		},
		`exactly one of board, model, or host should be specified
duration should be greater than 0
reason cannot exceed 30 characters`,
	},
	{ // Some flags raise errors
		leaseFlags{
			durationMins: 1441,
			reason:       "this desc is just short enough",
			host:         "sample-host",
			model:        "sample-model",
			board:        "sample-board",
		},
		`exactly one of board, model, or host should be specified
duration cannot exceed 1440 minutes (24 hours)`,
	},
	{ // No flags raise errors
		leaseFlags{
			durationMins: 1440,
			reason:       "this desc is just short enough",
			host:         "",
			model:        "sample-model",
		},
		"",
	},
}

func TestValidate(t *testing.T) {
	t.Parallel()
	for _, tt := range testValidateData {
		tt := tt
		t.Run(fmt.Sprintf("(%s)", tt.wantValidationErrString), func(t *testing.T) {
			t.Parallel()
			gotValidationErr := tt.leaseFlags.validate(&flag.FlagSet{})
			gotValidationErrString := common.ErrToString(gotValidationErr)
			if tt.wantValidationErrString != gotValidationErrString {
				t.Errorf("unexpected error: wanted %s, got %s", tt.wantValidationErrString, gotValidationErrString)
			}
		})
	}
}

// We avoid testing this function for a host-based lease since we'd have to
// fake a Swarming API call.
var testBotDimsAndBuildTagsData = []struct {
	leaseFlags
	wantDims, wantTags map[string]string
}{
	{ // Model-based lease with added dims
		leaseFlags{
			model:     "sample-model",
			reason:    "sample reason",
			addedDims: map[string]string{"added-key": "added-val"},
		},
		map[string]string{
			"added-key":   "added-val",
			"dut_state":   "ready",
			"label-model": "sample-model",
			"label-pool":  "DUT_POOL_QUOTA",
		},
		map[string]string{
			"added-key":      "added-val",
			"crosfleet-tool": "lease",
			"lease-by":       "model",
			"lease-reason":   "sample reason",
			"label-model":    "sample-model",
			"qs-account":     "leases",
		},
	},
	{ // Board-based lease without added dims
		leaseFlags{
			board:     "sample-board",
			reason:    "sample reason",
			addedDims: nil,
		},
		map[string]string{
			"dut_state":   "ready",
			"label-board": "sample-board",
			"label-pool":  "DUT_POOL_QUOTA",
		},
		map[string]string{
			"crosfleet-tool": "lease",
			"lease-by":       "board",
			"lease-reason":   "sample reason",
			"label-board":    "sample-board",
			"qs-account":     "leases",
		},
	},
}

func TestBotDimsAndBuildTagsData(t *testing.T) {
	t.Parallel()
	for _, tt := range testBotDimsAndBuildTagsData {
		tt := tt
		t.Run(fmt.Sprintf("(%s, %s)", tt.wantDims, tt.wantTags), func(t *testing.T) {
			ctx := context.Background()
			gotDims, gotTags, err := botDimsAndBuildTags(ctx, nil, tt.leaseFlags)
			if err != nil {
				t.Fatalf("unexpected error calling botDimsAndBuildTags: %v", err)
			}
			if dimDiff := cmp.Diff(tt.wantDims, gotDims); dimDiff != "" {
				t.Errorf("unexpected bot dimension diff (%s)", dimDiff)
			}
			if tagDiff := cmp.Diff(tt.wantTags, gotTags); tagDiff != "" {
				t.Errorf("unexpected build tag diff (%s)", tagDiff)
			}
		})
	}
}
