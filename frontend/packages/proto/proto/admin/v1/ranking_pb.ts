// @generated by protoc-gen-es v2.0.0 with parameter "target=ts"
// @generated from file admin/v1/ranking.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Team } from "./team_pb";
import { file_admin_v1_team } from "./team_pb";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/ranking.proto.
 */
export const file_admin_v1_ranking: GenFile = /*@__PURE__*/
  fileDesc("ChZhZG1pbi92MS9yYW5raW5nLnByb3RvEghhZG1pbi52MSJbCgRSYW5rEhUKBHJhbmsYASABKANCB7pIBCICIAASJAoEdGVhbRgCIAEoCzIOLmFkbWluLnYxLlRlYW1CBrpIA8gBARIWCgVwb2ludBgDIAEoA0IHukgEIgIoACIoChFHZXRSYW5raW5nUmVxdWVzdBITCgt1bnB1Ymxpc2hlZBgBIAEoCCI9ChJHZXRSYW5raW5nUmVzcG9uc2USJwoHcmFua2luZxgBIAMoCzIOLmFkbWluLnYxLlJhbmtCBrpIA8gBATJZCg5SYW5raW5nU2VydmljZRJHCgpHZXRSYW5raW5nEhsuYWRtaW4udjEuR2V0UmFua2luZ1JlcXVlc3QaHC5hZG1pbi52MS5HZXRSYW5raW5nUmVzcG9uc2VCngEKDGNvbS5hZG1pbi52MUIMUmFua2luZ1Byb3RvUAFaP2dpdGh1Yi5jb20vaWN0c2MvaWN0c2Mtb3V0bGFuZHMvYmFja2VuZC9pbnRlcm5hbC9wcm90by9hZG1pbi92MaICA0FYWKoCCEFkbWluLlYxygIIQWRtaW5cVjHiAhRBZG1pblxWMVxHUEJNZXRhZGF0YeoCCUFkbWluOjpWMWIGcHJvdG8z", [file_admin_v1_team, file_buf_validate_validate]);

/**
 * @generated from message admin.v1.Rank
 */
export type Rank = Message<"admin.v1.Rank"> & {
  /**
   * @generated from field: int64 rank = 1;
   */
  rank: bigint;

  /**
   * @generated from field: admin.v1.Team team = 2;
   */
  team?: Team;

  /**
   * @generated from field: int64 point = 3;
   */
  point: bigint;
};

/**
 * Describes the message admin.v1.Rank.
 * Use `create(RankSchema)` to create a new message.
 */
export const RankSchema: GenMessage<Rank> = /*@__PURE__*/
  messageDesc(file_admin_v1_ranking, 0);

/**
 * @generated from message admin.v1.GetRankingRequest
 */
export type GetRankingRequest = Message<"admin.v1.GetRankingRequest"> & {
  /**
   * @generated from field: bool unpublished = 1;
   */
  unpublished: boolean;
};

/**
 * Describes the message admin.v1.GetRankingRequest.
 * Use `create(GetRankingRequestSchema)` to create a new message.
 */
export const GetRankingRequestSchema: GenMessage<GetRankingRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_ranking, 1);

/**
 * @generated from message admin.v1.GetRankingResponse
 */
export type GetRankingResponse = Message<"admin.v1.GetRankingResponse"> & {
  /**
   * @generated from field: repeated admin.v1.Rank ranking = 1;
   */
  ranking: Rank[];
};

/**
 * Describes the message admin.v1.GetRankingResponse.
 * Use `create(GetRankingResponseSchema)` to create a new message.
 */
export const GetRankingResponseSchema: GenMessage<GetRankingResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_ranking, 2);

/**
 * @generated from service admin.v1.RankingService
 */
export const RankingService: GenService<{
  /**
   * @generated from rpc admin.v1.RankingService.GetRanking
   */
  getRanking: {
    methodKind: "unary";
    input: typeof GetRankingRequestSchema;
    output: typeof GetRankingResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_ranking, 0);

