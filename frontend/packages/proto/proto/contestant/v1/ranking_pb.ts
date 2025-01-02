// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file contestant/v1/ranking.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/ranking.proto.
 */
export const file_contestant_v1_ranking: GenFile = /*@__PURE__*/
  fileDesc("Chtjb250ZXN0YW50L3YxL3JhbmtpbmcucHJvdG8SDWNvbnRlc3RhbnQudjEiNgoEUmFuaxIMCgRyYW5rGAEgASgDEhEKCXRlYW1fbmFtZRgCIAEoCRINCgVzY29yZRgDIAEoAyITChFHZXRSYW5raW5nUmVxdWVzdCI6ChJHZXRSYW5raW5nUmVzcG9uc2USJAoHcmFua2luZxgBIAMoCzITLmNvbnRlc3RhbnQudjEuUmFuazJjCg5SYW5raW5nU2VydmljZRJRCgpHZXRSYW5raW5nEiAuY29udGVzdGFudC52MS5HZXRSYW5raW5nUmVxdWVzdBohLmNvbnRlc3RhbnQudjEuR2V0UmFua2luZ1Jlc3BvbnNlQsMBChFjb20uY29udGVzdGFudC52MUIMUmFua2luZ1Byb3RvUAFaS2dpdGh1Yi5jb20vaWN0c2MvaWN0c2MtcmVnYWxpYS9iYWNrZW5kL3BrZy9wcm90by9jb250ZXN0YW50L3YxO2NvbnRlc3RhbnR2MaICA0NYWKoCDUNvbnRlc3RhbnQuVjHKAg1Db250ZXN0YW50XFYx4gIZQ29udGVzdGFudFxWMVxHUEJNZXRhZGF0YeoCDkNvbnRlc3RhbnQ6OlYxYgZwcm90bzM");

/**
 * @generated from message contestant.v1.Rank
 */
export type Rank = Message<"contestant.v1.Rank"> & {
  /**
   * @generated from field: int64 rank = 1;
   */
  rank: bigint;

  /**
   * @generated from field: string team_name = 2;
   */
  teamName: string;

  /**
   * @generated from field: int64 score = 3;
   */
  score: bigint;
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

