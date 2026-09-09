import {
  afterAll,
  afterEach,
  beforeAll,
  beforeEach,
  describe,
  expect,
  it,
} from "vitest";
import {
  ApiError,
  createApiClient,
  type AdminDeployment,
  type ApiClient,
} from "@ictsc/api";
import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import {
  createAdminActions,
  createAdminMarkingResult,
  getAdminAnswer,
  getAdminViewer,
  listAdminMarkingResults,
  mergeAdminDeployments,
  signOutAdmin,
} from "./admin-api";

const origin = "http://example.test";
const commit = "1111111111111111111111111111111111111111";
const now = "2026-08-30T00:00:00Z";
const team = {
  code: 2,
  name: "Team A",
  organization: "ICTSC",
  member_limit: 4,
};
const content = {
  state: "SUCCEEDED" as const,
  active_commit: commit,
  activated_at: now,
  source_ref: "refs/heads/main",
  last_attempt_commit: commit,
  last_attempt_at: now,
  last_error: null,
  serving_last_known_good: false,
};
const deployment: AdminDeployment = {
  team_code: team.code,
  problem_code: "A01",
  revision: 1,
  latest_status: "QUEUED",
  events: [
    {
      event_id: "11111111-1111-4111-8111-111111111111",
      occurred_at: now,
      status: "QUEUED",
      message: null,
    },
  ],
  content_commit: commit,
};

const server = setupServer();
let client: ApiClient;

beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
afterAll(() => server.close());
afterEach(() => server.resetHandlers());
beforeEach(() => {
  // openapi-fetch captures fetch when the client is created, so create it
  // after MSW has installed its Node interceptor.
  client = createApiClient(origin);
});

function problem(status: number, code: string, retryAfter?: string) {
  return HttpResponse.json(
    {
      type: "about:blank",
      title: "Request failed",
      status,
      detail: code,
      code,
      errors: [],
    },
    {
      status,
      headers: retryAfter == null ? undefined : { "Retry-After": retryAfter },
    },
  );
}

describe("Admin REST loader and auth", () => {
  it("loads an Admin viewer and performs idempotent signout", async () => {
    server.use(
      http.get(`${origin}/api/v1/admin/viewer`, () =>
        HttpResponse.json({
          viewer: {
            state: "ADMIN",
            admin: {
              discord: {
                id: "123",
                username: "operator",
                display_name: "Operator",
              },
              guild_id: "guild",
              role_ids: ["admin-role"],
            },
          },
        }),
      ),
      http.post(
        `${origin}/api/v1/admin/auth/signout`,
        () => new HttpResponse(null, { status: 204 }),
      ),
    );

    await expect(getAdminViewer(client)).resolves.toMatchObject({
      state: "ADMIN",
    });
    await expect(signOutAdmin(client)).resolves.toBeUndefined();
  });

  it.each([401, 403, 404, 409, 422, 429, 502])(
    "surfaces RFC 9457 status %i",
    async (status: number) => {
      server.use(
        http.get(`${origin}/api/v1/admin/viewer`, () =>
          problem(
            status,
            `status_${status}`,
            status === 429 ? "17" : undefined,
          ),
        ),
      );

      const error = await getAdminViewer(client).catch(
        (caught: unknown) => caught,
      );
      expect(error).toBeInstanceOf(ApiError);
      expect(error).toMatchObject({ status, code: `status_${status}` });
      if (status === 429) {
        expect((error as ApiError).retryAfterSeconds).toBe(17);
      }
    },
  );
});

describe("Admin mutation actions", () => {
  it("sends contract-shaped bodies for content, teams, invitations, scores and deployments", async () => {
    const seen = new Map<string, unknown>();
    const remember = async (request: Request) => {
      const body: unknown = await request.json();
      seen.set(`${request.method} ${new URL(request.url).pathname}`, body);
    };
    server.use(
      http.get(
        `${origin}/api/v1/admin/problems/:problemCode`,
        ({ params, request }) => {
          expect(params.problemCode).toBe("A01");
          expect(new URL(request.url).searchParams.get("commit")).toBe(commit);
          return HttpResponse.json({
            problem: { code: "A01", content_commit: commit },
          });
        },
      ),
      http.get(`${origin}/api/v1/admin/announcements/:slug`, ({ params }) =>
        HttpResponse.json({
          announcement: { slug: params.slug, markdown: "notice" },
        }),
      ),
      http.post(
        `${origin}/api/v1/admin/content/actions/refresh`,
        async ({ request }) => {
          await remember(request);
          return HttpResponse.json({ content });
        },
      ),
      http.post(`${origin}/api/v1/admin/teams`, async ({ request }) => {
        await remember(request);
        return HttpResponse.json({ team }, { status: 201 });
      }),
      http.patch(
        `${origin}/api/v1/admin/teams/:teamCode`,
        async ({ request }) => {
          await remember(request);
          return HttpResponse.json({ team });
        },
      ),
      http.delete(
        `${origin}/api/v1/admin/teams/:teamCode`,
        () => new HttpResponse(null, { status: 204 }),
      ),
      http.post(`${origin}/api/v1/admin/invitations`, async ({ request }) => {
        await remember(request);
        return HttpResponse.json(
          {
            invitation: {
              code: "invite-a",
              team_code: team.code,
              created_at: now,
              expires_at: "2026-08-31T00:00:00Z",
            },
          },
          { status: 201 },
        );
      }),
      http.post(
        `${origin}/api/v1/admin/impersonations`,
        () => new HttpResponse(null, { status: 204 }),
      ),
      http.post(
        `${origin}/api/v1/admin/scores/actions/recalculate`,
        () => new HttpResponse(null, { status: 204 }),
      ),
      http.post(
        `${origin}/api/v1/admin/scores/actions/reveal-final`,
        () => new HttpResponse(null, { status: 204 }),
      ),
      http.post(`${origin}/api/v1/admin/deployments`, async ({ request }) => {
        await remember(request);
        return HttpResponse.json({ deployment }, { status: 201 });
      }),
      http.post(
        `${origin}/api/v1/admin/deployments/:teamCode/:problemCode/sync`,
        () => new HttpResponse(null, { status: 204 }),
      ),
      http.put(`${origin}/api/v1/admin/rule`, async ({ request }) => {
        await remember(request);
        return HttpResponse.json({ rule: { markdown: "# Updated" } });
      }),
      http.put(
        `${origin}/api/v1/admin/dashboard-schedule`,
        async ({ request }) => {
          await remember(request);
          return HttpResponse.json({
            dashboard_schedule: { ranking_freeze_at: null },
          });
        },
      ),
    );
    const actions = createAdminActions(client);

    await actions.getProblem("A01", commit);
    await actions.getAnnouncement("notice");
    await actions.refreshContent(commit);
    await actions.createTeam({
      code: 2,
      name: team.name,
      organization: team.organization,
      memberLimit: team.member_limit,
    });
    await actions.updateTeam(2, {
      name: "Updated",
      organization: team.organization,
      member_limit: 5,
    });
    await actions.deleteTeam(2);
    await actions.createInvitation(2, "2026-08-31T00:00:00Z");
    await actions.impersonate("alice");
    await actions.recalculate();
    await actions.reveal();
    await actions.createDeployment(2, "A01");
    await actions.syncDeployment(2, "A01");
    await actions.replaceRule("# Updated");
    await actions.replaceFreeze(null);

    expect(seen.get("POST /api/v1/admin/content/actions/refresh")).toEqual({
      commit,
    });
    expect(seen.get("POST /api/v1/admin/teams")).toEqual({
      code: 2,
      name: team.name,
      organization: team.organization,
      member_limit: 4,
    });
    expect(seen.get("POST /api/v1/admin/deployments")).toEqual({
      team_code: 2,
      problem_code: "A01",
    });
    expect(seen.get("PUT /api/v1/admin/dashboard-schedule")).toEqual({
      ranking_freeze_at: null,
    });
  });
});

describe("Admin marking and deployment SSE state", () => {
  const reference = {
    team_code: 2,
    problem_code: "A01",
    answer_number: 1,
  };

  it("preserves nullable score/content revision and posts a regrade", async () => {
    let markingBody: unknown;
    server.use(
      http.get(
        `${origin}/api/v1/admin/answers/:teamCode/:problemCode/:answerNumber`,
        () =>
          HttpResponse.json({
            answer: {
              reference,
              team,
              author: {
                name: "alice",
                display_name: "Alice",
                self_introduction: "",
              },
              problem: {
                code: "A01",
                title: "Problem",
                max_score: 100,
                category: "Network",
              },
              body: { type: "DESCRIPTIVE", body: "answer" },
              submitted_at: now,
              score: null,
              content_commit: commit,
            },
          }),
      ),
      http.get(`${origin}/api/v1/admin/marking-results`, () =>
        HttpResponse.json({ marking_results: [] }),
      ),
      http.post(
        `${origin}/api/v1/admin/marking-results`,
        async ({ request }) => {
          markingBody = await request.json();
          return HttpResponse.json(
            {
              marking_result: {
                id: "22222222-2222-4222-8222-222222222222",
                answer: reference,
                judge: { name: "operator" },
                score: 80,
                rationale: "ok",
                created_at: now,
                visibility: "PRIVATE",
              },
            },
            { status: 201 },
          );
        },
      ),
    );

    const answer = await getAdminAnswer(client, reference);
    expect(answer).toMatchObject({ score: null, content_commit: commit });
    await expect(listAdminMarkingResults(client)).resolves.toEqual([]);
    await createAdminMarkingResult(client, reference, 80, "ok");
    expect(markingBody).toEqual({
      answer: reference,
      score: 80,
      rationale: "ok",
    });
  });

  it("replaces reconnect snapshots and de-duplicates deployment events", () => {
    const completed: AdminDeployment = {
      ...deployment,
      latest_status: "COMPLETED",
    };
    expect(
      mergeAdminDeployments([deployment], {
        type: "snapshot",
        deployments: [completed],
      }),
    ).toEqual([completed]);
    expect(
      mergeAdminDeployments([deployment], {
        type: "deployment",
        deployment: completed,
      }),
    ).toEqual([completed]);
  });
});
