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

package inventory

import (
	"reflect"
	"testing"

	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
	"infra/appengine/crosskylabadmin/internal/app/config"
	"infra/appengine/crosskylabadmin/internal/app/gitstore/fakes"
	"infra/libs/skylab/inventory"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRemoveDutsFromDrones(t *testing.T) {
	Convey("With 2 DUTs assigned to drones (1 in prod, 1 in staging)", t, func() {
		tf, validate := newTestFixture(t)
		defer validate()

		dutID := "dut_id"
		dutHostname := "dut_hostname"
		serverID := "server_id"
		wrongEnvDutID := "wrong_env_dut"
		wrongEnvDutHostname := "wrong_env_dut_hostname"
		wrongEnvServer := "wrong_env_server"
		err := tf.FakeGitiles.SetInventory(config.Get(tf.C).Inventory, fakes.InventoryData{
			Lab: inventoryBytesFromDUTs([]testInventoryDut{
				{dutID, dutHostname, "link", "DUT_POOL_SUITES"},
				{wrongEnvDutID, wrongEnvDutHostname, "link", "DUT_POOL_SUITES"},
			}),
			Infrastructure: inventoryBytesFromServers([]testInventoryServer{
				{
					hostname:    serverID,
					environment: inventory.Environment_ENVIRONMENT_STAGING,
					dutIDs:      []string{dutID},
				},
				{
					hostname:    wrongEnvServer,
					environment: inventory.Environment_ENVIRONMENT_PROD,
					dutIDs:      []string{wrongEnvDutID},
				},
			}),
		})
		So(err, ShouldBeNil)

		Convey("RemoveDutsFromdrone for the staging dut removes it from drone.", func() {
			req := &fleet.RemoveDutsFromDronesRequest{
				Removals: []*fleet.RemoveDutsFromDronesRequest_Item{{DutId: dutID}},
			}
			resp, err := tf.Inventory.RemoveDutsFromDrones(tf.C, req)
			So(err, ShouldBeNil)
			So(resp.Removed, ShouldHaveLength, 1)
			So(resp.Removed[0].DutId, ShouldEqual, dutID)

			So(tf.FakeGerrit.Changes, ShouldHaveLength, 1)
			change := tf.FakeGerrit.Changes[0]
			p := "data/skylab/server_db.textpb"
			So(change.Files, ShouldContainKey, p)

			contents := change.Files[p]
			infra := &inventory.Infrastructure{}
			err = inventory.LoadInfrastructureFromString(contents, infra)
			So(err, ShouldBeNil)
			So(change.Subject, ShouldStartWith, "remove DUTs")
			So(infra.Servers, ShouldHaveLength, 2)

			var server *inventory.Server
			for _, s := range infra.Servers {
				if s.GetHostname() == serverID {
					server = s
					break
				}
			}
			So(server.DutUids, ShouldBeEmpty)
		})

		Convey("RemoveDutsFromdrone for the staging dut by name removes it from drone.", func() {
			req := &fleet.RemoveDutsFromDronesRequest{
				Removals: []*fleet.RemoveDutsFromDronesRequest_Item{{DutHostname: dutHostname}},
			}
			resp, err := tf.Inventory.RemoveDutsFromDrones(tf.C, req)
			So(err, ShouldBeNil)
			So(resp.Removed, ShouldHaveLength, 1)
			So(resp.Removed[0].DutId, ShouldEqual, dutID)

			So(tf.FakeGerrit.Changes, ShouldHaveLength, 1)
			change := tf.FakeGerrit.Changes[0]
			p := "data/skylab/server_db.textpb"
			So(change.Files, ShouldContainKey, p)

			contents := change.Files[p]
			infra := &inventory.Infrastructure{}
			err = inventory.LoadInfrastructureFromString(contents, infra)
			So(err, ShouldBeNil)
			So(change.Subject, ShouldStartWith, "remove DUTs")
			So(infra.Servers, ShouldHaveLength, 2)

			var server *inventory.Server
			for _, s := range infra.Servers {
				if s.GetHostname() == serverID {
					server = s
					break
				}
			}
			So(server.DutUids, ShouldBeEmpty)
		})

		Convey("RemoveDutsFromdrone for a nonexistant dut returns no results.", func() {
			req := &fleet.RemoveDutsFromDronesRequest{
				Removals: []*fleet.RemoveDutsFromDronesRequest_Item{{DutId: "foo"}},
			}
			resp, err := tf.Inventory.RemoveDutsFromDrones(tf.C, req)
			So(err, ShouldBeNil)
			So(resp.Removed, ShouldBeEmpty)
			So(resp.Url, ShouldEqual, "")
		})

		Convey("RemoveDutsFromdrone for prod dut returns no results.", func() {
			req := &fleet.RemoveDutsFromDronesRequest{
				Removals: []*fleet.RemoveDutsFromDronesRequest_Item{{DutId: wrongEnvDutID}},
			}
			resp, err := tf.Inventory.RemoveDutsFromDrones(tf.C, req)
			So(err, ShouldBeNil)
			So(resp.Removed, ShouldBeEmpty)
			So(resp.Url, ShouldEqual, "")
		})
	})
}

func TestRemoveSliceString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		input []string
		rem   string
		want  []string
	}{
		{
			name:  "from middle",
			input: []string{"a", "b", "c"},
			rem:   "b",
			want:  []string{"a", "c"},
		},
		{
			name:  "from end",
			input: []string{"a", "b", "c"},
			rem:   "c",
			want:  []string{"a", "b"},
		},
		{
			name:  "from beg",
			input: []string{"a", "b", "c"},
			rem:   "a",
			want:  []string{"b", "c"},
		},
		{
			name:  "missing",
			input: []string{"a", "b", "c"},
			rem:   "d",
			want:  []string{"a", "b", "c"},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := make([]string, len(c.input))
			copy(got, c.input)
			got = removeSliceString(got, c.rem)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("removeSliceString(%#v, %#v) = %#v; want %#v",
					c.input, c.rem, got, c.want)
			}
		})
	}
}
