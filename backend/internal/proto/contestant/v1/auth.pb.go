// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: contestant/v1/auth.proto

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

type GetCallbackRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetCallbackRequest) Reset() {
	*x = GetCallbackRequest{}
	mi := &file_contestant_v1_auth_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCallbackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCallbackRequest) ProtoMessage() {}

func (x *GetCallbackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_auth_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCallbackRequest.ProtoReflect.Descriptor instead.
func (*GetCallbackRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_auth_proto_rawDescGZIP(), []int{0}
}

type GetCallbackResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RedirectUri string `protobuf:"bytes,1,opt,name=redirect_uri,json=redirectUri,proto3" json:"redirect_uri,omitempty"`
}

func (x *GetCallbackResponse) Reset() {
	*x = GetCallbackResponse{}
	mi := &file_contestant_v1_auth_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCallbackResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCallbackResponse) ProtoMessage() {}

func (x *GetCallbackResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_auth_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCallbackResponse.ProtoReflect.Descriptor instead.
func (*GetCallbackResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_auth_proto_rawDescGZIP(), []int{1}
}

func (x *GetCallbackResponse) GetRedirectUri() string {
	if x != nil {
		return x.RedirectUri
	}
	return ""
}

type PostCodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *PostCodeRequest) Reset() {
	*x = PostCodeRequest{}
	mi := &file_contestant_v1_auth_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PostCodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostCodeRequest) ProtoMessage() {}

func (x *PostCodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_auth_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostCodeRequest.ProtoReflect.Descriptor instead.
func (*PostCodeRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_auth_proto_rawDescGZIP(), []int{2}
}

func (x *PostCodeRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type PostCodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PostCodeResponse) Reset() {
	*x = PostCodeResponse{}
	mi := &file_contestant_v1_auth_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PostCodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostCodeResponse) ProtoMessage() {}

func (x *PostCodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_auth_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostCodeResponse.ProtoReflect.Descriptor instead.
func (*PostCodeResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_auth_proto_rawDescGZIP(), []int{3}
}

var File_contestant_v1_auth_proto protoreflect.FileDescriptor

var file_contestant_v1_auth_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x14, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6c,
	0x6c, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x42, 0x0a, 0x13,
	0x47, 0x65, 0x74, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x0c, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f,
	0x75, 0x72, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03,
	0x88, 0x01, 0x01, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x55, 0x72, 0x69,
	0x22, 0x2e, 0x0a, 0x0f, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x22, 0x12, 0x0a, 0x10, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0xb4, 0x01, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x56, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6c, 0x6c, 0x62,
	0x61, 0x63, 0x6b, 0x12, 0x21, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61,
	0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x08,
	0x50, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x46, 0x5a, 0x44, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2f,
	0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x6f, 0x75, 0x74, 0x6c, 0x61, 0x6e, 0x64, 0x73, 0x2f, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_contestant_v1_auth_proto_rawDescOnce sync.Once
	file_contestant_v1_auth_proto_rawDescData = file_contestant_v1_auth_proto_rawDesc
)

func file_contestant_v1_auth_proto_rawDescGZIP() []byte {
	file_contestant_v1_auth_proto_rawDescOnce.Do(func() {
		file_contestant_v1_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_contestant_v1_auth_proto_rawDescData)
	})
	return file_contestant_v1_auth_proto_rawDescData
}

var file_contestant_v1_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_contestant_v1_auth_proto_goTypes = []any{
	(*GetCallbackRequest)(nil),  // 0: contestant.v1.GetCallbackRequest
	(*GetCallbackResponse)(nil), // 1: contestant.v1.GetCallbackResponse
	(*PostCodeRequest)(nil),     // 2: contestant.v1.PostCodeRequest
	(*PostCodeResponse)(nil),    // 3: contestant.v1.PostCodeResponse
}
var file_contestant_v1_auth_proto_depIdxs = []int32{
	0, // 0: contestant.v1.AuthService.GetCallback:input_type -> contestant.v1.GetCallbackRequest
	2, // 1: contestant.v1.AuthService.PostCode:input_type -> contestant.v1.PostCodeRequest
	1, // 2: contestant.v1.AuthService.GetCallback:output_type -> contestant.v1.GetCallbackResponse
	3, // 3: contestant.v1.AuthService.PostCode:output_type -> contestant.v1.PostCodeResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_contestant_v1_auth_proto_init() }
func file_contestant_v1_auth_proto_init() {
	if File_contestant_v1_auth_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_contestant_v1_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contestant_v1_auth_proto_goTypes,
		DependencyIndexes: file_contestant_v1_auth_proto_depIdxs,
		MessageInfos:      file_contestant_v1_auth_proto_msgTypes,
	}.Build()
	File_contestant_v1_auth_proto = out.File
	file_contestant_v1_auth_proto_rawDesc = nil
	file_contestant_v1_auth_proto_goTypes = nil
	file_contestant_v1_auth_proto_depIdxs = nil
}
