// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: infra/appengine/weetbix/proto/v1/common.proto

package weetbixpb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Status of a Verdict.
// It is determined by all the test results of the verdict, and exonerations are
// ignored(i.e. failure is treated as a failure, even if it is exonerated).
type VerdictStatus int32

const (
	// A verdict must not have this status.
	// This is only used when filtering verdicts.
	VerdictStatus_VERDICT_STATUS_UNSPECIFIED VerdictStatus = 0
	// All results of the verdict are unexpected.
	VerdictStatus_UNEXPECTED VerdictStatus = 10
	// The verdict has both expected and unexpected results.
	// To be differentiated with AnalyzedTestVariantStatus.FLAKY.
	VerdictStatus_VERDICT_FLAKY VerdictStatus = 30
	// All results of the verdict are expected.
	VerdictStatus_EXPECTED VerdictStatus = 50
)

// Enum value maps for VerdictStatus.
var (
	VerdictStatus_name = map[int32]string{
		0:  "VERDICT_STATUS_UNSPECIFIED",
		10: "UNEXPECTED",
		30: "VERDICT_FLAKY",
		50: "EXPECTED",
	}
	VerdictStatus_value = map[string]int32{
		"VERDICT_STATUS_UNSPECIFIED": 0,
		"UNEXPECTED":                 10,
		"VERDICT_FLAKY":              30,
		"EXPECTED":                   50,
	}
)

func (x VerdictStatus) Enum() *VerdictStatus {
	p := new(VerdictStatus)
	*p = x
	return p
}

func (x VerdictStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VerdictStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_infra_appengine_weetbix_proto_v1_common_proto_enumTypes[0].Descriptor()
}

func (VerdictStatus) Type() protoreflect.EnumType {
	return &file_infra_appengine_weetbix_proto_v1_common_proto_enumTypes[0]
}

func (x VerdictStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VerdictStatus.Descriptor instead.
func (VerdictStatus) EnumDescriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP(), []int{0}
}

// ExonerationStatus explains if and why a test failure was
// exonerated. Exonerated means the failure was ignored and did not
// have further impact, in terms of causing the build to fail or
// rejecting the CL being tested in a presubmit run.
type ExonerationStatus int32

const (
	// A test failure must not have this status.
	ExonerationStatus_EXONERATION_STATUS_UNSPECIFIED ExonerationStatus = 0
	// The test was not exonerated.
	ExonerationStatus_NOT_EXONERATED ExonerationStatus = 1
	// The unexpected failure was discounted despite
	// having an unexpected result and no exoneration recorded
	// in Result DB. For example, because the build passed or
	// was cancelled.
	ExonerationStatus_IMPLICIT ExonerationStatus = 2
	// The test was marked exonerated in ResultDB, for a reason
	// other than Weetbix or FindIt failure analysis.
	// If a test is exonerated in ResultDB for both reasons
	// other than Weetbix/FindIt and because of Weetbix/FindIt,
	// this status takes precedence.
	ExonerationStatus_EXPLICIT ExonerationStatus = 3
	// The test was exonerated based on Weetbix cluster analysis.
	// This status is only set if Weetbix is the only explicit
	// reason(s) given for the exoneration in ResultDB.
	//
	// This status is provided to avoid feedback loops in the
	// cluster analysis performed by Weetbix, by allowing Weetbix to
	// filter out exoneration decisions based on such analysis from
	// feeding back into the input of the analysis.
	//
	// Example of a situation we want to avoid:
	// Weetbix detects an impactful cluster of failures
	// affecting multiple CLs and cause a recipe to exonerate it.
	// As a result, Weetbix no longer detects the cluster as impactful.
	// As a result, the cluster is no longer exonerated.
	// As a result, the impact resumes.
	//
	// During the transition from FindIt to Weetbix for failure
	// exoneration, exonerations caused by FindIt will be treated
	// the same as exonerations caused by Weetbix, to ensure Weetbix
	// behaves as if FindIt no longer exists.
	ExonerationStatus_WEETBIX ExonerationStatus = 4
)

// Enum value maps for ExonerationStatus.
var (
	ExonerationStatus_name = map[int32]string{
		0: "EXONERATION_STATUS_UNSPECIFIED",
		1: "NOT_EXONERATED",
		2: "IMPLICIT",
		3: "EXPLICIT",
		4: "WEETBIX",
	}
	ExonerationStatus_value = map[string]int32{
		"EXONERATION_STATUS_UNSPECIFIED": 0,
		"NOT_EXONERATED":                 1,
		"IMPLICIT":                       2,
		"EXPLICIT":                       3,
		"WEETBIX":                        4,
	}
)

func (x ExonerationStatus) Enum() *ExonerationStatus {
	p := new(ExonerationStatus)
	*p = x
	return p
}

func (x ExonerationStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ExonerationStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_infra_appengine_weetbix_proto_v1_common_proto_enumTypes[1].Descriptor()
}

func (ExonerationStatus) Type() protoreflect.EnumType {
	return &file_infra_appengine_weetbix_proto_v1_common_proto_enumTypes[1]
}

func (x ExonerationStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ExonerationStatus.Descriptor instead.
func (ExonerationStatus) EnumDescriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP(), []int{1}
}

// A range of timestamps.
type TimeRange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The oldest timestamp to include in the range.
	Earliest *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=earliest,proto3" json:"earliest,omitempty"`
	// Include only timestamps that are strictly older than this.
	Latest *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=latest,proto3" json:"latest,omitempty"`
}

func (x *TimeRange) Reset() {
	*x = TimeRange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeRange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeRange) ProtoMessage() {}

func (x *TimeRange) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeRange.ProtoReflect.Descriptor instead.
func (*TimeRange) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP(), []int{0}
}

func (x *TimeRange) GetEarliest() *timestamppb.Timestamp {
	if x != nil {
		return x.Earliest
	}
	return nil
}

func (x *TimeRange) GetLatest() *timestamppb.Timestamp {
	if x != nil {
		return x.Latest
	}
	return nil
}

// Identity of a test result.
type TestResultId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The test results system.
	// Currently, the only valid value is "resultdb".
	System string `protobuf:"bytes,1,opt,name=system,proto3" json:"system,omitempty"`
	// ID for the test result in the test results system.
	// For test results in ResultDB, the format is:
	// "invocations/{INVOCATION_ID}/tests/{URL_ESCAPED_TEST_ID}/results/{RESULT_ID}"
	// Where INVOCATION_ID, URL_ESCAPED_TEST_ID and RESULT_ID are values defined
	// in ResultDB.
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *TestResultId) Reset() {
	*x = TestResultId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestResultId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestResultId) ProtoMessage() {}

func (x *TestResultId) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestResultId.ProtoReflect.Descriptor instead.
func (*TestResultId) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP(), []int{1}
}

func (x *TestResultId) GetSystem() string {
	if x != nil {
		return x.System
	}
	return ""
}

func (x *TestResultId) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Variant represents a way of running a test case.
//
// The same test case can be executed in different ways, for example on
// different OS, GPUs, with different compile options or runtime flags.
type Variant struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The definition of the variant. Each key-value pair represents a
	// parameter describing how the test was run (e.g. OS, GPU, etc.).
	Def map[string]string `protobuf:"bytes,1,rep,name=def,proto3" json:"def,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Variant) Reset() {
	*x = Variant{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Variant) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Variant) ProtoMessage() {}

func (x *Variant) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Variant.ProtoReflect.Descriptor instead.
func (*Variant) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP(), []int{2}
}

func (x *Variant) GetDef() map[string]string {
	if x != nil {
		return x.Def
	}
	return nil
}

type StringPair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Regex: ^[a-z][a-z0-9_]*(/[a-z][a-z0-9_]*)*$
	// Max length: 64.
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Max length: 256.
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *StringPair) Reset() {
	*x = StringPair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringPair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringPair) ProtoMessage() {}

func (x *StringPair) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringPair.ProtoReflect.Descriptor instead.
func (*StringPair) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP(), []int{3}
}

func (x *StringPair) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *StringPair) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// Identity of a bug tracking component in a bug tracking system.
type BugTrackingComponent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The bug tracking system corresponding to this test case, as identified
	// by the test results system.
	// Currently, the only valid value is "monorail".
	System string `protobuf:"bytes,1,opt,name=system,proto3" json:"system,omitempty"`
	// The bug tracking component corresponding to this test case, as identified
	// by the test results system.
	// If the bug tracking system is monorail, this is the component as the
	// user would see it, e.g. "Infra>Test>Flakiness". For monorail, the bug
	// tracking project (e.g. "chromium") is not encoded, but assumed to be
	// specified in the project's Weetbix configuration.
	Component string `protobuf:"bytes,2,opt,name=component,proto3" json:"component,omitempty"`
}

func (x *BugTrackingComponent) Reset() {
	*x = BugTrackingComponent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BugTrackingComponent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BugTrackingComponent) ProtoMessage() {}

func (x *BugTrackingComponent) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BugTrackingComponent.ProtoReflect.Descriptor instead.
func (*BugTrackingComponent) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP(), []int{4}
}

func (x *BugTrackingComponent) GetSystem() string {
	if x != nil {
		return x.System
	}
	return ""
}

func (x *BugTrackingComponent) GetComponent() string {
	if x != nil {
		return x.Component
	}
	return ""
}

// Identity of a presubmit run (also known as a "CQ Run" or "CV Run").
type PresubmitRunId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The system that was used to process the presubmit run.
	// Currently, the only valid value is "luci-cv" for LUCI Commit Verifier
	// (LUCI CV).
	System string `protobuf:"bytes,1,opt,name=system,proto3" json:"system,omitempty"`
	// Identity of the presubmit run.
	// If the presubmit system is LUCI CV, the format of this value is:
	//   "{LUCI_PROJECT}/{LUCI_CV_ID}", e.g.
	//   "infra/8988819463854-1-f94732fe20056fd1".
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PresubmitRunId) Reset() {
	*x = PresubmitRunId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PresubmitRunId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresubmitRunId) ProtoMessage() {}

func (x *PresubmitRunId) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresubmitRunId.ProtoReflect.Descriptor instead.
func (*PresubmitRunId) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP(), []int{5}
}

func (x *PresubmitRunId) GetSystem() string {
	if x != nil {
		return x.System
	}
	return ""
}

func (x *PresubmitRunId) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Identity of a bug in a bug-tracking system.
type AssociatedBug struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// System is the bug tracking system of the bug. This is either
	// "monorail" or "buganizer".
	System string `protobuf:"bytes,1,opt,name=system,proto3" json:"system,omitempty"`
	// Id is the bug tracking system-specific identity of the bug.
	// For monorail, the scheme is {project}/{numeric_id}, for
	// buganizer the scheme is {numeric_id}.
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// A human-readable name for the bug. This is typically the
	// bug shortlink (e.g. "crbug.com/1234567").
	LinkText string `protobuf:"bytes,3,opt,name=link_text,json=linkText,proto3" json:"link_text,omitempty"`
	// The resolved bug URL, e.g.
	// E.g. "https://bugs.chromium.org/p/chromium/issues/detail?id=123456".
	Url string `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *AssociatedBug) Reset() {
	*x = AssociatedBug{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssociatedBug) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssociatedBug) ProtoMessage() {}

func (x *AssociatedBug) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssociatedBug.ProtoReflect.Descriptor instead.
func (*AssociatedBug) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP(), []int{6}
}

func (x *AssociatedBug) GetSystem() string {
	if x != nil {
		return x.System
	}
	return ""
}

func (x *AssociatedBug) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AssociatedBug) GetLinkText() string {
	if x != nil {
		return x.LinkText
	}
	return ""
}

func (x *AssociatedBug) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

// ClusterId represents the identity of a cluster. The LUCI Project is
// omitted as it is assumed to be implicit from the context.
//
// This is often used in place of the resource name of the cluster
// (in the sense of https://google.aip.dev/122) as clients may need
// to access individual parts of the resource name (e.g. to determine
// the algorithm used) and it is not desirable to make clients parse
// the resource name.
type ClusterId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Algorithm is the name of the clustering algorithm that identified
	// the cluster.
	Algorithm string `protobuf:"bytes,1,opt,name=algorithm,proto3" json:"algorithm,omitempty"`
	// Id is the cluster identifier returned by the algorithm. The underlying
	// identifier is at most 16 bytes, but is represented here as a hexadecimal
	// string of up to 32 lowercase hexadecimal characters.
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ClusterId) Reset() {
	*x = ClusterId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterId) ProtoMessage() {}

func (x *ClusterId) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterId.ProtoReflect.Descriptor instead.
func (*ClusterId) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP(), []int{7}
}

func (x *ClusterId) GetAlgorithm() string {
	if x != nil {
		return x.Algorithm
	}
	return ""
}

func (x *ClusterId) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_infra_appengine_weetbix_proto_v1_common_proto protoreflect.FileDescriptor

var file_infra_appengine_weetbix_proto_v1_common_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2f, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65,
	0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x77, 0x0a,
	0x09, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x08, 0x65, 0x61,
	0x72, 0x6c, 0x69, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x65, 0x61, 0x72, 0x6c, 0x69, 0x65,
	0x73, 0x74, 0x12, 0x32, 0x0a, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06,
	0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x22, 0x36, 0x0a, 0x0c, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x71,
	0x0a, 0x07, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x03, 0x64, 0x65, 0x66,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78,
	0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x2e, 0x44, 0x65, 0x66, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x64, 0x65, 0x66, 0x1a, 0x36, 0x0a, 0x08, 0x44, 0x65, 0x66,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x34, 0x0a, 0x0a, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x50, 0x61, 0x69, 0x72, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x4c, 0x0a, 0x14, 0x42, 0x75, 0x67, 0x54, 0x72,
	0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6f,
	0x6e, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x70,
	0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x22, 0x38, 0x0a, 0x0e, 0x50, 0x72, 0x65, 0x73, 0x75, 0x62, 0x6d,
	0x69, 0x74, 0x52, 0x75, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x70, 0x0a, 0x0d, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x64, 0x42, 0x75, 0x67,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x09, 0x6c, 0x69, 0x6e, 0x6b,
	0x5f, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x03,
	0x52, 0x08, 0x6c, 0x69, 0x6e, 0x6b, 0x54, 0x65, 0x78, 0x74, 0x12, 0x15, 0x0a, 0x03, 0x75, 0x72,
	0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x03, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x22, 0x39, 0x0a, 0x09, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x2a, 0x60, 0x0a, 0x0d,
	0x56, 0x65, 0x72, 0x64, 0x69, 0x63, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a,
	0x1a, 0x56, 0x45, 0x52, 0x44, 0x49, 0x43, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a,
	0x0a, 0x55, 0x4e, 0x45, 0x58, 0x50, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x0a, 0x12, 0x11, 0x0a,
	0x0d, 0x56, 0x45, 0x52, 0x44, 0x49, 0x43, 0x54, 0x5f, 0x46, 0x4c, 0x41, 0x4b, 0x59, 0x10, 0x1e,
	0x12, 0x0c, 0x0a, 0x08, 0x45, 0x58, 0x50, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x32, 0x2a, 0x74,
	0x0a, 0x11, 0x45, 0x78, 0x6f, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x22, 0x0a, 0x1e, 0x45, 0x58, 0x4f, 0x4e, 0x45, 0x52, 0x41, 0x54, 0x49,
	0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x4e, 0x4f, 0x54, 0x5f, 0x45,
	0x58, 0x4f, 0x4e, 0x45, 0x52, 0x41, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x49,
	0x4d, 0x50, 0x4c, 0x49, 0x43, 0x49, 0x54, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x45, 0x58, 0x50,
	0x4c, 0x49, 0x43, 0x49, 0x54, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x45, 0x45, 0x54, 0x42,
	0x49, 0x58, 0x10, 0x04, 0x42, 0x2c, 0x5a, 0x2a, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x61, 0x70,
	0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x3b, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_appengine_weetbix_proto_v1_common_proto_rawDescOnce sync.Once
	file_infra_appengine_weetbix_proto_v1_common_proto_rawDescData = file_infra_appengine_weetbix_proto_v1_common_proto_rawDesc
)

func file_infra_appengine_weetbix_proto_v1_common_proto_rawDescGZIP() []byte {
	file_infra_appengine_weetbix_proto_v1_common_proto_rawDescOnce.Do(func() {
		file_infra_appengine_weetbix_proto_v1_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_appengine_weetbix_proto_v1_common_proto_rawDescData)
	})
	return file_infra_appengine_weetbix_proto_v1_common_proto_rawDescData
}

var file_infra_appengine_weetbix_proto_v1_common_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_infra_appengine_weetbix_proto_v1_common_proto_goTypes = []interface{}{
	(VerdictStatus)(0),            // 0: weetbix.v1.VerdictStatus
	(ExonerationStatus)(0),        // 1: weetbix.v1.ExonerationStatus
	(*TimeRange)(nil),             // 2: weetbix.v1.TimeRange
	(*TestResultId)(nil),          // 3: weetbix.v1.TestResultId
	(*Variant)(nil),               // 4: weetbix.v1.Variant
	(*StringPair)(nil),            // 5: weetbix.v1.StringPair
	(*BugTrackingComponent)(nil),  // 6: weetbix.v1.BugTrackingComponent
	(*PresubmitRunId)(nil),        // 7: weetbix.v1.PresubmitRunId
	(*AssociatedBug)(nil),         // 8: weetbix.v1.AssociatedBug
	(*ClusterId)(nil),             // 9: weetbix.v1.ClusterId
	nil,                           // 10: weetbix.v1.Variant.DefEntry
	(*timestamppb.Timestamp)(nil), // 11: google.protobuf.Timestamp
}
var file_infra_appengine_weetbix_proto_v1_common_proto_depIdxs = []int32{
	11, // 0: weetbix.v1.TimeRange.earliest:type_name -> google.protobuf.Timestamp
	11, // 1: weetbix.v1.TimeRange.latest:type_name -> google.protobuf.Timestamp
	10, // 2: weetbix.v1.Variant.def:type_name -> weetbix.v1.Variant.DefEntry
	3,  // [3:3] is the sub-list for method output_type
	3,  // [3:3] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_infra_appengine_weetbix_proto_v1_common_proto_init() }
func file_infra_appengine_weetbix_proto_v1_common_proto_init() {
	if File_infra_appengine_weetbix_proto_v1_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeRange); i {
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
		file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestResultId); i {
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
		file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Variant); i {
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
		file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringPair); i {
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
		file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BugTrackingComponent); i {
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
		file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PresubmitRunId); i {
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
		file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssociatedBug); i {
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
		file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterId); i {
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
			RawDescriptor: file_infra_appengine_weetbix_proto_v1_common_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_appengine_weetbix_proto_v1_common_proto_goTypes,
		DependencyIndexes: file_infra_appengine_weetbix_proto_v1_common_proto_depIdxs,
		EnumInfos:         file_infra_appengine_weetbix_proto_v1_common_proto_enumTypes,
		MessageInfos:      file_infra_appengine_weetbix_proto_v1_common_proto_msgTypes,
	}.Build()
	File_infra_appengine_weetbix_proto_v1_common_proto = out.File
	file_infra_appengine_weetbix_proto_v1_common_proto_rawDesc = nil
	file_infra_appengine_weetbix_proto_v1_common_proto_goTypes = nil
	file_infra_appengine_weetbix_proto_v1_common_proto_depIdxs = nil
}
