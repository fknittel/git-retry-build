// Copyright 2020 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// This file defines protobufs for features and related business
// objects, e.g., hotlists.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: api/v3/api_proto/permission_objects.proto

package api_proto

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// All possible permissions on the Monorail site.
// Next available tag: 6
type Permission int32

const (
	// Default value. This value is unused.
	Permission_PERMISSION_UNSPECIFIED Permission = 0
	// The permission needed to add and remove issues from a hotlist.
	Permission_HOTLIST_EDIT Permission = 1
	// The permission needed to delete a hotlist or change hotlist
	// settings/members.
	Permission_HOTLIST_ADMINISTER Permission = 2
	// The permission needed to edit an issue.
	Permission_ISSUE_EDIT Permission = 3
	// The permission needed to edit a custom field definition.
	Permission_FIELD_DEF_EDIT Permission = 4
	// The permission needed to edit the value of a custom field.
	// More permissions will be required in the specific issue
	// where the user plans to edit that value, e.g. ISSUE_EDIT.
	Permission_FIELD_DEF_VALUE_EDIT Permission = 5
)

// Enum value maps for Permission.
var (
	Permission_name = map[int32]string{
		0: "PERMISSION_UNSPECIFIED",
		1: "HOTLIST_EDIT",
		2: "HOTLIST_ADMINISTER",
		3: "ISSUE_EDIT",
		4: "FIELD_DEF_EDIT",
		5: "FIELD_DEF_VALUE_EDIT",
	}
	Permission_value = map[string]int32{
		"PERMISSION_UNSPECIFIED": 0,
		"HOTLIST_EDIT":           1,
		"HOTLIST_ADMINISTER":     2,
		"ISSUE_EDIT":             3,
		"FIELD_DEF_EDIT":         4,
		"FIELD_DEF_VALUE_EDIT":   5,
	}
)

func (x Permission) Enum() *Permission {
	p := new(Permission)
	*p = x
	return p
}

func (x Permission) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Permission) Descriptor() protoreflect.EnumDescriptor {
	return file_api_v3_api_proto_permission_objects_proto_enumTypes[0].Descriptor()
}

func (Permission) Type() protoreflect.EnumType {
	return &file_api_v3_api_proto_permission_objects_proto_enumTypes[0]
}

func (x Permission) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Permission.Descriptor instead.
func (Permission) EnumDescriptor() ([]byte, []int) {
	return file_api_v3_api_proto_permission_objects_proto_rawDescGZIP(), []int{0}
}

// The set of a user's permissions for a single resource.
// Next available tag: 3
type PermissionSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the resource `permissions` applies to.
	Resource string `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
	// All the permissions a user has for `resource`.
	Permissions []Permission `protobuf:"varint,2,rep,packed,name=permissions,proto3,enum=monorail.v3.Permission" json:"permissions,omitempty"`
}

func (x *PermissionSet) Reset() {
	*x = PermissionSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v3_api_proto_permission_objects_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PermissionSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PermissionSet) ProtoMessage() {}

func (x *PermissionSet) ProtoReflect() protoreflect.Message {
	mi := &file_api_v3_api_proto_permission_objects_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PermissionSet.ProtoReflect.Descriptor instead.
func (*PermissionSet) Descriptor() ([]byte, []int) {
	return file_api_v3_api_proto_permission_objects_proto_rawDescGZIP(), []int{0}
}

func (x *PermissionSet) GetResource() string {
	if x != nil {
		return x.Resource
	}
	return ""
}

func (x *PermissionSet) GetPermissions() []Permission {
	if x != nil {
		return x.Permissions
	}
	return nil
}

var File_api_v3_api_proto_permission_objects_proto protoreflect.FileDescriptor

var file_api_v3_api_proto_permission_objects_proto_rawDesc = []byte{
	0x0a, 0x29, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x33, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x6d, 0x6f, 0x6e,
	0x6f, 0x72, 0x61, 0x69, 0x6c, 0x2e, 0x76, 0x33, 0x22, 0x66, 0x0a, 0x0d, 0x50, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x39, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x6d, 0x6f, 0x6e,
	0x6f, 0x72, 0x61, 0x69, 0x6c, 0x2e, 0x76, 0x33, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2a, 0x90, 0x01, 0x0a, 0x0a, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x1a, 0x0a, 0x16, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x48,
	0x4f, 0x54, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x45, 0x44, 0x49, 0x54, 0x10, 0x01, 0x12, 0x16, 0x0a,
	0x12, 0x48, 0x4f, 0x54, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x49, 0x53,
	0x54, 0x45, 0x52, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x45,
	0x44, 0x49, 0x54, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x44,
	0x45, 0x46, 0x5f, 0x45, 0x44, 0x49, 0x54, 0x10, 0x04, 0x12, 0x18, 0x0a, 0x14, 0x46, 0x49, 0x45,
	0x4c, 0x44, 0x5f, 0x44, 0x45, 0x46, 0x5f, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x5f, 0x45, 0x44, 0x49,
	0x54, 0x10, 0x05, 0x42, 0x12, 0x5a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x33, 0x2f, 0x61, 0x70,
	0x69, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v3_api_proto_permission_objects_proto_rawDescOnce sync.Once
	file_api_v3_api_proto_permission_objects_proto_rawDescData = file_api_v3_api_proto_permission_objects_proto_rawDesc
)

func file_api_v3_api_proto_permission_objects_proto_rawDescGZIP() []byte {
	file_api_v3_api_proto_permission_objects_proto_rawDescOnce.Do(func() {
		file_api_v3_api_proto_permission_objects_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v3_api_proto_permission_objects_proto_rawDescData)
	})
	return file_api_v3_api_proto_permission_objects_proto_rawDescData
}

var file_api_v3_api_proto_permission_objects_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_v3_api_proto_permission_objects_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_v3_api_proto_permission_objects_proto_goTypes = []interface{}{
	(Permission)(0),       // 0: monorail.v3.Permission
	(*PermissionSet)(nil), // 1: monorail.v3.PermissionSet
}
var file_api_v3_api_proto_permission_objects_proto_depIdxs = []int32{
	0, // 0: monorail.v3.PermissionSet.permissions:type_name -> monorail.v3.Permission
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_v3_api_proto_permission_objects_proto_init() }
func file_api_v3_api_proto_permission_objects_proto_init() {
	if File_api_v3_api_proto_permission_objects_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v3_api_proto_permission_objects_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PermissionSet); i {
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
			RawDescriptor: file_api_v3_api_proto_permission_objects_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v3_api_proto_permission_objects_proto_goTypes,
		DependencyIndexes: file_api_v3_api_proto_permission_objects_proto_depIdxs,
		EnumInfos:         file_api_v3_api_proto_permission_objects_proto_enumTypes,
		MessageInfos:      file_api_v3_api_proto_permission_objects_proto_msgTypes,
	}.Build()
	File_api_v3_api_proto_permission_objects_proto = out.File
	file_api_v3_api_proto_permission_objects_proto_rawDesc = nil
	file_api_v3_api_proto_permission_objects_proto_goTypes = nil
	file_api_v3_api_proto_permission_objects_proto_depIdxs = nil
}
