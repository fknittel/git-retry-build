// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: infra/unifiedfleet/api/v1/proto/location.proto

package ufspb

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Lab refers to the different Labs under chrome org
// More labs to be added later if needed
// Next tag: 12
type Lab int32

const (
	Lab_LAB_UNSPECIFIED         Lab = 0
	Lab_LAB_CHROME_ATLANTA      Lab = 1
	Lab_LAB_CHROMEOS_SANTIAM    Lab = 2
	Lab_LAB_CHROMEOS_DESTINY    Lab = 3
	Lab_LAB_CHROMEOS_PROMETHEUS Lab = 4
	Lab_LAB_CHROMEOS_ATLANTIS   Lab = 5
	Lab_LAB_CHROMEOS_LINDAVISTA Lab = 6
	Lab_LAB_DATACENTER_ATL97    Lab = 7
	Lab_LAB_DATACENTER_IAD97    Lab = 8
	Lab_LAB_DATACENTER_MTV96    Lab = 9
	Lab_LAB_DATACENTER_MTV97    Lab = 10
	Lab_LAB_DATACENTER_FUCHSIA  Lab = 11
)

// Enum value maps for Lab.
var (
	Lab_name = map[int32]string{
		0:  "LAB_UNSPECIFIED",
		1:  "LAB_CHROME_ATLANTA",
		2:  "LAB_CHROMEOS_SANTIAM",
		3:  "LAB_CHROMEOS_DESTINY",
		4:  "LAB_CHROMEOS_PROMETHEUS",
		5:  "LAB_CHROMEOS_ATLANTIS",
		6:  "LAB_CHROMEOS_LINDAVISTA",
		7:  "LAB_DATACENTER_ATL97",
		8:  "LAB_DATACENTER_IAD97",
		9:  "LAB_DATACENTER_MTV96",
		10: "LAB_DATACENTER_MTV97",
		11: "LAB_DATACENTER_FUCHSIA",
	}
	Lab_value = map[string]int32{
		"LAB_UNSPECIFIED":         0,
		"LAB_CHROME_ATLANTA":      1,
		"LAB_CHROMEOS_SANTIAM":    2,
		"LAB_CHROMEOS_DESTINY":    3,
		"LAB_CHROMEOS_PROMETHEUS": 4,
		"LAB_CHROMEOS_ATLANTIS":   5,
		"LAB_CHROMEOS_LINDAVISTA": 6,
		"LAB_DATACENTER_ATL97":    7,
		"LAB_DATACENTER_IAD97":    8,
		"LAB_DATACENTER_MTV96":    9,
		"LAB_DATACENTER_MTV97":    10,
		"LAB_DATACENTER_FUCHSIA":  11,
	}
)

func (x Lab) Enum() *Lab {
	p := new(Lab)
	*p = x
	return p
}

func (x Lab) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Lab) Descriptor() protoreflect.EnumDescriptor {
	return file_infra_unifiedfleet_api_v1_proto_location_proto_enumTypes[0].Descriptor()
}

func (Lab) Type() protoreflect.EnumType {
	return &file_infra_unifiedfleet_api_v1_proto_location_proto_enumTypes[0]
}

func (x Lab) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Lab.Descriptor instead.
func (Lab) EnumDescriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_proto_location_proto_rawDescGZIP(), []int{0}
}

// Zone refers to the different network zones under chrome org
type Zone int32

const (
	Zone_ZONE_UNSPECIFIED Zone = 0
	Zone_ZONE_ATLANTA     Zone = 1  // "atl" Building ? Room ?
	Zone_ZONE_CHROMEOS1   Zone = 2  // "chromeos1" // US-MTV-1600 // Santiam
	Zone_ZONE_CHROMEOS2   Zone = 3  // "chromeos2" // US-MTV-2081 // Atlantis
	Zone_ZONE_CHROMEOS3   Zone = 4  // "chromeos3" // US-MTV-946 // Lindavista
	Zone_ZONE_CHROMEOS4   Zone = 5  // "chromeos4" // US-MTV-2081 // Destiny
	Zone_ZONE_CHROMEOS5   Zone = 6  // "chromeos5" // US-MTV-946 // Lindavista
	Zone_ZONE_CHROMEOS6   Zone = 7  // "chromeos6" // US-MTV-2081 // Prometheus
	Zone_ZONE_CHROMEOS7   Zone = 8  // "chromeos7" // US-MTV-946 // Lindavista
	Zone_ZONE_CHROMEOS15  Zone = 10 // "chromeos15" // US-MTV-946 // Lindavista
	Zone_ZONE_ATL97       Zone = 11 // "atl97" //  US-ATL-MET1 // Room ?
	Zone_ZONE_IAD97       Zone = 12 // "iad97" // Building ? Room ?
	Zone_ZONE_MTV96       Zone = 13 // "mtv96" // US-MTV-41 // 1-1M0
	Zone_ZONE_MTV97       Zone = 14 // "mtv97" // US-MTV-1950 // 1-144
	Zone_ZONE_FUCHSIA     Zone = 15 // "lab01" // Building ? Room ?
)

// Enum value maps for Zone.
var (
	Zone_name = map[int32]string{
		0:  "ZONE_UNSPECIFIED",
		1:  "ZONE_ATLANTA",
		2:  "ZONE_CHROMEOS1",
		3:  "ZONE_CHROMEOS2",
		4:  "ZONE_CHROMEOS3",
		5:  "ZONE_CHROMEOS4",
		6:  "ZONE_CHROMEOS5",
		7:  "ZONE_CHROMEOS6",
		8:  "ZONE_CHROMEOS7",
		10: "ZONE_CHROMEOS15",
		11: "ZONE_ATL97",
		12: "ZONE_IAD97",
		13: "ZONE_MTV96",
		14: "ZONE_MTV97",
		15: "ZONE_FUCHSIA",
	}
	Zone_value = map[string]int32{
		"ZONE_UNSPECIFIED": 0,
		"ZONE_ATLANTA":     1,
		"ZONE_CHROMEOS1":   2,
		"ZONE_CHROMEOS2":   3,
		"ZONE_CHROMEOS3":   4,
		"ZONE_CHROMEOS4":   5,
		"ZONE_CHROMEOS5":   6,
		"ZONE_CHROMEOS6":   7,
		"ZONE_CHROMEOS7":   8,
		"ZONE_CHROMEOS15":  10,
		"ZONE_ATL97":       11,
		"ZONE_IAD97":       12,
		"ZONE_MTV96":       13,
		"ZONE_MTV97":       14,
		"ZONE_FUCHSIA":     15,
	}
)

func (x Zone) Enum() *Zone {
	p := new(Zone)
	*p = x
	return p
}

func (x Zone) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Zone) Descriptor() protoreflect.EnumDescriptor {
	return file_infra_unifiedfleet_api_v1_proto_location_proto_enumTypes[1].Descriptor()
}

func (Zone) Type() protoreflect.EnumType {
	return &file_infra_unifiedfleet_api_v1_proto_location_proto_enumTypes[1]
}

func (x Zone) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Zone.Descriptor instead.
func (Zone) EnumDescriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_proto_location_proto_rawDescGZIP(), []int{1}
}

// Location of the asset(Rack/Machine) in the lab
// For Browser machine, lab and rack are the only field to fill in.
// The fine-grained location is mainly for OS machine as we care about rack, row, shelf.
type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Different labs in the chrome org. Required.
	//
	// Deprecated: Do not use.
	Lab Lab `protobuf:"varint,1,opt,name=lab,proto3,enum=unifiedfleet.api.v1.proto.Lab" json:"lab,omitempty"`
	// Each lab has many aisles.
	// This field refers to the aisle number/name in the lab.
	Aisle string `protobuf:"bytes,2,opt,name=aisle,proto3" json:"aisle,omitempty"`
	// Each aisle has many rows.
	// This field refers to the row number/name in the aisle.
	Row string `protobuf:"bytes,3,opt,name=row,proto3" json:"row,omitempty"`
	// Each row has many racks.
	// This field refers to the rack number/name in the row.
	Rack string `protobuf:"bytes,4,opt,name=rack,proto3" json:"rack,omitempty"`
	// The position of the rack in the row.
	RackNumber string `protobuf:"bytes,5,opt,name=rack_number,json=rackNumber,proto3" json:"rack_number,omitempty"`
	// Each rack has many shelves.
	// This field refers to the shelf number/name in the rack.
	Shelf string `protobuf:"bytes,6,opt,name=shelf,proto3" json:"shelf,omitempty"`
	// Each shelf has many positions where assets can be placed.
	// This field refers to the position number/name in the shelf
	Position string `protobuf:"bytes,7,opt,name=position,proto3" json:"position,omitempty"`
	// A string descriptor representing location. This can be to
	// store barcode values for location or user defined names.
	BarcodeName string `protobuf:"bytes,8,opt,name=barcode_name,json=barcodeName,proto3" json:"barcode_name,omitempty"`
	// Different zones in the chrome org. Required.
	Zone Zone `protobuf:"varint,9,opt,name=zone,proto3,enum=unifiedfleet.api.v1.proto.Zone" json:"zone,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_unifiedfleet_api_v1_proto_location_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_infra_unifiedfleet_api_v1_proto_location_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_proto_location_proto_rawDescGZIP(), []int{0}
}

// Deprecated: Do not use.
func (x *Location) GetLab() Lab {
	if x != nil {
		return x.Lab
	}
	return Lab_LAB_UNSPECIFIED
}

func (x *Location) GetAisle() string {
	if x != nil {
		return x.Aisle
	}
	return ""
}

func (x *Location) GetRow() string {
	if x != nil {
		return x.Row
	}
	return ""
}

func (x *Location) GetRack() string {
	if x != nil {
		return x.Rack
	}
	return ""
}

func (x *Location) GetRackNumber() string {
	if x != nil {
		return x.RackNumber
	}
	return ""
}

func (x *Location) GetShelf() string {
	if x != nil {
		return x.Shelf
	}
	return ""
}

func (x *Location) GetPosition() string {
	if x != nil {
		return x.Position
	}
	return ""
}

func (x *Location) GetBarcodeName() string {
	if x != nil {
		return x.BarcodeName
	}
	return ""
}

func (x *Location) GetZone() Zone {
	if x != nil {
		return x.Zone
	}
	return Zone_ZONE_UNSPECIFIED
}

var File_infra_unifiedfleet_api_v1_proto_location_proto protoreflect.FileDescriptor

var file_infra_unifiedfleet_api_v1_proto_location_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66,
	0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x19, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x39, 0x67, 0x6f, 0x2e,
	0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75, 0x6d, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x6c, 0x75, 0x63,
	0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd3, 0x02, 0x0a, 0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x34, 0x0a, 0x03, 0x6c, 0x61, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1e, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x61, 0x62,
	0x42, 0x02, 0x18, 0x01, 0x52, 0x03, 0x6c, 0x61, 0x62, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x69, 0x73,
	0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x69, 0x73, 0x6c, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x72, 0x6f, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x6f,
	0x77, 0x12, 0x3e, 0x0a, 0x04, 0x72, 0x61, 0x63, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x2a, 0xfa, 0x41, 0x27, 0x0a, 0x25, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x2d, 0x66, 0x6c,
	0x65, 0x65, 0x74, 0x2d, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x61, 0x70, 0x70, 0x73, 0x70,
	0x6f, 0x74, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x52, 0x61, 0x63, 0x6b, 0x52, 0x04, 0x72, 0x61, 0x63,
	0x6b, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x61, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x61, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x68, 0x65, 0x6c, 0x66, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x73, 0x68, 0x65, 0x6c, 0x66, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x61, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x62, 0x61, 0x72, 0x63,
	0x6f, 0x64, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x33, 0x0a, 0x04, 0x7a, 0x6f, 0x6e, 0x65, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66,
	0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x5a, 0x6f, 0x6e, 0x65, 0x52, 0x04, 0x7a, 0x6f, 0x6e, 0x65, 0x2a, 0xbf, 0x02, 0x0a,
	0x03, 0x4c, 0x61, 0x62, 0x12, 0x13, 0x0a, 0x0f, 0x4c, 0x41, 0x42, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x4c, 0x41, 0x42,
	0x5f, 0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45, 0x5f, 0x41, 0x54, 0x4c, 0x41, 0x4e, 0x54, 0x41, 0x10,
	0x01, 0x12, 0x18, 0x0a, 0x14, 0x4c, 0x41, 0x42, 0x5f, 0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45, 0x4f,
	0x53, 0x5f, 0x53, 0x41, 0x4e, 0x54, 0x49, 0x41, 0x4d, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x4c,
	0x41, 0x42, 0x5f, 0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45, 0x4f, 0x53, 0x5f, 0x44, 0x45, 0x53, 0x54,
	0x49, 0x4e, 0x59, 0x10, 0x03, 0x12, 0x1b, 0x0a, 0x17, 0x4c, 0x41, 0x42, 0x5f, 0x43, 0x48, 0x52,
	0x4f, 0x4d, 0x45, 0x4f, 0x53, 0x5f, 0x50, 0x52, 0x4f, 0x4d, 0x45, 0x54, 0x48, 0x45, 0x55, 0x53,
	0x10, 0x04, 0x12, 0x19, 0x0a, 0x15, 0x4c, 0x41, 0x42, 0x5f, 0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45,
	0x4f, 0x53, 0x5f, 0x41, 0x54, 0x4c, 0x41, 0x4e, 0x54, 0x49, 0x53, 0x10, 0x05, 0x12, 0x1b, 0x0a,
	0x17, 0x4c, 0x41, 0x42, 0x5f, 0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45, 0x4f, 0x53, 0x5f, 0x4c, 0x49,
	0x4e, 0x44, 0x41, 0x56, 0x49, 0x53, 0x54, 0x41, 0x10, 0x06, 0x12, 0x18, 0x0a, 0x14, 0x4c, 0x41,
	0x42, 0x5f, 0x44, 0x41, 0x54, 0x41, 0x43, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x41, 0x54, 0x4c,
	0x39, 0x37, 0x10, 0x07, 0x12, 0x18, 0x0a, 0x14, 0x4c, 0x41, 0x42, 0x5f, 0x44, 0x41, 0x54, 0x41,
	0x43, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x49, 0x41, 0x44, 0x39, 0x37, 0x10, 0x08, 0x12, 0x18,
	0x0a, 0x14, 0x4c, 0x41, 0x42, 0x5f, 0x44, 0x41, 0x54, 0x41, 0x43, 0x45, 0x4e, 0x54, 0x45, 0x52,
	0x5f, 0x4d, 0x54, 0x56, 0x39, 0x36, 0x10, 0x09, 0x12, 0x18, 0x0a, 0x14, 0x4c, 0x41, 0x42, 0x5f,
	0x44, 0x41, 0x54, 0x41, 0x43, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x4d, 0x54, 0x56, 0x39, 0x37,
	0x10, 0x0a, 0x12, 0x1a, 0x0a, 0x16, 0x4c, 0x41, 0x42, 0x5f, 0x44, 0x41, 0x54, 0x41, 0x43, 0x45,
	0x4e, 0x54, 0x45, 0x52, 0x5f, 0x46, 0x55, 0x43, 0x48, 0x53, 0x49, 0x41, 0x10, 0x0b, 0x2a, 0xa1,
	0x02, 0x0a, 0x04, 0x5a, 0x6f, 0x6e, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x5a, 0x4f, 0x4e, 0x45, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x10, 0x0a,
	0x0c, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x41, 0x54, 0x4c, 0x41, 0x4e, 0x54, 0x41, 0x10, 0x01, 0x12,
	0x12, 0x0a, 0x0e, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45, 0x4f, 0x53,
	0x31, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x43, 0x48, 0x52, 0x4f,
	0x4d, 0x45, 0x4f, 0x53, 0x32, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x5a, 0x4f, 0x4e, 0x45, 0x5f,
	0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45, 0x4f, 0x53, 0x33, 0x10, 0x04, 0x12, 0x12, 0x0a, 0x0e, 0x5a,
	0x4f, 0x4e, 0x45, 0x5f, 0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45, 0x4f, 0x53, 0x34, 0x10, 0x05, 0x12,
	0x12, 0x0a, 0x0e, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45, 0x4f, 0x53,
	0x35, 0x10, 0x06, 0x12, 0x12, 0x0a, 0x0e, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x43, 0x48, 0x52, 0x4f,
	0x4d, 0x45, 0x4f, 0x53, 0x36, 0x10, 0x07, 0x12, 0x12, 0x0a, 0x0e, 0x5a, 0x4f, 0x4e, 0x45, 0x5f,
	0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45, 0x4f, 0x53, 0x37, 0x10, 0x08, 0x12, 0x13, 0x0a, 0x0f, 0x5a,
	0x4f, 0x4e, 0x45, 0x5f, 0x43, 0x48, 0x52, 0x4f, 0x4d, 0x45, 0x4f, 0x53, 0x31, 0x35, 0x10, 0x0a,
	0x12, 0x0e, 0x0a, 0x0a, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x41, 0x54, 0x4c, 0x39, 0x37, 0x10, 0x0b,
	0x12, 0x0e, 0x0a, 0x0a, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x49, 0x41, 0x44, 0x39, 0x37, 0x10, 0x0c,
	0x12, 0x0e, 0x0a, 0x0a, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x4d, 0x54, 0x56, 0x39, 0x36, 0x10, 0x0d,
	0x12, 0x0e, 0x0a, 0x0a, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x4d, 0x54, 0x56, 0x39, 0x37, 0x10, 0x0e,
	0x12, 0x10, 0x0a, 0x0c, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x46, 0x55, 0x43, 0x48, 0x53, 0x49, 0x41,
	0x10, 0x0f, 0x42, 0x27, 0x5a, 0x25, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x75, 0x66, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_infra_unifiedfleet_api_v1_proto_location_proto_rawDescOnce sync.Once
	file_infra_unifiedfleet_api_v1_proto_location_proto_rawDescData = file_infra_unifiedfleet_api_v1_proto_location_proto_rawDesc
)

func file_infra_unifiedfleet_api_v1_proto_location_proto_rawDescGZIP() []byte {
	file_infra_unifiedfleet_api_v1_proto_location_proto_rawDescOnce.Do(func() {
		file_infra_unifiedfleet_api_v1_proto_location_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_unifiedfleet_api_v1_proto_location_proto_rawDescData)
	})
	return file_infra_unifiedfleet_api_v1_proto_location_proto_rawDescData
}

var file_infra_unifiedfleet_api_v1_proto_location_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_infra_unifiedfleet_api_v1_proto_location_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_infra_unifiedfleet_api_v1_proto_location_proto_goTypes = []interface{}{
	(Lab)(0),         // 0: unifiedfleet.api.v1.proto.Lab
	(Zone)(0),        // 1: unifiedfleet.api.v1.proto.Zone
	(*Location)(nil), // 2: unifiedfleet.api.v1.proto.Location
}
var file_infra_unifiedfleet_api_v1_proto_location_proto_depIdxs = []int32{
	0, // 0: unifiedfleet.api.v1.proto.Location.lab:type_name -> unifiedfleet.api.v1.proto.Lab
	1, // 1: unifiedfleet.api.v1.proto.Location.zone:type_name -> unifiedfleet.api.v1.proto.Zone
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_infra_unifiedfleet_api_v1_proto_location_proto_init() }
func file_infra_unifiedfleet_api_v1_proto_location_proto_init() {
	if File_infra_unifiedfleet_api_v1_proto_location_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_infra_unifiedfleet_api_v1_proto_location_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Location); i {
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
			RawDescriptor: file_infra_unifiedfleet_api_v1_proto_location_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_unifiedfleet_api_v1_proto_location_proto_goTypes,
		DependencyIndexes: file_infra_unifiedfleet_api_v1_proto_location_proto_depIdxs,
		EnumInfos:         file_infra_unifiedfleet_api_v1_proto_location_proto_enumTypes,
		MessageInfos:      file_infra_unifiedfleet_api_v1_proto_location_proto_msgTypes,
	}.Build()
	File_infra_unifiedfleet_api_v1_proto_location_proto = out.File
	file_infra_unifiedfleet_api_v1_proto_location_proto_rawDesc = nil
	file_infra_unifiedfleet_api_v1_proto_location_proto_goTypes = nil
	file_infra_unifiedfleet_api_v1_proto_location_proto_depIdxs = nil
}
