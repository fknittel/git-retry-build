// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package config

// setAllowFail updates allowFail property and return plan.
func setAllowFail(p *Plan, allowFail bool) *Plan {
	p.AllowFail = allowFail
	return p
}

// CrosRepairConfig provides config for repair cros setup in the lab task.
func CrosRepairConfig() *Configuration {
	return &Configuration{
		PlanNames: []string{
			PlanServo,
			PlanCrOS,
			PlanChameleon,
			PlanBluetoothPeer,
			PlanWifiRouter,
			PlanClosing,
		},
		Plans: map[string]*Plan{
			PlanServo:         setAllowFail(servoRepairPlan(), true),
			PlanCrOS:          setAllowFail(crosRepairPlan(), false),
			PlanChameleon:     setAllowFail(chameleonPlan(), true),
			PlanBluetoothPeer: setAllowFail(btpeerRepairPlan(), true),
			PlanWifiRouter:    setAllowFail(wifiRouterRepairPlan(), true),
			PlanClosing:       setAllowFail(crosClosePlan(), true),
		}}
}

// CrosDeployConfig provides config for deploy cros setup in the lab task.
func CrosDeployConfig() *Configuration {
	return &Configuration{
		PlanNames: []string{
			PlanServo,
			PlanCrOS,
			PlanChameleon,
			PlanBluetoothPeer,
			PlanWifiRouter,
			PlanClosing,
		},
		Plans: map[string]*Plan{
			PlanServo:         setAllowFail(servoRepairPlan(), false),
			PlanCrOS:          setAllowFail(crosDeployPlan(), false),
			PlanChameleon:     setAllowFail(chameleonPlan(), true),
			PlanBluetoothPeer: setAllowFail(btpeerRepairPlan(), true),
			PlanWifiRouter:    setAllowFail(wifiRouterRepairPlan(), true),
			PlanClosing:       setAllowFail(crosClosePlan(), true),
		},
	}
}

// crosClosePlan provides plan to close cros repair/deploy tasks.
func crosClosePlan() *Plan {
	return &Plan{
		CriticalActions: []string{
			"Remove in-use flag on servo-host",
			"Remove request to reboot is servo is good",
		},
		Actions: map[string]*Action{
			"servo_state_is_working": {
				Docs:          []string{"check the servo's state is ServoStateWorking."},
				ExecName:      "servo_match_state",
				ExecExtraArgs: []string{"state:WORKING"},
			},
			"Remove request to reboot is servo is good": {
				Conditions: []string{
					"is_not_flex_board",
					"servo_state_is_working",
				},
				ExecName:               "cros_remove_reboot_request",
				AllowFailAfterRecovery: true,
			},
			"Remove in-use flag on servo-host": {
				Conditions:             []string{"is_not_flex_board"},
				ExecName:               "cros_remove_servo_in_use",
				AllowFailAfterRecovery: true,
			},
			"is_not_flex_board": {
				Docs: []string{"Verify that device is belong Reven models"},
				ExecExtraArgs: []string{
					"string_values:x1c",
					"invert_result:true",
				},
				ExecName: "dut_check_model",
			},
		},
	}
}