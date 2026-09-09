import { describe, expect, it } from "vitest";
import { createApiClient } from "@ictsc/api";
import { setupMSW } from "../../__test__/msw/node";
import { HttpResponse, http } from "../../__test__/msw/rest";
import { fetchViewer } from "./index";

const server = setupMSW();

describe("fetchViewer", () => {
  it("maps a contestant viewer", async () => {
    const client = createApiClient("http://example.test");
    server.use(
      http.get("http://example.test/api/v1/viewer", () =>
        HttpResponse.json({
          viewer: {
            state: "CONTESTANT",
            profile: {
              name: "alice",
              display_name: "Alice",
              self_introduction: "hello",
            },
            team: {
              code: 2,
              name: "A",
              organization: "ICTSC",
              member_limit: 5,
            },
            impersonated_by: null,
          },
        }),
      ),
    );
    await expect(fetchViewer(client)).resolves.toEqual({
      type: "contestant",
      name: "alice",
      displayName: "Alice",
      admin: { canListContestants: false, canImpersonateContestants: false },
      impersonation: undefined,
    });
  });

  it("maps an anonymous viewer", async () => {
    const client = createApiClient("http://example.test");
    server.use(
      http.get("http://example.test/api/v1/viewer", () =>
        HttpResponse.json({ viewer: { state: "ANONYMOUS" } }),
      ),
    );
    await expect(fetchViewer(client)).resolves.toEqual({
      type: "unauthenticated",
      admin: { canListContestants: false, canImpersonateContestants: false },
    });
  });
});
