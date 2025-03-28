// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: contestant/v1/notice.proto

package contestantv1

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

type Notice struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Body          string                 `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	Slug          string                 `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Notice) Reset() {
	*x = Notice{}
	mi := &file_contestant_v1_notice_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Notice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notice) ProtoMessage() {}

func (x *Notice) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_notice_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notice.ProtoReflect.Descriptor instead.
func (*Notice) Descriptor() ([]byte, []int) {
	return file_contestant_v1_notice_proto_rawDescGZIP(), []int{0}
}

func (x *Notice) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Notice) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *Notice) GetSlug() string {
	if x != nil {
		return x.Slug
	}
	return ""
}

type ListNoticesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListNoticesRequest) Reset() {
	*x = ListNoticesRequest{}
	mi := &file_contestant_v1_notice_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListNoticesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListNoticesRequest) ProtoMessage() {}

func (x *ListNoticesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_notice_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListNoticesRequest.ProtoReflect.Descriptor instead.
func (*ListNoticesRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_notice_proto_rawDescGZIP(), []int{1}
}

type ListNoticesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Notices       []*Notice              `protobuf:"bytes,1,rep,name=notices,proto3" json:"notices,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListNoticesResponse) Reset() {
	*x = ListNoticesResponse{}
	mi := &file_contestant_v1_notice_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListNoticesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListNoticesResponse) ProtoMessage() {}

func (x *ListNoticesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_notice_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListNoticesResponse.ProtoReflect.Descriptor instead.
func (*ListNoticesResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_notice_proto_rawDescGZIP(), []int{2}
}

func (x *ListNoticesResponse) GetNotices() []*Notice {
	if x != nil {
		return x.Notices
	}
	return nil
}

var File_contestant_v1_notice_proto protoreflect.FileDescriptor

var file_contestant_v1_notice_proto_rawDesc = string([]byte{
	0x0a, 0x1a, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x22, 0x46, 0x0a, 0x06, 0x4e,
	0x6f, 0x74, 0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73,
	0x6c, 0x75, 0x67, 0x22, 0x14, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x46, 0x0a, 0x13, 0x4c, 0x69, 0x73,
	0x74, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x2f, 0x0a, 0x07, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x52, 0x07, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65,
	0x73, 0x32, 0x65, 0x0a, 0x0d, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x54, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65,
	0x73, 0x12, 0x21, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74,
	0x73, 0x63, 0x2d, 0x72, 0x65, 0x67, 0x61, 0x6c, 0x69, 0x61, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65,
	0x6e, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x73, 0x74, 0x61, 0x6e, 0x74, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_contestant_v1_notice_proto_rawDescOnce sync.Once
	file_contestant_v1_notice_proto_rawDescData []byte
)

func file_contestant_v1_notice_proto_rawDescGZIP() []byte {
	file_contestant_v1_notice_proto_rawDescOnce.Do(func() {
		file_contestant_v1_notice_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_contestant_v1_notice_proto_rawDesc), len(file_contestant_v1_notice_proto_rawDesc)))
	})
	return file_contestant_v1_notice_proto_rawDescData
}

var file_contestant_v1_notice_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_contestant_v1_notice_proto_goTypes = []any{
	(*Notice)(nil),              // 0: contestant.v1.Notice
	(*ListNoticesRequest)(nil),  // 1: contestant.v1.ListNoticesRequest
	(*ListNoticesResponse)(nil), // 2: contestant.v1.ListNoticesResponse
}
var file_contestant_v1_notice_proto_depIdxs = []int32{
	0, // 0: contestant.v1.ListNoticesResponse.notices:type_name -> contestant.v1.Notice
	1, // 1: contestant.v1.NoticeService.ListNotices:input_type -> contestant.v1.ListNoticesRequest
	2, // 2: contestant.v1.NoticeService.ListNotices:output_type -> contestant.v1.ListNoticesResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_contestant_v1_notice_proto_init() }
func file_contestant_v1_notice_proto_init() {
	if File_contestant_v1_notice_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_contestant_v1_notice_proto_rawDesc), len(file_contestant_v1_notice_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contestant_v1_notice_proto_goTypes,
		DependencyIndexes: file_contestant_v1_notice_proto_depIdxs,
		MessageInfos:      file_contestant_v1_notice_proto_msgTypes,
	}.Build()
	File_contestant_v1_notice_proto = out.File
	file_contestant_v1_notice_proto_goTypes = nil
	file_contestant_v1_notice_proto_depIdxs = nil
}
