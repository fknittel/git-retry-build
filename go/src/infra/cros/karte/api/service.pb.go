// Copyright 2021 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.1
// source: infra/cros/karte/api/service.proto

package kartepb

import prpc "go.chromium.org/luci/grpc/prpc"

import (
	context "context"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// CreateActionRequest creates a single action.
type CreateActionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action *Action `protobuf:"bytes,1,opt,name=action,proto3" json:"action,omitempty"`
}

func (x *CreateActionRequest) Reset() {
	*x = CreateActionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_cros_karte_api_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateActionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateActionRequest) ProtoMessage() {}

func (x *CreateActionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_infra_cros_karte_api_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateActionRequest.ProtoReflect.Descriptor instead.
func (*CreateActionRequest) Descriptor() ([]byte, []int) {
	return file_infra_cros_karte_api_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreateActionRequest) GetAction() *Action {
	if x != nil {
		return x.Action
	}
	return nil
}

// CreateObservationRequest creates a single action.
type CreateObservationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// An observation is the observation record being created.
	Observation *Observation `protobuf:"bytes,1,opt,name=observation,proto3" json:"observation,omitempty"`
}

func (x *CreateObservationRequest) Reset() {
	*x = CreateObservationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_cros_karte_api_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateObservationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateObservationRequest) ProtoMessage() {}

func (x *CreateObservationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_infra_cros_karte_api_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateObservationRequest.ProtoReflect.Descriptor instead.
func (*CreateObservationRequest) Descriptor() ([]byte, []int) {
	return file_infra_cros_karte_api_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateObservationRequest) GetObservation() *Observation {
	if x != nil {
		return x.Observation
	}
	return nil
}

// ListActionsRequest takes a page size and a token indicating where to start.
type ListActionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The maximum number of actions to return. The service may return fewer than
	// this value.
	// If unspecified, at most 50 actions will be returned.
	// The maximum value is 1000; values above 1000 will be coerced to 1000.
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// A page token, received from a previous `ListActions` call.
	// Provide this to retrieve the subsequent page.
	//
	// When paginating, all other parameters provided to `ListActions` must match
	// the call that provided the page token.
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	// Filter is a query using an expression syntax described in
	// filter_syntax.md.
	//
	// Currently supported filterable fields for actions are:
	// - kind
	Filter string `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *ListActionsRequest) Reset() {
	*x = ListActionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_cros_karte_api_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListActionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListActionsRequest) ProtoMessage() {}

func (x *ListActionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_infra_cros_karte_api_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListActionsRequest.ProtoReflect.Descriptor instead.
func (*ListActionsRequest) Descriptor() ([]byte, []int) {
	return file_infra_cros_karte_api_service_proto_rawDescGZIP(), []int{2}
}

func (x *ListActionsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListActionsRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListActionsRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

// ListActionsResponse returns the actions in question and returns a page token
// indicating where to start looking in the next search.
// The page token will be empty if and only if we have reached the end of the
// results.
type ListActionsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// These are all the actions fitting the criteria specified. Currently, no
	// criteria can be provided, so every action matches.
	Actions []*Action `protobuf:"bytes,1,rep,name=actions,proto3" json:"actions,omitempty"`
	// This is the page token that is needed for pagination. This token
	// must be supplied verbatim to subsequent calls to ListActions.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListActionsResponse) Reset() {
	*x = ListActionsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_cros_karte_api_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListActionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListActionsResponse) ProtoMessage() {}

func (x *ListActionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_infra_cros_karte_api_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListActionsResponse.ProtoReflect.Descriptor instead.
func (*ListActionsResponse) Descriptor() ([]byte, []int) {
	return file_infra_cros_karte_api_service_proto_rawDescGZIP(), []int{3}
}

func (x *ListActionsResponse) GetActions() []*Action {
	if x != nil {
		return x.Actions
	}
	return nil
}

func (x *ListActionsResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

// ListObservationsRequest take a page size and a token indicating where to
// start.
type ListObservationsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The maximum number of observations to return. The service may return fewer
	// than this value. If unspecified, at most 50 observations will be returned.
	// The maximum value is 1000; values above 1000 will be coerced to 1000.
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// A page token, received from a previous `ListObservations` call.
	// Provide this to retrieve the subsequent page.
	//
	// When paginating, all other parameters provided to `ListObservations` must
	// match the call that provided the page token.
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	// Filter is a query using an expression syntax described in
	// filter_syntax.md.
	//
	// Currently supported filterable values for actions are:
	// - metric_kind
	Filter string `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *ListObservationsRequest) Reset() {
	*x = ListObservationsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_cros_karte_api_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListObservationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListObservationsRequest) ProtoMessage() {}

func (x *ListObservationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_infra_cros_karte_api_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListObservationsRequest.ProtoReflect.Descriptor instead.
func (*ListObservationsRequest) Descriptor() ([]byte, []int) {
	return file_infra_cros_karte_api_service_proto_rawDescGZIP(), []int{4}
}

func (x *ListObservationsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListObservationsRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListObservationsRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

// ListObservationsResponse returns the observations in quetoin and returns a
// page token indicating where to start looking in the next search. The page
// token will be empty if and only if we have reached the end of the results.
type ListObservationsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// These are all the actions fitting the criteria specified. Currently, no
	// criteria can be provided, so every action matches.
	Observations []*Observation `protobuf:"bytes,1,rep,name=observations,proto3" json:"observations,omitempty"`
	// This is the page token that is needed for pagination. This token
	// must be supplied verbatim to subsequent calls to ListActions.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListObservationsResponse) Reset() {
	*x = ListObservationsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_cros_karte_api_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListObservationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListObservationsResponse) ProtoMessage() {}

func (x *ListObservationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_infra_cros_karte_api_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListObservationsResponse.ProtoReflect.Descriptor instead.
func (*ListObservationsResponse) Descriptor() ([]byte, []int) {
	return file_infra_cros_karte_api_service_proto_rawDescGZIP(), []int{5}
}

func (x *ListObservationsResponse) GetObservations() []*Observation {
	if x != nil {
		return x.Observations
	}
	return nil
}

func (x *ListObservationsResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

var File_infra_cros_karte_api_service_proto protoreflect.FileDescriptor

var file_infra_cros_karte_api_service_proto_rawDesc = []byte{
	0x0a, 0x22, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x2f, 0x6b, 0x61, 0x72,
	0x74, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6b,
	0x61, 0x72, 0x74, 0x65, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x2f,
	0x6b, 0x61, 0x72, 0x74, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x26, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x63, 0x72,
	0x6f, 0x73, 0x2f, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6f, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4a,
	0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73,
	0x2e, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x03, 0xe0,
	0x41, 0x02, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x5e, 0x0a, 0x18, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x42, 0x0a, 0x0b, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x68,
	0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x2e, 0x4f, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x0b, 0x6f,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x68, 0x0a, 0x12, 0x4c, 0x69,
	0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x22, 0x6f, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x07, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63,
	0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x2e, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x26, 0x0a,
	0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x6d, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x22, 0x83, 0x01, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3f, 0x0a, 0x0c, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65,
	0x6f, 0x73, 0x2e, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x2e, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78,
	0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0xde, 0x03, 0x0a, 0x05, 0x4b,
	0x61, 0x72, 0x74, 0x65, 0x12, 0x68, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e,
	0x6b, 0x61, 0x72, 0x74, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x63, 0x68, 0x72, 0x6f,
	0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x22, 0x0b, 0x2f, 0x76, 0x31, 0x2f, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x3a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x81,
	0x01, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e,
	0x6b, 0x61, 0x72, 0x74, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x62, 0x73, 0x65,
	0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b,
	0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x2e,
	0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x25, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x1f, 0x22, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x3a, 0x0b, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x6b, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x22, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6b, 0x61, 0x72,
	0x74, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73,
	0x2e, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x13, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x0d, 0x12, 0x0b, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x7a, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x27, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6b,
	0x61, 0x72, 0x74, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x63,
	0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x13, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x12, 0x0b,
	0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x1e, 0x5a, 0x1c, 0x69,
	0x6e, 0x66, 0x72, 0x61, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x2f, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x3b, 0x6b, 0x61, 0x72, 0x74, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_infra_cros_karte_api_service_proto_rawDescOnce sync.Once
	file_infra_cros_karte_api_service_proto_rawDescData = file_infra_cros_karte_api_service_proto_rawDesc
)

func file_infra_cros_karte_api_service_proto_rawDescGZIP() []byte {
	file_infra_cros_karte_api_service_proto_rawDescOnce.Do(func() {
		file_infra_cros_karte_api_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_cros_karte_api_service_proto_rawDescData)
	})
	return file_infra_cros_karte_api_service_proto_rawDescData
}

var file_infra_cros_karte_api_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_infra_cros_karte_api_service_proto_goTypes = []interface{}{
	(*CreateActionRequest)(nil),      // 0: chromeos.karte.CreateActionRequest
	(*CreateObservationRequest)(nil), // 1: chromeos.karte.CreateObservationRequest
	(*ListActionsRequest)(nil),       // 2: chromeos.karte.ListActionsRequest
	(*ListActionsResponse)(nil),      // 3: chromeos.karte.ListActionsResponse
	(*ListObservationsRequest)(nil),  // 4: chromeos.karte.ListObservationsRequest
	(*ListObservationsResponse)(nil), // 5: chromeos.karte.ListObservationsResponse
	(*Action)(nil),                   // 6: chromeos.karte.Action
	(*Observation)(nil),              // 7: chromeos.karte.Observation
}
var file_infra_cros_karte_api_service_proto_depIdxs = []int32{
	6, // 0: chromeos.karte.CreateActionRequest.action:type_name -> chromeos.karte.Action
	7, // 1: chromeos.karte.CreateObservationRequest.observation:type_name -> chromeos.karte.Observation
	6, // 2: chromeos.karte.ListActionsResponse.actions:type_name -> chromeos.karte.Action
	7, // 3: chromeos.karte.ListObservationsResponse.observations:type_name -> chromeos.karte.Observation
	0, // 4: chromeos.karte.Karte.CreateAction:input_type -> chromeos.karte.CreateActionRequest
	1, // 5: chromeos.karte.Karte.CreateObservation:input_type -> chromeos.karte.CreateObservationRequest
	2, // 6: chromeos.karte.Karte.ListActions:input_type -> chromeos.karte.ListActionsRequest
	4, // 7: chromeos.karte.Karte.ListObservations:input_type -> chromeos.karte.ListObservationsRequest
	6, // 8: chromeos.karte.Karte.CreateAction:output_type -> chromeos.karte.Action
	7, // 9: chromeos.karte.Karte.CreateObservation:output_type -> chromeos.karte.Observation
	3, // 10: chromeos.karte.Karte.ListActions:output_type -> chromeos.karte.ListActionsResponse
	5, // 11: chromeos.karte.Karte.ListObservations:output_type -> chromeos.karte.ListObservationsResponse
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_infra_cros_karte_api_service_proto_init() }
func file_infra_cros_karte_api_service_proto_init() {
	if File_infra_cros_karte_api_service_proto != nil {
		return
	}
	file_infra_cros_karte_api_action_proto_init()
	file_infra_cros_karte_api_observation_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_infra_cros_karte_api_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateActionRequest); i {
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
		file_infra_cros_karte_api_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateObservationRequest); i {
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
		file_infra_cros_karte_api_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListActionsRequest); i {
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
		file_infra_cros_karte_api_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListActionsResponse); i {
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
		file_infra_cros_karte_api_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListObservationsRequest); i {
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
		file_infra_cros_karte_api_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListObservationsResponse); i {
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
			RawDescriptor: file_infra_cros_karte_api_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_infra_cros_karte_api_service_proto_goTypes,
		DependencyIndexes: file_infra_cros_karte_api_service_proto_depIdxs,
		MessageInfos:      file_infra_cros_karte_api_service_proto_msgTypes,
	}.Build()
	File_infra_cros_karte_api_service_proto = out.File
	file_infra_cros_karte_api_service_proto_rawDesc = nil
	file_infra_cros_karte_api_service_proto_goTypes = nil
	file_infra_cros_karte_api_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// KarteClient is the client API for Karte service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type KarteClient interface {
	// CreateAction creates an action and returns the action just created.
	CreateAction(ctx context.Context, in *CreateActionRequest, opts ...grpc.CallOption) (*Action, error)
	// CreateObservation creates an observation and returns the observation
	// that was just created.
	// This API is based on https://google.aip.dev/133.
	CreateObservation(ctx context.Context, in *CreateObservationRequest, opts ...grpc.CallOption) (*Observation, error)
	// ListActions lists all the actions that Karte knows about.
	// The order in which the actions are returned is undefined.
	ListActions(ctx context.Context, in *ListActionsRequest, opts ...grpc.CallOption) (*ListActionsResponse, error)
	// ListObservations lists all the observations that Karte knows about.
	// The order in which the observations are returned is undefined.
	ListObservations(ctx context.Context, in *ListObservationsRequest, opts ...grpc.CallOption) (*ListObservationsResponse, error)
}
type kartePRPCClient struct {
	client *prpc.Client
}

func NewKartePRPCClient(client *prpc.Client) KarteClient {
	return &kartePRPCClient{client}
}

func (c *kartePRPCClient) CreateAction(ctx context.Context, in *CreateActionRequest, opts ...grpc.CallOption) (*Action, error) {
	out := new(Action)
	err := c.client.Call(ctx, "chromeos.karte.Karte", "CreateAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kartePRPCClient) CreateObservation(ctx context.Context, in *CreateObservationRequest, opts ...grpc.CallOption) (*Observation, error) {
	out := new(Observation)
	err := c.client.Call(ctx, "chromeos.karte.Karte", "CreateObservation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kartePRPCClient) ListActions(ctx context.Context, in *ListActionsRequest, opts ...grpc.CallOption) (*ListActionsResponse, error) {
	out := new(ListActionsResponse)
	err := c.client.Call(ctx, "chromeos.karte.Karte", "ListActions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kartePRPCClient) ListObservations(ctx context.Context, in *ListObservationsRequest, opts ...grpc.CallOption) (*ListObservationsResponse, error) {
	out := new(ListObservationsResponse)
	err := c.client.Call(ctx, "chromeos.karte.Karte", "ListObservations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type karteClient struct {
	cc grpc.ClientConnInterface
}

func NewKarteClient(cc grpc.ClientConnInterface) KarteClient {
	return &karteClient{cc}
}

func (c *karteClient) CreateAction(ctx context.Context, in *CreateActionRequest, opts ...grpc.CallOption) (*Action, error) {
	out := new(Action)
	err := c.cc.Invoke(ctx, "/chromeos.karte.Karte/CreateAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *karteClient) CreateObservation(ctx context.Context, in *CreateObservationRequest, opts ...grpc.CallOption) (*Observation, error) {
	out := new(Observation)
	err := c.cc.Invoke(ctx, "/chromeos.karte.Karte/CreateObservation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *karteClient) ListActions(ctx context.Context, in *ListActionsRequest, opts ...grpc.CallOption) (*ListActionsResponse, error) {
	out := new(ListActionsResponse)
	err := c.cc.Invoke(ctx, "/chromeos.karte.Karte/ListActions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *karteClient) ListObservations(ctx context.Context, in *ListObservationsRequest, opts ...grpc.CallOption) (*ListObservationsResponse, error) {
	out := new(ListObservationsResponse)
	err := c.cc.Invoke(ctx, "/chromeos.karte.Karte/ListObservations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KarteServer is the server API for Karte service.
type KarteServer interface {
	// CreateAction creates an action and returns the action just created.
	CreateAction(context.Context, *CreateActionRequest) (*Action, error)
	// CreateObservation creates an observation and returns the observation
	// that was just created.
	// This API is based on https://google.aip.dev/133.
	CreateObservation(context.Context, *CreateObservationRequest) (*Observation, error)
	// ListActions lists all the actions that Karte knows about.
	// The order in which the actions are returned is undefined.
	ListActions(context.Context, *ListActionsRequest) (*ListActionsResponse, error)
	// ListObservations lists all the observations that Karte knows about.
	// The order in which the observations are returned is undefined.
	ListObservations(context.Context, *ListObservationsRequest) (*ListObservationsResponse, error)
}

// UnimplementedKarteServer can be embedded to have forward compatible implementations.
type UnimplementedKarteServer struct {
}

func (*UnimplementedKarteServer) CreateAction(context.Context, *CreateActionRequest) (*Action, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAction not implemented")
}
func (*UnimplementedKarteServer) CreateObservation(context.Context, *CreateObservationRequest) (*Observation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateObservation not implemented")
}
func (*UnimplementedKarteServer) ListActions(context.Context, *ListActionsRequest) (*ListActionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListActions not implemented")
}
func (*UnimplementedKarteServer) ListObservations(context.Context, *ListObservationsRequest) (*ListObservationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListObservations not implemented")
}

func RegisterKarteServer(s prpc.Registrar, srv KarteServer) {
	s.RegisterService(&_Karte_serviceDesc, srv)
}

func _Karte_CreateAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KarteServer).CreateAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chromeos.karte.Karte/CreateAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KarteServer).CreateAction(ctx, req.(*CreateActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Karte_CreateObservation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateObservationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KarteServer).CreateObservation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chromeos.karte.Karte/CreateObservation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KarteServer).CreateObservation(ctx, req.(*CreateObservationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Karte_ListActions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListActionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KarteServer).ListActions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chromeos.karte.Karte/ListActions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KarteServer).ListActions(ctx, req.(*ListActionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Karte_ListObservations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListObservationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KarteServer).ListObservations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chromeos.karte.Karte/ListObservations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KarteServer).ListObservations(ctx, req.(*ListObservationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Karte_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chromeos.karte.Karte",
	HandlerType: (*KarteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAction",
			Handler:    _Karte_CreateAction_Handler,
		},
		{
			MethodName: "CreateObservation",
			Handler:    _Karte_CreateObservation_Handler,
		},
		{
			MethodName: "ListActions",
			Handler:    _Karte_ListActions_Handler,
		},
		{
			MethodName: "ListObservations",
			Handler:    _Karte_ListObservations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infra/cros/karte/api/service.proto",
}
