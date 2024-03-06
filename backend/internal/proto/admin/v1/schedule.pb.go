// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: admin/v1/schedule.proto

package v1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

type Rounds int32

const (
	Rounds_ROUNDS_UNSPECIFIED Rounds = 0
	Rounds_ROUNDS_QUALIFYING  Rounds = 1
	Rounds_ROUNDS_FINALS      Rounds = 2
)

// Enum value maps for Rounds.
var (
	Rounds_name = map[int32]string{
		0: "ROUNDS_UNSPECIFIED",
		1: "ROUNDS_QUALIFYING",
		2: "ROUNDS_FINALS",
	}
	Rounds_value = map[string]int32{
		"ROUNDS_UNSPECIFIED": 0,
		"ROUNDS_QUALIFYING":  1,
		"ROUNDS_FINALS":      2,
	}
)

func (x Rounds) Enum() *Rounds {
	p := new(Rounds)
	*p = x
	return p
}

func (x Rounds) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Rounds) Descriptor() protoreflect.EnumDescriptor {
	return file_admin_v1_schedule_proto_enumTypes[0].Descriptor()
}

func (Rounds) Type() protoreflect.EnumType {
	return &file_admin_v1_schedule_proto_enumTypes[0]
}

func (x Rounds) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Rounds.Descriptor instead.
func (Rounds) EnumDescriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{0}
}

type Schedule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name    string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Round   Rounds                 `protobuf:"varint,3,opt,name=round,proto3,enum=admin.v1.Rounds" json:"round,omitempty"`
	StartAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=start_at,json=startAt,proto3" json:"start_at,omitempty"`
	EndAt   *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=end_at,json=endAt,proto3" json:"end_at,omitempty"`
}

func (x *Schedule) Reset() {
	*x = Schedule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_schedule_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Schedule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Schedule) ProtoMessage() {}

func (x *Schedule) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Schedule.ProtoReflect.Descriptor instead.
func (*Schedule) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{0}
}

func (x *Schedule) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Schedule) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Schedule) GetRound() Rounds {
	if x != nil {
		return x.Round
	}
	return Rounds_ROUNDS_UNSPECIFIED
}

func (x *Schedule) GetStartAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StartAt
	}
	return nil
}

func (x *Schedule) GetEndAt() *timestamppb.Timestamp {
	if x != nil {
		return x.EndAt
	}
	return nil
}

type GetSchedulesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetSchedulesRequest) Reset() {
	*x = GetSchedulesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_schedule_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSchedulesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSchedulesRequest) ProtoMessage() {}

func (x *GetSchedulesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSchedulesRequest.ProtoReflect.Descriptor instead.
func (*GetSchedulesRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{1}
}

type GetSchedulesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Schedules []*Schedule `protobuf:"bytes,1,rep,name=schedules,proto3" json:"schedules,omitempty"`
}

func (x *GetSchedulesResponse) Reset() {
	*x = GetSchedulesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_schedule_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSchedulesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSchedulesResponse) ProtoMessage() {}

func (x *GetSchedulesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSchedulesResponse.ProtoReflect.Descriptor instead.
func (*GetSchedulesResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{2}
}

func (x *GetSchedulesResponse) GetSchedules() []*Schedule {
	if x != nil {
		return x.Schedules
	}
	return nil
}

type PatchScheduleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name    *string                `protobuf:"bytes,2,opt,name=name,proto3,oneof" json:"name,omitempty"`
	Round   *Rounds                `protobuf:"varint,3,opt,name=round,proto3,enum=admin.v1.Rounds,oneof" json:"round,omitempty"`
	StartAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=start_at,json=startAt,proto3,oneof" json:"start_at,omitempty"`
	EndAt   *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=end_at,json=endAt,proto3,oneof" json:"end_at,omitempty"`
}

func (x *PatchScheduleRequest) Reset() {
	*x = PatchScheduleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_schedule_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PatchScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchScheduleRequest) ProtoMessage() {}

func (x *PatchScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchScheduleRequest.ProtoReflect.Descriptor instead.
func (*PatchScheduleRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{3}
}

func (x *PatchScheduleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PatchScheduleRequest) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *PatchScheduleRequest) GetRound() Rounds {
	if x != nil && x.Round != nil {
		return *x.Round
	}
	return Rounds_ROUNDS_UNSPECIFIED
}

func (x *PatchScheduleRequest) GetStartAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StartAt
	}
	return nil
}

func (x *PatchScheduleRequest) GetEndAt() *timestamppb.Timestamp {
	if x != nil {
		return x.EndAt
	}
	return nil
}

type PatchScheduleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Schedule *Schedule `protobuf:"bytes,1,opt,name=schedule,proto3" json:"schedule,omitempty"`
}

func (x *PatchScheduleResponse) Reset() {
	*x = PatchScheduleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_schedule_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PatchScheduleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchScheduleResponse) ProtoMessage() {}

func (x *PatchScheduleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchScheduleResponse.ProtoReflect.Descriptor instead.
func (*PatchScheduleResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{4}
}

func (x *PatchScheduleResponse) GetSchedule() *Schedule {
	if x != nil {
		return x.Schedule
	}
	return nil
}

type PostScheduleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Round   Rounds                 `protobuf:"varint,2,opt,name=round,proto3,enum=admin.v1.Rounds" json:"round,omitempty"`
	StartAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=start_at,json=startAt,proto3" json:"start_at,omitempty"`
	EndAt   *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=end_at,json=endAt,proto3" json:"end_at,omitempty"`
}

func (x *PostScheduleRequest) Reset() {
	*x = PostScheduleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_schedule_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostScheduleRequest) ProtoMessage() {}

func (x *PostScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostScheduleRequest.ProtoReflect.Descriptor instead.
func (*PostScheduleRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{5}
}

func (x *PostScheduleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PostScheduleRequest) GetRound() Rounds {
	if x != nil {
		return x.Round
	}
	return Rounds_ROUNDS_UNSPECIFIED
}

func (x *PostScheduleRequest) GetStartAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StartAt
	}
	return nil
}

func (x *PostScheduleRequest) GetEndAt() *timestamppb.Timestamp {
	if x != nil {
		return x.EndAt
	}
	return nil
}

type PostScheduleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Schedule *Schedule `protobuf:"bytes,1,opt,name=schedule,proto3" json:"schedule,omitempty"`
}

func (x *PostScheduleResponse) Reset() {
	*x = PostScheduleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_schedule_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostScheduleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostScheduleResponse) ProtoMessage() {}

func (x *PostScheduleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostScheduleResponse.ProtoReflect.Descriptor instead.
func (*PostScheduleResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{6}
}

func (x *PostScheduleResponse) GetSchedule() *Schedule {
	if x != nil {
		return x.Schedule
	}
	return nil
}

type DeleteScheduleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteScheduleRequest) Reset() {
	*x = DeleteScheduleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_schedule_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteScheduleRequest) ProtoMessage() {}

func (x *DeleteScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteScheduleRequest.ProtoReflect.Descriptor instead.
func (*DeleteScheduleRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteScheduleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteScheduleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteScheduleResponse) Reset() {
	*x = DeleteScheduleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_schedule_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteScheduleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteScheduleResponse) ProtoMessage() {}

func (x *DeleteScheduleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteScheduleResponse.ProtoReflect.Descriptor instead.
func (*DeleteScheduleResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{8}
}

var File_admin_v1_schedule_proto protoreflect.FileDescriptor

var file_admin_v1_schedule_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xef, 0x01, 0x0a, 0x08, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x18,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72,
	0x03, 0x98, 0x01, 0x1a, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18,
	0x64, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x42, 0x08, 0xba, 0x48, 0x05, 0x82, 0x01, 0x02,
	0x10, 0x01, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x3d, 0x0a, 0x08, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52,
	0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x06, 0x65, 0x6e, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x05, 0x65, 0x6e,
	0x64, 0x41, 0x74, 0x22, 0x15, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x50, 0x0a, 0x14, 0x47, 0x65,
	0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01,
	0x01, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x22, 0xbf, 0x02, 0x0a,
	0x14, 0x50, 0x61, 0x74, 0x63, 0x68, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x98, 0x01, 0x1a, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x22, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba,
	0x48, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18, 0x64, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x35, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x10, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f,
	0x75, 0x6e, 0x64, 0x73, 0x42, 0x08, 0xba, 0x48, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x48, 0x01,
	0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x08, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01,
	0x48, 0x02, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x43,
	0x0a, 0x06, 0x65, 0x6e, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x0b, 0xba, 0x48, 0x08, 0xc8,
	0x01, 0x01, 0xb2, 0x01, 0x02, 0x40, 0x01, 0x48, 0x03, 0x52, 0x05, 0x65, 0x6e, 0x64, 0x41, 0x74,
	0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x08, 0x0a, 0x06,
	0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x5f, 0x61, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x61, 0x74, 0x22, 0x4f,
	0x0a, 0x15, 0x50, 0x61, 0x74, 0x63, 0x68, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x08, 0x73, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x42, 0x06, 0xba,
	0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x08, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x22,
	0xea, 0x01, 0x0a, 0x13, 0x50, 0x6f, 0x73, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18, 0x64,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x42, 0x08, 0xba, 0x48, 0x05, 0x82, 0x01, 0x02, 0x10,
	0x01, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x42, 0x0a, 0x08, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x0b, 0xba, 0x48, 0x08, 0xc8, 0x01, 0x01, 0xb2, 0x01,
	0x02, 0x40, 0x01, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x41, 0x74, 0x12, 0x3e, 0x0a, 0x06,
	0x65, 0x6e, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x0b, 0xba, 0x48, 0x08, 0xc8, 0x01, 0x01,
	0xb2, 0x01, 0x02, 0x40, 0x01, 0x52, 0x05, 0x65, 0x6e, 0x64, 0x41, 0x74, 0x22, 0x4e, 0x0a, 0x14,
	0x50, 0x6f, 0x73, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x08, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8,
	0x01, 0x01, 0x52, 0x08, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x22, 0x31, 0x0a, 0x15,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x98, 0x01, 0x1a, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x18, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x4a, 0x0a, 0x06, 0x52, 0x6f, 0x75,
	0x6e, 0x64, 0x73, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x4f, 0x55, 0x4e, 0x44, 0x53, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x52,
	0x4f, 0x55, 0x4e, 0x44, 0x53, 0x5f, 0x51, 0x55, 0x41, 0x4c, 0x49, 0x46, 0x59, 0x49, 0x4e, 0x47,
	0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x52, 0x4f, 0x55, 0x4e, 0x44, 0x53, 0x5f, 0x46, 0x49, 0x4e,
	0x41, 0x4c, 0x53, 0x10, 0x02, 0x32, 0xd6, 0x02, 0x0a, 0x0f, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x0c, 0x47, 0x65, 0x74,
	0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x1d, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x0d, 0x50, 0x61, 0x74, 0x63,
	0x68, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1e, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x74, 0x63, 0x68, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x74, 0x63, 0x68, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0c, 0x50, 0x6f,
	0x73, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1d, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x53, 0x0a, 0x0e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1f, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x41,
	0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63, 0x74,
	0x73, 0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x6f, 0x75, 0x74, 0x6c, 0x61, 0x6e, 0x64,
	0x73, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_admin_v1_schedule_proto_rawDescOnce sync.Once
	file_admin_v1_schedule_proto_rawDescData = file_admin_v1_schedule_proto_rawDesc
)

func file_admin_v1_schedule_proto_rawDescGZIP() []byte {
	file_admin_v1_schedule_proto_rawDescOnce.Do(func() {
		file_admin_v1_schedule_proto_rawDescData = protoimpl.X.CompressGZIP(file_admin_v1_schedule_proto_rawDescData)
	})
	return file_admin_v1_schedule_proto_rawDescData
}

var file_admin_v1_schedule_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_admin_v1_schedule_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_admin_v1_schedule_proto_goTypes = []interface{}{
	(Rounds)(0),                    // 0: admin.v1.Rounds
	(*Schedule)(nil),               // 1: admin.v1.Schedule
	(*GetSchedulesRequest)(nil),    // 2: admin.v1.GetSchedulesRequest
	(*GetSchedulesResponse)(nil),   // 3: admin.v1.GetSchedulesResponse
	(*PatchScheduleRequest)(nil),   // 4: admin.v1.PatchScheduleRequest
	(*PatchScheduleResponse)(nil),  // 5: admin.v1.PatchScheduleResponse
	(*PostScheduleRequest)(nil),    // 6: admin.v1.PostScheduleRequest
	(*PostScheduleResponse)(nil),   // 7: admin.v1.PostScheduleResponse
	(*DeleteScheduleRequest)(nil),  // 8: admin.v1.DeleteScheduleRequest
	(*DeleteScheduleResponse)(nil), // 9: admin.v1.DeleteScheduleResponse
	(*timestamppb.Timestamp)(nil),  // 10: google.protobuf.Timestamp
}
var file_admin_v1_schedule_proto_depIdxs = []int32{
	0,  // 0: admin.v1.Schedule.round:type_name -> admin.v1.Rounds
	10, // 1: admin.v1.Schedule.start_at:type_name -> google.protobuf.Timestamp
	10, // 2: admin.v1.Schedule.end_at:type_name -> google.protobuf.Timestamp
	1,  // 3: admin.v1.GetSchedulesResponse.schedules:type_name -> admin.v1.Schedule
	0,  // 4: admin.v1.PatchScheduleRequest.round:type_name -> admin.v1.Rounds
	10, // 5: admin.v1.PatchScheduleRequest.start_at:type_name -> google.protobuf.Timestamp
	10, // 6: admin.v1.PatchScheduleRequest.end_at:type_name -> google.protobuf.Timestamp
	1,  // 7: admin.v1.PatchScheduleResponse.schedule:type_name -> admin.v1.Schedule
	0,  // 8: admin.v1.PostScheduleRequest.round:type_name -> admin.v1.Rounds
	10, // 9: admin.v1.PostScheduleRequest.start_at:type_name -> google.protobuf.Timestamp
	10, // 10: admin.v1.PostScheduleRequest.end_at:type_name -> google.protobuf.Timestamp
	1,  // 11: admin.v1.PostScheduleResponse.schedule:type_name -> admin.v1.Schedule
	2,  // 12: admin.v1.ScheduleService.GetSchedules:input_type -> admin.v1.GetSchedulesRequest
	4,  // 13: admin.v1.ScheduleService.PatchSchedule:input_type -> admin.v1.PatchScheduleRequest
	6,  // 14: admin.v1.ScheduleService.PostSchedule:input_type -> admin.v1.PostScheduleRequest
	8,  // 15: admin.v1.ScheduleService.DeleteSchedule:input_type -> admin.v1.DeleteScheduleRequest
	3,  // 16: admin.v1.ScheduleService.GetSchedules:output_type -> admin.v1.GetSchedulesResponse
	5,  // 17: admin.v1.ScheduleService.PatchSchedule:output_type -> admin.v1.PatchScheduleResponse
	7,  // 18: admin.v1.ScheduleService.PostSchedule:output_type -> admin.v1.PostScheduleResponse
	9,  // 19: admin.v1.ScheduleService.DeleteSchedule:output_type -> admin.v1.DeleteScheduleResponse
	16, // [16:20] is the sub-list for method output_type
	12, // [12:16] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_admin_v1_schedule_proto_init() }
func file_admin_v1_schedule_proto_init() {
	if File_admin_v1_schedule_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_admin_v1_schedule_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Schedule); i {
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
		file_admin_v1_schedule_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSchedulesRequest); i {
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
		file_admin_v1_schedule_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSchedulesResponse); i {
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
		file_admin_v1_schedule_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PatchScheduleRequest); i {
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
		file_admin_v1_schedule_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PatchScheduleResponse); i {
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
		file_admin_v1_schedule_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostScheduleRequest); i {
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
		file_admin_v1_schedule_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostScheduleResponse); i {
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
		file_admin_v1_schedule_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteScheduleRequest); i {
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
		file_admin_v1_schedule_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteScheduleResponse); i {
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
	file_admin_v1_schedule_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_admin_v1_schedule_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_schedule_proto_goTypes,
		DependencyIndexes: file_admin_v1_schedule_proto_depIdxs,
		EnumInfos:         file_admin_v1_schedule_proto_enumTypes,
		MessageInfos:      file_admin_v1_schedule_proto_msgTypes,
	}.Build()
	File_admin_v1_schedule_proto = out.File
	file_admin_v1_schedule_proto_rawDesc = nil
	file_admin_v1_schedule_proto_goTypes = nil
	file_admin_v1_schedule_proto_depIdxs = nil
}
