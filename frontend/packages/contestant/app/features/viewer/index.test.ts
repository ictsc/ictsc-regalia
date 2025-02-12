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
    });
  });
});
