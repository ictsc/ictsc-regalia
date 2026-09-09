import { expectData, type ApiClient } from "@ictsc/api";

export type ViewerAdmin = {
  canListContestants: boolean;
  canImpersonateContestants: boolean;
};

export type Contestant = {
  type: "contestant";
  name: string;
  displayName: string;
  admin: ViewerAdmin;
  impersonation?: {
    adminName: string;
  };
};

export type PreSignUpUser = {
  type: "pre-signup";
  name: string;
  displayName: string;
  admin: ViewerAdmin;
};

export type User =
  Contestant | PreSignUpUser | { type: "unauthenticated"; admin: ViewerAdmin };

export async function fetchViewer(client: ApiClient): Promise<User> {
  const { viewer } = expectData(await client.GET("/api/v1/viewer"));
  switch (viewer.state) {
    case "CONTESTANT":
      return {
        type: "contestant",
        name: viewer.profile.name,
        displayName: viewer.profile.display_name,
        admin: defaultViewerAdmin(),
        impersonation:
          viewer.impersonated_by == null
            ? undefined
            : { adminName: viewer.impersonated_by },
      };
    case "DISCORD_AUTHENTICATED":
      return {
        type: "pre-signup",
        name: viewer.discord.username,
        displayName: viewer.discord.display_name,
        admin: defaultViewerAdmin(),
      };
    case "ANONYMOUS":
      return {
        type: "unauthenticated",
        admin: defaultViewerAdmin(),
      };
  }
}

function defaultViewerAdmin(): ViewerAdmin {
  return {
    canListContestants: false,
    canImpersonateContestants: false,
  };
}

export {
  listImpersonationCandidates,
  startImpersonation,
  type ImpersonationCandidate,
} from "./impersonation";
