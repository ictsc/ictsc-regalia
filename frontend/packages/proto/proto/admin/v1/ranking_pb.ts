// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file admin/v1/ranking.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Problem } from "./problem_pb";
import { file_admin_v1_problem } from "./problem_pb";
import type { Team } from "./team_pb";
import { file_admin_v1_team } from "./team_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/ranking.proto.
 */
export const file_admin_v1_ranking: GenFile = /*@__PURE__*/
  fileDesc("ChZhZG1pbi92MS9yYW5raW5nLnByb3RvEghhZG1pbi52MSJ/CgVTY29yZRIcCgR0ZWFtGAEgASgLMg4uYWRtaW4udjEuVGVhbRIiCgdwcm9ibGVtGAIgASgLMhEuYWRtaW4udjEuUHJvYmxlbRIUCgxtYXJrZWRfc2NvcmUYAyABKAMSDwoHcGVuYWx0eRgEIAEoAxINCgVzY29yZRgFIAEoAyJFCghUZWFtUmFuaxIcCgR0ZWFtGAEgASgLMg4uYWRtaW4udjEuVGVhbRIMCgRyYW5rGAIgASgDEg0KBXNjb3JlGAMgASgDIhIKEExpc3RTY29yZVJlcXVlc3QiNAoRTGlzdFNjb3JlUmVzcG9uc2USHwoGc2NvcmVzGAEgAygLMg8uYWRtaW4udjEuU2NvcmUiEwoRR2V0UmFua2luZ1JlcXVlc3QiOQoSR2V0UmFua2luZ1Jlc3BvbnNlEiMKB3JhbmtpbmcYASADKAsyEi5hZG1pbi52MS5UZWFtUmFuazKfAQoOUmFua2luZ1NlcnZpY2USRAoJTGlzdFNjb3JlEhouYWRtaW4udjEuTGlzdFNjb3JlUmVxdWVzdBobLmFkbWluLnYxLkxpc3RTY29yZVJlc3BvbnNlEkcKCkdldFJhbmtpbmcSGy5hZG1pbi52MS5HZXRSYW5raW5nUmVxdWVzdBocLmFkbWluLnYxLkdldFJhbmtpbmdSZXNwb25zZUKgAQoMY29tLmFkbWluLnYxQgxSYW5raW5nUHJvdG9QAVpBZ2l0aHViLmNvbS9pY3RzYy9pY3RzYy1yZWdhbGlhL2JhY2tlbmQvcGtnL3Byb3RvL2FkbWluL3YxO2FkbWludjGiAgNBWFiqAghBZG1pbi5WMcoCCEFkbWluXFYx4gIUQWRtaW5cVjFcR1BCTWV0YWRhdGHqAglBZG1pbjo6VjFiBnByb3RvMw", [file_admin_v1_problem, file_admin_v1_team]);

/**
 * @generated from message admin.v1.Score
 */
export type Score = Message<"admin.v1.Score"> & {
  /**
   * @generated from field: admin.v1.Team team = 1;
   */
  team?: Team;

  /**
   * @generated from field: admin.v1.Problem problem = 2;
   */
  problem?: Problem;

  /**
   * 採点による得点
   *
   * @generated from field: int64 marked_score = 3;
   */
  markedScore: bigint;

  /**
   * ペナルティによる減点
   *
   * @generated from field: int64 penalty = 4;
   */
  penalty: bigint;

  /**
   * 最終的な得点
   *
   * @generated from field: int64 score = 5;
   */
  score: bigint;
};

/**
 * Describes the message admin.v1.Score.
 * Use `create(ScoreSchema)` to create a new message.
 */
export const ScoreSchema: GenMessage<Score> = /*@__PURE__*/
  messageDesc(file_admin_v1_ranking, 0);

/**
 * @generated from message admin.v1.TeamRank
 */
export type TeamRank = Message<"admin.v1.TeamRank"> & {
  /**
   * @generated from field: admin.v1.Team team = 1;
   */
  team?: Team;

  /**
   * @generated from field: int64 rank = 2;
   */
  rank: bigint;

  /**
   * @generated from field: int64 score = 3;
   */
  score: bigint;
};

/**
 * Describes the message admin.v1.TeamRank.
 * Use `create(TeamRankSchema)` to create a new message.
 */
export const TeamRankSchema: GenMessage<TeamRank> = /*@__PURE__*/
  messageDesc(file_admin_v1_ranking, 1);

/**
 * @generated from message admin.v1.ListScoreRequest
 */
export type ListScoreRequest = Message<"admin.v1.ListScoreRequest"> & {
};

/**
 * Describes the message admin.v1.ListScoreRequest.
 * Use `create(ListScoreRequestSchema)` to create a new message.
 */
export const ListScoreRequestSchema: GenMessage<ListScoreRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_ranking, 2);

/**
 * @generated from message admin.v1.ListScoreResponse
 */
export type ListScoreResponse = Message<"admin.v1.ListScoreResponse"> & {
  /**
   * @generated from field: repeated admin.v1.Score scores = 1;
   */
  scores: Score[];
};

/**
 * Describes the message admin.v1.ListScoreResponse.
 * Use `create(ListScoreResponseSchema)` to create a new message.
 */
export const ListScoreResponseSchema: GenMessage<ListScoreResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_ranking, 3);

/**
 * @generated from message admin.v1.GetRankingRequest
 */
export type GetRankingRequest = Message<"admin.v1.GetRankingRequest"> & {
};

/**
 * Describes the message admin.v1.GetRankingRequest.
 * Use `create(GetRankingRequestSchema)` to create a new message.
 */
export const GetRankingRequestSchema: GenMessage<GetRankingRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_ranking, 4);

/**
 * @generated from message admin.v1.GetRankingResponse
 */
export type GetRankingResponse = Message<"admin.v1.GetRankingResponse"> & {
  /**
   * @generated from field: repeated admin.v1.TeamRank ranking = 1;
   */
  ranking: TeamRank[];
};

/**
 * Describes the message admin.v1.GetRankingResponse.
 * Use `create(GetRankingResponseSchema)` to create a new message.
 */
export const GetRankingResponseSchema: GenMessage<GetRankingResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_ranking, 5);

/**
 * @generated from service admin.v1.RankingService
 */
export const RankingService: GenService<{
  /**
   * @generated from rpc admin.v1.RankingService.ListScore
   */
  listScore: {
    methodKind: "unary";
    input: typeof ListScoreRequestSchema;
    output: typeof ListScoreResponseSchema;
  },
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

