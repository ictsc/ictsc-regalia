import { type Transport, createPromiseClient } from "@connectrpc/connect";
import { ContestantService } from "@ictsc/proto/contestant/v1";

export type User = {
  id: string;
  name: string;
  teamID: string;
};

export async function fetchMe(transport: Transport): Promise<User | undefined> {
  const client = createPromiseClient(ContestantService, transport);
  const { user } = await client.getMe({});
  if (user == null) {
    return;
  }
  return {
    id: user.id,
    name: user.name,
    teamID: user.teamId,
  };
}
