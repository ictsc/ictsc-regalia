import { expectData, type ApiClient } from "@ictsc/api";
import type { TeamProfile } from "./models";

export type { TeamProfile } from "./models";

export async function fetchTeams(client: ApiClient): Promise<TeamProfile[]> {
  const response = expectData(await client.GET("/api/v1/contestant/teams"));
  return response.teams.map(({ team, members }) => ({
    code: Number(team.code),
    name: team.name,
    organization: team.organization,
    memberLimit: team.member_limit,
    team: {
      code: Number(team.code),
      name: team.name,
      organization: team.organization,
      memberLimit: team.member_limit,
    },
    members: members.map((member) => ({
      name: member.name,
      displayName: member.display_name,
      selfIntroduction: member.self_introduction,
    })),
  }));
}
