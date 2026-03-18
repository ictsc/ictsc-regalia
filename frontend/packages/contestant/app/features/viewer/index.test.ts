import { describe, expect, it } from "vitest";
import { createConnectTransport } from "@connectrpc/connect-web";
import { setupMSW } from "@app/__test__/msw/node";
import { connect } from "@app/__test__/msw/connect";
import { ViewerService } from "@ictsc/proto/contestant/v1";
import { fetchViewer } from "./index";

const server = setupMSW();

describe("fetchMe", () => {
  it("fetches current user data while logged in", async () => {
    server.use(
      connect.rpc(ViewerService.method.getViewer, () => ({
        viewer: {
          name: "alice",
          admin: {
            canListContestants: false,
            canImpersonateContestants: false,
          },
          viewer: {
            case: "contestant",
            value: {
              displayName: "Alice",
            },
          },
        },
      })),
    );
    const transport = createConnectTransport({
      baseUrl: "http://example.test",
    });
    expect(await fetchViewer(transport)).toEqual({
      type: "contestant",
      name: "alice",
      displayName: "Alice",
      admin: {
        canListContestants: false,
        canImpersonateContestants: false,
      },
      impersonation: undefined,
    });
  });

  it("returns impersonation capability while signed out", async () => {
    server.use(
      connect.rpc(ViewerService.method.getViewer, () => ({
        viewer: {
          admin: {
            canListContestants: true,
            canImpersonateContestants: true,
          },
          viewer: {
            case: "unauthenticated",
            value: {},
          },
        },
      })),
    );
    const transport = createConnectTransport({
      baseUrl: "http://example.test",
    });
    expect(await fetchViewer(transport)).toEqual({
      type: "unauthenticated",
      admin: {
        canListContestants: true,
        canImpersonateContestants: true,
      },
    });
  });
});
