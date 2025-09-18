import { createClient, type Transport } from "@connectrpc/connect";
import {
  ProfileService,
  type TeamProfile as ProtoTeamProfile,
} from "@ictsc/proto/contestant/v1";

export type TeamProfile = ProtoTeamProfile;

export async function fetchTeams(transport: Transport): Promise<TeamProfile[]> {
  const client = createClient(ProfileService, transport);
  const teams = await client.listTeams({});
  return teams.teams;
}
