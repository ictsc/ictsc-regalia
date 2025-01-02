// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: contestant/v1/problem.proto

package contestantv1

import (
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

type DeploymentStatus int32

const (
	DeploymentStatus_DEPLOYMENT_STATUS_UNSPECIFIED DeploymentStatus = 0
	// 展開済み
	DeploymentStatus_DEPLOYMENT_STATUS_DEPLOYED DeploymentStatus = 1
	// 展開中
	DeploymentStatus_DEPLOYMENT_STATUS_DEPLOYING DeploymentStatus = 2
	// 展開失敗
	DeploymentStatus_DEPLOYMENT_STATUS_FAILED DeploymentStatus = 3
)

// Enum value maps for DeploymentStatus.
var (
	DeploymentStatus_name = map[int32]string{
		0: "DEPLOYMENT_STATUS_UNSPECIFIED",
		1: "DEPLOYMENT_STATUS_DEPLOYED",
		2: "DEPLOYMENT_STATUS_DEPLOYING",
		3: "DEPLOYMENT_STATUS_FAILED",
	}
	DeploymentStatus_value = map[string]int32{
		"DEPLOYMENT_STATUS_UNSPECIFIED": 0,
		"DEPLOYMENT_STATUS_DEPLOYED":    1,
		"DEPLOYMENT_STATUS_DEPLOYING":   2,
		"DEPLOYMENT_STATUS_FAILED":      3,
	}
)

func (x DeploymentStatus) Enum() *DeploymentStatus {
	p := new(DeploymentStatus)
	*p = x
	return p
}

func (x DeploymentStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DeploymentStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_contestant_v1_problem_proto_enumTypes[0].Descriptor()
}

func (DeploymentStatus) Type() protoreflect.EnumType {
	return &file_contestant_v1_problem_proto_enumTypes[0]
}

func (x DeploymentStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DeploymentStatus.Descriptor instead.
func (DeploymentStatus) EnumDescriptor() ([]byte, []int) {
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{0}
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
	return file_contestant_v1_problem_proto_enumTypes[1].Descriptor()
}

func (ProblemType) Type() protoreflect.EnumType {
	return &file_contestant_v1_problem_proto_enumTypes[1]
}

func (x ProblemType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProblemType.Descriptor instead.
func (ProblemType) EnumDescriptor() ([]byte, []int) {
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{1}
}

type Problem struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 問題コード
	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	// タイトル
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	// 最大得点
	MaxScore      int64        `protobuf:"varint,3,opt,name=max_score,json=maxScore,proto3" json:"max_score,omitempty"`
	Score         *Score       `protobuf:"bytes,4,opt,name=score,proto3,oneof" json:"score,omitempty"`
	Deployment    *Deployment  `protobuf:"bytes,5,opt,name=deployment,proto3" json:"deployment,omitempty"`
	Body          *ProblemBody `protobuf:"bytes,6,opt,name=body,proto3" json:"body,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Problem) Reset() {
	*x = Problem{}
	mi := &file_contestant_v1_problem_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Problem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Problem) ProtoMessage() {}

func (x *Problem) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[0]
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
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{0}
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

func (x *Problem) GetMaxScore() int64 {
	if x != nil {
		return x.MaxScore
	}
	return 0
}

func (x *Problem) GetScore() *Score {
	if x != nil {
		return x.Score
	}
	return nil
}

func (x *Problem) GetDeployment() *Deployment {
	if x != nil {
		return x.Deployment
	}
	return nil
}

func (x *Problem) GetBody() *ProblemBody {
	if x != nil {
		return x.Body
	}
	return nil
}

type Score struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 採点による得点
	MarkedScore int64 `protobuf:"varint,1,opt,name=marked_score,json=markedScore,proto3" json:"marked_score,omitempty"`
	// ペナルティによる減点
	Penalty int64 `protobuf:"varint,2,opt,name=penalty,proto3" json:"penalty,omitempty"`
	// 最終的な得点
	Score         int64 `protobuf:"varint,3,opt,name=score,proto3" json:"score,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Score) Reset() {
	*x = Score{}
	mi := &file_contestant_v1_problem_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Score) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Score) ProtoMessage() {}

func (x *Score) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Score.ProtoReflect.Descriptor instead.
func (*Score) Descriptor() ([]byte, []int) {
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{1}
}

func (x *Score) GetMarkedScore() int64 {
	if x != nil {
		return x.MarkedScore
	}
	return 0
}

func (x *Score) GetPenalty() int64 {
	if x != nil {
		return x.Penalty
	}
	return 0
}

func (x *Score) GetScore() int64 {
	if x != nil {
		return x.Score
	}
	return 0
}

type Deployment struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        DeploymentStatus       `protobuf:"varint,1,opt,name=status,proto3,enum=contestant.v1.DeploymentStatus" json:"status,omitempty"`
	Redeployable  bool                   `protobuf:"varint,2,opt,name=redeployable,proto3" json:"redeployable,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Deployment) Reset() {
	*x = Deployment{}
	mi := &file_contestant_v1_problem_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Deployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deployment) ProtoMessage() {}

func (x *Deployment) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deployment.ProtoReflect.Descriptor instead.
func (*Deployment) Descriptor() ([]byte, []int) {
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{2}
}

func (x *Deployment) GetStatus() DeploymentStatus {
	if x != nil {
		return x.Status
	}
	return DeploymentStatus_DEPLOYMENT_STATUS_UNSPECIFIED
}

func (x *Deployment) GetRedeployable() bool {
	if x != nil {
		return x.Redeployable
	}
	return false
}

type ProblemBody struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Type  ProblemType            `protobuf:"varint,1,opt,name=type,proto3,enum=contestant.v1.ProblemType" json:"type,omitempty"`
	// Types that are valid to be assigned to Body:
	//
	//	*ProblemBody_Descriptive
	Body          isProblemBody_Body `protobuf_oneof:"body"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProblemBody) Reset() {
	*x = ProblemBody{}
	mi := &file_contestant_v1_problem_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProblemBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProblemBody) ProtoMessage() {}

func (x *ProblemBody) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[3]
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
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{3}
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

type Connection struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ホスト名
	HostName string `protobuf:"bytes,1,opt,name=host_name,json=hostName,proto3" json:"host_name,omitempty"`
	// ホスト(IP アドレス or ドメイン)
	Host string `protobuf:"bytes,2,opt,name=host,proto3" json:"host,omitempty"`
	// ユーザ
	User *string `protobuf:"bytes,3,opt,name=user,proto3,oneof" json:"user,omitempty"`
	// パスワード
	Password      *string `protobuf:"bytes,4,opt,name=password,proto3,oneof" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Connection) Reset() {
	*x = Connection{}
	mi := &file_contestant_v1_problem_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Connection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Connection) ProtoMessage() {}

func (x *Connection) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Connection.ProtoReflect.Descriptor instead.
func (*Connection) Descriptor() ([]byte, []int) {
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{4}
}

func (x *Connection) GetHostName() string {
	if x != nil {
		return x.HostName
	}
	return ""
}

func (x *Connection) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Connection) GetUser() string {
	if x != nil && x.User != nil {
		return *x.User
	}
	return ""
}

func (x *Connection) GetPassword() string {
	if x != nil && x.Password != nil {
		return *x.Password
	}
	return ""
}

type DescriptiveProblem struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 問題文
	Body string `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
	// 接続情報
	Connection    []*Connection `protobuf:"bytes,2,rep,name=connection,proto3" json:"connection,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DescriptiveProblem) Reset() {
	*x = DescriptiveProblem{}
	mi := &file_contestant_v1_problem_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DescriptiveProblem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescriptiveProblem) ProtoMessage() {}

func (x *DescriptiveProblem) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[5]
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
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{5}
}

func (x *DescriptiveProblem) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *DescriptiveProblem) GetConnection() []*Connection {
	if x != nil {
		return x.Connection
	}
	return nil
}

type ListProblemsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListProblemsRequest) Reset() {
	*x = ListProblemsRequest{}
	mi := &file_contestant_v1_problem_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListProblemsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProblemsRequest) ProtoMessage() {}

func (x *ListProblemsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[6]
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
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{6}
}

type ListProblemsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Problems      []*Problem             `protobuf:"bytes,1,rep,name=problems,proto3" json:"problems,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListProblemsResponse) Reset() {
	*x = ListProblemsResponse{}
	mi := &file_contestant_v1_problem_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListProblemsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProblemsResponse) ProtoMessage() {}

func (x *ListProblemsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[7]
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
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{7}
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
	mi := &file_contestant_v1_problem_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProblemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProblemRequest) ProtoMessage() {}

func (x *GetProblemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[8]
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
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{8}
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
	mi := &file_contestant_v1_problem_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProblemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProblemResponse) ProtoMessage() {}

func (x *GetProblemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[9]
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
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{9}
}

func (x *GetProblemResponse) GetProblem() *Problem {
	if x != nil {
		return x.Problem
	}
	return nil
}

type DeployRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeployRequest) Reset() {
	*x = DeployRequest{}
	mi := &file_contestant_v1_problem_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeployRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeployRequest) ProtoMessage() {}

func (x *DeployRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeployRequest.ProtoReflect.Descriptor instead.
func (*DeployRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{10}
}

func (x *DeployRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type DeployResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeployResponse) Reset() {
	*x = DeployResponse{}
	mi := &file_contestant_v1_problem_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeployResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeployResponse) ProtoMessage() {}

func (x *DeployResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_problem_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeployResponse.ProtoReflect.Descriptor instead.
func (*DeployResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_problem_proto_rawDescGZIP(), []int{11}
}

var File_contestant_v1_problem_proto protoreflect.FileDescriptor

var file_contestant_v1_problem_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x22, 0xf6, 0x01, 0x0a,
	0x07, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6d, 0x61, 0x78, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12,
	0x2f, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x63, 0x6f, 0x72, 0x65, 0x48, 0x00, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x88, 0x01, 0x01,
	0x12, 0x39, 0x0a, 0x0a, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x0a, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65,
	0x6d, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x42, 0x08, 0x0a, 0x06, 0x5f,
	0x73, 0x63, 0x6f, 0x72, 0x65, 0x22, 0x5a, 0x0a, 0x05, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x21,
	0x0a, 0x0c, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x64, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x64, 0x53, 0x63, 0x6f, 0x72,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x70, 0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x73,
	0x63, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72,
	0x65, 0x22, 0x69, 0x0a, 0x0a, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x37, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x64, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c,
	0x72, 0x65, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x8c, 0x01, 0x0a,
	0x0b, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x2e, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c,
	0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x45, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x76, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x21, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x76, 0x65, 0x50, 0x72, 0x6f,
	0x62, 0x6c, 0x65, 0x6d, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x76, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x8d, 0x01, 0x0a, 0x0a,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6f,
	0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68,
	0x6f, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x04, 0x75,
	0x73, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x42, 0x0b,
	0x0a, 0x09, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x63, 0x0a, 0x12, 0x44,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x76, 0x65, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65,
	0x6d, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x15, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4a, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x50,
	0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x32, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x62, 0x6c,
	0x65, 0x6d, 0x73, 0x22, 0x27, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x46, 0x0a, 0x12,
	0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x30, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x52, 0x07, 0x70, 0x72, 0x6f,
	0x62, 0x6c, 0x65, 0x6d, 0x22, 0x23, 0x0a, 0x0d, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x44, 0x65, 0x70,
	0x6c, 0x6f, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x94, 0x01, 0x0a, 0x10,
	0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x21, 0x0a, 0x1d, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x1e, 0x0a, 0x1a, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x45,
	0x44, 0x10, 0x01, 0x12, 0x1f, 0x0a, 0x1b, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x49,
	0x4e, 0x47, 0x10, 0x02, 0x12, 0x1c, 0x0a, 0x18, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x4d, 0x45,
	0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44,
	0x10, 0x03, 0x2a, 0x49, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x1c, 0x0a, 0x18, 0x50, 0x52, 0x4f, 0x42, 0x4c, 0x45, 0x4d, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x1c, 0x0a, 0x18, 0x50, 0x52, 0x4f, 0x42, 0x4c, 0x45, 0x4d, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x44, 0x45, 0x53, 0x43, 0x52, 0x49, 0x50, 0x54, 0x49, 0x56, 0x45, 0x10, 0x01, 0x32, 0x83, 0x02,
	0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x57, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x73,
	0x12, 0x22, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x0a, 0x47, 0x65, 0x74,
	0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x12, 0x20, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73,
	0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x62, 0x6c,
	0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f,
	0x62, 0x6c, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x06,
	0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x12, 0x1c, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x72, 0x65,
	0x67, 0x61, 0x6c, 0x69, 0x61, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61,
	0x6e, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_contestant_v1_problem_proto_rawDescOnce sync.Once
	file_contestant_v1_problem_proto_rawDescData = file_contestant_v1_problem_proto_rawDesc
)

func file_contestant_v1_problem_proto_rawDescGZIP() []byte {
	file_contestant_v1_problem_proto_rawDescOnce.Do(func() {
		file_contestant_v1_problem_proto_rawDescData = protoimpl.X.CompressGZIP(file_contestant_v1_problem_proto_rawDescData)
	})
	return file_contestant_v1_problem_proto_rawDescData
}

var file_contestant_v1_problem_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_contestant_v1_problem_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_contestant_v1_problem_proto_goTypes = []any{
	(DeploymentStatus)(0),        // 0: contestant.v1.DeploymentStatus
	(ProblemType)(0),             // 1: contestant.v1.ProblemType
	(*Problem)(nil),              // 2: contestant.v1.Problem
	(*Score)(nil),                // 3: contestant.v1.Score
	(*Deployment)(nil),           // 4: contestant.v1.Deployment
	(*ProblemBody)(nil),          // 5: contestant.v1.ProblemBody
	(*Connection)(nil),           // 6: contestant.v1.Connection
	(*DescriptiveProblem)(nil),   // 7: contestant.v1.DescriptiveProblem
	(*ListProblemsRequest)(nil),  // 8: contestant.v1.ListProblemsRequest
	(*ListProblemsResponse)(nil), // 9: contestant.v1.ListProblemsResponse
	(*GetProblemRequest)(nil),    // 10: contestant.v1.GetProblemRequest
	(*GetProblemResponse)(nil),   // 11: contestant.v1.GetProblemResponse
	(*DeployRequest)(nil),        // 12: contestant.v1.DeployRequest
	(*DeployResponse)(nil),       // 13: contestant.v1.DeployResponse
}
var file_contestant_v1_problem_proto_depIdxs = []int32{
	3,  // 0: contestant.v1.Problem.score:type_name -> contestant.v1.Score
	4,  // 1: contestant.v1.Problem.deployment:type_name -> contestant.v1.Deployment
	5,  // 2: contestant.v1.Problem.body:type_name -> contestant.v1.ProblemBody
	0,  // 3: contestant.v1.Deployment.status:type_name -> contestant.v1.DeploymentStatus
	1,  // 4: contestant.v1.ProblemBody.type:type_name -> contestant.v1.ProblemType
	7,  // 5: contestant.v1.ProblemBody.descriptive:type_name -> contestant.v1.DescriptiveProblem
	6,  // 6: contestant.v1.DescriptiveProblem.connection:type_name -> contestant.v1.Connection
	2,  // 7: contestant.v1.ListProblemsResponse.problems:type_name -> contestant.v1.Problem
	2,  // 8: contestant.v1.GetProblemResponse.problem:type_name -> contestant.v1.Problem
	8,  // 9: contestant.v1.ProblemService.ListProblems:input_type -> contestant.v1.ListProblemsRequest
	10, // 10: contestant.v1.ProblemService.GetProblem:input_type -> contestant.v1.GetProblemRequest
	12, // 11: contestant.v1.ProblemService.Deploy:input_type -> contestant.v1.DeployRequest
	9,  // 12: contestant.v1.ProblemService.ListProblems:output_type -> contestant.v1.ListProblemsResponse
	11, // 13: contestant.v1.ProblemService.GetProblem:output_type -> contestant.v1.GetProblemResponse
	13, // 14: contestant.v1.ProblemService.Deploy:output_type -> contestant.v1.DeployResponse
	12, // [12:15] is the sub-list for method output_type
	9,  // [9:12] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_contestant_v1_problem_proto_init() }
func file_contestant_v1_problem_proto_init() {
	if File_contestant_v1_problem_proto != nil {
		return
	}
	file_contestant_v1_problem_proto_msgTypes[0].OneofWrappers = []any{}
	file_contestant_v1_problem_proto_msgTypes[3].OneofWrappers = []any{
		(*ProblemBody_Descriptive)(nil),
	}
	file_contestant_v1_problem_proto_msgTypes[4].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_contestant_v1_problem_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contestant_v1_problem_proto_goTypes,
		DependencyIndexes: file_contestant_v1_problem_proto_depIdxs,
		EnumInfos:         file_contestant_v1_problem_proto_enumTypes,
		MessageInfos:      file_contestant_v1_problem_proto_msgTypes,
	}.Build()
	File_contestant_v1_problem_proto = out.File
	file_contestant_v1_problem_proto_rawDesc = nil
	file_contestant_v1_problem_proto_goTypes = nil
	file_contestant_v1_problem_proto_depIdxs = nil
}
