// @generated by protoc-gen-es v2.0.0 with parameter "target=ts"
// @generated from file contestant/v1/rule.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/rule.proto.
 */
export const file_contestant_v1_rule: GenFile = /*@__PURE__*/
  fileDesc("Chhjb250ZXN0YW50L3YxL3J1bGUucHJvdG8SDWNvbnRlc3RhbnQudjEiZQoEUnVsZRIYCgRydWxlGAEgASgJQgq6SAdyBRABGOgHEh4KCnNob3J0X3J1bGUYAiABKAlCCrpIB3IFEAEYyAESIwoPcmVjcmVhdGlvbl9ydWxlGAMgASgJQgq6SAdyBRABGOgHIhAKDkdldFJ1bGVSZXF1ZXN0IjwKD0dldFJ1bGVSZXNwb25zZRIpCgRydWxlGAEgASgLMhMuY29udGVzdGFudC52MS5SdWxlQga6SAPIAQEyVwoLUnVsZVNlcnZpY2USSAoHR2V0UnVsZRIdLmNvbnRlc3RhbnQudjEuR2V0UnVsZVJlcXVlc3QaHi5jb250ZXN0YW50LnYxLkdldFJ1bGVSZXNwb25zZUK5AQoRY29tLmNvbnRlc3RhbnQudjFCCVJ1bGVQcm90b1ABWkRnaXRodWIuY29tL2ljdHNjL2ljdHNjLW91dGxhbmRzL2JhY2tlbmQvaW50ZXJuYWwvcHJvdG8vY29udGVzdGFudC92MaICA0NYWKoCDUNvbnRlc3RhbnQuVjHKAg1Db250ZXN0YW50XFYx4gIZQ29udGVzdGFudFxWMVxHUEJNZXRhZGF0YeoCDkNvbnRlc3RhbnQ6OlYxYgZwcm90bzM", [file_buf_validate_validate]);

/**
 * @generated from message contestant.v1.Rule
 */
export type Rule = Message<"contestant.v1.Rule"> & {
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
 * Describes the message contestant.v1.Rule.
 * Use `create(RuleSchema)` to create a new message.
 */
export const RuleSchema: GenMessage<Rule> = /*@__PURE__*/
  messageDesc(file_contestant_v1_rule, 0);

/**
 * @generated from message contestant.v1.GetRuleRequest
 */
export type GetRuleRequest = Message<"contestant.v1.GetRuleRequest"> & {
};

/**
 * Describes the message contestant.v1.GetRuleRequest.
 * Use `create(GetRuleRequestSchema)` to create a new message.
 */
export const GetRuleRequestSchema: GenMessage<GetRuleRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_rule, 1);

/**
 * @generated from message contestant.v1.GetRuleResponse
 */
export type GetRuleResponse = Message<"contestant.v1.GetRuleResponse"> & {
  /**
   * @generated from field: contestant.v1.Rule rule = 1;
   */
  rule?: Rule;
};

/**
 * Describes the message contestant.v1.GetRuleResponse.
 * Use `create(GetRuleResponseSchema)` to create a new message.
 */
export const GetRuleResponseSchema: GenMessage<GetRuleResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_rule, 2);

/**
 * @generated from service contestant.v1.RuleService
 */
export const RuleService: GenService<{
  /**
   * @generated from rpc contestant.v1.RuleService.GetRule
   */
  getRule: {
    methodKind: "unary";
    input: typeof GetRuleRequestSchema;
    output: typeof GetRuleResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_contestant_v1_rule, 0);

