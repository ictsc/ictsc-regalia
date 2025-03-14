// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: admin/v1/rule.proto

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

type Rule struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Markdown      string                 `protobuf:"bytes,2,opt,name=markdown,proto3" json:"markdown,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Rule) Reset() {
	*x = Rule{}
	mi := &file_admin_v1_rule_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Rule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule) ProtoMessage() {}

func (x *Rule) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_rule_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule.ProtoReflect.Descriptor instead.
func (*Rule) Descriptor() ([]byte, []int) {
	return file_admin_v1_rule_proto_rawDescGZIP(), []int{0}
}

func (x *Rule) GetMarkdown() string {
	if x != nil {
		return x.Markdown
	}
	return ""
}

type GetRuleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRuleRequest) Reset() {
	*x = GetRuleRequest{}
	mi := &file_admin_v1_rule_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRuleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRuleRequest) ProtoMessage() {}

func (x *GetRuleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_rule_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRuleRequest.ProtoReflect.Descriptor instead.
func (*GetRuleRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_rule_proto_rawDescGZIP(), []int{1}
}

type GetRuleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Rule          *Rule                  `protobuf:"bytes,1,opt,name=rule,proto3" json:"rule,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRuleResponse) Reset() {
	*x = GetRuleResponse{}
	mi := &file_admin_v1_rule_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRuleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRuleResponse) ProtoMessage() {}

func (x *GetRuleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_rule_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRuleResponse.ProtoReflect.Descriptor instead.
func (*GetRuleResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_rule_proto_rawDescGZIP(), []int{2}
}

func (x *GetRuleResponse) GetRule() *Rule {
	if x != nil {
		return x.Rule
	}
	return nil
}

type UpdateRuleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Rule          *Rule                  `protobuf:"bytes,1,opt,name=rule,proto3" json:"rule,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateRuleRequest) Reset() {
	*x = UpdateRuleRequest{}
	mi := &file_admin_v1_rule_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateRuleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRuleRequest) ProtoMessage() {}

func (x *UpdateRuleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_rule_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRuleRequest.ProtoReflect.Descriptor instead.
func (*UpdateRuleRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_rule_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateRuleRequest) GetRule() *Rule {
	if x != nil {
		return x.Rule
	}
	return nil
}

type UpdateRuleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateRuleResponse) Reset() {
	*x = UpdateRuleResponse{}
	mi := &file_admin_v1_rule_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateRuleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRuleResponse) ProtoMessage() {}

func (x *UpdateRuleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_rule_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRuleResponse.ProtoReflect.Descriptor instead.
func (*UpdateRuleResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_rule_proto_rawDescGZIP(), []int{4}
}

var File_admin_v1_rule_proto protoreflect.FileDescriptor

var file_admin_v1_rule_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x75, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x22,
	0x22, 0x0a, 0x04, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x61, 0x72, 0x6b, 0x64,
	0x6f, 0x77, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x61, 0x72, 0x6b, 0x64,
	0x6f, 0x77, 0x6e, 0x22, 0x10, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x35, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x52, 0x75, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x72, 0x75, 0x6c, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x22, 0x37, 0x0a, 0x11,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x22, 0x0a, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52,
	0x04, 0x72, 0x75, 0x6c, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x96, 0x01, 0x0a, 0x0b,
	0x52, 0x75, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x47,
	0x65, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x18, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52,
	0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0a, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x1b, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x72,
	0x65, 0x67, 0x61, 0x6c, 0x69, 0x61, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76,
	0x31, 0x3b, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_admin_v1_rule_proto_rawDescOnce sync.Once
	file_admin_v1_rule_proto_rawDescData []byte
)

func file_admin_v1_rule_proto_rawDescGZIP() []byte {
	file_admin_v1_rule_proto_rawDescOnce.Do(func() {
		file_admin_v1_rule_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_admin_v1_rule_proto_rawDesc), len(file_admin_v1_rule_proto_rawDesc)))
	})
	return file_admin_v1_rule_proto_rawDescData
}

var file_admin_v1_rule_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_admin_v1_rule_proto_goTypes = []any{
	(*Rule)(nil),               // 0: admin.v1.Rule
	(*GetRuleRequest)(nil),     // 1: admin.v1.GetRuleRequest
	(*GetRuleResponse)(nil),    // 2: admin.v1.GetRuleResponse
	(*UpdateRuleRequest)(nil),  // 3: admin.v1.UpdateRuleRequest
	(*UpdateRuleResponse)(nil), // 4: admin.v1.UpdateRuleResponse
}
var file_admin_v1_rule_proto_depIdxs = []int32{
	0, // 0: admin.v1.GetRuleResponse.rule:type_name -> admin.v1.Rule
	0, // 1: admin.v1.UpdateRuleRequest.rule:type_name -> admin.v1.Rule
	1, // 2: admin.v1.RuleService.GetRule:input_type -> admin.v1.GetRuleRequest
	3, // 3: admin.v1.RuleService.UpdateRule:input_type -> admin.v1.UpdateRuleRequest
	2, // 4: admin.v1.RuleService.GetRule:output_type -> admin.v1.GetRuleResponse
	4, // 5: admin.v1.RuleService.UpdateRule:output_type -> admin.v1.UpdateRuleResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_admin_v1_rule_proto_init() }
func file_admin_v1_rule_proto_init() {
	if File_admin_v1_rule_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_admin_v1_rule_proto_rawDesc), len(file_admin_v1_rule_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_rule_proto_goTypes,
		DependencyIndexes: file_admin_v1_rule_proto_depIdxs,
		MessageInfos:      file_admin_v1_rule_proto_msgTypes,
	}.Build()
	File_admin_v1_rule_proto = out.File
	file_admin_v1_rule_proto_goTypes = nil
	file_admin_v1_rule_proto_depIdxs = nil
}
