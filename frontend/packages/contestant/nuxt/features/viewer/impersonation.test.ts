import { describe, expect, it } from "vitest";
import { createApiClient } from "@ictsc/api";
import { setupMSW } from "../../__test__/msw/node";
import { HttpResponse, http } from "../../__test__/msw/rest";
import { listImpersonationCandidates } from "./impersonation";

const server = setupMSW();

describe("listImpersonationCandidates", () => {
  it("maps Admin contestants from REST", async () => {
    server.use(
      http.get("http://example.test/api/v1/admin/contestants", () =>
        HttpResponse.json({
          contestants: [
            {
              profile: {
                name: "alice",
                display_name: "Alice",
                self_introduction: "",
              },
              team: {
                code: 2,
                name: "Team A",
                organization: "ICTSC",
                member_limit: 5,
              },
              discord_id: "123",
            },
          ],
        }),
      ),
    );
    const client = createApiClient("http://example.test");
    await expect(listImpersonationCandidates(client)).resolves.toEqual([
      { name: "alice", displayName: "Alice", teamName: "Team A", teamCode: 2 },
    ]);
  });
});
