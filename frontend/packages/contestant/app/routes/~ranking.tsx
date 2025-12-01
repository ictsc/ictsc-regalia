import { createFileRoute } from "@tanstack/react-router";
import { fetchRanking } from "@app/features/ranking";
import { use, useDeferredValue } from "react";
import { RankingPage } from "./ranking.page";
import { timestampDate } from "@bufbuild/protobuf/wkt";

export const Route = createFileRoute("/ranking")({
  component: RouteComponent,
  loader: ({ context: { transport } }) => {
    return {
      ranking: fetchRanking(transport),
    };
  },
});

function RouteComponent() {
  const { ranking: rankingPromise } = Route.useLoaderData();
  const deferredRankingPromise = useDeferredValue(rankingPromise);
  const ranking = use(deferredRankingPromise);
  return (
    <RankingPage
      ranking={ranking.map((rank) => ({
        rank: Number(rank.rank),
        teamName: rank.teamName,
        organization: rank.organization,
        score: Number(rank.score),
        lastEffectiveSubmitAt:
          rank.timestamp != null
            ? timestampDate(rank.timestamp).toISOString()
            : undefined,
      }))}
    />
  );
}
