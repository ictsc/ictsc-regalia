import { describe, expect, it } from "vitest";
import { setupMSW } from "@app/__test__/msw/node";
import { connect } from "@app/__test__/msw/connect";
import { ViewerService } from "@ictsc/proto/contestant/v1";
import { listImpersonationCandidates } from "./impersonation";

const server = setupMSW();

describe("listImpersonationCandidates", () => {
  it("fetches contestants from viewer service", async () => {
    server.use(
      connect.rpc(ViewerService.method.listContestants, () => ({
        contestants: [
          {
            name: "alice",
            displayName: "Alice",
            teamName: "Team A",
            teamCode: 1n,
          },
        ],
      })),
    );

    await expect(listImpersonationCandidates()).resolves.toEqual([
      {
        name: "alice",
        displayName: "Alice",
        teamName: "Team A",
        teamCode: 1,
      },
    ]);
  });
});
