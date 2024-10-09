// @generated by protoc-gen-es v2.2.0 with parameter "target=ts"
// @generated from file admin/v1/auth.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/auth.proto.
 */
export const file_admin_v1_auth: GenFile = /*@__PURE__*/
  fileDesc("ChNhZG1pbi92MS9hdXRoLnByb3RvEghhZG1pbi52MSIUChJHZXRDYWxsYmFja1JlcXVlc3QiNQoTR2V0Q2FsbGJhY2tSZXNwb25zZRIeCgxyZWRpcmVjdF91cmkYASABKAlCCLpIBXIDiAEBIigKD1Bvc3RDb2RlUmVxdWVzdBIVCgRjb2RlGAEgASgJQge6SARyAhABIhIKEFBvc3RDb2RlUmVzcG9uc2UyoAEKC0F1dGhTZXJ2aWNlEkwKC0dldENhbGxiYWNrEhwuYWRtaW4udjEuR2V0Q2FsbGJhY2tSZXF1ZXN0Gh0uYWRtaW4udjEuR2V0Q2FsbGJhY2tSZXNwb25zZSIAEkMKCFBvc3RDb2RlEhkuYWRtaW4udjEuUG9zdENvZGVSZXF1ZXN0GhouYWRtaW4udjEuUG9zdENvZGVSZXNwb25zZSIAQpsBCgxjb20uYWRtaW4udjFCCUF1dGhQcm90b1ABWj9naXRodWIuY29tL2ljdHNjL2ljdHNjLW91dGxhbmRzL2JhY2tlbmQvaW50ZXJuYWwvcHJvdG8vYWRtaW4vdjGiAgNBWFiqAghBZG1pbi5WMcoCCEFkbWluXFYx4gIUQWRtaW5cVjFcR1BCTWV0YWRhdGHqAglBZG1pbjo6VjFiBnByb3RvMw", [file_buf_validate_validate]);

/**
 * @generated from message admin.v1.GetCallbackRequest
 */
export type GetCallbackRequest = Message<"admin.v1.GetCallbackRequest"> & {
};

/**
 * Describes the message admin.v1.GetCallbackRequest.
 * Use `create(GetCallbackRequestSchema)` to create a new message.
 */
export const GetCallbackRequestSchema: GenMessage<GetCallbackRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_auth, 0);

/**
 * @generated from message admin.v1.GetCallbackResponse
 */
export type GetCallbackResponse = Message<"admin.v1.GetCallbackResponse"> & {
  /**
   * @generated from field: string redirect_uri = 1;
   */
  redirectUri: string;
};

/**
 * Describes the message admin.v1.GetCallbackResponse.
 * Use `create(GetCallbackResponseSchema)` to create a new message.
 */
export const GetCallbackResponseSchema: GenMessage<GetCallbackResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_auth, 1);

/**
 * @generated from message admin.v1.PostCodeRequest
 */
export type PostCodeRequest = Message<"admin.v1.PostCodeRequest"> & {
  /**
   * @generated from field: string code = 1;
   */
  code: string;
};

/**
 * Describes the message admin.v1.PostCodeRequest.
 * Use `create(PostCodeRequestSchema)` to create a new message.
 */
export const PostCodeRequestSchema: GenMessage<PostCodeRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_auth, 2);

/**
 * @generated from message admin.v1.PostCodeResponse
 */
export type PostCodeResponse = Message<"admin.v1.PostCodeResponse"> & {
};

/**
 * Describes the message admin.v1.PostCodeResponse.
 * Use `create(PostCodeResponseSchema)` to create a new message.
 */
export const PostCodeResponseSchema: GenMessage<PostCodeResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_auth, 3);

/**
 * @generated from service admin.v1.AuthService
 */
export const AuthService: GenService<{
  /**
   * @generated from rpc admin.v1.AuthService.GetCallback
   */
  getCallback: {
    methodKind: "unary";
    input: typeof GetCallbackRequestSchema;
    output: typeof GetCallbackResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.AuthService.PostCode
   */
  postCode: {
    methodKind: "unary";
    input: typeof PostCodeRequestSchema;
    output: typeof PostCodeResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_auth, 0);

