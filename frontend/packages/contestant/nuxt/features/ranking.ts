import { expectData, type ApiClient } from "@ictsc/api";
import type { Rank } from "./models";

export type { Rank } from "./models";

export type RankingSnapshot = {
  ranking: Rank[];
  frozen: boolean;
  frozenAt?: string;
};

export type HeaderRankingItem =
  { kind: "rank"; entry: Rank } | { kind: "tie"; rank: number; count: number };

export function headerRankingItems(
  ranking: Rank[],
  teamCode: number | undefined,
): HeaderRankingItem[] {
  const own = ranking.find((entry) => entry.teamCode === teamCode);
  if (!own) {
    const first = ranking[0];
    return first ? [{ kind: "rank", entry: first }] : [];
  }

  const representatives = new Map<number, Rank>();
  for (const entry of ranking) {
    if (!representatives.has(entry.rank))
      representatives.set(entry.rank, entry);
  }
  representatives.set(own.rank, own);

  const lastRank = Math.max(...representatives.keys());
  const ranks =
    own.rank === 1
      ? [1, 2]
      : own.rank === lastRank
        ? [1, lastRank - 1, lastRank]
        : [1, own.rank - 1, own.rank, own.rank + 1];
  const tieCount = ranking.filter((entry) => entry.rank === own.rank).length;

  return [...new Set(ranks)].flatMap<HeaderRankingItem>((rank) => {
    const entry = representatives.get(rank);
    if (!entry) return [];
    const items: HeaderRankingItem[] = [{ kind: "rank", entry }];
    if (entry.teamCode === own.teamCode && tieCount > 1) {
      items.push({ kind: "tie", rank: own.rank, count: tieCount });
    }
    return items;
  });
}

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
