import { createFileRoute } from "@tanstack/react-router";
import { fetchRanking } from "@app/features/ranking";
import { use, useDeferredValue } from "react";
import { RankingPage } from "./page";

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
  return <RankingPage ranking={ranking} />;
}
