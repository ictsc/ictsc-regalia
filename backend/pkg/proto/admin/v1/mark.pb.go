// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: admin/v1/mark.proto

package adminv1

import (
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

// 解答
type Answer struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Team          *Team                  `protobuf:"bytes,2,opt,name=team,proto3" json:"team,omitempty"`
	Author        *Contestant            `protobuf:"bytes,3,opt,name=author,proto3" json:"author,omitempty"`
	Problem       *Problem               `protobuf:"bytes,4,opt,name=problem,proto3" json:"problem,omitempty"`
	Body          string                 `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Result        *MarkingResult         `protobuf:"bytes,7,opt,name=result,proto3,oneof" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Answer) Reset() {
	*x = Answer{}
	mi := &file_admin_v1_mark_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Answer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Answer) ProtoMessage() {}

func (x *Answer) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_mark_proto_msgTypes[0]
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
	return file_admin_v1_mark_proto_rawDescGZIP(), []int{0}
}

func (x *Answer) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Answer) GetTeam() *Team {
	if x != nil {
		return x.Team
	}
	return nil
}

func (x *Answer) GetAuthor() *Contestant {
	if x != nil {
		return x.Author
	}
	return nil
}

func (x *Answer) GetProblem() *Problem {
	if x != nil {
		return x.Problem
	}
	return nil
}

func (x *Answer) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *Answer) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Answer) GetResult() *MarkingResult {
	if x != nil {
		return x.Result
	}
	return nil
}

// 採点結果
type MarkingResult struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Answer        *Answer                `protobuf:"bytes,1,opt,name=answer,proto3" json:"answer,omitempty"`
	Judge         *Admin                 `protobuf:"bytes,2,opt,name=judge,proto3" json:"judge,omitempty"`
	Score         int64                  `protobuf:"varint,3,opt,name=score,proto3" json:"score,omitempty"`
	Comment       string                 `protobuf:"bytes,4,opt,name=comment,proto3" json:"comment,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MarkingResult) Reset() {
	*x = MarkingResult{}
	mi := &file_admin_v1_mark_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MarkingResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkingResult) ProtoMessage() {}

func (x *MarkingResult) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_mark_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkingResult.ProtoReflect.Descriptor instead.
func (*MarkingResult) Descriptor() ([]byte, []int) {
	return file_admin_v1_mark_proto_rawDescGZIP(), []int{1}
}

func (x *MarkingResult) GetAnswer() *Answer {
	if x != nil {
		return x.Answer
	}
	return nil
}

func (x *MarkingResult) GetJudge() *Admin {
	if x != nil {
		return x.Judge
	}
	return nil
}

func (x *MarkingResult) GetScore() int64 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *MarkingResult) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

func (x *MarkingResult) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type ListAnswersRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 採点が完了した提出を含めるかどうか
	IncludeScored bool `protobuf:"varint,1,opt,name=include_scored,json=includeScored,proto3" json:"include_scored,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListAnswersRequest) Reset() {
	*x = ListAnswersRequest{}
	mi := &file_admin_v1_mark_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListAnswersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAnswersRequest) ProtoMessage() {}

func (x *ListAnswersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_mark_proto_msgTypes[2]
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
	return file_admin_v1_mark_proto_rawDescGZIP(), []int{2}
}

func (x *ListAnswersRequest) GetIncludeScored() bool {
	if x != nil {
		return x.IncludeScored
	}
	return false
}

type ListAnswersResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Answers       []*Answer              `protobuf:"bytes,1,rep,name=answers,proto3" json:"answers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListAnswersResponse) Reset() {
	*x = ListAnswersResponse{}
	mi := &file_admin_v1_mark_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListAnswersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAnswersResponse) ProtoMessage() {}

func (x *ListAnswersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_mark_proto_msgTypes[3]
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
	return file_admin_v1_mark_proto_rawDescGZIP(), []int{3}
}

func (x *ListAnswersResponse) GetAnswers() []*Answer {
	if x != nil {
		return x.Answers
	}
	return nil
}

type ListMarkingResultsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListMarkingResultsRequest) Reset() {
	*x = ListMarkingResultsRequest{}
	mi := &file_admin_v1_mark_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListMarkingResultsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMarkingResultsRequest) ProtoMessage() {}

func (x *ListMarkingResultsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_mark_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMarkingResultsRequest.ProtoReflect.Descriptor instead.
func (*ListMarkingResultsRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_mark_proto_rawDescGZIP(), []int{4}
}

type ListMarkingResultsResponse struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	MarkingResults []*MarkingResult       `protobuf:"bytes,1,rep,name=marking_results,json=markingResults,proto3" json:"marking_results,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *ListMarkingResultsResponse) Reset() {
	*x = ListMarkingResultsResponse{}
	mi := &file_admin_v1_mark_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListMarkingResultsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMarkingResultsResponse) ProtoMessage() {}

func (x *ListMarkingResultsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_mark_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMarkingResultsResponse.ProtoReflect.Descriptor instead.
func (*ListMarkingResultsResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_mark_proto_rawDescGZIP(), []int{5}
}

func (x *ListMarkingResultsResponse) GetMarkingResults() []*MarkingResult {
	if x != nil {
		return x.MarkingResults
	}
	return nil
}

type CreateMarkingResultRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MarkingResult *MarkingResult         `protobuf:"bytes,1,opt,name=marking_result,json=markingResult,proto3" json:"marking_result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateMarkingResultRequest) Reset() {
	*x = CreateMarkingResultRequest{}
	mi := &file_admin_v1_mark_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateMarkingResultRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMarkingResultRequest) ProtoMessage() {}

func (x *CreateMarkingResultRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_mark_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMarkingResultRequest.ProtoReflect.Descriptor instead.
func (*CreateMarkingResultRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_mark_proto_rawDescGZIP(), []int{6}
}

func (x *CreateMarkingResultRequest) GetMarkingResult() *MarkingResult {
	if x != nil {
		return x.MarkingResult
	}
	return nil
}

type CreateMarkingResultResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MarkingResult *MarkingResult         `protobuf:"bytes,1,opt,name=marking_result,json=markingResult,proto3" json:"marking_result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateMarkingResultResponse) Reset() {
	*x = CreateMarkingResultResponse{}
	mi := &file_admin_v1_mark_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateMarkingResultResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMarkingResultResponse) ProtoMessage() {}

func (x *CreateMarkingResultResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_mark_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMarkingResultResponse.ProtoReflect.Descriptor instead.
func (*CreateMarkingResultResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_mark_proto_rawDescGZIP(), []int{7}
}

func (x *CreateMarkingResultResponse) GetMarkingResult() *MarkingResult {
	if x != nil {
		return x.MarkingResult
	}
	return nil
}

var File_admin_v1_mark_proto protoreflect.FileDescriptor

var file_admin_v1_mark_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x61, 0x72, 0x6b, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x1a,
	0x14, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x16, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x62, 0x6c,
	0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f,
	0x76, 0x31, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa7,
	0x02, 0x0a, 0x06, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x65, 0x61,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x04, 0x74, 0x65, 0x61, 0x6d, 0x12, 0x2c, 0x0a,
	0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x6e, 0x74, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x2b, 0x0a, 0x07, 0x70,
	0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x52,
	0x07, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x39, 0x0a, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x34, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x48, 0x00, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a,
	0x07, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0xcb, 0x01, 0x0a, 0x0d, 0x4d, 0x61, 0x72,
	0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x28, 0x0a, 0x06, 0x61, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x06, 0x61, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x05, 0x6a, 0x75, 0x64, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x52, 0x05, 0x6a, 0x75, 0x64, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73,
	0x63, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x3b, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e,
	0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x53, 0x63, 0x6f,
	0x72, 0x65, 0x64, 0x22, 0x41, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65,
	0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x61, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x07, 0x61,
	0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x22, 0x1b, 0x0a, 0x19, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x61,
	0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x5e, 0x0a, 0x1a, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x61, 0x72, 0x6b, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x40, 0x0a, 0x0f, 0x6d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x52, 0x0e, 0x6d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x73, 0x22, 0x5c, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x72,
	0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x3e, 0x0a, 0x0e, 0x6d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x52, 0x0d, 0x6d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x22, 0x5d, 0x0a, 0x1b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x3e, 0x0a, 0x0e, 0x6d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x52, 0x0d, 0x6d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x32, 0x9e, 0x02, 0x0a, 0x0b, 0x4d, 0x61, 0x72, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x4a, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x12,
	0x1c, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41,
	0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6e, 0x73,
	0x77, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5f, 0x0a, 0x12,
	0x4c, 0x69, 0x73, 0x74, 0x4d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x73, 0x12, 0x23, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x4d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x62, 0x0a,
	0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x12, 0x24, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x69, 0x63, 0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x72, 0x65, 0x67, 0x61,
	0x6c, 0x69, 0x61, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_admin_v1_mark_proto_rawDescOnce sync.Once
	file_admin_v1_mark_proto_rawDescData = file_admin_v1_mark_proto_rawDesc
)

func file_admin_v1_mark_proto_rawDescGZIP() []byte {
	file_admin_v1_mark_proto_rawDescOnce.Do(func() {
		file_admin_v1_mark_proto_rawDescData = protoimpl.X.CompressGZIP(file_admin_v1_mark_proto_rawDescData)
	})
	return file_admin_v1_mark_proto_rawDescData
}

var file_admin_v1_mark_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_admin_v1_mark_proto_goTypes = []any{
	(*Answer)(nil),                      // 0: admin.v1.Answer
	(*MarkingResult)(nil),               // 1: admin.v1.MarkingResult
	(*ListAnswersRequest)(nil),          // 2: admin.v1.ListAnswersRequest
	(*ListAnswersResponse)(nil),         // 3: admin.v1.ListAnswersResponse
	(*ListMarkingResultsRequest)(nil),   // 4: admin.v1.ListMarkingResultsRequest
	(*ListMarkingResultsResponse)(nil),  // 5: admin.v1.ListMarkingResultsResponse
	(*CreateMarkingResultRequest)(nil),  // 6: admin.v1.CreateMarkingResultRequest
	(*CreateMarkingResultResponse)(nil), // 7: admin.v1.CreateMarkingResultResponse
	(*Team)(nil),                        // 8: admin.v1.Team
	(*Contestant)(nil),                  // 9: admin.v1.Contestant
	(*Problem)(nil),                     // 10: admin.v1.Problem
	(*timestamppb.Timestamp)(nil),       // 11: google.protobuf.Timestamp
	(*Admin)(nil),                       // 12: admin.v1.Admin
}
var file_admin_v1_mark_proto_depIdxs = []int32{
	8,  // 0: admin.v1.Answer.team:type_name -> admin.v1.Team
	9,  // 1: admin.v1.Answer.author:type_name -> admin.v1.Contestant
	10, // 2: admin.v1.Answer.problem:type_name -> admin.v1.Problem
	11, // 3: admin.v1.Answer.created_at:type_name -> google.protobuf.Timestamp
	1,  // 4: admin.v1.Answer.result:type_name -> admin.v1.MarkingResult
	0,  // 5: admin.v1.MarkingResult.answer:type_name -> admin.v1.Answer
	12, // 6: admin.v1.MarkingResult.judge:type_name -> admin.v1.Admin
	11, // 7: admin.v1.MarkingResult.created_at:type_name -> google.protobuf.Timestamp
	0,  // 8: admin.v1.ListAnswersResponse.answers:type_name -> admin.v1.Answer
	1,  // 9: admin.v1.ListMarkingResultsResponse.marking_results:type_name -> admin.v1.MarkingResult
	1,  // 10: admin.v1.CreateMarkingResultRequest.marking_result:type_name -> admin.v1.MarkingResult
	1,  // 11: admin.v1.CreateMarkingResultResponse.marking_result:type_name -> admin.v1.MarkingResult
	2,  // 12: admin.v1.MarkService.ListAnswers:input_type -> admin.v1.ListAnswersRequest
	4,  // 13: admin.v1.MarkService.ListMarkingResults:input_type -> admin.v1.ListMarkingResultsRequest
	6,  // 14: admin.v1.MarkService.CreateMarkingResult:input_type -> admin.v1.CreateMarkingResultRequest
	3,  // 15: admin.v1.MarkService.ListAnswers:output_type -> admin.v1.ListAnswersResponse
	5,  // 16: admin.v1.MarkService.ListMarkingResults:output_type -> admin.v1.ListMarkingResultsResponse
	7,  // 17: admin.v1.MarkService.CreateMarkingResult:output_type -> admin.v1.CreateMarkingResultResponse
	15, // [15:18] is the sub-list for method output_type
	12, // [12:15] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_admin_v1_mark_proto_init() }
func file_admin_v1_mark_proto_init() {
	if File_admin_v1_mark_proto != nil {
		return
	}
	file_admin_v1_actor_proto_init()
	file_admin_v1_contestant_proto_init()
	file_admin_v1_problem_proto_init()
	file_admin_v1_team_proto_init()
	file_admin_v1_mark_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_admin_v1_mark_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_mark_proto_goTypes,
		DependencyIndexes: file_admin_v1_mark_proto_depIdxs,
		MessageInfos:      file_admin_v1_mark_proto_msgTypes,
	}.Build()
	File_admin_v1_mark_proto = out.File
	file_admin_v1_mark_proto_rawDesc = nil
	file_admin_v1_mark_proto_goTypes = nil
	file_admin_v1_mark_proto_depIdxs = nil
}
