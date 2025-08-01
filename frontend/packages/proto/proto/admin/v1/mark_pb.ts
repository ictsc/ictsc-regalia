// @generated by protoc-gen-es v2.6.1 with parameter "target=ts"
// @generated from file admin/v1/mark.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv2";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv2";
import type { Admin } from "./actor_pb";
import { file_admin_v1_actor } from "./actor_pb";
import type { Contestant } from "./contestant_pb";
import { file_admin_v1_contestant } from "./contestant_pb";
import type { Problem, ProblemType } from "./problem_pb";
import { file_admin_v1_problem } from "./problem_pb";
import type { Team } from "./team_pb";
import { file_admin_v1_team } from "./team_pb";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/mark.proto.
 */
export const file_admin_v1_mark: GenFile = /*@__PURE__*/
  fileDesc("ChNhZG1pbi92MS9tYXJrLnByb3RvEghhZG1pbi52MSKGAgoGQW5zd2VyEgoKAmlkGAEgASgNEhwKBHRlYW0YAiABKAsyDi5hZG1pbi52MS5UZWFtEiQKBmF1dGhvchgDIAEoCzIULmFkbWluLnYxLkNvbnRlc3RhbnQSIgoHcHJvYmxlbRgEIAEoCzIRLmFkbWluLnYxLlByb2JsZW0SIgoEYm9keRgFIAEoCzIULmFkbWluLnYxLkFuc3dlckJvZHkSLgoKY3JlYXRlZF9hdBgGIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXASKgoFc2NvcmUYByABKAsyFi5hZG1pbi52MS5NYXJraW5nU2NvcmVIAIgBAUIICgZfc2NvcmUibQoKQW5zd2VyQm9keRIjCgR0eXBlGAEgASgOMhUuYWRtaW4udjEuUHJvYmxlbVR5cGUSMgoLZGVzY3JpcHRpdmUYAiABKAsyGy5hZG1pbi52MS5EZXNjcmlwdGl2ZUFuc3dlckgAQgYKBGJvZHkiIQoRRGVzY3JpcHRpdmVBbnN3ZXISDAoEYm9keRgBIAEoCSL2AQoNTWFya2luZ1Jlc3VsdBIgCgZhbnN3ZXIYASABKAsyEC5hZG1pbi52MS5BbnN3ZXISHgoFanVkZ2UYAiABKAsyDy5hZG1pbi52MS5BZG1pbhINCgVzY29yZRgDIAEoDRItCglyYXRpb25hbGUYBCABKAsyGi5hZG1pbi52MS5NYXJraW5nUmF0aW9uYWxlEi4KCmNyZWF0ZWRfYXQYBSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wEjUKCnZpc2liaWxpdHkYByABKA4yIS5hZG1pbi52MS5NYXJraW5nUmVzdWx0VmlzaWJpbGl0eSJLCgxNYXJraW5nU2NvcmUSDQoFdG90YWwYASABKA0SDgoGbWFya2VkGAIgASgNEg8KB3BlbmFsdHkYAyABKA0SCwoDbWF4GAQgASgNIn0KEE1hcmtpbmdSYXRpb25hbGUSIwoEdHlwZRgBIAEoDjIVLmFkbWluLnYxLlByb2JsZW1UeXBlEjwKC2Rlc2NyaXB0aXZlGAIgASgLMiUuYWRtaW4udjEuRGVzY3JpcHRpdmVNYXJraW5nUmF0aW9uYWxlSABCBgoEYm9keSIuChtEZXNjcmlwdGl2ZU1hcmtpbmdSYXRpb25hbGUSDwoHY29tbWVudBgBIAEoCSIsChJMaXN0QW5zd2Vyc1JlcXVlc3QSFgoOaW5jbHVkZV9tYXJrZWQYASABKAgiOAoTTGlzdEFuc3dlcnNSZXNwb25zZRIhCgdhbnN3ZXJzGAEgAygLMhAuYWRtaW4udjEuQW5zd2VyIkcKEEdldEFuc3dlclJlcXVlc3QSEQoJdGVhbV9jb2RlGAEgASgNEhQKDHByb2JsZW1fY29kZRgCIAEoCRIKCgJpZBgDIAEoDSI1ChFHZXRBbnN3ZXJSZXNwb25zZRIgCgZhbnN3ZXIYASABKAsyEC5hZG1pbi52MS5BbnN3ZXIiUgoZTGlzdE1hcmtpbmdSZXN1bHRzUmVxdWVzdBI1Cgp2aXNpYmlsaXR5GAEgASgOMiEuYWRtaW4udjEuTWFya2luZ1Jlc3VsdFZpc2liaWxpdHkiTgoaTGlzdE1hcmtpbmdSZXN1bHRzUmVzcG9uc2USMAoPbWFya2luZ19yZXN1bHRzGAEgAygLMhcuYWRtaW4udjEuTWFya2luZ1Jlc3VsdCJNChpDcmVhdGVNYXJraW5nUmVzdWx0UmVxdWVzdBIvCg5tYXJraW5nX3Jlc3VsdBgBIAEoCzIXLmFkbWluLnYxLk1hcmtpbmdSZXN1bHQiTgobQ3JlYXRlTWFya2luZ1Jlc3VsdFJlc3BvbnNlEi8KDm1hcmtpbmdfcmVzdWx0GAEgASgLMhcuYWRtaW4udjEuTWFya2luZ1Jlc3VsdCIoCiZVcGRhdGVNYXJraW5nUmVzdWx0VmlzaWJpbGl0aWVzUmVxdWVzdCIpCidVcGRhdGVNYXJraW5nUmVzdWx0VmlzaWJpbGl0aWVzUmVzcG9uc2UiFQoTVXBkYXRlU2NvcmVzUmVxdWVzdCIWChRVcGRhdGVTY29yZXNSZXNwb25zZSqRAQoXTWFya2luZ1Jlc3VsdFZpc2liaWxpdHkSKQolTUFSS0lOR19SRVNVTFRfVklTSUJJTElUWV9VTlNQRUNJRklFRBAAEiQKIE1BUktJTkdfUkVTVUxUX1ZJU0lCSUxJVFlfUFVCTElDEAESJQohTUFSS0lOR19SRVNVTFRfVklTSUJJTElUWV9QUklWQVRFEAIyvAQKC01hcmtTZXJ2aWNlEkoKC0xpc3RBbnN3ZXJzEhwuYWRtaW4udjEuTGlzdEFuc3dlcnNSZXF1ZXN0Gh0uYWRtaW4udjEuTGlzdEFuc3dlcnNSZXNwb25zZRJECglHZXRBbnN3ZXISGi5hZG1pbi52MS5HZXRBbnN3ZXJSZXF1ZXN0GhsuYWRtaW4udjEuR2V0QW5zd2VyUmVzcG9uc2USXwoSTGlzdE1hcmtpbmdSZXN1bHRzEiMuYWRtaW4udjEuTGlzdE1hcmtpbmdSZXN1bHRzUmVxdWVzdBokLmFkbWluLnYxLkxpc3RNYXJraW5nUmVzdWx0c1Jlc3BvbnNlEmIKE0NyZWF0ZU1hcmtpbmdSZXN1bHQSJC5hZG1pbi52MS5DcmVhdGVNYXJraW5nUmVzdWx0UmVxdWVzdBolLmFkbWluLnYxLkNyZWF0ZU1hcmtpbmdSZXN1bHRSZXNwb25zZRKGAQofVXBkYXRlTWFya2luZ1Jlc3VsdFZpc2liaWxpdGllcxIwLmFkbWluLnYxLlVwZGF0ZU1hcmtpbmdSZXN1bHRWaXNpYmlsaXRpZXNSZXF1ZXN0GjEuYWRtaW4udjEuVXBkYXRlTWFya2luZ1Jlc3VsdFZpc2liaWxpdGllc1Jlc3BvbnNlEk0KDFVwZGF0ZVNjb3JlcxIdLmFkbWluLnYxLlVwZGF0ZVNjb3Jlc1JlcXVlc3QaHi5hZG1pbi52MS5VcGRhdGVTY29yZXNSZXNwb25zZUKdAQoMY29tLmFkbWluLnYxQglNYXJrUHJvdG9QAVpBZ2l0aHViLmNvbS9pY3RzYy9pY3RzYy1yZWdhbGlhL2JhY2tlbmQvcGtnL3Byb3RvL2FkbWluL3YxO2FkbWludjGiAgNBWFiqAghBZG1pbi5WMcoCCEFkbWluXFYx4gIUQWRtaW5cVjFcR1BCTWV0YWRhdGHqAglBZG1pbjo6VjFiBnByb3RvMw", [file_admin_v1_actor, file_admin_v1_contestant, file_admin_v1_problem, file_admin_v1_team, file_google_protobuf_timestamp]);

/**
 * 解答
 *
 * @generated from message admin.v1.Answer
 */
export type Answer = Message<"admin.v1.Answer"> & {
  /**
   * @generated from field: uint32 id = 1;
   */
  id: number;

  /**
   * @generated from field: admin.v1.Team team = 2;
   */
  team?: Team;

  /**
   * @generated from field: admin.v1.Contestant author = 3;
   */
  author?: Contestant;

  /**
   * @generated from field: admin.v1.Problem problem = 4;
   */
  problem?: Problem;

  /**
   * @generated from field: admin.v1.AnswerBody body = 5;
   */
  body?: AnswerBody;

  /**
   * @generated from field: google.protobuf.Timestamp created_at = 6;
   */
  createdAt?: Timestamp;

  /**
   * @generated from field: optional admin.v1.MarkingScore score = 7;
   */
  score?: MarkingScore;
};

/**
 * Describes the message admin.v1.Answer.
 * Use `create(AnswerSchema)` to create a new message.
 */
export const AnswerSchema: GenMessage<Answer> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 0);

/**
 * @generated from message admin.v1.AnswerBody
 */
export type AnswerBody = Message<"admin.v1.AnswerBody"> & {
  /**
   * @generated from field: admin.v1.ProblemType type = 1;
   */
  type: ProblemType;

  /**
   * @generated from oneof admin.v1.AnswerBody.body
   */
  body: {
    /**
     * @generated from field: admin.v1.DescriptiveAnswer descriptive = 2;
     */
    value: DescriptiveAnswer;
    case: "descriptive";
  } | { case: undefined; value?: undefined };
};

/**
 * Describes the message admin.v1.AnswerBody.
 * Use `create(AnswerBodySchema)` to create a new message.
 */
export const AnswerBodySchema: GenMessage<AnswerBody> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 1);

/**
 * @generated from message admin.v1.DescriptiveAnswer
 */
export type DescriptiveAnswer = Message<"admin.v1.DescriptiveAnswer"> & {
  /**
   * @generated from field: string body = 1;
   */
  body: string;
};

/**
 * Describes the message admin.v1.DescriptiveAnswer.
 * Use `create(DescriptiveAnswerSchema)` to create a new message.
 */
export const DescriptiveAnswerSchema: GenMessage<DescriptiveAnswer> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 2);

/**
 * 採点結果
 *
 * @generated from message admin.v1.MarkingResult
 */
export type MarkingResult = Message<"admin.v1.MarkingResult"> & {
  /**
   * @generated from field: admin.v1.Answer answer = 1;
   */
  answer?: Answer;

  /**
   * @generated from field: admin.v1.Admin judge = 2;
   */
  judge?: Admin;

  /**
   * @generated from field: uint32 score = 3;
   */
  score: number;

  /**
   * @generated from field: admin.v1.MarkingRationale rationale = 4;
   */
  rationale?: MarkingRationale;

  /**
   * @generated from field: google.protobuf.Timestamp created_at = 5;
   */
  createdAt?: Timestamp;

  /**
   * @generated from field: admin.v1.MarkingResultVisibility visibility = 7;
   */
  visibility: MarkingResultVisibility;
};

/**
 * Describes the message admin.v1.MarkingResult.
 * Use `create(MarkingResultSchema)` to create a new message.
 */
export const MarkingResultSchema: GenMessage<MarkingResult> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 3);

/**
 * @generated from message admin.v1.MarkingScore
 */
export type MarkingScore = Message<"admin.v1.MarkingScore"> & {
  /**
   * @generated from field: uint32 total = 1;
   */
  total: number;

  /**
   * @generated from field: uint32 marked = 2;
   */
  marked: number;

  /**
   * @generated from field: uint32 penalty = 3;
   */
  penalty: number;

  /**
   * @generated from field: uint32 max = 4;
   */
  max: number;
};

/**
 * Describes the message admin.v1.MarkingScore.
 * Use `create(MarkingScoreSchema)` to create a new message.
 */
export const MarkingScoreSchema: GenMessage<MarkingScore> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 4);

/**
 * @generated from message admin.v1.MarkingRationale
 */
export type MarkingRationale = Message<"admin.v1.MarkingRationale"> & {
  /**
   * @generated from field: admin.v1.ProblemType type = 1;
   */
  type: ProblemType;

  /**
   * @generated from oneof admin.v1.MarkingRationale.body
   */
  body: {
    /**
     * @generated from field: admin.v1.DescriptiveMarkingRationale descriptive = 2;
     */
    value: DescriptiveMarkingRationale;
    case: "descriptive";
  } | { case: undefined; value?: undefined };
};

/**
 * Describes the message admin.v1.MarkingRationale.
 * Use `create(MarkingRationaleSchema)` to create a new message.
 */
export const MarkingRationaleSchema: GenMessage<MarkingRationale> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 5);

/**
 * @generated from message admin.v1.DescriptiveMarkingRationale
 */
export type DescriptiveMarkingRationale = Message<"admin.v1.DescriptiveMarkingRationale"> & {
  /**
   * @generated from field: string comment = 1;
   */
  comment: string;
};

/**
 * Describes the message admin.v1.DescriptiveMarkingRationale.
 * Use `create(DescriptiveMarkingRationaleSchema)` to create a new message.
 */
export const DescriptiveMarkingRationaleSchema: GenMessage<DescriptiveMarkingRationale> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 6);

/**
 * @generated from message admin.v1.ListAnswersRequest
 */
export type ListAnswersRequest = Message<"admin.v1.ListAnswersRequest"> & {
  /**
   * 採点が完了した提出を含めるかどうか
   *
   * @generated from field: bool include_marked = 1;
   */
  includeMarked: boolean;
};

/**
 * Describes the message admin.v1.ListAnswersRequest.
 * Use `create(ListAnswersRequestSchema)` to create a new message.
 */
export const ListAnswersRequestSchema: GenMessage<ListAnswersRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 7);

/**
 * @generated from message admin.v1.ListAnswersResponse
 */
export type ListAnswersResponse = Message<"admin.v1.ListAnswersResponse"> & {
  /**
   * @generated from field: repeated admin.v1.Answer answers = 1;
   */
  answers: Answer[];
};

/**
 * Describes the message admin.v1.ListAnswersResponse.
 * Use `create(ListAnswersResponseSchema)` to create a new message.
 */
export const ListAnswersResponseSchema: GenMessage<ListAnswersResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 8);

/**
 * @generated from message admin.v1.GetAnswerRequest
 */
export type GetAnswerRequest = Message<"admin.v1.GetAnswerRequest"> & {
  /**
   * @generated from field: uint32 team_code = 1;
   */
  teamCode: number;

  /**
   * @generated from field: string problem_code = 2;
   */
  problemCode: string;

  /**
   * @generated from field: uint32 id = 3;
   */
  id: number;
};

/**
 * Describes the message admin.v1.GetAnswerRequest.
 * Use `create(GetAnswerRequestSchema)` to create a new message.
 */
export const GetAnswerRequestSchema: GenMessage<GetAnswerRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 9);

/**
 * @generated from message admin.v1.GetAnswerResponse
 */
export type GetAnswerResponse = Message<"admin.v1.GetAnswerResponse"> & {
  /**
   * @generated from field: admin.v1.Answer answer = 1;
   */
  answer?: Answer;
};

/**
 * Describes the message admin.v1.GetAnswerResponse.
 * Use `create(GetAnswerResponseSchema)` to create a new message.
 */
export const GetAnswerResponseSchema: GenMessage<GetAnswerResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 10);

/**
 * @generated from message admin.v1.ListMarkingResultsRequest
 */
export type ListMarkingResultsRequest = Message<"admin.v1.ListMarkingResultsRequest"> & {
  /**
   * @generated from field: admin.v1.MarkingResultVisibility visibility = 1;
   */
  visibility: MarkingResultVisibility;
};

/**
 * Describes the message admin.v1.ListMarkingResultsRequest.
 * Use `create(ListMarkingResultsRequestSchema)` to create a new message.
 */
export const ListMarkingResultsRequestSchema: GenMessage<ListMarkingResultsRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 11);

/**
 * @generated from message admin.v1.ListMarkingResultsResponse
 */
export type ListMarkingResultsResponse = Message<"admin.v1.ListMarkingResultsResponse"> & {
  /**
   * @generated from field: repeated admin.v1.MarkingResult marking_results = 1;
   */
  markingResults: MarkingResult[];
};

/**
 * Describes the message admin.v1.ListMarkingResultsResponse.
 * Use `create(ListMarkingResultsResponseSchema)` to create a new message.
 */
export const ListMarkingResultsResponseSchema: GenMessage<ListMarkingResultsResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 12);

/**
 * @generated from message admin.v1.CreateMarkingResultRequest
 */
export type CreateMarkingResultRequest = Message<"admin.v1.CreateMarkingResultRequest"> & {
  /**
   * @generated from field: admin.v1.MarkingResult marking_result = 1;
   */
  markingResult?: MarkingResult;
};

/**
 * Describes the message admin.v1.CreateMarkingResultRequest.
 * Use `create(CreateMarkingResultRequestSchema)` to create a new message.
 */
export const CreateMarkingResultRequestSchema: GenMessage<CreateMarkingResultRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 13);

/**
 * @generated from message admin.v1.CreateMarkingResultResponse
 */
export type CreateMarkingResultResponse = Message<"admin.v1.CreateMarkingResultResponse"> & {
  /**
   * @generated from field: admin.v1.MarkingResult marking_result = 1;
   */
  markingResult?: MarkingResult;
};

/**
 * Describes the message admin.v1.CreateMarkingResultResponse.
 * Use `create(CreateMarkingResultResponseSchema)` to create a new message.
 */
export const CreateMarkingResultResponseSchema: GenMessage<CreateMarkingResultResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 14);

/**
 * ビジネスロジックに沿って、採点結果の可視性を更新する
 *
 * @generated from message admin.v1.UpdateMarkingResultVisibilitiesRequest
 */
export type UpdateMarkingResultVisibilitiesRequest = Message<"admin.v1.UpdateMarkingResultVisibilitiesRequest"> & {
};

/**
 * Describes the message admin.v1.UpdateMarkingResultVisibilitiesRequest.
 * Use `create(UpdateMarkingResultVisibilitiesRequestSchema)` to create a new message.
 */
export const UpdateMarkingResultVisibilitiesRequestSchema: GenMessage<UpdateMarkingResultVisibilitiesRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 15);

/**
 * @generated from message admin.v1.UpdateMarkingResultVisibilitiesResponse
 */
export type UpdateMarkingResultVisibilitiesResponse = Message<"admin.v1.UpdateMarkingResultVisibilitiesResponse"> & {
};

/**
 * Describes the message admin.v1.UpdateMarkingResultVisibilitiesResponse.
 * Use `create(UpdateMarkingResultVisibilitiesResponseSchema)` to create a new message.
 */
export const UpdateMarkingResultVisibilitiesResponseSchema: GenMessage<UpdateMarkingResultVisibilitiesResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 16);

/**
 * @generated from message admin.v1.UpdateScoresRequest
 */
export type UpdateScoresRequest = Message<"admin.v1.UpdateScoresRequest"> & {
};

/**
 * Describes the message admin.v1.UpdateScoresRequest.
 * Use `create(UpdateScoresRequestSchema)` to create a new message.
 */
export const UpdateScoresRequestSchema: GenMessage<UpdateScoresRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 17);

/**
 * @generated from message admin.v1.UpdateScoresResponse
 */
export type UpdateScoresResponse = Message<"admin.v1.UpdateScoresResponse"> & {
};

/**
 * Describes the message admin.v1.UpdateScoresResponse.
 * Use `create(UpdateScoresResponseSchema)` to create a new message.
 */
export const UpdateScoresResponseSchema: GenMessage<UpdateScoresResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 18);

/**
 * @generated from enum admin.v1.MarkingResultVisibility
 */
export enum MarkingResultVisibility {
  /**
   * @generated from enum value: MARKING_RESULT_VISIBILITY_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * 参加者に見える
   *
   * @generated from enum value: MARKING_RESULT_VISIBILITY_PUBLIC = 1;
   */
  PUBLIC = 1,

  /**
   * 参加者に見えない
   *
   * @generated from enum value: MARKING_RESULT_VISIBILITY_PRIVATE = 2;
   */
  PRIVATE = 2,
}

/**
 * Describes the enum admin.v1.MarkingResultVisibility.
 */
export const MarkingResultVisibilitySchema: GenEnum<MarkingResultVisibility> = /*@__PURE__*/
  enumDesc(file_admin_v1_mark, 0);

/**
 * @generated from service admin.v1.MarkService
 */
export const MarkService: GenService<{
  /**
   * @generated from rpc admin.v1.MarkService.ListAnswers
   */
  listAnswers: {
    methodKind: "unary";
    input: typeof ListAnswersRequestSchema;
    output: typeof ListAnswersResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.MarkService.GetAnswer
   */
  getAnswer: {
    methodKind: "unary";
    input: typeof GetAnswerRequestSchema;
    output: typeof GetAnswerResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.MarkService.ListMarkingResults
   */
  listMarkingResults: {
    methodKind: "unary";
    input: typeof ListMarkingResultsRequestSchema;
    output: typeof ListMarkingResultsResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.MarkService.CreateMarkingResult
   */
  createMarkingResult: {
    methodKind: "unary";
    input: typeof CreateMarkingResultRequestSchema;
    output: typeof CreateMarkingResultResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.MarkService.UpdateMarkingResultVisibilities
   */
  updateMarkingResultVisibilities: {
    methodKind: "unary";
    input: typeof UpdateMarkingResultVisibilitiesRequestSchema;
    output: typeof UpdateMarkingResultVisibilitiesResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.MarkService.UpdateScores
   */
  updateScores: {
    methodKind: "unary";
    input: typeof UpdateScoresRequestSchema;
    output: typeof UpdateScoresResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_mark, 0);

