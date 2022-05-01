// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PinpointClient is the client API for Pinpoint service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PinpointClient interface {
	// Schedules a Pinpoint Job for execution.
	ScheduleJob(ctx context.Context, in *ScheduleJobRequest, opts ...grpc.CallOption) (*Job, error)
	// Retrieves details about a Pinpoint Job.
	GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*Job, error)
	// Lists jobs with filters.
	ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error)
	// Cancels an ongoing job.
	CancelJob(ctx context.Context, in *CancelJobRequest, opts ...grpc.CallOption) (*Job, error)
}

type pinpointClient struct {
	cc grpc.ClientConnInterface
}

func NewPinpointClient(cc grpc.ClientConnInterface) PinpointClient {
	return &pinpointClient{cc}
}

func (c *pinpointClient) ScheduleJob(ctx context.Context, in *ScheduleJobRequest, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := c.cc.Invoke(ctx, "/pinpoint.Pinpoint/ScheduleJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pinpointClient) GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := c.cc.Invoke(ctx, "/pinpoint.Pinpoint/GetJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pinpointClient) ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error) {
	out := new(ListJobsResponse)
	err := c.cc.Invoke(ctx, "/pinpoint.Pinpoint/ListJobs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pinpointClient) CancelJob(ctx context.Context, in *CancelJobRequest, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := c.cc.Invoke(ctx, "/pinpoint.Pinpoint/CancelJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PinpointServer is the server API for Pinpoint service.
// All implementations must embed UnimplementedPinpointServer
// for forward compatibility
type PinpointServer interface {
	// Schedules a Pinpoint Job for execution.
	ScheduleJob(context.Context, *ScheduleJobRequest) (*Job, error)
	// Retrieves details about a Pinpoint Job.
	GetJob(context.Context, *GetJobRequest) (*Job, error)
	// Lists jobs with filters.
	ListJobs(context.Context, *ListJobsRequest) (*ListJobsResponse, error)
	// Cancels an ongoing job.
	CancelJob(context.Context, *CancelJobRequest) (*Job, error)
	mustEmbedUnimplementedPinpointServer()
}

// UnimplementedPinpointServer must be embedded to have forward compatible implementations.
type UnimplementedPinpointServer struct {
}

func (UnimplementedPinpointServer) ScheduleJob(context.Context, *ScheduleJobRequest) (*Job, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScheduleJob not implemented")
}
func (UnimplementedPinpointServer) GetJob(context.Context, *GetJobRequest) (*Job, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJob not implemented")
}
func (UnimplementedPinpointServer) ListJobs(context.Context, *ListJobsRequest) (*ListJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListJobs not implemented")
}
func (UnimplementedPinpointServer) CancelJob(context.Context, *CancelJobRequest) (*Job, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelJob not implemented")
}
func (UnimplementedPinpointServer) mustEmbedUnimplementedPinpointServer() {}

// UnsafePinpointServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PinpointServer will
// result in compilation errors.
type UnsafePinpointServer interface {
	mustEmbedUnimplementedPinpointServer()
}

func RegisterPinpointServer(s grpc.ServiceRegistrar, srv PinpointServer) {
	s.RegisterService(&Pinpoint_ServiceDesc, srv)
}

func _Pinpoint_ScheduleJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScheduleJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PinpointServer).ScheduleJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pinpoint.Pinpoint/ScheduleJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PinpointServer).ScheduleJob(ctx, req.(*ScheduleJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pinpoint_GetJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PinpointServer).GetJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pinpoint.Pinpoint/GetJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PinpointServer).GetJob(ctx, req.(*GetJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pinpoint_ListJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListJobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PinpointServer).ListJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pinpoint.Pinpoint/ListJobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PinpointServer).ListJobs(ctx, req.(*ListJobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pinpoint_CancelJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PinpointServer).CancelJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pinpoint.Pinpoint/CancelJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PinpointServer).CancelJob(ctx, req.(*CancelJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Pinpoint_ServiceDesc is the grpc.ServiceDesc for Pinpoint service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pinpoint_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pinpoint.Pinpoint",
	HandlerType: (*PinpointServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ScheduleJob",
			Handler:    _Pinpoint_ScheduleJob_Handler,
		},
		{
			MethodName: "GetJob",
			Handler:    _Pinpoint_GetJob_Handler,
		},
		{
			MethodName: "ListJobs",
			Handler:    _Pinpoint_ListJobs_Handler,
		},
		{
			MethodName: "CancelJob",
			Handler:    _Pinpoint_CancelJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infra/chromeperf/pinpoint/pinpoint.proto",
}
