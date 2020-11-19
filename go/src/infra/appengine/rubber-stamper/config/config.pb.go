// Copyright 2020 The Chromium Authors. All Rights Reserved.
// Use of this source code is governed by the Apache v2.0 license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: infra/appengine/rubber-stamper/config/config.proto

package config

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

// Config is the service-wide configuration data for rubber-stamper.
type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A map stores configs for all the Gerrit hosts, where keys are names of
	// hosts (e.g. "chromium" or "chrome-internal"), values are corresponding
	// configs.
	HostConfigs map[string]*HostConfig `protobuf:"bytes,1,rep,name=host_configs,json=hostConfigs,proto3" json:"host_configs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_infra_appengine_rubber_stamper_config_config_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetHostConfigs() map[string]*HostConfig {
	if x != nil {
		return x.HostConfigs
	}
	return nil
}

// HostConfig describes the config to be used for a Gerrit host.
type HostConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BenignFilePattern *BenignFilePattern `protobuf:"bytes,1,opt,name=benign_file_pattern,json=benignFilePattern,proto3" json:"benign_file_pattern,omitempty"`
}

func (x *HostConfig) Reset() {
	*x = HostConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HostConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HostConfig) ProtoMessage() {}

func (x *HostConfig) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HostConfig.ProtoReflect.Descriptor instead.
func (*HostConfig) Descriptor() ([]byte, []int) {
	return file_infra_appengine_rubber_stamper_config_config_proto_rawDescGZIP(), []int{1}
}

func (x *HostConfig) GetBenignFilePattern() *BenignFilePattern {
	if x != nil {
		return x.BenignFilePattern
	}
	return nil
}

// BenignFilePattern describes pattern of changes to benign files.
type BenignFilePattern struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A map contains the information that which files are allowed under which
	// directories, where keys are file extensions, values are paths of files
	// or directories. For paths to specific files, these files can be considered
	// benign files; for paths to directories, files under these directories with
	// corresponding extensions can be considered as benign files. For files with
	// no extensions, their key should be an empty string "".
	FileExtensionMap map[string]*Paths `protobuf:"bytes,1,rep,name=file_extension_map,json=fileExtensionMap,proto3" json:"file_extension_map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *BenignFilePattern) Reset() {
	*x = BenignFilePattern{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BenignFilePattern) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BenignFilePattern) ProtoMessage() {}

func (x *BenignFilePattern) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BenignFilePattern.ProtoReflect.Descriptor instead.
func (*BenignFilePattern) Descriptor() ([]byte, []int) {
	return file_infra_appengine_rubber_stamper_config_config_proto_rawDescGZIP(), []int{2}
}

func (x *BenignFilePattern) GetFileExtensionMap() map[string]*Paths {
	if x != nil {
		return x.FileExtensionMap
	}
	return nil
}

// Paths contains a list of allowed paths.
type Paths struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Could be either a file or a directory. Directory path should end with '/'
	// or "/*". If end with '/', it only allows direct descendants of the
	// directory; if end with "/*", it means that all files under this path are
	// considered benign.
	Paths []string `protobuf:"bytes,1,rep,name=paths,proto3" json:"paths,omitempty"`
}

func (x *Paths) Reset() {
	*x = Paths{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Paths) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Paths) ProtoMessage() {}

func (x *Paths) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Paths.ProtoReflect.Descriptor instead.
func (*Paths) Descriptor() ([]byte, []int) {
	return file_infra_appengine_rubber_stamper_config_config_proto_rawDescGZIP(), []int{3}
}

func (x *Paths) GetPaths() []string {
	if x != nil {
		return x.Paths
	}
	return nil
}

var File_infra_appengine_rubber_stamper_config_config_proto protoreflect.FileDescriptor

var file_infra_appengine_rubber_stamper_config_config_proto_rawDesc = []byte{
	0x0a, 0x32, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2f, 0x72, 0x75, 0x62, 0x62, 0x65, 0x72, 0x2d, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x65, 0x72,
	0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x72, 0x75, 0x62, 0x62, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0xbe, 0x01, 0x0a, 0x06,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x51, 0x0a, 0x0c, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x72,
	0x75, 0x62, 0x62, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x65, 0x72, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x48, 0x6f, 0x73, 0x74,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x68, 0x6f,
	0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x1a, 0x61, 0x0a, 0x10, 0x48, 0x6f, 0x73,
	0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x37, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x72, 0x75, 0x62, 0x62, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x65, 0x72, 0x2e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x48, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x66, 0x0a, 0x0a,
	0x48, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x58, 0x0a, 0x13, 0x62, 0x65,
	0x6e, 0x69, 0x67, 0x6e, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x72, 0x75, 0x62, 0x62, 0x65, 0x72,
	0x5f, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x42, 0x65, 0x6e, 0x69, 0x67, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72,
	0x6e, 0x52, 0x11, 0x62, 0x65, 0x6e, 0x69, 0x67, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74,
	0x74, 0x65, 0x72, 0x6e, 0x22, 0xe4, 0x01, 0x0a, 0x11, 0x42, 0x65, 0x6e, 0x69, 0x67, 0x6e, 0x46,
	0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x12, 0x6c, 0x0a, 0x12, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x70,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3e, 0x2e, 0x72, 0x75, 0x62, 0x62, 0x65, 0x72, 0x5f,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x42,
	0x65, 0x6e, 0x69, 0x67, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x4d, 0x61,
	0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x10, 0x66, 0x69, 0x6c, 0x65, 0x45, 0x78, 0x74, 0x65,
	0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x70, 0x1a, 0x61, 0x0a, 0x15, 0x46, 0x69, 0x6c, 0x65,
	0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x32, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x72, 0x75, 0x62, 0x62, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x50, 0x61, 0x74, 0x68, 0x73,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x1d, 0x0a, 0x05, 0x50,
	0x61, 0x74, 0x68, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x61, 0x74, 0x68, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x05, 0x70, 0x61, 0x74, 0x68, 0x73, 0x42, 0x27, 0x5a, 0x25, 0x69, 0x6e,
	0x66, 0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x72, 0x75,
	0x62, 0x62, 0x65, 0x72, 0x2d, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x65, 0x72, 0x2f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_appengine_rubber_stamper_config_config_proto_rawDescOnce sync.Once
	file_infra_appengine_rubber_stamper_config_config_proto_rawDescData = file_infra_appengine_rubber_stamper_config_config_proto_rawDesc
)

func file_infra_appengine_rubber_stamper_config_config_proto_rawDescGZIP() []byte {
	file_infra_appengine_rubber_stamper_config_config_proto_rawDescOnce.Do(func() {
		file_infra_appengine_rubber_stamper_config_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_appengine_rubber_stamper_config_config_proto_rawDescData)
	})
	return file_infra_appengine_rubber_stamper_config_config_proto_rawDescData
}

var file_infra_appengine_rubber_stamper_config_config_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_infra_appengine_rubber_stamper_config_config_proto_goTypes = []interface{}{
	(*Config)(nil),            // 0: rubber_stamper.config.Config
	(*HostConfig)(nil),        // 1: rubber_stamper.config.HostConfig
	(*BenignFilePattern)(nil), // 2: rubber_stamper.config.BenignFilePattern
	(*Paths)(nil),             // 3: rubber_stamper.config.Paths
	nil,                       // 4: rubber_stamper.config.Config.HostConfigsEntry
	nil,                       // 5: rubber_stamper.config.BenignFilePattern.FileExtensionMapEntry
}
var file_infra_appengine_rubber_stamper_config_config_proto_depIdxs = []int32{
	4, // 0: rubber_stamper.config.Config.host_configs:type_name -> rubber_stamper.config.Config.HostConfigsEntry
	2, // 1: rubber_stamper.config.HostConfig.benign_file_pattern:type_name -> rubber_stamper.config.BenignFilePattern
	5, // 2: rubber_stamper.config.BenignFilePattern.file_extension_map:type_name -> rubber_stamper.config.BenignFilePattern.FileExtensionMapEntry
	1, // 3: rubber_stamper.config.Config.HostConfigsEntry.value:type_name -> rubber_stamper.config.HostConfig
	3, // 4: rubber_stamper.config.BenignFilePattern.FileExtensionMapEntry.value:type_name -> rubber_stamper.config.Paths
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_infra_appengine_rubber_stamper_config_config_proto_init() }
func file_infra_appengine_rubber_stamper_config_config_proto_init() {
	if File_infra_appengine_rubber_stamper_config_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
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
		file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HostConfig); i {
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
		file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BenignFilePattern); i {
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
		file_infra_appengine_rubber_stamper_config_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Paths); i {
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
			RawDescriptor: file_infra_appengine_rubber_stamper_config_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_appengine_rubber_stamper_config_config_proto_goTypes,
		DependencyIndexes: file_infra_appengine_rubber_stamper_config_config_proto_depIdxs,
		MessageInfos:      file_infra_appengine_rubber_stamper_config_config_proto_msgTypes,
	}.Build()
	File_infra_appengine_rubber_stamper_config_config_proto = out.File
	file_infra_appengine_rubber_stamper_config_config_proto_rawDesc = nil
	file_infra_appengine_rubber_stamper_config_config_proto_goTypes = nil
	file_infra_appengine_rubber_stamper_config_config_proto_depIdxs = nil
}
