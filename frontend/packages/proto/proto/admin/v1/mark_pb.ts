// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file admin/v1/mark.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
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
  fileDesc("ChNhZG1pbi92MS9tYXJrLnByb3RvEghhZG1pbi52MSKJAgoGQW5zd2VyEgoKAmlkGAEgASgNEhwKBHRlYW0YAiABKAsyDi5hZG1pbi52MS5UZWFtEiQKBmF1dGhvchgDIAEoCzIULmFkbWluLnYxLkNvbnRlc3RhbnQSIgoHcHJvYmxlbRgEIAEoCzIRLmFkbWluLnYxLlByb2JsZW0SIgoEYm9keRgFIAEoCzIULmFkbWluLnYxLkFuc3dlckJvZHkSLgoKY3JlYXRlZF9hdBgGIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXASLAoGcmVzdWx0GAcgASgLMhcuYWRtaW4udjEuTWFya2luZ1Jlc3VsdEgAiAEBQgkKB19yZXN1bHQibQoKQW5zd2VyQm9keRIjCgR0eXBlGAEgASgOMhUuYWRtaW4udjEuUHJvYmxlbVR5cGUSMgoLZGVzY3JpcHRpdmUYAiABKAsyGy5hZG1pbi52MS5EZXNjcmlwdGl2ZUFuc3dlckgAQgYKBGJvZHkiIQoRRGVzY3JpcHRpdmVBbnN3ZXISDAoEYm9keRgBIAEoCSK/AQoNTWFya2luZ1Jlc3VsdBIgCgZhbnN3ZXIYASABKAsyEC5hZG1pbi52MS5BbnN3ZXISHgoFanVkZ2UYAiABKAsyDy5hZG1pbi52MS5BZG1pbhINCgVzY29yZRgDIAEoDRItCglyYXRpb25hbGUYBCABKAsyGi5hZG1pbi52MS5NYXJraW5nUmF0aW9uYWxlEi4KCmNyZWF0ZWRfYXQYBSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wIn0KEE1hcmtpbmdSYXRpb25hbGUSIwoEdHlwZRgBIAEoDjIVLmFkbWluLnYxLlByb2JsZW1UeXBlEjwKC2Rlc2NyaXB0aXZlGAIgASgLMiUuYWRtaW4udjEuRGVzY3JpcHRpdmVNYXJraW5nUmF0aW9uYWxlSABCBgoEYm9keSIuChtEZXNjcmlwdGl2ZU1hcmtpbmdSYXRpb25hbGUSDwoHY29tbWVudBgBIAEoCSIsChJMaXN0QW5zd2Vyc1JlcXVlc3QSFgoOaW5jbHVkZV9tYXJrZWQYASABKAgiOAoTTGlzdEFuc3dlcnNSZXNwb25zZRIhCgdhbnN3ZXJzGAEgAygLMhAuYWRtaW4udjEuQW5zd2VyIkcKEEdldEFuc3dlclJlcXVlc3QSEQoJdGVhbV9jb2RlGAEgASgNEhQKDHByb2JsZW1fY29kZRgCIAEoCRIKCgJpZBgDIAEoDSI1ChFHZXRBbnN3ZXJSZXNwb25zZRIgCgZhbnN3ZXIYASABKAsyEC5hZG1pbi52MS5BbnN3ZXIiGwoZTGlzdE1hcmtpbmdSZXN1bHRzUmVxdWVzdCJOChpMaXN0TWFya2luZ1Jlc3VsdHNSZXNwb25zZRIwCg9tYXJraW5nX3Jlc3VsdHMYASADKAsyFy5hZG1pbi52MS5NYXJraW5nUmVzdWx0Ik0KGkNyZWF0ZU1hcmtpbmdSZXN1bHRSZXF1ZXN0Ei8KDm1hcmtpbmdfcmVzdWx0GAEgASgLMhcuYWRtaW4udjEuTWFya2luZ1Jlc3VsdCJOChtDcmVhdGVNYXJraW5nUmVzdWx0UmVzcG9uc2USLwoObWFya2luZ19yZXN1bHQYASABKAsyFy5hZG1pbi52MS5NYXJraW5nUmVzdWx0MuQCCgtNYXJrU2VydmljZRJKCgtMaXN0QW5zd2VycxIcLmFkbWluLnYxLkxpc3RBbnN3ZXJzUmVxdWVzdBodLmFkbWluLnYxLkxpc3RBbnN3ZXJzUmVzcG9uc2USRAoJR2V0QW5zd2VyEhouYWRtaW4udjEuR2V0QW5zd2VyUmVxdWVzdBobLmFkbWluLnYxLkdldEFuc3dlclJlc3BvbnNlEl8KEkxpc3RNYXJraW5nUmVzdWx0cxIjLmFkbWluLnYxLkxpc3RNYXJraW5nUmVzdWx0c1JlcXVlc3QaJC5hZG1pbi52MS5MaXN0TWFya2luZ1Jlc3VsdHNSZXNwb25zZRJiChNDcmVhdGVNYXJraW5nUmVzdWx0EiQuYWRtaW4udjEuQ3JlYXRlTWFya2luZ1Jlc3VsdFJlcXVlc3QaJS5hZG1pbi52MS5DcmVhdGVNYXJraW5nUmVzdWx0UmVzcG9uc2VCnQEKDGNvbS5hZG1pbi52MUIJTWFya1Byb3RvUAFaQWdpdGh1Yi5jb20vaWN0c2MvaWN0c2MtcmVnYWxpYS9iYWNrZW5kL3BrZy9wcm90by9hZG1pbi92MTthZG1pbnYxogIDQVhYqgIIQWRtaW4uVjHKAghBZG1pblxWMeICFEFkbWluXFYxXEdQQk1ldGFkYXRh6gIJQWRtaW46OlYxYgZwcm90bzM", [file_admin_v1_actor, file_admin_v1_contestant, file_admin_v1_problem, file_admin_v1_team, file_google_protobuf_timestamp]);

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
   * @generated from field: optional admin.v1.MarkingResult result = 7;
   */
  result?: MarkingResult;
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
};

/**
 * Describes the message admin.v1.MarkingResult.
 * Use `create(MarkingResultSchema)` to create a new message.
 */
export const MarkingResultSchema: GenMessage<MarkingResult> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 3);

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
  messageDesc(file_admin_v1_mark, 4);

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
  messageDesc(file_admin_v1_mark, 5);

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
  messageDesc(file_admin_v1_mark, 6);

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
  messageDesc(file_admin_v1_mark, 7);

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
  messageDesc(file_admin_v1_mark, 8);

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
  messageDesc(file_admin_v1_mark, 9);

/**
 * @generated from message admin.v1.ListMarkingResultsRequest
 */
export type ListMarkingResultsRequest = Message<"admin.v1.ListMarkingResultsRequest"> & {
};

/**
 * Describes the message admin.v1.ListMarkingResultsRequest.
 * Use `create(ListMarkingResultsRequestSchema)` to create a new message.
 */
export const ListMarkingResultsRequestSchema: GenMessage<ListMarkingResultsRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_mark, 10);

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
  messageDesc(file_admin_v1_mark, 11);

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
  messageDesc(file_admin_v1_mark, 12);

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
  messageDesc(file_admin_v1_mark, 13);

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
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_mark, 0);

