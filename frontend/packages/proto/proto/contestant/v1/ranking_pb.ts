// @generated by protoc-gen-es v2.1.0 with parameter "target=ts"
// @generated from file contestant/v1/ranking.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Team } from "./team_pb";
import { file_contestant_v1_team } from "./team_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/ranking.proto.
 */
export const file_contestant_v1_ranking: GenFile = /*@__PURE__*/
  fileDesc("Chtjb250ZXN0YW50L3YxL3JhbmtpbmcucHJvdG8SDWNvbnRlc3RhbnQudjEiYAoEUmFuaxIVCgRyYW5rGAEgASgDQge6SAQiAiAAEikKBHRlYW0YAiABKAsyEy5jb250ZXN0YW50LnYxLlRlYW1CBrpIA8gBARIWCgVwb2ludBgDIAEoA0IHukgEIgIoACITChFHZXRSYW5raW5nUmVxdWVzdCJCChJHZXRSYW5raW5nUmVzcG9uc2USLAoHcmFua2luZxgBIAMoCzITLmNvbnRlc3RhbnQudjEuUmFua0IGukgDyAEBMmMKDlJhbmtpbmdTZXJ2aWNlElEKCkdldFJhbmtpbmcSIC5jb250ZXN0YW50LnYxLkdldFJhbmtpbmdSZXF1ZXN0GiEuY29udGVzdGFudC52MS5HZXRSYW5raW5nUmVzcG9uc2VCvAEKEWNvbS5jb250ZXN0YW50LnYxQgxSYW5raW5nUHJvdG9QAVpEZ2l0aHViLmNvbS9pY3RzYy9pY3RzYy1vdXRsYW5kcy9iYWNrZW5kL2ludGVybmFsL3Byb3RvL2NvbnRlc3RhbnQvdjGiAgNDWFiqAg1Db250ZXN0YW50LlYxygINQ29udGVzdGFudFxWMeICGUNvbnRlc3RhbnRcVjFcR1BCTWV0YWRhdGHqAg5Db250ZXN0YW50OjpWMWIGcHJvdG8z", [file_buf_validate_validate, file_contestant_v1_team]);

/**
 * @generated from message contestant.v1.Rank
 */
export type Rank = Message<"contestant.v1.Rank"> & {
  /**
   * @generated from field: int64 rank = 1;
   */
  rank: bigint;

  /**
   * @generated from field: contestant.v1.Team team = 2;
   */
  team?: Team;

  /**
   * @generated from field: int64 point = 3;
   */
  point: bigint;
};

/**
 * Describes the message contestant.v1.Rank.
 * Use `create(RankSchema)` to create a new message.
 */
export const RankSchema: GenMessage<Rank> = /*@__PURE__*/
  messageDesc(file_contestant_v1_ranking, 0);

/**
 * @generated from message contestant.v1.GetRankingRequest
 */
export type GetRankingRequest = Message<"contestant.v1.GetRankingRequest"> & {
};

/**
 * Describes the message contestant.v1.GetRankingRequest.
 * Use `create(GetRankingRequestSchema)` to create a new message.
 */
export const GetRankingRequestSchema: GenMessage<GetRankingRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_ranking, 1);

/**
 * @generated from message contestant.v1.GetRankingResponse
 */
export type GetRankingResponse = Message<"contestant.v1.GetRankingResponse"> & {
  /**
   * @generated from field: repeated contestant.v1.Rank ranking = 1;
   */
  ranking: Rank[];
};

/**
 * Describes the message contestant.v1.GetRankingResponse.
 * Use `create(GetRankingResponseSchema)` to create a new message.
 */
export const GetRankingResponseSchema: GenMessage<GetRankingResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_ranking, 2);

/**
 * @generated from service contestant.v1.RankingService
 */
export const RankingService: GenService<{
  /**
   * @generated from rpc contestant.v1.RankingService.GetRanking
   */
  getRanking: {
    methodKind: "unary";
    input: typeof GetRankingRequestSchema;
    output: typeof GetRankingResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_contestant_v1_ranking, 0);

