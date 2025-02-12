import {
  type Transport,
  createClient,
  ConnectError,
  Code,
} from "@connectrpc/connect";
import * as contestantv1 from "@ictsc/proto/contestant/v1";

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
        };
      case "contestant.v1.SignUpViewer":
        return {
          type: "pre-signup",
          name: viewer.name,
          displayName: viewer.viewer.value.displayName,
        };
      default:
        throw new Error("unsupported type");
    }
  } catch (err: unknown) {
    if (err instanceof ConnectError && err.code == Code.Unauthenticated) {
      return { type: "unauthenticated" };
    }
    throw err;
  }
}
