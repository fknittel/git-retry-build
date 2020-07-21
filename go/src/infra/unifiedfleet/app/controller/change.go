// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package controller

import (
	"context"
	"fmt"

	"go.chromium.org/luci/common/logging"

	ufspb "infra/unifiedfleet/api/v1/proto"
	"infra/unifiedfleet/app/model/history"
	"infra/unifiedfleet/app/util"
)

// SaveChangeEvents saves change events to database
func SaveChangeEvents(ctx context.Context, changes []*ufspb.ChangeEvent) {
	_, err := history.CreateBatchChangeEvents(ctx, changes)
	if err != nil {
		logging.Debugf(ctx, "fail to log changes: %s", err.Error())
		for _, c := range changes {
			logging.Debugf(ctx, "\t%#v", c)
		}
	}
}

const (
	// LifeCycleRegistration indicates the registration of a resource.
	LifeCycleRegistration = "REGISTRATION"

	// LifeCycleDeployment indicates the deployment of a device.
	LifeCycleDeployment = "DEPLOYMENT"

	// LifeCycleDecomm indicates the decommission of a device.
	LifeCycleDecomm = "DECOMMISSION"

	// LifeCycleRetire indicates the retirement of a device.
	LifeCycleRetire = "RETIREMENT"
)

// logLifeCycle logs the life cycle event of a ChromeOSDevice.
func logLifeCycle(resourceName, resourceType, lifeCycleType string) (changes []*ufspb.ChangeEvent) {
	return []*ufspb.ChangeEvent{
		{
			Name:       resourceName,
			EventLabel: resourceType,
			OldValue:   lifeCycleType,
			NewValue:   lifeCycleType,
		},
	}
}

// LogMachineChanges logs the change of the given machine.
func LogMachineChanges(oldData *ufspb.Machine, newData *ufspb.Machine) []*ufspb.ChangeEvent {
	changes := make([]*ufspb.ChangeEvent, 0)
	if oldData == nil && newData == nil {
		return changes
	}
	if oldData == nil {
		return append(changes, logLifeCycle(util.AddPrefix(util.MachineCollection, newData.GetName()), "machine", LifeCycleRegistration)...)
	}
	if newData == nil {
		return append(changes, logLifeCycle(util.AddPrefix(util.MachineCollection, oldData.GetName()), "machine", LifeCycleRetire)...)
	}
	changes = append(changes, logCommon("", "machine.serial_number", oldData.GetSerialNumber(), newData.GetSerialNumber())...)
	changes = append(changes, logCommon("", "machine.location", oldData.GetLocation(), newData.GetLocation())...)
	changes = append(changes, logCommon("", "machine.realm", oldData.GetRealm(), newData.GetRealm())...)
	if newData.GetChromeBrowserMachine() != nil {
		changes = append(changes, logChromeBrowserMachine(oldData.GetChromeBrowserMachine(), newData.GetChromeBrowserMachine())...)
	}

	// Set resource name for all changes.
	for i := range changes {
		changes[i].Name = util.AddPrefix(util.MachineCollection, oldData.GetName())
	}
	return changes
}

// LogMachineLSEChanges logs the change of the given machine lse.
func LogMachineLSEChanges(oldData *ufspb.MachineLSE, newData *ufspb.MachineLSE) []*ufspb.ChangeEvent {
	changes := make([]*ufspb.ChangeEvent, 0)
	if oldData == nil && newData == nil {
		return changes
	}
	if oldData == nil {
		return append(changes, logLifeCycle(util.AddPrefix(util.MachineLSECollection, newData.GetName()), "machine_lse", LifeCycleRegistration)...)
	}
	if newData == nil {
		return append(changes, logLifeCycle(util.AddPrefix(util.MachineLSECollection, oldData.GetName()), "machine_lse", LifeCycleRetire)...)
	}
	changes = append(changes, logCommon("", "machine_lse.machine_lse_prototype", oldData.GetMachineLsePrototype(), newData.GetMachineLsePrototype())...)
	changes = append(changes, logCommon("", "machine_lse.hostname", oldData.GetHostname(), newData.GetHostname())...)
	changes = append(changes, logCommon("", "machine_lse.machines", oldData.GetMachines(), newData.GetMachines())...)
	if newData.GetChromeBrowserMachineLse() != nil {
		changes = append(changes, logChromeBrowserMachineLse(oldData.GetChromeBrowserMachineLse(), newData.GetChromeBrowserMachineLse())...)
	}

	// Set resource name for all changes.
	for i := range changes {
		changes[i].Name = util.AddPrefix(util.MachineLSECollection, oldData.GetName())
	}

	// Log VM differences
	vmChanges := make([]*ufspb.ChangeEvent, 0)
	if newData.GetChromeBrowserMachineLse() != nil {
		newVMs := newData.GetChromeBrowserMachineLse().GetVms()
		if oldData.GetChromeBrowserMachineLse() == nil {
			for _, vm := range newVMs {
				vmChanges = append(vmChanges, logLifeCycle(util.AddPrefix(util.VMCollection, vm.GetName()), "vm", LifeCycleRegistration)...)
			}
		} else {
			oldVMs := oldData.GetChromeBrowserMachineLse().GetVms()
			vmChanges = append(vmChanges, logVMs(oldVMs, newVMs)...)
		}
	} else {
		if oldData.GetChromeBrowserMachineLse() != nil {
			for _, vm := range oldData.GetChromeBrowserMachineLse().GetVms() {
				vmChanges = append(vmChanges, logLifeCycle(util.AddPrefix(util.VMCollection, vm.GetName()), "vm", LifeCycleRetire)...)
			}
		}
	}
	return append(changes, vmChanges...)
}

func logChromeBrowserMachineLse(oldData, newData *ufspb.ChromeBrowserMachineLSE) []*ufspb.ChangeEvent {
	changes := make([]*ufspb.ChangeEvent, 0)
	changes = append(changes, logCommon("", "machine_lse.chrome_browser_machine_lse.vm_capacity", oldData.GetVmCapacity(), newData.GetVmCapacity())...)
	changes = append(changes, logCommon("", "machine_lse.chrome_browser_machine_lse.os_version", oldData.GetOsVersion(), newData.GetOsVersion())...)
	return changes
}

func logVMs(oldData, newData []*ufspb.VM) []*ufspb.ChangeEvent {
	changes := make([]*ufspb.ChangeEvent, 0)
	vmMap := make(map[string]*ufspb.VM, len(oldData))
	for _, vm := range oldData {
		vmMap[vm.GetName()] = vm
	}
	newVMMap := make(map[string]*ufspb.VM, len(newData))
	for _, vm := range newData {
		newVMMap[vm.GetName()] = vm
		resourceName := util.AddPrefix(util.VMCollection, vm.GetName())
		if old, ok := vmMap[vm.GetName()]; ok {
			changes = append(changes, logCommon(resourceName, "vm.mac_address", old.GetMacAddress(), vm.GetMacAddress())...)
			changes = append(changes, logCommon(resourceName, "vm.os_version", old.GetOsVersion(), vm.GetOsVersion())...)
		} else {
			changes = append(changes, logLifeCycle(resourceName, "vm", LifeCycleRegistration)...)
		}
	}
	for _, vm := range oldData {
		resourceName := util.AddPrefix(util.VMCollection, vm.GetName())
		if _, ok := newVMMap[vm.GetName()]; !ok {
			changes = append(changes, logLifeCycle(resourceName, "vm", LifeCycleRetire)...)
		}
	}
	return changes
}

func logChromeBrowserMachine(oldData, newData *ufspb.ChromeBrowserMachine) []*ufspb.ChangeEvent {
	changes := make([]*ufspb.ChangeEvent, 0)
	changes = append(changes, logCommon("", "machine.chrome_browser_machine.display_name", oldData.GetDisplayName(), newData.GetDisplayName())...)
	changes = append(changes, logCommon("", "machine.chrome_browser_machine.chrome_platform", oldData.GetChromePlatform(), newData.GetChromePlatform())...)
	changes = append(changes, logCommon("", "machine.chrome_browser_machine.nics", oldData.GetNics(), newData.GetNics())...)
	changes = append(changes, logCommon("", "machine.chrome_browser_machine.drac", oldData.GetDrac(), newData.GetDrac())...)
	changes = append(changes, logCommon("", "machine.chrome_browser_machine.deployment_ticket", oldData.GetDeploymentTicket(), newData.GetDeploymentTicket())...)
	changes = append(changes, logKVMInterface(oldData.GetKvmInterface(), newData.GetKvmInterface())...)
	changes = append(changes, logRPMInterface(oldData.GetRpmInterface(), newData.GetRpmInterface())...)
	return changes
}

func logKVMInterface(oldData, newData *ufspb.KVMInterface) []*ufspb.ChangeEvent {
	changes := make([]*ufspb.ChangeEvent, 0)
	changes = append(changes, logCommon("", "machine.chrome_browser_machine.kvm_interface.kvm", oldData.GetKvm(), newData.GetKvm())...)
	return append(changes, logCommon("", "machine.chrome_browser_machine.kvm_interface.port", oldData.GetPort(), newData.GetPort())...)
}

func logRPMInterface(oldData, newData *ufspb.RPMInterface) []*ufspb.ChangeEvent {
	changes := make([]*ufspb.ChangeEvent, 0)
	changes = append(changes, logCommon("", "machine.chrome_browser_machine.rpm_interface.rpm", oldData.GetRpm(), newData.GetRpm())...)
	return append(changes, logCommon("", "machine.chrome_browser_machine.rpm_interface.port", oldData.GetPort(), newData.GetPort())...)
}

func logCommon(name, label string, oldValue interface{}, newValue interface{}) []*ufspb.ChangeEvent {
	oldValueStr := fmt.Sprintf("%v", oldValue)
	newValueStr := fmt.Sprintf("%v", newValue)
	if oldValueStr == newValueStr {
		return nil
	}
	return []*ufspb.ChangeEvent{
		{
			Name:       name,
			EventLabel: label,
			OldValue:   oldValueStr,
			NewValue:   newValueStr,
		},
	}
}
