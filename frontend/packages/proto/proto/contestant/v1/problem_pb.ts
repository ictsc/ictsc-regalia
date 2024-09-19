// @generated by protoc-gen-es v2.1.0 with parameter "target=ts"
// @generated from file contestant/v1/problem.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/problem.proto.
 */
export const file_contestant_v1_problem: GenFile = /*@__PURE__*/
  fileDesc("Chtjb250ZXN0YW50L3YxL3Byb2JsZW0ucHJvdG8SDWNvbnRlc3RhbnQudjEiOgoGQ2hvaWNlEhYKBWluZGV4GAEgASgDQge6SAQiAigAEhgKBGJvZHkYAiABKAlCCrpIB3IFEAEY6AciuwEKCFF1ZXN0aW9uEhQKAmlkGAEgASgJQgi6SAVyA5gBGhIYCgRib2R5GAIgASgJQgq6SAdyBRABGOgHEjMKBHR5cGUYAyABKA4yGy5jb250ZXN0YW50LnYxLlF1ZXN0aW9uVHlwZUIIukgFggECEAESMgoHY2hvaWNlcxgEIAMoCzIVLmNvbnRlc3RhbnQudjEuQ2hvaWNlQgq6SAeSAQQIAhAKEhYKBXBvaW50GAUgASgDQge6SAQiAiAAIncKFU11bHRpcGxlQ2hvaWNlUHJvYmxlbRIdCgRib2R5GAEgASgJQgq6SAdyBRABGOgHSACIAQESNgoJcXVlc3Rpb25zGAIgAygLMhcuY29udGVzdGFudC52MS5RdWVzdGlvbkIKukgHkgEECAEQCkIHCgVfYm9keSJ/Cg5Db25uZWN0aW9uSW5mbxIbCghob3N0bmFtZRgBIAEoCUIJukgGcgQQARhkEhoKB2NvbW1hbmQYAiABKAlCCbpIBnIEEAEYZBIbCghwYXNzd29yZBgDIAEoCUIJukgGcgQQARgUEhcKBHR5cGUYBCABKAlCCbpIBnIEEAEYFCJzChJEZXNjcmlwdGl2ZVByb2JsZW0SGAoEYm9keRgBIAEoCUIKukgHcgUQARjoBxJDChBjb25uZWN0aW9uX2luZm9zGAIgAygLMh0uY29udGVzdGFudC52MS5Db25uZWN0aW9uSW5mb0IKukgHkgEECAAQCiKnAgoHUHJvYmxlbRIUCgJpZBgBIAEoCUIIukgFcgOYARoSGAoFdGl0bGUYAiABKAlCCbpIBnIEEAEYZBIWCgRjb2RlGAMgASgJQgi6SAVyA5gBAxIWCgVwb2ludBgEIAEoA0IHukgEIgIgABIyCgR0eXBlGAUgASgOMhouY29udGVzdGFudC52MS5Qcm9ibGVtVHlwZUIIukgFggECEAESOAoLZGVzY3JpcHRpdmUYBiABKAsyIS5jb250ZXN0YW50LnYxLkRlc2NyaXB0aXZlUHJvYmxlbUgAEj8KD211bHRpcGxlX2Nob2ljZRgHIAEoCzIkLmNvbnRlc3RhbnQudjEuTXVsdGlwbGVDaG9pY2VQcm9ibGVtSABCDQoEYm9keRIFukgCCAEiFAoSR2V0UHJvYmxlbXNSZXF1ZXN0IkcKE0dldFByb2JsZW1zUmVzcG9uc2USMAoIcHJvYmxlbXMYASADKAsyFi5jb250ZXN0YW50LnYxLlByb2JsZW1CBrpIA8gBASpiCgxRdWVzdGlvblR5cGUSHQoZUVVFU1RJT05fVFlQRV9VTlNQRUNJRklFRBAAEhcKE1FVRVNUSU9OX1RZUEVfUkFESU8QARIaChZRVUVTVElPTl9UWVBFX0NIRUNLQk9YEAIqawoLUHJvYmxlbVR5cGUSHAoYUFJPQkxFTV9UWVBFX1VOU1BFQ0lGSUVEEAASHAoYUFJPQkxFTV9UWVBFX0RFU0NSSVBUSVZFEAESIAocUFJPQkxFTV9UWVBFX01VTFRJUExFX0NIT0lDRRACMmYKDlByb2JsZW1TZXJ2aWNlElQKC0dldFByb2JsZW1zEiEuY29udGVzdGFudC52MS5HZXRQcm9ibGVtc1JlcXVlc3QaIi5jb250ZXN0YW50LnYxLkdldFByb2JsZW1zUmVzcG9uc2VCvAEKEWNvbS5jb250ZXN0YW50LnYxQgxQcm9ibGVtUHJvdG9QAVpEZ2l0aHViLmNvbS9pY3RzYy9pY3RzYy1vdXRsYW5kcy9iYWNrZW5kL2ludGVybmFsL3Byb3RvL2NvbnRlc3RhbnQvdjGiAgNDWFiqAg1Db250ZXN0YW50LlYxygINQ29udGVzdGFudFxWMeICGUNvbnRlc3RhbnRcVjFcR1BCTWV0YWRhdGHqAg5Db250ZXN0YW50OjpWMWIGcHJvdG8z", [file_buf_validate_validate]);

/**
 * @generated from message contestant.v1.Choice
 */
export type Choice = Message<"contestant.v1.Choice"> & {
  /**
   * @generated from field: int64 index = 1;
   */
  index: bigint;

  /**
   * @generated from field: string body = 2;
   */
  body: string;
};

/**
 * Describes the message contestant.v1.Choice.
 * Use `create(ChoiceSchema)` to create a new message.
 */
export const ChoiceSchema: GenMessage<Choice> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 0);

/**
 * @generated from message contestant.v1.Question
 */
export type Question = Message<"contestant.v1.Question"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string body = 2;
   */
  body: string;

  /**
   * @generated from field: contestant.v1.QuestionType type = 3;
   */
  type: QuestionType;

  /**
   * @generated from field: repeated contestant.v1.Choice choices = 4;
   */
  choices: Choice[];

  /**
   * @generated from field: int64 point = 5;
   */
  point: bigint;
};

/**
 * Describes the message contestant.v1.Question.
 * Use `create(QuestionSchema)` to create a new message.
 */
export const QuestionSchema: GenMessage<Question> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 1);

/**
 * @generated from message contestant.v1.MultipleChoiceProblem
 */
export type MultipleChoiceProblem = Message<"contestant.v1.MultipleChoiceProblem"> & {
  /**
   * @generated from field: optional string body = 1;
   */
  body?: string;

  /**
   * @generated from field: repeated contestant.v1.Question questions = 2;
   */
  questions: Question[];
};

/**
 * Describes the message contestant.v1.MultipleChoiceProblem.
 * Use `create(MultipleChoiceProblemSchema)` to create a new message.
 */
export const MultipleChoiceProblemSchema: GenMessage<MultipleChoiceProblem> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 2);

/**
 * @generated from message contestant.v1.ConnectionInfo
 */
export type ConnectionInfo = Message<"contestant.v1.ConnectionInfo"> & {
  /**
   * @generated from field: string hostname = 1;
   */
  hostname: string;

  /**
   * @generated from field: string command = 2;
   */
  command: string;

  /**
   * @generated from field: string password = 3;
   */
  password: string;

  /**
   * @generated from field: string type = 4;
   */
  type: string;
};

/**
 * Describes the message contestant.v1.ConnectionInfo.
 * Use `create(ConnectionInfoSchema)` to create a new message.
 */
export const ConnectionInfoSchema: GenMessage<ConnectionInfo> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 3);

/**
 * @generated from message contestant.v1.DescriptiveProblem
 */
export type DescriptiveProblem = Message<"contestant.v1.DescriptiveProblem"> & {
  /**
   * @generated from field: string body = 1;
   */
  body: string;

  /**
   * @generated from field: repeated contestant.v1.ConnectionInfo connection_infos = 2;
   */
  connectionInfos: ConnectionInfo[];
};

/**
 * Describes the message contestant.v1.DescriptiveProblem.
 * Use `create(DescriptiveProblemSchema)` to create a new message.
 */
export const DescriptiveProblemSchema: GenMessage<DescriptiveProblem> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 4);

/**
 * @generated from message contestant.v1.Problem
 */
export type Problem = Message<"contestant.v1.Problem"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string title = 2;
   */
  title: string;

  /**
   * @generated from field: string code = 3;
   */
  code: string;

  /**
   * @generated from field: int64 point = 4;
   */
  point: bigint;

  /**
   * @generated from field: contestant.v1.ProblemType type = 5;
   */
  type: ProblemType;

  /**
   * @generated from oneof contestant.v1.Problem.body
   */
  body: {
    /**
     * @generated from field: contestant.v1.DescriptiveProblem descriptive = 6;
     */
    value: DescriptiveProblem;
    case: "descriptive";
  } | {
    /**
     * @generated from field: contestant.v1.MultipleChoiceProblem multiple_choice = 7;
     */
    value: MultipleChoiceProblem;
    case: "multipleChoice";
  } | { case: undefined; value?: undefined };
};

/**
 * Describes the message contestant.v1.Problem.
 * Use `create(ProblemSchema)` to create a new message.
 */
export const ProblemSchema: GenMessage<Problem> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 5);

/**
 * @generated from message contestant.v1.GetProblemsRequest
 */
export type GetProblemsRequest = Message<"contestant.v1.GetProblemsRequest"> & {
};

/**
 * Describes the message contestant.v1.GetProblemsRequest.
 * Use `create(GetProblemsRequestSchema)` to create a new message.
 */
export const GetProblemsRequestSchema: GenMessage<GetProblemsRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 6);

/**
 * @generated from message contestant.v1.GetProblemsResponse
 */
export type GetProblemsResponse = Message<"contestant.v1.GetProblemsResponse"> & {
  /**
   * @generated from field: repeated contestant.v1.Problem problems = 1;
   */
  problems: Problem[];
};

/**
 * Describes the message contestant.v1.GetProblemsResponse.
 * Use `create(GetProblemsResponseSchema)` to create a new message.
 */
export const GetProblemsResponseSchema: GenMessage<GetProblemsResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_problem, 7);

/**
 * @generated from enum contestant.v1.QuestionType
 */
export enum QuestionType {
  /**
   * @generated from enum value: QUESTION_TYPE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: QUESTION_TYPE_RADIO = 1;
   */
  RADIO = 1,

  /**
   * @generated from enum value: QUESTION_TYPE_CHECKBOX = 2;
   */
  CHECKBOX = 2,
}

/**
 * Describes the enum contestant.v1.QuestionType.
 */
export const QuestionTypeSchema: GenEnum<QuestionType> = /*@__PURE__*/
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

  /**
   * @generated from enum value: PROBLEM_TYPE_MULTIPLE_CHOICE = 2;
   */
  MULTIPLE_CHOICE = 2,
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
   * @generated from rpc contestant.v1.ProblemService.GetProblems
   */
  getProblems: {
    methodKind: "unary";
    input: typeof GetProblemsRequestSchema;
    output: typeof GetProblemsResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_contestant_v1_problem, 0);

