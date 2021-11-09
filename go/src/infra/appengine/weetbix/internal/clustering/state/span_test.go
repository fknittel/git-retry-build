// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package state

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.chromium.org/luci/server/span"

	"infra/appengine/weetbix/internal/clustering"
	"infra/appengine/weetbix/internal/testutil"

	. "github.com/smartystreets/goconvey/convey"
	. "go.chromium.org/luci/common/testing/assertions"
)

func TestSpanner(t *testing.T) {
	Convey(`Create`, t, func() {
		ctx := testutil.SpannerTestContext(t)
		testCreate := func(e *Entry) error {
			_, err := span.ReadWriteTransaction(ctx, func(ctx context.Context) error {
				return Create(ctx, e)
			})
			return err
		}
		e := newEntry(100)
		Convey(`Valid`, func() {
			err := testCreate(e)
			So(err, ShouldBeNil)

			txn := span.Single(ctx)
			actual, err := Read(txn, e.Project, e.ChunkID)
			So(err, ShouldBeNil)

			// Check the LastUpdated time is set, but ignore it for
			// further comparisons.
			So(actual.LastUpdated, ShouldNotBeZeroValue)
			actual.LastUpdated = time.Time{}

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, e)
		})
		Convey(`Invalid`, func() {
			Convey(`Project missing`, func() {
				e.Project = ""
				err := testCreate(e)
				So(err, ShouldErrLike, `project "" is not valid`)
			})
			Convey(`Chunk ID missing`, func() {
				e.ChunkID = ""
				err := testCreate(e)
				So(err, ShouldErrLike, `chunk ID "" is not valid`)
			})
			Convey(`Partition Time missing`, func() {
				var t time.Time
				e.PartitionTime = t
				err := testCreate(e)
				So(err, ShouldErrLike, "partition time must be specified")
			})
			Convey(`Object ID missing`, func() {
				e.ObjectID = ""
				err := testCreate(e)
				So(err, ShouldErrLike, "object ID must be specified")
			})
			Convey(`Rule Version missing`, func() {
				var t time.Time
				e.Clustering.RulesVersion = t
				err := testCreate(e)
				So(err, ShouldErrLike, "rule version must be specified")
			})
			Convey(`Algorithms Version missing`, func() {
				e.Clustering.AlgorithmsVersion = 0
				err := testCreate(e)
				So(err, ShouldErrLike, "algorithms version must be specified")
			})
			Convey(`Clusters missing`, func() {
				e.Clustering.Clusters = nil
				err := testCreate(e)
				So(err, ShouldErrLike, "there must be clustered test results in the chunk")
			})
			Convey(`Algorithms invalid`, func() {
				Convey(`Empty algorithm`, func() {
					e.Clustering.Algorithms[""] = struct{}{}
					err := testCreate(e)
					So(err, ShouldErrLike, `algorithm "" is not valid`)
				})
				Convey("Algorithm invalid", func() {
					e.Clustering.Algorithms["!!!"] = struct{}{}
					err := testCreate(e)
					So(err, ShouldErrLike, `algorithm "!!!" is not valid`)
				})
			})
			Convey(`Clusters invalid`, func() {
				Convey(`Algorithm missing`, func() {
					e.Clustering.Clusters[1][1].Algorithm = ""
					err := testCreate(e)
					So(err, ShouldErrLike, `clusters: test result 1: cluster 1: cluster ID is not valid: algorithm not valid`)
				})
				Convey("Algorithm invalid", func() {
					e.Clustering.Clusters[1][1].Algorithm = "!!!"
					err := testCreate(e)
					So(err, ShouldErrLike, `clusters: test result 1: cluster 1: cluster ID is not valid: algorithm not valid`)
				})
				Convey("Algorithm not in algorithms set", func() {
					e.Clustering.Algorithms = map[string]struct{}{
						"alg-extra": {},
					}
					err := testCreate(e)
					So(err, ShouldErrLike, `a test result was clustered with an unregistered algorithm`)
				})
				Convey("ID missing", func() {
					e.Clustering.Clusters[1][1].ID = ""
					err := testCreate(e)
					So(err, ShouldErrLike, `clusters: test result 1: cluster 1: cluster ID is not valid: ID is empty`)
				})
			})
		})
	})
}

func newEntry(uniqifier int) *Entry {
	return &Entry{
		Project:       fmt.Sprintf("project-%v", uniqifier),
		ChunkID:       fmt.Sprintf("c%v", uniqifier),
		PartitionTime: time.Date(2030, 1, 1, 1, 1, 1, uniqifier, time.UTC),
		ObjectID:      "abcdef1234567890abcdef1234567890",
		Clustering: clustering.ClusterResults{
			AlgorithmsVersion: int64(uniqifier),
			RulesVersion:      time.Date(2025, 1, 1, 1, 1, 1, uniqifier, time.UTC),
			Algorithms: map[string]struct{}{
				fmt.Sprintf("alg-%v", uniqifier): {},
				"alg-extra":                      {},
			},
			Clusters: [][]*clustering.ClusterID{
				{
					{
						Algorithm: fmt.Sprintf("alg-%v", uniqifier),
						ID:        "00112233445566778899aabbccddeeff",
					},
				},
				{
					{
						Algorithm: fmt.Sprintf("alg-%v", uniqifier),
						ID:        "00112233445566778899aabbccddeeff",
					},
					{
						Algorithm: fmt.Sprintf("alg-%v", uniqifier),
						ID:        "22",
					},
				},
			},
		},
	}
}
