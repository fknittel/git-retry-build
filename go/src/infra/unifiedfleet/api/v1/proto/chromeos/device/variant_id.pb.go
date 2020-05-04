// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/unifiedfleet/api/v1/proto/chromeos/device/variant_id.proto

package ufspb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Globally unique identifier.
type VariantId struct {
	// Required. Source: 'mosys platform sku', aka Device-SKU.
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VariantId) Reset()         { *m = VariantId{} }
func (m *VariantId) String() string { return proto.CompactTextString(m) }
func (*VariantId) ProtoMessage()    {}
func (*VariantId) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2516ad08d0135ae, []int{0}
}

func (m *VariantId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VariantId.Unmarshal(m, b)
}
func (m *VariantId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VariantId.Marshal(b, m, deterministic)
}
func (m *VariantId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VariantId.Merge(m, src)
}
func (m *VariantId) XXX_Size() int {
	return xxx_messageInfo_VariantId.Size(m)
}
func (m *VariantId) XXX_DiscardUnknown() {
	xxx_messageInfo_VariantId.DiscardUnknown(m)
}

var xxx_messageInfo_VariantId proto.InternalMessageInfo

func (m *VariantId) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*VariantId)(nil), "unifiedfleet.api.v1.proto.chromeos.device.VariantId")
}

func init() {
	proto.RegisterFile("infra/unifiedfleet/api/v1/proto/chromeos/device/variant_id.proto", fileDescriptor_e2516ad08d0135ae)
}

var fileDescriptor_e2516ad08d0135ae = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x72, 0xc8, 0xcc, 0x4b, 0x2b,
	0x4a, 0xd4, 0x2f, 0xcd, 0xcb, 0x4c, 0xcb, 0x4c, 0x4d, 0x49, 0xcb, 0x49, 0x4d, 0x2d, 0xd1, 0x4f,
	0x2c, 0xc8, 0xd4, 0x2f, 0x33, 0xd4, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x4f, 0xce, 0x28, 0xca,
	0xcf, 0x4d, 0xcd, 0x2f, 0xd6, 0x4f, 0x49, 0x2d, 0xcb, 0x4c, 0x4e, 0xd5, 0x2f, 0x4b, 0x2c, 0xca,
	0x4c, 0xcc, 0x2b, 0x89, 0xcf, 0x4c, 0xd1, 0x03, 0x2b, 0x10, 0xd2, 0x44, 0xd6, 0xab, 0x97, 0x58,
	0x90, 0xa9, 0x57, 0x66, 0x08, 0x91, 0xd2, 0x83, 0xe9, 0xd5, 0x83, 0xe8, 0x55, 0x52, 0xe4, 0xe2,
	0x0c, 0x83, 0x68, 0xf7, 0x4c, 0x11, 0x12, 0xe1, 0x62, 0x2d, 0x4b, 0xcc, 0x29, 0x4d, 0x95, 0x60,
	0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x70, 0x9c, 0xcc, 0xa3, 0x4c, 0x49, 0x74, 0x91, 0x75, 0x69,
	0x5a, 0x71, 0x41, 0x52, 0x12, 0x1b, 0x58, 0xd2, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xaf, 0xb0,
	0x89, 0x64, 0xd1, 0x00, 0x00, 0x00,
}
