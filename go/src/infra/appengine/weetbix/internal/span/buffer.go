// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package span

import (
	"fmt"
	"reflect"

	"cloud.google.com/go/spanner"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.chromium.org/luci/common/errors"

	"infra/appengine/weetbix/internal"
	"infra/appengine/weetbix/pbutil"
	pb "infra/appengine/weetbix/proto/v1"
)

// Value can be converted to a Spanner value.
// Typically if type T implements Value, then *T implements Ptr.
type Value interface {
	// ToSpanner returns a value of a type supported by Spanner client.
	ToSpanner() interface{}
}

// Ptr can be used a destination of reading a Spanner cell.
// Typically if type *T implements Ptr, then T implements Value.
type Ptr interface {
	// SpannerPtr returns to a pointer of a type supported by Spanner client.
	// SpannerPtr must use one of typed buffers in b.
	SpannerPtr(b *Buffer) interface{}
	// FromSpanner replaces Ptr value with the value in the typed buffer returned
	// by SpannerPtr.
	FromSpanner(b *Buffer) error
}

// Buffer can convert a value from a Spanner type to a Go type.
// Supported types:
//   - Value and Ptr
//   - string
//   - timestamppb.Timestamp
//   - pb.AnalyzedTestVariantStatus
//   - pb.Variant
//   - pb.StringPair
//   - proto.Message
//   - internal.VerdictStatus
type Buffer struct {
	NullString  spanner.NullString
	NullTime    spanner.NullTime
	Int64       int64
	StringSlice []string
	ByteSlice   []byte
	Int64Slice  []int64
}

// FromSpanner is a shortcut for (&Buffer{}).FromSpanner.
// Appropriate when FromSpanner is called only once, whereas Buffer is reusable
// throughout function.
func FromSpanner(row *spanner.Row, ptrs ...interface{}) error {
	return (&Buffer{}).FromSpanner(row, ptrs...)
}

// FromSpanner reads values from row to ptrs, converting types from Spanner
// to Go along the way.
func (b *Buffer) FromSpanner(row *spanner.Row, ptrs ...interface{}) error {
	if len(ptrs) != row.Size() {
		panic("len(ptrs) != row.Size()")
	}

	for i, goPtr := range ptrs {
		if err := b.fromSpanner(row, i, goPtr); err != nil {
			return err
		}
	}
	return nil
}

func (b *Buffer) fromSpanner(row *spanner.Row, col int, goPtr interface{}) error {
	b.StringSlice = b.StringSlice[:0]
	b.ByteSlice = b.ByteSlice[:0]
	b.Int64Slice = b.Int64Slice[:0]

	var spanPtr interface{}
	switch goPtr := goPtr.(type) {
	case Ptr:
		spanPtr = goPtr.SpannerPtr(b)
	case *string:
		spanPtr = &b.NullString
	case **timestamppb.Timestamp:
		spanPtr = &b.NullTime
	case *pb.AnalyzedTestVariantStatus:
		spanPtr = &b.Int64
	case **pb.Variant:
		spanPtr = &b.StringSlice
	case *[]*pb.StringPair:
		spanPtr = &b.StringSlice
	case proto.Message:
		spanPtr = &b.ByteSlice
	case *[]pb.AnalyzedTestVariantStatus:
		spanPtr = &b.Int64Slice
	case *internal.VerdictStatus:
		spanPtr = &b.Int64
	default:
		spanPtr = goPtr
	}

	if err := row.Column(col, spanPtr); err != nil {
		return errors.Annotate(err, "failed to read column %q", row.ColumnName(col)).Err()
	}

	if spanPtr == goPtr {
		return nil
	}

	var err error
	switch goPtr := goPtr.(type) {
	case Ptr:
		if err := goPtr.FromSpanner(b); err != nil {
			return err
		}

	case *string:
		*goPtr = ""
		if b.NullString.Valid {
			*goPtr = b.NullString.StringVal
		}

	case **timestamppb.Timestamp:
		*goPtr = nil
		if b.NullTime.Valid {
			*goPtr = pbutil.MustTimestampProto(b.NullTime.Time)
		}

	case *pb.AnalyzedTestVariantStatus:
		*goPtr = pb.AnalyzedTestVariantStatus(b.Int64)

	case **pb.Variant:
		if *goPtr, err = pbutil.VariantFromStrings(b.StringSlice); err != nil {
			// If it was written to Spanner, it should have been validated.
			panic(err)
		}

	case *[]*pb.StringPair:
		*goPtr = make([]*pb.StringPair, len(b.StringSlice))
		for i, p := range b.StringSlice {
			if (*goPtr)[i], err = pbutil.StringPairFromString(p); err != nil {
				// If it was written to Spanner, it should have been validated.
				panic(err)
			}
		}

	case proto.Message:
		if reflect.ValueOf(goPtr).IsNil() {
			return errors.Reason("nil pointer encountered").Err()
		}
		if err := proto.Unmarshal(b.ByteSlice, goPtr); err != nil {
			// If it was written to Spanner, it should have been validated.
			panic(err)
		}

	case *internal.VerdictStatus:
		*goPtr = internal.VerdictStatus(b.Int64)

	case *[]pb.AnalyzedTestVariantStatus:
		*goPtr = make([]pb.AnalyzedTestVariantStatus, len(b.Int64Slice))
		for i, p := range b.Int64Slice {
			(*goPtr)[i] = pb.AnalyzedTestVariantStatus(p)
		}

	default:
		panic(fmt.Sprintf("impossible %q", goPtr))
	}
	return nil
}

// ToSpanner converts values from Go types to Spanner types. In addition to
// supported types in FromSpanner, also supports []interface{} and
// map[string]interface{}.
func ToSpanner(v interface{}) interface{} {
	switch v := v.(type) {
	case Value:
		return v.ToSpanner()

	case *timestamppb.Timestamp:
		if v == nil {
			return spanner.NullTime{}
		}
		ret := spanner.NullTime{Valid: true}
		ret.Time = v.AsTime()
		// ptypes.Timestamp always returns a timestamp, even
		// if the returned err is non-nil, see its documentation.
		// The error is returned only if the timestamp violates its
		// own invariants, which will be caught on the attempt to
		// insert this value into Spanner.
		// This is consistent with the behavior of spanner.Insert() and
		// other mutating functions that ignore invalid time.Time
		// until it is time to apply the mutation.
		// Not returning an error here significantly simplifies usage
		// of this function and functions based on this one.
		return ret

	case pb.AnalyzedTestVariantStatus:
		return int64(v)

	case *pb.Variant:
		return pbutil.VariantToStrings(v)

	case []*pb.StringPair:
		return pbutil.StringPairsToStrings(v...)

	case proto.Message:
		if reflect.ValueOf(v).IsNil() {
			// Do not store empty messages.
			return []byte(nil)
		}

		ret, err := proto.Marshal(v)
		if err != nil {
			panic(err)
		}
		return ret

	case []interface{}:
		ret := make([]interface{}, len(v))
		for i, el := range v {
			ret[i] = ToSpanner(el)
		}
		return ret

	case map[string]interface{}:
		ret := make(map[string]interface{}, len(v))
		for key, el := range v {
			ret[key] = ToSpanner(el)
		}
		return ret

	case internal.VerdictStatus:
		return int64(v)

	case []pb.AnalyzedTestVariantStatus:
		spanPtr := make([]int64, len(v))
		for i, s := range v {
			spanPtr[i] = int64(s)
		}
		return spanPtr

	default:
		return v
	}
}

// ToSpannerMap converts a map of Go values to a map of Spanner values.
// See also ToSpanner.
func ToSpannerMap(values map[string]interface{}) map[string]interface{} {
	return ToSpanner(values).(map[string]interface{})
}

// UpdateMap is a shortcut for spanner.UpdateMap with ToSpannerMap applied to
// in.
func UpdateMap(table string, in map[string]interface{}) *spanner.Mutation {
	return spanner.UpdateMap(table, ToSpannerMap(in))
}

// InsertMap is a shortcut for spanner.InsertMap with ToSpannerMap applied to
// in.
func InsertMap(table string, in map[string]interface{}) *spanner.Mutation {
	return spanner.InsertMap(table, ToSpannerMap(in))
}

// InsertOrUpdateMap is a shortcut for spanner.InsertOrUpdateMap with ToSpannerMap applied to
// in.
func InsertOrUpdateMap(table string, in map[string]interface{}) *spanner.Mutation {
	return spanner.InsertOrUpdateMap(table, ToSpannerMap(in))
}
