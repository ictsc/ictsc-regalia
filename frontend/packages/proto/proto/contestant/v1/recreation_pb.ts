// @generated by protoc-gen-es v2.1.0 with parameter "target=ts"
// @generated from file contestant/v1/recreation.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/recreation.proto.
 */
export const file_contestant_v1_recreation: GenFile = /*@__PURE__*/
  fileDesc("Ch5jb250ZXN0YW50L3YxL3JlY3JlYXRpb24ucHJvdG8SDWNvbnRlc3RhbnQudjEiqQEKClJlY3JlYXRpb24SFAoCaWQYASABKAlCCLpIBXIDmAEaEhwKCnByb2JsZW1faWQYAiABKAlCCLpIBXIDmAEaEi8KBnN0YXR1cxgDIAEoDjIVLmNvbnRlc3RhbnQudjEuU3RhdHVzQgi6SAWCAQIQARI2CgpjcmVhdGVkX2F0GAQgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcEIGukgDyAEBIjUKFUdldFJlY3JlYXRpb25zUmVxdWVzdBIcCgpwcm9ibGVtX2lkGAEgASgJQgi6SAVyA5gBGiJQChZHZXRSZWNyZWF0aW9uc1Jlc3BvbnNlEjYKC3JlY3JlYXRpb25zGAEgAygLMhkuY29udGVzdGFudC52MS5SZWNyZWF0aW9uQga6SAPIAQEiNQoVUG9zdFJlY3JlYXRpb25SZXF1ZXN0EhwKCnByb2JsZW1faWQYASABKAlCCLpIBXIDmAEaIk8KFlBvc3RSZWNyZWF0aW9uUmVzcG9uc2USNQoKcmVjcmVhdGlvbhgBIAEoCzIZLmNvbnRlc3RhbnQudjEuUmVjcmVhdGlvbkIGukgDyAEBKk4KBlN0YXR1cxIWChJTVEFUVVNfVU5TUEVDSUZJRUQQABIWChJTVEFUVVNfSU5fUFJPR1JFU1MQARIUChBTVEFUVVNfQ09NUExFVEVEEAIy0QEKEVJlY3JlYXRpb25TZXJ2aWNlEl0KDkdldFJlY3JlYXRpb25zEiQuY29udGVzdGFudC52MS5HZXRSZWNyZWF0aW9uc1JlcXVlc3QaJS5jb250ZXN0YW50LnYxLkdldFJlY3JlYXRpb25zUmVzcG9uc2USXQoOUG9zdFJlY3JlYXRpb24SJC5jb250ZXN0YW50LnYxLlBvc3RSZWNyZWF0aW9uUmVxdWVzdBolLmNvbnRlc3RhbnQudjEuUG9zdFJlY3JlYXRpb25SZXNwb25zZUK/AQoRY29tLmNvbnRlc3RhbnQudjFCD1JlY3JlYXRpb25Qcm90b1ABWkRnaXRodWIuY29tL2ljdHNjL2ljdHNjLW91dGxhbmRzL2JhY2tlbmQvaW50ZXJuYWwvcHJvdG8vY29udGVzdGFudC92MaICA0NYWKoCDUNvbnRlc3RhbnQuVjHKAg1Db250ZXN0YW50XFYx4gIZQ29udGVzdGFudFxWMVxHUEJNZXRhZGF0YeoCDkNvbnRlc3RhbnQ6OlYxYgZwcm90bzM", [file_buf_validate_validate, file_google_protobuf_timestamp]);

/**
 * @generated from message contestant.v1.Recreation
 */
export type Recreation = Message<"contestant.v1.Recreation"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string problem_id = 2;
   */
  problemId: string;

  /**
   * @generated from field: contestant.v1.Status status = 3;
   */
  status: Status;

  /**
   * @generated from field: google.protobuf.Timestamp created_at = 4;
   */
  createdAt?: Timestamp;
};

/**
 * Describes the message contestant.v1.Recreation.
 * Use `create(RecreationSchema)` to create a new message.
 */
export const RecreationSchema: GenMessage<Recreation> = /*@__PURE__*/
  messageDesc(file_contestant_v1_recreation, 0);

/**
 * @generated from message contestant.v1.GetRecreationsRequest
 */
export type GetRecreationsRequest = Message<"contestant.v1.GetRecreationsRequest"> & {
  /**
   * @generated from field: string problem_id = 1;
   */
  problemId: string;
};

/**
 * Describes the message contestant.v1.GetRecreationsRequest.
 * Use `create(GetRecreationsRequestSchema)` to create a new message.
 */
export const GetRecreationsRequestSchema: GenMessage<GetRecreationsRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_recreation, 1);

/**
 * @generated from message contestant.v1.GetRecreationsResponse
 */
export type GetRecreationsResponse = Message<"contestant.v1.GetRecreationsResponse"> & {
  /**
   * @generated from field: repeated contestant.v1.Recreation recreations = 1;
   */
  recreations: Recreation[];
};

/**
 * Describes the message contestant.v1.GetRecreationsResponse.
 * Use `create(GetRecreationsResponseSchema)` to create a new message.
 */
export const GetRecreationsResponseSchema: GenMessage<GetRecreationsResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_recreation, 2);

/**
 * @generated from message contestant.v1.PostRecreationRequest
 */
export type PostRecreationRequest = Message<"contestant.v1.PostRecreationRequest"> & {
  /**
   * @generated from field: string problem_id = 1;
   */
  problemId: string;
};

/**
 * Describes the message contestant.v1.PostRecreationRequest.
 * Use `create(PostRecreationRequestSchema)` to create a new message.
 */
export const PostRecreationRequestSchema: GenMessage<PostRecreationRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_recreation, 3);

/**
 * @generated from message contestant.v1.PostRecreationResponse
 */
export type PostRecreationResponse = Message<"contestant.v1.PostRecreationResponse"> & {
  /**
   * @generated from field: contestant.v1.Recreation recreation = 1;
   */
  recreation?: Recreation;
};

/**
 * Describes the message contestant.v1.PostRecreationResponse.
 * Use `create(PostRecreationResponseSchema)` to create a new message.
 */
export const PostRecreationResponseSchema: GenMessage<PostRecreationResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_recreation, 4);

/**
 * @generated from enum contestant.v1.Status
 */
export enum Status {
  /**
   * @generated from enum value: STATUS_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: STATUS_IN_PROGRESS = 1;
   */
  IN_PROGRESS = 1,

  /**
   * @generated from enum value: STATUS_COMPLETED = 2;
   */
  COMPLETED = 2,
}

/**
 * Describes the enum contestant.v1.Status.
 */
export const StatusSchema: GenEnum<Status> = /*@__PURE__*/
  enumDesc(file_contestant_v1_recreation, 0);

/**
 * @generated from service contestant.v1.RecreationService
 */
export const RecreationService: GenService<{
  /**
   * @generated from rpc contestant.v1.RecreationService.GetRecreations
   */
  getRecreations: {
    methodKind: "unary";
    input: typeof GetRecreationsRequestSchema;
    output: typeof GetRecreationsResponseSchema;
  },
  /**
   * @generated from rpc contestant.v1.RecreationService.PostRecreation
   */
  postRecreation: {
    methodKind: "unary";
    input: typeof PostRecreationRequestSchema;
    output: typeof PostRecreationResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_contestant_v1_recreation, 0);

