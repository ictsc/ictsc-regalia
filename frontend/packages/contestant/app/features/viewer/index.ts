import {
  type Transport,
  createClient,
  ConnectError,
  Code,
} from "@connectrpc/connect";
import * as contestantv1 from "@ictsc/proto/contestant/v1";

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
  | Contestant
  | PreSignUpUser
  | { type: "unauthenticated"; admin: ViewerAdmin };

export async function fetchViewer(transport: Transport): Promise<User> {
  const client = createClient(contestantv1.ViewerService, transport);
  try {
    const { viewer } = await client.getViewer({});
    if (viewer == null) {
      throw new Error("getViewer must return viewer");
    }

    switch (viewer.viewer.value?.$typeName) {
      case "contestant.v1.ContestantViewer":
        return {
          type: "contestant",
          name: viewer.name,
          displayName: viewer.viewer.value.displayName,
          admin: parseViewerAdmin(viewer.admin),
          impersonation:
            viewer.viewer.value.impersonation == null
              ? undefined
              : {
                  adminName: viewer.viewer.value.impersonation.adminName,
                },
        };
      case "contestant.v1.SignUpViewer":
        return {
          type: "pre-signup",
          name: viewer.name,
          displayName: viewer.viewer.value.displayName,
          admin: parseViewerAdmin(viewer.admin),
        };
      case "contestant.v1.UnauthenticatedViewer":
        return {
          type: "unauthenticated",
          admin: parseViewerAdmin(viewer.admin),
        };
      default:
        throw new Error("unsupported type");
    }
  } catch (err: unknown) {
    if (err instanceof ConnectError && err.code == Code.Unauthenticated) {
      return { type: "unauthenticated", admin: defaultViewerAdmin() };
    }
    throw err;
  }
}

function parseViewerAdmin(admin?: contestantv1.ViewerAdmin): ViewerAdmin {
  if (admin == null) {
    return defaultViewerAdmin();
  }
  return {
    canListContestants: admin.canListContestants,
    canImpersonateContestants: admin.canImpersonateContestants,
  };
}

function defaultViewerAdmin(): ViewerAdmin {
  return {
    canListContestants: false,
    canImpersonateContestants: false,
  };
}

export { useSignOut, type SignOutAction } from "./use-signout";
export {
  listImpersonationCandidates,
  startImpersonation,
  type ImpersonationCandidate,
} from "./impersonation";
