// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file admin/v1/invitation.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/invitation.proto.
 */
export const file_admin_v1_invitation: GenFile = /*@__PURE__*/
  fileDesc("ChlhZG1pbi92MS9pbnZpdGF0aW9uLnByb3RvEghhZG1pbi52MSK9AQoOSW52aXRhdGlvbkNvZGUSDAoEY29kZRgBIAEoCRIRCgl0ZWFtX2NvZGUYAiABKAMSEgoKdG90YWxfdXNlcxgDIAEoBBIWCg5yZW1haW5pbmdfdXNlcxgEIAEoBBIuCgpjcmVhdGVkX2F0GAUgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcBIuCgpleHBpcmVzX2F0GAYgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcCJIChpMaXN0SW52aXRhdGlvbkNvZGVzUmVxdWVzdBIRCgl0ZWFtX2NvZGUYASABKAMSFwoPaW5jbHVkZV9leHBpcmVkGAIgASgIIlEKG0xpc3RJbnZpdGF0aW9uQ29kZXNSZXNwb25zZRIyChBpbnZpdGF0aW9uX2NvZGVzGAEgAygLMhguYWRtaW4udjEuSW52aXRhdGlvbkNvZGUiUAobQ3JlYXRlSW52aXRhdGlvbkNvZGVSZXF1ZXN0EjEKD2ludml0YXRpb25fY29kZRgBIAEoCzIYLmFkbWluLnYxLkludml0YXRpb25Db2RlIlEKHENyZWF0ZUludml0YXRpb25Db2RlUmVzcG9uc2USMQoPaW52aXRhdGlvbl9jb2RlGAEgASgLMhguYWRtaW4udjEuSW52aXRhdGlvbkNvZGUiUAobVXBkYXRlSW52aXRhdGlvbkNvZGVSZXF1ZXN0EjEKD2ludml0YXRpb25fY29kZRgBIAEoCzIYLmFkbWluLnYxLkludml0YXRpb25Db2RlIlEKHFVwZGF0ZUludml0YXRpb25Db2RlUmVzcG9uc2USMQoPaW52aXRhdGlvbl9jb2RlGAEgASgLMhguYWRtaW4udjEuSW52aXRhdGlvbkNvZGUiKwobRGVsZXRlSW52aXRhdGlvbkNvZGVSZXF1ZXN0EgwKBGNvZGUYASABKAkiHgocRGVsZXRlSW52aXRhdGlvbkNvZGVSZXNwb25zZTKsAwoRSW52aXRhdGlvblNlcnZpY2USYgoTTGlzdEludml0YXRpb25Db2RlcxIkLmFkbWluLnYxLkxpc3RJbnZpdGF0aW9uQ29kZXNSZXF1ZXN0GiUuYWRtaW4udjEuTGlzdEludml0YXRpb25Db2Rlc1Jlc3BvbnNlEmUKFENyZWF0ZUludml0YXRpb25Db2RlEiUuYWRtaW4udjEuQ3JlYXRlSW52aXRhdGlvbkNvZGVSZXF1ZXN0GiYuYWRtaW4udjEuQ3JlYXRlSW52aXRhdGlvbkNvZGVSZXNwb25zZRJlChRVcGRhdGVJbnZpdGF0aW9uQ29kZRIlLmFkbWluLnYxLlVwZGF0ZUludml0YXRpb25Db2RlUmVxdWVzdBomLmFkbWluLnYxLlVwZGF0ZUludml0YXRpb25Db2RlUmVzcG9uc2USZQoURGVsZXRlSW52aXRhdGlvbkNvZGUSJS5hZG1pbi52MS5EZWxldGVJbnZpdGF0aW9uQ29kZVJlcXVlc3QaJi5hZG1pbi52MS5EZWxldGVJbnZpdGF0aW9uQ29kZVJlc3BvbnNlQqMBCgxjb20uYWRtaW4udjFCD0ludml0YXRpb25Qcm90b1ABWkFnaXRodWIuY29tL2ljdHNjL2ljdHNjLXJlZ2FsaWEvYmFja2VuZC9wa2cvcHJvdG8vYWRtaW4vdjE7YWRtaW52MaICA0FYWKoCCEFkbWluLlYxygIIQWRtaW5cVjHiAhRBZG1pblxWMVxHUEJNZXRhZGF0YeoCCUFkbWluOjpWMWIGcHJvdG8z", [file_google_protobuf_timestamp]);

/**
 * @generated from message admin.v1.InvitationCode
 */
export type InvitationCode = Message<"admin.v1.InvitationCode"> & {
  /**
   * @generated from field: string code = 1;
   */
  code: string;

  /**
   * @generated from field: int64 team_code = 2;
   */
  teamCode: bigint;

  /**
   * @generated from field: uint64 total_uses = 3;
   */
  totalUses: bigint;

  /**
   * @generated from field: uint64 remaining_uses = 4;
   */
  remainingUses: bigint;

  /**
   * @generated from field: google.protobuf.Timestamp created_at = 5;
   */
  createdAt?: Timestamp;

  /**
   * @generated from field: google.protobuf.Timestamp expires_at = 6;
   */
  expiresAt?: Timestamp;
};

/**
 * Describes the message admin.v1.InvitationCode.
 * Use `create(InvitationCodeSchema)` to create a new message.
 */
export const InvitationCodeSchema: GenMessage<InvitationCode> = /*@__PURE__*/
  messageDesc(file_admin_v1_invitation, 0);

/**
 * @generated from message admin.v1.ListInvitationCodesRequest
 */
export type ListInvitationCodesRequest = Message<"admin.v1.ListInvitationCodesRequest"> & {
  /**
   * @generated from field: int64 team_code = 1;
   */
  teamCode: bigint;

  /**
   * @generated from field: bool include_expired = 2;
   */
  includeExpired: boolean;
};

/**
 * Describes the message admin.v1.ListInvitationCodesRequest.
 * Use `create(ListInvitationCodesRequestSchema)` to create a new message.
 */
export const ListInvitationCodesRequestSchema: GenMessage<ListInvitationCodesRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_invitation, 1);

/**
 * @generated from message admin.v1.ListInvitationCodesResponse
 */
export type ListInvitationCodesResponse = Message<"admin.v1.ListInvitationCodesResponse"> & {
  /**
   * @generated from field: repeated admin.v1.InvitationCode invitation_codes = 1;
   */
  invitationCodes: InvitationCode[];
};

/**
 * Describes the message admin.v1.ListInvitationCodesResponse.
 * Use `create(ListInvitationCodesResponseSchema)` to create a new message.
 */
export const ListInvitationCodesResponseSchema: GenMessage<ListInvitationCodesResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_invitation, 2);

/**
 * @generated from message admin.v1.CreateInvitationCodeRequest
 */
export type CreateInvitationCodeRequest = Message<"admin.v1.CreateInvitationCodeRequest"> & {
  /**
   * @generated from field: admin.v1.InvitationCode invitation_code = 1;
   */
  invitationCode?: InvitationCode;
};

/**
 * Describes the message admin.v1.CreateInvitationCodeRequest.
 * Use `create(CreateInvitationCodeRequestSchema)` to create a new message.
 */
export const CreateInvitationCodeRequestSchema: GenMessage<CreateInvitationCodeRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_invitation, 3);

/**
 * @generated from message admin.v1.CreateInvitationCodeResponse
 */
export type CreateInvitationCodeResponse = Message<"admin.v1.CreateInvitationCodeResponse"> & {
  /**
   * @generated from field: admin.v1.InvitationCode invitation_code = 1;
   */
  invitationCode?: InvitationCode;
};

/**
 * Describes the message admin.v1.CreateInvitationCodeResponse.
 * Use `create(CreateInvitationCodeResponseSchema)` to create a new message.
 */
export const CreateInvitationCodeResponseSchema: GenMessage<CreateInvitationCodeResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_invitation, 4);

/**
 * @generated from message admin.v1.UpdateInvitationCodeRequest
 */
export type UpdateInvitationCodeRequest = Message<"admin.v1.UpdateInvitationCodeRequest"> & {
  /**
   * @generated from field: admin.v1.InvitationCode invitation_code = 1;
   */
  invitationCode?: InvitationCode;
};

/**
 * Describes the message admin.v1.UpdateInvitationCodeRequest.
 * Use `create(UpdateInvitationCodeRequestSchema)` to create a new message.
 */
export const UpdateInvitationCodeRequestSchema: GenMessage<UpdateInvitationCodeRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_invitation, 5);

/**
 * @generated from message admin.v1.UpdateInvitationCodeResponse
 */
export type UpdateInvitationCodeResponse = Message<"admin.v1.UpdateInvitationCodeResponse"> & {
  /**
   * @generated from field: admin.v1.InvitationCode invitation_code = 1;
   */
  invitationCode?: InvitationCode;
};

/**
 * Describes the message admin.v1.UpdateInvitationCodeResponse.
 * Use `create(UpdateInvitationCodeResponseSchema)` to create a new message.
 */
export const UpdateInvitationCodeResponseSchema: GenMessage<UpdateInvitationCodeResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_invitation, 6);

/**
 * @generated from message admin.v1.DeleteInvitationCodeRequest
 */
export type DeleteInvitationCodeRequest = Message<"admin.v1.DeleteInvitationCodeRequest"> & {
  /**
   * @generated from field: string code = 1;
   */
  code: string;
};

/**
 * Describes the message admin.v1.DeleteInvitationCodeRequest.
 * Use `create(DeleteInvitationCodeRequestSchema)` to create a new message.
 */
export const DeleteInvitationCodeRequestSchema: GenMessage<DeleteInvitationCodeRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_invitation, 7);

/**
 * @generated from message admin.v1.DeleteInvitationCodeResponse
 */
export type DeleteInvitationCodeResponse = Message<"admin.v1.DeleteInvitationCodeResponse"> & {
};

/**
 * Describes the message admin.v1.DeleteInvitationCodeResponse.
 * Use `create(DeleteInvitationCodeResponseSchema)` to create a new message.
 */
export const DeleteInvitationCodeResponseSchema: GenMessage<DeleteInvitationCodeResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_invitation, 8);

/**
 * @generated from service admin.v1.InvitationService
 */
export const InvitationService: GenService<{
  /**
   * @generated from rpc admin.v1.InvitationService.ListInvitationCodes
   */
  listInvitationCodes: {
    methodKind: "unary";
    input: typeof ListInvitationCodesRequestSchema;
    output: typeof ListInvitationCodesResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.InvitationService.CreateInvitationCode
   */
  createInvitationCode: {
    methodKind: "unary";
    input: typeof CreateInvitationCodeRequestSchema;
    output: typeof CreateInvitationCodeResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.InvitationService.UpdateInvitationCode
   */
  updateInvitationCode: {
    methodKind: "unary";
    input: typeof UpdateInvitationCodeRequestSchema;
    output: typeof UpdateInvitationCodeResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.InvitationService.DeleteInvitationCode
   */
  deleteInvitationCode: {
    methodKind: "unary";
    input: typeof DeleteInvitationCodeRequestSchema;
    output: typeof DeleteInvitationCodeResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_invitation, 0);

