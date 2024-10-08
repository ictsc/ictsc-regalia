// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: admin/v1/admin_env.proto

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

type GetAdminConnectionInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAdminConnectionInfoRequest) Reset() {
	*x = GetAdminConnectionInfoRequest{}
	mi := &file_admin_v1_admin_env_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAdminConnectionInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAdminConnectionInfoRequest) ProtoMessage() {}

func (x *GetAdminConnectionInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_admin_env_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAdminConnectionInfoRequest.ProtoReflect.Descriptor instead.
func (*GetAdminConnectionInfoRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_admin_env_proto_rawDescGZIP(), []int{0}
}

type GetAdminConnectionInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bastion *Bastion `protobuf:"bytes,1,opt,name=bastion,proto3" json:"bastion,omitempty"`
}

func (x *GetAdminConnectionInfoResponse) Reset() {
	*x = GetAdminConnectionInfoResponse{}
	mi := &file_admin_v1_admin_env_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAdminConnectionInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAdminConnectionInfoResponse) ProtoMessage() {}

func (x *GetAdminConnectionInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_admin_env_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAdminConnectionInfoResponse.ProtoReflect.Descriptor instead.
func (*GetAdminConnectionInfoResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_admin_env_proto_rawDescGZIP(), []int{1}
}

func (x *GetAdminConnectionInfoResponse) GetBastion() *Bastion {
	if x != nil {
		return x.Bastion
	}
	return nil
}

type PutAdminConnectionInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bastion *Bastion `protobuf:"bytes,1,opt,name=bastion,proto3" json:"bastion,omitempty"`
}

func (x *PutAdminConnectionInfoRequest) Reset() {
	*x = PutAdminConnectionInfoRequest{}
	mi := &file_admin_v1_admin_env_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PutAdminConnectionInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutAdminConnectionInfoRequest) ProtoMessage() {}

func (x *PutAdminConnectionInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_admin_env_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutAdminConnectionInfoRequest.ProtoReflect.Descriptor instead.
func (*PutAdminConnectionInfoRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_admin_env_proto_rawDescGZIP(), []int{2}
}

func (x *PutAdminConnectionInfoRequest) GetBastion() *Bastion {
	if x != nil {
		return x.Bastion
	}
	return nil
}

type PutAdminConnectionInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PutAdminConnectionInfoResponse) Reset() {
	*x = PutAdminConnectionInfoResponse{}
	mi := &file_admin_v1_admin_env_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PutAdminConnectionInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutAdminConnectionInfoResponse) ProtoMessage() {}

func (x *PutAdminConnectionInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_admin_env_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutAdminConnectionInfoResponse.ProtoReflect.Descriptor instead.
func (*PutAdminConnectionInfoResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_admin_env_proto_rawDescGZIP(), []int{3}
}

var File_admin_v1_admin_env_proto protoreflect.FileDescriptor

var file_admin_v1_admin_env_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x5f, 0x65, 0x6e, 0x76, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x74,
	0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1f, 0x0a, 0x1d, 0x47, 0x65, 0x74, 0x41, 0x64, 0x6d,
	0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x55, 0x0a, 0x1e, 0x47, 0x65, 0x74, 0x41, 0x64,
	0x6d, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x07, 0x62, 0x61, 0x73,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x06, 0xba,
	0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x07, 0x62, 0x61, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x4c,
	0x0a, 0x1d, 0x50, 0x75, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2b, 0x0a, 0x07, 0x62, 0x61, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x73, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x07, 0x62, 0x61, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x20, 0x0a, 0x1e,
	0x50, 0x75, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xeb,
	0x01, 0x0a, 0x0f, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x45, 0x6e, 0x76, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x6b, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x43, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x27, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x6b, 0x0a, 0x16, 0x50, 0x75, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x27, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x28, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75,
	0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x41, 0x5a, 0x3f,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63,
	0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x6f, 0x75, 0x74, 0x6c, 0x61, 0x6e, 0x64, 0x73, 0x2f,
	0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_admin_v1_admin_env_proto_rawDescOnce sync.Once
	file_admin_v1_admin_env_proto_rawDescData = file_admin_v1_admin_env_proto_rawDesc
)

func file_admin_v1_admin_env_proto_rawDescGZIP() []byte {
	file_admin_v1_admin_env_proto_rawDescOnce.Do(func() {
		file_admin_v1_admin_env_proto_rawDescData = protoimpl.X.CompressGZIP(file_admin_v1_admin_env_proto_rawDescData)
	})
	return file_admin_v1_admin_env_proto_rawDescData
}

var file_admin_v1_admin_env_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_admin_v1_admin_env_proto_goTypes = []any{
	(*GetAdminConnectionInfoRequest)(nil),  // 0: admin.v1.GetAdminConnectionInfoRequest
	(*GetAdminConnectionInfoResponse)(nil), // 1: admin.v1.GetAdminConnectionInfoResponse
	(*PutAdminConnectionInfoRequest)(nil),  // 2: admin.v1.PutAdminConnectionInfoRequest
	(*PutAdminConnectionInfoResponse)(nil), // 3: admin.v1.PutAdminConnectionInfoResponse
	(*Bastion)(nil),                        // 4: admin.v1.Bastion
}
var file_admin_v1_admin_env_proto_depIdxs = []int32{
	4, // 0: admin.v1.GetAdminConnectionInfoResponse.bastion:type_name -> admin.v1.Bastion
	4, // 1: admin.v1.PutAdminConnectionInfoRequest.bastion:type_name -> admin.v1.Bastion
	0, // 2: admin.v1.AdminEnvService.GetAdminConnectionInfo:input_type -> admin.v1.GetAdminConnectionInfoRequest
	2, // 3: admin.v1.AdminEnvService.PutAdminConnectionInfo:input_type -> admin.v1.PutAdminConnectionInfoRequest
	1, // 4: admin.v1.AdminEnvService.GetAdminConnectionInfo:output_type -> admin.v1.GetAdminConnectionInfoResponse
	3, // 5: admin.v1.AdminEnvService.PutAdminConnectionInfo:output_type -> admin.v1.PutAdminConnectionInfoResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_admin_v1_admin_env_proto_init() }
func file_admin_v1_admin_env_proto_init() {
	if File_admin_v1_admin_env_proto != nil {
		return
	}
	file_admin_v1_team_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_admin_v1_admin_env_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_admin_env_proto_goTypes,
		DependencyIndexes: file_admin_v1_admin_env_proto_depIdxs,
		MessageInfos:      file_admin_v1_admin_env_proto_msgTypes,
	}.Build()
	File_admin_v1_admin_env_proto = out.File
	file_admin_v1_admin_env_proto_rawDesc = nil
	file_admin_v1_admin_env_proto_goTypes = nil
	file_admin_v1_admin_env_proto_depIdxs = nil
}
