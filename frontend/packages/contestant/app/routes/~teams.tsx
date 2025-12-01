import { createFileRoute } from "@tanstack/react-router";
import { fetchTeams } from "@app/features/teams";
import { use, useDeferredValue } from "react";
import { TeamsPage } from "./teams.page";

export const Route = createFileRoute("/teams")({
  component: RouteComponent,
  loader: ({ context: { transport } }) => {
    return {
      teams: fetchTeams(transport),
    };
  },
});

function RouteComponent() {
  const { teams: teamsPromise } = Route.useLoaderData();
  const deferredTeamsPromise = useDeferredValue(teamsPromise);
  const teams = use(deferredTeamsPromise);
  return <TeamsPage teamProfile={teams} />;
}
