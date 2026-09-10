import {
  expectData,
  impersonateContestant,
  expectNoContent,
  type AdminDeployment,
  type AdminDeploymentStreamMessage,
  type ApiClient,
  type components,
} from "@ictsc/api";

type Team = components["schemas"]["Team"];

export async function getAdminViewer(api: ApiClient) {
  return expectData(await api.GET("/api/v1/admin/viewer")).viewer;
}

export async function signOutAdmin(api: ApiClient): Promise<void> {
  expectNoContent(await api.POST("/api/v1/admin/auth/signout"));
}

export async function getAdminAnswer(
  api: ApiClient,
  reference: components["schemas"]["AdminAnswerReference"],
) {
  return expectData(
    await api.GET(
      "/api/v1/admin/answers/{team_code}/{problem_code}/{answer_number}",
      {
        params: {
          path: reference,
        },
      },
    ),
  ).answer;
}

export async function listAdminMarkingResults(api: ApiClient) {
  return expectData(await api.GET("/api/v1/admin/marking-results"))
    .marking_results;
}

export async function createAdminMarkingResult(
  api: ApiClient,
  answer: components["schemas"]["AdminAnswerReference"],
  score: number,
  rationale: string,
) {
  return expectData(
    await api.POST("/api/v1/admin/marking-results", {
      body: { answer, score, rationale },
    }),
  ).marking_result;
}

export function createAdminActions(api: ApiClient) {
  return {
    getProblem: async (problemCode: string, commit?: string) =>
      expectData(
        await api.GET("/api/v1/admin/problems/{problem_code}", {
          params: {
            path: { problem_code: problemCode },
            ...(commit == null || commit === "" ? {} : { query: { commit } }),
          },
        }),
      ),
    getAnnouncement: async (announcementSlug: string) =>
      expectData(
        await api.GET("/api/v1/admin/announcements/{announcement_slug}", {
          params: { path: { announcement_slug: announcementSlug } },
        }),
      ),
    refreshContent: async (commit: string) =>
      expectData(
        await api.POST("/api/v1/admin/content/actions/refresh", {
          body: { commit },
        }),
      ),
    createTeam: async (team: {
      code: number;
      name: string;
      organization: string;
      memberLimit: number;
    }) =>
      expectData(
        await api.POST("/api/v1/admin/teams", {
          body: {
            code: team.code,
            name: team.name,
            organization: team.organization,
            member_limit: team.memberLimit,
          },
        }),
      ),
    updateTeam: async (code: number, team: Omit<Team, "code">) =>
      expectData(
        await api.PATCH("/api/v1/admin/teams/{team_code}", {
          params: { path: { team_code: code } },
          body: team,
        }),
      ),
    deleteTeam: async (code: number) =>
      expectNoContent(
        await api.DELETE("/api/v1/admin/teams/{team_code}", {
          params: { path: { team_code: code } },
        }),
      ),
    createInvitation: async (teamCode: number, expiresAt: string) =>
      expectData(
        await api.POST("/api/v1/admin/invitations", {
          body: { team_code: teamCode, expires_at: expiresAt },
        }),
      ),
    impersonate: (contestantName: string) =>
      impersonateContestant(api, contestantName),
    recalculate: async () =>
      expectNoContent(
        await api.POST("/api/v1/admin/scores/actions/recalculate"),
      ),
    reveal: async () =>
      expectNoContent(
        await api.POST("/api/v1/admin/scores/actions/reveal-final"),
      ),
    createDeployment: async (teamCode: number, problemCode: string) =>
      expectData(
        await api.POST("/api/v1/admin/deployments", {
          body: { team_code: teamCode, problem_code: problemCode },
        }),
      ),
    syncDeployment: async (teamCode: number, problemCode: string) =>
      expectNoContent(
        await api.POST(
          "/api/v1/admin/deployments/{team_code}/{problem_code}/sync",
          {
            params: {
              path: { team_code: teamCode, problem_code: problemCode },
            },
          },
        ),
      ),
    replaceRule: async (markdown: string) =>
      expectData(await api.PUT("/api/v1/admin/rule", { body: { markdown } })),
    replaceFreeze: async (rankingFreezeAt: string | null) =>
      expectData(
        await api.PUT("/api/v1/admin/dashboard-schedule", {
          body: { ranking_freeze_at: rankingFreezeAt },
        }),
      ),
  };
}

export function mergeAdminDeployments(
  current: AdminDeployment[],
  message: AdminDeploymentStreamMessage,
): AdminDeployment[] {
  if (message.type === "snapshot") return message.deployments;
  const next = message.deployment;
  return [
    next,
    ...current.filter(
      (item) =>
        item.team_code !== next.team_code ||
        item.problem_code !== next.problem_code ||
        item.revision !== next.revision,
    ),
  ];
}
