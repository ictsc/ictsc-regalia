// @generated by protoc-gen-es v2.1.0 with parameter "target=ts"
// @generated from file admin/v1/contestant.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/contestant.proto.
 */
export const file_admin_v1_contestant: GenFile = /*@__PURE__*/
  fileDesc("ChlhZG1pbi92MS9jb250ZXN0YW50LnByb3RvEghhZG1pbi52MSJWCgpDb250ZXN0YW50EhQKAmlkGAEgASgJQgi6SAVyA5gBGhIXCgRuYW1lGAIgASgJQgm6SAZyBBABGBQSGQoHdGVhbV9pZBgDIAEoCUIIukgFcgOYARoiLAoUR2V0Q29udGVzdGFudFJlcXVlc3QSFAoCaWQYASABKAlCCLpIBXIDmAEaIkMKFUdldENvbnRlc3RhbnRSZXNwb25zZRIqCgR1c2VyGAEgASgLMhQuYWRtaW4udjEuQ29udGVzdGFudEIGukgDyAEBIhcKFUdldENvbnRlc3RhbnRzUmVxdWVzdCJFChZHZXRDb250ZXN0YW50c1Jlc3BvbnNlEisKBXVzZXJzGAEgAygLMhQuYWRtaW4udjEuQ29udGVzdGFudEIGukgDyAEBIi8KF0RlbGV0ZUNvbnRlc3RhbnRSZXF1ZXN0EhQKAmlkGAEgASgJQgi6SAVyA5gBGiIaChhEZWxldGVDb250ZXN0YW50UmVzcG9uc2UymwIKEUNvbnRlc3RhbnRTZXJ2aWNlElIKDUdldENvbnRlc3RhbnQSHi5hZG1pbi52MS5HZXRDb250ZXN0YW50UmVxdWVzdBofLmFkbWluLnYxLkdldENvbnRlc3RhbnRSZXNwb25zZSIAElUKDkdldENvbnRlc3RhbnRzEh8uYWRtaW4udjEuR2V0Q29udGVzdGFudHNSZXF1ZXN0GiAuYWRtaW4udjEuR2V0Q29udGVzdGFudHNSZXNwb25zZSIAElsKEERlbGV0ZUNvbnRlc3RhbnQSIS5hZG1pbi52MS5EZWxldGVDb250ZXN0YW50UmVxdWVzdBoiLmFkbWluLnYxLkRlbGV0ZUNvbnRlc3RhbnRSZXNwb25zZSIAQqEBCgxjb20uYWRtaW4udjFCD0NvbnRlc3RhbnRQcm90b1ABWj9naXRodWIuY29tL2ljdHNjL2ljdHNjLW91dGxhbmRzL2JhY2tlbmQvaW50ZXJuYWwvcHJvdG8vYWRtaW4vdjGiAgNBWFiqAghBZG1pbi5WMcoCCEFkbWluXFYx4gIUQWRtaW5cVjFcR1BCTWV0YWRhdGHqAglBZG1pbjo6VjFiBnByb3RvMw", [file_buf_validate_validate]);

/**
 * @generated from message admin.v1.Contestant
 */
export type Contestant = Message<"admin.v1.Contestant"> & {
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
 * Describes the message admin.v1.Contestant.
 * Use `create(ContestantSchema)` to create a new message.
 */
export const ContestantSchema: GenMessage<Contestant> = /*@__PURE__*/
  messageDesc(file_admin_v1_contestant, 0);

/**
 * @generated from message admin.v1.GetContestantRequest
 */
export type GetContestantRequest = Message<"admin.v1.GetContestantRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message admin.v1.GetContestantRequest.
 * Use `create(GetContestantRequestSchema)` to create a new message.
 */
export const GetContestantRequestSchema: GenMessage<GetContestantRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_contestant, 1);

/**
 * @generated from message admin.v1.GetContestantResponse
 */
export type GetContestantResponse = Message<"admin.v1.GetContestantResponse"> & {
  /**
   * @generated from field: admin.v1.Contestant user = 1;
   */
  user?: Contestant;
};

/**
 * Describes the message admin.v1.GetContestantResponse.
 * Use `create(GetContestantResponseSchema)` to create a new message.
 */
export const GetContestantResponseSchema: GenMessage<GetContestantResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_contestant, 2);

/**
 * @generated from message admin.v1.GetContestantsRequest
 */
export type GetContestantsRequest = Message<"admin.v1.GetContestantsRequest"> & {
};

/**
 * Describes the message admin.v1.GetContestantsRequest.
 * Use `create(GetContestantsRequestSchema)` to create a new message.
 */
export const GetContestantsRequestSchema: GenMessage<GetContestantsRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_contestant, 3);

/**
 * @generated from message admin.v1.GetContestantsResponse
 */
export type GetContestantsResponse = Message<"admin.v1.GetContestantsResponse"> & {
  /**
   * @generated from field: repeated admin.v1.Contestant users = 1;
   */
  users: Contestant[];
};

/**
 * Describes the message admin.v1.GetContestantsResponse.
 * Use `create(GetContestantsResponseSchema)` to create a new message.
 */
export const GetContestantsResponseSchema: GenMessage<GetContestantsResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_contestant, 4);

/**
 * @generated from message admin.v1.DeleteContestantRequest
 */
export type DeleteContestantRequest = Message<"admin.v1.DeleteContestantRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message admin.v1.DeleteContestantRequest.
 * Use `create(DeleteContestantRequestSchema)` to create a new message.
 */
export const DeleteContestantRequestSchema: GenMessage<DeleteContestantRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_contestant, 5);

/**
 * @generated from message admin.v1.DeleteContestantResponse
 */
export type DeleteContestantResponse = Message<"admin.v1.DeleteContestantResponse"> & {
};

/**
 * Describes the message admin.v1.DeleteContestantResponse.
 * Use `create(DeleteContestantResponseSchema)` to create a new message.
 */
export const DeleteContestantResponseSchema: GenMessage<DeleteContestantResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_contestant, 6);

/**
 * @generated from service admin.v1.ContestantService
 */
export const ContestantService: GenService<{
  /**
   * @generated from rpc admin.v1.ContestantService.GetContestant
   */
  getContestant: {
    methodKind: "unary";
    input: typeof GetContestantRequestSchema;
    output: typeof GetContestantResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.ContestantService.GetContestants
   */
  getContestants: {
    methodKind: "unary";
    input: typeof GetContestantsRequestSchema;
    output: typeof GetContestantsResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.ContestantService.DeleteContestant
   */
  deleteContestant: {
    methodKind: "unary";
    input: typeof DeleteContestantRequestSchema;
    output: typeof DeleteContestantResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_contestant, 0);

