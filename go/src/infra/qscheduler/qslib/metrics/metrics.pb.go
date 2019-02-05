// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/qscheduler/qslib/metrics/metrics.proto

package metrics

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type TaskEvent_EventType int32

const (
	// Invalid or unspecified event type.
	TaskEvent_UNSPECIFIED TaskEvent_EventType = 0
	// Task was enqueued.
	TaskEvent_ENQUEUED TaskEvent_EventType = 1
	// Task was assigned to a bot.
	TaskEvent_ASSIGNED TaskEvent_EventType = 2
	// Task (which was previously assigned to a bot) was interrupted by another
	// task.
	TaskEvent_PREEMPTED TaskEvent_EventType = 3
	// Task (which was previously assigned to a bot) changed its running
	// priority.
	TaskEvent_REPRIORITIZED TaskEvent_EventType = 4
	// Task (which was previously assigned to a bot) completed on that bot,
	// because the bot reported itself as idle.
	TaskEvent_COMPLETED TaskEvent_EventType = 5
)

var TaskEvent_EventType_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "ENQUEUED",
	2: "ASSIGNED",
	3: "PREEMPTED",
	4: "REPRIORITIZED",
	5: "COMPLETED",
}

var TaskEvent_EventType_value = map[string]int32{
	"UNSPECIFIED":   0,
	"ENQUEUED":      1,
	"ASSIGNED":      2,
	"PREEMPTED":     3,
	"REPRIORITIZED": 4,
	"COMPLETED":     5,
}

func (x TaskEvent_EventType) String() string {
	return proto.EnumName(TaskEvent_EventType_name, int32(x))
}

func (TaskEvent_EventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3f5c8147aadbdf8a, []int{0, 0}
}

// TaskEvent represents a quotascheduler event that happened to a particular
// task, for metrics purposes.
//
// This proto is intended to be used as a schema for a BigQuery table, in which
// events are logged.
type TaskEvent struct {
	// EventType is the type of event that occurred.
	EventType TaskEvent_EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=metrics.TaskEvent_EventType" json:"event_type,omitempty"`
	// SchedulerId is the ID of the scheduler in which the event occurred.
	SchedulerId string `protobuf:"bytes,2,opt,name=scheduler_id,json=schedulerId,proto3" json:"scheduler_id,omitempty"`
	// TaskId is the task ID that the event happened to.
	TaskId string `protobuf:"bytes,3,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	// Time is the time at which the event happened.
	Time *timestamp.Timestamp `protobuf:"bytes,4,opt,name=time,proto3" json:"time,omitempty"`
	// BaseLabels are the base labels of the task.
	BaseLabels []string `protobuf:"bytes,5,rep,name=base_labels,json=baseLabels,proto3" json:"base_labels,omitempty"`
	// ProvisionableLabels are the provisionable labels of the task.
	ProvisionableLabels []string `protobuf:"bytes,6,rep,name=provisionable_labels,json=provisionableLabels,proto3" json:"provisionable_labels,omitempty"`
	// AccountId is the quotascheduler account that the task will be charged to.
	AccountId string `protobuf:"bytes,7,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// AccountValid indicates whether this task's quotascheduler account is valid
	// (exists) at this time.
	AccountValid bool `protobuf:"varint,8,opt,name=account_valid,json=accountValid,proto3" json:"account_valid,omitempty"`
	// AccountBalance is the task's quotascheduler account's balance at this time.
	AccountBalance []float64 `protobuf:"fixed64,9,rep,packed,name=account_balance,json=accountBalance,proto3" json:"account_balance,omitempty"`
	// Cost is the total quota cost accumulated so far to the quotascheduler
	// account, when running this task.
	Cost []float64 `protobuf:"fixed64,10,rep,packed,name=cost,proto3" json:"cost,omitempty"`
	// BotId is the bot that the event occurred on (relevant for all event
	// types except for ENQUEUED).
	BotId string `protobuf:"bytes,11,opt,name=bot_id,json=botId,proto3" json:"bot_id,omitempty"`
	// BotDimensions are the dimensions of the bot (if relevant).
	BotDimensions []string `protobuf:"bytes,12,rep,name=bot_dimensions,json=botDimensions,proto3" json:"bot_dimensions,omitempty"`
	// Types that are valid to be assigned to Details:
	//	*TaskEvent_EnqueuedDetails_
	//	*TaskEvent_AssignedDetails_
	//	*TaskEvent_PreemptedDetails_
	//	*TaskEvent_ReprioritizedDetails_
	//	*TaskEvent_CompletedDetails_
	Details              isTaskEvent_Details `protobuf_oneof:"details"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *TaskEvent) Reset()         { *m = TaskEvent{} }
func (m *TaskEvent) String() string { return proto.CompactTextString(m) }
func (*TaskEvent) ProtoMessage()    {}
func (*TaskEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f5c8147aadbdf8a, []int{0}
}

func (m *TaskEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent.Unmarshal(m, b)
}
func (m *TaskEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent.Marshal(b, m, deterministic)
}
func (m *TaskEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent.Merge(m, src)
}
func (m *TaskEvent) XXX_Size() int {
	return xxx_messageInfo_TaskEvent.Size(m)
}
func (m *TaskEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent proto.InternalMessageInfo

func (m *TaskEvent) GetEventType() TaskEvent_EventType {
	if m != nil {
		return m.EventType
	}
	return TaskEvent_UNSPECIFIED
}

func (m *TaskEvent) GetSchedulerId() string {
	if m != nil {
		return m.SchedulerId
	}
	return ""
}

func (m *TaskEvent) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *TaskEvent) GetTime() *timestamp.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *TaskEvent) GetBaseLabels() []string {
	if m != nil {
		return m.BaseLabels
	}
	return nil
}

func (m *TaskEvent) GetProvisionableLabels() []string {
	if m != nil {
		return m.ProvisionableLabels
	}
	return nil
}

func (m *TaskEvent) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *TaskEvent) GetAccountValid() bool {
	if m != nil {
		return m.AccountValid
	}
	return false
}

func (m *TaskEvent) GetAccountBalance() []float64 {
	if m != nil {
		return m.AccountBalance
	}
	return nil
}

func (m *TaskEvent) GetCost() []float64 {
	if m != nil {
		return m.Cost
	}
	return nil
}

func (m *TaskEvent) GetBotId() string {
	if m != nil {
		return m.BotId
	}
	return ""
}

func (m *TaskEvent) GetBotDimensions() []string {
	if m != nil {
		return m.BotDimensions
	}
	return nil
}

type isTaskEvent_Details interface {
	isTaskEvent_Details()
}

type TaskEvent_EnqueuedDetails_ struct {
	EnqueuedDetails *TaskEvent_EnqueuedDetails `protobuf:"bytes,100,opt,name=enqueued_details,json=enqueuedDetails,proto3,oneof"`
}

type TaskEvent_AssignedDetails_ struct {
	AssignedDetails *TaskEvent_AssignedDetails `protobuf:"bytes,101,opt,name=assigned_details,json=assignedDetails,proto3,oneof"`
}

type TaskEvent_PreemptedDetails_ struct {
	PreemptedDetails *TaskEvent_PreemptedDetails `protobuf:"bytes,102,opt,name=preempted_details,json=preemptedDetails,proto3,oneof"`
}

type TaskEvent_ReprioritizedDetails_ struct {
	ReprioritizedDetails *TaskEvent_ReprioritizedDetails `protobuf:"bytes,103,opt,name=reprioritized_details,json=reprioritizedDetails,proto3,oneof"`
}

type TaskEvent_CompletedDetails_ struct {
	CompletedDetails *TaskEvent_CompletedDetails `protobuf:"bytes,104,opt,name=completed_details,json=completedDetails,proto3,oneof"`
}

func (*TaskEvent_EnqueuedDetails_) isTaskEvent_Details() {}

func (*TaskEvent_AssignedDetails_) isTaskEvent_Details() {}

func (*TaskEvent_PreemptedDetails_) isTaskEvent_Details() {}

func (*TaskEvent_ReprioritizedDetails_) isTaskEvent_Details() {}

func (*TaskEvent_CompletedDetails_) isTaskEvent_Details() {}

func (m *TaskEvent) GetDetails() isTaskEvent_Details {
	if m != nil {
		return m.Details
	}
	return nil
}

func (m *TaskEvent) GetEnqueuedDetails() *TaskEvent_EnqueuedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_EnqueuedDetails_); ok {
		return x.EnqueuedDetails
	}
	return nil
}

func (m *TaskEvent) GetAssignedDetails() *TaskEvent_AssignedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_AssignedDetails_); ok {
		return x.AssignedDetails
	}
	return nil
}

func (m *TaskEvent) GetPreemptedDetails() *TaskEvent_PreemptedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_PreemptedDetails_); ok {
		return x.PreemptedDetails
	}
	return nil
}

func (m *TaskEvent) GetReprioritizedDetails() *TaskEvent_ReprioritizedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_ReprioritizedDetails_); ok {
		return x.ReprioritizedDetails
	}
	return nil
}

func (m *TaskEvent) GetCompletedDetails() *TaskEvent_CompletedDetails {
	if x, ok := m.GetDetails().(*TaskEvent_CompletedDetails_); ok {
		return x.CompletedDetails
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*TaskEvent) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TaskEvent_EnqueuedDetails_)(nil),
		(*TaskEvent_AssignedDetails_)(nil),
		(*TaskEvent_PreemptedDetails_)(nil),
		(*TaskEvent_ReprioritizedDetails_)(nil),
		(*TaskEvent_CompletedDetails_)(nil),
	}
}

// EnqueuedDetails represents event details that are specific to the
// ENQUEUED event type.
type TaskEvent_EnqueuedDetails struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent_EnqueuedDetails) Reset()         { *m = TaskEvent_EnqueuedDetails{} }
func (m *TaskEvent_EnqueuedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_EnqueuedDetails) ProtoMessage()    {}
func (*TaskEvent_EnqueuedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f5c8147aadbdf8a, []int{0, 0}
}

func (m *TaskEvent_EnqueuedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_EnqueuedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_EnqueuedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_EnqueuedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_EnqueuedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_EnqueuedDetails.Merge(m, src)
}
func (m *TaskEvent_EnqueuedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_EnqueuedDetails.Size(m)
}
func (m *TaskEvent_EnqueuedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_EnqueuedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_EnqueuedDetails proto.InternalMessageInfo

// AssignedDetails represents event details that are specific to the
// ASSIGNED event type.
type TaskEvent_AssignedDetails struct {
	// ProvisionRequired is whether provision is required to run this task
	// on the bot (i.e. if a slice other than the 0th slice was selected).
	ProvisionRequired bool `protobuf:"varint,1,opt,name=provision_required,json=provisionRequired,proto3" json:"provision_required,omitempty"`
	// Priority is the qscheduler priority that the task is running at.
	Priority int32 `protobuf:"varint,2,opt,name=priority,proto3" json:"priority,omitempty"`
	// Preempting is true if this task preempted another one that was already
	// running on the bot.
	Preempting bool `protobuf:"varint,3,opt,name=preempting,proto3" json:"preempting,omitempty"`
	// PreemptionCost is the cost paid by this task's account in order to
	// preempt the previous task on this bot, if this was a preempting
	// assignment.
	PreemptionCost       []float64 `protobuf:"fixed64,4,rep,packed,name=preemption_cost,json=preemptionCost,proto3" json:"preemption_cost,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *TaskEvent_AssignedDetails) Reset()         { *m = TaskEvent_AssignedDetails{} }
func (m *TaskEvent_AssignedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_AssignedDetails) ProtoMessage()    {}
func (*TaskEvent_AssignedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f5c8147aadbdf8a, []int{0, 1}
}

func (m *TaskEvent_AssignedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_AssignedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_AssignedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_AssignedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_AssignedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_AssignedDetails.Merge(m, src)
}
func (m *TaskEvent_AssignedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_AssignedDetails.Size(m)
}
func (m *TaskEvent_AssignedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_AssignedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_AssignedDetails proto.InternalMessageInfo

func (m *TaskEvent_AssignedDetails) GetProvisionRequired() bool {
	if m != nil {
		return m.ProvisionRequired
	}
	return false
}

func (m *TaskEvent_AssignedDetails) GetPriority() int32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *TaskEvent_AssignedDetails) GetPreempting() bool {
	if m != nil {
		return m.Preempting
	}
	return false
}

func (m *TaskEvent_AssignedDetails) GetPreemptionCost() []float64 {
	if m != nil {
		return m.PreemptionCost
	}
	return nil
}

// PreemptedDetails represents event details that are specific to the
// PREEMPTED event type.
type TaskEvent_PreemptedDetails struct {
	// PreemptingAccountId is the account id of the task that preempted this
	// task.
	PreemptingAccountId string `protobuf:"bytes,1,opt,name=preempting_account_id,json=preemptingAccountId,proto3" json:"preempting_account_id,omitempty"`
	// PreemptingTaskId is the task id of the task that preempted this task.
	PreemptingTaskId string `protobuf:"bytes,2,opt,name=preempting_task_id,json=preemptingTaskId,proto3" json:"preempting_task_id,omitempty"`
	// Priority is the priority that this task was running at prior to being
	// preempted.
	Priority int32 `protobuf:"varint,3,opt,name=priority,proto3" json:"priority,omitempty"`
	// PreemptingPriority is the priority of the task that preempted this task.
	PreemptingPriority   int32    `protobuf:"varint,4,opt,name=preempting_priority,json=preemptingPriority,proto3" json:"preempting_priority,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent_PreemptedDetails) Reset()         { *m = TaskEvent_PreemptedDetails{} }
func (m *TaskEvent_PreemptedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_PreemptedDetails) ProtoMessage()    {}
func (*TaskEvent_PreemptedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f5c8147aadbdf8a, []int{0, 2}
}

func (m *TaskEvent_PreemptedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_PreemptedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_PreemptedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_PreemptedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_PreemptedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_PreemptedDetails.Merge(m, src)
}
func (m *TaskEvent_PreemptedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_PreemptedDetails.Size(m)
}
func (m *TaskEvent_PreemptedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_PreemptedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_PreemptedDetails proto.InternalMessageInfo

func (m *TaskEvent_PreemptedDetails) GetPreemptingAccountId() string {
	if m != nil {
		return m.PreemptingAccountId
	}
	return ""
}

func (m *TaskEvent_PreemptedDetails) GetPreemptingTaskId() string {
	if m != nil {
		return m.PreemptingTaskId
	}
	return ""
}

func (m *TaskEvent_PreemptedDetails) GetPriority() int32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *TaskEvent_PreemptedDetails) GetPreemptingPriority() int32 {
	if m != nil {
		return m.PreemptingPriority
	}
	return 0
}

// ReprioritizedDetails represents event details that are specific to the
// PREPRIORITIZED event type.
type TaskEvent_ReprioritizedDetails struct {
	// OldPriority is the previous priority the task was running at.
	OldPriority int32 `protobuf:"varint,1,opt,name=old_priority,json=oldPriority,proto3" json:"old_priority,omitempty"`
	// NewPriority is the new priority the task is running at.
	NewPriority          int32    `protobuf:"varint,2,opt,name=new_priority,json=newPriority,proto3" json:"new_priority,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent_ReprioritizedDetails) Reset()         { *m = TaskEvent_ReprioritizedDetails{} }
func (m *TaskEvent_ReprioritizedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_ReprioritizedDetails) ProtoMessage()    {}
func (*TaskEvent_ReprioritizedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f5c8147aadbdf8a, []int{0, 3}
}

func (m *TaskEvent_ReprioritizedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_ReprioritizedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_ReprioritizedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_ReprioritizedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_ReprioritizedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_ReprioritizedDetails.Merge(m, src)
}
func (m *TaskEvent_ReprioritizedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_ReprioritizedDetails.Size(m)
}
func (m *TaskEvent_ReprioritizedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_ReprioritizedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_ReprioritizedDetails proto.InternalMessageInfo

func (m *TaskEvent_ReprioritizedDetails) GetOldPriority() int32 {
	if m != nil {
		return m.OldPriority
	}
	return 0
}

func (m *TaskEvent_ReprioritizedDetails) GetNewPriority() int32 {
	if m != nil {
		return m.NewPriority
	}
	return 0
}

// ReprioritizedDetails represents event details that are specific to the
// COMPLETED event type.
type TaskEvent_CompletedDetails struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent_CompletedDetails) Reset()         { *m = TaskEvent_CompletedDetails{} }
func (m *TaskEvent_CompletedDetails) String() string { return proto.CompactTextString(m) }
func (*TaskEvent_CompletedDetails) ProtoMessage()    {}
func (*TaskEvent_CompletedDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f5c8147aadbdf8a, []int{0, 4}
}

func (m *TaskEvent_CompletedDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskEvent_CompletedDetails.Unmarshal(m, b)
}
func (m *TaskEvent_CompletedDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskEvent_CompletedDetails.Marshal(b, m, deterministic)
}
func (m *TaskEvent_CompletedDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent_CompletedDetails.Merge(m, src)
}
func (m *TaskEvent_CompletedDetails) XXX_Size() int {
	return xxx_messageInfo_TaskEvent_CompletedDetails.Size(m)
}
func (m *TaskEvent_CompletedDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent_CompletedDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent_CompletedDetails proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("metrics.TaskEvent_EventType", TaskEvent_EventType_name, TaskEvent_EventType_value)
	proto.RegisterType((*TaskEvent)(nil), "metrics.TaskEvent")
	proto.RegisterType((*TaskEvent_EnqueuedDetails)(nil), "metrics.TaskEvent.EnqueuedDetails")
	proto.RegisterType((*TaskEvent_AssignedDetails)(nil), "metrics.TaskEvent.AssignedDetails")
	proto.RegisterType((*TaskEvent_PreemptedDetails)(nil), "metrics.TaskEvent.PreemptedDetails")
	proto.RegisterType((*TaskEvent_ReprioritizedDetails)(nil), "metrics.TaskEvent.ReprioritizedDetails")
	proto.RegisterType((*TaskEvent_CompletedDetails)(nil), "metrics.TaskEvent.CompletedDetails")
}

func init() {
	proto.RegisterFile("infra/qscheduler/qslib/metrics/metrics.proto", fileDescriptor_3f5c8147aadbdf8a)
}

var fileDescriptor_3f5c8147aadbdf8a = []byte{
	// 756 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0xdd, 0x4e, 0xdb, 0x30,
	0x18, 0x25, 0xf4, 0x37, 0x5f, 0x5a, 0x9a, 0x1a, 0xd0, 0xa2, 0x6a, 0x1b, 0x1d, 0x68, 0xa2, 0x17,
	0xac, 0xd5, 0xd8, 0xe5, 0xae, 0xa0, 0xcd, 0xb6, 0x4c, 0xfc, 0x74, 0xa6, 0xec, 0x62, 0x9a, 0x16,
	0xe5, 0xc7, 0x14, 0x8b, 0x34, 0x4e, 0x93, 0x14, 0xc4, 0xde, 0x67, 0x0f, 0xb2, 0x17, 0xd9, 0xb3,
	0x4c, 0x76, 0x93, 0x34, 0x44, 0xe5, 0xa6, 0xb5, 0xcf, 0x39, 0x3e, 0x3e, 0x76, 0x3e, 0x7f, 0x70,
	0x44, 0xfd, 0x9b, 0xd0, 0x1a, 0xcc, 0x23, 0xe7, 0x96, 0xb8, 0x0b, 0x8f, 0x84, 0x83, 0x79, 0xe4,
	0x51, 0x7b, 0x30, 0x23, 0x71, 0x48, 0x9d, 0x28, 0xfd, 0xef, 0x07, 0x21, 0x8b, 0x19, 0xaa, 0x25,
	0xd3, 0xce, 0xde, 0x94, 0xb1, 0xa9, 0x47, 0x06, 0x02, 0xb6, 0x17, 0x37, 0x83, 0x98, 0xce, 0x48,
	0x14, 0x5b, 0xb3, 0x60, 0xa9, 0xdc, 0xff, 0xa7, 0x80, 0x3c, 0xb1, 0xa2, 0x3b, 0xfd, 0x9e, 0xf8,
	0x31, 0xfa, 0x08, 0x40, 0xf8, 0xc0, 0x8c, 0x1f, 0x03, 0xa2, 0x49, 0x5d, 0xa9, 0xb7, 0x75, 0xfc,
	0xb2, 0x9f, 0x7a, 0x67, 0xba, 0xbe, 0xf8, 0x9d, 0x3c, 0x06, 0x04, 0xcb, 0x24, 0x1d, 0xa2, 0x37,
	0xd0, 0xc8, 0xd2, 0x99, 0xd4, 0xd5, 0x36, 0xbb, 0x52, 0x4f, 0xc6, 0x4a, 0x86, 0x19, 0x2e, 0x7a,
	0x01, 0xb5, 0xd8, 0x8a, 0xee, 0x38, 0x5b, 0x12, 0x6c, 0x95, 0x4f, 0x0d, 0x17, 0xf5, 0xa1, 0xcc,
	0x93, 0x69, 0xe5, 0xae, 0xd4, 0x53, 0x8e, 0x3b, 0xfd, 0x65, 0xec, 0x7e, 0x1a, 0xbb, 0x3f, 0x49,
	0x63, 0x63, 0xa1, 0x43, 0x7b, 0xa0, 0xd8, 0x56, 0x44, 0x4c, 0xcf, 0xb2, 0x89, 0x17, 0x69, 0x95,
	0x6e, 0xa9, 0x27, 0x63, 0xe0, 0xd0, 0x99, 0x40, 0xd0, 0x7b, 0xd8, 0x09, 0x42, 0x76, 0x4f, 0x23,
	0xca, 0x7c, 0xcb, 0xf6, 0x32, 0x65, 0x55, 0x28, 0xb7, 0x9f, 0x70, 0xc9, 0x92, 0x57, 0x00, 0x96,
	0xe3, 0xb0, 0x85, 0x1f, 0xf3, 0x7c, 0x35, 0x91, 0x4f, 0x4e, 0x10, 0xc3, 0x45, 0x07, 0xd0, 0x4c,
	0xe9, 0x7b, 0xcb, 0xa3, 0xae, 0x56, 0xef, 0x4a, 0xbd, 0x3a, 0x6e, 0x24, 0xe0, 0x77, 0x8e, 0xa1,
	0x43, 0x68, 0xa5, 0x22, 0xdb, 0xf2, 0x2c, 0xdf, 0x21, 0x9a, 0xdc, 0x2d, 0xf5, 0x24, 0xbc, 0x95,
	0xc0, 0xa7, 0x4b, 0x14, 0x21, 0x28, 0x3b, 0x2c, 0x8a, 0x35, 0x10, 0xac, 0x18, 0xa3, 0x5d, 0xa8,
	0xda, 0x4c, 0x6c, 0xae, 0x88, 0xcd, 0x2b, 0x36, 0xe3, 0x1b, 0xbf, 0x85, 0x2d, 0x0e, 0xbb, 0x74,
	0x46, 0x7c, 0x1e, 0x39, 0xd2, 0x1a, 0xe2, 0x10, 0x4d, 0x9b, 0xc5, 0xa3, 0x0c, 0x44, 0x97, 0xa0,
	0x12, 0x7f, 0xbe, 0x20, 0x0b, 0xe2, 0x9a, 0x2e, 0x89, 0x2d, 0xea, 0x45, 0x9a, 0x2b, 0xae, 0x73,
	0x7f, 0xdd, 0x17, 0x4c, 0xa4, 0xa3, 0xa5, 0xf2, 0xcb, 0x06, 0x6e, 0x91, 0xa7, 0x10, 0x37, 0xb4,
	0xa2, 0x88, 0x4e, 0xfd, 0x9c, 0x21, 0x79, 0xd6, 0xf0, 0x24, 0x91, 0xe6, 0x0c, 0xad, 0xa7, 0x10,
	0xc2, 0xd0, 0x0e, 0x42, 0x42, 0x66, 0x41, 0x9c, 0x73, 0xbc, 0x11, 0x8e, 0x07, 0x6b, 0x1c, 0xc7,
	0xa9, 0x76, 0x65, 0xa9, 0x06, 0x05, 0x0c, 0xfd, 0x82, 0xdd, 0x90, 0x04, 0x21, 0x65, 0x21, 0x8d,
	0xe9, 0xef, 0x9c, 0xef, 0x54, 0xf8, 0x1e, 0xae, 0xf1, 0xc5, 0x79, 0xfd, 0xca, 0x7b, 0x27, 0x5c,
	0x83, 0xf3, 0xcc, 0x0e, 0x9b, 0x05, 0x1e, 0xc9, 0x67, 0xbe, 0x7d, 0x36, 0xf3, 0x30, 0xd5, 0xe6,
	0x32, 0x3b, 0x05, 0xac, 0xd3, 0x86, 0x56, 0xe1, 0xfa, 0x3b, 0x7f, 0x24, 0x68, 0x15, 0x6e, 0x10,
	0xbd, 0x03, 0x94, 0x95, 0xa9, 0x19, 0x92, 0xf9, 0x82, 0x86, 0xc4, 0x15, 0x8f, 0xb2, 0x8e, 0xdb,
	0x19, 0x83, 0x13, 0x02, 0x75, 0xa0, 0x9e, 0xe4, 0x7f, 0x14, 0x4f, 0xaf, 0x82, 0xb3, 0x39, 0x7a,
	0x0d, 0x90, 0xdc, 0x1c, 0xf5, 0xa7, 0xe2, 0xe9, 0xd5, 0x71, 0x0e, 0xe1, 0x65, 0x9b, 0xce, 0x98,
	0x6f, 0x8a, 0xc2, 0x2c, 0x2f, 0xcb, 0x76, 0x05, 0x0f, 0x59, 0x14, 0x77, 0xfe, 0x4a, 0xa0, 0x16,
	0xbf, 0x0b, 0x3a, 0x86, 0xdd, 0x95, 0x97, 0x99, 0x7b, 0x43, 0x92, 0x28, 0xe3, 0xed, 0x15, 0x79,
	0x92, 0xbd, 0xa6, 0x23, 0x7e, 0xb8, 0x6c, 0x4d, 0xda, 0x14, 0x96, 0x2d, 0x43, 0x5d, 0x31, 0x93,
	0x65, 0x7b, 0xc8, 0x9f, 0xad, 0x54, 0x38, 0xdb, 0x00, 0x72, 0x1b, 0x98, 0x99, 0xac, 0x2c, 0x64,
	0xb9, 0x4d, 0xc6, 0x09, 0xd3, 0xf9, 0x09, 0x3b, 0xeb, 0x4a, 0x80, 0xf7, 0x2f, 0xe6, 0xb9, 0x2b,
	0x07, 0x49, 0x38, 0x28, 0xcc, 0x73, 0xd3, 0xa5, 0x5c, 0xe2, 0x93, 0x07, 0xb3, 0x70, 0xcf, 0x8a,
	0x4f, 0x1e, 0x32, 0x77, 0x04, 0x6a, 0xb1, 0x08, 0xf6, 0x29, 0xc8, 0x59, 0xc7, 0x44, 0x2d, 0x50,
	0xae, 0x2f, 0xae, 0xc6, 0xfa, 0xd0, 0xf8, 0x64, 0xe8, 0x23, 0x75, 0x03, 0x35, 0xa0, 0xae, 0x5f,
	0x7c, 0xbb, 0xd6, 0xaf, 0xf5, 0x91, 0x2a, 0xf1, 0xd9, 0xc9, 0xd5, 0x95, 0xf1, 0xf9, 0x42, 0x1f,
	0xa9, 0x9b, 0xa8, 0x09, 0xf2, 0x18, 0xeb, 0xfa, 0xf9, 0x78, 0xa2, 0x8f, 0xd4, 0x12, 0x6a, 0x43,
	0x13, 0xeb, 0x63, 0x6c, 0x5c, 0x62, 0x63, 0x62, 0xfc, 0xd0, 0x47, 0x6a, 0x99, 0x2b, 0x86, 0x97,
	0xe7, 0xe3, 0x33, 0x9d, 0x2b, 0x2a, 0xa7, 0x32, 0xd4, 0x92, 0x2a, 0xfd, 0x5a, 0xae, 0x37, 0x55,
	0xd7, 0xae, 0x8a, 0x1e, 0xfa, 0xe1, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xbc, 0xe1, 0x64, 0xd3,
	0x41, 0x06, 0x00, 0x00,
}
