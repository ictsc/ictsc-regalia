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
  fileDesc("ChlhZG1pbi92MS9pbnZpdGF0aW9uLnByb3RvEghhZG1pbi52MSKRAQoOSW52aXRhdGlvbkNvZGUSDAoEY29kZRgBIAEoCRIRCgl0ZWFtX2NvZGUYAiABKAMSLgoKY3JlYXRlZF9hdBgDIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXASLgoKZXhwaXJlc19hdBgEIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXAiSAoaTGlzdEludml0YXRpb25Db2Rlc1JlcXVlc3QSEQoJdGVhbV9jb2RlGAEgASgDEhcKD2luY2x1ZGVfZXhwaXJlZBgCIAEoCCJRChtMaXN0SW52aXRhdGlvbkNvZGVzUmVzcG9uc2USMgoQaW52aXRhdGlvbl9jb2RlcxgBIAMoCzIYLmFkbWluLnYxLkludml0YXRpb25Db2RlIlAKG0NyZWF0ZUludml0YXRpb25Db2RlUmVxdWVzdBIxCg9pbnZpdGF0aW9uX2NvZGUYASABKAsyGC5hZG1pbi52MS5JbnZpdGF0aW9uQ29kZSJRChxDcmVhdGVJbnZpdGF0aW9uQ29kZVJlc3BvbnNlEjEKD2ludml0YXRpb25fY29kZRgBIAEoCzIYLmFkbWluLnYxLkludml0YXRpb25Db2RlMt4BChFJbnZpdGF0aW9uU2VydmljZRJiChNMaXN0SW52aXRhdGlvbkNvZGVzEiQuYWRtaW4udjEuTGlzdEludml0YXRpb25Db2Rlc1JlcXVlc3QaJS5hZG1pbi52MS5MaXN0SW52aXRhdGlvbkNvZGVzUmVzcG9uc2USZQoUQ3JlYXRlSW52aXRhdGlvbkNvZGUSJS5hZG1pbi52MS5DcmVhdGVJbnZpdGF0aW9uQ29kZVJlcXVlc3QaJi5hZG1pbi52MS5DcmVhdGVJbnZpdGF0aW9uQ29kZVJlc3BvbnNlQqMBCgxjb20uYWRtaW4udjFCD0ludml0YXRpb25Qcm90b1ABWkFnaXRodWIuY29tL2ljdHNjL2ljdHNjLXJlZ2FsaWEvYmFja2VuZC9wa2cvcHJvdG8vYWRtaW4vdjE7YWRtaW52MaICA0FYWKoCCEFkbWluLlYxygIIQWRtaW5cVjHiAhRBZG1pblxWMVxHUEJNZXRhZGF0YeoCCUFkbWluOjpWMWIGcHJvdG8z", [file_google_protobuf_timestamp]);

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
   * @generated from field: google.protobuf.Timestamp created_at = 3;
   */
  createdAt?: Timestamp;

  /**
   * @generated from field: google.protobuf.Timestamp expires_at = 4;
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
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_invitation, 0);

