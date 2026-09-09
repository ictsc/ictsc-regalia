import { expect, test, type Request } from "@playwright/test";
import {
  deferred,
  emitSse,
  fulfillJson,
  installEventSourceMock,
} from "./support/http";

const commit = "a".repeat(40);
const requestedAt = "2026-09-01T12:00:00+09:00";

const queuedDeployment = {
  revision: 1,
  status: "QUEUED",
  requested_at: requestedAt,
  penalty: 0,
  allowed_request_count: 1,
  content_commit: commit,
};

const completedDeployment = {
  ...queuedDeployment,
  status: "COMPLETED",
};

function requestBody(request: Request): unknown {
  const raw = request.postData();
  return raw == null ? null : JSON.parse(raw);
}

test("viewerから問題を開き、再展開のQUEUED表示をSSEで完了へ更新する", async ({
  page,
}) => {
  const postStarted = deferred();
  const allowPostResponse = deferred();
  const deploymentRequests: unknown[] = [];
  let deploymentListRequests = 0;

  await installEventSourceMock(page);
  await page.route("**/api/v1/**", async (route) => {
    const request = route.request();
    const path = new URL(request.url()).pathname;

    if (
      path === "/api/v1/contestant/problems/A01/deployments" &&
      request.method() === "GET"
    ) {
      deploymentListRequests += 1;
      return fulfillJson(route, { deployments: [] });
    }

    if (
      path === "/api/v1/contestant/problems/A01/deployments" &&
      request.method() === "POST"
    ) {
      deploymentRequests.push(requestBody(request));
      postStarted.resolve();
      await allowPostResponse.promise;
      return fulfillJson(route, { deployment: queuedDeployment }, 201);
    }

    const responses: Record<string, unknown> = {
      "/api/v1/viewer": {
        viewer: {
          state: "CONTESTANT",
          profile: {
            name: "alice",
            display_name: "Alice",
            self_introduction: "E2E contestant",
          },
          team: {
            code: 12,
            name: "E2E Team",
            organization: "ICTSC University",
            member_limit: 4,
          },
          impersonated_by: null,
        },
      },
      "/api/v1/contestant/schedule": {
        sections: [
          {
            slug: "day1",
            beginning: "2026-01-01T00:00:00+09:00",
            ending: "2027-01-01T00:00:00+09:00",
          },
        ],
      },
      "/api/v1/contestant/announcements": { announcements: [] },
      "/api/v1/contestant/problems/A01": {
        problem: {
          code: "A01",
          title: "Reliable DNS",
          max_score: 100,
          category: "Network",
          section_slug: "day1",
          type: "NORMAL",
          body: "# Restore service\n\nKeep DNS available.",
          score: null,
          deployment: {
            redeployable: true,
            penalty_threshold: 2,
          },
          submission_status: {
            is_submittable: true,
            submittable_from: "2026-01-01T00:00:00+09:00",
            submittable_until: "2027-01-01T00:00:00+09:00",
          },
        },
      },
      "/api/v1/contestant/problems/A01/answers": {
        answers: [
          {
            number: 1,
            type: "DESCRIPTIVE",
            submitted_at: requestedAt,
            score: null,
            content_commit: commit,
          },
        ],
        submit_interval_seconds: 1200,
        last_submitted_at: null,
      },
    };
    const response = responses[path];
    if (response != null) return fulfillJson(route, response);
    return fulfillJson(
      route,
      {
        type: "about:blank",
        title: "Unexpected E2E request",
        status: 404,
        detail: `${request.method()} ${path}`,
        code: "not_found",
      },
      404,
    );
  });

  await page.goto("/problems/A01");
  await expect(page).toHaveTitle("A01 Reliable DNS / ICTSC REGALIA", {
    timeout: 20000,
  });
  await expect(
    page.getByRole("heading", { name: "A01:Reliable DNS." }),
  ).toBeVisible();
  await expect(page.getByText("Keep DNS available.")).toBeVisible();
  page.on("dialog", (dialog) => dialog.accept());
  const redeployButton = page.getByRole("button", { name: "環境をリセット" });
  await redeployButton.click();

  await postStarted.promise;
  await expect(page.getByText("QUEUED — 再展開を要求中")).toBeVisible();
  await expect(redeployButton).toBeDisabled();
  expect(deploymentRequests).toEqual([null]);

  allowPostResponse.resolve();
  await expect(redeployButton).toBeDisabled();

  await emitSse(
    page,
    "/api/v1/contestant/problems/A01/deployments/stream",
    "deployment",
    { deployment: completedDeployment },
  );
  await expect(page.getByText(/#1 COMPLETED/)).toBeVisible();
  await expect(redeployButton).toBeEnabled();
  expect(deploymentListRequests).toBe(1);
});

test("ANONYMOUSとDISCORD_AUTHENTICATEDをsignin/signupへ誘導する", async ({
  page,
}) => {
  let viewer: unknown = { state: "ANONYMOUS" };
  await page.route("**/api/v1/**", async (route) => {
    const path = new URL(route.request().url()).pathname;
    if (path === "/api/v1/viewer") return fulfillJson(route, { viewer });
    if (path === "/api/v1/contestant/schedule") {
      return fulfillJson(route, { sections: [] });
    }
    return fulfillJson(
      route,
      {
        type: "about:blank",
        title: "not found",
        status: 404,
        detail: path,
        code: "not_found",
      },
      404,
    );
  });

  await page.goto("/signin/");
  await expect(
    page.getByRole("link", { name: /Discordでログイン/ }),
  ).toBeVisible();

  viewer = {
    state: "DISCORD_AUTHENTICATED",
    discord: {
      id: "123456789012345678",
      username: "pre-user",
      display_name: "Pre User",
    },
  };
  await page.goto("/signup");
  await expect(page.getByRole("heading", { name: "参加登録" })).toBeVisible();
  await expect(page.getByLabel("競技者名")).toHaveValue("pre-user");
});

test("profile更新と凍結rankingを表示する", async ({ page }) => {
  let profile = {
    name: "alice",
    display_name: "Alice",
    self_introduction: "before",
  };
  await page.route("**/api/v1/**", async (route) => {
    const request = route.request();
    const path = new URL(request.url()).pathname;
    if (path === "/api/v1/viewer") {
      return fulfillJson(route, {
        viewer: {
          state: "CONTESTANT",
          profile,
          team: {
            code: 12,
            name: "E2E Team",
            organization: "ICTSC University",
            member_limit: 4,
          },
          impersonated_by: null,
        },
      });
    }
    if (path === "/api/v1/contestant/schedule") {
      return fulfillJson(route, { sections: [] });
    }
    if (path === "/api/v1/contestant/profile" && request.method() === "GET") {
      return fulfillJson(route, { profile });
    }
    if (path === "/api/v1/contestant/profile" && request.method() === "PATCH") {
      profile = { ...profile, ...(requestBody(request) as object) };
      return fulfillJson(route, { profile });
    }
    if (path === "/api/v1/contestant/ranking") {
      return fulfillJson(route, {
        ranking: [
          {
            rank: 1,
            team_code: 12,
            team_name: "E2E Team",
            organization: "ICTSC University",
            score: 350,
            last_effective_submission_at: requestedAt,
          },
        ],
        frozen: true,
        frozen_at: requestedAt,
      });
    }
    return fulfillJson(
      route,
      {
        type: "about:blank",
        title: "not found",
        status: 404,
        detail: path,
        code: "not_found",
      },
      404,
    );
  });

  await page.goto("/profile");
  await page.getByLabel("表示名").fill("Alice Updated");
  await page.getByLabel("自己紹介").fill("after");
  await page.getByRole("button", { name: "更新する" }).click();
  await expect(page.getByText("プロフィールを更新しました")).toBeVisible();

  await page.goto("/ranking");
  await expect(page.getByText(/順位表は凍結中です/)).toBeVisible();
  await expect(
    page.locator(".ranking-team").getByText("E2E Team", { exact: true }),
  ).toBeVisible();
});

test("回答429のRetry-Afterをcountdownへ反映する", async ({ page }) => {
  await page.route("**/api/v1/**", async (route) => {
    const request = route.request();
    const path = new URL(request.url()).pathname;
    if (path === "/api/v1/viewer") {
      return fulfillJson(route, {
        viewer: {
          state: "CONTESTANT",
          profile: {
            name: "alice",
            display_name: "Alice",
            self_introduction: "",
          },
          team: {
            code: 12,
            name: "E2E Team",
            organization: "ICTSC University",
            member_limit: 4,
          },
          impersonated_by: null,
        },
      });
    }
    if (path === "/api/v1/contestant/schedule") {
      return fulfillJson(route, {
        sections: [
          {
            slug: "day1",
            beginning: "2026-01-01T00:00:00+09:00",
            ending: "2027-01-01T00:00:00+09:00",
          },
        ],
      });
    }
    if (path === "/api/v1/contestant/announcements") {
      return fulfillJson(route, { announcements: [] });
    }
    if (path === "/api/v1/contestant/problems/A01") {
      return fulfillJson(route, {
        problem: {
          code: "A01",
          title: "Reliable DNS",
          max_score: 100,
          category: "Network",
          section_slug: "day1",
          type: "DESCRIPTIVE",
          body: "# Restore service",
          score: null,
          deployment: {
            redeployable: true,
            penalty_threshold: 2,
          },
          submission_status: {
            is_submittable: true,
            submittable_from: "2026-01-01T00:00:00+09:00",
            submittable_until: "2027-01-01T00:00:00+09:00",
          },
        },
      });
    }
    if (
      path === "/api/v1/contestant/problems/A01/answers" &&
      request.method() === "GET"
    ) {
      return fulfillJson(route, {
        answers: [],
        submit_interval_seconds: 1200,
        last_submitted_at: null,
      });
    }
    if (
      path === "/api/v1/contestant/problems/A01/answers" &&
      request.method() === "POST"
    ) {
      return route.fulfill({
        status: 429,
        headers: {
          "Content-Type": "application/problem+json",
          "Retry-After": "2",
        },
        body: JSON.stringify({
          type: "about:blank",
          title: "Too many answers",
          status: 429,
          detail: "次の回答まで待ってください",
          code: "answer_rate_limited",
        }),
      });
    }
    if (path === "/api/v1/contestant/problems/A01/deployments") {
      return fulfillJson(route, { deployments: [] });
    }
    return fulfillJson(
      route,
      {
        type: "about:blank",
        title: "not found",
        status: 404,
        detail: path,
        code: "not_found",
      },
      404,
    );
  });

  await page.goto("/problems/A01");
  await page.getByLabel("回答", { exact: true }).fill("設定を修正しました");
  await page.getByRole("button", { name: "回答を提出 ↗" }).click();
  await expect(
    page.getByRole("status").filter({ hasText: /再提出まで/ }),
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "回答を提出 ↗" }),
  ).toBeDisabled();
});

test("Nuxtの問題一覧、下書き復元、利用者分離とモバイル表示", async ({
  page,
}, testInfo) => {
  const { installCompetitionApi } = await import("./support/competition");
  const state = await installCompetitionApi(page);
  await page.goto("/problems");
  await expect(page.locator(".problem-row:not(.problem-heading)")).toHaveCount(
    11,
  );
  await page.getByText("絞り込み", { exact: true }).click();
  await page
    .getByRole("combobox", { name: "カテゴリ", exact: true })
    .selectOption("無線");
  await expect(page.locator(".problem-row:not(.problem-heading)")).toHaveCount(
    2,
  );
  await page
    .getByRole("combobox", { name: "カテゴリ", exact: true })
    .selectOption("all");
  await page.getByText("絞り込み", { exact: true }).click();
  await page.screenshot({
    path: testInfo.outputPath("problems-desktop.png"),
    fullPage: true,
  });
  await page.goto("/problems/A04");
  await page
    .getByLabel("回答", { exact: true })
    .fill("原因: 設定の不整合\n復旧確認済み");
  await expect(page.getByText(/このブラウザに保存済み/)).toBeVisible();
  await page.reload();
  await expect(page.getByLabel("回答", { exact: true })).toHaveValue(
    "原因: 設定の不整合\n復旧確認済み",
  );
  await page.getByRole("button", { name: "回答を提出 ↗" }).click();
  await expect(
    page.getByText("回答を提出しました", { exact: true }),
  ).toBeVisible();
  expect(state.submissions).toEqual(["原因: 設定の不整合\n復旧確認済み"]);
  await page.screenshot({
    path: testInfo.outputPath("problem-desktop.png"),
    fullPage: true,
  });
  state.name = "bob";
  await page.reload();
  await expect(page.getByLabel("回答", { exact: true })).toHaveValue("");
  await page.setViewportSize({ width: 390, height: 844 });
  await page.screenshot({
    path: testInfo.outputPath("problem-mobile.png"),
    fullPage: true,
  });
  expect(
    await page.evaluate(() => document.documentElement.scrollWidth),
  ).toBeLessThanOrEqual(390);
  await page.goto("/problems/A01");
  await expect(
    page.getByRole("button", { name: "回答を提出 ↗" }),
  ).toBeDisabled();
  await page.goto("/ranking");
  await expect(page.locator(".ranking-row.is-team")).toContainText(
    "KERNEL PANIC",
  );
  await page.goto("/announces");
  await expect(
    page.getByRole("heading", { name: "競技を開始しました" }),
  ).toBeVisible();
  await page.screenshot({
    path: testInfo.outputPath("notices-mobile.png"),
    fullPage: true,
  });
});

test("下書き保存に失敗した場合は保存済みと表示しない", async ({ page }) => {
  const { installCompetitionApi } = await import("./support/competition");
  await installCompetitionApi(page);
  await page.addInitScript(() => {
    Storage.prototype.setItem = () => {
      throw new DOMException("quota", "QuotaExceededError");
    };
  });
  await page.goto("/problems/A04");
  await page.getByLabel("回答", { exact: true }).fill("消してはいけない回答");
  await expect(page.getByText(/下書きを保存できません/)).toBeVisible();
  await expect(page.getByLabel("回答", { exact: true })).toHaveValue(
    "消してはいけない回答",
  );
  await expect(page.getByText(/このブラウザに保存済み/)).toHaveCount(0);
});
