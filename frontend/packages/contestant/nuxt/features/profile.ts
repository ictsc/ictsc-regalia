import { expectData, type ApiClient } from "@ictsc/api";
import type { ContestantProfile } from "./models";

export async function fetchProfile(
  client: ApiClient,
): Promise<ContestantProfile> {
  const { profile } = expectData(
    await client.GET("/api/v1/contestant/profile"),
  );
  return {
    name: profile.name,
    displayName: profile.display_name,
    selfIntroduction: profile.self_introduction,
  };
}

export async function updateProfile(
  client: ApiClient,
  input: Pick<ContestantProfile, "displayName" | "selfIntroduction">,
): Promise<ContestantProfile> {
  const { profile } = expectData(
    await client.PATCH("/api/v1/contestant/profile", {
      body: {
        display_name: input.displayName,
        self_introduction: input.selfIntroduction,
      },
    }),
  );
  return {
    name: profile.name,
    displayName: profile.display_name,
    selfIntroduction: profile.self_introduction,
  };
}
