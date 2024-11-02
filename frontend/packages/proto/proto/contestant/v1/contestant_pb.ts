// @generated by protoc-gen-es v2.2.2 with parameter "target=ts"
// @generated from file contestant/v1/contestant.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/contestant.proto.
 */
export const file_contestant_v1_contestant: GenFile = /*@__PURE__*/
  fileDesc("Ch5jb250ZXN0YW50L3YxL2NvbnRlc3RhbnQucHJvdG8SDWNvbnRlc3RhbnQudjEiVgoKQ29udGVzdGFudBIUCgJpZBgBIAEoCUIIukgFcgOYARoSFwoEbmFtZRgCIAEoCUIJukgGcgQQARgUEhkKB3RlYW1faWQYAyABKAlCCLpIBXIDmAEaIg4KDEdldE1lUmVxdWVzdCJACg1HZXRNZVJlc3BvbnNlEi8KBHVzZXIYASABKAsyGS5jb250ZXN0YW50LnYxLkNvbnRlc3RhbnRCBrpIA8gBASIsChRHZXRDb250ZXN0YW50UmVxdWVzdBIUCgJpZBgBIAEoCUIIukgFcgOYARoiSAoVR2V0Q29udGVzdGFudFJlc3BvbnNlEi8KBHVzZXIYASABKAsyGS5jb250ZXN0YW50LnYxLkNvbnRlc3RhbnRCBrpIA8gBASJTChVQb3N0Q29udGVzdGFudFJlcXVlc3QSFwoEbmFtZRgBIAEoCUIJukgGcgQQARgUEiEKD2ludml0YXRpb25fY29kZRgCIAEoCUIIukgFcgOYASAiSQoWUG9zdENvbnRlc3RhbnRSZXNwb25zZRIvCgR1c2VyGAEgASgLMhkuY29udGVzdGFudC52MS5Db250ZXN0YW50Qga6SAPIAQEiKQoOUGF0Y2hNZVJlcXVlc3QSFwoEbmFtZRgBIAEoCUIJukgGcgQQARgUIkIKD1BhdGNoTWVSZXNwb25zZRIvCgR1c2VyGAEgASgLMhkuY29udGVzdGFudC52MS5Db250ZXN0YW50Qga6SAPIAQEy5AIKEUNvbnRlc3RhbnRTZXJ2aWNlEkQKBUdldE1lEhsuY29udGVzdGFudC52MS5HZXRNZVJlcXVlc3QaHC5jb250ZXN0YW50LnYxLkdldE1lUmVzcG9uc2UiABJcCg1HZXRDb250ZXN0YW50EiMuY29udGVzdGFudC52MS5HZXRDb250ZXN0YW50UmVxdWVzdBokLmNvbnRlc3RhbnQudjEuR2V0Q29udGVzdGFudFJlc3BvbnNlIgASXwoOUG9zdENvbnRlc3RhbnQSJC5jb250ZXN0YW50LnYxLlBvc3RDb250ZXN0YW50UmVxdWVzdBolLmNvbnRlc3RhbnQudjEuUG9zdENvbnRlc3RhbnRSZXNwb25zZSIAEkoKB1BhdGNoTWUSHS5jb250ZXN0YW50LnYxLlBhdGNoTWVSZXF1ZXN0Gh4uY29udGVzdGFudC52MS5QYXRjaE1lUmVzcG9uc2UiAEK/AQoRY29tLmNvbnRlc3RhbnQudjFCD0NvbnRlc3RhbnRQcm90b1ABWkRnaXRodWIuY29tL2ljdHNjL2ljdHNjLW91dGxhbmRzL2JhY2tlbmQvaW50ZXJuYWwvcHJvdG8vY29udGVzdGFudC92MaICA0NYWKoCDUNvbnRlc3RhbnQuVjHKAg1Db250ZXN0YW50XFYx4gIZQ29udGVzdGFudFxWMVxHUEJNZXRhZGF0YeoCDkNvbnRlc3RhbnQ6OlYxYgZwcm90bzM", [file_buf_validate_validate]);

/**
 * @generated from message contestant.v1.Contestant
 */
export type Contestant = Message<"contestant.v1.Contestant"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string name = 2;
   */
  name: string;

  /**
   * @generated from field: string team_id = 3;
   */
  teamId: string;
};

/**
 * Describes the message contestant.v1.Contestant.
 * Use `create(ContestantSchema)` to create a new message.
 */
export const ContestantSchema: GenMessage<Contestant> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contestant, 0);

/**
 * @generated from message contestant.v1.GetMeRequest
 */
export type GetMeRequest = Message<"contestant.v1.GetMeRequest"> & {
};

/**
 * Describes the message contestant.v1.GetMeRequest.
 * Use `create(GetMeRequestSchema)` to create a new message.
 */
export const GetMeRequestSchema: GenMessage<GetMeRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contestant, 1);

/**
 * @generated from message contestant.v1.GetMeResponse
 */
export type GetMeResponse = Message<"contestant.v1.GetMeResponse"> & {
  /**
   * @generated from field: contestant.v1.Contestant user = 1;
   */
  user?: Contestant;
};

/**
 * Describes the message contestant.v1.GetMeResponse.
 * Use `create(GetMeResponseSchema)` to create a new message.
 */
export const GetMeResponseSchema: GenMessage<GetMeResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contestant, 2);

/**
 * @generated from message contestant.v1.GetContestantRequest
 */
export type GetContestantRequest = Message<"contestant.v1.GetContestantRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message contestant.v1.GetContestantRequest.
 * Use `create(GetContestantRequestSchema)` to create a new message.
 */
export const GetContestantRequestSchema: GenMessage<GetContestantRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contestant, 3);

/**
 * @generated from message contestant.v1.GetContestantResponse
 */
export type GetContestantResponse = Message<"contestant.v1.GetContestantResponse"> & {
  /**
   * @generated from field: contestant.v1.Contestant user = 1;
   */
  user?: Contestant;
};

/**
 * Describes the message contestant.v1.GetContestantResponse.
 * Use `create(GetContestantResponseSchema)` to create a new message.
 */
export const GetContestantResponseSchema: GenMessage<GetContestantResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contestant, 4);

/**
 * @generated from message contestant.v1.PostContestantRequest
 */
export type PostContestantRequest = Message<"contestant.v1.PostContestantRequest"> & {
  /**
   * @generated from field: string name = 1;
   */
  name: string;

  /**
   * @generated from field: string invitation_code = 2;
   */
  invitationCode: string;
};

/**
 * Describes the message contestant.v1.PostContestantRequest.
 * Use `create(PostContestantRequestSchema)` to create a new message.
 */
export const PostContestantRequestSchema: GenMessage<PostContestantRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contestant, 5);

/**
 * @generated from message contestant.v1.PostContestantResponse
 */
export type PostContestantResponse = Message<"contestant.v1.PostContestantResponse"> & {
  /**
   * @generated from field: contestant.v1.Contestant user = 1;
   */
  user?: Contestant;
};

/**
 * Describes the message contestant.v1.PostContestantResponse.
 * Use `create(PostContestantResponseSchema)` to create a new message.
 */
export const PostContestantResponseSchema: GenMessage<PostContestantResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contestant, 6);

/**
 * @generated from message contestant.v1.PatchMeRequest
 */
export type PatchMeRequest = Message<"contestant.v1.PatchMeRequest"> & {
  /**
   * @generated from field: string name = 1;
   */
  name: string;
};

/**
 * Describes the message contestant.v1.PatchMeRequest.
 * Use `create(PatchMeRequestSchema)` to create a new message.
 */
export const PatchMeRequestSchema: GenMessage<PatchMeRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contestant, 7);

/**
 * @generated from message contestant.v1.PatchMeResponse
 */
export type PatchMeResponse = Message<"contestant.v1.PatchMeResponse"> & {
  /**
   * @generated from field: contestant.v1.Contestant user = 1;
   */
  user?: Contestant;
};

/**
 * Describes the message contestant.v1.PatchMeResponse.
 * Use `create(PatchMeResponseSchema)` to create a new message.
 */
export const PatchMeResponseSchema: GenMessage<PatchMeResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contestant, 8);

/**
 * @generated from service contestant.v1.ContestantService
 */
export const ContestantService: GenService<{
  /**
   * @generated from rpc contestant.v1.ContestantService.GetMe
   */
  getMe: {
    methodKind: "unary";
    input: typeof GetMeRequestSchema;
    output: typeof GetMeResponseSchema;
  },
  /**
   * @generated from rpc contestant.v1.ContestantService.GetContestant
   */
  getContestant: {
    methodKind: "unary";
    input: typeof GetContestantRequestSchema;
    output: typeof GetContestantResponseSchema;
  },
  /**
   * @generated from rpc contestant.v1.ContestantService.PostContestant
   */
  postContestant: {
    methodKind: "unary";
    input: typeof PostContestantRequestSchema;
    output: typeof PostContestantResponseSchema;
  },
  /**
   * @generated from rpc contestant.v1.ContestantService.PatchMe
   */
  patchMe: {
    methodKind: "unary";
    input: typeof PatchMeRequestSchema;
    output: typeof PatchMeResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_contestant_v1_contestant, 0);

