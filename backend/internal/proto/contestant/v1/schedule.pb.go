// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: contestant/v1/schedule.proto

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

type Phases int32

const (
	Phases_PHASES_UNSPECIFIED Phases = 0
	Phases_PHASES_QUALIFYING  Phases = 1
	Phases_PHASES_FINALS      Phases = 2
)

// Enum value maps for Phases.
var (
	Phases_name = map[int32]string{
		0: "PHASES_UNSPECIFIED",
		1: "PHASES_QUALIFYING",
		2: "PHASES_FINALS",
	}
	Phases_value = map[string]int32{
		"PHASES_UNSPECIFIED": 0,
		"PHASES_QUALIFYING":  1,
		"PHASES_FINALS":      2,
	}
)

func (x Phases) Enum() *Phases {
	p := new(Phases)
	*p = x
	return p
}

func (x Phases) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Phases) Descriptor() protoreflect.EnumDescriptor {
	return file_contestant_v1_schedule_proto_enumTypes[0].Descriptor()
}

func (Phases) Type() protoreflect.EnumType {
	return &file_contestant_v1_schedule_proto_enumTypes[0]
}

func (x Phases) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Phases.Descriptor instead.
func (Phases) EnumDescriptor() ([]byte, []int) {
	return file_contestant_v1_schedule_proto_rawDescGZIP(), []int{0}
}

type Schedule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrentPhase Phases                 `protobuf:"varint,1,opt,name=current_phase,json=currentPhase,proto3,enum=contestant.v1.Phases" json:"current_phase,omitempty"`
	EndAt        *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end_at,json=endAt,proto3" json:"end_at,omitempty"`
	NextPhase    Phases                 `protobuf:"varint,3,opt,name=next_phase,json=nextPhase,proto3,enum=contestant.v1.Phases" json:"next_phase,omitempty"`
}

func (x *Schedule) Reset() {
	*x = Schedule{}
	mi := &file_contestant_v1_schedule_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Schedule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Schedule) ProtoMessage() {}

func (x *Schedule) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_schedule_proto_msgTypes[0]
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
	return file_contestant_v1_schedule_proto_rawDescGZIP(), []int{0}
}

func (x *Schedule) GetCurrentPhase() Phases {
	if x != nil {
		return x.CurrentPhase
	}
	return Phases_PHASES_UNSPECIFIED
}

func (x *Schedule) GetEndAt() *timestamppb.Timestamp {
	if x != nil {
		return x.EndAt
	}
	return nil
}

func (x *Schedule) GetNextPhase() Phases {
	if x != nil {
		return x.NextPhase
	}
	return Phases_PHASES_UNSPECIFIED
}

type GetScheduleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetScheduleRequest) Reset() {
	*x = GetScheduleRequest{}
	mi := &file_contestant_v1_schedule_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScheduleRequest) ProtoMessage() {}

func (x *GetScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_schedule_proto_msgTypes[1]
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
	return file_contestant_v1_schedule_proto_rawDescGZIP(), []int{1}
}

type GetScheduleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Schedules []*Schedule `protobuf:"bytes,1,rep,name=schedules,proto3" json:"schedules,omitempty"`
}

func (x *GetScheduleResponse) Reset() {
	*x = GetScheduleResponse{}
	mi := &file_contestant_v1_schedule_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetScheduleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScheduleResponse) ProtoMessage() {}

func (x *GetScheduleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_schedule_proto_msgTypes[2]
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
	return file_contestant_v1_schedule_proto_rawDescGZIP(), []int{2}
}

func (x *GetScheduleResponse) GetSchedules() []*Schedule {
	if x != nil {
		return x.Schedules
	}
	return nil
}

var File_contestant_v1_schedule_proto protoreflect.FileDescriptor

var file_contestant_v1_schedule_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62,
	0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcb, 0x01, 0x0a, 0x08,
	0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x44, 0x0a, 0x0d, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x68, 0x61, 0x73, 0x65, 0x73, 0x42, 0x08, 0xba, 0x48, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01,
	0x52, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12, 0x39,
	0x0a, 0x06, 0x65, 0x6e, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8,
	0x01, 0x01, 0x52, 0x05, 0x65, 0x6e, 0x64, 0x41, 0x74, 0x12, 0x3e, 0x0a, 0x0a, 0x6e, 0x65, 0x78,
	0x74, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x68,
	0x61, 0x73, 0x65, 0x73, 0x42, 0x08, 0xba, 0x48, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x09,
	0x6e, 0x65, 0x78, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x47, 0x65, 0x74,
	0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x54, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65,
	0x64, 0x75, 0x6c, 0x65, 0x73, 0x2a, 0x4a, 0x0a, 0x06, 0x50, 0x68, 0x61, 0x73, 0x65, 0x73, 0x12,
	0x16, 0x0a, 0x12, 0x50, 0x48, 0x41, 0x53, 0x45, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x50, 0x48, 0x41, 0x53, 0x45,
	0x53, 0x5f, 0x51, 0x55, 0x41, 0x4c, 0x49, 0x46, 0x59, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x11,
	0x0a, 0x0d, 0x50, 0x48, 0x41, 0x53, 0x45, 0x53, 0x5f, 0x46, 0x49, 0x4e, 0x41, 0x4c, 0x53, 0x10,
	0x02, 0x32, 0x67, 0x0a, 0x0f, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x12, 0x21, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x46, 0x5a, 0x44, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2f, 0x69,
	0x63, 0x74, 0x73, 0x63, 0x2d, 0x6f, 0x75, 0x74, 0x6c, 0x61, 0x6e, 0x64, 0x73, 0x2f, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_contestant_v1_schedule_proto_rawDescOnce sync.Once
	file_contestant_v1_schedule_proto_rawDescData = file_contestant_v1_schedule_proto_rawDesc
)

func file_contestant_v1_schedule_proto_rawDescGZIP() []byte {
	file_contestant_v1_schedule_proto_rawDescOnce.Do(func() {
		file_contestant_v1_schedule_proto_rawDescData = protoimpl.X.CompressGZIP(file_contestant_v1_schedule_proto_rawDescData)
	})
	return file_contestant_v1_schedule_proto_rawDescData
}

var file_contestant_v1_schedule_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_contestant_v1_schedule_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_contestant_v1_schedule_proto_goTypes = []any{
	(Phases)(0),                   // 0: contestant.v1.Phases
	(*Schedule)(nil),              // 1: contestant.v1.Schedule
	(*GetScheduleRequest)(nil),    // 2: contestant.v1.GetScheduleRequest
	(*GetScheduleResponse)(nil),   // 3: contestant.v1.GetScheduleResponse
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_contestant_v1_schedule_proto_depIdxs = []int32{
	0, // 0: contestant.v1.Schedule.current_phase:type_name -> contestant.v1.Phases
	4, // 1: contestant.v1.Schedule.end_at:type_name -> google.protobuf.Timestamp
	0, // 2: contestant.v1.Schedule.next_phase:type_name -> contestant.v1.Phases
	1, // 3: contestant.v1.GetScheduleResponse.schedules:type_name -> contestant.v1.Schedule
	2, // 4: contestant.v1.ScheduleService.GetSchedule:input_type -> contestant.v1.GetScheduleRequest
	3, // 5: contestant.v1.ScheduleService.GetSchedule:output_type -> contestant.v1.GetScheduleResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_contestant_v1_schedule_proto_init() }
func file_contestant_v1_schedule_proto_init() {
	if File_contestant_v1_schedule_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_contestant_v1_schedule_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contestant_v1_schedule_proto_goTypes,
		DependencyIndexes: file_contestant_v1_schedule_proto_depIdxs,
		EnumInfos:         file_contestant_v1_schedule_proto_enumTypes,
		MessageInfos:      file_contestant_v1_schedule_proto_msgTypes,
	}.Build()
	File_contestant_v1_schedule_proto = out.File
	file_contestant_v1_schedule_proto_rawDesc = nil
	file_contestant_v1_schedule_proto_goTypes = nil
	file_contestant_v1_schedule_proto_depIdxs = nil
}
