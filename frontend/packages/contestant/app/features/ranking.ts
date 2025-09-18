import { type Transport, createClient } from "@connectrpc/connect";

import {
  type Rank as ProtoRank,
  RankingService,
} from "@ictsc/proto/contestant/v1";

export type Rank = ProtoRank;

export async function fetchRanking(transport: Transport): Promise<Rank[]> {
  const client = createClient(RankingService, transport);
  const ranking = await client.getRanking({});
  return ranking.ranking;
}
