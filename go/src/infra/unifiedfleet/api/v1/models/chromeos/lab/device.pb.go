// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: infra/unifiedfleet/api/v1/models/chromeos/lab/device.proto

package ufspb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	device "infra/unifiedfleet/api/v1/models/chromeos/device"
	manufacturing "infra/unifiedfleet/api/v1/models/chromeos/manufacturing"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// critical_pools are pool labels that the builders are dependent on, and
// that the cros-infra team is responsible for managing explicitly. All other
// pool labels used for adhoc labeling of DUTs go into self_serve_pools.
// TO BE DELETED
type DeviceUnderTest_DUTPool int32

const (
	DeviceUnderTest_DUT_POOL_INVALID       DeviceUnderTest_DUTPool = 0
	DeviceUnderTest_DUT_POOL_CQ            DeviceUnderTest_DUTPool = 1
	DeviceUnderTest_DUT_POOL_BVT           DeviceUnderTest_DUTPool = 2
	DeviceUnderTest_DUT_POOL_SUITES        DeviceUnderTest_DUTPool = 3
	DeviceUnderTest_DUT_POOL_CTS           DeviceUnderTest_DUTPool = 4
	DeviceUnderTest_DUT_POOL_CTS_PERBUILD  DeviceUnderTest_DUTPool = 5
	DeviceUnderTest_DUT_POOL_CONTINUOUS    DeviceUnderTest_DUTPool = 6
	DeviceUnderTest_DUT_POOL_ARC_PRESUBMIT DeviceUnderTest_DUTPool = 7
	DeviceUnderTest_DUT_POOL_QUOTA         DeviceUnderTest_DUTPool = 8
)

// Enum value maps for DeviceUnderTest_DUTPool.
var (
	DeviceUnderTest_DUTPool_name = map[int32]string{
		0: "DUT_POOL_INVALID",
		1: "DUT_POOL_CQ",
		2: "DUT_POOL_BVT",
		3: "DUT_POOL_SUITES",
		4: "DUT_POOL_CTS",
		5: "DUT_POOL_CTS_PERBUILD",
		6: "DUT_POOL_CONTINUOUS",
		7: "DUT_POOL_ARC_PRESUBMIT",
		8: "DUT_POOL_QUOTA",
	}
	DeviceUnderTest_DUTPool_value = map[string]int32{
		"DUT_POOL_INVALID":       0,
		"DUT_POOL_CQ":            1,
		"DUT_POOL_BVT":           2,
		"DUT_POOL_SUITES":        3,
		"DUT_POOL_CTS":           4,
		"DUT_POOL_CTS_PERBUILD":  5,
		"DUT_POOL_CONTINUOUS":    6,
		"DUT_POOL_ARC_PRESUBMIT": 7,
		"DUT_POOL_QUOTA":         8,
	}
)

func (x DeviceUnderTest_DUTPool) Enum() *DeviceUnderTest_DUTPool {
	p := new(DeviceUnderTest_DUTPool)
	*p = x
	return p
}

func (x DeviceUnderTest_DUTPool) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DeviceUnderTest_DUTPool) Descriptor() protoreflect.EnumDescriptor {
	return file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_enumTypes[0].Descriptor()
}

func (DeviceUnderTest_DUTPool) Type() protoreflect.EnumType {
	return &file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_enumTypes[0]
}

func (x DeviceUnderTest_DUTPool) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DeviceUnderTest_DUTPool.Descriptor instead.
func (DeviceUnderTest_DUTPool) EnumDescriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescGZIP(), []int{1, 0}
}

// Next Tag: 7
type ChromeOSDevice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A unique ID for chromeos device, a randomly generated uuid or AssetTag.
	Id              *ChromeOSDeviceID       `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	SerialNumber    string                  `protobuf:"bytes,2,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	ManufacturingId *manufacturing.ConfigID `protobuf:"bytes,3,opt,name=manufacturing_id,json=manufacturingId,proto3" json:"manufacturing_id,omitempty"`
	// Device config identifiers.
	// These values will be extracted from DUT and joinable to device config.
	DeviceConfigId *device.ConfigId `protobuf:"bytes,4,opt,name=device_config_id,json=deviceConfigId,proto3" json:"device_config_id,omitempty"`
	// Types that are assignable to Device:
	//	*ChromeOSDevice_Dut
	//	*ChromeOSDevice_Labstation
	Device isChromeOSDevice_Device `protobuf_oneof:"device"`
}

func (x *ChromeOSDevice) Reset() {
	*x = ChromeOSDevice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChromeOSDevice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChromeOSDevice) ProtoMessage() {}

func (x *ChromeOSDevice) ProtoReflect() protoreflect.Message {
	mi := &file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChromeOSDevice.ProtoReflect.Descriptor instead.
func (*ChromeOSDevice) Descriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescGZIP(), []int{0}
}

func (x *ChromeOSDevice) GetId() *ChromeOSDeviceID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ChromeOSDevice) GetSerialNumber() string {
	if x != nil {
		return x.SerialNumber
	}
	return ""
}

func (x *ChromeOSDevice) GetManufacturingId() *manufacturing.ConfigID {
	if x != nil {
		return x.ManufacturingId
	}
	return nil
}

func (x *ChromeOSDevice) GetDeviceConfigId() *device.ConfigId {
	if x != nil {
		return x.DeviceConfigId
	}
	return nil
}

func (m *ChromeOSDevice) GetDevice() isChromeOSDevice_Device {
	if m != nil {
		return m.Device
	}
	return nil
}

func (x *ChromeOSDevice) GetDut() *DeviceUnderTest {
	if x, ok := x.GetDevice().(*ChromeOSDevice_Dut); ok {
		return x.Dut
	}
	return nil
}

func (x *ChromeOSDevice) GetLabstation() *Labstation {
	if x, ok := x.GetDevice().(*ChromeOSDevice_Labstation); ok {
		return x.Labstation
	}
	return nil
}

type isChromeOSDevice_Device interface {
	isChromeOSDevice_Device()
}

type ChromeOSDevice_Dut struct {
	Dut *DeviceUnderTest `protobuf:"bytes,5,opt,name=dut,proto3,oneof"`
}

type ChromeOSDevice_Labstation struct {
	Labstation *Labstation `protobuf:"bytes,6,opt,name=labstation,proto3,oneof"`
}

func (*ChromeOSDevice_Dut) isChromeOSDevice_Device() {}

func (*ChromeOSDevice_Labstation) isChromeOSDevice_Device() {}

// Next Tag: 5
type DeviceUnderTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hostname      string                    `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Peripherals   *Peripherals              `protobuf:"bytes,2,opt,name=peripherals,proto3" json:"peripherals,omitempty"`
	CriticalPools []DeviceUnderTest_DUTPool `protobuf:"varint,3,rep,packed,name=critical_pools,json=criticalPools,proto3,enum=unifiedfleet.api.v1.models.chromeos.lab.DeviceUnderTest_DUTPool" json:"critical_pools,omitempty"`
	Pools         []string                  `protobuf:"bytes,4,rep,name=pools,proto3" json:"pools,omitempty"`
	Licenses      []*License                `protobuf:"bytes,5,rep,name=licenses,proto3" json:"licenses,omitempty"`
}

func (x *DeviceUnderTest) Reset() {
	*x = DeviceUnderTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceUnderTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceUnderTest) ProtoMessage() {}

func (x *DeviceUnderTest) ProtoReflect() protoreflect.Message {
	mi := &file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceUnderTest.ProtoReflect.Descriptor instead.
func (*DeviceUnderTest) Descriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescGZIP(), []int{1}
}

func (x *DeviceUnderTest) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *DeviceUnderTest) GetPeripherals() *Peripherals {
	if x != nil {
		return x.Peripherals
	}
	return nil
}

func (x *DeviceUnderTest) GetCriticalPools() []DeviceUnderTest_DUTPool {
	if x != nil {
		return x.CriticalPools
	}
	return nil
}

func (x *DeviceUnderTest) GetPools() []string {
	if x != nil {
		return x.Pools
	}
	return nil
}

func (x *DeviceUnderTest) GetLicenses() []*License {
	if x != nil {
		return x.Licenses
	}
	return nil
}

// Next Tag: 5
type Labstation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hostname string   `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Servos   []*Servo `protobuf:"bytes,2,rep,name=servos,proto3" json:"servos,omitempty"`
	Rpm      *OSRPM   `protobuf:"bytes,3,opt,name=rpm,proto3" json:"rpm,omitempty"`
	Pools    []string `protobuf:"bytes,4,rep,name=pools,proto3" json:"pools,omitempty"`
}

func (x *Labstation) Reset() {
	*x = Labstation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Labstation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Labstation) ProtoMessage() {}

func (x *Labstation) ProtoReflect() protoreflect.Message {
	mi := &file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Labstation.ProtoReflect.Descriptor instead.
func (*Labstation) Descriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescGZIP(), []int{2}
}

func (x *Labstation) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *Labstation) GetServos() []*Servo {
	if x != nil {
		return x.Servos
	}
	return nil
}

func (x *Labstation) GetRpm() *OSRPM {
	if x != nil {
		return x.Rpm
	}
	return nil
}

func (x *Labstation) GetPools() []string {
	if x != nil {
		return x.Pools
	}
	return nil
}

var File_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto protoreflect.FileDescriptor

var file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDesc = []byte{
	0x0a, 0x3a, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66,
	0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x6c, 0x61, 0x62, 0x2f,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x27, 0x75, 0x6e,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f,
	0x73, 0x2e, 0x6c, 0x61, 0x62, 0x1a, 0x40, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69,
	0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73,
	0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x69,
	0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x46, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75,
	0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65,
	0x6f, 0x73, 0x2f, 0x6c, 0x61, 0x62, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x5f,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x3b, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c,
	0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x6c, 0x61, 0x62, 0x2f, 0x6c,
	0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x3f, 0x69, 0x6e,
	0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x63,
	0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x6c, 0x61, 0x62, 0x2f, 0x70, 0x65, 0x72, 0x69,
	0x70, 0x68, 0x65, 0x72, 0x61, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x39, 0x69,
	0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65,
	0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f,
	0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x6c, 0x61, 0x62, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x47, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f,
	0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d,
	0x65, 0x6f, 0x73, 0x2f, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69, 0x6e,
	0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xf7, 0x03, 0x0a, 0x0e, 0x43, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x4f, 0x53, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x39, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68,
	0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6c, 0x61, 0x62, 0x2e, 0x43, 0x68, 0x72, 0x6f, 0x6d,
	0x65, 0x4f, 0x53, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x44, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x23, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x66, 0x0a, 0x10, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74,
	0x75, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x3b,
	0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f,
	0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69,
	0x6e, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x44, 0x52, 0x0f, 0x6d, 0x61, 0x6e,
	0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69, 0x6e, 0x67, 0x49, 0x64, 0x12, 0x5e, 0x0a, 0x10,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64,
	0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x64, 0x52, 0x0e, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x64, 0x12, 0x4c, 0x0a, 0x03,
	0x64, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x38, 0x2e, 0x75, 0x6e, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e,
	0x6c, 0x61, 0x62, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x55, 0x6e, 0x64, 0x65, 0x72, 0x54,
	0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x03, 0x64, 0x75, 0x74, 0x12, 0x55, 0x0a, 0x0a, 0x6c, 0x61,
	0x62, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x33,
	0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f,
	0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6c, 0x61, 0x62, 0x2e, 0x4c, 0x61, 0x62, 0x73, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0a, 0x6c, 0x61, 0x62, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x08, 0x0a, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x22, 0xa2, 0x04, 0x0a, 0x0f,
	0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x55, 0x6e, 0x64, 0x65, 0x72, 0x54, 0x65, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x56, 0x0a, 0x0b, 0x70,
	0x65, 0x72, 0x69, 0x70, 0x68, 0x65, 0x72, 0x61, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x34, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68,
	0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6c, 0x61, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x69, 0x70,
	0x68, 0x65, 0x72, 0x61, 0x6c, 0x73, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x69, 0x70, 0x68, 0x65, 0x72,
	0x61, 0x6c, 0x73, 0x12, 0x67, 0x0a, 0x0e, 0x63, 0x72, 0x69, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x5f,
	0x70, 0x6f, 0x6f, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x40, 0x2e, 0x75, 0x6e,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f,
	0x73, 0x2e, 0x6c, 0x61, 0x62, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x55, 0x6e, 0x64, 0x65,
	0x72, 0x54, 0x65, 0x73, 0x74, 0x2e, 0x44, 0x55, 0x54, 0x50, 0x6f, 0x6f, 0x6c, 0x52, 0x0d, 0x63,
	0x72, 0x69, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x50, 0x6f, 0x6f, 0x6c, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x70, 0x6f, 0x6f, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x70, 0x6f, 0x6f,
	0x6c, 0x73, 0x12, 0x4c, 0x0a, 0x08, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c,
	0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6c, 0x61, 0x62, 0x2e, 0x4c,
	0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x73,
	0x22, 0xcd, 0x01, 0x0a, 0x07, 0x44, 0x55, 0x54, 0x50, 0x6f, 0x6f, 0x6c, 0x12, 0x14, 0x0a, 0x10,
	0x44, 0x55, 0x54, 0x5f, 0x50, 0x4f, 0x4f, 0x4c, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44,
	0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x55, 0x54, 0x5f, 0x50, 0x4f, 0x4f, 0x4c, 0x5f, 0x43,
	0x51, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x44, 0x55, 0x54, 0x5f, 0x50, 0x4f, 0x4f, 0x4c, 0x5f,
	0x42, 0x56, 0x54, 0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x44, 0x55, 0x54, 0x5f, 0x50, 0x4f, 0x4f,
	0x4c, 0x5f, 0x53, 0x55, 0x49, 0x54, 0x45, 0x53, 0x10, 0x03, 0x12, 0x10, 0x0a, 0x0c, 0x44, 0x55,
	0x54, 0x5f, 0x50, 0x4f, 0x4f, 0x4c, 0x5f, 0x43, 0x54, 0x53, 0x10, 0x04, 0x12, 0x19, 0x0a, 0x15,
	0x44, 0x55, 0x54, 0x5f, 0x50, 0x4f, 0x4f, 0x4c, 0x5f, 0x43, 0x54, 0x53, 0x5f, 0x50, 0x45, 0x52,
	0x42, 0x55, 0x49, 0x4c, 0x44, 0x10, 0x05, 0x12, 0x17, 0x0a, 0x13, 0x44, 0x55, 0x54, 0x5f, 0x50,
	0x4f, 0x4f, 0x4c, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x49, 0x4e, 0x55, 0x4f, 0x55, 0x53, 0x10, 0x06,
	0x12, 0x1a, 0x0a, 0x16, 0x44, 0x55, 0x54, 0x5f, 0x50, 0x4f, 0x4f, 0x4c, 0x5f, 0x41, 0x52, 0x43,
	0x5f, 0x50, 0x52, 0x45, 0x53, 0x55, 0x42, 0x4d, 0x49, 0x54, 0x10, 0x07, 0x12, 0x12, 0x0a, 0x0e,
	0x44, 0x55, 0x54, 0x5f, 0x50, 0x4f, 0x4f, 0x4c, 0x5f, 0x51, 0x55, 0x4f, 0x54, 0x41, 0x10, 0x08,
	0x22, 0xc8, 0x01, 0x0a, 0x0a, 0x4c, 0x61, 0x62, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x46, 0x0a, 0x06, 0x73,
	0x65, 0x72, 0x76, 0x6f, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x75, 0x6e,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f,
	0x73, 0x2e, 0x6c, 0x61, 0x62, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x6f, 0x52, 0x06, 0x73, 0x65, 0x72,
	0x76, 0x6f, 0x73, 0x12, 0x40, 0x0a, 0x03, 0x72, 0x70, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x2e, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68,
	0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6c, 0x61, 0x62, 0x2e, 0x4f, 0x53, 0x52, 0x50, 0x4d,
	0x52, 0x03, 0x72, 0x70, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x6f, 0x6f, 0x6c, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x70, 0x6f, 0x6f, 0x6c, 0x73, 0x42, 0x35, 0x5a, 0x33, 0x69,
	0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65,
	0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f,
	0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x6c, 0x61, 0x62, 0x3b, 0x75, 0x66, 0x73,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescOnce sync.Once
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescData = file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDesc
)

func file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescGZIP() []byte {
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescOnce.Do(func() {
		file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescData)
	})
	return file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDescData
}

var file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_goTypes = []interface{}{
	(DeviceUnderTest_DUTPool)(0),   // 0: unifiedfleet.api.v1.models.chromeos.lab.DeviceUnderTest.DUTPool
	(*ChromeOSDevice)(nil),         // 1: unifiedfleet.api.v1.models.chromeos.lab.ChromeOSDevice
	(*DeviceUnderTest)(nil),        // 2: unifiedfleet.api.v1.models.chromeos.lab.DeviceUnderTest
	(*Labstation)(nil),             // 3: unifiedfleet.api.v1.models.chromeos.lab.Labstation
	(*ChromeOSDeviceID)(nil),       // 4: unifiedfleet.api.v1.models.chromeos.lab.ChromeOSDeviceID
	(*manufacturing.ConfigID)(nil), // 5: unifiedfleet.api.v1.models.chromeos.manufacturing.ConfigID
	(*device.ConfigId)(nil),        // 6: unifiedfleet.api.v1.models.chromeos.device.ConfigId
	(*Peripherals)(nil),            // 7: unifiedfleet.api.v1.models.chromeos.lab.Peripherals
	(*License)(nil),                // 8: unifiedfleet.api.v1.models.chromeos.lab.License
	(*Servo)(nil),                  // 9: unifiedfleet.api.v1.models.chromeos.lab.Servo
	(*OSRPM)(nil),                  // 10: unifiedfleet.api.v1.models.chromeos.lab.OSRPM
}
var file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_depIdxs = []int32{
	4,  // 0: unifiedfleet.api.v1.models.chromeos.lab.ChromeOSDevice.id:type_name -> unifiedfleet.api.v1.models.chromeos.lab.ChromeOSDeviceID
	5,  // 1: unifiedfleet.api.v1.models.chromeos.lab.ChromeOSDevice.manufacturing_id:type_name -> unifiedfleet.api.v1.models.chromeos.manufacturing.ConfigID
	6,  // 2: unifiedfleet.api.v1.models.chromeos.lab.ChromeOSDevice.device_config_id:type_name -> unifiedfleet.api.v1.models.chromeos.device.ConfigId
	2,  // 3: unifiedfleet.api.v1.models.chromeos.lab.ChromeOSDevice.dut:type_name -> unifiedfleet.api.v1.models.chromeos.lab.DeviceUnderTest
	3,  // 4: unifiedfleet.api.v1.models.chromeos.lab.ChromeOSDevice.labstation:type_name -> unifiedfleet.api.v1.models.chromeos.lab.Labstation
	7,  // 5: unifiedfleet.api.v1.models.chromeos.lab.DeviceUnderTest.peripherals:type_name -> unifiedfleet.api.v1.models.chromeos.lab.Peripherals
	0,  // 6: unifiedfleet.api.v1.models.chromeos.lab.DeviceUnderTest.critical_pools:type_name -> unifiedfleet.api.v1.models.chromeos.lab.DeviceUnderTest.DUTPool
	8,  // 7: unifiedfleet.api.v1.models.chromeos.lab.DeviceUnderTest.licenses:type_name -> unifiedfleet.api.v1.models.chromeos.lab.License
	9,  // 8: unifiedfleet.api.v1.models.chromeos.lab.Labstation.servos:type_name -> unifiedfleet.api.v1.models.chromeos.lab.Servo
	10, // 9: unifiedfleet.api.v1.models.chromeos.lab.Labstation.rpm:type_name -> unifiedfleet.api.v1.models.chromeos.lab.OSRPM
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_init() }
func file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_init() {
	if File_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto != nil {
		return
	}
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_chromeos_device_id_proto_init()
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_license_proto_init()
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_peripherals_proto_init()
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_servo_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChromeOSDevice); i {
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
		file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceUnderTest); i {
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
		file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Labstation); i {
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
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*ChromeOSDevice_Dut)(nil),
		(*ChromeOSDevice_Labstation)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_goTypes,
		DependencyIndexes: file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_depIdxs,
		EnumInfos:         file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_enumTypes,
		MessageInfos:      file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_msgTypes,
	}.Build()
	File_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto = out.File
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_rawDesc = nil
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_goTypes = nil
	file_infra_unifiedfleet_api_v1_models_chromeos_lab_device_proto_depIdxs = nil
}
