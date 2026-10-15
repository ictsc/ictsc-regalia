import { expectData, type ApiClient, type components } from "@ictsc/api";

export type Team = components["schemas"]["Team"];

export async function fetchTeams(client: ApiClient): Promise<Team[]> {
  return expectData(await client.GET("/api/v1/admin/teams")).teams;
}
