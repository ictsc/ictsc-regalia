// @generated by protoc-gen-es v2.6.1 with parameter "target=ts"
// @generated from file admin/v1/rule.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv2";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv2";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/rule.proto.
 */
export const file_admin_v1_rule: GenFile = /*@__PURE__*/
  fileDesc("ChNhZG1pbi92MS9ydWxlLnByb3RvEghhZG1pbi52MSIYCgRSdWxlEhAKCG1hcmtkb3duGAIgASgJIhAKDkdldFJ1bGVSZXF1ZXN0Ii8KD0dldFJ1bGVSZXNwb25zZRIcCgRydWxlGAEgASgLMg4uYWRtaW4udjEuUnVsZSIxChFVcGRhdGVSdWxlUmVxdWVzdBIcCgRydWxlGAEgASgLMg4uYWRtaW4udjEuUnVsZSIUChJVcGRhdGVSdWxlUmVzcG9uc2UylgEKC1J1bGVTZXJ2aWNlEj4KB0dldFJ1bGUSGC5hZG1pbi52MS5HZXRSdWxlUmVxdWVzdBoZLmFkbWluLnYxLkdldFJ1bGVSZXNwb25zZRJHCgpVcGRhdGVSdWxlEhsuYWRtaW4udjEuVXBkYXRlUnVsZVJlcXVlc3QaHC5hZG1pbi52MS5VcGRhdGVSdWxlUmVzcG9uc2VCnQEKDGNvbS5hZG1pbi52MUIJUnVsZVByb3RvUAFaQWdpdGh1Yi5jb20vaWN0c2MvaWN0c2MtcmVnYWxpYS9iYWNrZW5kL3BrZy9wcm90by9hZG1pbi92MTthZG1pbnYxogIDQVhYqgIIQWRtaW4uVjHKAghBZG1pblxWMeICFEFkbWluXFYxXEdQQk1ldGFkYXRh6gIJQWRtaW46OlYxYgZwcm90bzM");

/**
 * @generated from message admin.v1.Rule
 */
export type Rule = Message<"admin.v1.Rule"> & {
  /**
   * @generated from field: string markdown = 2;
   */
  markdown: string;
};

/**
 * Describes the message admin.v1.Rule.
 * Use `create(RuleSchema)` to create a new message.
 */
export const RuleSchema: GenMessage<Rule> = /*@__PURE__*/
  messageDesc(file_admin_v1_rule, 0);

/**
 * @generated from message admin.v1.GetRuleRequest
 */
export type GetRuleRequest = Message<"admin.v1.GetRuleRequest"> & {
};

/**
 * Describes the message admin.v1.GetRuleRequest.
 * Use `create(GetRuleRequestSchema)` to create a new message.
 */
export const GetRuleRequestSchema: GenMessage<GetRuleRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_rule, 1);

/**
 * @generated from message admin.v1.GetRuleResponse
 */
export type GetRuleResponse = Message<"admin.v1.GetRuleResponse"> & {
  /**
   * @generated from field: admin.v1.Rule rule = 1;
   */
  rule?: Rule;
};

/**
 * Describes the message admin.v1.GetRuleResponse.
 * Use `create(GetRuleResponseSchema)` to create a new message.
 */
export const GetRuleResponseSchema: GenMessage<GetRuleResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_rule, 2);

/**
 * @generated from message admin.v1.UpdateRuleRequest
 */
export type UpdateRuleRequest = Message<"admin.v1.UpdateRuleRequest"> & {
  /**
   * @generated from field: admin.v1.Rule rule = 1;
   */
  rule?: Rule;
};

/**
 * Describes the message admin.v1.UpdateRuleRequest.
 * Use `create(UpdateRuleRequestSchema)` to create a new message.
 */
export const UpdateRuleRequestSchema: GenMessage<UpdateRuleRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_rule, 3);

/**
 * @generated from message admin.v1.UpdateRuleResponse
 */
export type UpdateRuleResponse = Message<"admin.v1.UpdateRuleResponse"> & {
};

/**
 * Describes the message admin.v1.UpdateRuleResponse.
 * Use `create(UpdateRuleResponseSchema)` to create a new message.
 */
export const UpdateRuleResponseSchema: GenMessage<UpdateRuleResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_rule, 4);

/**
 * @generated from service admin.v1.RuleService
 */
export const RuleService: GenService<{
  /**
   * @generated from rpc admin.v1.RuleService.GetRule
   */
  getRule: {
    methodKind: "unary";
    input: typeof GetRuleRequestSchema;
    output: typeof GetRuleResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.RuleService.UpdateRule
   */
  updateRule: {
    methodKind: "unary";
    input: typeof UpdateRuleRequestSchema;
    output: typeof UpdateRuleResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_rule, 0);

