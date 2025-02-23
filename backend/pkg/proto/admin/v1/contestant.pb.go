// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: admin/v1/contestant.proto

package adminv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type Contestant struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	DisplayName   string                 `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	Team          *Team                  `protobuf:"bytes,3,opt,name=team,proto3" json:"team,omitempty"`
	Profile       *Profile               `protobuf:"bytes,4,opt,name=profile,proto3" json:"profile,omitempty"`
	DiscordId     string                 `protobuf:"bytes,5,opt,name=discord_id,json=discordId,proto3" json:"discord_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Contestant) Reset() {
	*x = Contestant{}
	mi := &file_admin_v1_contestant_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Contestant) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Contestant) ProtoMessage() {}

func (x *Contestant) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_contestant_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Contestant.ProtoReflect.Descriptor instead.
func (*Contestant) Descriptor() ([]byte, []int) {
	return file_admin_v1_contestant_proto_rawDescGZIP(), []int{0}
}

func (x *Contestant) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Contestant) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *Contestant) GetTeam() *Team {
	if x != nil {
		return x.Team
	}
	return nil
}

func (x *Contestant) GetProfile() *Profile {
	if x != nil {
		return x.Profile
	}
	return nil
}

func (x *Contestant) GetDiscordId() string {
	if x != nil {
		return x.DiscordId
	}
	return ""
}

type Profile struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	SelfIntroduction string                 `protobuf:"bytes,1,opt,name=self_introduction,json=selfIntroduction,proto3" json:"self_introduction,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *Profile) Reset() {
	*x = Profile{}
	mi := &file_admin_v1_contestant_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Profile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Profile) ProtoMessage() {}

func (x *Profile) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_contestant_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Profile.ProtoReflect.Descriptor instead.
func (*Profile) Descriptor() ([]byte, []int) {
	return file_admin_v1_contestant_proto_rawDescGZIP(), []int{1}
}

func (x *Profile) GetSelfIntroduction() string {
	if x != nil {
		return x.SelfIntroduction
	}
	return ""
}

type ListContestantsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TeamCode      int64                  `protobuf:"varint,1,opt,name=team_code,json=teamCode,proto3" json:"team_code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListContestantsRequest) Reset() {
	*x = ListContestantsRequest{}
	mi := &file_admin_v1_contestant_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListContestantsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListContestantsRequest) ProtoMessage() {}

func (x *ListContestantsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_contestant_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListContestantsRequest.ProtoReflect.Descriptor instead.
func (*ListContestantsRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_contestant_proto_rawDescGZIP(), []int{2}
}

func (x *ListContestantsRequest) GetTeamCode() int64 {
	if x != nil {
		return x.TeamCode
	}
	return 0
}

type ListContestantsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Contestants   []*Contestant          `protobuf:"bytes,1,rep,name=contestants,proto3" json:"contestants,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListContestantsResponse) Reset() {
	*x = ListContestantsResponse{}
	mi := &file_admin_v1_contestant_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListContestantsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListContestantsResponse) ProtoMessage() {}

func (x *ListContestantsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_contestant_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListContestantsResponse.ProtoReflect.Descriptor instead.
func (*ListContestantsResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_contestant_proto_rawDescGZIP(), []int{3}
}

func (x *ListContestantsResponse) GetContestants() []*Contestant {
	if x != nil {
		return x.Contestants
	}
	return nil
}

var File_admin_v1_contestant_proto protoreflect.FileDescriptor

var file_admin_v1_contestant_proto_rawDesc = string([]byte{
	0x0a, 0x19, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f,
	0x74, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb3, 0x01, 0x0a, 0x0a, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x22, 0x0a, 0x04, 0x74, 0x65, 0x61, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x04,
	0x74, 0x65, 0x61, 0x6d, 0x12, 0x2b, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x49, 0x64,
	0x22, 0x36, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x73,
	0x65, 0x6c, 0x66, 0x5f, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x73, 0x65, 0x6c, 0x66, 0x49, 0x6e, 0x74, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x35, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x61, 0x6d, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x74, 0x65, 0x61, 0x6d, 0x43, 0x6f, 0x64, 0x65, 0x22,
	0x51, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x0b, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x73, 0x74, 0x61, 0x6e, 0x74, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x73, 0x32, 0x6b, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x56, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x12, 0x20, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73,
	0x74, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63,
	0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x72, 0x65, 0x67, 0x61, 0x6c, 0x69,
	0x61, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_admin_v1_contestant_proto_rawDescOnce sync.Once
	file_admin_v1_contestant_proto_rawDescData []byte
)

func file_admin_v1_contestant_proto_rawDescGZIP() []byte {
	file_admin_v1_contestant_proto_rawDescOnce.Do(func() {
		file_admin_v1_contestant_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_admin_v1_contestant_proto_rawDesc), len(file_admin_v1_contestant_proto_rawDesc)))
	})
	return file_admin_v1_contestant_proto_rawDescData
}

var file_admin_v1_contestant_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_admin_v1_contestant_proto_goTypes = []any{
	(*Contestant)(nil),              // 0: admin.v1.Contestant
	(*Profile)(nil),                 // 1: admin.v1.Profile
	(*ListContestantsRequest)(nil),  // 2: admin.v1.ListContestantsRequest
	(*ListContestantsResponse)(nil), // 3: admin.v1.ListContestantsResponse
	(*Team)(nil),                    // 4: admin.v1.Team
}
var file_admin_v1_contestant_proto_depIdxs = []int32{
	4, // 0: admin.v1.Contestant.team:type_name -> admin.v1.Team
	1, // 1: admin.v1.Contestant.profile:type_name -> admin.v1.Profile
	0, // 2: admin.v1.ListContestantsResponse.contestants:type_name -> admin.v1.Contestant
	2, // 3: admin.v1.ContestantService.ListContestants:input_type -> admin.v1.ListContestantsRequest
	3, // 4: admin.v1.ContestantService.ListContestants:output_type -> admin.v1.ListContestantsResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_admin_v1_contestant_proto_init() }
func file_admin_v1_contestant_proto_init() {
	if File_admin_v1_contestant_proto != nil {
		return
	}
	file_admin_v1_team_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_admin_v1_contestant_proto_rawDesc), len(file_admin_v1_contestant_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_contestant_proto_goTypes,
		DependencyIndexes: file_admin_v1_contestant_proto_depIdxs,
		MessageInfos:      file_admin_v1_contestant_proto_msgTypes,
	}.Build()
	File_admin_v1_contestant_proto = out.File
	file_admin_v1_contestant_proto_goTypes = nil
	file_admin_v1_contestant_proto_depIdxs = nil
}
