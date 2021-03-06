// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package controller

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "go.chromium.org/luci/common/testing/assertions"

	ufspb "infra/unifiedfleet/api/v1/models"
	"infra/unifiedfleet/app/model/configuration"
	. "infra/unifiedfleet/app/model/datastore"
	"infra/unifiedfleet/app/model/inventory"
	"infra/unifiedfleet/app/model/registration"
)

func mockChromePlatform(id, desc string) *ufspb.ChromePlatform {
	return &ufspb.ChromePlatform{
		Name:        id,
		Description: desc,
	}
}

func TestListChromePlatforms(t *testing.T) {
	t.Parallel()
	ctx := testingContext()
	chromePlatforms := make([]*ufspb.ChromePlatform, 0, 4)
	for i := 0; i < 4; i++ {
		chromePlatform1 := mockChromePlatform("", "Camera")
		chromePlatform1.Name = fmt.Sprintf("chromePlatform-%d", i)
		resp, _ := configuration.CreateChromePlatform(ctx, chromePlatform1)
		chromePlatforms = append(chromePlatforms, resp)
	}
	Convey("ListChromePlatforms", t, func() {
		Convey("List chromePlatforms - filter invalid", func() {
			_, _, err := ListChromePlatforms(ctx, 5, "", "machine=mx-1", false)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Failed to read filter for listing chromeplatforms")
		})

		Convey("ListChromePlatforms - Full listing - happy path", func() {
			resp, _, _ := ListChromePlatforms(ctx, 5, "", "", false)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResembleProto, chromePlatforms)
		})
	})
}

func TestDeleteChromePlatform(t *testing.T) {
	t.Parallel()
	ctx := testingContext()
	chromePlatform1 := mockChromePlatform("chromePlatform-1", "Camera")
	chromePlatform2 := mockChromePlatform("chromePlatform-2", "Camera")
	chromePlatform3 := mockChromePlatform("chromePlatform-3", "Sensor")
	Convey("DeleteChromePlatform", t, func() {
		Convey("Delete chromePlatform by existing ID with machine reference", func() {
			resp, cerr := configuration.CreateChromePlatform(ctx, chromePlatform1)
			So(cerr, ShouldBeNil)
			So(resp, ShouldResembleProto, chromePlatform1)

			chromeBrowserMachine1 := &ufspb.Machine{
				Name: "machine-1",
				Device: &ufspb.Machine_ChromeBrowserMachine{
					ChromeBrowserMachine: &ufspb.ChromeBrowserMachine{
						ChromePlatform: "chromePlatform-1",
					},
				},
			}
			mresp, merr := registration.CreateMachine(ctx, chromeBrowserMachine1)
			So(merr, ShouldBeNil)
			So(mresp, ShouldResembleProto, chromeBrowserMachine1)

			err := DeleteChromePlatform(ctx, "chromePlatform-1")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, CannotDelete)

			resp, cerr = configuration.GetChromePlatform(ctx, "chromePlatform-1")
			So(resp, ShouldNotBeNil)
			So(cerr, ShouldBeNil)
			So(resp, ShouldResembleProto, chromePlatform1)
		})
		Convey("Delete chromePlatform by existing ID with KVM reference", func() {
			resp, cerr := configuration.CreateChromePlatform(ctx, chromePlatform3)
			So(cerr, ShouldBeNil)
			So(resp, ShouldResembleProto, chromePlatform3)

			kvm1 := &ufspb.KVM{
				Name:           "kvm-1",
				ChromePlatform: "chromePlatform-3",
			}
			kresp, kerr := registration.CreateKVM(ctx, kvm1)
			So(kerr, ShouldBeNil)
			So(kresp, ShouldResembleProto, kvm1)

			err := DeleteChromePlatform(ctx, "chromePlatform-3")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, CannotDelete)

			resp, cerr = configuration.GetChromePlatform(ctx, "chromePlatform-3")
			So(resp, ShouldNotBeNil)
			So(cerr, ShouldBeNil)
			So(resp, ShouldResembleProto, chromePlatform3)
		})
		Convey("Delete chromePlatform successfully by existing ID without references", func() {
			resp, cerr := configuration.CreateChromePlatform(ctx, chromePlatform2)
			So(cerr, ShouldBeNil)
			So(resp, ShouldResembleProto, chromePlatform2)

			err := DeleteChromePlatform(ctx, "chromePlatform-2")
			So(err, ShouldBeNil)

			resp, cerr = configuration.GetChromePlatform(ctx, "chromePlatform-2")
			So(resp, ShouldBeNil)
			So(cerr, ShouldNotBeNil)
			So(cerr.Error(), ShouldContainSubstring, NotFound)
		})
	})
}

func TestUpdateChromePlatforms(t *testing.T) {
	t.Parallel()
	ctx := testingContext()
	chromePlatform1 := mockChromePlatform("chromePlatform-update", "Camera")
	chromePlatform1.Manufacturer = "fake"
	chromePlatform2 := mockChromePlatform("chromePlatform-update-2", "Camera")
	chromePlatform2.Manufacturer = "apple"
	configuration.BatchUpdateChromePlatforms(ctx, []*ufspb.ChromePlatform{chromePlatform1, chromePlatform2})
	Convey("UpdateChromePlatforms", t, func() {
		Convey("happy path", func() {
			p2 := mockChromePlatform("chromePlatform-update", "Camera")
			p2.Manufacturer = "non-fake"
			newP, err := UpdateChromePlatform(ctx, p2, nil)
			So(err, ShouldBeNil)
			So(newP, ShouldResembleProto, p2)
		})

		Convey("happy path with updating manufacturer", func() {
			inventory.CreateMachineLSE(ctx, &ufspb.MachineLSE{
				Name:         "platform-host",
				Hostname:     "platform-host",
				Manufacturer: "apple",
				Machines:     []string{"platform-machine"},
			})
			registration.CreateMachine(ctx, &ufspb.Machine{
				Name: "platform-machine",
				Device: &ufspb.Machine_ChromeBrowserMachine{
					ChromeBrowserMachine: &ufspb.ChromeBrowserMachine{
						ChromePlatform: "chromePlatform-update-2",
					},
				},
			})
			p2 := mockChromePlatform("chromePlatform-update-2", "Camera")
			p2.Manufacturer = "dell"
			newP, err := UpdateChromePlatform(ctx, p2, nil)
			So(err, ShouldBeNil)
			So(newP, ShouldResembleProto, p2)

			lse, err := inventory.GetMachineLSE(ctx, "platform-host")
			So(err, ShouldBeNil)
			So(lse.GetManufacturer(), ShouldEqual, "dell")
		})
	})
}

func TestBatchGetChromePlatforms(t *testing.T) {
	t.Parallel()
	ctx := testingContext()
	Convey("BatchGetChromePlatforms", t, func() {
		Convey("Batch get chrome platforms - happy path", func() {
			platforms := make([]*ufspb.ChromePlatform, 4)
			for i := 0; i < 4; i++ {
				platforms[i] = &ufspb.ChromePlatform{
					Name: fmt.Sprintf("platform-batchGet-%d", i),
				}
			}
			_, err := configuration.BatchUpdateChromePlatforms(ctx, platforms)
			So(err, ShouldBeNil)
			resp, err := configuration.BatchGetChromePlatforms(ctx, []string{"platform-batchGet-0", "platform-batchGet-1", "platform-batchGet-2", "platform-batchGet-3"})
			So(err, ShouldBeNil)
			So(resp, ShouldHaveLength, 4)
			So(resp, ShouldResembleProto, platforms)
		})
		Convey("Batch get chrome platforms  - missing id", func() {
			resp, err := configuration.BatchGetChromePlatforms(ctx, []string{"platform-batchGet-non-existing"})
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(err.Error(), ShouldContainSubstring, "platform-batchGet-non-existing")
		})
		Convey("Batch get chrome platforms  - empty input", func() {
			resp, err := configuration.BatchGetChromePlatforms(ctx, nil)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveLength, 0)

			input := make([]string, 0)
			resp, err = configuration.BatchGetChromePlatforms(ctx, input)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveLength, 0)
		})
	})
}
