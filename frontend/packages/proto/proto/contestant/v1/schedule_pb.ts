// @generated by protoc-gen-es v2.2.0 with parameter "target=ts"
// @generated from file contestant/v1/schedule.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/schedule.proto.
 */
export const file_contestant_v1_schedule: GenFile = /*@__PURE__*/
  fileDesc("Chxjb250ZXN0YW50L3YxL3NjaGVkdWxlLnByb3RvEg1jb250ZXN0YW50LnYxIqsBCghTY2hlZHVsZRI2Cg1jdXJyZW50X3BoYXNlGAEgASgOMhUuY29udGVzdGFudC52MS5QaGFzZXNCCLpIBYIBAhABEjIKBmVuZF9hdBgCIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBCBrpIA8gBARIzCgpuZXh0X3BoYXNlGAMgASgOMhUuY29udGVzdGFudC52MS5QaGFzZXNCCLpIBYIBAhABIhQKEkdldFNjaGVkdWxlUmVxdWVzdCJJChNHZXRTY2hlZHVsZVJlc3BvbnNlEjIKCXNjaGVkdWxlcxgBIAMoCzIXLmNvbnRlc3RhbnQudjEuU2NoZWR1bGVCBrpIA8gBASpKCgZQaGFzZXMSFgoSUEhBU0VTX1VOU1BFQ0lGSUVEEAASFQoRUEhBU0VTX1FVQUxJRllJTkcQARIRCg1QSEFTRVNfRklOQUxTEAIyZwoPU2NoZWR1bGVTZXJ2aWNlElQKC0dldFNjaGVkdWxlEiEuY29udGVzdGFudC52MS5HZXRTY2hlZHVsZVJlcXVlc3QaIi5jb250ZXN0YW50LnYxLkdldFNjaGVkdWxlUmVzcG9uc2VCvQEKEWNvbS5jb250ZXN0YW50LnYxQg1TY2hlZHVsZVByb3RvUAFaRGdpdGh1Yi5jb20vaWN0c2MvaWN0c2Mtb3V0bGFuZHMvYmFja2VuZC9pbnRlcm5hbC9wcm90by9jb250ZXN0YW50L3YxogIDQ1hYqgINQ29udGVzdGFudC5WMcoCDUNvbnRlc3RhbnRcVjHiAhlDb250ZXN0YW50XFYxXEdQQk1ldGFkYXRh6gIOQ29udGVzdGFudDo6VjFiBnByb3RvMw", [file_buf_validate_validate, file_google_protobuf_timestamp]);

/**
 * @generated from message contestant.v1.Schedule
 */
export type Schedule = Message<"contestant.v1.Schedule"> & {
  /**
   * @generated from field: contestant.v1.Phases current_phase = 1;
   */
  currentPhase: Phases;

  /**
   * @generated from field: google.protobuf.Timestamp end_at = 2;
   */
  endAt?: Timestamp;

  /**
   * @generated from field: contestant.v1.Phases next_phase = 3;
   */
  nextPhase: Phases;
};

/**
 * Describes the message contestant.v1.Schedule.
 * Use `create(ScheduleSchema)` to create a new message.
 */
export const ScheduleSchema: GenMessage<Schedule> = /*@__PURE__*/
  messageDesc(file_contestant_v1_schedule, 0);

/**
 * @generated from message contestant.v1.GetScheduleRequest
 */
export type GetScheduleRequest = Message<"contestant.v1.GetScheduleRequest"> & {
};

/**
 * Describes the message contestant.v1.GetScheduleRequest.
 * Use `create(GetScheduleRequestSchema)` to create a new message.
 */
export const GetScheduleRequestSchema: GenMessage<GetScheduleRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_schedule, 1);

/**
 * @generated from message contestant.v1.GetScheduleResponse
 */
export type GetScheduleResponse = Message<"contestant.v1.GetScheduleResponse"> & {
  /**
   * @generated from field: repeated contestant.v1.Schedule schedules = 1;
   */
  schedules: Schedule[];
};

/**
 * Describes the message contestant.v1.GetScheduleResponse.
 * Use `create(GetScheduleResponseSchema)` to create a new message.
 */
export const GetScheduleResponseSchema: GenMessage<GetScheduleResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_schedule, 2);

/**
 * @generated from enum contestant.v1.Phases
 */
export enum Phases {
  /**
   * @generated from enum value: PHASES_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: PHASES_QUALIFYING = 1;
   */
  QUALIFYING = 1,

  /**
   * @generated from enum value: PHASES_FINALS = 2;
   */
  FINALS = 2,
}

/**
 * Describes the enum contestant.v1.Phases.
 */
export const PhasesSchema: GenEnum<Phases> = /*@__PURE__*/
  enumDesc(file_contestant_v1_schedule, 0);

/**
 * @generated from service contestant.v1.ScheduleService
 */
export const ScheduleService: GenService<{
  /**
   * @generated from rpc contestant.v1.ScheduleService.GetSchedule
   */
  getSchedule: {
    methodKind: "unary";
    input: typeof GetScheduleRequestSchema;
    output: typeof GetScheduleResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_contestant_v1_schedule, 0);

