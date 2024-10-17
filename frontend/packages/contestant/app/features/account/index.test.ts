import { describe, expect, it } from "vitest";
import { createConnectTransport } from "@connectrpc/connect-web";
import { setupMSW } from "@app/__test__/msw/node";
import { connect } from "@app/__test__/msw/connect";
import { ContestantService } from "@ictsc/proto/contestant/v1";
import { fetchMe } from "./index";

const server = setupMSW();

describe("fetchMe", () => {
  it("fetches current user data while logged in", async () => {
    server.use(
      connect.rpc(ContestantService.method.getMe, () => ({
        user: {
          id: "1",
          name: "Alice",
          teamId: "1",
        },
      })),
    );
    const transport = createConnectTransport({
      baseUrl: "http://example.test",
    });
    expect(await fetchMe(transport)).toEqual({
      id: "1",
      name: "Alice",
      teamID: "1",
    });
  });
});
