// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: infra/tricium/api/admin/v1/driver.proto

package admin

import prpc "go.chromium.org/luci/grpc/prpc"

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// TriggerRequest contains the details for launching a build for a worker.
type TriggerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RunId  int64  `protobuf:"varint,1,opt,name=run_id,json=runId,proto3" json:"run_id,omitempty"`
	Worker string `protobuf:"bytes,3,opt,name=worker,proto3" json:"worker,omitempty"`
}

func (x *TriggerRequest) Reset() {
	*x = TriggerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_tricium_api_admin_v1_driver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TriggerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerRequest) ProtoMessage() {}

func (x *TriggerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_infra_tricium_api_admin_v1_driver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TriggerRequest.ProtoReflect.Descriptor instead.
func (*TriggerRequest) Descriptor() ([]byte, []int) {
	return file_infra_tricium_api_admin_v1_driver_proto_rawDescGZIP(), []int{0}
}

func (x *TriggerRequest) GetRunId() int64 {
	if x != nil {
		return x.RunId
	}
	return 0
}

func (x *TriggerRequest) GetWorker() string {
	if x != nil {
		return x.Worker
	}
	return ""
}

type TriggerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TriggerResponse) Reset() {
	*x = TriggerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_tricium_api_admin_v1_driver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TriggerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerResponse) ProtoMessage() {}

func (x *TriggerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_infra_tricium_api_admin_v1_driver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TriggerResponse.ProtoReflect.Descriptor instead.
func (*TriggerResponse) Descriptor() ([]byte, []int) {
	return file_infra_tricium_api_admin_v1_driver_proto_rawDescGZIP(), []int{1}
}

// CollectRequest contains the details needed to collect results from a worker.
type CollectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RunId int64 `protobuf:"varint,1,opt,name=run_id,json=runId,proto3" json:"run_id,omitempty"`
	// Worker name of the worker to collect results for.
	Worker string `protobuf:"bytes,3,opt,name=worker,proto3" json:"worker,omitempty"`
	// The Buildbucket build ID.
	//
	// Used to collect results from the completed buildbucket worker task.
	BuildId int64 `protobuf:"varint,5,opt,name=build_id,json=buildId,proto3" json:"build_id,omitempty"`
}

func (x *CollectRequest) Reset() {
	*x = CollectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_tricium_api_admin_v1_driver_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CollectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CollectRequest) ProtoMessage() {}

func (x *CollectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_infra_tricium_api_admin_v1_driver_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CollectRequest.ProtoReflect.Descriptor instead.
func (*CollectRequest) Descriptor() ([]byte, []int) {
	return file_infra_tricium_api_admin_v1_driver_proto_rawDescGZIP(), []int{2}
}

func (x *CollectRequest) GetRunId() int64 {
	if x != nil {
		return x.RunId
	}
	return 0
}

func (x *CollectRequest) GetWorker() string {
	if x != nil {
		return x.Worker
	}
	return ""
}

func (x *CollectRequest) GetBuildId() int64 {
	if x != nil {
		return x.BuildId
	}
	return 0
}

type CollectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CollectResponse) Reset() {
	*x = CollectResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_tricium_api_admin_v1_driver_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CollectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CollectResponse) ProtoMessage() {}

func (x *CollectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_infra_tricium_api_admin_v1_driver_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CollectResponse.ProtoReflect.Descriptor instead.
func (*CollectResponse) Descriptor() ([]byte, []int) {
	return file_infra_tricium_api_admin_v1_driver_proto_rawDescGZIP(), []int{3}
}

var File_infra_tricium_api_admin_v1_driver_proto protoreflect.FileDescriptor

var file_infra_tricium_api_admin_v1_driver_proto_rawDesc = []byte{
	0x0a, 0x27, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x74, 0x72, 0x69, 0x63, 0x69, 0x75, 0x6d, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x72, 0x69,
	0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x22, 0x45, 0x0a, 0x0e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x75, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x72, 0x75, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x22, 0x11, 0x0a, 0x0f, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x66, 0x0a, 0x0e, 0x43, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06,
	0x72, 0x75, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x72, 0x75,
	0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x19, 0x0a, 0x08, 0x62,
	0x75, 0x69, 0x6c, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x62,
	0x75, 0x69, 0x6c, 0x64, 0x49, 0x64, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x04,
	0x10, 0x05, 0x22, 0x11, 0x0a, 0x0f, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x7c, 0x0a, 0x06, 0x44, 0x72, 0x69, 0x76, 0x65, 0x72, 0x12,
	0x38, 0x0a, 0x07, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12, 0x15, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x07, 0x43, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x12, 0x15, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x43, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_tricium_api_admin_v1_driver_proto_rawDescOnce sync.Once
	file_infra_tricium_api_admin_v1_driver_proto_rawDescData = file_infra_tricium_api_admin_v1_driver_proto_rawDesc
)

func file_infra_tricium_api_admin_v1_driver_proto_rawDescGZIP() []byte {
	file_infra_tricium_api_admin_v1_driver_proto_rawDescOnce.Do(func() {
		file_infra_tricium_api_admin_v1_driver_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_tricium_api_admin_v1_driver_proto_rawDescData)
	})
	return file_infra_tricium_api_admin_v1_driver_proto_rawDescData
}

var file_infra_tricium_api_admin_v1_driver_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_infra_tricium_api_admin_v1_driver_proto_goTypes = []interface{}{
	(*TriggerRequest)(nil),  // 0: admin.TriggerRequest
	(*TriggerResponse)(nil), // 1: admin.TriggerResponse
	(*CollectRequest)(nil),  // 2: admin.CollectRequest
	(*CollectResponse)(nil), // 3: admin.CollectResponse
}
var file_infra_tricium_api_admin_v1_driver_proto_depIdxs = []int32{
	0, // 0: admin.Driver.Trigger:input_type -> admin.TriggerRequest
	2, // 1: admin.Driver.Collect:input_type -> admin.CollectRequest
	1, // 2: admin.Driver.Trigger:output_type -> admin.TriggerResponse
	3, // 3: admin.Driver.Collect:output_type -> admin.CollectResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_infra_tricium_api_admin_v1_driver_proto_init() }
func file_infra_tricium_api_admin_v1_driver_proto_init() {
	if File_infra_tricium_api_admin_v1_driver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_infra_tricium_api_admin_v1_driver_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TriggerRequest); i {
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
		file_infra_tricium_api_admin_v1_driver_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TriggerResponse); i {
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
		file_infra_tricium_api_admin_v1_driver_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CollectRequest); i {
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
		file_infra_tricium_api_admin_v1_driver_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CollectResponse); i {
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
			RawDescriptor: file_infra_tricium_api_admin_v1_driver_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_infra_tricium_api_admin_v1_driver_proto_goTypes,
		DependencyIndexes: file_infra_tricium_api_admin_v1_driver_proto_depIdxs,
		MessageInfos:      file_infra_tricium_api_admin_v1_driver_proto_msgTypes,
	}.Build()
	File_infra_tricium_api_admin_v1_driver_proto = out.File
	file_infra_tricium_api_admin_v1_driver_proto_rawDesc = nil
	file_infra_tricium_api_admin_v1_driver_proto_goTypes = nil
	file_infra_tricium_api_admin_v1_driver_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DriverClient is the client API for Driver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DriverClient interface {
	// Trigger triggers a build for a Tricium worker.
	Trigger(ctx context.Context, in *TriggerRequest, opts ...grpc.CallOption) (*TriggerResponse, error)
	// Collect collects results from a build running a Tricium worker.
	Collect(ctx context.Context, in *CollectRequest, opts ...grpc.CallOption) (*CollectResponse, error)
}
type driverPRPCClient struct {
	client *prpc.Client
}

func NewDriverPRPCClient(client *prpc.Client) DriverClient {
	return &driverPRPCClient{client}
}

func (c *driverPRPCClient) Trigger(ctx context.Context, in *TriggerRequest, opts ...grpc.CallOption) (*TriggerResponse, error) {
	out := new(TriggerResponse)
	err := c.client.Call(ctx, "admin.Driver", "Trigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverPRPCClient) Collect(ctx context.Context, in *CollectRequest, opts ...grpc.CallOption) (*CollectResponse, error) {
	out := new(CollectResponse)
	err := c.client.Call(ctx, "admin.Driver", "Collect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type driverClient struct {
	cc grpc.ClientConnInterface
}

func NewDriverClient(cc grpc.ClientConnInterface) DriverClient {
	return &driverClient{cc}
}

func (c *driverClient) Trigger(ctx context.Context, in *TriggerRequest, opts ...grpc.CallOption) (*TriggerResponse, error) {
	out := new(TriggerResponse)
	err := c.cc.Invoke(ctx, "/admin.Driver/Trigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) Collect(ctx context.Context, in *CollectRequest, opts ...grpc.CallOption) (*CollectResponse, error) {
	out := new(CollectResponse)
	err := c.cc.Invoke(ctx, "/admin.Driver/Collect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DriverServer is the server API for Driver service.
type DriverServer interface {
	// Trigger triggers a build for a Tricium worker.
	Trigger(context.Context, *TriggerRequest) (*TriggerResponse, error)
	// Collect collects results from a build running a Tricium worker.
	Collect(context.Context, *CollectRequest) (*CollectResponse, error)
}

// UnimplementedDriverServer can be embedded to have forward compatible implementations.
type UnimplementedDriverServer struct {
}

func (*UnimplementedDriverServer) Trigger(context.Context, *TriggerRequest) (*TriggerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Trigger not implemented")
}
func (*UnimplementedDriverServer) Collect(context.Context, *CollectRequest) (*CollectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Collect not implemented")
}

func RegisterDriverServer(s prpc.Registrar, srv DriverServer) {
	s.RegisterService(&_Driver_serviceDesc, srv)
}

func _Driver_Trigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Trigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin.Driver/Trigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Trigger(ctx, req.(*TriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_Collect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Collect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin.Driver/Collect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Collect(ctx, req.(*CollectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Driver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "admin.Driver",
	HandlerType: (*DriverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Trigger",
			Handler:    _Driver_Trigger_Handler,
		},
		{
			MethodName: "Collect",
			Handler:    _Driver_Collect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infra/tricium/api/admin/v1/driver.proto",
}
