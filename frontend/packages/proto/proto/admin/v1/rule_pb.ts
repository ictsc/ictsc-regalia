// @generated by protoc-gen-es v2.2.1 with parameter "target=ts"
// @generated from file admin/v1/rule.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/rule.proto.
 */
export const file_admin_v1_rule: GenFile = /*@__PURE__*/
  fileDesc("ChNhZG1pbi92MS9ydWxlLnByb3RvEghhZG1pbi52MSJlCgRSdWxlEhgKBHJ1bGUYASABKAlCCrpIB3IFEAEY6AcSHgoKc2hvcnRfcnVsZRgCIAEoCUIKukgHcgUQARjIARIjCg9yZWNyZWF0aW9uX3J1bGUYAyABKAlCCrpIB3IFEAEY6AciEAoOR2V0UnVsZVJlcXVlc3QiNwoPR2V0UnVsZVJlc3BvbnNlEiQKBHJ1bGUYASABKAsyDi5hZG1pbi52MS5SdWxlQga6SAPIAQEirAEKEFBhdGNoUnVsZVJlcXVlc3QSHQoEcnVsZRgBIAEoCUIKukgHcgUQARjoB0gAiAEBEiMKCnNob3J0X3J1bGUYAiABKAlCCrpIB3IFEAEYyAFIAYgBARIoCg9yZWNyZWF0aW9uX3J1bGUYAyABKAlCCrpIB3IFEAEY6AdIAogBAUIHCgVfcnVsZUINCgtfc2hvcnRfcnVsZUISChBfcmVjcmVhdGlvbl9ydWxlIjkKEVBhdGNoUnVsZVJlc3BvbnNlEiQKBHJ1bGUYASABKAsyDi5hZG1pbi52MS5SdWxlQga6SAPIAQEykwEKC1J1bGVTZXJ2aWNlEj4KB0dldFJ1bGUSGC5hZG1pbi52MS5HZXRSdWxlUmVxdWVzdBoZLmFkbWluLnYxLkdldFJ1bGVSZXNwb25zZRJECglQYXRjaFJ1bGUSGi5hZG1pbi52MS5QYXRjaFJ1bGVSZXF1ZXN0GhsuYWRtaW4udjEuUGF0Y2hSdWxlUmVzcG9uc2VCmwEKDGNvbS5hZG1pbi52MUIJUnVsZVByb3RvUAFaP2dpdGh1Yi5jb20vaWN0c2MvaWN0c2Mtb3V0bGFuZHMvYmFja2VuZC9pbnRlcm5hbC9wcm90by9hZG1pbi92MaICA0FYWKoCCEFkbWluLlYxygIIQWRtaW5cVjHiAhRBZG1pblxWMVxHUEJNZXRhZGF0YeoCCUFkbWluOjpWMWIGcHJvdG8z", [file_buf_validate_validate]);

/**
 * @generated from message admin.v1.Rule
 */
export type Rule = Message<"admin.v1.Rule"> & {
  /**
   * @generated from field: string rule = 1;
   */
  rule: string;

  /**
   * @generated from field: string short_rule = 2;
   */
  shortRule: string;

  /**
   * @generated from field: string recreation_rule = 3;
   */
  recreationRule: string;
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
 * @generated from message admin.v1.PatchRuleRequest
 */
export type PatchRuleRequest = Message<"admin.v1.PatchRuleRequest"> & {
  /**
   * @generated from field: optional string rule = 1;
   */
  rule?: string;

  /**
   * @generated from field: optional string short_rule = 2;
   */
  shortRule?: string;

  /**
   * @generated from field: optional string recreation_rule = 3;
   */
  recreationRule?: string;
};

/**
 * Describes the message admin.v1.PatchRuleRequest.
 * Use `create(PatchRuleRequestSchema)` to create a new message.
 */
export const PatchRuleRequestSchema: GenMessage<PatchRuleRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_rule, 3);

/**
 * @generated from message admin.v1.PatchRuleResponse
 */
export type PatchRuleResponse = Message<"admin.v1.PatchRuleResponse"> & {
  /**
   * @generated from field: admin.v1.Rule rule = 1;
   */
  rule?: Rule;
};

/**
 * Describes the message admin.v1.PatchRuleResponse.
 * Use `create(PatchRuleResponseSchema)` to create a new message.
 */
export const PatchRuleResponseSchema: GenMessage<PatchRuleResponse> = /*@__PURE__*/
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
   * @generated from rpc admin.v1.RuleService.PatchRule
   */
  patchRule: {
    methodKind: "unary";
    input: typeof PatchRuleRequestSchema;
    output: typeof PatchRuleResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_rule, 0);

