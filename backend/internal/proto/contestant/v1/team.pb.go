// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: contestant/v1/team.proto

package v1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

type Team struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Code         int64  `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Name         string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Organization string `protobuf:"bytes,4,opt,name=organization,proto3" json:"organization,omitempty"`
}

func (x *Team) Reset() {
	*x = Team{}
	mi := &file_contestant_v1_team_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Team) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Team) ProtoMessage() {}

func (x *Team) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_team_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Team.ProtoReflect.Descriptor instead.
func (*Team) Descriptor() ([]byte, []int) {
	return file_contestant_v1_team_proto_rawDescGZIP(), []int{0}
}

func (x *Team) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Team) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Team) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Team) GetOrganization() string {
	if x != nil {
		return x.Organization
	}
	return ""
}

type GetTeamsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetTeamsRequest) Reset() {
	*x = GetTeamsRequest{}
	mi := &file_contestant_v1_team_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTeamsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTeamsRequest) ProtoMessage() {}

func (x *GetTeamsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_team_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTeamsRequest.ProtoReflect.Descriptor instead.
func (*GetTeamsRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_team_proto_rawDescGZIP(), []int{1}
}

type GetTeamsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Teams []*Team `protobuf:"bytes,1,rep,name=teams,proto3" json:"teams,omitempty"`
}

func (x *GetTeamsResponse) Reset() {
	*x = GetTeamsResponse{}
	mi := &file_contestant_v1_team_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTeamsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTeamsResponse) ProtoMessage() {}

func (x *GetTeamsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_team_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTeamsResponse.ProtoReflect.Descriptor instead.
func (*GetTeamsResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_team_proto_rawDescGZIP(), []int{2}
}

func (x *GetTeamsResponse) GetTeams() []*Team {
	if x != nil {
		return x.Teams
	}
	return nil
}

type GetTeamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetTeamRequest) Reset() {
	*x = GetTeamRequest{}
	mi := &file_contestant_v1_team_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTeamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTeamRequest) ProtoMessage() {}

func (x *GetTeamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_team_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTeamRequest.ProtoReflect.Descriptor instead.
func (*GetTeamRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_team_proto_rawDescGZIP(), []int{3}
}

func (x *GetTeamRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetTeamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Team *Team `protobuf:"bytes,1,opt,name=team,proto3" json:"team,omitempty"`
}

func (x *GetTeamResponse) Reset() {
	*x = GetTeamResponse{}
	mi := &file_contestant_v1_team_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTeamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTeamResponse) ProtoMessage() {}

func (x *GetTeamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_team_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTeamResponse.ProtoReflect.Descriptor instead.
func (*GetTeamResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_team_proto_rawDescGZIP(), []int{4}
}

func (x *GetTeamResponse) GetTeam() *Team {
	if x != nil {
		return x.Team
	}
	return nil
}

type GetMembersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetMembersRequest) Reset() {
	*x = GetMembersRequest{}
	mi := &file_contestant_v1_team_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMembersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMembersRequest) ProtoMessage() {}

func (x *GetMembersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_team_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMembersRequest.ProtoReflect.Descriptor instead.
func (*GetMembersRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_team_proto_rawDescGZIP(), []int{5}
}

func (x *GetMembersRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetMembersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Members []*Contestant `protobuf:"bytes,1,rep,name=members,proto3" json:"members,omitempty"`
}

func (x *GetMembersResponse) Reset() {
	*x = GetMembersResponse{}
	mi := &file_contestant_v1_team_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMembersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMembersResponse) ProtoMessage() {}

func (x *GetMembersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_team_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMembersResponse.ProtoReflect.Descriptor instead.
func (*GetMembersResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_team_proto_rawDescGZIP(), []int{6}
}

func (x *GetMembersResponse) GetMembers() []*Contestant {
	if x != nil {
		return x.Members
	}
	return nil
}

type Bastion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User     string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Host     string `protobuf:"bytes,3,opt,name=host,proto3" json:"host,omitempty"`
	Port     int64  `protobuf:"varint,4,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *Bastion) Reset() {
	*x = Bastion{}
	mi := &file_contestant_v1_team_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Bastion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bastion) ProtoMessage() {}

func (x *Bastion) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_team_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bastion.ProtoReflect.Descriptor instead.
func (*Bastion) Descriptor() ([]byte, []int) {
	return file_contestant_v1_team_proto_rawDescGZIP(), []int{7}
}

func (x *Bastion) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *Bastion) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Bastion) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Bastion) GetPort() int64 {
	if x != nil {
		return x.Port
	}
	return 0
}

type GetConnectionInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetConnectionInfoRequest) Reset() {
	*x = GetConnectionInfoRequest{}
	mi := &file_contestant_v1_team_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetConnectionInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConnectionInfoRequest) ProtoMessage() {}

func (x *GetConnectionInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_team_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConnectionInfoRequest.ProtoReflect.Descriptor instead.
func (*GetConnectionInfoRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_team_proto_rawDescGZIP(), []int{8}
}

type GetConnectionInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bastion *Bastion `protobuf:"bytes,1,opt,name=bastion,proto3" json:"bastion,omitempty"`
}

func (x *GetConnectionInfoResponse) Reset() {
	*x = GetConnectionInfoResponse{}
	mi := &file_contestant_v1_team_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetConnectionInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConnectionInfoResponse) ProtoMessage() {}

func (x *GetConnectionInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_team_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConnectionInfoResponse.ProtoReflect.Descriptor instead.
func (*GetConnectionInfoResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_team_proto_rawDescGZIP(), []int{9}
}

func (x *GetConnectionInfoResponse) GetBastion() *Bastion {
	if x != nil {
		return x.Bastion
	}
	return nil
}

var File_contestant_v1_team_proto protoreflect.FileDescriptor

var file_contestant_v1_team_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x74, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61,
	0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8d, 0x01, 0x0a, 0x04, 0x54, 0x65, 0x61, 0x6d, 0x12,
	0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05,
	0x72, 0x03, 0x98, 0x01, 0x1a, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x42, 0x09, 0xba, 0x48, 0x06, 0x22, 0x04, 0x10, 0x64,
	0x20, 0x01, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18,
	0x14, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2d, 0x0a, 0x0c, 0x6f, 0x72, 0x67, 0x61, 0x6e,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba,
	0x48, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18, 0x32, 0x52, 0x0c, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x11, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x54, 0x65, 0x61,
	0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x45, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x54, 0x65, 0x61, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a,
	0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x61,
	0x6d, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73,
	0x22, 0x2a, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xba, 0x48, 0x05, 0x72, 0x03, 0x98, 0x01, 0x1a, 0x52, 0x02, 0x69, 0x64, 0x22, 0x42, 0x0a, 0x0f,
	0x47, 0x65, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2f, 0x0a, 0x04, 0x74, 0x65, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65,
	0x61, 0x6d, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x04, 0x74, 0x65, 0x61, 0x6d,
	0x22, 0x2d, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x98, 0x01, 0x1a, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x55, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x42, 0x0a, 0xba, 0x48, 0x07, 0x92, 0x01, 0x04, 0x08, 0x00, 0x10, 0x05, 0x52, 0x07, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x22, 0x8b, 0x01, 0x0a, 0x07, 0x42, 0x61, 0x73, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18, 0x14, 0x52, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x12, 0x25, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18, 0x14, 0x52, 0x08,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1d, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18,
	0x64, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x42, 0x07, 0xba, 0x48, 0x04, 0x22, 0x02, 0x28, 0x00, 0x52, 0x04,
	0x70, 0x6f, 0x72, 0x74, 0x22, 0x1a, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x4d, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a,
	0x07, 0x62, 0x61, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x42,
	0x61, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x62, 0x61, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x32,
	0xdf, 0x02, 0x0a, 0x0b, 0x54, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x4b, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x73, 0x12, 0x1e, 0x2e, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x54,
	0x65, 0x61, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x54,
	0x65, 0x61, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x07,
	0x47, 0x65, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x12, 0x1d, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73,
	0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x12, 0x20, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x66, 0x0a, 0x11, 0x47, 0x65, 0x74,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x27,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73,
	0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x46, 0x5a, 0x44, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x69, 0x63, 0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x6f, 0x75, 0x74, 0x6c,
	0x61, 0x6e, 0x64, 0x73, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_contestant_v1_team_proto_rawDescOnce sync.Once
	file_contestant_v1_team_proto_rawDescData = file_contestant_v1_team_proto_rawDesc
)

func file_contestant_v1_team_proto_rawDescGZIP() []byte {
	file_contestant_v1_team_proto_rawDescOnce.Do(func() {
		file_contestant_v1_team_proto_rawDescData = protoimpl.X.CompressGZIP(file_contestant_v1_team_proto_rawDescData)
	})
	return file_contestant_v1_team_proto_rawDescData
}

var file_contestant_v1_team_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_contestant_v1_team_proto_goTypes = []any{
	(*Team)(nil),                      // 0: contestant.v1.Team
	(*GetTeamsRequest)(nil),           // 1: contestant.v1.GetTeamsRequest
	(*GetTeamsResponse)(nil),          // 2: contestant.v1.GetTeamsResponse
	(*GetTeamRequest)(nil),            // 3: contestant.v1.GetTeamRequest
	(*GetTeamResponse)(nil),           // 4: contestant.v1.GetTeamResponse
	(*GetMembersRequest)(nil),         // 5: contestant.v1.GetMembersRequest
	(*GetMembersResponse)(nil),        // 6: contestant.v1.GetMembersResponse
	(*Bastion)(nil),                   // 7: contestant.v1.Bastion
	(*GetConnectionInfoRequest)(nil),  // 8: contestant.v1.GetConnectionInfoRequest
	(*GetConnectionInfoResponse)(nil), // 9: contestant.v1.GetConnectionInfoResponse
	(*Contestant)(nil),                // 10: contestant.v1.Contestant
}
var file_contestant_v1_team_proto_depIdxs = []int32{
	0,  // 0: contestant.v1.GetTeamsResponse.teams:type_name -> contestant.v1.Team
	0,  // 1: contestant.v1.GetTeamResponse.team:type_name -> contestant.v1.Team
	10, // 2: contestant.v1.GetMembersResponse.members:type_name -> contestant.v1.Contestant
	7,  // 3: contestant.v1.GetConnectionInfoResponse.bastion:type_name -> contestant.v1.Bastion
	1,  // 4: contestant.v1.TeamService.GetTeams:input_type -> contestant.v1.GetTeamsRequest
	3,  // 5: contestant.v1.TeamService.GetTeam:input_type -> contestant.v1.GetTeamRequest
	5,  // 6: contestant.v1.TeamService.GetMembers:input_type -> contestant.v1.GetMembersRequest
	8,  // 7: contestant.v1.TeamService.GetConnectionInfo:input_type -> contestant.v1.GetConnectionInfoRequest
	2,  // 8: contestant.v1.TeamService.GetTeams:output_type -> contestant.v1.GetTeamsResponse
	4,  // 9: contestant.v1.TeamService.GetTeam:output_type -> contestant.v1.GetTeamResponse
	6,  // 10: contestant.v1.TeamService.GetMembers:output_type -> contestant.v1.GetMembersResponse
	9,  // 11: contestant.v1.TeamService.GetConnectionInfo:output_type -> contestant.v1.GetConnectionInfoResponse
	8,  // [8:12] is the sub-list for method output_type
	4,  // [4:8] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_contestant_v1_team_proto_init() }
func file_contestant_v1_team_proto_init() {
	if File_contestant_v1_team_proto != nil {
		return
	}
	file_contestant_v1_contestant_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_contestant_v1_team_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contestant_v1_team_proto_goTypes,
		DependencyIndexes: file_contestant_v1_team_proto_depIdxs,
		MessageInfos:      file_contestant_v1_team_proto_msgTypes,
	}.Build()
	File_contestant_v1_team_proto = out.File
	file_contestant_v1_team_proto_rawDesc = nil
	file_contestant_v1_team_proto_goTypes = nil
	file_contestant_v1_team_proto_depIdxs = nil
}
