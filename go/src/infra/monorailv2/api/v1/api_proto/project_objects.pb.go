// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/v1/api_proto/project_objects.proto

package monorail_v1

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// The top level organization of issues in Monorail.
//
// See monorail/doc/userguide/concepts.md#Projects-and-roles.
// and monorail/doc/userguide/project-owners.md#why-does-monorail-have-projects
// Next available tag: 2
type Project struct {
	// Resource name of the project.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Project) Reset()         { *m = Project{} }
func (m *Project) String() string { return proto.CompactTextString(m) }
func (*Project) ProtoMessage()    {}
func (*Project) Descriptor() ([]byte, []int) {
	return fileDescriptor_63b58a6c62ab6503, []int{0}
}

func (m *Project) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Project.Unmarshal(m, b)
}
func (m *Project) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Project.Marshal(b, m, deterministic)
}
func (m *Project) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Project.Merge(m, src)
}
func (m *Project) XXX_Size() int {
	return xxx_messageInfo_Project.Size(m)
}
func (m *Project) XXX_DiscardUnknown() {
	xxx_messageInfo_Project.DiscardUnknown(m)
}

var xxx_messageInfo_Project proto.InternalMessageInfo

func (m *Project) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Potential steps along the development process that an issue can be in.
//
// See monorail/doc/userguide/project-owners.md#How-to-configure-statuses
// (-- aip.dev/not-precedent: "Status" should be reserved for HTTP/gRPC codes
//     per aip.dev/216. Monorail's Status  preceded the AIP standards, and is
//     used extensively throughout the system.)
// Next available tag: 2
type StatusDef struct {
	// Resource name of the status.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusDef) Reset()         { *m = StatusDef{} }
func (m *StatusDef) String() string { return proto.CompactTextString(m) }
func (*StatusDef) ProtoMessage()    {}
func (*StatusDef) Descriptor() ([]byte, []int) {
	return fileDescriptor_63b58a6c62ab6503, []int{1}
}

func (m *StatusDef) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusDef.Unmarshal(m, b)
}
func (m *StatusDef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusDef.Marshal(b, m, deterministic)
}
func (m *StatusDef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusDef.Merge(m, src)
}
func (m *StatusDef) XXX_Size() int {
	return xxx_messageInfo_StatusDef.Size(m)
}
func (m *StatusDef) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusDef.DiscardUnknown(m)
}

var xxx_messageInfo_StatusDef proto.InternalMessageInfo

func (m *StatusDef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Well-known labels that can be applied to issues within the project.
//
// See monorail/doc/userguide/concepts.md#issue-fields-and-labels.
// Next available tag: 2
type LabelDef struct {
	// Resource name of the label.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LabelDef) Reset()         { *m = LabelDef{} }
func (m *LabelDef) String() string { return proto.CompactTextString(m) }
func (*LabelDef) ProtoMessage()    {}
func (*LabelDef) Descriptor() ([]byte, []int) {
	return fileDescriptor_63b58a6c62ab6503, []int{2}
}

func (m *LabelDef) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LabelDef.Unmarshal(m, b)
}
func (m *LabelDef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LabelDef.Marshal(b, m, deterministic)
}
func (m *LabelDef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LabelDef.Merge(m, src)
}
func (m *LabelDef) XXX_Size() int {
	return xxx_messageInfo_LabelDef.Size(m)
}
func (m *LabelDef) XXX_DiscardUnknown() {
	xxx_messageInfo_LabelDef.DiscardUnknown(m)
}

var xxx_messageInfo_LabelDef proto.InternalMessageInfo

func (m *LabelDef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Custom fields defined for the project.
//
// See monorail/doc/userguide/concepts.md#issue-fields-and-labels.
// Next available tag: 2
type FieldDef struct {
	// Resource name of the field.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FieldDef) Reset()         { *m = FieldDef{} }
func (m *FieldDef) String() string { return proto.CompactTextString(m) }
func (*FieldDef) ProtoMessage()    {}
func (*FieldDef) Descriptor() ([]byte, []int) {
	return fileDescriptor_63b58a6c62ab6503, []int{3}
}

func (m *FieldDef) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FieldDef.Unmarshal(m, b)
}
func (m *FieldDef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FieldDef.Marshal(b, m, deterministic)
}
func (m *FieldDef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FieldDef.Merge(m, src)
}
func (m *FieldDef) XXX_Size() int {
	return xxx_messageInfo_FieldDef.Size(m)
}
func (m *FieldDef) XXX_DiscardUnknown() {
	xxx_messageInfo_FieldDef.DiscardUnknown(m)
}

var xxx_messageInfo_FieldDef proto.InternalMessageInfo

func (m *FieldDef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// A high level definition of the part of the software affected by an issue.
//
// See monorail/doc/userguide/project-owners.md#how-to-configure-components.
// Next available tag: 2
type ComponentDef struct {
	// Resource name of the component.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ComponentDef) Reset()         { *m = ComponentDef{} }
func (m *ComponentDef) String() string { return proto.CompactTextString(m) }
func (*ComponentDef) ProtoMessage()    {}
func (*ComponentDef) Descriptor() ([]byte, []int) {
	return fileDescriptor_63b58a6c62ab6503, []int{4}
}

func (m *ComponentDef) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComponentDef.Unmarshal(m, b)
}
func (m *ComponentDef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComponentDef.Marshal(b, m, deterministic)
}
func (m *ComponentDef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComponentDef.Merge(m, src)
}
func (m *ComponentDef) XXX_Size() int {
	return xxx_messageInfo_ComponentDef.Size(m)
}
func (m *ComponentDef) XXX_DiscardUnknown() {
	xxx_messageInfo_ComponentDef.DiscardUnknown(m)
}

var xxx_messageInfo_ComponentDef proto.InternalMessageInfo

func (m *ComponentDef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Defines approvals that issues within the project may need.
//
// TODO(monorail:7193): Add documentation for approvals.
// Next available tag: 2
type ApprovalDef struct {
	// Resource name of the approval.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ApprovalDef) Reset()         { *m = ApprovalDef{} }
func (m *ApprovalDef) String() string { return proto.CompactTextString(m) }
func (*ApprovalDef) ProtoMessage()    {}
func (*ApprovalDef) Descriptor() ([]byte, []int) {
	return fileDescriptor_63b58a6c62ab6503, []int{5}
}

func (m *ApprovalDef) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApprovalDef.Unmarshal(m, b)
}
func (m *ApprovalDef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApprovalDef.Marshal(b, m, deterministic)
}
func (m *ApprovalDef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApprovalDef.Merge(m, src)
}
func (m *ApprovalDef) XXX_Size() int {
	return xxx_messageInfo_ApprovalDef.Size(m)
}
func (m *ApprovalDef) XXX_DiscardUnknown() {
	xxx_messageInfo_ApprovalDef.DiscardUnknown(m)
}

var xxx_messageInfo_ApprovalDef proto.InternalMessageInfo

func (m *ApprovalDef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Project)(nil), "monorail.v1.Project")
	proto.RegisterType((*StatusDef)(nil), "monorail.v1.StatusDef")
	proto.RegisterType((*LabelDef)(nil), "monorail.v1.LabelDef")
	proto.RegisterType((*FieldDef)(nil), "monorail.v1.FieldDef")
	proto.RegisterType((*ComponentDef)(nil), "monorail.v1.ComponentDef")
	proto.RegisterType((*ApprovalDef)(nil), "monorail.v1.ApprovalDef")
}

func init() {
	proto.RegisterFile("api/v1/api_proto/project_objects.proto", fileDescriptor_63b58a6c62ab6503)
}

var fileDescriptor_63b58a6c62ab6503 = []byte{
	// 327 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0x29, 0x88, 0xda, 0xa9, 0xa7, 0x82, 0xff, 0x7a, 0x92, 0x1c, 0xaa, 0xa2, 0x66, 0x28,
	0xde, 0xbc, 0x2d, 0x16, 0x05, 0xb1, 0x50, 0xf4, 0x03, 0x94, 0x4d, 0xba, 0x6d, 0x23, 0x49, 0x66,
	0xd9, 0xdd, 0xe6, 0xb2, 0xe4, 0x0b, 0xe7, 0x53, 0x48, 0x93, 0x34, 0x4d, 0x30, 0x39, 0x78, 0xca,
	0x9b, 0xcc, 0x7b, 0xef, 0x17, 0x86, 0xc0, 0x98, 0xcb, 0x00, 0x93, 0x09, 0x72, 0x19, 0x2c, 0xa4,
	0x22, 0x43, 0x28, 0x15, 0xfd, 0x08, 0xdf, 0x2c, 0xc8, 0xdb, 0x3d, 0xb4, 0x9b, 0xbf, 0x1d, 0x0e,
	0x22, 0x8a, 0x49, 0xf1, 0x20, 0x74, 0x93, 0xc9, 0xe8, 0x71, 0x4d, 0xb4, 0x0e, 0x45, 0x19, 0x28,
	0x86, 0x5d, 0x03, 0xae, 0x02, 0x11, 0x2e, 0x17, 0x9e, 0xd8, 0xf0, 0x24, 0x20, 0x55, 0x44, 0x47,
	0xe3, 0x2e, 0xb7, 0x12, 0x9a, 0xb6, 0xca, 0x17, 0x85, 0xcf, 0x99, 0xc1, 0xc9, 0xbc, 0x60, 0x0f,
	0x87, 0x70, 0x14, 0xf3, 0x48, 0x5c, 0xf5, 0x6e, 0x7a, 0x77, 0xfd, 0xaf, 0x5c, 0xbf, 0xb8, 0x19,
	0x7b, 0x80, 0x73, 0x2e, 0x03, 0xd7, 0x57, 0xde, 0x76, 0xed, 0xfa, 0x14, 0x61, 0xe5, 0x2f, 0x3f,
	0x5a, 0xa3, 0x2d, 0x55, 0xea, 0x6c, 0xa0, 0xff, 0x6d, 0xb8, 0xd9, 0xea, 0xa9, 0x58, 0xb5, 0x16,
	0xbe, 0x67, 0x6c, 0x0a, 0x97, 0xcd, 0xc2, 0x43, 0xe2, 0xfe, 0x6f, 0x25, 0xea, 0xfd, 0x56, 0xa3,
	0xad, 0x74, 0xea, 0x2c, 0xe1, 0xf4, 0x93, 0x7b, 0x22, 0xec, 0x02, 0x4d, 0x33, 0xc6, 0xe0, 0xa2,
	0x09, 0xaa, 0x02, 0xb7, 0x2d, 0x9c, 0xb0, 0x5c, 0x6a, 0xb4, 0x7b, 0x99, 0x53, 0xde, 0x76, 0xe7,
	0xfd, 0x0f, 0xa5, 0x0a, 0xb4, 0x51, 0x56, 0xe5, 0x52, 0xa3, 0xdd, 0xcb, 0xd4, 0x31, 0x70, 0xf6,
	0x4a, 0x91, 0xa4, 0x58, 0xc4, 0xa6, 0x8b, 0x34, 0xcf, 0xd8, 0x0c, 0x46, 0x4d, 0x52, 0x23, 0x84,
	0x2d, 0x34, 0xbf, 0x66, 0xd0, 0x68, 0xeb, 0x63, 0xea, 0x48, 0x18, 0x30, 0x29, 0x15, 0x25, 0xbc,
	0xf3, 0x88, 0xb3, 0x8c, 0x7d, 0xc0, 0x75, 0x13, 0x5a, 0xcf, 0x3c, 0xb5, 0x30, 0xf9, 0x61, 0xaf,
	0xd1, 0xd6, 0xa6, 0xd4, 0x3b, 0xce, 0xff, 0xb9, 0xe7, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x36,
	0x91, 0x3a, 0x8b, 0x00, 0x03, 0x00, 0x00,
}
