// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.0
// source: infra/cros/cmd/labpack/steps/steps.proto

package stepspb

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

// LabpackInput represents list of input parameters.
type LabpackInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unit name represents some device setup against which running the task.
	UnitName string `protobuf:"bytes,1,opt,name=unit_name,json=unitName,proto3" json:"unit_name,omitempty"`
	// Task name running against unit.
	TaskName string `protobuf:"bytes,2,opt,name=task_name,json=taskName,proto3" json:"task_name,omitempty"`
	// Enable recovery tells if recovery actions are enabled.
	EnableRecovery bool `protobuf:"varint,3,opt,name=enable_recovery,json=enableRecovery,proto3" json:"enable_recovery,omitempty"`
	// Update inventory tells if process ellow update inventory during execution.
	UpdateInventory bool `protobuf:"varint,4,opt,name=update_inventory,json=updateInventory,proto3" json:"update_inventory,omitempty"`
	// Admin service path to initialie local TLW.
	AdminService string `protobuf:"bytes,5,opt,name=admin_service,json=adminService,proto3" json:"admin_service,omitempty"`
	// Inventory service path to initialie local TLW.
	InventoryService string `protobuf:"bytes,6,opt,name=inventory_service,json=inventoryService,proto3" json:"inventory_service,omitempty"`
	// Do not use stepper during execution.
	NoStepper bool `protobuf:"varint,7,opt,name=no_stepper,json=noStepper,proto3" json:"no_stepper,omitempty"`
	// Do not use metrics during execution.
	NoMetrics bool `protobuf:"varint,9,opt,name=no_metrics,json=noMetrics,proto3" json:"no_metrics,omitempty"`
	// Custom configuration.
	Configuration string `protobuf:"bytes,8,opt,name=configuration,proto3" json:"configuration,omitempty"`
}

func (x *LabpackInput) Reset() {
	*x = LabpackInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_cros_cmd_labpack_steps_steps_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabpackInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabpackInput) ProtoMessage() {}

func (x *LabpackInput) ProtoReflect() protoreflect.Message {
	mi := &file_infra_cros_cmd_labpack_steps_steps_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabpackInput.ProtoReflect.Descriptor instead.
func (*LabpackInput) Descriptor() ([]byte, []int) {
	return file_infra_cros_cmd_labpack_steps_steps_proto_rawDescGZIP(), []int{0}
}

func (x *LabpackInput) GetUnitName() string {
	if x != nil {
		return x.UnitName
	}
	return ""
}

func (x *LabpackInput) GetTaskName() string {
	if x != nil {
		return x.TaskName
	}
	return ""
}

func (x *LabpackInput) GetEnableRecovery() bool {
	if x != nil {
		return x.EnableRecovery
	}
	return false
}

func (x *LabpackInput) GetUpdateInventory() bool {
	if x != nil {
		return x.UpdateInventory
	}
	return false
}

func (x *LabpackInput) GetAdminService() string {
	if x != nil {
		return x.AdminService
	}
	return ""
}

func (x *LabpackInput) GetInventoryService() string {
	if x != nil {
		return x.InventoryService
	}
	return ""
}

func (x *LabpackInput) GetNoStepper() bool {
	if x != nil {
		return x.NoStepper
	}
	return false
}

func (x *LabpackInput) GetNoMetrics() bool {
	if x != nil {
		return x.NoMetrics
	}
	return false
}

func (x *LabpackInput) GetConfiguration() string {
	if x != nil {
		return x.Configuration
	}
	return ""
}

// LabpackResponse represents result of execution the task on unit.
type LabpackResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	// Tells what was the reason of failure.
	FailReason string `protobuf:"bytes,2,opt,name=fail_reason,json=failReason,proto3" json:"fail_reason,omitempty"`
}

func (x *LabpackResponse) Reset() {
	*x = LabpackResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_cros_cmd_labpack_steps_steps_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabpackResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabpackResponse) ProtoMessage() {}

func (x *LabpackResponse) ProtoReflect() protoreflect.Message {
	mi := &file_infra_cros_cmd_labpack_steps_steps_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabpackResponse.ProtoReflect.Descriptor instead.
func (*LabpackResponse) Descriptor() ([]byte, []int) {
	return file_infra_cros_cmd_labpack_steps_steps_proto_rawDescGZIP(), []int{1}
}

func (x *LabpackResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *LabpackResponse) GetFailReason() string {
	if x != nil {
		return x.FailReason
	}
	return ""
}

var File_infra_cros_cmd_labpack_steps_steps_proto protoreflect.FileDescriptor

var file_infra_cros_cmd_labpack_steps_steps_proto_rawDesc = []byte{
	0x0a, 0x28, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x2f, 0x63, 0x6d, 0x64,
	0x2f, 0x6c, 0x61, 0x62, 0x70, 0x61, 0x63, 0x6b, 0x2f, 0x73, 0x74, 0x65, 0x70, 0x73, 0x2f, 0x73,
	0x74, 0x65, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73, 0x74, 0x65, 0x70,
	0x73, 0x22, 0xd2, 0x02, 0x0a, 0x0c, 0x4c, 0x61, 0x62, 0x70, 0x61, 0x63, 0x6b, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x6e, 0x69, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x74, 0x61, 0x73, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f,
	0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x12, 0x29, 0x0a, 0x10, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f,
	0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79,
	0x12, 0x23, 0x0a, 0x0d, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x10, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x6f, 0x5f, 0x73, 0x74, 0x65, 0x70, 0x70, 0x65, 0x72,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x6e, 0x6f, 0x53, 0x74, 0x65, 0x70, 0x70, 0x65,
	0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x6e, 0x6f, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x4c, 0x0a, 0x0f, 0x4c, 0x61, 0x62, 0x70, 0x61, 0x63,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x66, 0x61, 0x69, 0x6c, 0x5f, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x66, 0x61, 0x69, 0x6c, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x42, 0x26, 0x5a, 0x24, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x63, 0x72,
	0x6f, 0x73, 0x2f, 0x63, 0x6d, 0x64, 0x2f, 0x6c, 0x61, 0x62, 0x70, 0x61, 0x63, 0x6b, 0x2f, 0x73,
	0x74, 0x65, 0x70, 0x73, 0x3b, 0x73, 0x74, 0x65, 0x70, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_cros_cmd_labpack_steps_steps_proto_rawDescOnce sync.Once
	file_infra_cros_cmd_labpack_steps_steps_proto_rawDescData = file_infra_cros_cmd_labpack_steps_steps_proto_rawDesc
)

func file_infra_cros_cmd_labpack_steps_steps_proto_rawDescGZIP() []byte {
	file_infra_cros_cmd_labpack_steps_steps_proto_rawDescOnce.Do(func() {
		file_infra_cros_cmd_labpack_steps_steps_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_cros_cmd_labpack_steps_steps_proto_rawDescData)
	})
	return file_infra_cros_cmd_labpack_steps_steps_proto_rawDescData
}

var file_infra_cros_cmd_labpack_steps_steps_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_infra_cros_cmd_labpack_steps_steps_proto_goTypes = []interface{}{
	(*LabpackInput)(nil),    // 0: steps.LabpackInput
	(*LabpackResponse)(nil), // 1: steps.LabpackResponse
}
var file_infra_cros_cmd_labpack_steps_steps_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_infra_cros_cmd_labpack_steps_steps_proto_init() }
func file_infra_cros_cmd_labpack_steps_steps_proto_init() {
	if File_infra_cros_cmd_labpack_steps_steps_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_infra_cros_cmd_labpack_steps_steps_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabpackInput); i {
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
		file_infra_cros_cmd_labpack_steps_steps_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabpackResponse); i {
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
			RawDescriptor: file_infra_cros_cmd_labpack_steps_steps_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_cros_cmd_labpack_steps_steps_proto_goTypes,
		DependencyIndexes: file_infra_cros_cmd_labpack_steps_steps_proto_depIdxs,
		MessageInfos:      file_infra_cros_cmd_labpack_steps_steps_proto_msgTypes,
	}.Build()
	File_infra_cros_cmd_labpack_steps_steps_proto = out.File
	file_infra_cros_cmd_labpack_steps_steps_proto_rawDesc = nil
	file_infra_cros_cmd_labpack_steps_steps_proto_goTypes = nil
	file_infra_cros_cmd_labpack_steps_steps_proto_depIdxs = nil
}
