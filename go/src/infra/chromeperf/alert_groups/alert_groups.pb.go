// Copyright 2021 The Chromium Authors.
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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: infra/chromeperf/alert_groups/alert_groups.proto

package alert_groups

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type AlertGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// AlertGroup id in Datastore.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *AlertGroup) Reset() {
	*x = AlertGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_chromeperf_alert_groups_alert_groups_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlertGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlertGroup) ProtoMessage() {}

func (x *AlertGroup) ProtoReflect() protoreflect.Message {
	mi := &file_infra_chromeperf_alert_groups_alert_groups_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlertGroup.ProtoReflect.Descriptor instead.
func (*AlertGroup) Descriptor() ([]byte, []int) {
	return file_infra_chromeperf_alert_groups_alert_groups_proto_rawDescGZIP(), []int{0}
}

func (x *AlertGroup) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type MergeAlertGroupsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of groups that will be merged into the destination_group.
	// The source groups will be deleted after merge is complete.
	SourceGroups []*AlertGroup `protobuf:"bytes,1,rep,name=source_groups,json=sourceGroups,proto3" json:"source_groups,omitempty"`
	// Group that will contain all the anomalies from the source groups.
	DestinationGroup *AlertGroup `protobuf:"bytes,2,opt,name=destination_group,json=destinationGroup,proto3" json:"destination_group,omitempty"`
}

func (x *MergeAlertGroupsRequest) Reset() {
	*x = MergeAlertGroupsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_chromeperf_alert_groups_alert_groups_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MergeAlertGroupsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MergeAlertGroupsRequest) ProtoMessage() {}

func (x *MergeAlertGroupsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_infra_chromeperf_alert_groups_alert_groups_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MergeAlertGroupsRequest.ProtoReflect.Descriptor instead.
func (*MergeAlertGroupsRequest) Descriptor() ([]byte, []int) {
	return file_infra_chromeperf_alert_groups_alert_groups_proto_rawDescGZIP(), []int{1}
}

func (x *MergeAlertGroupsRequest) GetSourceGroups() []*AlertGroup {
	if x != nil {
		return x.SourceGroups
	}
	return nil
}

func (x *MergeAlertGroupsRequest) GetDestinationGroup() *AlertGroup {
	if x != nil {
		return x.DestinationGroup
	}
	return nil
}

var File_infra_chromeperf_alert_groups_alert_groups_proto protoreflect.FileDescriptor

var file_infra_chromeperf_alert_groups_alert_groups_proto_rawDesc = []byte{
	0x0a, 0x30, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x70, 0x65,
	0x72, 0x66, 0x2f, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f,
	0x61, 0x6c, 0x65, 0x72, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0c, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x25, 0x0a, 0x0a, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x17, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0xa9, 0x01, 0x0a, 0x17, 0x4d, 0x65, 0x72, 0x67, 0x65,
	0x41, 0x6c, 0x65, 0x72, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x42, 0x0a, 0x0d, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x6c, 0x65, 0x72,
	0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x0c, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x4a, 0x0a, 0x11, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73,
	0x2e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x42, 0x03, 0xe0, 0x41, 0x02,
	0x52, 0x10, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x32, 0x82, 0x01, 0x0a, 0x0b, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x73, 0x12, 0x73, 0x0a, 0x10, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x41, 0x6c, 0x65, 0x72, 0x74,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x25, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x5f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x41, 0x6c, 0x65, 0x72, 0x74,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e,
	0x61, 0x6c, 0x65, 0x72, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e, 0x41, 0x6c, 0x65,
	0x72, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x22,
	0x13, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x5f, 0x61, 0x6e, 0x6f, 0x6d, 0x61,
	0x6c, 0x69, 0x65, 0x73, 0x3a, 0x01, 0x2a, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_chromeperf_alert_groups_alert_groups_proto_rawDescOnce sync.Once
	file_infra_chromeperf_alert_groups_alert_groups_proto_rawDescData = file_infra_chromeperf_alert_groups_alert_groups_proto_rawDesc
)

func file_infra_chromeperf_alert_groups_alert_groups_proto_rawDescGZIP() []byte {
	file_infra_chromeperf_alert_groups_alert_groups_proto_rawDescOnce.Do(func() {
		file_infra_chromeperf_alert_groups_alert_groups_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_chromeperf_alert_groups_alert_groups_proto_rawDescData)
	})
	return file_infra_chromeperf_alert_groups_alert_groups_proto_rawDescData
}

var file_infra_chromeperf_alert_groups_alert_groups_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_infra_chromeperf_alert_groups_alert_groups_proto_goTypes = []interface{}{
	(*AlertGroup)(nil),              // 0: alert_groups.AlertGroup
	(*MergeAlertGroupsRequest)(nil), // 1: alert_groups.MergeAlertGroupsRequest
}
var file_infra_chromeperf_alert_groups_alert_groups_proto_depIdxs = []int32{
	0, // 0: alert_groups.MergeAlertGroupsRequest.source_groups:type_name -> alert_groups.AlertGroup
	0, // 1: alert_groups.MergeAlertGroupsRequest.destination_group:type_name -> alert_groups.AlertGroup
	1, // 2: alert_groups.AlertGroups.MergeAlertGroups:input_type -> alert_groups.MergeAlertGroupsRequest
	0, // 3: alert_groups.AlertGroups.MergeAlertGroups:output_type -> alert_groups.AlertGroup
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_infra_chromeperf_alert_groups_alert_groups_proto_init() }
func file_infra_chromeperf_alert_groups_alert_groups_proto_init() {
	if File_infra_chromeperf_alert_groups_alert_groups_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_infra_chromeperf_alert_groups_alert_groups_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlertGroup); i {
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
		file_infra_chromeperf_alert_groups_alert_groups_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MergeAlertGroupsRequest); i {
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
			RawDescriptor: file_infra_chromeperf_alert_groups_alert_groups_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_infra_chromeperf_alert_groups_alert_groups_proto_goTypes,
		DependencyIndexes: file_infra_chromeperf_alert_groups_alert_groups_proto_depIdxs,
		MessageInfos:      file_infra_chromeperf_alert_groups_alert_groups_proto_msgTypes,
	}.Build()
	File_infra_chromeperf_alert_groups_alert_groups_proto = out.File
	file_infra_chromeperf_alert_groups_alert_groups_proto_rawDesc = nil
	file_infra_chromeperf_alert_groups_alert_groups_proto_goTypes = nil
	file_infra_chromeperf_alert_groups_alert_groups_proto_depIdxs = nil
}
