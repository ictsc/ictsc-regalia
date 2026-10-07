import { expectData, type ApiClient } from "@ictsc/api";

export type Contestant = {
  type: "contestant";
  name: string;
  displayName: string;
};

export type PreSignUpUser = {
  type: "pre-signup";
  name: string;
  displayName: string;
};

export type User = Contestant | PreSignUpUser | { type: "unauthenticated" };

export async function fetchViewer(client: ApiClient): Promise<User> {
  const { viewer } = expectData(await client.GET("/api/v1/viewer"));
  switch (viewer.state) {
    case "CONTESTANT":
      return {
        type: "contestant",
        name: viewer.profile.name,
        displayName: viewer.profile.display_name,
      };
    case "DISCORD_AUTHENTICATED":
      return {
        type: "pre-signup",
        name: viewer.discord.username,
        displayName: viewer.discord.display_name,
      };
    case "ANONYMOUS":
      return { type: "unauthenticated" };
  }
}
