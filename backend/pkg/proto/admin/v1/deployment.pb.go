// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: admin/v1/deployment.proto

package adminv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type DeploymentEventType int32

const (
	DeploymentEventType_DEPLOYMENT_EVENT_TYPE_UNSPECIFIED DeploymentEventType = 0
	DeploymentEventType_DEPLOYMENT_EVENT_TYPE_QUEUED      DeploymentEventType = 1
	DeploymentEventType_DEPLOYMENT_EVENT_TYPE_CREATING    DeploymentEventType = 2
	DeploymentEventType_DEPLOYMENT_EVENT_TYPE_FINISHED    DeploymentEventType = 3
	DeploymentEventType_DEPLOYMENT_EVENT_TYPE_ERROR       DeploymentEventType = 4
)

// Enum value maps for DeploymentEventType.
var (
	DeploymentEventType_name = map[int32]string{
		0: "DEPLOYMENT_EVENT_TYPE_UNSPECIFIED",
		1: "DEPLOYMENT_EVENT_TYPE_QUEUED",
		2: "DEPLOYMENT_EVENT_TYPE_CREATING",
		3: "DEPLOYMENT_EVENT_TYPE_FINISHED",
		4: "DEPLOYMENT_EVENT_TYPE_ERROR",
	}
	DeploymentEventType_value = map[string]int32{
		"DEPLOYMENT_EVENT_TYPE_UNSPECIFIED": 0,
		"DEPLOYMENT_EVENT_TYPE_QUEUED":      1,
		"DEPLOYMENT_EVENT_TYPE_CREATING":    2,
		"DEPLOYMENT_EVENT_TYPE_FINISHED":    3,
		"DEPLOYMENT_EVENT_TYPE_ERROR":       4,
	}
)

func (x DeploymentEventType) Enum() *DeploymentEventType {
	p := new(DeploymentEventType)
	*p = x
	return p
}

func (x DeploymentEventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DeploymentEventType) Descriptor() protoreflect.EnumDescriptor {
	return file_admin_v1_deployment_proto_enumTypes[0].Descriptor()
}

func (DeploymentEventType) Type() protoreflect.EnumType {
	return &file_admin_v1_deployment_proto_enumTypes[0]
}

func (x DeploymentEventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DeploymentEventType.Descriptor instead.
func (DeploymentEventType) EnumDescriptor() ([]byte, []int) {
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{0}
}

// 問題の展開状態
type Deployment struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// チームコード
	TeamCode int64 `protobuf:"varint,1,opt,name=team_code,json=teamCode,proto3" json:"team_code,omitempty"`
	// 問題コード
	ProblemCode string `protobuf:"bytes,2,opt,name=problem_code,json=problemCode,proto3" json:"problem_code,omitempty"`
	// リビジョン - 0 が初期状態で再展開される度にインクリメントされる
	Revision int64 `protobuf:"varint,3,opt,name=revision,proto3" json:"revision,omitempty"`
	// 最新のイベント
	LatestEvent DeploymentEventType `protobuf:"varint,4,opt,name=latest_event,json=latestEvent,proto3,enum=admin.v1.DeploymentEventType" json:"latest_event,omitempty"`
	// イベント
	Events        []*DeploymentEvent `protobuf:"bytes,5,rep,name=events,proto3" json:"events,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Deployment) Reset() {
	*x = Deployment{}
	mi := &file_admin_v1_deployment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Deployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deployment) ProtoMessage() {}

func (x *Deployment) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_deployment_proto_msgTypes[0]
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
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{0}
}

func (x *Deployment) GetTeamCode() int64 {
	if x != nil {
		return x.TeamCode
	}
	return 0
}

func (x *Deployment) GetProblemCode() string {
	if x != nil {
		return x.ProblemCode
	}
	return ""
}

func (x *Deployment) GetRevision() int64 {
	if x != nil {
		return x.Revision
	}
	return 0
}

func (x *Deployment) GetLatestEvent() DeploymentEventType {
	if x != nil {
		return x.LatestEvent
	}
	return DeploymentEventType_DEPLOYMENT_EVENT_TYPE_UNSPECIFIED
}

func (x *Deployment) GetEvents() []*DeploymentEvent {
	if x != nil {
		return x.Events
	}
	return nil
}

// 問題展開に関するイベント
type DeploymentEvent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OccurredAt    *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=occurred_at,json=occurredAt,proto3" json:"occurred_at,omitempty"`
	Type          DeploymentEventType    `protobuf:"varint,2,opt,name=type,proto3,enum=admin.v1.DeploymentEventType" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeploymentEvent) Reset() {
	*x = DeploymentEvent{}
	mi := &file_admin_v1_deployment_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeploymentEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentEvent) ProtoMessage() {}

func (x *DeploymentEvent) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_deployment_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentEvent.ProtoReflect.Descriptor instead.
func (*DeploymentEvent) Descriptor() ([]byte, []int) {
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{1}
}

func (x *DeploymentEvent) GetOccurredAt() *timestamppb.Timestamp {
	if x != nil {
		return x.OccurredAt
	}
	return nil
}

func (x *DeploymentEvent) GetType() DeploymentEventType {
	if x != nil {
		return x.Type
	}
	return DeploymentEventType_DEPLOYMENT_EVENT_TYPE_UNSPECIFIED
}

type ListDeploymentsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TeamCode      int64                  `protobuf:"varint,1,opt,name=team_code,json=teamCode,proto3" json:"team_code,omitempty"`
	ProblemCode   string                 `protobuf:"bytes,2,opt,name=problem_code,json=problemCode,proto3" json:"problem_code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListDeploymentsRequest) Reset() {
	*x = ListDeploymentsRequest{}
	mi := &file_admin_v1_deployment_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDeploymentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDeploymentsRequest) ProtoMessage() {}

func (x *ListDeploymentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_deployment_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDeploymentsRequest.ProtoReflect.Descriptor instead.
func (*ListDeploymentsRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{2}
}

func (x *ListDeploymentsRequest) GetTeamCode() int64 {
	if x != nil {
		return x.TeamCode
	}
	return 0
}

func (x *ListDeploymentsRequest) GetProblemCode() string {
	if x != nil {
		return x.ProblemCode
	}
	return ""
}

type ListDeploymentsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Deployments   []*Deployment          `protobuf:"bytes,1,rep,name=deployments,proto3" json:"deployments,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListDeploymentsResponse) Reset() {
	*x = ListDeploymentsResponse{}
	mi := &file_admin_v1_deployment_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDeploymentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDeploymentsResponse) ProtoMessage() {}

func (x *ListDeploymentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_deployment_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDeploymentsResponse.ProtoReflect.Descriptor instead.
func (*ListDeploymentsResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{3}
}

func (x *ListDeploymentsResponse) GetDeployments() []*Deployment {
	if x != nil {
		return x.Deployments
	}
	return nil
}

type UpdateDeploymentStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TeamCode      int64                  `protobuf:"varint,1,opt,name=team_code,json=teamCode,proto3" json:"team_code,omitempty"`
	ProblemCode   string                 `protobuf:"bytes,2,opt,name=problem_code,json=problemCode,proto3" json:"problem_code,omitempty"`
	Revision      uint32                 `protobuf:"varint,3,opt,name=revision,proto3" json:"revision,omitempty"`
	Status        DeploymentEventType    `protobuf:"varint,4,opt,name=status,proto3,enum=admin.v1.DeploymentEventType" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateDeploymentStatusRequest) Reset() {
	*x = UpdateDeploymentStatusRequest{}
	mi := &file_admin_v1_deployment_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateDeploymentStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDeploymentStatusRequest) ProtoMessage() {}

func (x *UpdateDeploymentStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_deployment_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDeploymentStatusRequest.ProtoReflect.Descriptor instead.
func (*UpdateDeploymentStatusRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateDeploymentStatusRequest) GetTeamCode() int64 {
	if x != nil {
		return x.TeamCode
	}
	return 0
}

func (x *UpdateDeploymentStatusRequest) GetProblemCode() string {
	if x != nil {
		return x.ProblemCode
	}
	return ""
}

func (x *UpdateDeploymentStatusRequest) GetRevision() uint32 {
	if x != nil {
		return x.Revision
	}
	return 0
}

func (x *UpdateDeploymentStatusRequest) GetStatus() DeploymentEventType {
	if x != nil {
		return x.Status
	}
	return DeploymentEventType_DEPLOYMENT_EVENT_TYPE_UNSPECIFIED
}

type UpdateDeploymentStatusResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateDeploymentStatusResponse) Reset() {
	*x = UpdateDeploymentStatusResponse{}
	mi := &file_admin_v1_deployment_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateDeploymentStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDeploymentStatusResponse) ProtoMessage() {}

func (x *UpdateDeploymentStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_deployment_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDeploymentStatusResponse.ProtoReflect.Descriptor instead.
func (*UpdateDeploymentStatusResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{5}
}

type SyncDeploymentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TeamCode      int64                  `protobuf:"varint,1,opt,name=team_code,json=teamCode,proto3" json:"team_code,omitempty"`
	ProblemCode   string                 `protobuf:"bytes,2,opt,name=problem_code,json=problemCode,proto3" json:"problem_code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncDeploymentRequest) Reset() {
	*x = SyncDeploymentRequest{}
	mi := &file_admin_v1_deployment_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncDeploymentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncDeploymentRequest) ProtoMessage() {}

func (x *SyncDeploymentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_deployment_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncDeploymentRequest.ProtoReflect.Descriptor instead.
func (*SyncDeploymentRequest) Descriptor() ([]byte, []int) {
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{6}
}

func (x *SyncDeploymentRequest) GetTeamCode() int64 {
	if x != nil {
		return x.TeamCode
	}
	return 0
}

func (x *SyncDeploymentRequest) GetProblemCode() string {
	if x != nil {
		return x.ProblemCode
	}
	return ""
}

type SyncDeploymentResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncDeploymentResponse) Reset() {
	*x = SyncDeploymentResponse{}
	mi := &file_admin_v1_deployment_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncDeploymentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncDeploymentResponse) ProtoMessage() {}

func (x *SyncDeploymentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_deployment_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncDeploymentResponse.ProtoReflect.Descriptor instead.
func (*SyncDeploymentResponse) Descriptor() ([]byte, []int) {
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{7}
}

type DeployRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TeamCode      int64                  `protobuf:"varint,1,opt,name=team_code,json=teamCode,proto3" json:"team_code,omitempty"`
	ProblemCode   string                 `protobuf:"bytes,2,opt,name=problem_code,json=problemCode,proto3" json:"problem_code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeployRequest) Reset() {
	*x = DeployRequest{}
	mi := &file_admin_v1_deployment_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeployRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeployRequest) ProtoMessage() {}

func (x *DeployRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_deployment_proto_msgTypes[8]
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
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{8}
}

func (x *DeployRequest) GetTeamCode() int64 {
	if x != nil {
		return x.TeamCode
	}
	return 0
}

func (x *DeployRequest) GetProblemCode() string {
	if x != nil {
		return x.ProblemCode
	}
	return ""
}

type DeployResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Deployment    *Deployment            `protobuf:"bytes,1,opt,name=deployment,proto3" json:"deployment,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeployResponse) Reset() {
	*x = DeployResponse{}
	mi := &file_admin_v1_deployment_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeployResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeployResponse) ProtoMessage() {}

func (x *DeployResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_deployment_proto_msgTypes[9]
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
	return file_admin_v1_deployment_proto_rawDescGZIP(), []int{9}
}

func (x *DeployResponse) GetDeployment() *Deployment {
	if x != nil {
		return x.Deployment
	}
	return nil
}

var File_admin_v1_deployment_proto protoreflect.FileDescriptor

const file_admin_v1_deployment_proto_rawDesc = "" +
	"\n" +
	"\x19admin/v1/deployment.proto\x12\badmin.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"\xdd\x01\n" +
	"\n" +
	"Deployment\x12\x1b\n" +
	"\tteam_code\x18\x01 \x01(\x03R\bteamCode\x12!\n" +
	"\fproblem_code\x18\x02 \x01(\tR\vproblemCode\x12\x1a\n" +
	"\brevision\x18\x03 \x01(\x03R\brevision\x12@\n" +
	"\flatest_event\x18\x04 \x01(\x0e2\x1d.admin.v1.DeploymentEventTypeR\vlatestEvent\x121\n" +
	"\x06events\x18\x05 \x03(\v2\x19.admin.v1.DeploymentEventR\x06events\"\x81\x01\n" +
	"\x0fDeploymentEvent\x12;\n" +
	"\voccurred_at\x18\x01 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"occurredAt\x121\n" +
	"\x04type\x18\x02 \x01(\x0e2\x1d.admin.v1.DeploymentEventTypeR\x04type\"X\n" +
	"\x16ListDeploymentsRequest\x12\x1b\n" +
	"\tteam_code\x18\x01 \x01(\x03R\bteamCode\x12!\n" +
	"\fproblem_code\x18\x02 \x01(\tR\vproblemCode\"Q\n" +
	"\x17ListDeploymentsResponse\x126\n" +
	"\vdeployments\x18\x01 \x03(\v2\x14.admin.v1.DeploymentR\vdeployments\"\xb2\x01\n" +
	"\x1dUpdateDeploymentStatusRequest\x12\x1b\n" +
	"\tteam_code\x18\x01 \x01(\x03R\bteamCode\x12!\n" +
	"\fproblem_code\x18\x02 \x01(\tR\vproblemCode\x12\x1a\n" +
	"\brevision\x18\x03 \x01(\rR\brevision\x125\n" +
	"\x06status\x18\x04 \x01(\x0e2\x1d.admin.v1.DeploymentEventTypeR\x06status\" \n" +
	"\x1eUpdateDeploymentStatusResponse\"W\n" +
	"\x15SyncDeploymentRequest\x12\x1b\n" +
	"\tteam_code\x18\x01 \x01(\x03R\bteamCode\x12!\n" +
	"\fproblem_code\x18\x02 \x01(\tR\vproblemCode\"\x18\n" +
	"\x16SyncDeploymentResponse\"O\n" +
	"\rDeployRequest\x12\x1b\n" +
	"\tteam_code\x18\x01 \x01(\x03R\bteamCode\x12!\n" +
	"\fproblem_code\x18\x02 \x01(\tR\vproblemCode\"F\n" +
	"\x0eDeployResponse\x124\n" +
	"\n" +
	"deployment\x18\x01 \x01(\v2\x14.admin.v1.DeploymentR\n" +
	"deployment*\xc7\x01\n" +
	"\x13DeploymentEventType\x12%\n" +
	"!DEPLOYMENT_EVENT_TYPE_UNSPECIFIED\x10\x00\x12 \n" +
	"\x1cDEPLOYMENT_EVENT_TYPE_QUEUED\x10\x01\x12\"\n" +
	"\x1eDEPLOYMENT_EVENT_TYPE_CREATING\x10\x02\x12\"\n" +
	"\x1eDEPLOYMENT_EVENT_TYPE_FINISHED\x10\x03\x12\x1f\n" +
	"\x1bDEPLOYMENT_EVENT_TYPE_ERROR\x10\x042\xea\x02\n" +
	"\x11DeploymentService\x12V\n" +
	"\x0fListDeployments\x12 .admin.v1.ListDeploymentsRequest\x1a!.admin.v1.ListDeploymentsResponse\x12k\n" +
	"\x16UpdateDeploymentStatus\x12'.admin.v1.UpdateDeploymentStatusRequest\x1a(.admin.v1.UpdateDeploymentStatusResponse\x12S\n" +
	"\x0eSyncDeployment\x12\x1f.admin.v1.SyncDeploymentRequest\x1a .admin.v1.SyncDeploymentResponse\x12;\n" +
	"\x06Deploy\x12\x17.admin.v1.DeployRequest\x1a\x18.admin.v1.DeployResponseBCZAgithub.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1;adminv1b\x06proto3"

var (
	file_admin_v1_deployment_proto_rawDescOnce sync.Once
	file_admin_v1_deployment_proto_rawDescData []byte
)

func file_admin_v1_deployment_proto_rawDescGZIP() []byte {
	file_admin_v1_deployment_proto_rawDescOnce.Do(func() {
		file_admin_v1_deployment_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_admin_v1_deployment_proto_rawDesc), len(file_admin_v1_deployment_proto_rawDesc)))
	})
	return file_admin_v1_deployment_proto_rawDescData
}

var file_admin_v1_deployment_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_admin_v1_deployment_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_admin_v1_deployment_proto_goTypes = []any{
	(DeploymentEventType)(0),               // 0: admin.v1.DeploymentEventType
	(*Deployment)(nil),                     // 1: admin.v1.Deployment
	(*DeploymentEvent)(nil),                // 2: admin.v1.DeploymentEvent
	(*ListDeploymentsRequest)(nil),         // 3: admin.v1.ListDeploymentsRequest
	(*ListDeploymentsResponse)(nil),        // 4: admin.v1.ListDeploymentsResponse
	(*UpdateDeploymentStatusRequest)(nil),  // 5: admin.v1.UpdateDeploymentStatusRequest
	(*UpdateDeploymentStatusResponse)(nil), // 6: admin.v1.UpdateDeploymentStatusResponse
	(*SyncDeploymentRequest)(nil),          // 7: admin.v1.SyncDeploymentRequest
	(*SyncDeploymentResponse)(nil),         // 8: admin.v1.SyncDeploymentResponse
	(*DeployRequest)(nil),                  // 9: admin.v1.DeployRequest
	(*DeployResponse)(nil),                 // 10: admin.v1.DeployResponse
	(*timestamppb.Timestamp)(nil),          // 11: google.protobuf.Timestamp
}
var file_admin_v1_deployment_proto_depIdxs = []int32{
	0,  // 0: admin.v1.Deployment.latest_event:type_name -> admin.v1.DeploymentEventType
	2,  // 1: admin.v1.Deployment.events:type_name -> admin.v1.DeploymentEvent
	11, // 2: admin.v1.DeploymentEvent.occurred_at:type_name -> google.protobuf.Timestamp
	0,  // 3: admin.v1.DeploymentEvent.type:type_name -> admin.v1.DeploymentEventType
	1,  // 4: admin.v1.ListDeploymentsResponse.deployments:type_name -> admin.v1.Deployment
	0,  // 5: admin.v1.UpdateDeploymentStatusRequest.status:type_name -> admin.v1.DeploymentEventType
	1,  // 6: admin.v1.DeployResponse.deployment:type_name -> admin.v1.Deployment
	3,  // 7: admin.v1.DeploymentService.ListDeployments:input_type -> admin.v1.ListDeploymentsRequest
	5,  // 8: admin.v1.DeploymentService.UpdateDeploymentStatus:input_type -> admin.v1.UpdateDeploymentStatusRequest
	7,  // 9: admin.v1.DeploymentService.SyncDeployment:input_type -> admin.v1.SyncDeploymentRequest
	9,  // 10: admin.v1.DeploymentService.Deploy:input_type -> admin.v1.DeployRequest
	4,  // 11: admin.v1.DeploymentService.ListDeployments:output_type -> admin.v1.ListDeploymentsResponse
	6,  // 12: admin.v1.DeploymentService.UpdateDeploymentStatus:output_type -> admin.v1.UpdateDeploymentStatusResponse
	8,  // 13: admin.v1.DeploymentService.SyncDeployment:output_type -> admin.v1.SyncDeploymentResponse
	10, // 14: admin.v1.DeploymentService.Deploy:output_type -> admin.v1.DeployResponse
	11, // [11:15] is the sub-list for method output_type
	7,  // [7:11] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_admin_v1_deployment_proto_init() }
func file_admin_v1_deployment_proto_init() {
	if File_admin_v1_deployment_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_admin_v1_deployment_proto_rawDesc), len(file_admin_v1_deployment_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_v1_deployment_proto_goTypes,
		DependencyIndexes: file_admin_v1_deployment_proto_depIdxs,
		EnumInfos:         file_admin_v1_deployment_proto_enumTypes,
		MessageInfos:      file_admin_v1_deployment_proto_msgTypes,
	}.Build()
	File_admin_v1_deployment_proto = out.File
	file_admin_v1_deployment_proto_goTypes = nil
	file_admin_v1_deployment_proto_depIdxs = nil
}
