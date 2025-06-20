// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: admin/v1/schedule.proto

package adminv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Phase int32

const (
	Phase_PHASE_UNSPECIFIED    Phase = 0
	Phase_PHASE_OUT_OF_CONTEST Phase = 1
	Phase_PHASE_IN_CONTEST     Phase = 2
	Phase_PHASE_BREAK          Phase = 3
	Phase_PHASE_AFTER_CONTEST  Phase = 4
)

// Enum value maps for Phase.
var (
	Phase_name = map[int32]string{
		0: "PHASE_UNSPECIFIED",
		1: "PHASE_OUT_OF_CONTEST",
		2: "PHASE_IN_CONTEST",
		3: "PHASE_BREAK",
		4: "PHASE_AFTER_CONTEST",
	}
	Phase_value = map[string]int32{
		"PHASE_UNSPECIFIED":    0,
		"PHASE_OUT_OF_CONTEST": 1,
		"PHASE_IN_CONTEST":     2,
		"PHASE_BREAK":          3,
		"PHASE_AFTER_CONTEST":  4,
	}
)

func (x Phase) Enum() *Phase {
	p := new(Phase)
	*p = x
	return p
}

func (x Phase) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Phase) Descriptor() protoreflect.EnumDescriptor {
	return file_admin_v1_schedule_proto_enumTypes[0].Descriptor()
}

func (Phase) Type() protoreflect.EnumType {
	return &file_admin_v1_schedule_proto_enumTypes[0]
}

func (x Phase) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Phase.Descriptor instead.
func (Phase) EnumDescriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{0}
}

type Schedule struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Phase         Phase                  `protobuf:"varint,1,opt,name=phase,proto3,enum=admin.v1.Phase" json:"phase,omitempty"`
	StartAt       *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=start_at,json=startAt,proto3" json:"start_at,omitempty"`
	EndAt         *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=end_at,json=endAt,proto3" json:"end_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Schedule) Reset() {
	*x = Schedule{}
	mi := &file_admin_v1_schedule_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Schedule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Schedule) ProtoMessage() {}

func (x *Schedule) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[0]
	if x != nil {
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

func (x *Schedule) GetPhase() Phase {
	if x != nil {
		return x.Phase
	}
	return Phase_PHASE_UNSPECIFIED
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

type GetScheduleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetScheduleRequest) Reset() {
	*x = GetScheduleRequest{}
	mi := &file_admin_v1_schedule_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScheduleRequest) ProtoMessage() {}

func (x *GetScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetScheduleRequest.ProtoReflect.Descriptor instead.
func (*GetScheduleRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{1}
}

type GetScheduleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Schedule      []*Schedule            `protobuf:"bytes,1,rep,name=schedule,proto3" json:"schedule,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetScheduleResponse) Reset() {
	*x = GetScheduleResponse{}
	mi := &file_admin_v1_schedule_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetScheduleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScheduleResponse) ProtoMessage() {}

func (x *GetScheduleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetScheduleResponse.ProtoReflect.Descriptor instead.
func (*GetScheduleResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{2}
}

func (x *GetScheduleResponse) GetSchedule() []*Schedule {
	if x != nil {
		return x.Schedule
	}
	return nil
}

type UpdateScheduleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Schedule      []*Schedule            `protobuf:"bytes,1,rep,name=schedule,proto3" json:"schedule,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateScheduleRequest) Reset() {
	*x = UpdateScheduleRequest{}
	mi := &file_admin_v1_schedule_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateScheduleRequest) ProtoMessage() {}

func (x *UpdateScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateScheduleRequest.ProtoReflect.Descriptor instead.
func (*UpdateScheduleRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateScheduleRequest) GetSchedule() []*Schedule {
	if x != nil {
		return x.Schedule
	}
	return nil
}

type UpdateScheduleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateScheduleResponse) Reset() {
	*x = UpdateScheduleResponse{}
	mi := &file_admin_v1_schedule_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateScheduleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateScheduleResponse) ProtoMessage() {}

func (x *UpdateScheduleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_schedule_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateScheduleResponse.ProtoReflect.Descriptor instead.
func (*UpdateScheduleResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_schedule_proto_rawDescGZIP(), []int{4}
}

var File_admin_v1_schedule_proto protoreflect.FileDescriptor

const file_admin_v1_schedule_proto_rawDesc = "" +
	"\n" +
	"\x17admin/v1/schedule.proto\x12\badmin.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"\x9b\x01\n" +
	"\bSchedule\x12%\n" +
	"\x05phase\x18\x01 \x01(\x0e2\x0f.admin.v1.PhaseR\x05phase\x125\n" +
	"\bstart_at\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\astartAt\x121\n" +
	"\x06end_at\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\x05endAt\"\x14\n" +
	"\x12GetScheduleRequest\"E\n" +
	"\x13GetScheduleResponse\x12.\n" +
	"\bschedule\x18\x01 \x03(\v2\x12.admin.v1.ScheduleR\bschedule\"G\n" +
	"\x15UpdateScheduleRequest\x12.\n" +
	"\bschedule\x18\x01 \x03(\v2\x12.admin.v1.ScheduleR\bschedule\"\x18\n" +
	"\x16UpdateScheduleResponse*x\n" +
	"\x05Phase\x12\x15\n" +
	"\x11PHASE_UNSPECIFIED\x10\x00\x12\x18\n" +
	"\x14PHASE_OUT_OF_CONTEST\x10\x01\x12\x14\n" +
	"\x10PHASE_IN_CONTEST\x10\x02\x12\x0f\n" +
	"\vPHASE_BREAK\x10\x03\x12\x17\n" +
	"\x13PHASE_AFTER_CONTEST\x10\x042\xb2\x01\n" +
	"\x0fScheduleService\x12J\n" +
	"\vGetSchedule\x12\x1c.admin.v1.GetScheduleRequest\x1a\x1d.admin.v1.GetScheduleResponse\x12S\n" +
	"\x0eUpdateSchedule\x12\x1f.admin.v1.UpdateScheduleRequest\x1a .admin.v1.UpdateScheduleResponseBCZAgithub.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1;adminv1b\x06proto3"

var (
	file_admin_v1_schedule_proto_rawDescOnce sync.Once
	file_admin_v1_schedule_proto_rawDescData []byte
)

func file_admin_v1_schedule_proto_rawDescGZIP() []byte {
	file_admin_v1_schedule_proto_rawDescOnce.Do(func() {
		file_admin_v1_schedule_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_admin_v1_schedule_proto_rawDesc), len(file_admin_v1_schedule_proto_rawDesc)))
	})
	return file_admin_v1_schedule_proto_rawDescData
}

var file_admin_v1_schedule_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_admin_v1_schedule_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_admin_v1_schedule_proto_goTypes = []any{
	(Phase)(0),                     // 0: admin.v1.Phase
	(*Schedule)(nil),               // 1: admin.v1.Schedule
	(*GetScheduleRequest)(nil),     // 2: admin.v1.GetScheduleRequest
	(*GetScheduleResponse)(nil),    // 3: admin.v1.GetScheduleResponse
	(*UpdateScheduleRequest)(nil),  // 4: admin.v1.UpdateScheduleRequest
	(*UpdateScheduleResponse)(nil), // 5: admin.v1.UpdateScheduleResponse
	(*timestamppb.Timestamp)(nil),  // 6: google.protobuf.Timestamp
}
var file_admin_v1_schedule_proto_depIdxs = []int32{
	0, // 0: admin.v1.Schedule.phase:type_name -> admin.v1.Phase
	6, // 1: admin.v1.Schedule.start_at:type_name -> google.protobuf.Timestamp
	6, // 2: admin.v1.Schedule.end_at:type_name -> google.protobuf.Timestamp
	1, // 3: admin.v1.GetScheduleResponse.schedule:type_name -> admin.v1.Schedule
	1, // 4: admin.v1.UpdateScheduleRequest.schedule:type_name -> admin.v1.Schedule
	2, // 5: admin.v1.ScheduleService.GetSchedule:input_type -> admin.v1.GetScheduleRequest
	4, // 6: admin.v1.ScheduleService.UpdateSchedule:input_type -> admin.v1.UpdateScheduleRequest
	3, // 7: admin.v1.ScheduleService.GetSchedule:output_type -> admin.v1.GetScheduleResponse
	5, // 8: admin.v1.ScheduleService.UpdateSchedule:output_type -> admin.v1.UpdateScheduleResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_admin_v1_schedule_proto_init() }
func file_admin_v1_schedule_proto_init() {
	if File_admin_v1_schedule_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_admin_v1_schedule_proto_rawDesc), len(file_admin_v1_schedule_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_schedule_proto_goTypes,
		DependencyIndexes: file_admin_v1_schedule_proto_depIdxs,
		EnumInfos:         file_admin_v1_schedule_proto_enumTypes,
		MessageInfos:      file_admin_v1_schedule_proto_msgTypes,
	}.Build()
	File_admin_v1_schedule_proto = out.File
	file_admin_v1_schedule_proto_goTypes = nil
	file_admin_v1_schedule_proto_depIdxs = nil
}
