// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/api_proto/features_objects.proto

package monorail

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

// TODO(jojwang): add editors and followers
// Next available tag: 5
type Hotlist struct {
	OwnerRef             *UserRef `protobuf:"bytes,1,opt,name=owner_ref,json=ownerRef,proto3" json:"owner_ref,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Summary              string   `protobuf:"bytes,3,opt,name=summary,proto3" json:"summary,omitempty"`
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Hotlist) Reset()         { *m = Hotlist{} }
func (m *Hotlist) String() string { return proto.CompactTextString(m) }
func (*Hotlist) ProtoMessage()    {}
func (*Hotlist) Descriptor() ([]byte, []int) {
	return fileDescriptor_806b6b78af767289, []int{0}
}

func (m *Hotlist) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hotlist.Unmarshal(m, b)
}
func (m *Hotlist) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hotlist.Marshal(b, m, deterministic)
}
func (m *Hotlist) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hotlist.Merge(m, src)
}
func (m *Hotlist) XXX_Size() int {
	return xxx_messageInfo_Hotlist.Size(m)
}
func (m *Hotlist) XXX_DiscardUnknown() {
	xxx_messageInfo_Hotlist.DiscardUnknown(m)
}

var xxx_messageInfo_Hotlist proto.InternalMessageInfo

func (m *Hotlist) GetOwnerRef() *UserRef {
	if m != nil {
		return m.OwnerRef
	}
	return nil
}

func (m *Hotlist) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Hotlist) GetSummary() string {
	if m != nil {
		return m.Summary
	}
	return ""
}

func (m *Hotlist) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

// Next available tag: 6
type HotlistItem struct {
	Issue                *Issue   `protobuf:"bytes,1,opt,name=issue,proto3" json:"issue,omitempty"`
	Rank                 uint32   `protobuf:"varint,2,opt,name=rank,proto3" json:"rank,omitempty"`
	AdderRef             *UserRef `protobuf:"bytes,3,opt,name=adder_ref,json=adderRef,proto3" json:"adder_ref,omitempty"`
	AddedTimestamp       uint32   `protobuf:"varint,4,opt,name=added_timestamp,json=addedTimestamp,proto3" json:"added_timestamp,omitempty"`
	Note                 string   `protobuf:"bytes,5,opt,name=note,proto3" json:"note,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HotlistItem) Reset()         { *m = HotlistItem{} }
func (m *HotlistItem) String() string { return proto.CompactTextString(m) }
func (*HotlistItem) ProtoMessage()    {}
func (*HotlistItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_806b6b78af767289, []int{1}
}

func (m *HotlistItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HotlistItem.Unmarshal(m, b)
}
func (m *HotlistItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HotlistItem.Marshal(b, m, deterministic)
}
func (m *HotlistItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HotlistItem.Merge(m, src)
}
func (m *HotlistItem) XXX_Size() int {
	return xxx_messageInfo_HotlistItem.Size(m)
}
func (m *HotlistItem) XXX_DiscardUnknown() {
	xxx_messageInfo_HotlistItem.DiscardUnknown(m)
}

var xxx_messageInfo_HotlistItem proto.InternalMessageInfo

func (m *HotlistItem) GetIssue() *Issue {
	if m != nil {
		return m.Issue
	}
	return nil
}

func (m *HotlistItem) GetRank() uint32 {
	if m != nil {
		return m.Rank
	}
	return 0
}

func (m *HotlistItem) GetAdderRef() *UserRef {
	if m != nil {
		return m.AdderRef
	}
	return nil
}

func (m *HotlistItem) GetAddedTimestamp() uint32 {
	if m != nil {
		return m.AddedTimestamp
	}
	return 0
}

func (m *HotlistItem) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

// Next available tag: 5
type HotlistPeopleDelta struct {
	NewOwnerRef          *UserRef   `protobuf:"bytes,1,opt,name=new_owner_ref,json=newOwnerRef,proto3" json:"new_owner_ref,omitempty"`
	AddEditorRefs        []*UserRef `protobuf:"bytes,2,rep,name=add_editor_refs,json=addEditorRefs,proto3" json:"add_editor_refs,omitempty"`
	AddFollowerRefs      []*UserRef `protobuf:"bytes,3,rep,name=add_follower_refs,json=addFollowerRefs,proto3" json:"add_follower_refs,omitempty"`
	RemoveUserRefs       []*UserRef `protobuf:"bytes,4,rep,name=remove_user_refs,json=removeUserRefs,proto3" json:"remove_user_refs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *HotlistPeopleDelta) Reset()         { *m = HotlistPeopleDelta{} }
func (m *HotlistPeopleDelta) String() string { return proto.CompactTextString(m) }
func (*HotlistPeopleDelta) ProtoMessage()    {}
func (*HotlistPeopleDelta) Descriptor() ([]byte, []int) {
	return fileDescriptor_806b6b78af767289, []int{2}
}

func (m *HotlistPeopleDelta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HotlistPeopleDelta.Unmarshal(m, b)
}
func (m *HotlistPeopleDelta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HotlistPeopleDelta.Marshal(b, m, deterministic)
}
func (m *HotlistPeopleDelta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HotlistPeopleDelta.Merge(m, src)
}
func (m *HotlistPeopleDelta) XXX_Size() int {
	return xxx_messageInfo_HotlistPeopleDelta.Size(m)
}
func (m *HotlistPeopleDelta) XXX_DiscardUnknown() {
	xxx_messageInfo_HotlistPeopleDelta.DiscardUnknown(m)
}

var xxx_messageInfo_HotlistPeopleDelta proto.InternalMessageInfo

func (m *HotlistPeopleDelta) GetNewOwnerRef() *UserRef {
	if m != nil {
		return m.NewOwnerRef
	}
	return nil
}

func (m *HotlistPeopleDelta) GetAddEditorRefs() []*UserRef {
	if m != nil {
		return m.AddEditorRefs
	}
	return nil
}

func (m *HotlistPeopleDelta) GetAddFollowerRefs() []*UserRef {
	if m != nil {
		return m.AddFollowerRefs
	}
	return nil
}

func (m *HotlistPeopleDelta) GetRemoveUserRefs() []*UserRef {
	if m != nil {
		return m.RemoveUserRefs
	}
	return nil
}

func init() {
	proto.RegisterType((*Hotlist)(nil), "monorail.Hotlist")
	proto.RegisterType((*HotlistItem)(nil), "monorail.HotlistItem")
	proto.RegisterType((*HotlistPeopleDelta)(nil), "monorail.HotlistPeopleDelta")
}

func init() {
	proto.RegisterFile("api/api_proto/features_objects.proto", fileDescriptor_806b6b78af767289)
}

var fileDescriptor_806b6b78af767289 = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x41, 0xab, 0xd3, 0x40,
	0x14, 0x85, 0x49, 0xd3, 0xda, 0x76, 0x42, 0x5a, 0x3b, 0xab, 0xd0, 0x55, 0x2c, 0x8a, 0x5d, 0xa5,
	0xa0, 0xb8, 0x10, 0x71, 0xa7, 0x62, 0x57, 0xca, 0xa0, 0xeb, 0x30, 0x6d, 0x6e, 0x60, 0x34, 0x33,
	0x37, 0xcc, 0x4c, 0x0c, 0x6e, 0xdd, 0xf9, 0x83, 0xfc, 0x7f, 0x92, 0xdb, 0x09, 0xef, 0x15, 0x5e,
	0xe1, 0xed, 0xee, 0x9c, 0x7b, 0x4e, 0xe6, 0x3b, 0x43, 0xd8, 0x73, 0xd9, 0xaa, 0x83, 0x6c, 0x55,
	0xd9, 0x5a, 0xf4, 0x78, 0xa8, 0x41, 0xfa, 0xce, 0x82, 0x2b, 0xf1, 0xf4, 0x03, 0xce, 0xde, 0x15,
	0x24, 0xf3, 0x85, 0x46, 0x83, 0x56, 0xaa, 0x66, 0xbb, 0xbd, 0xf6, 0x9f, 0x51, 0x6b, 0x34, 0x17,
	0xd7, 0xf6, 0xd9, 0xf5, 0x4e, 0x39, 0xd7, 0xc1, 0xf5, 0x87, 0x76, 0x7f, 0x23, 0x36, 0xff, 0x8c,
	0xbe, 0x51, 0xce, 0xf3, 0x82, 0x2d, 0xb1, 0x37, 0x60, 0x4b, 0x0b, 0x75, 0x16, 0xe5, 0xd1, 0x3e,
	0x79, 0xb5, 0x29, 0xc6, 0x8b, 0x8a, 0xef, 0x0e, 0xac, 0x80, 0x5a, 0x2c, 0xc8, 0x23, 0xa0, 0xe6,
	0x9c, 0x4d, 0x8d, 0xd4, 0x90, 0x4d, 0xf2, 0x68, 0xbf, 0x14, 0x34, 0xf3, 0x8c, 0xcd, 0x5d, 0xa7,
	0xb5, 0xb4, 0xbf, 0xb3, 0x98, 0xe4, 0xf1, 0xc8, 0x73, 0x96, 0x54, 0xe0, 0xce, 0x56, 0xb5, 0x5e,
	0xa1, 0xc9, 0xa6, 0xb4, 0xbd, 0x2f, 0xed, 0xfe, 0x45, 0x2c, 0x09, 0x2c, 0x47, 0x0f, 0x9a, 0xbf,
	0x60, 0x33, 0x42, 0x0e, 0x2c, 0xeb, 0x3b, 0x96, 0xe3, 0x20, 0x8b, 0xcb, 0x76, 0xc0, 0xb0, 0xd2,
	0xfc, 0x24, 0x8c, 0x54, 0xd0, 0x3c, 0x54, 0x91, 0x55, 0x15, 0xaa, 0xc4, 0x37, 0xab, 0x90, 0x67,
	0xa8, 0xf2, 0x92, 0xad, 0x87, 0xb9, 0x2a, 0xbd, 0xd2, 0xe0, 0xbc, 0xd4, 0x2d, 0x01, 0xa6, 0x62,
	0x45, 0xf2, 0xb7, 0x51, 0xa5, 0xce, 0xe8, 0x21, 0x9b, 0x85, 0xce, 0xe8, 0x61, 0xf7, 0x67, 0xc2,
	0x78, 0xe0, 0xfe, 0x0a, 0xd8, 0x36, 0xf0, 0x01, 0x1a, 0x2f, 0xf9, 0x1b, 0x96, 0x1a, 0xe8, 0xcb,
	0x47, 0x3c, 0x69, 0x62, 0xa0, 0xff, 0x32, 0xbe, 0xea, 0x5b, 0x42, 0x29, 0xa1, 0x52, 0x1e, 0x29,
	0xe7, 0xb2, 0x49, 0x1e, 0x3f, 0x1c, 0x4c, 0x65, 0x55, 0x7d, 0x24, 0xa3, 0x80, 0xda, 0xf1, 0xf7,
	0x6c, 0x33, 0x44, 0x6b, 0x6c, 0x1a, 0xec, 0x21, 0x84, 0xe3, 0x5b, 0xe1, 0xe1, 0x9a, 0x4f, 0xc1,
	0x4a, 0xf1, 0x77, 0xec, 0xa9, 0x05, 0x8d, 0xbf, 0xa0, 0xec, 0xdc, 0x98, 0x9e, 0xde, 0x4a, 0xaf,
	0x2e, 0xd6, 0x70, 0x74, 0xa7, 0x27, 0xf4, 0x3f, 0xbd, 0xfe, 0x1f, 0x00, 0x00, 0xff, 0xff, 0xd5,
	0x1d, 0xf1, 0xe2, 0xc0, 0x02, 0x00, 0x00,
}
