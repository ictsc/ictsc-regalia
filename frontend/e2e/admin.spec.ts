import { expect, test, type Page, type Request } from "@playwright/test";
import { emitSse, fulfillJson, installEventSourceMock } from "./support/http";

const activeCommit = "a".repeat(40);
const submittedCommit = "b".repeat(40);
const historicalCommit = "c".repeat(40);
const now = "2026-09-01T12:00:00+09:00";

const team = {
  code: 12,
  name: "Existing Team",
  organization: "ICTSC University",
  member_limit: 4,
};
const problem = {
  code: "A01",
  title: "Reliable DNS",
  max_score: 100,
  category: "Network",
  section_slug: "day1",
  type: "DESCRIPTIVE",
  default_grader_discord_id: "123456789012345678",
  body: "# Restore service",
  explanation: "# 解説\n\nDNSを復旧します。",
  redeploy_rule: {
    type: "PERCENTAGE_PENALTY",
    penalty_threshold: 1,
    penalty_percentage: 10,
  },
  content_commit: activeCommit,
};
const answer = {
  reference: { team_code: 12, problem_code: "A01", answer_number: 1 },
  team,
  author: {
    name: "alice",
    display_name: "Alice",
    self_introduction: "operator",
  },
  problem: {
    code: "A01",
    title: "Reliable DNS",
    max_score: 100,
    category: "Network",
  },
  body: { type: "DESCRIPTIVE", body: "# 回答\n\n設定を修正しました。" },
  submitted_at: now,
  score: { total: 72, marked: 82, penalty: 10, max: 100 },
  content_commit: submittedCommit,
  workflow: {
    status: "COMPLETED",
    assignee_discord_id: "123456789012345678",
    assignment_source: "DEFAULT",
    revision: 1,
  },
};
const existingMark = {
  id: "11111111-1111-4111-8111-111111111111",
  answer: answer.reference,
  judge: { name: "judge-a" },
  score: 82,
  rationale: "初回採点",
  line_comments: [],
  created_at: now,
  visibility: "PRIVATE",
};
const queuedDeployment = {
  team_code: 12,
  problem_code: "A01",
  revision: 1,
  latest_status: "QUEUED",
  content_commit: activeCommit,
  events: [
    {
      event_id: "22222222-2222-4222-8222-222222222222",
      occurred_at: now,
      status: "QUEUED",
      message: null,
    },
  ],
};
const completedDeployment = {
  ...queuedDeployment,
  latest_status: "COMPLETED",
  events: [
    ...queuedDeployment.events,
    {
      event_id: "33333333-3333-4333-8333-333333333333",
      occurred_at: "2026-09-01T12:00:05+09:00",
      status: "COMPLETED",
      message: "ready",
    },
  ],
};

function requestBody(request: Request): unknown {
  const raw = request.postData();
  return raw == null ? null : JSON.parse(raw);
}

type ApiState = {
  answers: Array<Record<string, unknown>>;
  teams: Array<typeof team>;
  invitations: Array<Record<string, unknown>>;
  marks: Array<Record<string, unknown>>;
  deployments: Array<typeof queuedDeployment>;
  requests: Array<{
    method: string;
    path: string;
    body: unknown;
    search: string;
  }>;
  deploymentGets: number;
};

async function installAdminApi(page: Page): Promise<ApiState> {
  const state: ApiState = {
    answers: [structuredClone(answer)],
    teams: [team],
    invitations: [],
    marks: [existingMark],
    deployments: [],
    requests: [],
    deploymentGets: 0,
  };
  await page.route("**/api/v1/**", async (route) => {
    const request = route.request();
    const url = new URL(request.url());
    const path = url.pathname;
    const method = request.method();
    const body = requestBody(request);
    state.requests.push({ method, path, body, search: url.search });

    if (path === "/api/v1/admin/viewer") {
      return fulfillJson(route, {
        viewer: {
          state: "ADMIN",
          admin: {
            discord: {
              id: "123456789012345678",
              username: "ops-admin",
              display_name: "Ops Admin",
            },
            guild_id: "234567890123456789",
            role_ids: ["345678901234567890"],
          },
        },
      });
    }
    if (path === "/api/v1/admin/teams" && method === "GET") {
      return fulfillJson(route, { teams: state.teams });
    }
    if (path === "/api/v1/admin/teams" && method === "POST") {
      const input = body as typeof team;
      state.teams.push(input);
      return fulfillJson(route, { team: input }, 201);
    }
    if (path === "/api/v1/admin/teams/12" && method === "GET") {
      return fulfillJson(route, {
        team: state.teams.find((t) => t.code === 12),
      });
    }
    if (path === "/api/v1/admin/teams/12" && method === "PATCH") {
      const index = state.teams.findIndex((t) => t.code === 12);
      state.teams[index] = { ...state.teams[index]!, ...(body as object) };
      return fulfillJson(route, { team: state.teams[index] });
    }
    if (path === "/api/v1/admin/teams/12" && method === "DELETE") {
      return route.fulfill({ status: 204, body: "" });
    }
    if (path === "/api/v1/admin/invitations" && method === "GET") {
      return fulfillJson(route, { invitations: state.invitations });
    }
    if (path === "/api/v1/admin/invitations" && method === "POST") {
      const invitation = {
        code: "invite-deterministic",
        team_code: 12,
        created_at: now,
        expires_at: (body as { expires_at: string }).expires_at,
      };
      state.invitations.push(invitation);
      return fulfillJson(route, { invitation }, 201);
    }
    if (path === "/api/v1/admin/contestants") {
      return fulfillJson(route, {
        contestants: [
          {
            profile: answer.author,
            team,
            discord_id: "456789012345678901",
          },
        ],
      });
    }
    if (path === "/api/v1/admin/impersonations" && method === "POST") {
      return route.fulfill({ status: 204, body: "" });
    }
    if (path === "/api/v1/admin/content/status") {
      return fulfillJson(route, {
        content: {
          state: "SUCCEEDED",
          active_commit: activeCommit,
          activated_at: now,
          source_ref: "refs/heads/main",
          last_attempt_commit: activeCommit,
          last_attempt_at: now,
          last_error: null,
          serving_last_known_good: false,
        },
      });
    }
    if (path === "/api/v1/admin/content/actions/refresh" && method === "POST") {
      return fulfillJson(route, {
        content: {
          state: "SUCCEEDED",
          active_commit: (body as { commit: string }).commit,
          activated_at: now,
          source_ref: "refs/heads/main",
          last_attempt_commit: (body as { commit: string }).commit,
          last_attempt_at: now,
          last_error: null,
          serving_last_known_good: false,
        },
      });
    }
    if (path === "/api/v1/admin/sections") {
      return fulfillJson(route, {
        sections: [
          { slug: "day1", beginning: now, ending: "2026-09-01T18:00:00+09:00" },
        ],
      });
    }
    if (path === "/api/v1/admin/problems") {
      return fulfillJson(route, { problems: [problem] });
    }
    if (path === "/api/v1/admin/problems/A01") {
      const commit = url.searchParams.get("commit");
      return fulfillJson(route, {
        problem: {
          ...problem,
          max_score: commit === submittedCommit ? 90 : 100,
          content_commit: commit ?? activeCommit,
        },
      });
    }
    if (path === "/api/v1/admin/announcements") {
      return fulfillJson(route, {
        announcements: [
          {
            slug: "notice-1",
            title: "競技開始",
            effective_from: now,
            body: "# 開始",
            content_commit: activeCommit,
          },
        ],
      });
    }
    if (path === "/api/v1/admin/announcements/notice-1") {
      return fulfillJson(route, {
        announcement: {
          slug: "notice-1",
          title: "競技開始",
          effective_from: now,
          body: "# 開始",
          content_commit: activeCommit,
        },
      });
    }
    if (path === "/api/v1/admin/answers" && method === "GET") {
      const teamCode = url.searchParams.get("team_code");
      const problemCode = url.searchParams.get("problem_code");
      return fulfillJson(route, {
        answers: state.answers.filter((item) => {
          const ref = item.reference as typeof answer.reference;
          return (
            (!teamCode || ref.team_code === Number(teamCode)) &&
            (!problemCode || ref.problem_code === problemCode)
          );
        }),
      });
    }
    if (path === "/api/v1/admin/answers/12/A01/1") {
      return fulfillJson(route, { answer: state.answers[0] });
    }
    if (path === "/api/v1/admin/answers/12/A01/2") {
      return fulfillJson(route, { answer: state.answers[1] });
    }
    if (path === "/api/v1/admin/answers/12/A01/3") {
      return fulfillJson(route, { answer: state.answers[2] });
    }
    if (path.endsWith("/workflow") && method === "PATCH") {
      const input = body as {
        expected_revision: number;
        status?: string;
        assignment?: string;
      };
      const index = Number(path.split("/")[7]) - 1;
      const current = state.answers[index]!;
      const workflow = current.workflow as typeof answer.workflow;
      if (workflow.revision !== input.expected_revision)
        return fulfillJson(
          route,
          { title: "Conflict", status: 409, code: "workflow_conflict" },
          409,
        );
      current.workflow = {
        ...workflow,
        revision: workflow.revision + 1,
        status: input.status ?? workflow.status,
        assignee_discord_id:
          input.assignment === "CLAIM_SELF"
            ? "123456789012345678"
            : input.assignment === "RESET_TO_DEFAULT"
              ? "123456789012345678"
              : workflow.assignee_discord_id,
        assignment_source:
          input.assignment === "CLAIM_SELF"
            ? "CLAIMED"
            : input.assignment === "RESET_TO_DEFAULT"
              ? "DEFAULT"
              : workflow.assignment_source,
      };
      return fulfillJson(route, { answer: current });
    }
    if (path.startsWith("/api/v1/admin/answers/") && path.endsWith("/999")) {
      return fulfillJson(
        route,
        {
          type: "about:blank",
          title: "Answer not found",
          status: 404,
          detail: "指定された回答は存在しません",
          code: "answer_not_found",
        },
        404,
      );
    }
    if (path === "/api/v1/admin/marking-results" && method === "GET") {
      return fulfillJson(route, { marking_results: state.marks });
    }
    const deletedComment = path.match(
      /^\/api\/v1\/admin\/marking-results\/([^/]+)\/line-comments\/(\d+)$/,
    );
    if (deletedComment && method === "DELETE") {
      const mark = state.marks.find((item) => item.id === deletedComment[1]);
      const comments = mark?.line_comments as
        | Array<{ body: string; deleted_at?: string; deleted_by?: string }>
        | undefined;
      const note = comments?.[Number(deletedComment[2])];
      if (!note) return fulfillJson(route, { status: 404 }, 404);
      if (note.deleted_at) return fulfillJson(route, { status: 409 }, 409);
      note.deleted_at = "2026-09-01T12:30:00+09:00";
      note.deleted_by = "ops-admin";
      return fulfillJson(route, { marking_result: mark });
    }
    if (path === "/api/v1/admin/marking-results" && method === "POST") {
      const input = body as {
        answer: typeof answer.reference;
        score: number;
        rationale: string;
        line_comments?: Array<{
          line_number: number;
          end_line_number?: number;
          body: string;
        }>;
      };
      const mark = {
        id: "44444444-4444-4444-8444-444444444444",
        answer: input.answer,
        judge: { name: "ops-admin" },
        score: input.score,
        rationale: input.rationale,
        line_comments: input.line_comments ?? [],
        created_at: "2026-09-01T12:10:00+09:00",
        visibility: "PRIVATE",
      };
      state.marks.push(mark);
      const current = state.answers[input.answer.answer_number - 1];
      if (current) {
        const workflow = current.workflow as typeof answer.workflow;
        current.workflow = {
          ...workflow,
          status: "COMPLETED",
          revision: workflow.revision + 1,
        };
        current.score = {
          total: input.score,
          marked: input.score,
          penalty: 0,
          max: 100,
        };
      }
      return fulfillJson(route, { marking_result: mark }, 201);
    }
    if (path === "/api/v1/admin/scores" && method === "GET") {
      return fulfillJson(route, {
        scores: [
          {
            team,
            problem: answer.problem,
            marked_score: 82,
            penalty: 10,
            score: 72,
          },
        ],
      });
    }
    if (path === "/api/v1/admin/ranking") {
      return fulfillJson(route, {
        ranking: [
          {
            rank: 1,
            team_code: 12,
            team_name: team.name,
            organization: team.organization,
            score: 72,
            last_effective_submission_at: now,
          },
        ],
        frozen: true,
        frozen_at: now,
      });
    }
    if (
      path === "/api/v1/admin/scores/actions/recalculate" ||
      path === "/api/v1/admin/scores/actions/reveal-final"
    ) {
      return route.fulfill({ status: 204, body: "" });
    }
    if (path === "/api/v1/admin/deployments" && method === "GET") {
      state.deploymentGets += 1;
      return fulfillJson(route, { deployments: state.deployments });
    }
    if (path === "/api/v1/admin/deployments" && method === "POST") {
      state.deployments.splice(0, state.deployments.length, queuedDeployment);
      return fulfillJson(route, { deployment: queuedDeployment }, 201);
    }
    if (path === "/api/v1/admin/deployments/12/A01/sync" && method === "POST") {
      return route.fulfill({ status: 204, body: "" });
    }
    if (path === "/api/v1/admin/rule" && method === "GET") {
      return fulfillJson(route, { rule: { markdown: "# E2E rule" } });
    }
    if (path === "/api/v1/admin/rule" && method === "PUT") {
      return fulfillJson(route, { rule: body });
    }
    if (path === "/api/v1/admin/dashboard-schedule" && method === "GET") {
      return fulfillJson(route, {
        dashboard_schedule: { ranking_freeze_at: now },
      });
    }
    if (path === "/api/v1/admin/dashboard-schedule" && method === "PUT") {
      return fulfillJson(route, { dashboard_schedule: body });
    }

    return fulfillJson(
      route,
      {
        type: "about:blank",
        title: "Unexpected E2E request",
        status: 404,
        detail: method + " " + path,
        code: "not_found",
      },
      404,
    );
  });
  return state;
}

test("Admin navigation、team・招待・参加者・content履歴を操作する", async ({
  page,
}) => {
  const state = await installAdminApi(page);
  await page.goto("/admin/");
  await expect(page.getByRole("heading", { name: "概要" })).toBeVisible();
  await expect(page.getByText("Ops Admin", { exact: true })).toBeVisible();

  await page.getByRole("link", { name: "チーム", exact: true }).click();
  await page.getByLabel("チームコード").fill("13");
  await page.getByLabel("チーム名").fill("New E2E Team");
  await page.getByLabel("所属").fill("E2E College");
  await page.getByLabel("定員").fill("5");
  await page.getByRole("button", { name: "チームを作成" }).click();
  await expect(page.getByText("New E2E Team", { exact: true })).toBeVisible();

  await page.goto("/admin/teams/12");
  await page.getByLabel("有効期限").fill("2026-09-02T12:00");
  await page.getByRole("button", { name: "招待を作成" }).click();
  await expect(page.getByText("invite-deterministic")).toBeVisible();

  await page.getByRole("link", { name: "コンテンツ", exact: true }).click();
  await page.getByLabel("Git commit SHA").fill(historicalCommit);
  await page.getByRole("button", { name: "指定commitを取得" }).click();
  await expect(page.getByText("コンテンツ更新しました")).toBeVisible();
  await page.goto("/admin/content/problems/A01");
  await page.getByLabel("参照commit（空欄でactive）").fill(historicalCommit);
  await page.getByRole("button", { name: "表示" }).click();
  await expect(page.getByText(historicalCommit, { exact: true })).toBeVisible();
  expect(
    state.requests.some(
      (item) =>
        item.path === "/api/v1/admin/problems/A01" &&
        item.search.includes("commit=" + historicalCommit),
    ),
  ).toBe(true);

  await page.getByRole("link", { name: "参加者", exact: true }).click();
  await page.getByLabel("検索").fill("Alice");
  await page.getByRole("button", { name: "代理ログイン" }).click();
  const impersonationConfirm = page.getByRole("dialog", {
    name: "代理ログイン",
  });
  await expect(impersonationConfirm).toBeVisible();
  await expect(impersonationConfirm).toHaveClass(/confirmation-dialog/);
  await impersonationConfirm
    .getByRole("button", { name: "キャンセル" })
    .click();
  await expect(impersonationConfirm).not.toBeVisible();
  expect(
    state.requests.some((item) => item.path === "/api/v1/admin/impersonations"),
  ).toBe(false);
});

test("回答詳細・再採点・得点操作・freeze・404を検証する", async ({ page }) => {
  test.setTimeout(90_000);
  await page.setViewportSize({ width: 1800, height: 900 });
  const state = await installAdminApi(page);
  await page.goto("/admin/submissions/");
  await expect(
    page.getByRole("heading", { name: "採点", exact: true }),
  ).toBeVisible();
  await page.getByRole("link", { name: /A01: Reliable DNS/ }).click();

  await expect(page.locator(".grading-pane")).toHaveCount(3);
  await expect(page.getByRole("heading", { name: "過去の回答" })).toHaveCount(
    0,
  );
  expect(
    (await page.locator(".grading-panes").boundingBox())!.width,
  ).toBeGreaterThan(1400);

  await expect(page.getByText(/Existing Team · Alice/)).toBeVisible();
  await expect(page.getByText("提出時点の版を表示中")).toBeVisible();
  await expect(page.getByText("初回採点")).toBeVisible();
  await page.getByRole("button", { name: "行 1 を選択" }).click();
  await expect(page.getByLabel("1 行目へのコメント")).toBeVisible();
  await page.locator(".grading-line-text").first().click();
  await expect(page.getByLabel("1 行目へのコメント")).toHaveCount(0);
  await page.getByRole("button", { name: "行 1 を選択" }).click();
  await page
    .getByRole("button", { name: "行 3 を選択" })
    .click({ modifiers: ["Shift"] });
  await page.getByLabel("1〜3 行目へのコメント").fill("DNS設定を確認");
  await expect(page.locator(".grading-line-source.is-selected")).toHaveCount(3);
  await page.getByRole("button", { name: "行 2 を選択" }).click();
  await expect(page.getByLabel("1〜3 行目へのコメント")).toHaveCount(0);
  await expect(
    page.getByRole("button", { name: "1〜3 行目の下書き" }),
  ).toBeVisible();
  await page.getByRole("button", { name: "行 2 を選択" }).click();
  await expect(page.getByLabel("1〜3 行目へのコメント")).toHaveValue(
    "DNS設定を確認",
  );
  await page.getByLabel("得点").fill("88");
  await page.getByLabel("全体コメント").fill("再採点");
  await page.getByRole("button", { name: "採点結果を送信" }).click();
  const markingConfirm = page.getByRole("dialog", {
    name: "採点結果を保存",
  });
  await expect(markingConfirm).toContainText("行コメント: 1件");
  await markingConfirm.getByRole("button", { name: "キャンセル" }).click();
  expect(
    state.requests.some(
      (item) =>
        item.path === "/api/v1/admin/marking-results" && item.method === "POST",
    ),
  ).toBe(false);
  await page.getByRole("button", { name: "採点結果を送信" }).click();
  await markingConfirm.getByRole("button", { name: "採点結果を保存" }).click();
  await expect(page.getByText("再採点", { exact: true })).toBeVisible();
  await expect(page.locator(".grading-line-history")).toContainText(
    "DNS設定を確認",
  );
  await expect(page.locator(".grading-line-source.has-history")).toHaveCount(3);
  const historyComments = page.locator(".grading-history-comments");
  await historyComments.getByText("行コメント 1 件").click();
  await expect(historyComments).toContainText("DNS設定を確認");
  await historyComments
    .getByRole("button", { name: "1〜3 行目のコメントを削除" })
    .click();
  const deleteConfirm = page.getByRole("dialog", { name: "行コメントを削除" });
  await deleteConfirm.getByRole("button", { name: "キャンセル" }).click();
  await expect(page.locator(".grading-line-history")).toContainText(
    "DNS設定を確認",
  );
  await historyComments
    .getByRole("button", { name: "1〜3 行目のコメントを削除" })
    .click();
  await deleteConfirm.getByRole("button", { name: "行コメントを削除" }).click();
  await expect(page.locator(".grading-line-history")).toHaveCount(0);
  await expect(historyComments).toContainText("削除済み 1 件");
  await expect(historyComments).not.toContainText("DNS設定を確認");
  expect(
    state.requests.some(
      (item) =>
        item.method === "DELETE" && item.path.endsWith("/line-comments/0"),
    ),
  ).toBe(true);
  expect(
    state.requests.some(
      (item) =>
        item.method === "POST" &&
        item.path === "/api/v1/admin/marking-results" &&
        (item.body as { score?: number; line_comments?: unknown[] })?.score ===
          88 &&
        (
          item.body as {
            line_comments?: Array<{
              line_number: number;
              end_line_number?: number;
              body: string;
            }>;
          }
        )?.line_comments?.[0]?.line_number === 1 &&
        (
          item.body as {
            line_comments?: Array<{ end_line_number?: number }>;
          }
        )?.line_comments?.[0]?.end_line_number === 3,
    ),
  ).toBe(true);

  await page.goto("/admin/scores");
  await expect(page.getByText("凍結中")).toBeVisible();
  await page.getByRole("button", { name: "得点を再計算" }).click();
  await page
    .getByRole("dialog", { name: "得点再計算" })
    .getByRole("button", { name: "得点再計算" })
    .click();
  await expect(page.getByText("得点再計算しました")).toBeVisible();
  await page.getByRole("button", { name: "最終得点を公開" }).click();
  const revealConfirm = page.getByRole("dialog", { name: "最終得点公開" });
  await expect(revealConfirm).toContainText("取り消せない操作");
  await revealConfirm.getByRole("button", { name: "キャンセル" }).click();
  expect(
    state.requests.some(
      (item) => item.path === "/api/v1/admin/scores/actions/reveal-final",
    ),
  ).toBe(false);

  await page.goto("/admin/settings");
  await page.getByLabel("Markdown").fill("# Updated rule");
  await page.getByRole("button", { name: "ルールを更新" }).click();
  await expect(page.getByText("ルール更新しました")).toBeVisible();
  await page.getByLabel("凍結時刻").fill("");
  await page.getByRole("button", { name: "凍結時刻を保存" }).click();
  await expect(page.getByText("凍結時刻更新しました")).toBeVisible();

  await page.goto("/admin/submissions/A01/12/999");
  await expect(page.getByText("指定された回答は存在しません")).toBeVisible();
});

test("採点看板の担当・移動と過去回答の比較", async ({ page }) => {
  const state = await installAdminApi(page);
  await page.clock.install({ time: new Date("2026-09-01T12:25:00+09:00") });
  state.answers.push({
    ...structuredClone(answer),
    reference: { ...answer.reference, answer_number: 2 },
    body: {
      type: "DESCRIPTIVE",
      body: "# 2回目の回答\n\n設定を修正しました。\n追記しました。",
    },
    score: null,
    workflow: {
      status: "WAITING",
      assignee_discord_id: null,
      assignment_source: "DEFAULT",
      revision: 0,
    },
  });
  state.answers.push({
    ...structuredClone(answer),
    reference: { ...answer.reference, answer_number: 3 },
    submitted_at: "2026-09-01T12:20:00+09:00",
    score: null,
    workflow: {
      status: "WAITING",
      assignee_discord_id: null,
      assignment_source: "DEFAULT",
      revision: 0,
    },
  });
  await page.goto("/admin/submissions/");
  const waiting = page.getByRole("region", { name: "採点待ち" });
  const inProgress = page.getByRole("region", { name: "採点中" });
  await expect(waiting.getByText("Existing Team / #2")).toBeVisible();
  const secondCard = waiting
    .locator(".grading-card")
    .filter({ hasText: "Existing Team / #2" });
  await expect(secondCard).toContainText("提出から 25 分");
  await page.clock.fastForward(60_000);
  await expect(secondCard).toContainText("提出から 26 分");
  await secondCard.getByRole("button", { name: "自分が引き受ける" }).click();
  await expect(secondCard.getByText("自分の担当")).toBeVisible();
  await expect(waiting.locator(".grading-card").first()).toContainText(
    "Existing Team / #2",
  );
  await secondCard.dragTo(inProgress);
  await expect(inProgress.getByText("Existing Team / #2")).toBeVisible();
  await inProgress.getByRole("link", { name: /A01: Reliable DNS/ }).click();
  await expect(page.locator(".grading-detail-age")).toHaveText(
    "提出から 26 分",
  );
  await expect(
    page.locator(".grading-line-text").getByText("# 2回目の回答"),
  ).toBeVisible();
  await expect(page.locator(".grading-pane")).toHaveCount(3);
  await expect(page.getByLabel("比較する回答")).toHaveValue("1");
  await expect(page.locator(".grading-diff-hunk")).toContainText(
    "@@ -1,3 +1,4 @@",
  );
  await expect(page.locator(".grading-answer-line.is-removed")).toContainText(
    "# 回答",
  );
  await expect(
    page
      .locator(".grading-answer-line.is-added")
      .filter({ hasText: "# 2回目の回答" }),
  ).toBeVisible();
  await expect(
    page
      .locator(".grading-answer-line.is-added")
      .filter({ hasText: "追記しました。" }),
  ).toBeVisible();
  await page.getByRole("button", { name: "今回の回答のみ" }).click();
  await expect(page.locator(".grading-answer-line.is-removed")).toHaveCount(0);
  await expect(page.locator(".grading-line-source")).toHaveCount(4);
  await page.getByRole("button", { name: "差分", exact: true }).click();
  await page.getByRole("button", { name: "行 4 を選択" }).click();
  await expect(page.getByLabel("4 行目へのコメント")).toBeVisible();
  const problemPane = page.locator(".grading-pane").first();
  const initialWidth = (await problemPane.boundingBox())!.width;
  const divider = (await page
    .locator(".grading-resizer")
    .first()
    .boundingBox())!;
  const dragY = Math.max(100, divider.y + 30);
  await page.mouse.move(divider.x + divider.width / 2, dragY);
  await page.mouse.down();
  await page.mouse.move(divider.x + divider.width / 2 + 70, dragY);
  await page.mouse.up();
  expect((await problemPane.boundingBox())!.width).toBeGreaterThan(
    initialWidth + 30,
  );
  await page.getByRole("button", { name: "回答を折りたたむ" }).click();
  await page.getByRole("button", { name: "回答を展開" }).click();
  await expect(page.getByLabel("比較する回答")).toHaveValue("1");
  await page.goto("/admin/submissions/A01/12/3");
  await expect(page.getByLabel("比較する回答")).toHaveValue("2");
  await expect(page.locator(".grading-answer-line.is-removed")).toHaveCount(2);
  await page.getByLabel("比較する回答").selectOption("1");
  await expect(page.locator(".grading-no-diff")).toBeVisible();
  await expect(page.locator(".grading-answer-line.is-added")).toHaveCount(0);
  await page.getByRole("button", { name: "今回の回答のみ" }).click();
  await expect(page.getByRole("button", { name: "行 3 を選択" })).toBeVisible();

  const longLines = Array.from(
    { length: 20 },
    (_, index) => `line ${index + 1}`,
  );
  state.answers[0]!.body = { type: "DESCRIPTIVE", body: longLines.join("\n") };
  longLines[9] = "changed line 10";
  state.answers[1]!.body = { type: "DESCRIPTIVE", body: longLines.join("\n") };
  await page.goto("/admin/submissions/A01/12/2");
  await expect(page.locator(".grading-diff-omitted").first()).toBeVisible();
  await expect(page.getByRole("button", { name: "行 1 を選択" })).toHaveCount(
    0,
  );
  await page.locator(".grading-diff-omitted").first().click();
  await expect(page.getByRole("button", { name: "行 1 を選択" })).toBeVisible();
});

test("redeployはQUEUEDを楽観表示しSSE更新・履歴・一回syncのみを使う", async ({
  page,
}) => {
  await installEventSourceMock(page);
  const state = await installAdminApi(page);
  await page.goto("/admin/deployments");
  await page
    .getByRole("combobox", { name: "チーム", exact: true })
    .selectOption("12");
  await page
    .getByRole("combobox", { name: "問題", exact: true })
    .selectOption("A01");
  await page
    .getByRole("button", { name: "選択した組み合わせを再展開" })
    .click();
  await expect(page.getByText("待機中", { exact: true })).toBeVisible();

  const getsBeforeSse = state.deploymentGets;
  await emitSse(page, "/api/v1/admin/deployments/stream", "deployment", {
    deployment: completedDeployment,
  });
  await expect(page.getByText("完了", { exact: true })).toBeVisible();
  expect(state.deploymentGets).toBe(getsBeforeSse);

  await page.getByRole("button", { name: "履歴" }).click();
  await expect(page.getByText("ready", { exact: true })).toBeVisible();
  await page.getByRole("button", { name: "閉じる" }).click();
  await page.getByRole("button", { name: "一回同期" }).click();
  await page
    .getByRole("dialog", { name: "手動同期" })
    .getByRole("button", { name: "手動同期" })
    .click();
  await expect(page.getByText("手動同期しました")).toBeVisible();
  expect(
    state.requests.filter(
      (item) => item.method === "POST" && item.path.endsWith("/sync"),
    ),
  ).toHaveLength(1);
  const getsAfterSync = state.deploymentGets;
  await page.waitForTimeout(1_000);
  expect(state.deploymentGets).toBe(getsAfterSync);
});

test("運営がチームカラーを保存し再読込できる", async ({ page }, testInfo) => {
  const state = await installAdminApi(page);
  await page.goto("/admin/teams/12");
  await page
    .getByRole("combobox", { name: "チームカラー", exact: true })
    .selectOption("#0083C3");
  await page.getByRole("button", { name: "チーム情報を保存" }).click();
  await expect(
    page.getByText("チームを更新しました", { exact: true }),
  ).toBeVisible();
  expect(
    state.requests.some(
      (r) =>
        r.method === "PATCH" &&
        r.path === "/api/v1/admin/teams/12" &&
        (r.body as { color?: string }).color === "#0083C3",
    ),
  ).toBe(true);
  await page.reload();
  await expect(
    page.getByRole("combobox", { name: "チームカラー", exact: true }),
  ).toHaveValue("#0083C3");
  await page.screenshot({
    path: testInfo.outputPath("admin-team-desktop.png"),
    fullPage: true,
  });
  await page.setViewportSize({ width: 390, height: 844 });
  expect(
    await page.evaluate(() => document.documentElement.scrollWidth),
  ).toBeLessThanOrEqual(390);
  await page.screenshot({
    path: testInfo.outputPath("admin-team-mobile.png"),
    fullPage: true,
  });
});
