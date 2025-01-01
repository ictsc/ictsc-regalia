// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file admin/v1/schedule.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/schedule.proto.
 */
export const file_admin_v1_schedule: GenFile = /*@__PURE__*/
  fileDesc("ChdhZG1pbi92MS9zY2hlZHVsZS5wcm90bxIIYWRtaW4udjEihAEKCFNjaGVkdWxlEh4KBXBoYXNlGAEgASgOMg8uYWRtaW4udjEuUGhhc2USLAoIc3RhcnRfYXQYAiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wEioKBmVuZF9hdBgDIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXAiFAoSR2V0U2NoZWR1bGVSZXF1ZXN0IjsKE0dldFNjaGVkdWxlUmVzcG9uc2USJAoIc2NoZWR1bGUYASADKAsyEi5hZG1pbi52MS5TY2hlZHVsZSI9ChVVcGRhdGVTY2hlZHVsZVJlcXVlc3QSJAoIc2NoZWR1bGUYASADKAsyEi5hZG1pbi52MS5TY2hlZHVsZSIYChZVcGRhdGVTY2hlZHVsZVJlc3BvbnNlKngKBVBoYXNlEhUKEVBIQVNFX1VOU1BFQ0lGSUVEEAASGAoUUEhBU0VfT1VUX09GX0NPTlRFU1QQARIUChBQSEFTRV9JTl9DT05URVNUEAISDwoLUEhBU0VfQlJFQUsQAxIXChNQSEFTRV9BRlRFUl9DT05URVNUEAQysgEKD1NjaGVkdWxlU2VydmljZRJKCgtHZXRTY2hlZHVsZRIcLmFkbWluLnYxLkdldFNjaGVkdWxlUmVxdWVzdBodLmFkbWluLnYxLkdldFNjaGVkdWxlUmVzcG9uc2USUwoOVXBkYXRlU2NoZWR1bGUSHy5hZG1pbi52MS5VcGRhdGVTY2hlZHVsZVJlcXVlc3QaIC5hZG1pbi52MS5VcGRhdGVTY2hlZHVsZVJlc3BvbnNlQqEBCgxjb20uYWRtaW4udjFCDVNjaGVkdWxlUHJvdG9QAVpBZ2l0aHViLmNvbS9pY3RzYy9pY3RzYy1yZWdhbGlhL2JhY2tlbmQvcGtnL3Byb3RvL2FkbWluL3YxO2FkbWludjGiAgNBWFiqAghBZG1pbi5WMcoCCEFkbWluXFYx4gIUQWRtaW5cVjFcR1BCTWV0YWRhdGHqAglBZG1pbjo6VjFiBnByb3RvMw", [file_google_protobuf_timestamp]);

/**
 * @generated from message admin.v1.Schedule
 */
export type Schedule = Message<"admin.v1.Schedule"> & {
  /**
   * @generated from field: admin.v1.Phase phase = 1;
   */
  phase: Phase;

  /**
   * @generated from field: google.protobuf.Timestamp start_at = 2;
   */
  startAt?: Timestamp;

  /**
   * @generated from field: google.protobuf.Timestamp end_at = 3;
   */
  endAt?: Timestamp;
};

/**
 * Describes the message admin.v1.Schedule.
 * Use `create(ScheduleSchema)` to create a new message.
 */
export const ScheduleSchema: GenMessage<Schedule> = /*@__PURE__*/
  messageDesc(file_admin_v1_schedule, 0);

/**
 * @generated from message admin.v1.GetScheduleRequest
 */
export type GetScheduleRequest = Message<"admin.v1.GetScheduleRequest"> & {
};

/**
 * Describes the message admin.v1.GetScheduleRequest.
 * Use `create(GetScheduleRequestSchema)` to create a new message.
 */
export const GetScheduleRequestSchema: GenMessage<GetScheduleRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_schedule, 1);

/**
 * @generated from message admin.v1.GetScheduleResponse
 */
export type GetScheduleResponse = Message<"admin.v1.GetScheduleResponse"> & {
  /**
   * @generated from field: repeated admin.v1.Schedule schedule = 1;
   */
  schedule: Schedule[];
};

/**
 * Describes the message admin.v1.GetScheduleResponse.
 * Use `create(GetScheduleResponseSchema)` to create a new message.
 */
export const GetScheduleResponseSchema: GenMessage<GetScheduleResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_schedule, 2);

/**
 * @generated from message admin.v1.UpdateScheduleRequest
 */
export type UpdateScheduleRequest = Message<"admin.v1.UpdateScheduleRequest"> & {
  /**
   * @generated from field: repeated admin.v1.Schedule schedule = 1;
   */
  schedule: Schedule[];
};

/**
 * Describes the message admin.v1.UpdateScheduleRequest.
 * Use `create(UpdateScheduleRequestSchema)` to create a new message.
 */
export const UpdateScheduleRequestSchema: GenMessage<UpdateScheduleRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_schedule, 3);

/**
 * @generated from message admin.v1.UpdateScheduleResponse
 */
export type UpdateScheduleResponse = Message<"admin.v1.UpdateScheduleResponse"> & {
};

/**
 * Describes the message admin.v1.UpdateScheduleResponse.
 * Use `create(UpdateScheduleResponseSchema)` to create a new message.
 */
export const UpdateScheduleResponseSchema: GenMessage<UpdateScheduleResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_schedule, 4);

/**
 * @generated from enum admin.v1.Phase
 */
export enum Phase {
  /**
   * @generated from enum value: PHASE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: PHASE_OUT_OF_CONTEST = 1;
   */
  OUT_OF_CONTEST = 1,

  /**
   * @generated from enum value: PHASE_IN_CONTEST = 2;
   */
  IN_CONTEST = 2,

  /**
   * @generated from enum value: PHASE_BREAK = 3;
   */
  BREAK = 3,

  /**
   * @generated from enum value: PHASE_AFTER_CONTEST = 4;
   */
  AFTER_CONTEST = 4,
}

/**
 * Describes the enum admin.v1.Phase.
 */
export const PhaseSchema: GenEnum<Phase> = /*@__PURE__*/
  enumDesc(file_admin_v1_schedule, 0);

/**
 * @generated from service admin.v1.ScheduleService
 */
export const ScheduleService: GenService<{
  /**
   * @generated from rpc admin.v1.ScheduleService.GetSchedule
   */
  getSchedule: {
    methodKind: "unary";
    input: typeof GetScheduleRequestSchema;
    output: typeof GetScheduleResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.ScheduleService.UpdateSchedule
   */
  updateSchedule: {
    methodKind: "unary";
    input: typeof UpdateScheduleRequestSchema;
    output: typeof UpdateScheduleResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_schedule, 0);

