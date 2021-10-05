// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package analyzedtestvariants

import (
	"context"

	"cloud.google.com/go/spanner"

	"go.chromium.org/luci/server/span"

	spanutil "infra/appengine/weetbix/internal/span"
	pb "infra/appengine/weetbix/proto/v1"
)

// Read reads AnalyzedTestVariant rows by keys.
func Read(ctx context.Context, ks spanner.KeySet, f func(*pb.AnalyzedTestVariant) error) error {
	fields := []string{"Realm", "TestId", "VariantHash", "Status"}
	var b spanutil.Buffer
	return span.Read(ctx, "AnalyzedTestVariants", ks, fields).Do(
		func(row *spanner.Row) error {
			tv := &pb.AnalyzedTestVariant{}
			if err := b.FromSpanner(row, &tv.Realm, &tv.TestId, &tv.VariantHash, &tv.Status); err != nil {
				return err
			}
			return f(tv)
		},
	)
}
