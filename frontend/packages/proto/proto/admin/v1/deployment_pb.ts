// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file admin/v1/deployment.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/deployment.proto.
 */
export const file_admin_v1_deployment: GenFile = /*@__PURE__*/
  fileDesc("ChlhZG1pbi92MS9kZXBsb3ltZW50LnByb3RvEghhZG1pbi52MSKnAQoKRGVwbG95bWVudBIRCgl0ZWFtX2NvZGUYASABKAMSFAoMcHJvYmxlbV9jb2RlGAIgASgJEhAKCHJldmlzaW9uGAMgASgDEjMKDGxhdGVzdF9ldmVudBgEIAEoDjIdLmFkbWluLnYxLkRlcGxveW1lbnRFdmVudFR5cGUSKQoGZXZlbnRzGAUgAygLMhkuYWRtaW4udjEuRGVwbG95bWVudEV2ZW50Im8KD0RlcGxveW1lbnRFdmVudBIvCgtvY2N1cnJlZF9hdBgBIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXASKwoEdHlwZRgCIAEoDjIdLmFkbWluLnYxLkRlcGxveW1lbnRFdmVudFR5cGUiQQoWTGlzdERlcGxveW1lbnRzUmVxdWVzdBIRCgl0ZWFtX2NvZGUYASABKAMSFAoMcHJvYmxlbV9jb2RlGAIgASgJIkQKF0xpc3REZXBsb3ltZW50c1Jlc3BvbnNlEikKC2RlcGxveW1lbnRzGAEgAygLMhQuYWRtaW4udjEuRGVwbG95bWVudCJAChVTeW5jRGVwbG95bWVudFJlcXVlc3QSEQoJdGVhbV9jb2RlGAEgASgDEhQKDHByb2JsZW1fY29kZRgCIAEoCSIYChZTeW5jRGVwbG95bWVudFJlc3BvbnNlIjgKDURlcGxveVJlcXVlc3QSEQoJdGVhbV9jb2RlGAEgASgDEhQKDHByb2JsZW1fY29kZRgCIAEoCSI6Cg5EZXBsb3lSZXNwb25zZRIoCgpkZXBsb3ltZW50GAEgASgLMhQuYWRtaW4udjEuRGVwbG95bWVudCrHAQoTRGVwbG95bWVudEV2ZW50VHlwZRIlCiFERVBMT1lNRU5UX0VWRU5UX1RZUEVfVU5TUEVDSUZJRUQQABIgChxERVBMT1lNRU5UX0VWRU5UX1RZUEVfUVVFVUVEEAESIgoeREVQTE9ZTUVOVF9FVkVOVF9UWVBFX0NSRUFUSU5HEAISIgoeREVQTE9ZTUVOVF9FVkVOVF9UWVBFX0ZJTklTSEVEEAMSHwobREVQTE9ZTUVOVF9FVkVOVF9UWVBFX0VSUk9SEAQy/QEKEURlcGxveW1lbnRTZXJ2aWNlElYKD0xpc3REZXBsb3ltZW50cxIgLmFkbWluLnYxLkxpc3REZXBsb3ltZW50c1JlcXVlc3QaIS5hZG1pbi52MS5MaXN0RGVwbG95bWVudHNSZXNwb25zZRJTCg5TeW5jRGVwbG95bWVudBIfLmFkbWluLnYxLlN5bmNEZXBsb3ltZW50UmVxdWVzdBogLmFkbWluLnYxLlN5bmNEZXBsb3ltZW50UmVzcG9uc2USOwoGRGVwbG95EhcuYWRtaW4udjEuRGVwbG95UmVxdWVzdBoYLmFkbWluLnYxLkRlcGxveVJlc3BvbnNlQqMBCgxjb20uYWRtaW4udjFCD0RlcGxveW1lbnRQcm90b1ABWkFnaXRodWIuY29tL2ljdHNjL2ljdHNjLXJlZ2FsaWEvYmFja2VuZC9wa2cvcHJvdG8vYWRtaW4vdjE7YWRtaW52MaICA0FYWKoCCEFkbWluLlYxygIIQWRtaW5cVjHiAhRBZG1pblxWMVxHUEJNZXRhZGF0YeoCCUFkbWluOjpWMWIGcHJvdG8z", [file_google_protobuf_timestamp]);

/**
 * 問題の展開状態
 *
 * @generated from message admin.v1.Deployment
 */
export type Deployment = Message<"admin.v1.Deployment"> & {
  /**
   * チームコード
   *
   * @generated from field: int64 team_code = 1;
   */
  teamCode: bigint;

  /**
   * 問題コード
   *
   * @generated from field: string problem_code = 2;
   */
  problemCode: string;

  /**
   * リビジョン - 0 が初期状態で再展開される度にインクリメントされる
   *
   * @generated from field: int64 revision = 3;
   */
  revision: bigint;

  /**
   * 最新のイベント
   *
   * @generated from field: admin.v1.DeploymentEventType latest_event = 4;
   */
  latestEvent: DeploymentEventType;

  /**
   * イベント
   *
   * @generated from field: repeated admin.v1.DeploymentEvent events = 5;
   */
  events: DeploymentEvent[];
};

/**
 * Describes the message admin.v1.Deployment.
 * Use `create(DeploymentSchema)` to create a new message.
 */
export const DeploymentSchema: GenMessage<Deployment> = /*@__PURE__*/
  messageDesc(file_admin_v1_deployment, 0);

/**
 * 問題展開に関するイベント
 *
 * @generated from message admin.v1.DeploymentEvent
 */
export type DeploymentEvent = Message<"admin.v1.DeploymentEvent"> & {
  /**
   * @generated from field: google.protobuf.Timestamp occurred_at = 1;
   */
  occurredAt?: Timestamp;

  /**
   * @generated from field: admin.v1.DeploymentEventType type = 2;
   */
  type: DeploymentEventType;
};

/**
 * Describes the message admin.v1.DeploymentEvent.
 * Use `create(DeploymentEventSchema)` to create a new message.
 */
export const DeploymentEventSchema: GenMessage<DeploymentEvent> = /*@__PURE__*/
  messageDesc(file_admin_v1_deployment, 1);

/**
 * @generated from message admin.v1.ListDeploymentsRequest
 */
export type ListDeploymentsRequest = Message<"admin.v1.ListDeploymentsRequest"> & {
  /**
   * @generated from field: int64 team_code = 1;
   */
  teamCode: bigint;

  /**
   * @generated from field: string problem_code = 2;
   */
  problemCode: string;
};

/**
 * Describes the message admin.v1.ListDeploymentsRequest.
 * Use `create(ListDeploymentsRequestSchema)` to create a new message.
 */
export const ListDeploymentsRequestSchema: GenMessage<ListDeploymentsRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_deployment, 2);

/**
 * @generated from message admin.v1.ListDeploymentsResponse
 */
export type ListDeploymentsResponse = Message<"admin.v1.ListDeploymentsResponse"> & {
  /**
   * @generated from field: repeated admin.v1.Deployment deployments = 1;
   */
  deployments: Deployment[];
};

/**
 * Describes the message admin.v1.ListDeploymentsResponse.
 * Use `create(ListDeploymentsResponseSchema)` to create a new message.
 */
export const ListDeploymentsResponseSchema: GenMessage<ListDeploymentsResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_deployment, 3);

/**
 * @generated from message admin.v1.SyncDeploymentRequest
 */
export type SyncDeploymentRequest = Message<"admin.v1.SyncDeploymentRequest"> & {
  /**
   * @generated from field: int64 team_code = 1;
   */
  teamCode: bigint;

  /**
   * @generated from field: string problem_code = 2;
   */
  problemCode: string;
};

/**
 * Describes the message admin.v1.SyncDeploymentRequest.
 * Use `create(SyncDeploymentRequestSchema)` to create a new message.
 */
export const SyncDeploymentRequestSchema: GenMessage<SyncDeploymentRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_deployment, 4);

/**
 * @generated from message admin.v1.SyncDeploymentResponse
 */
export type SyncDeploymentResponse = Message<"admin.v1.SyncDeploymentResponse"> & {
};

/**
 * Describes the message admin.v1.SyncDeploymentResponse.
 * Use `create(SyncDeploymentResponseSchema)` to create a new message.
 */
export const SyncDeploymentResponseSchema: GenMessage<SyncDeploymentResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_deployment, 5);

/**
 * @generated from message admin.v1.DeployRequest
 */
export type DeployRequest = Message<"admin.v1.DeployRequest"> & {
  /**
   * @generated from field: int64 team_code = 1;
   */
  teamCode: bigint;

  /**
   * @generated from field: string problem_code = 2;
   */
  problemCode: string;
};

/**
 * Describes the message admin.v1.DeployRequest.
 * Use `create(DeployRequestSchema)` to create a new message.
 */
export const DeployRequestSchema: GenMessage<DeployRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_deployment, 6);

/**
 * @generated from message admin.v1.DeployResponse
 */
export type DeployResponse = Message<"admin.v1.DeployResponse"> & {
  /**
   * @generated from field: admin.v1.Deployment deployment = 1;
   */
  deployment?: Deployment;
};

/**
 * Describes the message admin.v1.DeployResponse.
 * Use `create(DeployResponseSchema)` to create a new message.
 */
export const DeployResponseSchema: GenMessage<DeployResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_deployment, 7);

/**
 * @generated from enum admin.v1.DeploymentEventType
 */
export enum DeploymentEventType {
  /**
   * @generated from enum value: DEPLOYMENT_EVENT_TYPE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: DEPLOYMENT_EVENT_TYPE_QUEUED = 1;
   */
  QUEUED = 1,

  /**
   * @generated from enum value: DEPLOYMENT_EVENT_TYPE_CREATING = 2;
   */
  CREATING = 2,

  /**
   * @generated from enum value: DEPLOYMENT_EVENT_TYPE_FINISHED = 3;
   */
  FINISHED = 3,

  /**
   * @generated from enum value: DEPLOYMENT_EVENT_TYPE_ERROR = 4;
   */
  ERROR = 4,
}

/**
 * Describes the enum admin.v1.DeploymentEventType.
 */
export const DeploymentEventTypeSchema: GenEnum<DeploymentEventType> = /*@__PURE__*/
  enumDesc(file_admin_v1_deployment, 0);

/**
 * @generated from service admin.v1.DeploymentService
 */
export const DeploymentService: GenService<{
  /**
   * @generated from rpc admin.v1.DeploymentService.ListDeployments
   */
  listDeployments: {
    methodKind: "unary";
    input: typeof ListDeploymentsRequestSchema;
    output: typeof ListDeploymentsResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.DeploymentService.SyncDeployment
   */
  syncDeployment: {
    methodKind: "unary";
    input: typeof SyncDeploymentRequestSchema;
    output: typeof SyncDeploymentResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.DeploymentService.Deploy
   */
  deploy: {
    methodKind: "unary";
    input: typeof DeployRequestSchema;
    output: typeof DeployResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_deployment, 0);
