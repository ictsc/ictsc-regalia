// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: admin/v1/problem.proto

package adminv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

type RedeployRuleType int32

const (
	RedeployRuleType_REDEPLOY_RULE_TYPE_UNSPECIFIED RedeployRuleType = 0
	// 自動での再展開ができない問題
	RedeployRuleType_REDEPLOY_RULE_TYPE_UNREDEPLOYABLE RedeployRuleType = 1
	// 再展開に最大点数への割合ペナルティがある問題
	RedeployRuleType_REDEPLOY_RULE_TYPE_PERCENTAGE_PENALTY RedeployRuleType = 2
)

// Enum value maps for RedeployRuleType.
var (
	RedeployRuleType_name = map[int32]string{
		0: "REDEPLOY_RULE_TYPE_UNSPECIFIED",
		1: "REDEPLOY_RULE_TYPE_UNREDEPLOYABLE",
		2: "REDEPLOY_RULE_TYPE_PERCENTAGE_PENALTY",
	}
	RedeployRuleType_value = map[string]int32{
		"REDEPLOY_RULE_TYPE_UNSPECIFIED":        0,
		"REDEPLOY_RULE_TYPE_UNREDEPLOYABLE":     1,
		"REDEPLOY_RULE_TYPE_PERCENTAGE_PENALTY": 2,
	}
)

func (x RedeployRuleType) Enum() *RedeployRuleType {
	p := new(RedeployRuleType)
	*p = x
	return p
}

func (x RedeployRuleType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RedeployRuleType) Descriptor() protoreflect.EnumDescriptor {
	return file_admin_v1_problem_proto_enumTypes[0].Descriptor()
}

func (RedeployRuleType) Type() protoreflect.EnumType {
	return &file_admin_v1_problem_proto_enumTypes[0]
}

func (x RedeployRuleType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RedeployRuleType.Descriptor instead.
func (RedeployRuleType) EnumDescriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{0}
}

type ProblemType int32

const (
	ProblemType_PROBLEM_TYPE_UNSPECIFIED ProblemType = 0
	ProblemType_PROBLEM_TYPE_DESCRIPTIVE ProblemType = 1
)

// Enum value maps for ProblemType.
var (
	ProblemType_name = map[int32]string{
		0: "PROBLEM_TYPE_UNSPECIFIED",
		1: "PROBLEM_TYPE_DESCRIPTIVE",
	}
	ProblemType_value = map[string]int32{
		"PROBLEM_TYPE_UNSPECIFIED": 0,
		"PROBLEM_TYPE_DESCRIPTIVE": 1,
	}
)

func (x ProblemType) Enum() *ProblemType {
	p := new(ProblemType)
	*p = x
	return p
}

func (x ProblemType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProblemType) Descriptor() protoreflect.EnumDescriptor {
	return file_admin_v1_problem_proto_enumTypes[1].Descriptor()
}

func (ProblemType) Type() protoreflect.EnumType {
	return &file_admin_v1_problem_proto_enumTypes[1]
}

func (x ProblemType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProblemType.Descriptor instead.
func (ProblemType) EnumDescriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{1}
}

type Problem struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 問題コード
	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	// タイトル
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	// 最大得点
	MaxScore uint32 `protobuf:"varint,3,opt,name=max_score,json=maxScore,proto3" json:"max_score,omitempty"`
	// 問題カテゴリー
	Category      string        `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	RedeployRule  *RedeployRule `protobuf:"bytes,5,opt,name=redeploy_rule,json=redeployRule,proto3" json:"redeploy_rule,omitempty"`
	Body          *ProblemBody  `protobuf:"bytes,6,opt,name=body,proto3" json:"body,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Problem) Reset() {
	*x = Problem{}
	mi := &file_admin_v1_problem_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Problem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Problem) ProtoMessage() {}

func (x *Problem) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Problem.ProtoReflect.Descriptor instead.
func (*Problem) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{0}
}

func (x *Problem) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *Problem) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Problem) GetMaxScore() uint32 {
	if x != nil {
		return x.MaxScore
	}
	return 0
}

func (x *Problem) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Problem) GetRedeployRule() *RedeployRule {
	if x != nil {
		return x.RedeployRule
	}
	return nil
}

func (x *Problem) GetBody() *ProblemBody {
	if x != nil {
		return x.Body
	}
	return nil
}

type RedeployRule struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Type  RedeployRuleType       `protobuf:"varint,1,opt,name=type,proto3,enum=admin.v1.RedeployRuleType" json:"type,omitempty"`
	// 再展開ペナルティの発生する再展開回数の閾値
	PenaltyThreshold uint32 `protobuf:"varint,2,opt,name=penalty_threshold,json=penaltyThreshold,proto3" json:"penalty_threshold,omitempty"`
	// 再展開ペナルティの割合
	PenaltyPercentage uint32 `protobuf:"varint,3,opt,name=penalty_percentage,json=penaltyPercentage,proto3" json:"penalty_percentage,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *RedeployRule) Reset() {
	*x = RedeployRule{}
	mi := &file_admin_v1_problem_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RedeployRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedeployRule) ProtoMessage() {}

func (x *RedeployRule) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedeployRule.ProtoReflect.Descriptor instead.
func (*RedeployRule) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{1}
}

func (x *RedeployRule) GetType() RedeployRuleType {
	if x != nil {
		return x.Type
	}
	return RedeployRuleType_REDEPLOY_RULE_TYPE_UNSPECIFIED
}

func (x *RedeployRule) GetPenaltyThreshold() uint32 {
	if x != nil {
		return x.PenaltyThreshold
	}
	return 0
}

func (x *RedeployRule) GetPenaltyPercentage() uint32 {
	if x != nil {
		return x.PenaltyPercentage
	}
	return 0
}

type ProblemBody struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Type  ProblemType            `protobuf:"varint,1,opt,name=type,proto3,enum=admin.v1.ProblemType" json:"type,omitempty"`
	// Types that are valid to be assigned to Body:
	//
	//	*ProblemBody_Descriptive
	Body          isProblemBody_Body `protobuf_oneof:"body"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProblemBody) Reset() {
	*x = ProblemBody{}
	mi := &file_admin_v1_problem_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProblemBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProblemBody) ProtoMessage() {}

func (x *ProblemBody) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProblemBody.ProtoReflect.Descriptor instead.
func (*ProblemBody) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{2}
}

func (x *ProblemBody) GetType() ProblemType {
	if x != nil {
		return x.Type
	}
	return ProblemType_PROBLEM_TYPE_UNSPECIFIED
}

func (x *ProblemBody) GetBody() isProblemBody_Body {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *ProblemBody) GetDescriptive() *DescriptiveProblem {
	if x != nil {
		if x, ok := x.Body.(*ProblemBody_Descriptive); ok {
			return x.Descriptive
		}
	}
	return nil
}

type isProblemBody_Body interface {
	isProblemBody_Body()
}

type ProblemBody_Descriptive struct {
	Descriptive *DescriptiveProblem `protobuf:"bytes,2,opt,name=descriptive,proto3,oneof"`
}

func (*ProblemBody_Descriptive) isProblemBody_Body() {}

type DescriptiveProblem struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Markdown 形式の問題文
	ProblemMarkdown string `protobuf:"bytes,3,opt,name=problem_markdown,json=problemMarkdown,proto3" json:"problem_markdown,omitempty"`
	// Markdown 形式の解説文
	ExplanationMarkdown string `protobuf:"bytes,4,opt,name=explanation_markdown,json=explanationMarkdown,proto3" json:"explanation_markdown,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *DescriptiveProblem) Reset() {
	*x = DescriptiveProblem{}
	mi := &file_admin_v1_problem_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DescriptiveProblem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescriptiveProblem) ProtoMessage() {}

func (x *DescriptiveProblem) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescriptiveProblem.ProtoReflect.Descriptor instead.
func (*DescriptiveProblem) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{3}
}

func (x *DescriptiveProblem) GetProblemMarkdown() string {
	if x != nil {
		return x.ProblemMarkdown
	}
	return ""
}

func (x *DescriptiveProblem) GetExplanationMarkdown() string {
	if x != nil {
		return x.ExplanationMarkdown
	}
	return ""
}

type ListProblemsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListProblemsRequest) Reset() {
	*x = ListProblemsRequest{}
	mi := &file_admin_v1_problem_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListProblemsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProblemsRequest) ProtoMessage() {}

func (x *ListProblemsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProblemsRequest.ProtoReflect.Descriptor instead.
func (*ListProblemsRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{4}
}

type ListProblemsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Problems      []*Problem             `protobuf:"bytes,1,rep,name=problems,proto3" json:"problems,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListProblemsResponse) Reset() {
	*x = ListProblemsResponse{}
	mi := &file_admin_v1_problem_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListProblemsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProblemsResponse) ProtoMessage() {}

func (x *ListProblemsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProblemsResponse.ProtoReflect.Descriptor instead.
func (*ListProblemsResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{5}
}

func (x *ListProblemsResponse) GetProblems() []*Problem {
	if x != nil {
		return x.Problems
	}
	return nil
}

type GetProblemRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetProblemRequest) Reset() {
	*x = GetProblemRequest{}
	mi := &file_admin_v1_problem_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProblemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProblemRequest) ProtoMessage() {}

func (x *GetProblemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProblemRequest.ProtoReflect.Descriptor instead.
func (*GetProblemRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{6}
}

func (x *GetProblemRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type GetProblemResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Problem       *Problem               `protobuf:"bytes,1,opt,name=problem,proto3" json:"problem,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetProblemResponse) Reset() {
	*x = GetProblemResponse{}
	mi := &file_admin_v1_problem_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProblemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProblemResponse) ProtoMessage() {}

func (x *GetProblemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProblemResponse.ProtoReflect.Descriptor instead.
func (*GetProblemResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{7}
}

func (x *GetProblemResponse) GetProblem() *Problem {
	if x != nil {
		return x.Problem
	}
	return nil
}

type CreateProblemRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Problem       *Problem               `protobuf:"bytes,1,opt,name=problem,proto3" json:"problem,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateProblemRequest) Reset() {
	*x = CreateProblemRequest{}
	mi := &file_admin_v1_problem_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateProblemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateProblemRequest) ProtoMessage() {}

func (x *CreateProblemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateProblemRequest.ProtoReflect.Descriptor instead.
func (*CreateProblemRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{8}
}

func (x *CreateProblemRequest) GetProblem() *Problem {
	if x != nil {
		return x.Problem
	}
	return nil
}

type CreateProblemResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Problem       *Problem               `protobuf:"bytes,1,opt,name=problem,proto3" json:"problem,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateProblemResponse) Reset() {
	*x = CreateProblemResponse{}
	mi := &file_admin_v1_problem_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateProblemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateProblemResponse) ProtoMessage() {}

func (x *CreateProblemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateProblemResponse.ProtoReflect.Descriptor instead.
func (*CreateProblemResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{9}
}

func (x *CreateProblemResponse) GetProblem() *Problem {
	if x != nil {
		return x.Problem
	}
	return nil
}

type UpdateProblemRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Problem       *Problem               `protobuf:"bytes,1,opt,name=problem,proto3" json:"problem,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateProblemRequest) Reset() {
	*x = UpdateProblemRequest{}
	mi := &file_admin_v1_problem_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateProblemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProblemRequest) ProtoMessage() {}

func (x *UpdateProblemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProblemRequest.ProtoReflect.Descriptor instead.
func (*UpdateProblemRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{10}
}

func (x *UpdateProblemRequest) GetProblem() *Problem {
	if x != nil {
		return x.Problem
	}
	return nil
}

type UpdateProblemResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Problem       *Problem               `protobuf:"bytes,1,opt,name=problem,proto3" json:"problem,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateProblemResponse) Reset() {
	*x = UpdateProblemResponse{}
	mi := &file_admin_v1_problem_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateProblemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProblemResponse) ProtoMessage() {}

func (x *UpdateProblemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProblemResponse.ProtoReflect.Descriptor instead.
func (*UpdateProblemResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{11}
}

func (x *UpdateProblemResponse) GetProblem() *Problem {
	if x != nil {
		return x.Problem
	}
	return nil
}

type DeleteProblemRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteProblemRequest) Reset() {
	*x = DeleteProblemRequest{}
	mi := &file_admin_v1_problem_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteProblemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteProblemRequest) ProtoMessage() {}

func (x *DeleteProblemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteProblemRequest.ProtoReflect.Descriptor instead.
func (*DeleteProblemRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{12}
}

func (x *DeleteProblemRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type DeleteProblemResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteProblemResponse) Reset() {
	*x = DeleteProblemResponse{}
	mi := &file_admin_v1_problem_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteProblemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteProblemResponse) ProtoMessage() {}

func (x *DeleteProblemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_problem_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteProblemResponse.ProtoReflect.Descriptor instead.
func (*DeleteProblemResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_problem_proto_rawDescGZIP(), []int{13}
}

var File_admin_v1_problem_proto protoreflect.FileDescriptor

const file_admin_v1_problem_proto_rawDesc = "" +
	"\n" +
	"\x16admin/v1/problem.proto\x12\badmin.v1\x1a\x1bbuf/validate/validate.proto\"\xdd\x01\n" +
	"\aProblem\x12\x1b\n" +
	"\x04code\x18\x01 \x01(\tB\a\xbaH\x04r\x02\x10\x01R\x04code\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12\x1b\n" +
	"\tmax_score\x18\x03 \x01(\rR\bmaxScore\x12\x1a\n" +
	"\bcategory\x18\x04 \x01(\tR\bcategory\x12;\n" +
	"\rredeploy_rule\x18\x05 \x01(\v2\x16.admin.v1.RedeployRuleR\fredeployRule\x12)\n" +
	"\x04body\x18\x06 \x01(\v2\x15.admin.v1.ProblemBodyR\x04body\"\x9a\x01\n" +
	"\fRedeployRule\x12.\n" +
	"\x04type\x18\x01 \x01(\x0e2\x1a.admin.v1.RedeployRuleTypeR\x04type\x12+\n" +
	"\x11penalty_threshold\x18\x02 \x01(\rR\x10penaltyThreshold\x12-\n" +
	"\x12penalty_percentage\x18\x03 \x01(\rR\x11penaltyPercentage\"\x82\x01\n" +
	"\vProblemBody\x12)\n" +
	"\x04type\x18\x01 \x01(\x0e2\x15.admin.v1.ProblemTypeR\x04type\x12@\n" +
	"\vdescriptive\x18\x02 \x01(\v2\x1c.admin.v1.DescriptiveProblemH\x00R\vdescriptiveB\x06\n" +
	"\x04body\"r\n" +
	"\x12DescriptiveProblem\x12)\n" +
	"\x10problem_markdown\x18\x03 \x01(\tR\x0fproblemMarkdown\x121\n" +
	"\x14explanation_markdown\x18\x04 \x01(\tR\x13explanationMarkdown\"\x15\n" +
	"\x13ListProblemsRequest\"E\n" +
	"\x14ListProblemsResponse\x12-\n" +
	"\bproblems\x18\x01 \x03(\v2\x11.admin.v1.ProblemR\bproblems\"'\n" +
	"\x11GetProblemRequest\x12\x12\n" +
	"\x04code\x18\x01 \x01(\tR\x04code\"A\n" +
	"\x12GetProblemResponse\x12+\n" +
	"\aproblem\x18\x01 \x01(\v2\x11.admin.v1.ProblemR\aproblem\"C\n" +
	"\x14CreateProblemRequest\x12+\n" +
	"\aproblem\x18\x01 \x01(\v2\x11.admin.v1.ProblemR\aproblem\"D\n" +
	"\x15CreateProblemResponse\x12+\n" +
	"\aproblem\x18\x01 \x01(\v2\x11.admin.v1.ProblemR\aproblem\"C\n" +
	"\x14UpdateProblemRequest\x12+\n" +
	"\aproblem\x18\x01 \x01(\v2\x11.admin.v1.ProblemR\aproblem\"D\n" +
	"\x15UpdateProblemResponse\x12+\n" +
	"\aproblem\x18\x01 \x01(\v2\x11.admin.v1.ProblemR\aproblem\"*\n" +
	"\x14DeleteProblemRequest\x12\x12\n" +
	"\x04code\x18\x01 \x01(\tR\x04code\"\x17\n" +
	"\x15DeleteProblemResponse*\x88\x01\n" +
	"\x10RedeployRuleType\x12\"\n" +
	"\x1eREDEPLOY_RULE_TYPE_UNSPECIFIED\x10\x00\x12%\n" +
	"!REDEPLOY_RULE_TYPE_UNREDEPLOYABLE\x10\x01\x12)\n" +
	"%REDEPLOY_RULE_TYPE_PERCENTAGE_PENALTY\x10\x02*I\n" +
	"\vProblemType\x12\x1c\n" +
	"\x18PROBLEM_TYPE_UNSPECIFIED\x10\x00\x12\x1c\n" +
	"\x18PROBLEM_TYPE_DESCRIPTIVE\x10\x012\x9e\x03\n" +
	"\x0eProblemService\x12M\n" +
	"\fListProblems\x12\x1d.admin.v1.ListProblemsRequest\x1a\x1e.admin.v1.ListProblemsResponse\x12G\n" +
	"\n" +
	"GetProblem\x12\x1b.admin.v1.GetProblemRequest\x1a\x1c.admin.v1.GetProblemResponse\x12P\n" +
	"\rCreateProblem\x12\x1e.admin.v1.CreateProblemRequest\x1a\x1f.admin.v1.CreateProblemResponse\x12P\n" +
	"\rUpdateProblem\x12\x1e.admin.v1.UpdateProblemRequest\x1a\x1f.admin.v1.UpdateProblemResponse\x12P\n" +
	"\rDeleteProblem\x12\x1e.admin.v1.DeleteProblemRequest\x1a\x1f.admin.v1.DeleteProblemResponseBCZAgithub.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1;adminv1b\x06proto3"

var (
	file_admin_v1_problem_proto_rawDescOnce sync.Once
	file_admin_v1_problem_proto_rawDescData []byte
)

func file_admin_v1_problem_proto_rawDescGZIP() []byte {
	file_admin_v1_problem_proto_rawDescOnce.Do(func() {
		file_admin_v1_problem_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_admin_v1_problem_proto_rawDesc), len(file_admin_v1_problem_proto_rawDesc)))
	})
	return file_admin_v1_problem_proto_rawDescData
}

var file_admin_v1_problem_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_admin_v1_problem_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_admin_v1_problem_proto_goTypes = []any{
	(RedeployRuleType)(0),         // 0: admin.v1.RedeployRuleType
	(ProblemType)(0),              // 1: admin.v1.ProblemType
	(*Problem)(nil),               // 2: admin.v1.Problem
	(*RedeployRule)(nil),          // 3: admin.v1.RedeployRule
	(*ProblemBody)(nil),           // 4: admin.v1.ProblemBody
	(*DescriptiveProblem)(nil),    // 5: admin.v1.DescriptiveProblem
	(*ListProblemsRequest)(nil),   // 6: admin.v1.ListProblemsRequest
	(*ListProblemsResponse)(nil),  // 7: admin.v1.ListProblemsResponse
	(*GetProblemRequest)(nil),     // 8: admin.v1.GetProblemRequest
	(*GetProblemResponse)(nil),    // 9: admin.v1.GetProblemResponse
	(*CreateProblemRequest)(nil),  // 10: admin.v1.CreateProblemRequest
	(*CreateProblemResponse)(nil), // 11: admin.v1.CreateProblemResponse
	(*UpdateProblemRequest)(nil),  // 12: admin.v1.UpdateProblemRequest
	(*UpdateProblemResponse)(nil), // 13: admin.v1.UpdateProblemResponse
	(*DeleteProblemRequest)(nil),  // 14: admin.v1.DeleteProblemRequest
	(*DeleteProblemResponse)(nil), // 15: admin.v1.DeleteProblemResponse
}
var file_admin_v1_problem_proto_depIdxs = []int32{
	3,  // 0: admin.v1.Problem.redeploy_rule:type_name -> admin.v1.RedeployRule
	4,  // 1: admin.v1.Problem.body:type_name -> admin.v1.ProblemBody
	0,  // 2: admin.v1.RedeployRule.type:type_name -> admin.v1.RedeployRuleType
	1,  // 3: admin.v1.ProblemBody.type:type_name -> admin.v1.ProblemType
	5,  // 4: admin.v1.ProblemBody.descriptive:type_name -> admin.v1.DescriptiveProblem
	2,  // 5: admin.v1.ListProblemsResponse.problems:type_name -> admin.v1.Problem
	2,  // 6: admin.v1.GetProblemResponse.problem:type_name -> admin.v1.Problem
	2,  // 7: admin.v1.CreateProblemRequest.problem:type_name -> admin.v1.Problem
	2,  // 8: admin.v1.CreateProblemResponse.problem:type_name -> admin.v1.Problem
	2,  // 9: admin.v1.UpdateProblemRequest.problem:type_name -> admin.v1.Problem
	2,  // 10: admin.v1.UpdateProblemResponse.problem:type_name -> admin.v1.Problem
	6,  // 11: admin.v1.ProblemService.ListProblems:input_type -> admin.v1.ListProblemsRequest
	8,  // 12: admin.v1.ProblemService.GetProblem:input_type -> admin.v1.GetProblemRequest
	10, // 13: admin.v1.ProblemService.CreateProblem:input_type -> admin.v1.CreateProblemRequest
	12, // 14: admin.v1.ProblemService.UpdateProblem:input_type -> admin.v1.UpdateProblemRequest
	14, // 15: admin.v1.ProblemService.DeleteProblem:input_type -> admin.v1.DeleteProblemRequest
	7,  // 16: admin.v1.ProblemService.ListProblems:output_type -> admin.v1.ListProblemsResponse
	9,  // 17: admin.v1.ProblemService.GetProblem:output_type -> admin.v1.GetProblemResponse
	11, // 18: admin.v1.ProblemService.CreateProblem:output_type -> admin.v1.CreateProblemResponse
	13, // 19: admin.v1.ProblemService.UpdateProblem:output_type -> admin.v1.UpdateProblemResponse
	15, // 20: admin.v1.ProblemService.DeleteProblem:output_type -> admin.v1.DeleteProblemResponse
	16, // [16:21] is the sub-list for method output_type
	11, // [11:16] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_admin_v1_problem_proto_init() }
func file_admin_v1_problem_proto_init() {
	if File_admin_v1_problem_proto != nil {
		return
	}
	file_admin_v1_problem_proto_msgTypes[2].OneofWrappers = []any{
		(*ProblemBody_Descriptive)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_admin_v1_problem_proto_rawDesc), len(file_admin_v1_problem_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_problem_proto_goTypes,
		DependencyIndexes: file_admin_v1_problem_proto_depIdxs,
		EnumInfos:         file_admin_v1_problem_proto_enumTypes,
		MessageInfos:      file_admin_v1_problem_proto_msgTypes,
	}.Build()
	File_admin_v1_problem_proto = out.File
	file_admin_v1_problem_proto_goTypes = nil
	file_admin_v1_problem_proto_depIdxs = nil
}
