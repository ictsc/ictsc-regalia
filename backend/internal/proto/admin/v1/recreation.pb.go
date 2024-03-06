// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: admin/v1/recreation.proto

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

type Status int32

const (
	Status_STATUS_UNSPECIFIED Status = 0
	Status_STATUS_IN_PROGRESS Status = 1
	Status_STATUS_COMPLETED   Status = 2
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "STATUS_UNSPECIFIED",
		1: "STATUS_IN_PROGRESS",
		2: "STATUS_COMPLETED",
	}
	Status_value = map[string]int32{
		"STATUS_UNSPECIFIED": 0,
		"STATUS_IN_PROGRESS": 1,
		"STATUS_COMPLETED":   2,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_admin_v1_recreation_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_admin_v1_recreation_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_admin_v1_recreation_proto_rawDescGZIP(), []int{0}
}

type Recreation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TeamId    string                 `protobuf:"bytes,2,opt,name=team_id,json=teamId,proto3" json:"team_id,omitempty"`
	ProblemId string                 `protobuf:"bytes,3,opt,name=problem_id,json=problemId,proto3" json:"problem_id,omitempty"`
	Status    Status                 `protobuf:"varint,4,opt,name=status,proto3,enum=admin.v1.Status" json:"status,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *Recreation) Reset() {
	*x = Recreation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_recreation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Recreation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Recreation) ProtoMessage() {}

func (x *Recreation) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_recreation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Recreation.ProtoReflect.Descriptor instead.
func (*Recreation) Descriptor() ([]byte, []int) {
	return file_admin_v1_recreation_proto_rawDescGZIP(), []int{0}
}

func (x *Recreation) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Recreation) GetTeamId() string {
	if x != nil {
		return x.TeamId
	}
	return ""
}

func (x *Recreation) GetProblemId() string {
	if x != nil {
		return x.ProblemId
	}
	return ""
}

func (x *Recreation) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_STATUS_UNSPECIFIED
}

func (x *Recreation) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type GetRecreationsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TeamId string `protobuf:"bytes,1,opt,name=team_id,json=teamId,proto3" json:"team_id,omitempty"`
}

func (x *GetRecreationsRequest) Reset() {
	*x = GetRecreationsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_recreation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRecreationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecreationsRequest) ProtoMessage() {}

func (x *GetRecreationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_recreation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRecreationsRequest.ProtoReflect.Descriptor instead.
func (*GetRecreationsRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_recreation_proto_rawDescGZIP(), []int{1}
}

func (x *GetRecreationsRequest) GetTeamId() string {
	if x != nil {
		return x.TeamId
	}
	return ""
}

type GetRecreationsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Recreations []*Recreation `protobuf:"bytes,1,rep,name=recreations,proto3" json:"recreations,omitempty"`
}

func (x *GetRecreationsResponse) Reset() {
	*x = GetRecreationsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_recreation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRecreationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecreationsResponse) ProtoMessage() {}

func (x *GetRecreationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_recreation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRecreationsResponse.ProtoReflect.Descriptor instead.
func (*GetRecreationsResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_recreation_proto_rawDescGZIP(), []int{2}
}

func (x *GetRecreationsResponse) GetRecreations() []*Recreation {
	if x != nil {
		return x.Recreations
	}
	return nil
}

type PostAdminRecreationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProblemId string `protobuf:"bytes,1,opt,name=problem_id,json=problemId,proto3" json:"problem_id,omitempty"`
}

func (x *PostAdminRecreationRequest) Reset() {
	*x = PostAdminRecreationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_recreation_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostAdminRecreationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostAdminRecreationRequest) ProtoMessage() {}

func (x *PostAdminRecreationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_recreation_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostAdminRecreationRequest.ProtoReflect.Descriptor instead.
func (*PostAdminRecreationRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_recreation_proto_rawDescGZIP(), []int{3}
}

func (x *PostAdminRecreationRequest) GetProblemId() string {
	if x != nil {
		return x.ProblemId
	}
	return ""
}

type PostAdminRecreationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Recreation *Recreation `protobuf:"bytes,1,opt,name=recreation,proto3" json:"recreation,omitempty"`
}

func (x *PostAdminRecreationResponse) Reset() {
	*x = PostAdminRecreationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_v1_recreation_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostAdminRecreationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostAdminRecreationResponse) ProtoMessage() {}

func (x *PostAdminRecreationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_recreation_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostAdminRecreationResponse.ProtoReflect.Descriptor instead.
func (*PostAdminRecreationResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_recreation_proto_rawDescGZIP(), []int{4}
}

func (x *PostAdminRecreationResponse) GetRecreation() *Recreation {
	if x != nil {
		return x.Recreation
	}
	return nil
}

var File_admin_v1_recreation_proto protoreflect.FileDescriptor

var file_admin_v1_recreation_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xe9, 0x01, 0x0a, 0x0a, 0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xba, 0x48, 0x05, 0x72, 0x03, 0x98, 0x01, 0x1a, 0x52, 0x02, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x07,
	0x74, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba,
	0x48, 0x05, 0x72, 0x03, 0x98, 0x01, 0x1a, 0x52, 0x06, 0x74, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x12,
	0x27, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x98, 0x01, 0x1a, 0x52, 0x09, 0x70,
	0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x32, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x08, 0xba, 0x48, 0x05, 0x82,
	0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x41, 0x0a, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x06, 0xba, 0x48,
	0x03, 0xc8, 0x01, 0x01, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22,
	0x3a, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x74, 0x65, 0x61, 0x6d,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03,
	0x98, 0x01, 0x1a, 0x52, 0x06, 0x74, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x22, 0x58, 0x0a, 0x16, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x45, 0x0a, 0x1a, 0x50, 0x6f, 0x73, 0x74, 0x41, 0x64, 0x6d,
	0x69, 0x6e, 0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x98, 0x01,
	0x1a, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x49, 0x64, 0x22, 0x5b, 0x0a, 0x1b,
	0x50, 0x6f, 0x73, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x0a, 0x72,
	0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x0a, 0x72,
	0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2a, 0x4e, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x12, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x53,
	0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x49, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x47, 0x52, 0x45, 0x53,
	0x53, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x43, 0x4f,
	0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x02, 0x32, 0xcc, 0x01, 0x0a, 0x11, 0x52, 0x65,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x53, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x1f, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x20, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x62, 0x0a, 0x13, 0x50, 0x6f, 0x73, 0x74, 0x41, 0x64, 0x6d, 0x69,
	0x6e, 0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x25, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73,
	0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x52, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74,
	0x73, 0x63, 0x2d, 0x6f, 0x75, 0x74, 0x6c, 0x61, 0x6e, 0x64, 0x73, 0x2f, 0x62, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_admin_v1_recreation_proto_rawDescOnce sync.Once
	file_admin_v1_recreation_proto_rawDescData = file_admin_v1_recreation_proto_rawDesc
)

func file_admin_v1_recreation_proto_rawDescGZIP() []byte {
	file_admin_v1_recreation_proto_rawDescOnce.Do(func() {
		file_admin_v1_recreation_proto_rawDescData = protoimpl.X.CompressGZIP(file_admin_v1_recreation_proto_rawDescData)
	})
	return file_admin_v1_recreation_proto_rawDescData
}

var file_admin_v1_recreation_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_admin_v1_recreation_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_admin_v1_recreation_proto_goTypes = []interface{}{
	(Status)(0),                         // 0: admin.v1.Status
	(*Recreation)(nil),                  // 1: admin.v1.Recreation
	(*GetRecreationsRequest)(nil),       // 2: admin.v1.GetRecreationsRequest
	(*GetRecreationsResponse)(nil),      // 3: admin.v1.GetRecreationsResponse
	(*PostAdminRecreationRequest)(nil),  // 4: admin.v1.PostAdminRecreationRequest
	(*PostAdminRecreationResponse)(nil), // 5: admin.v1.PostAdminRecreationResponse
	(*timestamppb.Timestamp)(nil),       // 6: google.protobuf.Timestamp
}
var file_admin_v1_recreation_proto_depIdxs = []int32{
	0, // 0: admin.v1.Recreation.status:type_name -> admin.v1.Status
	6, // 1: admin.v1.Recreation.created_at:type_name -> google.protobuf.Timestamp
	1, // 2: admin.v1.GetRecreationsResponse.recreations:type_name -> admin.v1.Recreation
	1, // 3: admin.v1.PostAdminRecreationResponse.recreation:type_name -> admin.v1.Recreation
	2, // 4: admin.v1.RecreationService.GetRecreations:input_type -> admin.v1.GetRecreationsRequest
	4, // 5: admin.v1.RecreationService.PostAdminRecreation:input_type -> admin.v1.PostAdminRecreationRequest
	3, // 6: admin.v1.RecreationService.GetRecreations:output_type -> admin.v1.GetRecreationsResponse
	5, // 7: admin.v1.RecreationService.PostAdminRecreation:output_type -> admin.v1.PostAdminRecreationResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_admin_v1_recreation_proto_init() }
func file_admin_v1_recreation_proto_init() {
	if File_admin_v1_recreation_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_admin_v1_recreation_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Recreation); i {
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
		file_admin_v1_recreation_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRecreationsRequest); i {
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
		file_admin_v1_recreation_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRecreationsResponse); i {
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
		file_admin_v1_recreation_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostAdminRecreationRequest); i {
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
		file_admin_v1_recreation_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostAdminRecreationResponse); i {
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
			RawDescriptor: file_admin_v1_recreation_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_recreation_proto_goTypes,
		DependencyIndexes: file_admin_v1_recreation_proto_depIdxs,
		EnumInfos:         file_admin_v1_recreation_proto_enumTypes,
		MessageInfos:      file_admin_v1_recreation_proto_msgTypes,
	}.Build()
	File_admin_v1_recreation_proto = out.File
	file_admin_v1_recreation_proto_rawDesc = nil
	file_admin_v1_recreation_proto_goTypes = nil
	file_admin_v1_recreation_proto_depIdxs = nil
}
