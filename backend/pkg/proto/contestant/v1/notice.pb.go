// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
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

const file_contestant_v1_notice_proto_rawDesc = "" +
	"\n" +
	"\x1acontestant/v1/notice.proto\x12\rcontestant.v1\"F\n" +
	"\x06Notice\x12\x14\n" +
	"\x05title\x18\x01 \x01(\tR\x05title\x12\x12\n" +
	"\x04body\x18\x02 \x01(\tR\x04body\x12\x12\n" +
	"\x04slug\x18\x03 \x01(\tR\x04slug\"\x14\n" +
	"\x12ListNoticesRequest\"F\n" +
	"\x13ListNoticesResponse\x12/\n" +
	"\anotices\x18\x01 \x03(\v2\x15.contestant.v1.NoticeR\anotices2e\n" +
	"\rNoticeService\x12T\n" +
	"\vListNotices\x12!.contestant.v1.ListNoticesRequest\x1a\".contestant.v1.ListNoticesResponseBMZKgithub.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1;contestantv1b\x06proto3"

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
