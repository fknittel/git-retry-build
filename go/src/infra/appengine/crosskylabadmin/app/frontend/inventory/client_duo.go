// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package inventory

import (
	"context"
	"math/rand"
	"time"

	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"

	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
	"infra/appengine/crosskylabadmin/app/frontend/internal/gitstore"
	"infra/libs/skylab/inventory"
)

var timeoutForEachDUT = 1 * time.Second

type duoClient struct {
	gc *gitStoreClient
	ic *invServiceClient

	// A number in [0, 100] indicate the write traffic (deploy/update)
	// duplicated to inventory v2 service.
	writeTrafficRatio int
	// A number in [0, 100] indicate the read traffic fanning out to inventory
	// v2 service.
	readTrafficRatio int
}

func newDuoClient(ctx context.Context, gs *gitstore.InventoryStore, host string, readTrafficRatio, writeTrafficRatio int) (inventoryClient, error) {
	gc, err := newGitStoreClient(ctx, gs)
	if err != nil {
		return nil, errors.Annotate(err, "create git client").Err()
	}
	ic, err := newInvServiceClient(ctx, host)
	if err != nil {
		logging.Infof(ctx, "Failed to create inventory client of the duo client. Just return the git store client")
		return gc, nil
	}
	return &duoClient{
		gc:                gc.(*gitStoreClient),
		ic:                ic.(*invServiceClient),
		readTrafficRatio:  readTrafficRatio,
		writeTrafficRatio: writeTrafficRatio,
	}, nil
}

func (client *duoClient) willWriteToV2() bool {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(100) < client.writeTrafficRatio
}

func (client *duoClient) addManyDUTsToFleet(ctx context.Context, nds []*inventory.CommonDeviceSpecs, pickServoPort bool) (string, []*inventory.CommonDeviceSpecs, error) {

	// Add DUTs to v1 first as it backfill DUT id and servo port. Then pass new
	// specs to v2 to make it sync with v1.
	url, ds, err := client.gc.addManyDUTsToFleet(ctx, nds, pickServoPort)
	logging.Infof(ctx, "[v1] add dut result: %s, %s", url, err)
	logging.Infof(ctx, "[v1] spec returned: %s", ds)

	if client.willWriteToV2() {
		go func() {
			// Set timeout for RPC call to inventory v2.
			// The timeout should correlated to how many DUTs being operated.
			ctx2, cancel := context.WithTimeout(ctx, time.Duration(len(nds))*timeoutForEachDUT)
			defer cancel()

			url2, ds2, err2 := client.ic.addManyDUTsToFleet(ctx2, ds, pickServoPort)
			logging.Infof(ctx2, "[v2] add dut result: %s, %s", url2, err2)
			logging.Infof(ctx2, "[v2] spec returned: %s", ds2)
		}()
	}

	return url, ds, err
}

func (client *duoClient) updateDUTSpecs(ctx context.Context, od, nd *inventory.CommonDeviceSpecs, pickServoPort bool) (string, error) {
	if client.willWriteToV2() {
		go func() {
			ctx2, cancel := context.WithTimeout(ctx, timeoutForEachDUT)
			defer cancel()

			url2, err2 := client.ic.updateDUTSpecs(ctx2, od, nd, pickServoPort)
			logging.Infof(ctx2, "[v2] add dut result: %s, %s", url2, err2)
		}()
	}

	url, err := client.gc.updateDUTSpecs(ctx, od, nd, pickServoPort)
	logging.Infof(ctx, "[v1] update dut result: %s, %s", url, err)
	return url, err
}

func (client *duoClient) deleteDUTsFromFleet(ctx context.Context, ids []string) (string, []string, error) {
	if client.willWriteToV2() {
		go func() {
			ctx2, cancel := context.WithTimeout(ctx, time.Duration(len(ids))*timeoutForEachDUT)
			defer cancel()
			url2, deletedIds2, err2 := client.ic.deleteDUTsFromFleet(ctx2, ids)
			logging.Infof(ctx2, "[v2] delete dut result: %s, %s, %s", url2, deletedIds2, err2)
		}()
	}
	url, deletedIds, err := client.gc.deleteDUTsFromFleet(ctx, ids)
	logging.Infof(ctx, "[v1] delete dut result: %s, %s, %s", url, deletedIds, err)

	return url, deletedIds, err
}

func (client *duoClient) selectDutsFromInventory(ctx context.Context, sel *fleet.DutSelector) ([]*inventory.DeviceUnderTest, error) {
	if client.willWriteToV2() {
		go func() {
			// Cannot know how many duts will be reutrned, so hard code the
			// timeout.
			ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			duts, _ := client.ic.selectDutsFromInventory(ctx2, sel)
			logging.Infof(ctx2, "[v2] select duts by %v", sel)
			if len(duts) > 0 {
				logging.Infof(ctx2, "[v2] selecting returns '%s'...(total %d duts)", duts[0].GetCommon().GetHostname(), len(duts))
			} else {
				logging.Infof(ctx2, "[v2] selecting returns 0 duts")
			}

		}()
	}
	return client.gc.selectDutsFromInventory(ctx, sel)
}

func (client *duoClient) commitBalancePoolChanges(ctx context.Context, changes []*fleet.PoolChange) (string, error) {
	if client.willWriteToV2() {
		go func() {
			ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			u, err := client.ic.commitBalancePoolChanges(ctx2, changes)
			logging.Infof(ctx2, "[v2] Commit balancing pool result: %s: %s", u, err)
		}()
	}
	return client.gc.commitBalancePoolChanges(ctx, changes)
}
