// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: infra/appengine/weetbix/proto/v1/predicate.proto

package weetbixpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Represents a function Variant -> bool.
type VariantPredicate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Predicate:
	//	*VariantPredicate_Equals
	//	*VariantPredicate_Contains
	Predicate isVariantPredicate_Predicate `protobuf_oneof:"predicate"`
}

func (x *VariantPredicate) Reset() {
	*x = VariantPredicate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_predicate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VariantPredicate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VariantPredicate) ProtoMessage() {}

func (x *VariantPredicate) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_predicate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VariantPredicate.ProtoReflect.Descriptor instead.
func (*VariantPredicate) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDescGZIP(), []int{0}
}

func (m *VariantPredicate) GetPredicate() isVariantPredicate_Predicate {
	if m != nil {
		return m.Predicate
	}
	return nil
}

func (x *VariantPredicate) GetEquals() *Variant {
	if x, ok := x.GetPredicate().(*VariantPredicate_Equals); ok {
		return x.Equals
	}
	return nil
}

func (x *VariantPredicate) GetContains() *Variant {
	if x, ok := x.GetPredicate().(*VariantPredicate_Contains); ok {
		return x.Contains
	}
	return nil
}

type isVariantPredicate_Predicate interface {
	isVariantPredicate_Predicate()
}

type VariantPredicate_Equals struct {
	// A variant must be equal this definition exactly.
	Equals *Variant `protobuf:"bytes,1,opt,name=equals,proto3,oneof"`
}

type VariantPredicate_Contains struct {
	// A variant's key-value pairs must contain those in this one.
	Contains *Variant `protobuf:"bytes,2,opt,name=contains,proto3,oneof"`
}

func (*VariantPredicate_Equals) isVariantPredicate_Predicate() {}

func (*VariantPredicate_Contains) isVariantPredicate_Predicate() {}

// Represents a function AnalyzedTestVariant -> bool.
type AnalyzedTestVariantPredicate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A test variant must have a test id matching this regular expression
	// entirely, i.e. the expression is implicitly wrapped with ^ and $.
	TestIdRegexp string `protobuf:"bytes,1,opt,name=test_id_regexp,json=testIdRegexp,proto3" json:"test_id_regexp,omitempty"`
	// A test variant must have a variant satisfying this predicate.
	Variant *VariantPredicate `protobuf:"bytes,2,opt,name=variant,proto3" json:"variant,omitempty"`
	// A test variant must have this status.
	Status AnalyzedTestVariantStatus `protobuf:"varint,3,opt,name=status,proto3,enum=weetbix.v1.AnalyzedTestVariantStatus" json:"status,omitempty"`
}

func (x *AnalyzedTestVariantPredicate) Reset() {
	*x = AnalyzedTestVariantPredicate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_predicate_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzedTestVariantPredicate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzedTestVariantPredicate) ProtoMessage() {}

func (x *AnalyzedTestVariantPredicate) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_predicate_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzedTestVariantPredicate.ProtoReflect.Descriptor instead.
func (*AnalyzedTestVariantPredicate) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDescGZIP(), []int{1}
}

func (x *AnalyzedTestVariantPredicate) GetTestIdRegexp() string {
	if x != nil {
		return x.TestIdRegexp
	}
	return ""
}

func (x *AnalyzedTestVariantPredicate) GetVariant() *VariantPredicate {
	if x != nil {
		return x.Variant
	}
	return nil
}

func (x *AnalyzedTestVariantPredicate) GetStatus() AnalyzedTestVariantStatus {
	if x != nil {
		return x.Status
	}
	return AnalyzedTestVariantStatus_STATUS_UNSPECIFIED
}

var File_infra_appengine_weetbix_proto_v1_predicate_proto protoreflect.FileDescriptor

var file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDesc = []byte{
	0x0a, 0x30, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2f, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x76, 0x31, 0x2f, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x1a, 0x3c,
	0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f,
	0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x64, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x76,
	0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2d, 0x69, 0x6e,
	0x66, 0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x77, 0x65,
	0x65, 0x74, 0x62, 0x69, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x81, 0x01, 0x0a, 0x10,
	0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x12, 0x2d, 0x0a, 0x06, 0x65, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61,
	0x72, 0x69, 0x61, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x06, 0x65, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x12,
	0x31, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x56,
	0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x73, 0x42, 0x0b, 0x0a, 0x09, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x22,
	0xbb, 0x01, 0x0a, 0x1c, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x64, 0x54, 0x65, 0x73, 0x74,
	0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x12, 0x24, 0x0a, 0x0e, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x5f, 0x72, 0x65, 0x67, 0x65,
	0x78, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x65, 0x73, 0x74, 0x49, 0x64,
	0x52, 0x65, 0x67, 0x65, 0x78, 0x70, 0x12, 0x36, 0x0a, 0x07, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69,
	0x78, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x50, 0x72, 0x65, 0x64,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x07, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x12, 0x3d,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25,
	0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6e, 0x61, 0x6c,
	0x79, 0x7a, 0x65, 0x64, 0x54, 0x65, 0x73, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x2c, 0x5a,
	0x2a, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65,
	0x2f, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76,
	0x31, 0x3b, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDescOnce sync.Once
	file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDescData = file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDesc
)

func file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDescGZIP() []byte {
	file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDescOnce.Do(func() {
		file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDescData)
	})
	return file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDescData
}

var file_infra_appengine_weetbix_proto_v1_predicate_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_infra_appengine_weetbix_proto_v1_predicate_proto_goTypes = []interface{}{
	(*VariantPredicate)(nil),             // 0: weetbix.v1.VariantPredicate
	(*AnalyzedTestVariantPredicate)(nil), // 1: weetbix.v1.AnalyzedTestVariantPredicate
	(*Variant)(nil),                      // 2: weetbix.v1.Variant
	(AnalyzedTestVariantStatus)(0),       // 3: weetbix.v1.AnalyzedTestVariantStatus
}
var file_infra_appengine_weetbix_proto_v1_predicate_proto_depIdxs = []int32{
	2, // 0: weetbix.v1.VariantPredicate.equals:type_name -> weetbix.v1.Variant
	2, // 1: weetbix.v1.VariantPredicate.contains:type_name -> weetbix.v1.Variant
	0, // 2: weetbix.v1.AnalyzedTestVariantPredicate.variant:type_name -> weetbix.v1.VariantPredicate
	3, // 3: weetbix.v1.AnalyzedTestVariantPredicate.status:type_name -> weetbix.v1.AnalyzedTestVariantStatus
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_infra_appengine_weetbix_proto_v1_predicate_proto_init() }
func file_infra_appengine_weetbix_proto_v1_predicate_proto_init() {
	if File_infra_appengine_weetbix_proto_v1_predicate_proto != nil {
		return
	}
	file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_init()
	file_infra_appengine_weetbix_proto_v1_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_infra_appengine_weetbix_proto_v1_predicate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VariantPredicate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_infra_appengine_weetbix_proto_v1_predicate_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzedTestVariantPredicate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_infra_appengine_weetbix_proto_v1_predicate_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*VariantPredicate_Equals)(nil),
		(*VariantPredicate_Contains)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_appengine_weetbix_proto_v1_predicate_proto_goTypes,
		DependencyIndexes: file_infra_appengine_weetbix_proto_v1_predicate_proto_depIdxs,
		MessageInfos:      file_infra_appengine_weetbix_proto_v1_predicate_proto_msgTypes,
	}.Build()
	File_infra_appengine_weetbix_proto_v1_predicate_proto = out.File
	file_infra_appengine_weetbix_proto_v1_predicate_proto_rawDesc = nil
	file_infra_appengine_weetbix_proto_v1_predicate_proto_goTypes = nil
	file_infra_appengine_weetbix_proto_v1_predicate_proto_depIdxs = nil
}
