// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: contestant/v1/answer.proto

package contestantv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

// 解答
type Answer struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Body          *AnswerBody            `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	SubmittedAt   *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=submitted_at,json=submittedAt,proto3" json:"submitted_at,omitempty"`
	Score         *Score                 `protobuf:"bytes,4,opt,name=score,proto3,oneof" json:"score,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Answer) Reset() {
	*x = Answer{}
	mi := &file_contestant_v1_answer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Answer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Answer) ProtoMessage() {}

func (x *Answer) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_answer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Answer.ProtoReflect.Descriptor instead.
func (*Answer) Descriptor() ([]byte, []int) {
	return file_contestant_v1_answer_proto_rawDescGZIP(), []int{0}
}

func (x *Answer) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Answer) GetBody() *AnswerBody {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *Answer) GetSubmittedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.SubmittedAt
	}
	return nil
}

func (x *Answer) GetScore() *Score {
	if x != nil {
		return x.Score
	}
	return nil
}

type AnswerBody struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Type  ProblemType            `protobuf:"varint,1,opt,name=type,proto3,enum=contestant.v1.ProblemType" json:"type,omitempty"`
	// Types that are valid to be assigned to Body:
	//
	//	*AnswerBody_Descriptive
	Body          isAnswerBody_Body `protobuf_oneof:"body"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AnswerBody) Reset() {
	*x = AnswerBody{}
	mi := &file_contestant_v1_answer_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AnswerBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnswerBody) ProtoMessage() {}

func (x *AnswerBody) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_answer_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnswerBody.ProtoReflect.Descriptor instead.
func (*AnswerBody) Descriptor() ([]byte, []int) {
	return file_contestant_v1_answer_proto_rawDescGZIP(), []int{1}
}

func (x *AnswerBody) GetType() ProblemType {
	if x != nil {
		return x.Type
	}
	return ProblemType_PROBLEM_TYPE_UNSPECIFIED
}

func (x *AnswerBody) GetBody() isAnswerBody_Body {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *AnswerBody) GetDescriptive() *DescriptiveAnswer {
	if x != nil {
		if x, ok := x.Body.(*AnswerBody_Descriptive); ok {
			return x.Descriptive
		}
	}
	return nil
}

type isAnswerBody_Body interface {
	isAnswerBody_Body()
}

type AnswerBody_Descriptive struct {
	Descriptive *DescriptiveAnswer `protobuf:"bytes,2,opt,name=descriptive,proto3,oneof"`
}

func (*AnswerBody_Descriptive) isAnswerBody_Body() {}

type DescriptiveAnswer struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Body          string                 `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DescriptiveAnswer) Reset() {
	*x = DescriptiveAnswer{}
	mi := &file_contestant_v1_answer_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DescriptiveAnswer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescriptiveAnswer) ProtoMessage() {}

func (x *DescriptiveAnswer) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_answer_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescriptiveAnswer.ProtoReflect.Descriptor instead.
func (*DescriptiveAnswer) Descriptor() ([]byte, []int) {
	return file_contestant_v1_answer_proto_rawDescGZIP(), []int{2}
}

func (x *DescriptiveAnswer) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type ListAnswersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProblemCode   string                 `protobuf:"bytes,1,opt,name=problem_code,json=problemCode,proto3" json:"problem_code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListAnswersRequest) Reset() {
	*x = ListAnswersRequest{}
	mi := &file_contestant_v1_answer_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListAnswersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAnswersRequest) ProtoMessage() {}

func (x *ListAnswersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_answer_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAnswersRequest.ProtoReflect.Descriptor instead.
func (*ListAnswersRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_answer_proto_rawDescGZIP(), []int{3}
}

func (x *ListAnswersRequest) GetProblemCode() string {
	if x != nil {
		return x.ProblemCode
	}
	return ""
}

type ListAnswersResponse struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Answers         []*Answer              `protobuf:"bytes,1,rep,name=answers,proto3" json:"answers,omitempty"`
	SubmitInterval  *durationpb.Duration   `protobuf:"bytes,2,opt,name=submit_interval,json=submitInterval,proto3" json:"submit_interval,omitempty"`
	LastSubmittedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=last_submitted_at,json=lastSubmittedAt,proto3,oneof" json:"last_submitted_at,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *ListAnswersResponse) Reset() {
	*x = ListAnswersResponse{}
	mi := &file_contestant_v1_answer_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListAnswersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAnswersResponse) ProtoMessage() {}

func (x *ListAnswersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_answer_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAnswersResponse.ProtoReflect.Descriptor instead.
func (*ListAnswersResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_answer_proto_rawDescGZIP(), []int{4}
}

func (x *ListAnswersResponse) GetAnswers() []*Answer {
	if x != nil {
		return x.Answers
	}
	return nil
}

func (x *ListAnswersResponse) GetSubmitInterval() *durationpb.Duration {
	if x != nil {
		return x.SubmitInterval
	}
	return nil
}

func (x *ListAnswersResponse) GetLastSubmittedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.LastSubmittedAt
	}
	return nil
}

type SubmitAnswerRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProblemCode   string                 `protobuf:"bytes,1,opt,name=problem_code,json=problemCode,proto3" json:"problem_code,omitempty"`
	Body          string                 `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitAnswerRequest) Reset() {
	*x = SubmitAnswerRequest{}
	mi := &file_contestant_v1_answer_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitAnswerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitAnswerRequest) ProtoMessage() {}

func (x *SubmitAnswerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_answer_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitAnswerRequest.ProtoReflect.Descriptor instead.
func (*SubmitAnswerRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_answer_proto_rawDescGZIP(), []int{5}
}

func (x *SubmitAnswerRequest) GetProblemCode() string {
	if x != nil {
		return x.ProblemCode
	}
	return ""
}

func (x *SubmitAnswerRequest) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type SubmitAnswerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Answer        *Answer                `protobuf:"bytes,1,opt,name=answer,proto3" json:"answer,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitAnswerResponse) Reset() {
	*x = SubmitAnswerResponse{}
	mi := &file_contestant_v1_answer_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitAnswerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitAnswerResponse) ProtoMessage() {}

func (x *SubmitAnswerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_answer_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitAnswerResponse.ProtoReflect.Descriptor instead.
func (*SubmitAnswerResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_answer_proto_rawDescGZIP(), []int{6}
}

func (x *SubmitAnswerResponse) GetAnswer() *Answer {
	if x != nil {
		return x.Answer
	}
	return nil
}

type GetAnswerRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProblemCode   string                 `protobuf:"bytes,1,opt,name=problem_code,json=problemCode,proto3" json:"problem_code,omitempty"`
	Id            uint32                 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAnswerRequest) Reset() {
	*x = GetAnswerRequest{}
	mi := &file_contestant_v1_answer_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAnswerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAnswerRequest) ProtoMessage() {}

func (x *GetAnswerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_answer_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAnswerRequest.ProtoReflect.Descriptor instead.
func (*GetAnswerRequest) Descriptor() ([]byte, []int) {
	return file_contestant_v1_answer_proto_rawDescGZIP(), []int{7}
}

func (x *GetAnswerRequest) GetProblemCode() string {
	if x != nil {
		return x.ProblemCode
	}
	return ""
}

func (x *GetAnswerRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetAnswerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Answer        *Answer                `protobuf:"bytes,1,opt,name=answer,proto3" json:"answer,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAnswerResponse) Reset() {
	*x = GetAnswerResponse{}
	mi := &file_contestant_v1_answer_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAnswerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAnswerResponse) ProtoMessage() {}

func (x *GetAnswerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contestant_v1_answer_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAnswerResponse.ProtoReflect.Descriptor instead.
func (*GetAnswerResponse) Descriptor() ([]byte, []int) {
	return file_contestant_v1_answer_proto_rawDescGZIP(), []int{8}
}

func (x *GetAnswerResponse) GetAnswer() *Answer {
	if x != nil {
		return x.Answer
	}
	return nil
}

var File_contestant_v1_answer_proto protoreflect.FileDescriptor

var file_contestant_v1_answer_proto_rawDesc = string([]byte{
	0x0a, 0x1a, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x62, 0x6c,
	0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc1, 0x01, 0x0a, 0x06, 0x41, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x2d, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x12, 0x3d, 0x0a, 0x0c, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x2f, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x48, 0x00, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x22, 0x8a, 0x01,
	0x0a, 0x0a, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x2e, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c,
	0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x44, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x76, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x20, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x76, 0x65, 0x41, 0x6e, 0x73,
	0x77, 0x65, 0x72, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x76, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x27, 0x0a, 0x11, 0x44, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x76, 0x65, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12,
	0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x22, 0x37, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65,
	0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f,
	0x62, 0x6c, 0x65, 0x6d, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x43, 0x6f, 0x64, 0x65, 0x22, 0xed, 0x01, 0x0a,
	0x13, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x07, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x07, 0x61, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x73, 0x12, 0x42, 0x0a, 0x0f, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x5f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0e, 0x73, 0x75, 0x62, 0x6d, 0x69,
	0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x4b, 0x0a, 0x11, 0x6c, 0x61, 0x73,
	0x74, 0x5f, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x48, 0x00, 0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x42, 0x14, 0x0a, 0x12, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f,
	0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x22, 0x4c, 0x0a, 0x13,
	0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x5f, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x62, 0x6c,
	0x65, 0x6d, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x45, 0x0a, 0x14, 0x53, 0x75,
	0x62, 0x6d, 0x69, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65,
	0x72, 0x22, 0x45, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f,
	0x62, 0x6c, 0x65, 0x6d, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x22, 0x42, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x41,
	0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a,
	0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x52, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x32, 0x8e, 0x02, 0x0a,
	0x0d, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x54,
	0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x12, 0x21, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x22, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x57, 0x0a, 0x0c, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x41, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x12, 0x22, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x41,
	0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4e, 0x0a,
	0x09, 0x47, 0x65, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x4d, 0x5a,
	0x4b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73,
	0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x72, 0x65, 0x67, 0x61, 0x6c, 0x69, 0x61, 0x2f,
	0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x3b,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_contestant_v1_answer_proto_rawDescOnce sync.Once
	file_contestant_v1_answer_proto_rawDescData []byte
)

func file_contestant_v1_answer_proto_rawDescGZIP() []byte {
	file_contestant_v1_answer_proto_rawDescOnce.Do(func() {
		file_contestant_v1_answer_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_contestant_v1_answer_proto_rawDesc), len(file_contestant_v1_answer_proto_rawDesc)))
	})
	return file_contestant_v1_answer_proto_rawDescData
}

var file_contestant_v1_answer_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_contestant_v1_answer_proto_goTypes = []any{
	(*Answer)(nil),                // 0: contestant.v1.Answer
	(*AnswerBody)(nil),            // 1: contestant.v1.AnswerBody
	(*DescriptiveAnswer)(nil),     // 2: contestant.v1.DescriptiveAnswer
	(*ListAnswersRequest)(nil),    // 3: contestant.v1.ListAnswersRequest
	(*ListAnswersResponse)(nil),   // 4: contestant.v1.ListAnswersResponse
	(*SubmitAnswerRequest)(nil),   // 5: contestant.v1.SubmitAnswerRequest
	(*SubmitAnswerResponse)(nil),  // 6: contestant.v1.SubmitAnswerResponse
	(*GetAnswerRequest)(nil),      // 7: contestant.v1.GetAnswerRequest
	(*GetAnswerResponse)(nil),     // 8: contestant.v1.GetAnswerResponse
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
	(*Score)(nil),                 // 10: contestant.v1.Score
	(ProblemType)(0),              // 11: contestant.v1.ProblemType
	(*durationpb.Duration)(nil),   // 12: google.protobuf.Duration
}
var file_contestant_v1_answer_proto_depIdxs = []int32{
	1,  // 0: contestant.v1.Answer.body:type_name -> contestant.v1.AnswerBody
	9,  // 1: contestant.v1.Answer.submitted_at:type_name -> google.protobuf.Timestamp
	10, // 2: contestant.v1.Answer.score:type_name -> contestant.v1.Score
	11, // 3: contestant.v1.AnswerBody.type:type_name -> contestant.v1.ProblemType
	2,  // 4: contestant.v1.AnswerBody.descriptive:type_name -> contestant.v1.DescriptiveAnswer
	0,  // 5: contestant.v1.ListAnswersResponse.answers:type_name -> contestant.v1.Answer
	12, // 6: contestant.v1.ListAnswersResponse.submit_interval:type_name -> google.protobuf.Duration
	9,  // 7: contestant.v1.ListAnswersResponse.last_submitted_at:type_name -> google.protobuf.Timestamp
	0,  // 8: contestant.v1.SubmitAnswerResponse.answer:type_name -> contestant.v1.Answer
	0,  // 9: contestant.v1.GetAnswerResponse.answer:type_name -> contestant.v1.Answer
	3,  // 10: contestant.v1.AnswerService.ListAnswers:input_type -> contestant.v1.ListAnswersRequest
	5,  // 11: contestant.v1.AnswerService.SubmitAnswer:input_type -> contestant.v1.SubmitAnswerRequest
	7,  // 12: contestant.v1.AnswerService.GetAnswer:input_type -> contestant.v1.GetAnswerRequest
	4,  // 13: contestant.v1.AnswerService.ListAnswers:output_type -> contestant.v1.ListAnswersResponse
	6,  // 14: contestant.v1.AnswerService.SubmitAnswer:output_type -> contestant.v1.SubmitAnswerResponse
	8,  // 15: contestant.v1.AnswerService.GetAnswer:output_type -> contestant.v1.GetAnswerResponse
	13, // [13:16] is the sub-list for method output_type
	10, // [10:13] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_contestant_v1_answer_proto_init() }
func file_contestant_v1_answer_proto_init() {
	if File_contestant_v1_answer_proto != nil {
		return
	}
	file_contestant_v1_problem_proto_init()
	file_contestant_v1_answer_proto_msgTypes[0].OneofWrappers = []any{}
	file_contestant_v1_answer_proto_msgTypes[1].OneofWrappers = []any{
		(*AnswerBody_Descriptive)(nil),
	}
	file_contestant_v1_answer_proto_msgTypes[4].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_contestant_v1_answer_proto_rawDesc), len(file_contestant_v1_answer_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contestant_v1_answer_proto_goTypes,
		DependencyIndexes: file_contestant_v1_answer_proto_depIdxs,
		MessageInfos:      file_contestant_v1_answer_proto_msgTypes,
	}.Build()
	File_contestant_v1_answer_proto = out.File
	file_contestant_v1_answer_proto_goTypes = nil
	file_contestant_v1_answer_proto_depIdxs = nil
}
