// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: admin/v1/contestant.proto

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

type Contestant struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	TeamId string `protobuf:"bytes,3,opt,name=team_id,json=teamId,proto3" json:"team_id,omitempty"`
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

func (x *Contestant) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Contestant) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Contestant) GetTeamId() string {
	if x != nil {
		return x.TeamId
	}
	return ""
}

type GetContestantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetContestantRequest) Reset() {
	*x = GetContestantRequest{}
	mi := &file_admin_v1_contestant_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetContestantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetContestantRequest) ProtoMessage() {}

func (x *GetContestantRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetContestantRequest.ProtoReflect.Descriptor instead.
func (*GetContestantRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_contestant_proto_rawDescGZIP(), []int{1}
}

func (x *GetContestantRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetContestantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *Contestant `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *GetContestantResponse) Reset() {
	*x = GetContestantResponse{}
	mi := &file_admin_v1_contestant_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetContestantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetContestantResponse) ProtoMessage() {}

func (x *GetContestantResponse) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetContestantResponse.ProtoReflect.Descriptor instead.
func (*GetContestantResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_contestant_proto_rawDescGZIP(), []int{2}
}

func (x *GetContestantResponse) GetUser() *Contestant {
	if x != nil {
		return x.User
	}
	return nil
}

type GetContestantsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetContestantsRequest) Reset() {
	*x = GetContestantsRequest{}
	mi := &file_admin_v1_contestant_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetContestantsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetContestantsRequest) ProtoMessage() {}

func (x *GetContestantsRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetContestantsRequest.ProtoReflect.Descriptor instead.
func (*GetContestantsRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_contestant_proto_rawDescGZIP(), []int{3}
}

type GetContestantsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*Contestant `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *GetContestantsResponse) Reset() {
	*x = GetContestantsResponse{}
	mi := &file_admin_v1_contestant_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetContestantsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetContestantsResponse) ProtoMessage() {}

func (x *GetContestantsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_contestant_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetContestantsResponse.ProtoReflect.Descriptor instead.
func (*GetContestantsResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_contestant_proto_rawDescGZIP(), []int{4}
}

func (x *GetContestantsResponse) GetUsers() []*Contestant {
	if x != nil {
		return x.Users
	}
	return nil
}

type DeleteContestantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteContestantRequest) Reset() {
	*x = DeleteContestantRequest{}
	mi := &file_admin_v1_contestant_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteContestantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteContestantRequest) ProtoMessage() {}

func (x *DeleteContestantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_contestant_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteContestantRequest.ProtoReflect.Descriptor instead.
func (*DeleteContestantRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_contestant_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteContestantRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteContestantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteContestantResponse) Reset() {
	*x = DeleteContestantResponse{}
	mi := &file_admin_v1_contestant_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteContestantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteContestantResponse) ProtoMessage() {}

func (x *DeleteContestantResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_contestant_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteContestantResponse.ProtoReflect.Descriptor instead.
func (*DeleteContestantResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_contestant_proto_rawDescGZIP(), []int{6}
}

var File_admin_v1_contestant_proto protoreflect.FileDescriptor

var file_admin_v1_contestant_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x68, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48,
	0x05, 0x72, 0x03, 0x98, 0x01, 0x1a, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10,
	0x01, 0x18, 0x14, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x07, 0x74, 0x65, 0x61,
	0x6d, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72,
	0x03, 0x98, 0x01, 0x1a, 0x52, 0x06, 0x74, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x22, 0x30, 0x0a, 0x14,
	0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x98, 0x01, 0x1a, 0x52, 0x02, 0x69, 0x64, 0x22, 0x49,
	0x0a, 0x15, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x42, 0x06, 0xba, 0x48, 0x03,
	0xc8, 0x01, 0x01, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x17, 0x0a, 0x15, 0x47, 0x65, 0x74,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x4c, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x05,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x22, 0x33, 0x0a, 0x17, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73,
	0x74, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x98, 0x01,
	0x1a, 0x52, 0x02, 0x69, 0x64, 0x22, 0x1a, 0x0a, 0x18, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x32, 0x9b, 0x02, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x52, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x12, 0x1e, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x0e, 0x47,
	0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x12, 0x1f, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20,
	0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x5b, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x12, 0x21, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x73, 0x74, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63,
	0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x6f, 0x75, 0x74, 0x6c, 0x61, 0x6e,
	0x64, 0x73, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_admin_v1_contestant_proto_rawDescOnce sync.Once
	file_admin_v1_contestant_proto_rawDescData = file_admin_v1_contestant_proto_rawDesc
)

func file_admin_v1_contestant_proto_rawDescGZIP() []byte {
	file_admin_v1_contestant_proto_rawDescOnce.Do(func() {
		file_admin_v1_contestant_proto_rawDescData = protoimpl.X.CompressGZIP(file_admin_v1_contestant_proto_rawDescData)
	})
	return file_admin_v1_contestant_proto_rawDescData
}

var file_admin_v1_contestant_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_admin_v1_contestant_proto_goTypes = []any{
	(*Contestant)(nil),               // 0: admin.v1.Contestant
	(*GetContestantRequest)(nil),     // 1: admin.v1.GetContestantRequest
	(*GetContestantResponse)(nil),    // 2: admin.v1.GetContestantResponse
	(*GetContestantsRequest)(nil),    // 3: admin.v1.GetContestantsRequest
	(*GetContestantsResponse)(nil),   // 4: admin.v1.GetContestantsResponse
	(*DeleteContestantRequest)(nil),  // 5: admin.v1.DeleteContestantRequest
	(*DeleteContestantResponse)(nil), // 6: admin.v1.DeleteContestantResponse
}
var file_admin_v1_contestant_proto_depIdxs = []int32{
	0, // 0: admin.v1.GetContestantResponse.user:type_name -> admin.v1.Contestant
	0, // 1: admin.v1.GetContestantsResponse.users:type_name -> admin.v1.Contestant
	1, // 2: admin.v1.ContestantService.GetContestant:input_type -> admin.v1.GetContestantRequest
	3, // 3: admin.v1.ContestantService.GetContestants:input_type -> admin.v1.GetContestantsRequest
	5, // 4: admin.v1.ContestantService.DeleteContestant:input_type -> admin.v1.DeleteContestantRequest
	2, // 5: admin.v1.ContestantService.GetContestant:output_type -> admin.v1.GetContestantResponse
	4, // 6: admin.v1.ContestantService.GetContestants:output_type -> admin.v1.GetContestantsResponse
	6, // 7: admin.v1.ContestantService.DeleteContestant:output_type -> admin.v1.DeleteContestantResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_admin_v1_contestant_proto_init() }
func file_admin_v1_contestant_proto_init() {
	if File_admin_v1_contestant_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_admin_v1_contestant_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_contestant_proto_goTypes,
		DependencyIndexes: file_admin_v1_contestant_proto_depIdxs,
		MessageInfos:      file_admin_v1_contestant_proto_msgTypes,
	}.Build()
	File_admin_v1_contestant_proto = out.File
	file_admin_v1_contestant_proto_rawDesc = nil
	file_admin_v1_contestant_proto_goTypes = nil
	file_admin_v1_contestant_proto_depIdxs = nil
}
