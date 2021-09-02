// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: infra/unifiedfleet/api/v1/models/chromeos/device/variant_id.proto

package ufspb

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

// Globally unique identifier.
type VariantId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required. Source: 'mosys platform sku', aka Device-SKU.
	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *VariantId) Reset() {
	*x = VariantId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VariantId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VariantId) ProtoMessage() {}

func (x *VariantId) ProtoReflect() protoreflect.Message {
	mi := &file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VariantId.ProtoReflect.Descriptor instead.
func (*VariantId) Descriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDescGZIP(), []int{0}
}

func (x *VariantId) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto protoreflect.FileDescriptor

var file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDesc = []byte{
	0x0a, 0x41, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66,
	0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x2f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x2a, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65,
	0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x22,
	0x21, 0x0a, 0x09, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x42, 0x38, 0x5a, 0x36, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x3b, 0x75, 0x66, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDescOnce sync.Once
	file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDescData = file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDesc
)

func file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDescGZIP() []byte {
	file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDescOnce.Do(func() {
		file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDescData)
	})
	return file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDescData
}

var file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_goTypes = []interface{}{
	(*VariantId)(nil), // 0: unifiedfleet.api.v1.models.chromeos.device.VariantId
}
var file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_init() }
func file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_init() {
	if File_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VariantId); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_goTypes,
		DependencyIndexes: file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_depIdxs,
		MessageInfos:      file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_msgTypes,
	}.Build()
	File_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto = out.File
	file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_rawDesc = nil
	file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_goTypes = nil
	file_infra_unifiedfleet_api_v1_models_chromeos_device_variant_id_proto_depIdxs = nil
}
