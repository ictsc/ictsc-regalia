// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: admin/v1/ranking.proto

package adminv1

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

type Score struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	Team    *Team                  `protobuf:"bytes,1,opt,name=team,proto3" json:"team,omitempty"`
	Problem *Problem               `protobuf:"bytes,2,opt,name=problem,proto3" json:"problem,omitempty"`
	// 採点による得点
	MarkedScore int64 `protobuf:"varint,3,opt,name=marked_score,json=markedScore,proto3" json:"marked_score,omitempty"`
	// ペナルティによる減点
	Penalty int64 `protobuf:"varint,4,opt,name=penalty,proto3" json:"penalty,omitempty"`
	// 最終的な得点
	Score         int64 `protobuf:"varint,5,opt,name=score,proto3" json:"score,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Score) Reset() {
	*x = Score{}
	mi := &file_admin_v1_ranking_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Score) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Score) ProtoMessage() {}

func (x *Score) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_ranking_proto_msgTypes[0]
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
	return file_admin_v1_ranking_proto_rawDescGZIP(), []int{0}
}

func (x *Score) GetTeam() *Team {
	if x != nil {
		return x.Team
	}
	return nil
}

func (x *Score) GetProblem() *Problem {
	if x != nil {
		return x.Problem
	}
	return nil
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

type TeamRank struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Team          *Team                  `protobuf:"bytes,1,opt,name=team,proto3" json:"team,omitempty"`
	Rank          int64                  `protobuf:"varint,2,opt,name=rank,proto3" json:"rank,omitempty"`
	Score         int64                  `protobuf:"varint,3,opt,name=score,proto3" json:"score,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TeamRank) Reset() {
	*x = TeamRank{}
	mi := &file_admin_v1_ranking_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TeamRank) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TeamRank) ProtoMessage() {}

func (x *TeamRank) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_ranking_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TeamRank.ProtoReflect.Descriptor instead.
func (*TeamRank) Descriptor() ([]byte, []int) {
	return file_admin_v1_ranking_proto_rawDescGZIP(), []int{1}
}

func (x *TeamRank) GetTeam() *Team {
	if x != nil {
		return x.Team
	}
	return nil
}

func (x *TeamRank) GetRank() int64 {
	if x != nil {
		return x.Rank
	}
	return 0
}

func (x *TeamRank) GetScore() int64 {
	if x != nil {
		return x.Score
	}
	return 0
}

type ListScoreRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListScoreRequest) Reset() {
	*x = ListScoreRequest{}
	mi := &file_admin_v1_ranking_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListScoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListScoreRequest) ProtoMessage() {}

func (x *ListScoreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_ranking_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListScoreRequest.ProtoReflect.Descriptor instead.
func (*ListScoreRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_ranking_proto_rawDescGZIP(), []int{2}
}

type ListScoreResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Scores        []*Score               `protobuf:"bytes,1,rep,name=scores,proto3" json:"scores,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListScoreResponse) Reset() {
	*x = ListScoreResponse{}
	mi := &file_admin_v1_ranking_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListScoreResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListScoreResponse) ProtoMessage() {}

func (x *ListScoreResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_ranking_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListScoreResponse.ProtoReflect.Descriptor instead.
func (*ListScoreResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_ranking_proto_rawDescGZIP(), []int{3}
}

func (x *ListScoreResponse) GetScores() []*Score {
	if x != nil {
		return x.Scores
	}
	return nil
}

type GetRankingRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRankingRequest) Reset() {
	*x = GetRankingRequest{}
	mi := &file_admin_v1_ranking_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRankingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRankingRequest) ProtoMessage() {}

func (x *GetRankingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_ranking_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRankingRequest.ProtoReflect.Descriptor instead.
func (*GetRankingRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_ranking_proto_rawDescGZIP(), []int{4}
}

type GetRankingResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ranking       []*TeamRank            `protobuf:"bytes,1,rep,name=ranking,proto3" json:"ranking,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRankingResponse) Reset() {
	*x = GetRankingResponse{}
	mi := &file_admin_v1_ranking_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRankingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRankingResponse) ProtoMessage() {}

func (x *GetRankingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_ranking_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRankingResponse.ProtoReflect.Descriptor instead.
func (*GetRankingResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_ranking_proto_rawDescGZIP(), []int{5}
}

func (x *GetRankingResponse) GetRanking() []*TeamRank {
	if x != nil {
		return x.Ranking
	}
	return nil
}

var File_admin_v1_ranking_proto protoreflect.FileDescriptor

var file_admin_v1_ranking_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x61, 0x6e, 0x6b, 0x69,
	0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x76, 0x31, 0x1a, 0x16, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f,
	0x62, 0x6c, 0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xab, 0x01, 0x0a, 0x05, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x65, 0x61,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x04, 0x74, 0x65, 0x61, 0x6d, 0x12, 0x2b, 0x0a,
	0x07, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65,
	0x6d, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x61,
	0x72, 0x6b, 0x65, 0x64, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0b, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x64, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x70, 0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x70, 0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x22, 0x58, 0x0a,
	0x08, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x61, 0x6e, 0x6b, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x65, 0x61,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x04, 0x74, 0x65, 0x61, 0x6d, 0x12, 0x12, 0x0a,
	0x04, 0x72, 0x61, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x72, 0x61, 0x6e,
	0x6b, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x22, 0x12, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x53,
	0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3c, 0x0a, 0x11, 0x4c,
	0x69, 0x73, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x27, 0x0a, 0x06, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x6f, 0x72,
	0x65, 0x52, 0x06, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x22, 0x13, 0x0a, 0x11, 0x47, 0x65, 0x74,
	0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x42,
	0x0a, 0x12, 0x47, 0x65, 0x74, 0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x07, 0x72, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x61, 0x6e, 0x6b, 0x52, 0x07, 0x72, 0x61, 0x6e, 0x6b, 0x69,
	0x6e, 0x67, 0x32, 0x9f, 0x01, 0x0a, 0x0e, 0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x44, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x63, 0x6f,
	0x72, 0x65, 0x12, 0x1a, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b,
	0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x63,
	0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0a, 0x47,
	0x65, 0x74, 0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x12, 0x1b, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2f, 0x69, 0x63, 0x74, 0x73, 0x63, 0x2d, 0x72,
	0x65, 0x67, 0x61, 0x6c, 0x69, 0x61, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76,
	0x31, 0x3b, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_admin_v1_ranking_proto_rawDescOnce sync.Once
	file_admin_v1_ranking_proto_rawDescData = file_admin_v1_ranking_proto_rawDesc
)

func file_admin_v1_ranking_proto_rawDescGZIP() []byte {
	file_admin_v1_ranking_proto_rawDescOnce.Do(func() {
		file_admin_v1_ranking_proto_rawDescData = protoimpl.X.CompressGZIP(file_admin_v1_ranking_proto_rawDescData)
	})
	return file_admin_v1_ranking_proto_rawDescData
}

var file_admin_v1_ranking_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_admin_v1_ranking_proto_goTypes = []any{
	(*Score)(nil),              // 0: admin.v1.Score
	(*TeamRank)(nil),           // 1: admin.v1.TeamRank
	(*ListScoreRequest)(nil),   // 2: admin.v1.ListScoreRequest
	(*ListScoreResponse)(nil),  // 3: admin.v1.ListScoreResponse
	(*GetRankingRequest)(nil),  // 4: admin.v1.GetRankingRequest
	(*GetRankingResponse)(nil), // 5: admin.v1.GetRankingResponse
	(*Team)(nil),               // 6: admin.v1.Team
	(*Problem)(nil),            // 7: admin.v1.Problem
}
var file_admin_v1_ranking_proto_depIdxs = []int32{
	6, // 0: admin.v1.Score.team:type_name -> admin.v1.Team
	7, // 1: admin.v1.Score.problem:type_name -> admin.v1.Problem
	6, // 2: admin.v1.TeamRank.team:type_name -> admin.v1.Team
	0, // 3: admin.v1.ListScoreResponse.scores:type_name -> admin.v1.Score
	1, // 4: admin.v1.GetRankingResponse.ranking:type_name -> admin.v1.TeamRank
	2, // 5: admin.v1.RankingService.ListScore:input_type -> admin.v1.ListScoreRequest
	4, // 6: admin.v1.RankingService.GetRanking:input_type -> admin.v1.GetRankingRequest
	3, // 7: admin.v1.RankingService.ListScore:output_type -> admin.v1.ListScoreResponse
	5, // 8: admin.v1.RankingService.GetRanking:output_type -> admin.v1.GetRankingResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_admin_v1_ranking_proto_init() }
func file_admin_v1_ranking_proto_init() {
	if File_admin_v1_ranking_proto != nil {
		return
	}
	file_admin_v1_problem_proto_init()
	file_admin_v1_team_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_admin_v1_ranking_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_ranking_proto_goTypes,
		DependencyIndexes: file_admin_v1_ranking_proto_depIdxs,
		MessageInfos:      file_admin_v1_ranking_proto_msgTypes,
	}.Build()
	File_admin_v1_ranking_proto = out.File
	file_admin_v1_ranking_proto_rawDesc = nil
	file_admin_v1_ranking_proto_goTypes = nil
	file_admin_v1_ranking_proto_depIdxs = nil
}