import type { Page } from "@playwright/test";
import { fulfillJson, installEventSourceMock } from "./http";
export async function installCompetitionApi(page: Page) {
  await page.clock.setFixedTime(new Date("2026-09-02T12:00:00+09:00"));
  await installEventSourceMock(page);
  const state = {
    name: "alice",
    color: "#0083C3",
    submissions: [] as string[],
  };
  const definitions = [
    ["A01", "DNSの名前解決ができない", "ネットワーク", 300, "day1", 300],
    ["A02", "Webサーバーに接続できない", "サーバー", 150, "both", 80],
    ["A03", "BGP経路が広報されない", "ネットワーク", 350, "day1", 0],
    ["A04", "無線LANに接続できない", "無線", 200, "day2", null],
    ["A05", "ログが保存されていない", "サーバー", 200, "both", 200],
    ["A06", "監視通知が届かない", "その他", 100, "day1", null],
    ["B01", "IPv6で外部へ通信できない", "ネットワーク", 250, "day1", 130],
    ["B02", "コンテナが起動を繰り返す", "サーバー", 300, "day2", 300],
    ["B03", "AP間でローミングできない", "無線", 250, "both", null],
    ["B04", "VPNトンネルが確立しない", "ネットワーク", 200, "day2", 70],
    ["B05", "時刻同期がずれている", "その他", 150, "both", 150],
  ] as const;
  const sections = [
    {
      slug: "day1",
      beginning: "2026-09-01T00:00:00+09:00",
      ending: "2026-09-02T00:00:00+09:00",
    },
    {
      slug: "day2",
      beginning: "2026-09-02T00:00:00+09:00",
      ending: "2026-09-02T14:13:48+09:00",
    },
    {
      slug: "both",
      beginning: "2026-09-01T00:00:00+09:00",
      ending: "2026-09-02T14:13:48+09:00",
    },
  ];
  const problems = definitions.map(
    ([code, title, category, max_score, section_slug, score]) => ({
      code,
      title,
      category,
      max_score,
      section_slug,
      type: "NORMAL",
      body: "## 状況\n\n指定されたサービスに接続できません。原因を特定し、必要な設定を修正してください。\n\n## 環境\n\n| ホスト | アドレス |\n| --- | --- |\n| router-a | 10.20.1.1 |\n\n```bash\nip route show\n```",
      score:
        score === null
          ? null
          : { marked_score: score, penalty: 0, score, max_score },
      deployment: { redeployable: true, penalty_threshold: 2 },
      submission_status: {
        is_submittable: section_slug !== "day1",
        submittable_from: "2026-09-01T00:00:00+09:00",
        submittable_until:
          section_slug === "day1"
            ? "2026-09-02T00:00:00+09:00"
            : "2026-09-02T14:13:48+09:00",
      },
    }),
  );
  await page.route("**/api/v1/**", async (route) => {
    const path = new URL(route.request().url()).pathname;
    if (path === "/api/v1/viewer")
      return fulfillJson(route, {
        viewer: {
          state: "CONTESTANT",
          profile: {
            name: state.name,
            display_name: state.name,
            self_introduction: "",
          },
          team: {
            code: 12,
            name: "KERNEL PANIC",
            organization: "ICTSC University",
            member_limit: 4,
            color: state.color,
          },
          impersonated_by: null,
        },
      });
    if (path === "/api/v1/contestant/problems")
      return fulfillJson(route, { problems });
    if (
      ["/api/v1/contestant/sections", "/api/v1/contestant/schedule"].includes(
        path,
      )
    )
      return fulfillJson(route, { sections });
    if (path === "/api/v1/contestant/ranking")
      return fulfillJson(route, {
        ranking: [
          {
            rank: 1,
            team_code: 2,
            team_name: "NETWORK LAB",
            organization: "ICTSC",
            score: 2130,
            last_effective_submission_at: null,
          },
          {
            rank: 11,
            team_code: 11,
            team_name: "PACKET",
            organization: "ICTSC",
            score: 1290,
            last_effective_submission_at: null,
          },
          {
            rank: 12,
            team_code: 12,
            team_name: "KERNEL PANIC",
            organization: "ICTSC",
            score: 1230,
            last_effective_submission_at: null,
          },
          {
            rank: 13,
            team_code: 13,
            team_name: "ROUTE",
            organization: "ICTSC",
            score: 1180,
            last_effective_submission_at: null,
          },
        ],
        frozen: false,
        frozen_at: null,
      });
    if (path === "/api/v1/contestant/announcements")
      return fulfillJson(route, {
        announcements: [
          {
            slug: "start",
            title: "競技を開始しました",
            markdown: "問題一覧から問題を確認し、回答を提出してください。",
            effective_from: "2026-09-01T12:00:00+09:00",
          },
        ],
      });
    const match = path.match(
      /^\/api\/v1\/contestant\/problems\/([^/]+)(?:\/(answers|deployments))?$/,
    );
    if (match) {
      const problem = problems.find((p) => p.code === match[1]);
      if (match[2] === "deployments")
        return fulfillJson(route, { deployments: [] });
      if (match[2] === "answers") {
        if (route.request().method() === "POST") {
          state.submissions.push(route.request().postDataJSON().body);
          return fulfillJson(route, { answer: { number: 1 } }, 201);
        }
        return fulfillJson(route, {
          answers: [],
          last_submitted_at: null,
          submit_interval_seconds: 1200,
        });
      }
      return fulfillJson(route, { problem });
    }
    return fulfillJson(
      route,
      { title: "Not found", status: 404, detail: path },
      404,
    );
  });
  return state;
}
