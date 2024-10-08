// @generated by protoc-gen-es v2.2.0 with parameter "target=ts"
// @generated from file admin/v1/team.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Contestant } from "./contestant_pb";
import { file_admin_v1_contestant } from "./contestant_pb";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/team.proto.
 */
export const file_admin_v1_team: GenFile = /*@__PURE__*/
  fileDesc("ChNhZG1pbi92MS90ZWFtLnByb3RvEghhZG1pbi52MSK1AQoEVGVhbRIUCgJpZBgBIAEoCUIIukgFcgOYARoSFwoEY29kZRgCIAEoA0IJukgGIgQQZCABEhcKBG5hbWUYAyABKAlCCbpIBnIEEAEYFBIfCgxvcmdhbml6YXRpb24YBCABKAlCCbpIBnIEEAEYMhIhCg9pbnZpdGF0aW9uX2NvZGUYBSABKAlCCLpIBXIDmAEgEiEKDmNvZGVfcmVtYWluaW5nGAYgASgDQgm6SAYiBBgFKAAiEQoPR2V0VGVhbXNSZXF1ZXN0IjkKEEdldFRlYW1zUmVzcG9uc2USJQoFdGVhbXMYASADKAsyDi5hZG1pbi52MS5UZWFtQga6SAPIAQEiJgoOR2V0VGVhbVJlcXVlc3QSFAoCaWQYASABKAlCCLpIBXIDmAEaIjcKD0dldFRlYW1SZXNwb25zZRIkCgR0ZWFtGAEgASgLMg4uYWRtaW4udjEuVGVhbUIGukgDyAEBIm8KB0Jhc3Rpb24SFwoEdXNlchgBIAEoCUIJukgGcgQQARgUEhsKCHBhc3N3b3JkGAIgASgJQgm6SAZyBBABGBQSFwoEaG9zdBgDIAEoCUIJukgGcgQQARhkEhUKBHBvcnQYBCABKANCB7pIBCICKAAiMAoYR2V0Q29ubmVjdGlvbkluZm9SZXF1ZXN0EhQKAmlkGAEgASgJQgi6SAVyA5gBGiI/ChlHZXRDb25uZWN0aW9uSW5mb1Jlc3BvbnNlEiIKB2Jhc3Rpb24YASABKAsyES5hZG1pbi52MS5CYXN0aW9uIikKEUdldE1lbWJlcnNSZXF1ZXN0EhQKAmlkGAEgASgJQgi6SAVyA5gBGiJHChJHZXRNZW1iZXJzUmVzcG9uc2USMQoHbWVtYmVycxgBIAMoCzIULmFkbWluLnYxLkNvbnRlc3RhbnRCCrpIB5IBBAgAEAUi6AEKEFBhdGNoVGVhbVJlcXVlc3QSFAoCaWQYASABKAlCCLpIBXIDmAEaEhwKBGNvZGUYAiABKANCCbpIBiIEEGQgAUgAiAEBEhwKBG5hbWUYAyABKAlCCbpIBnIEEAEYFEgBiAEBEiQKDG9yZ2FuaXphdGlvbhgEIAEoCUIJukgGcgQQARgySAKIAQESJgoOY29kZV9yZW1haW5pbmcYBSABKANCCbpIBiIEGAUoAEgDiAEBQgcKBV9jb2RlQgcKBV9uYW1lQg8KDV9vcmdhbml6YXRpb25CEQoPX2NvZGVfcmVtYWluaW5nIjkKEVBhdGNoVGVhbVJlc3BvbnNlEiQKBHRlYW0YASABKAsyDi5hZG1pbi52MS5UZWFtQga6SAPIAQEiXAoYUHV0Q29ubmVjdGlvbkluZm9SZXF1ZXN0EhQKAmlkGAEgASgJQgi6SAVyA5gBGhIqCgdiYXN0aW9uGAIgASgLMhEuYWRtaW4udjEuQmFzdGlvbkIGukgDyAEBIhsKGVB1dENvbm5lY3Rpb25JbmZvUmVzcG9uc2UihwEKD1Bvc3RUZWFtUmVxdWVzdBIXCgRjb2RlGAEgASgDQgm6SAYiBBBkIAESFwoEbmFtZRgCIAEoCUIJukgGcgQQARgUEh8KDG9yZ2FuaXphdGlvbhgDIAEoCUIJukgGcgQQARgyEiEKDmNvZGVfcmVtYWluaW5nGAQgASgDQgm6SAYiBBgFKAAiOAoQUG9zdFRlYW1SZXNwb25zZRIkCgR0ZWFtGAEgASgLMg4uYWRtaW4udjEuVGVhbUIGukgDyAEBIkMKEEFkZE1lbWJlclJlcXVlc3QSFAoCaWQYASABKAlCCLpIBXIDmAEaEhkKB3VzZXJfaWQYAiABKAlCCLpIBXIDmAEaIhMKEUFkZE1lbWJlclJlc3BvbnNlIikKEURlbGV0ZVRlYW1SZXF1ZXN0EhQKAmlkGAEgASgJQgi6SAVyA5gBGiIUChJEZWxldGVUZWFtUmVzcG9uc2UiTAoRTW92ZU1lbWJlclJlcXVlc3QSHAoKdG9fdGVhbV9pZBgBIAEoCUIIukgFcgOYARoSGQoHdXNlcl9pZBgCIAEoCUIIukgFcgOYARoiFAoSTW92ZU1lbWJlclJlc3BvbnNlMrAFCgtUZWFtU2VydmljZRJBCghHZXRUZWFtcxIZLmFkbWluLnYxLkdldFRlYW1zUmVxdWVzdBoaLmFkbWluLnYxLkdldFRlYW1zUmVzcG9uc2USPgoHR2V0VGVhbRIYLmFkbWluLnYxLkdldFRlYW1SZXF1ZXN0GhkuYWRtaW4udjEuR2V0VGVhbVJlc3BvbnNlElwKEUdldENvbm5lY3Rpb25JbmZvEiIuYWRtaW4udjEuR2V0Q29ubmVjdGlvbkluZm9SZXF1ZXN0GiMuYWRtaW4udjEuR2V0Q29ubmVjdGlvbkluZm9SZXNwb25zZRJHCgpHZXRNZW1iZXJzEhsuYWRtaW4udjEuR2V0TWVtYmVyc1JlcXVlc3QaHC5hZG1pbi52MS5HZXRNZW1iZXJzUmVzcG9uc2USRAoJUGF0Y2hUZWFtEhouYWRtaW4udjEuUGF0Y2hUZWFtUmVxdWVzdBobLmFkbWluLnYxLlBhdGNoVGVhbVJlc3BvbnNlElwKEVB1dENvbm5lY3Rpb25JbmZvEiIuYWRtaW4udjEuUHV0Q29ubmVjdGlvbkluZm9SZXF1ZXN0GiMuYWRtaW4udjEuUHV0Q29ubmVjdGlvbkluZm9SZXNwb25zZRJBCghQb3N0VGVhbRIZLmFkbWluLnYxLlBvc3RUZWFtUmVxdWVzdBoaLmFkbWluLnYxLlBvc3RUZWFtUmVzcG9uc2USRwoKRGVsZXRlVGVhbRIbLmFkbWluLnYxLkRlbGV0ZVRlYW1SZXF1ZXN0GhwuYWRtaW4udjEuRGVsZXRlVGVhbVJlc3BvbnNlEkcKCk1vdmVNZW1iZXISGy5hZG1pbi52MS5Nb3ZlTWVtYmVyUmVxdWVzdBocLmFkbWluLnYxLk1vdmVNZW1iZXJSZXNwb25zZUKbAQoMY29tLmFkbWluLnYxQglUZWFtUHJvdG9QAVo/Z2l0aHViLmNvbS9pY3RzYy9pY3RzYy1vdXRsYW5kcy9iYWNrZW5kL2ludGVybmFsL3Byb3RvL2FkbWluL3YxogIDQVhYqgIIQWRtaW4uVjHKAghBZG1pblxWMeICFEFkbWluXFYxXEdQQk1ldGFkYXRh6gIJQWRtaW46OlYxYgZwcm90bzM", [file_admin_v1_contestant, file_buf_validate_validate]);

/**
 * @generated from message admin.v1.Team
 */
export type Team = Message<"admin.v1.Team"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: int64 code = 2;
   */
  code: bigint;

  /**
   * @generated from field: string name = 3;
   */
  name: string;

  /**
   * @generated from field: string organization = 4;
   */
  organization: string;

  /**
   * @generated from field: string invitation_code = 5;
   */
  invitationCode: string;

  /**
   * @generated from field: int64 code_remaining = 6;
   */
  codeRemaining: bigint;
};

/**
 * Describes the message admin.v1.Team.
 * Use `create(TeamSchema)` to create a new message.
 */
export const TeamSchema: GenMessage<Team> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 0);

/**
 * @generated from message admin.v1.GetTeamsRequest
 */
export type GetTeamsRequest = Message<"admin.v1.GetTeamsRequest"> & {
};

/**
 * Describes the message admin.v1.GetTeamsRequest.
 * Use `create(GetTeamsRequestSchema)` to create a new message.
 */
export const GetTeamsRequestSchema: GenMessage<GetTeamsRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 1);

/**
 * @generated from message admin.v1.GetTeamsResponse
 */
export type GetTeamsResponse = Message<"admin.v1.GetTeamsResponse"> & {
  /**
   * @generated from field: repeated admin.v1.Team teams = 1;
   */
  teams: Team[];
};

/**
 * Describes the message admin.v1.GetTeamsResponse.
 * Use `create(GetTeamsResponseSchema)` to create a new message.
 */
export const GetTeamsResponseSchema: GenMessage<GetTeamsResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 2);

/**
 * @generated from message admin.v1.GetTeamRequest
 */
export type GetTeamRequest = Message<"admin.v1.GetTeamRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message admin.v1.GetTeamRequest.
 * Use `create(GetTeamRequestSchema)` to create a new message.
 */
export const GetTeamRequestSchema: GenMessage<GetTeamRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 3);

/**
 * @generated from message admin.v1.GetTeamResponse
 */
export type GetTeamResponse = Message<"admin.v1.GetTeamResponse"> & {
  /**
   * @generated from field: admin.v1.Team team = 1;
   */
  team?: Team;
};

/**
 * Describes the message admin.v1.GetTeamResponse.
 * Use `create(GetTeamResponseSchema)` to create a new message.
 */
export const GetTeamResponseSchema: GenMessage<GetTeamResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 4);

/**
 * @generated from message admin.v1.Bastion
 */
export type Bastion = Message<"admin.v1.Bastion"> & {
  /**
   * @generated from field: string user = 1;
   */
  user: string;

  /**
   * @generated from field: string password = 2;
   */
  password: string;

  /**
   * @generated from field: string host = 3;
   */
  host: string;

  /**
   * @generated from field: int64 port = 4;
   */
  port: bigint;
};

/**
 * Describes the message admin.v1.Bastion.
 * Use `create(BastionSchema)` to create a new message.
 */
export const BastionSchema: GenMessage<Bastion> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 5);

/**
 * @generated from message admin.v1.GetConnectionInfoRequest
 */
export type GetConnectionInfoRequest = Message<"admin.v1.GetConnectionInfoRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message admin.v1.GetConnectionInfoRequest.
 * Use `create(GetConnectionInfoRequestSchema)` to create a new message.
 */
export const GetConnectionInfoRequestSchema: GenMessage<GetConnectionInfoRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 6);

/**
 * @generated from message admin.v1.GetConnectionInfoResponse
 */
export type GetConnectionInfoResponse = Message<"admin.v1.GetConnectionInfoResponse"> & {
  /**
   * @generated from field: admin.v1.Bastion bastion = 1;
   */
  bastion?: Bastion;
};

/**
 * Describes the message admin.v1.GetConnectionInfoResponse.
 * Use `create(GetConnectionInfoResponseSchema)` to create a new message.
 */
export const GetConnectionInfoResponseSchema: GenMessage<GetConnectionInfoResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 7);

/**
 * @generated from message admin.v1.GetMembersRequest
 */
export type GetMembersRequest = Message<"admin.v1.GetMembersRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message admin.v1.GetMembersRequest.
 * Use `create(GetMembersRequestSchema)` to create a new message.
 */
export const GetMembersRequestSchema: GenMessage<GetMembersRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 8);

/**
 * @generated from message admin.v1.GetMembersResponse
 */
export type GetMembersResponse = Message<"admin.v1.GetMembersResponse"> & {
  /**
   * @generated from field: repeated admin.v1.Contestant members = 1;
   */
  members: Contestant[];
};

/**
 * Describes the message admin.v1.GetMembersResponse.
 * Use `create(GetMembersResponseSchema)` to create a new message.
 */
export const GetMembersResponseSchema: GenMessage<GetMembersResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 9);

/**
 * @generated from message admin.v1.PatchTeamRequest
 */
export type PatchTeamRequest = Message<"admin.v1.PatchTeamRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: optional int64 code = 2;
   */
  code?: bigint;

  /**
   * @generated from field: optional string name = 3;
   */
  name?: string;

  /**
   * @generated from field: optional string organization = 4;
   */
  organization?: string;

  /**
   * @generated from field: optional int64 code_remaining = 5;
   */
  codeRemaining?: bigint;
};

/**
 * Describes the message admin.v1.PatchTeamRequest.
 * Use `create(PatchTeamRequestSchema)` to create a new message.
 */
export const PatchTeamRequestSchema: GenMessage<PatchTeamRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 10);

/**
 * @generated from message admin.v1.PatchTeamResponse
 */
export type PatchTeamResponse = Message<"admin.v1.PatchTeamResponse"> & {
  /**
   * @generated from field: admin.v1.Team team = 1;
   */
  team?: Team;
};

/**
 * Describes the message admin.v1.PatchTeamResponse.
 * Use `create(PatchTeamResponseSchema)` to create a new message.
 */
export const PatchTeamResponseSchema: GenMessage<PatchTeamResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 11);

/**
 * @generated from message admin.v1.PutConnectionInfoRequest
 */
export type PutConnectionInfoRequest = Message<"admin.v1.PutConnectionInfoRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: admin.v1.Bastion bastion = 2;
   */
  bastion?: Bastion;
};

/**
 * Describes the message admin.v1.PutConnectionInfoRequest.
 * Use `create(PutConnectionInfoRequestSchema)` to create a new message.
 */
export const PutConnectionInfoRequestSchema: GenMessage<PutConnectionInfoRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 12);

/**
 * @generated from message admin.v1.PutConnectionInfoResponse
 */
export type PutConnectionInfoResponse = Message<"admin.v1.PutConnectionInfoResponse"> & {
};

/**
 * Describes the message admin.v1.PutConnectionInfoResponse.
 * Use `create(PutConnectionInfoResponseSchema)` to create a new message.
 */
export const PutConnectionInfoResponseSchema: GenMessage<PutConnectionInfoResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 13);

/**
 * @generated from message admin.v1.PostTeamRequest
 */
export type PostTeamRequest = Message<"admin.v1.PostTeamRequest"> & {
  /**
   * @generated from field: int64 code = 1;
   */
  code: bigint;

  /**
   * @generated from field: string name = 2;
   */
  name: string;

  /**
   * @generated from field: string organization = 3;
   */
  organization: string;

  /**
   * @generated from field: int64 code_remaining = 4;
   */
  codeRemaining: bigint;
};

/**
 * Describes the message admin.v1.PostTeamRequest.
 * Use `create(PostTeamRequestSchema)` to create a new message.
 */
export const PostTeamRequestSchema: GenMessage<PostTeamRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 14);

/**
 * @generated from message admin.v1.PostTeamResponse
 */
export type PostTeamResponse = Message<"admin.v1.PostTeamResponse"> & {
  /**
   * @generated from field: admin.v1.Team team = 1;
   */
  team?: Team;
};

/**
 * Describes the message admin.v1.PostTeamResponse.
 * Use `create(PostTeamResponseSchema)` to create a new message.
 */
export const PostTeamResponseSchema: GenMessage<PostTeamResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 15);

/**
 * @generated from message admin.v1.AddMemberRequest
 */
export type AddMemberRequest = Message<"admin.v1.AddMemberRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string user_id = 2;
   */
  userId: string;
};

/**
 * Describes the message admin.v1.AddMemberRequest.
 * Use `create(AddMemberRequestSchema)` to create a new message.
 */
export const AddMemberRequestSchema: GenMessage<AddMemberRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 16);

/**
 * @generated from message admin.v1.AddMemberResponse
 */
export type AddMemberResponse = Message<"admin.v1.AddMemberResponse"> & {
};

/**
 * Describes the message admin.v1.AddMemberResponse.
 * Use `create(AddMemberResponseSchema)` to create a new message.
 */
export const AddMemberResponseSchema: GenMessage<AddMemberResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 17);

/**
 * @generated from message admin.v1.DeleteTeamRequest
 */
export type DeleteTeamRequest = Message<"admin.v1.DeleteTeamRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message admin.v1.DeleteTeamRequest.
 * Use `create(DeleteTeamRequestSchema)` to create a new message.
 */
export const DeleteTeamRequestSchema: GenMessage<DeleteTeamRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 18);

/**
 * @generated from message admin.v1.DeleteTeamResponse
 */
export type DeleteTeamResponse = Message<"admin.v1.DeleteTeamResponse"> & {
};

/**
 * Describes the message admin.v1.DeleteTeamResponse.
 * Use `create(DeleteTeamResponseSchema)` to create a new message.
 */
export const DeleteTeamResponseSchema: GenMessage<DeleteTeamResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 19);

/**
 * @generated from message admin.v1.MoveMemberRequest
 */
export type MoveMemberRequest = Message<"admin.v1.MoveMemberRequest"> & {
  /**
   * @generated from field: string to_team_id = 1;
   */
  toTeamId: string;

  /**
   * @generated from field: string user_id = 2;
   */
  userId: string;
};

/**
 * Describes the message admin.v1.MoveMemberRequest.
 * Use `create(MoveMemberRequestSchema)` to create a new message.
 */
export const MoveMemberRequestSchema: GenMessage<MoveMemberRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 20);

/**
 * @generated from message admin.v1.MoveMemberResponse
 */
export type MoveMemberResponse = Message<"admin.v1.MoveMemberResponse"> & {
};

/**
 * Describes the message admin.v1.MoveMemberResponse.
 * Use `create(MoveMemberResponseSchema)` to create a new message.
 */
export const MoveMemberResponseSchema: GenMessage<MoveMemberResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_team, 21);

/**
 * @generated from service admin.v1.TeamService
 */
export const TeamService: GenService<{
  /**
   * @generated from rpc admin.v1.TeamService.GetTeams
   */
  getTeams: {
    methodKind: "unary";
    input: typeof GetTeamsRequestSchema;
    output: typeof GetTeamsResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.TeamService.GetTeam
   */
  getTeam: {
    methodKind: "unary";
    input: typeof GetTeamRequestSchema;
    output: typeof GetTeamResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.TeamService.GetConnectionInfo
   */
  getConnectionInfo: {
    methodKind: "unary";
    input: typeof GetConnectionInfoRequestSchema;
    output: typeof GetConnectionInfoResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.TeamService.GetMembers
   */
  getMembers: {
    methodKind: "unary";
    input: typeof GetMembersRequestSchema;
    output: typeof GetMembersResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.TeamService.PatchTeam
   */
  patchTeam: {
    methodKind: "unary";
    input: typeof PatchTeamRequestSchema;
    output: typeof PatchTeamResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.TeamService.PutConnectionInfo
   */
  putConnectionInfo: {
    methodKind: "unary";
    input: typeof PutConnectionInfoRequestSchema;
    output: typeof PutConnectionInfoResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.TeamService.PostTeam
   */
  postTeam: {
    methodKind: "unary";
    input: typeof PostTeamRequestSchema;
    output: typeof PostTeamResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.TeamService.DeleteTeam
   */
  deleteTeam: {
    methodKind: "unary";
    input: typeof DeleteTeamRequestSchema;
    output: typeof DeleteTeamResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.TeamService.MoveMember
   */
  moveMember: {
    methodKind: "unary";
    input: typeof MoveMemberRequestSchema;
    output: typeof MoveMemberResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_team, 0);

