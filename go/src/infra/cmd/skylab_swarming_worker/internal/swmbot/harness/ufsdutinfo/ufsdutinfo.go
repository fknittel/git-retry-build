// Copyright 2019 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package ufsdutinfo implement loading Skylab DUT inventory(UFS) info for the
// worker.
package ufsdutinfo

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/retry"
	"infra/appengine/crosskylabadmin/api/fleet/v1"
	"infra/cmd/skylab_swarming_worker/internal/swmbot"
	"infra/libs/skylab/inventory"
	ufspb "infra/unifiedfleet/api/v1/models"
	ufsAPI "infra/unifiedfleet/api/v1/rpc"
	ufsutil "infra/unifiedfleet/app/util"
)

type DeviceType int64

const (
	ChromeOSDevice DeviceType = iota
	AndroidDevice
)

// Store holds a DUT's inventory info and adds a Close method.
type Store struct {
	DeviceType     DeviceType
	DUT            *inventory.DeviceUnderTest
	oldDUT         *inventory.DeviceUnderTest
	StableVersions map[string]string
	updateFunc     UpdateFunc
}

// Close updates the DUT's inventory info.  This method does nothing on
// subsequent calls.  This method is safe to call on a nil pointer.
func (s *Store) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}
	if s.updateFunc == nil {
		return nil
	}
	c := s.DUT.GetCommon()
	labels := c.GetLabels()
	inventory.SortLabels(labels)
	old := s.oldDUT.GetCommon().GetLabels()
	inventory.SortLabels(old)
	if labels.GetUselessSwitch() {
		*labels.UselessSwitch = false
	}

	log.Printf("Calling label update function")
	if err := s.updateFunc(ctx, c.GetId(), s.oldDUT, s.DUT); err != nil {
		return errors.Annotate(err, "close DUT inventory").Err()
	}
	s.updateFunc = nil
	return nil
}

// UpdateFunc is used to implement inventory updating for any changes
// to the loaded DUT info.
type UpdateFunc func(ctx context.Context, dutID string, oldDut *inventory.DeviceUnderTest, newDut *inventory.DeviceUnderTest) error

// LoadByID loads the DUT's info from the UFS by ID.
// This function returns a Store that should be closed to update the inventory
// with any changes to the info, using a supplied UpdateFunc. If UpdateFunc is
// nil, the inventory is not updated.
func LoadByID(ctx context.Context, b *swmbot.Info, dutID string, f UpdateFunc) (*Store, error) {
	return load(ctx, b, &ufsAPI.GetDeviceDataRequest{DeviceId: dutID}, f)
}

// LoadByHostname loads the DUT's info from the UFS by hostname.
// This function returns a Store that should be closed to update the inventory
// with any changes to the info, using a supplied UpdateFunc. If UpdateFunc is
// nil, the inventory is not updated.
func LoadByHostname(ctx context.Context, b *swmbot.Info, hostname string, f UpdateFunc) (*Store, error) {
	return load(ctx, b, &ufsAPI.GetDeviceDataRequest{Hostname: hostname}, f)
}

// GetDeviceData fetches a device entry from UFS.
// This function returns GetDeviceDataResponse which contains info about DUT,
// attached device or scheduling unit.
func GetDeviceData(ctx context.Context, b *swmbot.Info, req *ufsAPI.GetDeviceDataRequest) (*ufsAPI.GetDeviceDataResponse, error) {
	client, err := swmbot.UFSClient(ctx, b)
	if err != nil {
		return nil, errors.Annotate(err, "get device data: initialize UFS client").Err()
	}
	var resp *ufsAPI.GetDeviceDataResponse
	f := func() (err error) {
		osCtx := swmbot.SetupContext(ctx, ufsutil.OSNamespace)
		resp, err = client.GetDeviceData(osCtx, req)
		return err
	}
	if err := retry.Retry(ctx, retry.Default, f, retry.LogCallback(ctx, "ufsdutinfo.GetDeviceData")); err != nil {
		return nil, errors.Annotate(err, "get device data: retry GetDeviceData").Err()
	}
	return resp, nil
}

// getStableVersion fetches the current stable version from an inventory client
func getStableVersion(ctx context.Context, client fleet.InventoryClient, hostname string) (map[string]string, error) {
	log.Printf("getStableVersion: hostname (%s)", hostname)
	if hostname == "" {
		log.Printf("getStableVersion: failed validation for hostname")
		return nil, errors.Reason("getStableVersion: hostname cannot be \"\"").Err()
	}
	req := &fleet.GetStableVersionRequest{
		Hostname: hostname,
	}
	log.Printf("getStableVersion: client request (%v) with retries", req)
	res, err := retryGetStableVersion(ctx, client, req)
	log.Printf("getStableVersion: client response (%v)", res)
	if err != nil {
		return nil, err
	}
	s := map[string]string{
		"cros":       res.GetCrosVersion(),
		"faft":       res.GetFaftVersion(),
		"firmware":   res.GetFirmwareVersion(),
		"servo-cros": res.GetServoCrosVersion(),
	}
	log.Printf("getStableVersion: stable version map (%v)", s)
	return s, nil
}

func load(ctx context.Context, b *swmbot.Info, req *ufsAPI.GetDeviceDataRequest, uf UpdateFunc) (*Store, error) {
	ctx, err := swmbot.WithSystemAccount(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "setup system account").Err()
	}
	log.Printf("Loading DUT info from UFS")
	deviceData, err := GetDeviceData(ctx, b, req)
	if err != nil {
		return nil, errors.Annotate(err, "load DUT info from UFS").Err()
	}
	switch deviceData.GetResourceType() {
	case ufsAPI.GetDeviceDataResponse_RESOURCE_TYPE_CHROMEOS_DEVICE:
		return createChromeOSDeviceStore(ctx, b, deviceData.GetChromeOsDeviceData(), uf)
	case ufsAPI.GetDeviceDataResponse_RESOURCE_TYPE_ATTACHED_DEVICE:
		// Attached devices are not supported.
		return &Store{DeviceType: AndroidDevice}, nil
	}
	return nil, fmt.Errorf("load from UFS: invalid DUT type - %s", deviceData.GetResourceType())
}

func createChromeOSDeviceStore(ctx context.Context, b *swmbot.Info, deviceData *ufspb.ChromeOSDeviceData, uf UpdateFunc) (*Store, error) {
	dut := deviceData.GetDutV1()
	// TODO(gregorynisbet): should failure to get the stableversion information
	// cause the entire request to error out?
	c, err := swmbot.InventoryClient(ctx, b)
	if err != nil {
		return nil, errors.Annotate(err, "setup inventory client for stable_version").Err()
	}
	sv, err := getStableVersion(ctx, c, dut.GetCommon().GetHostname())
	if err != nil {
		sv = map[string]string{}
		log.Printf("create ChromeOS device store: getting stable version: sv (%v) err (%v)", sv, err)
	}
	// once we reach this point, sv is guaranteed to be non-nil
	return &Store{
		DeviceType:     ChromeOSDevice,
		DUT:            dut,
		oldDUT:         proto.Clone(dut).(*inventory.DeviceUnderTest),
		updateFunc:     uf,
		StableVersions: sv,
	}, nil
}

func retryGetStableVersion(ctx context.Context, client fleet.InventoryClient, req *fleet.GetStableVersionRequest) (*fleet.GetStableVersionResponse, error) {
	var resp *fleet.GetStableVersionResponse
	var err error
	f := func() error {
		resp, err = client.GetStableVersion(ctx, req)
		if err != nil {
			return err
		}
		return nil
	}
	if err := retry.Retry(ctx, retry.Default, f, retry.LogCallback(ctx, "ufsdutinfo.retryGetStableVersion")); err != nil {
		return nil, errors.Annotate(err, "retry getStableVersion").Err()
	}
	return resp, nil
}
