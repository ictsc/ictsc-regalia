// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file contestant/v1/contest.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/contest.proto.
 */
export const file_contestant_v1_contest: GenFile = /*@__PURE__*/
  fileDesc("Chtjb250ZXN0YW50L3YxL2NvbnRlc3QucHJvdG8SDWNvbnRlc3RhbnQudjEiwwEKCFNjaGVkdWxlEiMKBXBoYXNlGAEgASgOMhQuY29udGVzdGFudC52MS5QaGFzZRIoCgpuZXh0X3BoYXNlGAIgASgOMhQuY29udGVzdGFudC52MS5QaGFzZRIsCghzdGFydF9hdBgDIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXASLwoGZW5kX2F0GAQgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcEgAiAEBQgkKB19lbmRfYXQiFAoSR2V0U2NoZWR1bGVSZXF1ZXN0IkAKE0dldFNjaGVkdWxlUmVzcG9uc2USKQoIc2NoZWR1bGUYASABKAsyFy5jb250ZXN0YW50LnYxLlNjaGVkdWxlKngKBVBoYXNlEhUKEVBIQVNFX1VOU1BFQ0lGSUVEEAASGAoUUEhBU0VfT1VUX09GX0NPTlRFU1QQARIUChBQSEFTRV9JTl9DT05URVNUEAISDwoLUEhBU0VfQlJFQUsQAxIXChNQSEFTRV9BRlRFUl9DT05URVNUEAQyZgoOQ29udGVzdFNlcnZpY2USVAoLR2V0U2NoZWR1bGUSIS5jb250ZXN0YW50LnYxLkdldFNjaGVkdWxlUmVxdWVzdBoiLmNvbnRlc3RhbnQudjEuR2V0U2NoZWR1bGVSZXNwb25zZULDAQoRY29tLmNvbnRlc3RhbnQudjFCDENvbnRlc3RQcm90b1ABWktnaXRodWIuY29tL2ljdHNjL2ljdHNjLXJlZ2FsaWEvYmFja2VuZC9wa2cvcHJvdG8vY29udGVzdGFudC92MTtjb250ZXN0YW50djGiAgNDWFiqAg1Db250ZXN0YW50LlYxygINQ29udGVzdGFudFxWMeICGUNvbnRlc3RhbnRcVjFcR1BCTWV0YWRhdGHqAg5Db250ZXN0YW50OjpWMWIGcHJvdG8z", [file_google_protobuf_timestamp]);

/**
 * @generated from message contestant.v1.Schedule
 */
export type Schedule = Message<"contestant.v1.Schedule"> & {
  /**
   * @generated from field: contestant.v1.Phase phase = 1;
   */
  phase: Phase;

  /**
   * @generated from field: contestant.v1.Phase next_phase = 2;
   */
  nextPhase: Phase;

  /**
   * @generated from field: google.protobuf.Timestamp start_at = 3;
   */
  startAt?: Timestamp;

  /**
   * @generated from field: optional google.protobuf.Timestamp end_at = 4;
   */
  endAt?: Timestamp;
};

/**
 * Describes the message contestant.v1.Schedule.
 * Use `create(ScheduleSchema)` to create a new message.
 */
export const ScheduleSchema: GenMessage<Schedule> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contest, 0);

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
  messageDesc(file_contestant_v1_contest, 1);

/**
 * @generated from message contestant.v1.GetScheduleResponse
 */
export type GetScheduleResponse = Message<"contestant.v1.GetScheduleResponse"> & {
  /**
   * @generated from field: contestant.v1.Schedule schedule = 1;
   */
  schedule?: Schedule;
};

/**
 * Describes the message contestant.v1.GetScheduleResponse.
 * Use `create(GetScheduleResponseSchema)` to create a new message.
 */
export const GetScheduleResponseSchema: GenMessage<GetScheduleResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_contest, 2);

/**
 * @generated from enum contestant.v1.Phase
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
 * Describes the enum contestant.v1.Phase.
 */
export const PhaseSchema: GenEnum<Phase> = /*@__PURE__*/
  enumDesc(file_contestant_v1_contest, 0);

/**
 * @generated from service contestant.v1.ContestService
 */
export const ContestService: GenService<{
  /**
   * @generated from rpc contestant.v1.ContestService.GetSchedule
   */
  getSchedule: {
    methodKind: "unary";
    input: typeof GetScheduleRequestSchema;
    output: typeof GetScheduleResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_contestant_v1_contest, 0);

