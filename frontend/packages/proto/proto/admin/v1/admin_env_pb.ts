// @generated by protoc-gen-es v2.2.1 with parameter "target=ts"
// @generated from file admin/v1/admin_env.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Bastion } from "./team_pb";
import { file_admin_v1_team } from "./team_pb";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/admin_env.proto.
 */
export const file_admin_v1_admin_env: GenFile = /*@__PURE__*/
  fileDesc("ChhhZG1pbi92MS9hZG1pbl9lbnYucHJvdG8SCGFkbWluLnYxIh8KHUdldEFkbWluQ29ubmVjdGlvbkluZm9SZXF1ZXN0IkwKHkdldEFkbWluQ29ubmVjdGlvbkluZm9SZXNwb25zZRIqCgdiYXN0aW9uGAEgASgLMhEuYWRtaW4udjEuQmFzdGlvbkIGukgDyAEBIkMKHVB1dEFkbWluQ29ubmVjdGlvbkluZm9SZXF1ZXN0EiIKB2Jhc3Rpb24YASABKAsyES5hZG1pbi52MS5CYXN0aW9uIiAKHlB1dEFkbWluQ29ubmVjdGlvbkluZm9SZXNwb25zZTLrAQoPQWRtaW5FbnZTZXJ2aWNlEmsKFkdldEFkbWluQ29ubmVjdGlvbkluZm8SJy5hZG1pbi52MS5HZXRBZG1pbkNvbm5lY3Rpb25JbmZvUmVxdWVzdBooLmFkbWluLnYxLkdldEFkbWluQ29ubmVjdGlvbkluZm9SZXNwb25zZRJrChZQdXRBZG1pbkNvbm5lY3Rpb25JbmZvEicuYWRtaW4udjEuUHV0QWRtaW5Db25uZWN0aW9uSW5mb1JlcXVlc3QaKC5hZG1pbi52MS5QdXRBZG1pbkNvbm5lY3Rpb25JbmZvUmVzcG9uc2VCnwEKDGNvbS5hZG1pbi52MUINQWRtaW5FbnZQcm90b1ABWj9naXRodWIuY29tL2ljdHNjL2ljdHNjLW91dGxhbmRzL2JhY2tlbmQvaW50ZXJuYWwvcHJvdG8vYWRtaW4vdjGiAgNBWFiqAghBZG1pbi5WMcoCCEFkbWluXFYx4gIUQWRtaW5cVjFcR1BCTWV0YWRhdGHqAglBZG1pbjo6VjFiBnByb3RvMw", [file_admin_v1_team, file_buf_validate_validate]);

/**
 * @generated from message admin.v1.GetAdminConnectionInfoRequest
 */
export type GetAdminConnectionInfoRequest = Message<"admin.v1.GetAdminConnectionInfoRequest"> & {
};

/**
 * Describes the message admin.v1.GetAdminConnectionInfoRequest.
 * Use `create(GetAdminConnectionInfoRequestSchema)` to create a new message.
 */
export const GetAdminConnectionInfoRequestSchema: GenMessage<GetAdminConnectionInfoRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_admin_env, 0);

/**
 * @generated from message admin.v1.GetAdminConnectionInfoResponse
 */
export type GetAdminConnectionInfoResponse = Message<"admin.v1.GetAdminConnectionInfoResponse"> & {
  /**
   * @generated from field: admin.v1.Bastion bastion = 1;
   */
  bastion?: Bastion;
};

/**
 * Describes the message admin.v1.GetAdminConnectionInfoResponse.
 * Use `create(GetAdminConnectionInfoResponseSchema)` to create a new message.
 */
export const GetAdminConnectionInfoResponseSchema: GenMessage<GetAdminConnectionInfoResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_admin_env, 1);

/**
 * @generated from message admin.v1.PutAdminConnectionInfoRequest
 */
export type PutAdminConnectionInfoRequest = Message<"admin.v1.PutAdminConnectionInfoRequest"> & {
  /**
   * @generated from field: admin.v1.Bastion bastion = 1;
   */
  bastion?: Bastion;
};

/**
 * Describes the message admin.v1.PutAdminConnectionInfoRequest.
 * Use `create(PutAdminConnectionInfoRequestSchema)` to create a new message.
 */
export const PutAdminConnectionInfoRequestSchema: GenMessage<PutAdminConnectionInfoRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_admin_env, 2);

/**
 * @generated from message admin.v1.PutAdminConnectionInfoResponse
 */
export type PutAdminConnectionInfoResponse = Message<"admin.v1.PutAdminConnectionInfoResponse"> & {
};

/**
 * Describes the message admin.v1.PutAdminConnectionInfoResponse.
 * Use `create(PutAdminConnectionInfoResponseSchema)` to create a new message.
 */
export const PutAdminConnectionInfoResponseSchema: GenMessage<PutAdminConnectionInfoResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_admin_env, 3);

/**
 * @generated from service admin.v1.AdminEnvService
 */
export const AdminEnvService: GenService<{
  /**
   * @generated from rpc admin.v1.AdminEnvService.GetAdminConnectionInfo
   */
  getAdminConnectionInfo: {
    methodKind: "unary";
    input: typeof GetAdminConnectionInfoRequestSchema;
    output: typeof GetAdminConnectionInfoResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.AdminEnvService.PutAdminConnectionInfo
   */
  putAdminConnectionInfo: {
    methodKind: "unary";
    input: typeof PutAdminConnectionInfoRequestSchema;
    output: typeof PutAdminConnectionInfoResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_admin_env, 0);

