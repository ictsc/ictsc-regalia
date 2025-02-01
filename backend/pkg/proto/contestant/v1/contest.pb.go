// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: contestant/v1/contest.proto

package contestantv1

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
	return file_contestant_v1_contest_proto_enumTypes[0].Descriptor()
}

func (Phase) Type() protoreflect.EnumType {
	return &file_contestant_v1_contest_proto_enumTypes[0]
}

func (x Phase) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Phase.Descriptor instead.
func (Phase) EnumDescriptor() ([]byte, []int) {
	return file_contestant_v1_contest_proto_rawDescGZIP(), []int{0}
}

type Schedule struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Phase         Phase                  `protobuf:"varint,1,opt,name=phase,proto3,enum=contestant.v1.Phase" json:"phase,omitempty"`
	NextPhase     Phase                  `protobuf:"varint,2,opt,name=next_phase,json=nextPhase,proto3,enum=contestant.v1.Phase" json:"next_phase,omitempty"`
	StartAt       *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=start_at,json=startAt,proto3" json:"start_at,omitempty"`
	EndAt         *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=end_at,json=endAt,proto3,oneof" json:"end_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Schedule) Reset() {
	*x = Schedule{}
	mi := &file_contestant_v1_contest_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Schedule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Schedule) ProtoMessage() {}

func (x *Schedule) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_contest_proto_msgTypes[0]
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
	return file_contestant_v1_contest_proto_rawDescGZIP(), []int{0}
}

func (x *Schedule) GetPhase() Phase {
	if x != nil {
		return x.Phase
	}
	return Phase_PHASE_UNSPECIFIED
}

func (x *Schedule) GetNextPhase() Phase {
	if x != nil {
		return x.NextPhase
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
	mi := &file_contestant_v1_contest_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScheduleRequest) ProtoMessage() {}

func (x *GetScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_contest_proto_msgTypes[1]
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
	return file_contestant_v1_contest_proto_rawDescGZIP(), []int{1}
}

type GetScheduleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Schedule      *Schedule              `protobuf:"bytes,1,opt,name=schedule,proto3" json:"schedule,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetScheduleResponse) Reset() {
	*x = GetScheduleResponse{}
	mi := &file_contestant_v1_contest_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetScheduleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScheduleResponse) ProtoMessage() {}

func (x *GetScheduleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_contest_proto_msgTypes[2]
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
	return file_contestant_v1_contest_proto_rawDescGZIP(), []int{2}
}

func (x *GetScheduleResponse) GetSchedule() *Schedule {
	if x != nil {
		return x.Schedule
	}
	return nil
}

var File_contestant_v1_contest_proto protoreflect.FileDescriptor

var file_contestant_v1_contest_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe5, 0x01,
	0x0a, 0x08, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x70, 0x68,
	0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x68, 0x61, 0x73, 0x65, 0x52,
	0x05, 0x70, 0x68, 0x61, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x0a, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70,
	0x68, 0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x68, 0x61, 0x73, 0x65,
	0x52, 0x09, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x41, 0x74, 0x12, 0x36, 0x0a, 0x06, 0x65, 0x6e, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00,
	0x52, 0x05, 0x65, 0x6e, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x65,
	0x6e, 0x64, 0x5f, 0x61, 0x74, 0x22, 0x14, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65,
	0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4a, 0x0a, 0x13, 0x47,
	0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x33, 0x0a, 0x08, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x08, 0x73,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x2a, 0x78, 0x0a, 0x05, 0x50, 0x68, 0x61, 0x73, 0x65,
	0x12, 0x15, 0x0a, 0x11, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x50, 0x48, 0x41, 0x53, 0x45,
	0x5f, 0x4f, 0x55, 0x54, 0x5f, 0x4f, 0x46, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x53, 0x54, 0x10,
	0x01, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x43, 0x4f,
	0x4e, 0x54, 0x45, 0x53, 0x54, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x50, 0x48, 0x41, 0x53, 0x45,
	0x5f, 0x42, 0x52, 0x45, 0x41, 0x4b, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x50, 0x48, 0x41, 0x53,
	0x45, 0x5f, 0x41, 0x46, 0x54, 0x45, 0x52, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x53, 0x54, 0x10,
	0x04, 0x32, 0x66, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x54, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x12, 0x21, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2f, 0x69, 0x63,
	0x74, 0x73, 0x63, 0x2d, 0x72, 0x65, 0x67, 0x61, 0x6c, 0x69, 0x61, 0x2f, 0x62, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_contestant_v1_contest_proto_rawDescOnce sync.Once
	file_contestant_v1_contest_proto_rawDescData []byte
)

func file_contestant_v1_contest_proto_rawDescGZIP() []byte {
	file_contestant_v1_contest_proto_rawDescOnce.Do(func() {
		file_contestant_v1_contest_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_contestant_v1_contest_proto_rawDesc), len(file_contestant_v1_contest_proto_rawDesc)))
	})
	return file_contestant_v1_contest_proto_rawDescData
}

var file_contestant_v1_contest_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_contestant_v1_contest_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_contestant_v1_contest_proto_goTypes = []any{
	(Phase)(0),                    // 0: contestant.v1.Phase
	(*Schedule)(nil),              // 1: contestant.v1.Schedule
	(*GetScheduleRequest)(nil),    // 2: contestant.v1.GetScheduleRequest
	(*GetScheduleResponse)(nil),   // 3: contestant.v1.GetScheduleResponse
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_contestant_v1_contest_proto_depIdxs = []int32{
	0, // 0: contestant.v1.Schedule.phase:type_name -> contestant.v1.Phase
	0, // 1: contestant.v1.Schedule.next_phase:type_name -> contestant.v1.Phase
	4, // 2: contestant.v1.Schedule.start_at:type_name -> google.protobuf.Timestamp
	4, // 3: contestant.v1.Schedule.end_at:type_name -> google.protobuf.Timestamp
	1, // 4: contestant.v1.GetScheduleResponse.schedule:type_name -> contestant.v1.Schedule
	2, // 5: contestant.v1.ContestService.GetSchedule:input_type -> contestant.v1.GetScheduleRequest
	3, // 6: contestant.v1.ContestService.GetSchedule:output_type -> contestant.v1.GetScheduleResponse
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_contestant_v1_contest_proto_init() }
func file_contestant_v1_contest_proto_init() {
	if File_contestant_v1_contest_proto != nil {
		return
	}
	file_contestant_v1_contest_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_contestant_v1_contest_proto_rawDesc), len(file_contestant_v1_contest_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contestant_v1_contest_proto_goTypes,
		DependencyIndexes: file_contestant_v1_contest_proto_depIdxs,
		EnumInfos:         file_contestant_v1_contest_proto_enumTypes,
		MessageInfos:      file_contestant_v1_contest_proto_msgTypes,
	}.Build()
	File_contestant_v1_contest_proto = out.File
	file_contestant_v1_contest_proto_goTypes = nil
	file_contestant_v1_contest_proto_depIdxs = nil
}
