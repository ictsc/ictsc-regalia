// @generated by protoc-gen-es v1.7.2 with parameter "target=ts"
// @generated from file admin/v1/team.proto (package admin.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { User } from "./user_pb.js";

/**
 * @generated from message admin.v1.Team
 */
export class Team extends Message<Team> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: int32 code = 2;
   */
  code = 0;

  /**
   * @generated from field: string name = 3;
   */
  name = "";

  /**
   * @generated from field: string organization = 4;
   */
  organization = "";

  /**
   * @generated from field: string invitation_code = 5;
   */
  invitationCode = "";

  /**
   * @generated from field: int32 code_remaining = 6;
   */
  codeRemaining = 0;

  constructor(data?: PartialMessage<Team>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.Team";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "code", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 3, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "organization", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "invitation_code", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "code_remaining", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Team {
    return new Team().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Team {
    return new Team().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Team {
    return new Team().fromJsonString(jsonString, options);
  }

  static equals(a: Team | PlainMessage<Team> | undefined, b: Team | PlainMessage<Team> | undefined): boolean {
    return proto3.util.equals(Team, a, b);
  }
}

/**
 * @generated from message admin.v1.GetTeamsRequest
 */
export class GetTeamsRequest extends Message<GetTeamsRequest> {
  constructor(data?: PartialMessage<GetTeamsRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.GetTeamsRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetTeamsRequest {
    return new GetTeamsRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetTeamsRequest {
    return new GetTeamsRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetTeamsRequest {
    return new GetTeamsRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetTeamsRequest | PlainMessage<GetTeamsRequest> | undefined, b: GetTeamsRequest | PlainMessage<GetTeamsRequest> | undefined): boolean {
    return proto3.util.equals(GetTeamsRequest, a, b);
  }
}

/**
 * @generated from message admin.v1.GetTeamsResponse
 */
export class GetTeamsResponse extends Message<GetTeamsResponse> {
  /**
   * @generated from field: repeated admin.v1.Team teams = 1;
   */
  teams: Team[] = [];

  constructor(data?: PartialMessage<GetTeamsResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.GetTeamsResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "teams", kind: "message", T: Team, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetTeamsResponse {
    return new GetTeamsResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetTeamsResponse {
    return new GetTeamsResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetTeamsResponse {
    return new GetTeamsResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetTeamsResponse | PlainMessage<GetTeamsResponse> | undefined, b: GetTeamsResponse | PlainMessage<GetTeamsResponse> | undefined): boolean {
    return proto3.util.equals(GetTeamsResponse, a, b);
  }
}

/**
 * @generated from message admin.v1.GetTeamRequest
 */
export class GetTeamRequest extends Message<GetTeamRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  constructor(data?: PartialMessage<GetTeamRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.GetTeamRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetTeamRequest {
    return new GetTeamRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetTeamRequest {
    return new GetTeamRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetTeamRequest {
    return new GetTeamRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetTeamRequest | PlainMessage<GetTeamRequest> | undefined, b: GetTeamRequest | PlainMessage<GetTeamRequest> | undefined): boolean {
    return proto3.util.equals(GetTeamRequest, a, b);
  }
}

/**
 * @generated from message admin.v1.GetTeamResponse
 */
export class GetTeamResponse extends Message<GetTeamResponse> {
  /**
   * @generated from field: admin.v1.Team team = 1;
   */
  team?: Team;

  constructor(data?: PartialMessage<GetTeamResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.GetTeamResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "team", kind: "message", T: Team },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetTeamResponse {
    return new GetTeamResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetTeamResponse {
    return new GetTeamResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetTeamResponse {
    return new GetTeamResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetTeamResponse | PlainMessage<GetTeamResponse> | undefined, b: GetTeamResponse | PlainMessage<GetTeamResponse> | undefined): boolean {
    return proto3.util.equals(GetTeamResponse, a, b);
  }
}

/**
 * @generated from message admin.v1.BaStation
 */
export class BaStation extends Message<BaStation> {
  /**
   * @generated from field: string user = 1;
   */
  user = "";

  /**
   * @generated from field: string password = 2;
   */
  password = "";

  /**
   * @generated from field: string host = 3;
   */
  host = "";

  /**
   * @generated from field: int32 port = 4;
   */
  port = 0;

  constructor(data?: PartialMessage<BaStation>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.BaStation";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "user", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "password", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "host", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "port", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): BaStation {
    return new BaStation().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): BaStation {
    return new BaStation().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): BaStation {
    return new BaStation().fromJsonString(jsonString, options);
  }

  static equals(a: BaStation | PlainMessage<BaStation> | undefined, b: BaStation | PlainMessage<BaStation> | undefined): boolean {
    return proto3.util.equals(BaStation, a, b);
  }
}

/**
 * @generated from message admin.v1.GetTeamConnectionInfoRequest
 */
export class GetTeamConnectionInfoRequest extends Message<GetTeamConnectionInfoRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  constructor(data?: PartialMessage<GetTeamConnectionInfoRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.GetTeamConnectionInfoRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetTeamConnectionInfoRequest {
    return new GetTeamConnectionInfoRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetTeamConnectionInfoRequest {
    return new GetTeamConnectionInfoRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetTeamConnectionInfoRequest {
    return new GetTeamConnectionInfoRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetTeamConnectionInfoRequest | PlainMessage<GetTeamConnectionInfoRequest> | undefined, b: GetTeamConnectionInfoRequest | PlainMessage<GetTeamConnectionInfoRequest> | undefined): boolean {
    return proto3.util.equals(GetTeamConnectionInfoRequest, a, b);
  }
}

/**
 * @generated from message admin.v1.GetTeamConnectionInfoResponse
 */
export class GetTeamConnectionInfoResponse extends Message<GetTeamConnectionInfoResponse> {
  /**
   * @generated from field: admin.v1.BaStation ba_station = 1;
   */
  baStation?: BaStation;

  constructor(data?: PartialMessage<GetTeamConnectionInfoResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.GetTeamConnectionInfoResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "ba_station", kind: "message", T: BaStation },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetTeamConnectionInfoResponse {
    return new GetTeamConnectionInfoResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetTeamConnectionInfoResponse {
    return new GetTeamConnectionInfoResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetTeamConnectionInfoResponse {
    return new GetTeamConnectionInfoResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetTeamConnectionInfoResponse | PlainMessage<GetTeamConnectionInfoResponse> | undefined, b: GetTeamConnectionInfoResponse | PlainMessage<GetTeamConnectionInfoResponse> | undefined): boolean {
    return proto3.util.equals(GetTeamConnectionInfoResponse, a, b);
  }
}

/**
 * @generated from message admin.v1.GetMembersRequest
 */
export class GetMembersRequest extends Message<GetMembersRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  constructor(data?: PartialMessage<GetMembersRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.GetMembersRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetMembersRequest {
    return new GetMembersRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetMembersRequest {
    return new GetMembersRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetMembersRequest {
    return new GetMembersRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetMembersRequest | PlainMessage<GetMembersRequest> | undefined, b: GetMembersRequest | PlainMessage<GetMembersRequest> | undefined): boolean {
    return proto3.util.equals(GetMembersRequest, a, b);
  }
}

/**
 * @generated from message admin.v1.GetMembersResponse
 */
export class GetMembersResponse extends Message<GetMembersResponse> {
  /**
   * @generated from field: repeated admin.v1.User members = 1;
   */
  members: User[] = [];

  constructor(data?: PartialMessage<GetMembersResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.GetMembersResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "members", kind: "message", T: User, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetMembersResponse {
    return new GetMembersResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetMembersResponse {
    return new GetMembersResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetMembersResponse {
    return new GetMembersResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetMembersResponse | PlainMessage<GetMembersResponse> | undefined, b: GetMembersResponse | PlainMessage<GetMembersResponse> | undefined): boolean {
    return proto3.util.equals(GetMembersResponse, a, b);
  }
}

/**
 * @generated from message admin.v1.PatchTeamRequest
 */
export class PatchTeamRequest extends Message<PatchTeamRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: optional int32 code = 2;
   */
  code?: number;

  /**
   * @generated from field: optional string name = 3;
   */
  name?: string;

  /**
   * @generated from field: optional string organization = 4;
   */
  organization?: string;

  /**
   * @generated from field: optional int32 code_remaining = 6;
   */
  codeRemaining?: number;

  constructor(data?: PartialMessage<PatchTeamRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.PatchTeamRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "code", kind: "scalar", T: 5 /* ScalarType.INT32 */, opt: true },
    { no: 3, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 4, name: "organization", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 6, name: "code_remaining", kind: "scalar", T: 5 /* ScalarType.INT32 */, opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PatchTeamRequest {
    return new PatchTeamRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PatchTeamRequest {
    return new PatchTeamRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PatchTeamRequest {
    return new PatchTeamRequest().fromJsonString(jsonString, options);
  }

  static equals(a: PatchTeamRequest | PlainMessage<PatchTeamRequest> | undefined, b: PatchTeamRequest | PlainMessage<PatchTeamRequest> | undefined): boolean {
    return proto3.util.equals(PatchTeamRequest, a, b);
  }
}

/**
 * @generated from message admin.v1.PatchTeamResponse
 */
export class PatchTeamResponse extends Message<PatchTeamResponse> {
  /**
   * @generated from field: admin.v1.Team team = 1;
   */
  team?: Team;

  constructor(data?: PartialMessage<PatchTeamResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.PatchTeamResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "team", kind: "message", T: Team },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PatchTeamResponse {
    return new PatchTeamResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PatchTeamResponse {
    return new PatchTeamResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PatchTeamResponse {
    return new PatchTeamResponse().fromJsonString(jsonString, options);
  }

  static equals(a: PatchTeamResponse | PlainMessage<PatchTeamResponse> | undefined, b: PatchTeamResponse | PlainMessage<PatchTeamResponse> | undefined): boolean {
    return proto3.util.equals(PatchTeamResponse, a, b);
  }
}

/**
 * @generated from message admin.v1.PutTeamConnectionInfoRequest
 */
export class PutTeamConnectionInfoRequest extends Message<PutTeamConnectionInfoRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: admin.v1.BaStation ba_station = 2;
   */
  baStation?: BaStation;

  constructor(data?: PartialMessage<PutTeamConnectionInfoRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.PutTeamConnectionInfoRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "ba_station", kind: "message", T: BaStation },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PutTeamConnectionInfoRequest {
    return new PutTeamConnectionInfoRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PutTeamConnectionInfoRequest {
    return new PutTeamConnectionInfoRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PutTeamConnectionInfoRequest {
    return new PutTeamConnectionInfoRequest().fromJsonString(jsonString, options);
  }

  static equals(a: PutTeamConnectionInfoRequest | PlainMessage<PutTeamConnectionInfoRequest> | undefined, b: PutTeamConnectionInfoRequest | PlainMessage<PutTeamConnectionInfoRequest> | undefined): boolean {
    return proto3.util.equals(PutTeamConnectionInfoRequest, a, b);
  }
}

/**
 * @generated from message admin.v1.PutTeamConnectionInfoResponse
 */
export class PutTeamConnectionInfoResponse extends Message<PutTeamConnectionInfoResponse> {
  constructor(data?: PartialMessage<PutTeamConnectionInfoResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.PutTeamConnectionInfoResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PutTeamConnectionInfoResponse {
    return new PutTeamConnectionInfoResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PutTeamConnectionInfoResponse {
    return new PutTeamConnectionInfoResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PutTeamConnectionInfoResponse {
    return new PutTeamConnectionInfoResponse().fromJsonString(jsonString, options);
  }

  static equals(a: PutTeamConnectionInfoResponse | PlainMessage<PutTeamConnectionInfoResponse> | undefined, b: PutTeamConnectionInfoResponse | PlainMessage<PutTeamConnectionInfoResponse> | undefined): boolean {
    return proto3.util.equals(PutTeamConnectionInfoResponse, a, b);
  }
}

/**
 * @generated from message admin.v1.PostTeamRequest
 */
export class PostTeamRequest extends Message<PostTeamRequest> {
  /**
   * @generated from field: int32 code = 1;
   */
  code = 0;

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: string organization = 3;
   */
  organization = "";

  /**
   * @generated from field: int32 code_remaining = 6;
   */
  codeRemaining = 0;

  constructor(data?: PartialMessage<PostTeamRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.PostTeamRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "code", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "organization", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "code_remaining", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PostTeamRequest {
    return new PostTeamRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PostTeamRequest {
    return new PostTeamRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PostTeamRequest {
    return new PostTeamRequest().fromJsonString(jsonString, options);
  }

  static equals(a: PostTeamRequest | PlainMessage<PostTeamRequest> | undefined, b: PostTeamRequest | PlainMessage<PostTeamRequest> | undefined): boolean {
    return proto3.util.equals(PostTeamRequest, a, b);
  }
}

/**
 * @generated from message admin.v1.PostTeamResponse
 */
export class PostTeamResponse extends Message<PostTeamResponse> {
  /**
   * @generated from field: admin.v1.Team team = 1;
   */
  team?: Team;

  constructor(data?: PartialMessage<PostTeamResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.PostTeamResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "team", kind: "message", T: Team },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PostTeamResponse {
    return new PostTeamResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PostTeamResponse {
    return new PostTeamResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PostTeamResponse {
    return new PostTeamResponse().fromJsonString(jsonString, options);
  }

  static equals(a: PostTeamResponse | PlainMessage<PostTeamResponse> | undefined, b: PostTeamResponse | PlainMessage<PostTeamResponse> | undefined): boolean {
    return proto3.util.equals(PostTeamResponse, a, b);
  }
}

/**
 * @generated from message admin.v1.AddMemberRequest
 */
export class AddMemberRequest extends Message<AddMemberRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string user_id = 2;
   */
  userId = "";

  constructor(data?: PartialMessage<AddMemberRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.AddMemberRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "user_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AddMemberRequest {
    return new AddMemberRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AddMemberRequest {
    return new AddMemberRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AddMemberRequest {
    return new AddMemberRequest().fromJsonString(jsonString, options);
  }

  static equals(a: AddMemberRequest | PlainMessage<AddMemberRequest> | undefined, b: AddMemberRequest | PlainMessage<AddMemberRequest> | undefined): boolean {
    return proto3.util.equals(AddMemberRequest, a, b);
  }
}

/**
 * @generated from message admin.v1.AddMemberResponse
 */
export class AddMemberResponse extends Message<AddMemberResponse> {
  constructor(data?: PartialMessage<AddMemberResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.AddMemberResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AddMemberResponse {
    return new AddMemberResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AddMemberResponse {
    return new AddMemberResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AddMemberResponse {
    return new AddMemberResponse().fromJsonString(jsonString, options);
  }

  static equals(a: AddMemberResponse | PlainMessage<AddMemberResponse> | undefined, b: AddMemberResponse | PlainMessage<AddMemberResponse> | undefined): boolean {
    return proto3.util.equals(AddMemberResponse, a, b);
  }
}

/**
 * @generated from message admin.v1.DeleteTeamRequest
 */
export class DeleteTeamRequest extends Message<DeleteTeamRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  constructor(data?: PartialMessage<DeleteTeamRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.DeleteTeamRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeleteTeamRequest {
    return new DeleteTeamRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeleteTeamRequest {
    return new DeleteTeamRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeleteTeamRequest {
    return new DeleteTeamRequest().fromJsonString(jsonString, options);
  }

  static equals(a: DeleteTeamRequest | PlainMessage<DeleteTeamRequest> | undefined, b: DeleteTeamRequest | PlainMessage<DeleteTeamRequest> | undefined): boolean {
    return proto3.util.equals(DeleteTeamRequest, a, b);
  }
}

/**
 * @generated from message admin.v1.DeleteTeamResponse
 */
export class DeleteTeamResponse extends Message<DeleteTeamResponse> {
  constructor(data?: PartialMessage<DeleteTeamResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.DeleteTeamResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeleteTeamResponse {
    return new DeleteTeamResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeleteTeamResponse {
    return new DeleteTeamResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeleteTeamResponse {
    return new DeleteTeamResponse().fromJsonString(jsonString, options);
  }

  static equals(a: DeleteTeamResponse | PlainMessage<DeleteTeamResponse> | undefined, b: DeleteTeamResponse | PlainMessage<DeleteTeamResponse> | undefined): boolean {
    return proto3.util.equals(DeleteTeamResponse, a, b);
  }
}

/**
 * @generated from message admin.v1.DeleteMemberRequest
 */
export class DeleteMemberRequest extends Message<DeleteMemberRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string user_id = 2;
   */
  userId = "";

  constructor(data?: PartialMessage<DeleteMemberRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.DeleteMemberRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "user_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeleteMemberRequest {
    return new DeleteMemberRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeleteMemberRequest {
    return new DeleteMemberRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeleteMemberRequest {
    return new DeleteMemberRequest().fromJsonString(jsonString, options);
  }

  static equals(a: DeleteMemberRequest | PlainMessage<DeleteMemberRequest> | undefined, b: DeleteMemberRequest | PlainMessage<DeleteMemberRequest> | undefined): boolean {
    return proto3.util.equals(DeleteMemberRequest, a, b);
  }
}

/**
 * @generated from message admin.v1.DeleteMemberResponse
 */
export class DeleteMemberResponse extends Message<DeleteMemberResponse> {
  constructor(data?: PartialMessage<DeleteMemberResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "admin.v1.DeleteMemberResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeleteMemberResponse {
    return new DeleteMemberResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeleteMemberResponse {
    return new DeleteMemberResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeleteMemberResponse {
    return new DeleteMemberResponse().fromJsonString(jsonString, options);
  }

  static equals(a: DeleteMemberResponse | PlainMessage<DeleteMemberResponse> | undefined, b: DeleteMemberResponse | PlainMessage<DeleteMemberResponse> | undefined): boolean {
    return proto3.util.equals(DeleteMemberResponse, a, b);
  }
}

