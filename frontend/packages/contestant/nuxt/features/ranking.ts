import { expectData, type ApiClient } from "@ictsc/api";
import type { Rank } from "./models";

export type { Rank } from "./models";

export type RankingSnapshot = {
  ranking: Rank[];
  frozen: boolean;
  frozenAt?: string;
};

export async function fetchRanking(
  client: ApiClient,
): Promise<RankingSnapshot> {
  const response = expectData(await client.GET("/api/v1/contestant/ranking"));
  return {
    ranking: response.ranking.map((rank) => ({
      rank: rank.rank,
      teamCode: Number(rank.team_code),
      teamName: rank.team_name,
      organization: rank.organization,
      score: Number(rank.score),
      lastEffectiveSubmissionAt: rank.last_effective_submission_at ?? undefined,
    })),
    frozen: response.frozen,
    frozenAt: response.frozen_at ?? undefined,
  };
}
