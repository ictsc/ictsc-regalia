// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: admin/v1/notice.proto

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

type Notice struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Path          string                 `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Markdown      string                 `protobuf:"bytes,3,opt,name=markdown,proto3" json:"markdown,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Notice) Reset() {
	*x = Notice{}
	mi := &file_admin_v1_notice_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Notice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notice) ProtoMessage() {}

func (x *Notice) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_notice_proto_msgTypes[0]
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
	return file_admin_v1_notice_proto_rawDescGZIP(), []int{0}
}

func (x *Notice) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Notice) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Notice) GetMarkdown() string {
	if x != nil {
		return x.Markdown
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
	mi := &file_admin_v1_notice_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListNoticesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListNoticesRequest) ProtoMessage() {}

func (x *ListNoticesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_notice_proto_msgTypes[1]
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
	return file_admin_v1_notice_proto_rawDescGZIP(), []int{1}
}

type ListNoticesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Notices       []*Notice              `protobuf:"bytes,1,rep,name=notices,proto3" json:"notices,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListNoticesResponse) Reset() {
	*x = ListNoticesResponse{}
	mi := &file_admin_v1_notice_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListNoticesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListNoticesResponse) ProtoMessage() {}

func (x *ListNoticesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_notice_proto_msgTypes[2]
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
	return file_admin_v1_notice_proto_rawDescGZIP(), []int{2}
}

func (x *ListNoticesResponse) GetNotices() []*Notice {
	if x != nil {
		return x.Notices
	}
	return nil
}

type SyncNoticesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Path          string                 `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncNoticesRequest) Reset() {
	*x = SyncNoticesRequest{}
	mi := &file_admin_v1_notice_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncNoticesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncNoticesRequest) ProtoMessage() {}

func (x *SyncNoticesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_notice_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncNoticesRequest.ProtoReflect.Descriptor instead.
func (*SyncNoticesRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_notice_proto_rawDescGZIP(), []int{3}
}

func (x *SyncNoticesRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type SyncNoticesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Notice        *Notice                `protobuf:"bytes,1,opt,name=notice,proto3" json:"notice,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncNoticesResponse) Reset() {
	*x = SyncNoticesResponse{}
	mi := &file_admin_v1_notice_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncNoticesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncNoticesResponse) ProtoMessage() {}

func (x *SyncNoticesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_notice_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncNoticesResponse.ProtoReflect.Descriptor instead.
func (*SyncNoticesResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_notice_proto_rawDescGZIP(), []int{4}
}

func (x *SyncNoticesResponse) GetNotice() *Notice {
	if x != nil {
		return x.Notice
	}
	return nil
}

var File_admin_v1_notice_proto protoreflect.FileDescriptor

var file_admin_v1_notice_proto_rawDesc = string([]byte{
	0x0a, 0x15, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x22, 0x4e, 0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x61, 0x72, 0x6b, 0x64, 0x6f, 0x77,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x61, 0x72, 0x6b, 0x64, 0x6f, 0x77,
	0x6e, 0x22, 0x14, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x41, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x4e,
	0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a,
	0x0a, 0x07, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x52, 0x07, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x22, 0x28, 0x0a, 0x12, 0x53, 0x79,
	0x6e, 0x63, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x70, 0x61, 0x74, 0x68, 0x22, 0x3f, 0x0a, 0x13, 0x53, 0x79, 0x6e, 0x63, 0x4e, 0x6f, 0x74, 0x69,
	0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x6e,
	0x6f, 0x74, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x52, 0x06, 0x6e,
	0x6f, 0x74, 0x69, 0x63, 0x65, 0x32, 0xa7, 0x01, 0x0a, 0x0d, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4a, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x4e,
	0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x12, 0x1c, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0b, 0x53, 0x79, 0x6e, 0x63, 0x4e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x73, 0x12, 0x1c, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x79,
	0x6e, 0x63, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x79, 0x6e, 0x63,
	0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63,
	0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x72, 0x65, 0x67, 0x61, 0x6c, 0x69,
	0x61, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_admin_v1_notice_proto_rawDescOnce sync.Once
	file_admin_v1_notice_proto_rawDescData []byte
)

func file_admin_v1_notice_proto_rawDescGZIP() []byte {
	file_admin_v1_notice_proto_rawDescOnce.Do(func() {
		file_admin_v1_notice_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_admin_v1_notice_proto_rawDesc), len(file_admin_v1_notice_proto_rawDesc)))
	})
	return file_admin_v1_notice_proto_rawDescData
}

var file_admin_v1_notice_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_admin_v1_notice_proto_goTypes = []any{
	(*Notice)(nil),              // 0: admin.v1.Notice
	(*ListNoticesRequest)(nil),  // 1: admin.v1.ListNoticesRequest
	(*ListNoticesResponse)(nil), // 2: admin.v1.ListNoticesResponse
	(*SyncNoticesRequest)(nil),  // 3: admin.v1.SyncNoticesRequest
	(*SyncNoticesResponse)(nil), // 4: admin.v1.SyncNoticesResponse
}
var file_admin_v1_notice_proto_depIdxs = []int32{
	0, // 0: admin.v1.ListNoticesResponse.notices:type_name -> admin.v1.Notice
	0, // 1: admin.v1.SyncNoticesResponse.notice:type_name -> admin.v1.Notice
	1, // 2: admin.v1.NoticeService.ListNotices:input_type -> admin.v1.ListNoticesRequest
	3, // 3: admin.v1.NoticeService.SyncNotices:input_type -> admin.v1.SyncNoticesRequest
	2, // 4: admin.v1.NoticeService.ListNotices:output_type -> admin.v1.ListNoticesResponse
	4, // 5: admin.v1.NoticeService.SyncNotices:output_type -> admin.v1.SyncNoticesResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_admin_v1_notice_proto_init() }
func file_admin_v1_notice_proto_init() {
	if File_admin_v1_notice_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_admin_v1_notice_proto_rawDesc), len(file_admin_v1_notice_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_notice_proto_goTypes,
		DependencyIndexes: file_admin_v1_notice_proto_depIdxs,
		MessageInfos:      file_admin_v1_notice_proto_msgTypes,
	}.Build()
	File_admin_v1_notice_proto = out.File
	file_admin_v1_notice_proto_goTypes = nil
	file_admin_v1_notice_proto_depIdxs = nil
}
