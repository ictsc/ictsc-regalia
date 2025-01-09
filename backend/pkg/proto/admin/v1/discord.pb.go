// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: admin/v1/discord.proto

package adminv1

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

type SyncTeamsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncTeamsRequest) Reset() {
	*x = SyncTeamsRequest{}
	mi := &file_admin_v1_discord_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncTeamsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncTeamsRequest) ProtoMessage() {}

func (x *SyncTeamsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_discord_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncTeamsRequest.ProtoReflect.Descriptor instead.
func (*SyncTeamsRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_discord_proto_rawDescGZIP(), []int{0}
}

type SyncTeamsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncTeamsResponse) Reset() {
	*x = SyncTeamsResponse{}
	mi := &file_admin_v1_discord_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncTeamsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncTeamsResponse) ProtoMessage() {}

func (x *SyncTeamsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_discord_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncTeamsResponse.ProtoReflect.Descriptor instead.
func (*SyncTeamsResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_discord_proto_rawDescGZIP(), []int{1}
}

type SyncUsersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncUsersRequest) Reset() {
	*x = SyncUsersRequest{}
	mi := &file_admin_v1_discord_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncUsersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncUsersRequest) ProtoMessage() {}

func (x *SyncUsersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_discord_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncUsersRequest.ProtoReflect.Descriptor instead.
func (*SyncUsersRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_discord_proto_rawDescGZIP(), []int{2}
}

type SyncUsersResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncUsersResponse) Reset() {
	*x = SyncUsersResponse{}
	mi := &file_admin_v1_discord_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncUsersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncUsersResponse) ProtoMessage() {}

func (x *SyncUsersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_discord_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncUsersResponse.ProtoReflect.Descriptor instead.
func (*SyncUsersResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_discord_proto_rawDescGZIP(), []int{3}
}

var File_admin_v1_discord_proto protoreflect.FileDescriptor

var file_admin_v1_discord_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x76, 0x31, 0x22, 0x12, 0x0a, 0x10, 0x53, 0x79, 0x6e, 0x63, 0x54, 0x65, 0x61, 0x6d, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x13, 0x0a, 0x11, 0x53, 0x79, 0x6e, 0x63, 0x54, 0x65,
	0x61, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x12, 0x0a, 0x10, 0x53,
	0x79, 0x6e, 0x63, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x13, 0x0a, 0x11, 0x53, 0x79, 0x6e, 0x63, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0x9c, 0x01, 0x0a, 0x0e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x44, 0x0a, 0x09, 0x53, 0x79, 0x6e, 0x63, 0x54,
	0x65, 0x61, 0x6d, 0x73, 0x12, 0x1a, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x79, 0x6e, 0x63, 0x54, 0x65, 0x61, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1b, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x79, 0x6e, 0x63,
	0x54, 0x65, 0x61, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a,
	0x09, 0x53, 0x79, 0x6e, 0x63, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x1a, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x72, 0x65,
	0x67, 0x61, 0x6c, 0x69, 0x61, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31,
	0x3b, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_admin_v1_discord_proto_rawDescOnce sync.Once
	file_admin_v1_discord_proto_rawDescData = file_admin_v1_discord_proto_rawDesc
)

func file_admin_v1_discord_proto_rawDescGZIP() []byte {
	file_admin_v1_discord_proto_rawDescOnce.Do(func() {
		file_admin_v1_discord_proto_rawDescData = protoimpl.X.CompressGZIP(file_admin_v1_discord_proto_rawDescData)
	})
	return file_admin_v1_discord_proto_rawDescData
}

var file_admin_v1_discord_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_admin_v1_discord_proto_goTypes = []any{
	(*SyncTeamsRequest)(nil),  // 0: admin.v1.SyncTeamsRequest
	(*SyncTeamsResponse)(nil), // 1: admin.v1.SyncTeamsResponse
	(*SyncUsersRequest)(nil),  // 2: admin.v1.SyncUsersRequest
	(*SyncUsersResponse)(nil), // 3: admin.v1.SyncUsersResponse
}
var file_admin_v1_discord_proto_depIdxs = []int32{
	0, // 0: admin.v1.DiscordService.SyncTeams:input_type -> admin.v1.SyncTeamsRequest
	2, // 1: admin.v1.DiscordService.SyncUsers:input_type -> admin.v1.SyncUsersRequest
	1, // 2: admin.v1.DiscordService.SyncTeams:output_type -> admin.v1.SyncTeamsResponse
	3, // 3: admin.v1.DiscordService.SyncUsers:output_type -> admin.v1.SyncUsersResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_admin_v1_discord_proto_init() }
func file_admin_v1_discord_proto_init() {
	if File_admin_v1_discord_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_admin_v1_discord_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_discord_proto_goTypes,
		DependencyIndexes: file_admin_v1_discord_proto_depIdxs,
		MessageInfos:      file_admin_v1_discord_proto_msgTypes,
	}.Build()
	File_admin_v1_discord_proto = out.File
	file_admin_v1_discord_proto_rawDesc = nil
	file_admin_v1_discord_proto_goTypes = nil
	file_admin_v1_discord_proto_depIdxs = nil
}
