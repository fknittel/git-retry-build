// Copyright 2022 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package gclient

import (
	"context"
	"infra/chromium/bootstrapper/cipd"
	fakecipd "infra/chromium/bootstrapper/fakes/cipd"
	. "infra/chromium/util"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "go.chromium.org/luci/common/testing/assertions"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	Convey("NewClient", t, func() {

		pkg := &fakecipd.Package{
			Refs:      map[string]string{},
			Instances: map[string]*fakecipd.PackageInstance{},
		}
		ctx = cipd.UseCipdClientFactory(ctx, fakecipd.Factory(map[string]*fakecipd.Package{
			depotToolsPackage: pkg,
		}))

		cipdRoot := t.TempDir()
		cipdClient, err := cipd.NewClient(ctx, cipdRoot)
		PanicOnError(err)

		Convey("fails if resolving depot tools version fails", func() {
			pkg.Refs[depotToolsPackageVersion] = ""

			client, err := NewClient(ctx, cipdClient)

			So(err, ShouldErrLike, "failed to resolve depot tools package version")
			So(client, ShouldBeNil)
		})

		Convey("fails if downloading depot tools package fails", func() {
			pkg.Refs[depotToolsPackageVersion] = "fake-instance"
			pkg.Instances["fake-instance"] = nil

			client, err := NewClient(ctx, cipdClient)

			So(err, ShouldErrLike, "failed to download/install depot tools")
			So(client, ShouldBeNil)
		})

		Convey("returns the path to gclient if successful", func() {
			client, err := NewClient(ctx, cipdClient)

			So(err, ShouldBeNil)
			So(client, ShouldNotBeNil)
			So(client.gclientPath, ShouldEqual, filepath.Join(cipdRoot, "depot-tools", "depot_tools", "gclient"))
		})

	})
}

func TestGetDep(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, _ := NewClientForTesting()

	Convey("getDep", t, func() {

		Convey("returns the revision for the specified path", func() {
			depsContents := `deps = {
				'foo': 'https://chromium.googlesource.com/foo.git@foo-revision',
			}`

			revision, err := client.GetDep(ctx, depsContents, "foo")

			So(err, ShouldBeNil)
			So(revision, ShouldEqual, "foo-revision")
		})

		Convey("fails for unknown path", func() {
			depsContents := `deps = {
				'foo': 'https://chromium.googlesource.com/foo.git@foo-revision',
			}`

			revision, err := client.GetDep(ctx, depsContents, "bar")

			So(err, ShouldErrLike, "Could not find any dependency called bar")
			So(revision, ShouldBeEmpty)
		})

	})
}
