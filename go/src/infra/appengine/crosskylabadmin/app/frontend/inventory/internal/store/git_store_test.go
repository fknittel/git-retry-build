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

package store

import (
	"infra/appengine/crosskylabadmin/app/config"
	"infra/appengine/crosskylabadmin/app/frontend/inventory/internal/fakes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.chromium.org/luci/appengine/gaetesting"
)

func TestStoreValidity(t *testing.T) {
	Convey("With 1 known DUT", t, func() {
		ctx := gaetesting.TestingContextWithAppID("dev~infra-crosskylabadmin")
		ctx = config.Use(ctx, &config.Config{
			AccessGroup: "fake-access-group",
			Inventory: &config.Inventory{
				GitilesHost:            "some-gitiles-host",
				GerritHost:             "some-gerrit-host",
				Project:                "some-project",
				Branch:                 "master",
				LabDataPath:            "data/skylab/lab.textpb",
				InfrastructureDataPath: "data/skylab/server_db.textpb",
				Environment:            "ENVIRONMENT_STAGING",
			},
		})
		gerritC := &fakes.GerritClient{}
		gitilesC := fakes.NewGitilesClient()

		err := gitilesC.AddArchive(config.Get(ctx).Inventory, []byte{}, []byte{})
		So(err, ShouldBeNil)

		Convey("store initially contains no data", func() {
			store := NewGitStore(gerritC, gitilesC)
			So(store.Lab, ShouldBeNil)

			Convey("and initial Commit() fails", func() {
				_, err := store.Commit(ctx, "no reason")
				So(err, ShouldNotBeNil)
			})

			Convey("on Refresh(), store obtains data", func() {
				err := store.Refresh(ctx)
				So(err, ShouldBeNil)
				So(store.Lab, ShouldNotBeNil)

				Convey("on Commit(), store is flushed", func() {
					_, err := store.Commit(ctx, "no reason")
					So(err, ShouldBeNil)
					So(store.Lab, ShouldBeNil)

					Convey("and invalidated", func() {
						_, err := store.Commit(ctx, "no reason")
						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}
