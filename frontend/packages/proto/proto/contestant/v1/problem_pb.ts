// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file contestant/v1/problem.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/problem.proto.
 */
export const file_contestant_v1_problem: GenFile = /*@__PURE__*/
  fileDesc("Chtjb250ZXN0YW50L3YxL3Byb2JsZW0ucHJvdG8SDWNvbnRlc3RhbnQudjEixgEKB1Byb2JsZW0SDAoEY29kZRgBIAEoCRINCgV0aXRsZRgCIAEoCRIRCgltYXhfc2NvcmUYAyABKAMSKAoFc2NvcmUYBCABKAsyFC5jb250ZXN0YW50LnYxLlNjb3JlSACIAQESLQoKZGVwbG95bWVudBgFIAEoCzIZLmNvbnRlc3RhbnQudjEuRGVwbG95bWVudBIoCgRib2R5GAYgASgLMhouY29udGVzdGFudC52MS5Qcm9ibGVtQm9keUIICgZfc2NvcmUiPQoFU2NvcmUSFAoMbWFya2VkX3Njb3JlGAEgASgDEg8KB3BlbmFsdHkYAiABKAMSDQoFc2NvcmUYAyABKAMiUwoKRGVwbG95bWVudBIvCgZzdGF0dXMYASABKA4yHy5jb250ZXN0YW50LnYxLkRlcGxveW1lbnRTdGF0dXMSFAoMcmVkZXBsb3lhYmxlGAIgASgIInkKC1Byb2JsZW1Cb2R5EigKBHR5cGUYASABKA4yGi5jb250ZXN0YW50LnYxLlByb2JsZW1UeXBlEjgKC2Rlc2NyaXB0aXZlGAIgASgLMiEuY29udGVzdGFudC52MS5EZXNjcmlwdGl2ZVByb2JsZW1IAEIGCgRib2R5Im0KCkNvbm5lY3Rpb24SEQoJaG9zdF9uYW1lGAEgASgJEgwKBGhvc3QYAiABKAkSEQoEdXNlchgDIAEoCUgAiAEBEhUKCHBhc3N3b3JkGAQgASgJSAGIAQFCBwoFX3VzZXJCCwoJX3Bhc3N3b3JkIlEKEkRlc2NyaXB0aXZlUHJvYmxlbRIMCgRib2R5GAEgASgJEi0KCmNvbm5lY3Rpb24YAiADKAsyGS5jb250ZXN0YW50LnYxLkNvbm5lY3Rpb24iFQoTTGlzdFByb2JsZW1zUmVxdWVzdCJAChRMaXN0UHJvYmxlbXNSZXNwb25zZRIoCghwcm9ibGVtcxgBIAMoCzIWLmNvbnRlc3RhbnQudjEuUHJvYmxlbSIhChFHZXRQcm9ibGVtUmVxdWVzdBIMCgRjb2RlGAEgASgJIj0KEkdldFByb2JsZW1SZXNwb25zZRInCgdwcm9ibGVtGAEgASgLMhYuY29udGVzdGFudC52MS5Qcm9ibGVtIh0KDURlcGxveVJlcXVlc3QSDAoEY29kZRgBIAEoCSIQCg5EZXBsb3lSZXNwb25zZSqUAQoQRGVwbG95bWVudFN0YXR1cxIhCh1ERVBMT1lNRU5UX1NUQVRVU19VTlNQRUNJRklFRBAAEh4KGkRFUExPWU1FTlRfU1RBVFVTX0RFUExPWUVEEAESHwobREVQTE9ZTUVOVF9TVEFUVVNfREVQTE9ZSU5HEAISHAoYREVQTE9ZTUVOVF9TVEFUVVNfRkFJTEVEEAMqSQoLUHJvYmxlbVR5cGUSHAoYUFJPQkxFTV9UWVBFX1VOU1BFQ0lGSUVEEAASHAoYUFJPQkxFTV9UWVBFX0RFU0NSSVBUSVZFEAEygwIKDlByb2JsZW1TZXJ2aWNlElcKDExpc3RQcm9ibGVtcxIiLmNvbnRlc3RhbnQudjEuTGlzdFByb2JsZW1zUmVxdWVzdBojLmNvbnRlc3RhbnQudjEuTGlzdFByb2JsZW1zUmVzcG9uc2USUQoKR2V0UHJvYmxlbRIgLmNvbnRlc3RhbnQudjEuR2V0UHJvYmxlbVJlcXVlc3QaIS5jb250ZXN0YW50LnYxLkdldFByb2JsZW1SZXNwb25zZRJFCgZEZXBsb3kSHC5jb250ZXN0YW50LnYxLkRlcGxveVJlcXVlc3QaHS5jb250ZXN0YW50LnYxLkRlcGxveVJlc3BvbnNlQsMBChFjb20uY29udGVzdGFudC52MUIMUHJvYmxlbVByb3RvUAFaS2dpdGh1Yi5jb20vaWN0c2MvaWN0c2MtcmVnYWxpYS9iYWNrZW5kL3BrZy9wcm90by9jb250ZXN0YW50L3YxO2NvbnRlc3RhbnR2MaICA0NYWKoCDUNvbnRlc3RhbnQuVjHKAg1Db250ZXN0YW50XFYx4gIZQ29udGVzdGFudFxWMVxHUEJNZXRhZGF0YeoCDkNvbnRlc3RhbnQ6OlYxYgZwcm90bzM");

/**
 * @generated from message contestant.v1.Problem
 */
export type Problem = Message<"contestant.v1.Problem"> & {
  /**
   * 問題コード
   *
   * @generated from field: string code = 1;
   */
  code: string;

  /**
   * タイトル
   *
   * @generated from field: string title = 2;
   */
  title: string;

  /**
   * 最大得点
   *
   * @generated from field: int64 max_score = 3;
   */
  maxScore: bigint;

  /**
   * @generated from field: optional contestant.v1.Score score = 4;
   */
  score?: Score;

  /**
   * @generated from field: contestant.v1.Deployment deployment = 5;
   */
  deployment?: Deployment;

  /**
   * @generated from field: contestant.v1.ProblemBody body = 6;
   */
  body?: ProblemBody;
};

/**
 * Describes the message contestant.v1.Problem.
 * Use `create(ProblemSchema)` to create a new message.
 */
export const ProblemSchema: GenMessage<Problem> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 0);

/**
 * @generated from message contestant.v1.Score
 */
export type Score = Message<"contestant.v1.Score"> & {
  /**
   * 採点による得点
   *
   * @generated from field: int64 marked_score = 1;
   */
  markedScore: bigint;

  /**
   * ペナルティによる減点
   *
   * @generated from field: int64 penalty = 2;
   */
  penalty: bigint;

  /**
   * 最終的な得点
   *
   * @generated from field: int64 score = 3;
   */
  score: bigint;
};

/**
 * Describes the message contestant.v1.Score.
 * Use `create(ScoreSchema)` to create a new message.
 */
export const ScoreSchema: GenMessage<Score> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 1);

/**
 * @generated from message contestant.v1.Deployment
 */
export type Deployment = Message<"contestant.v1.Deployment"> & {
  /**
   * @generated from field: contestant.v1.DeploymentStatus status = 1;
   */
  status: DeploymentStatus;

  /**
   * @generated from field: bool redeployable = 2;
   */
  redeployable: boolean;
};

/**
 * Describes the message contestant.v1.Deployment.
 * Use `create(DeploymentSchema)` to create a new message.
 */
export const DeploymentSchema: GenMessage<Deployment> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 2);

/**
 * @generated from message contestant.v1.ProblemBody
 */
export type ProblemBody = Message<"contestant.v1.ProblemBody"> & {
  /**
   * @generated from field: contestant.v1.ProblemType type = 1;
   */
  type: ProblemType;

  /**
   * @generated from oneof contestant.v1.ProblemBody.body
   */
  body: {
    /**
     * @generated from field: contestant.v1.DescriptiveProblem descriptive = 2;
     */
    value: DescriptiveProblem;
    case: "descriptive";
  } | { case: undefined; value?: undefined };
};

/**
 * Describes the message contestant.v1.ProblemBody.
 * Use `create(ProblemBodySchema)` to create a new message.
 */
export const ProblemBodySchema: GenMessage<ProblemBody> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 3);

/**
 * @generated from message contestant.v1.Connection
 */
export type Connection = Message<"contestant.v1.Connection"> & {
  /**
   * ホスト名
   *
   * @generated from field: string host_name = 1;
   */
  hostName: string;

  /**
   * ホスト(IP アドレス or ドメイン)
   *
   * @generated from field: string host = 2;
   */
  host: string;

  /**
   * ユーザ
   *
   * @generated from field: optional string user = 3;
   */
  user?: string;

  /**
   * パスワード
   *
   * @generated from field: optional string password = 4;
   */
  password?: string;
};

/**
 * Describes the message contestant.v1.Connection.
 * Use `create(ConnectionSchema)` to create a new message.
 */
export const ConnectionSchema: GenMessage<Connection> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 4);

/**
 * @generated from message contestant.v1.DescriptiveProblem
 */
export type DescriptiveProblem = Message<"contestant.v1.DescriptiveProblem"> & {
  /**
   * 問題文
   *
   * @generated from field: string body = 1;
   */
  body: string;

  /**
   * 接続情報
   *
   * @generated from field: repeated contestant.v1.Connection connection = 2;
   */
  connection: Connection[];
};

/**
 * Describes the message contestant.v1.DescriptiveProblem.
 * Use `create(DescriptiveProblemSchema)` to create a new message.
 */
export const DescriptiveProblemSchema: GenMessage<DescriptiveProblem> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 5);

/**
 * @generated from message contestant.v1.ListProblemsRequest
 */
export type ListProblemsRequest = Message<"contestant.v1.ListProblemsRequest"> & {
};

/**
 * Describes the message contestant.v1.ListProblemsRequest.
 * Use `create(ListProblemsRequestSchema)` to create a new message.
 */
export const ListProblemsRequestSchema: GenMessage<ListProblemsRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 6);

/**
 * @generated from message contestant.v1.ListProblemsResponse
 */
export type ListProblemsResponse = Message<"contestant.v1.ListProblemsResponse"> & {
  /**
   * @generated from field: repeated contestant.v1.Problem problems = 1;
   */
  problems: Problem[];
};

/**
 * Describes the message contestant.v1.ListProblemsResponse.
 * Use `create(ListProblemsResponseSchema)` to create a new message.
 */
export const ListProblemsResponseSchema: GenMessage<ListProblemsResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 7);

/**
 * @generated from message contestant.v1.GetProblemRequest
 */
export type GetProblemRequest = Message<"contestant.v1.GetProblemRequest"> & {
  /**
   * @generated from field: string code = 1;
   */
  code: string;
};

/**
 * Describes the message contestant.v1.GetProblemRequest.
 * Use `create(GetProblemRequestSchema)` to create a new message.
 */
export const GetProblemRequestSchema: GenMessage<GetProblemRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 8);

/**
 * @generated from message contestant.v1.GetProblemResponse
 */
export type GetProblemResponse = Message<"contestant.v1.GetProblemResponse"> & {
  /**
   * @generated from field: contestant.v1.Problem problem = 1;
   */
  problem?: Problem;
};

/**
 * Describes the message contestant.v1.GetProblemResponse.
 * Use `create(GetProblemResponseSchema)` to create a new message.
 */
export const GetProblemResponseSchema: GenMessage<GetProblemResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 9);

/**
 * @generated from message contestant.v1.DeployRequest
 */
export type DeployRequest = Message<"contestant.v1.DeployRequest"> & {
  /**
   * @generated from field: string code = 1;
   */
  code: string;
};

/**
 * Describes the message contestant.v1.DeployRequest.
 * Use `create(DeployRequestSchema)` to create a new message.
 */
export const DeployRequestSchema: GenMessage<DeployRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 10);

/**
 * @generated from message contestant.v1.DeployResponse
 */
export type DeployResponse = Message<"contestant.v1.DeployResponse"> & {
};

/**
 * Describes the message contestant.v1.DeployResponse.
 * Use `create(DeployResponseSchema)` to create a new message.
 */
export const DeployResponseSchema: GenMessage<DeployResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 11);

/**
 * @generated from enum contestant.v1.DeploymentStatus
 */
export enum DeploymentStatus {
  /**
   * @generated from enum value: DEPLOYMENT_STATUS_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * 展開済み
   *
   * @generated from enum value: DEPLOYMENT_STATUS_DEPLOYED = 1;
   */
  DEPLOYED = 1,

  /**
   * 展開中
   *
   * @generated from enum value: DEPLOYMENT_STATUS_DEPLOYING = 2;
   */
  DEPLOYING = 2,

  /**
   * 展開失敗
   *
   * @generated from enum value: DEPLOYMENT_STATUS_FAILED = 3;
   */
  FAILED = 3,
}

/**
 * Describes the enum contestant.v1.DeploymentStatus.
 */
export const DeploymentStatusSchema: GenEnum<DeploymentStatus> = /*@__PURE__*/
  enumDesc(file_contestant_v1_problem, 0);

/**
 * @generated from enum contestant.v1.ProblemType
 */
export enum ProblemType {
  /**
   * @generated from enum value: PROBLEM_TYPE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: PROBLEM_TYPE_DESCRIPTIVE = 1;
   */
  DESCRIPTIVE = 1,
}

/**
 * Describes the enum contestant.v1.ProblemType.
 */
export const ProblemTypeSchema: GenEnum<ProblemType> = /*@__PURE__*/
  enumDesc(file_contestant_v1_problem, 1);

/**
 * @generated from service contestant.v1.ProblemService
 */
export const ProblemService: GenService<{
  /**
   * @generated from rpc contestant.v1.ProblemService.ListProblems
   */
  listProblems: {
    methodKind: "unary";
    input: typeof ListProblemsRequestSchema;
    output: typeof ListProblemsResponseSchema;
  },
  /**
   * @generated from rpc contestant.v1.ProblemService.GetProblem
   */
  getProblem: {
    methodKind: "unary";
    input: typeof GetProblemRequestSchema;
    output: typeof GetProblemResponseSchema;
  },
  /**
   * @generated from rpc contestant.v1.ProblemService.Deploy
   */
  deploy: {
    methodKind: "unary";
    input: typeof DeployRequestSchema;
    output: typeof DeployResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_contestant_v1_problem, 0);

