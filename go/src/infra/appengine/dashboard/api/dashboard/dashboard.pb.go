// Code generated by protoc-gen-go.
// source: infra/appengine/dashboard/api/dashboard/dashboard.proto
// DO NOT EDIT!

/*
Package dashboard is a generated protocol buffer package.

It is generated from these files:
	infra/appengine/dashboard/api/dashboard/dashboard.proto

It has these top-level messages:
	UpdateOpenIncidentsRequest
	UpdateOpenIncidentsResponse
	Incident
	ChopsService
*/
package dashboard

import prpc "github.com/luci/luci-go/grpc/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Severity int32

const (
	Severity_RED    Severity = 0
	Severity_YELLOW Severity = 1
)

var Severity_name = map[int32]string{
	0: "RED",
	1: "YELLOW",
}
var Severity_value = map[string]int32{
	"RED":    0,
	"YELLOW": 1,
}

func (x Severity) String() string {
	return proto.EnumName(Severity_name, int32(x))
}
func (Severity) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type UpdateOpenIncidentsRequest struct {
	ChopsService *ChopsService `protobuf:"bytes,1,opt,name=chops_service,json=chopsService" json:"chops_service,omitempty"`
}

func (m *UpdateOpenIncidentsRequest) Reset()                    { *m = UpdateOpenIncidentsRequest{} }
func (m *UpdateOpenIncidentsRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateOpenIncidentsRequest) ProtoMessage()               {}
func (*UpdateOpenIncidentsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *UpdateOpenIncidentsRequest) GetChopsService() *ChopsService {
	if m != nil {
		return m.ChopsService
	}
	return nil
}

type UpdateOpenIncidentsResponse struct {
	OpenIncidents []*Incident `protobuf:"bytes,1,rep,name=open_incidents,json=openIncidents" json:"open_incidents,omitempty"`
}

func (m *UpdateOpenIncidentsResponse) Reset()                    { *m = UpdateOpenIncidentsResponse{} }
func (m *UpdateOpenIncidentsResponse) String() string            { return proto.CompactTextString(m) }
func (*UpdateOpenIncidentsResponse) ProtoMessage()               {}
func (*UpdateOpenIncidentsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UpdateOpenIncidentsResponse) GetOpenIncidents() []*Incident {
	if m != nil {
		return m.OpenIncidents
	}
	return nil
}

type Incident struct {
	Id        string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Open      bool     `protobuf:"varint,2,opt,name=open" json:"open,omitempty"`
	StartTime int64    `protobuf:"varint,3,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	EndTime   int64    `protobuf:"varint,4,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
	Severity  Severity `protobuf:"varint,5,opt,name=severity,enum=dashboard.Severity" json:"severity,omitempty"`
}

func (m *Incident) Reset()                    { *m = Incident{} }
func (m *Incident) String() string            { return proto.CompactTextString(m) }
func (*Incident) ProtoMessage()               {}
func (*Incident) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Incident) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Incident) GetOpen() bool {
	if m != nil {
		return m.Open
	}
	return false
}

func (m *Incident) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *Incident) GetEndTime() int64 {
	if m != nil {
		return m.EndTime
	}
	return 0
}

func (m *Incident) GetSeverity() Severity {
	if m != nil {
		return m.Severity
	}
	return Severity_RED
}

type ChopsService struct {
	Name      string      `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Incidents []*Incident `protobuf:"bytes,2,rep,name=incidents" json:"incidents,omitempty"`
	Sla       string      `protobuf:"bytes,3,opt,name=sla" json:"sla,omitempty"`
}

func (m *ChopsService) Reset()                    { *m = ChopsService{} }
func (m *ChopsService) String() string            { return proto.CompactTextString(m) }
func (*ChopsService) ProtoMessage()               {}
func (*ChopsService) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ChopsService) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ChopsService) GetIncidents() []*Incident {
	if m != nil {
		return m.Incidents
	}
	return nil
}

func (m *ChopsService) GetSla() string {
	if m != nil {
		return m.Sla
	}
	return ""
}

func init() {
	proto.RegisterType((*UpdateOpenIncidentsRequest)(nil), "dashboard.UpdateOpenIncidentsRequest")
	proto.RegisterType((*UpdateOpenIncidentsResponse)(nil), "dashboard.UpdateOpenIncidentsResponse")
	proto.RegisterType((*Incident)(nil), "dashboard.Incident")
	proto.RegisterType((*ChopsService)(nil), "dashboard.ChopsService")
	proto.RegisterEnum("dashboard.Severity", Severity_name, Severity_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ChopsServiceStatus service

type ChopsServiceStatusClient interface {
	UpdateOpenIncidents(ctx context.Context, in *UpdateOpenIncidentsRequest, opts ...grpc.CallOption) (*UpdateOpenIncidentsResponse, error)
}
type chopsServiceStatusPRPCClient struct {
	client *prpc.Client
}

func NewChopsServiceStatusPRPCClient(client *prpc.Client) ChopsServiceStatusClient {
	return &chopsServiceStatusPRPCClient{client}
}

func (c *chopsServiceStatusPRPCClient) UpdateOpenIncidents(ctx context.Context, in *UpdateOpenIncidentsRequest, opts ...grpc.CallOption) (*UpdateOpenIncidentsResponse, error) {
	out := new(UpdateOpenIncidentsResponse)
	err := c.client.Call(ctx, "dashboard.ChopsServiceStatus", "UpdateOpenIncidents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type chopsServiceStatusClient struct {
	cc *grpc.ClientConn
}

func NewChopsServiceStatusClient(cc *grpc.ClientConn) ChopsServiceStatusClient {
	return &chopsServiceStatusClient{cc}
}

func (c *chopsServiceStatusClient) UpdateOpenIncidents(ctx context.Context, in *UpdateOpenIncidentsRequest, opts ...grpc.CallOption) (*UpdateOpenIncidentsResponse, error) {
	out := new(UpdateOpenIncidentsResponse)
	err := grpc.Invoke(ctx, "/dashboard.ChopsServiceStatus/UpdateOpenIncidents", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ChopsServiceStatus service

type ChopsServiceStatusServer interface {
	UpdateOpenIncidents(context.Context, *UpdateOpenIncidentsRequest) (*UpdateOpenIncidentsResponse, error)
}

func RegisterChopsServiceStatusServer(s prpc.Registrar, srv ChopsServiceStatusServer) {
	s.RegisterService(&_ChopsServiceStatus_serviceDesc, srv)
}

func _ChopsServiceStatus_UpdateOpenIncidents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateOpenIncidentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChopsServiceStatusServer).UpdateOpenIncidents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.ChopsServiceStatus/UpdateOpenIncidents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChopsServiceStatusServer).UpdateOpenIncidents(ctx, req.(*UpdateOpenIncidentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ChopsServiceStatus_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dashboard.ChopsServiceStatus",
	HandlerType: (*ChopsServiceStatusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateOpenIncidents",
			Handler:    _ChopsServiceStatus_UpdateOpenIncidents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infra/appengine/dashboard/api/dashboard/dashboard.proto",
}

func init() {
	proto.RegisterFile("infra/appengine/dashboard/api/dashboard/dashboard.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 371 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xcf, 0xcb, 0xd3, 0x40,
	0x10, 0xfd, 0x36, 0xf9, 0xfc, 0x9a, 0x4c, 0x7f, 0x50, 0xb6, 0x07, 0x63, 0x45, 0x0c, 0x01, 0x25,
	0x78, 0x68, 0xb0, 0x1e, 0x04, 0xf1, 0xa6, 0x3d, 0x08, 0x85, 0xc2, 0x56, 0x91, 0x7a, 0x29, 0xdb,
	0xec, 0xb4, 0x5d, 0xb0, 0x9b, 0x35, 0xbb, 0x2d, 0x08, 0xfe, 0x23, 0xfe, 0xb7, 0x92, 0x6d, 0xd3,
	0x46, 0x68, 0xf9, 0x6e, 0x6f, 0xf6, 0xbd, 0xbc, 0x79, 0x93, 0x19, 0x78, 0x2f, 0xd5, 0xba, 0xe4,
	0x19, 0xd7, 0x1a, 0xd5, 0x46, 0x2a, 0xcc, 0x04, 0x37, 0xdb, 0x55, 0xc1, 0x4b, 0x91, 0x71, 0x2d,
	0x1b, 0xd5, 0x19, 0x8d, 0x74, 0x59, 0xd8, 0x82, 0x86, 0xe7, 0x87, 0xe4, 0x07, 0x0c, 0xbf, 0x69,
	0xc1, 0x2d, 0xce, 0x34, 0xaa, 0x2f, 0x2a, 0x97, 0x02, 0x95, 0x35, 0x0c, 0x7f, 0xed, 0xd1, 0x58,
	0xfa, 0x11, 0xba, 0xf9, 0xb6, 0xd0, 0x66, 0x69, 0xb0, 0x3c, 0xc8, 0x1c, 0x23, 0x12, 0x93, 0xb4,
	0x3d, 0x7e, 0x3a, 0xba, 0x38, 0x7e, 0xaa, 0xf8, 0xf9, 0x91, 0x66, 0x9d, 0xbc, 0x51, 0x25, 0x0b,
	0x78, 0x7e, 0xd5, 0xdb, 0xe8, 0x42, 0x19, 0xa4, 0x1f, 0xa0, 0x57, 0x68, 0x54, 0x4b, 0x59, 0x33,
	0x11, 0x89, 0xfd, 0xb4, 0x3d, 0x1e, 0x34, 0xdc, 0xeb, 0xaf, 0x58, 0xb7, 0x68, 0x7a, 0x24, 0x7f,
	0x09, 0x04, 0x75, 0x45, 0x7b, 0xe0, 0x49, 0xe1, 0xa2, 0x85, 0xcc, 0x93, 0x82, 0x52, 0xb8, 0xaf,
	0xd4, 0x91, 0x17, 0x93, 0x34, 0x60, 0x0e, 0xd3, 0x17, 0x00, 0xc6, 0xf2, 0xd2, 0x2e, 0xad, 0xdc,
	0x61, 0xe4, 0xc7, 0x24, 0xf5, 0x59, 0xe8, 0x5e, 0xbe, 0xca, 0x1d, 0xd2, 0x67, 0x10, 0xa0, 0x12,
	0x47, 0xf2, 0xde, 0x91, 0x2d, 0x54, 0xc2, 0x51, 0x19, 0x04, 0x06, 0x0f, 0x58, 0x4a, 0xfb, 0x3b,
	0x7a, 0x12, 0x93, 0xb4, 0xf7, 0x5f, 0xc0, 0xf9, 0x89, 0x62, 0x67, 0x51, 0xb2, 0x81, 0x4e, 0xf3,
	0xa7, 0x54, 0x71, 0x14, 0xdf, 0xe1, 0x29, 0xa0, 0xc3, 0xf4, 0x2d, 0x84, 0x97, 0xb1, 0xbd, 0xdb,
	0x63, 0x5f, 0x54, 0xb4, 0x0f, 0xbe, 0xf9, 0xc9, 0x5d, 0xf4, 0x90, 0x55, 0xf0, 0xcd, 0x4b, 0x08,
	0xea, 0xf6, 0xb4, 0x05, 0x3e, 0x9b, 0x7c, 0xee, 0xdf, 0x51, 0x80, 0x87, 0xc5, 0x64, 0x3a, 0x9d,
	0x7d, 0xef, 0x93, 0xf1, 0x1f, 0xa0, 0xcd, 0x24, 0x73, 0xcb, 0xed, 0xde, 0xd0, 0x35, 0x0c, 0xae,
	0xac, 0x85, 0xbe, 0x6a, 0xf4, 0xbf, 0x7d, 0x12, 0xc3, 0xd7, 0x8f, 0xc9, 0x8e, 0xdb, 0x4d, 0xee,
	0x56, 0x0f, 0xee, 0xd8, 0xde, 0xfd, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x92, 0x39, 0xe1, 0x16, 0xa7,
	0x02, 0x00, 0x00,
}
