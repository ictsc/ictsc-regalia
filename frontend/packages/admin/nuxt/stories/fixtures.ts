// Fake data for admin Storybook stories. The Storybook is published, so use
// only made-up names.
import type { components } from "@ictsc/api";

type Schemas = components["schemas"];

export const fixedNow = Date.parse("2026-09-02T12:00:00+09:00");
export const adminName = "運営スタッフ";

export const contentStatus: Schemas["ContentStatus"] = {
  state: "SUCCEEDED",
  active_commit: "9f2c4e1a7b3d5c6e8f0a1b2c3d4e5f6a7b8c9d0e",
  activated_at: "2026-09-02T08:30:00+09:00",
  source_ref: "main",
  last_attempt_commit: "9f2c4e1a7b3d5c6e8f0a1b2c3d4e5f6a7b8c9d0e",
  last_attempt_at: "2026-09-02T08:30:00+09:00",
  last_error: null,
  serving_last_known_good: false,
};

const teamNames = [
  "NETWORK LAB",
  "PACKET",
  "KERNEL PANIC",
  "ROUTE",
  "SUBNET",
  "LOOPBACK",
];
const colors = [
  "#0083C3",
  "#FF6C2F",
  "#397E58",
  "#BF53B6",
  "#FAAF3F",
  "#00CDC2",
] as const;

export const teams: Schemas["Team"][] = teamNames.map((name, index) => ({
  code: index + 2,
  name,
  organization: index % 2 ? "ICTSC University" : "ICTSC 高専",
  member_limit: 4,
  color: colors[index],
}));

const memberNames = ["alice", "bob", "carol", "dave", "erin", "frank"];

export const contestants: Schemas["AdminContestant"][] = memberNames.map(
  (name, index) => ({
    profile: {
      name,
      display_name: name[0]!.toUpperCase() + name.slice(1),
      self_introduction: "",
    },
    team: teams[index % 3]!,
    discord_id: `11000000000000000${index}`,
  }),
);

export const invitations: Schemas["Invitation"][] = [
  {
    code: "KP-7Q2M-9XAD",
    team_code: 4,
    created_at: "2026-08-20T10:00:00+09:00",
    expires_at: "2026-09-10T23:59:00+09:00",
  },
  {
    code: "KP-3RZK-1BCE",
    team_code: 4,
    created_at: "2026-08-01T10:00:00+09:00",
    expires_at: "2026-08-15T23:59:00+09:00",
  },
];

export const catalog: Schemas["ProblemCatalogEntry"][] = [
  {
    code: "A01",
    title: "DNSの名前解決ができない",
    max_score: 300,
    category: "ネットワーク",
  },
  {
    code: "A02",
    title: "Webサーバーに接続できない",
    max_score: 150,
    category: "サーバー",
  },
  {
    code: "B02",
    title: "コンテナが起動を繰り返す",
    max_score: 300,
    category: "サーバー",
  },
  {
    code: "B04",
    title: "VPNトンネルが確立しない",
    max_score: 200,
    category: "ネットワーク",
  },
];

export const sections: Schemas["Section"][] = [
  {
    slug: "day1",
    beginning: "2026-09-01T00:00:00+09:00",
    ending: "2026-09-02T00:00:00+09:00",
    problems: catalog.slice(0, 2),
  },
  {
    slug: "day2",
    beginning: "2026-09-02T00:00:00+09:00",
    ending: "2026-09-02T14:13:48+09:00",
    problems: catalog.slice(2),
  },
];

export const announcements: Schemas["Announcement"][] = [
  {
    slug: "day2-start",
    title: "2日目の競技を開始しました",
    markdown:
      "2日目の問題を公開しました。\n\n- 受付終了: **14:13**\n- 1日目の問題は閲覧のみです",
    effective_from: "2026-09-02T09:00:00+09:00",
  },
  {
    slug: "start",
    title: "競技を開始しました",
    markdown: "問題一覧から問題を確認し、回答を提出してください。",
    effective_from: "2026-09-01T09:00:00+09:00",
  },
];

export const adminProblem: Schemas["AdminProblem"] = {
  code: "B02",
  title: "コンテナが起動を繰り返す",
  max_score: 300,
  category: "サーバー",
  section_slug: "day2",
  type: "DESCRIPTIVE",
  body: "## 状況\n\n監視対象のコンテナ`api-gateway`が起動して数秒後に終了し、再起動を繰り返しています。\n\n| ホスト | アドレス |\n| --- | --- |\n| app-01 | 10.20.1.10 |",
  explanation:
    "## 想定解\n\n`UPSTREAM_URL` の誤りを修正し、`docker compose up -d` で再作成する。\n\n## 採点基準\n\n- 原因の特定: 150点\n- 復旧と確認: 150点",
  redeploy_rule: {
    type: "PERCENTAGE_PENALTY",
    penalty_threshold: 2,
    penalty_percentage: 10,
  },
  content_commit: contentStatus.active_commit!,
  default_grader_discord_id: null,
};

export const ranking: Schemas["RankingResponse"] = {
  ranking: teams.map((team, index) => ({
    rank: index + 1,
    team_code: team.code,
    team_name: team.name,
    organization: team.organization,
    score: 2100 - index * 230,
    last_effective_submission_at: null,
  })),
  frozen: false,
  frozen_at: null,
  hidden: false,
};

export const scores: Schemas["AdminScore"][] = teams
  .slice(0, 3)
  .flatMap((team, t) =>
    catalog.slice(0, 2).map((problem, p) => ({
      team,
      problem,
      marked_score: problem.max_score - (t + p) * 30,
      penalty: t === 1 && p === 0 ? 30 : 0,
      score: problem.max_score - (t + p) * 30 - (t === 1 && p === 0 ? 30 : 0),
    })),
  );

const event = (
  id: number,
  status: Schemas["DeploymentStatus"],
  minutesAgo: number,
  message: string | null = null,
): Schemas["DeploymentEvent"] => ({
  event_id: `evt-${id}`,
  occurred_at: new Date(fixedNow - minutesAgo * 60_000).toISOString(),
  status,
  message,
});

export const deployments: Schemas["AdminDeployment"][] = [
  {
    team_code: 4,
    problem_code: "B02",
    revision: 2,
    latest_status: "DEPLOYING",
    events: [event(4, "DEPLOYING", 2), event(3, "QUEUED", 3)],
    content_commit: contentStatus.active_commit!,
  },
  {
    team_code: 3,
    problem_code: "A02",
    revision: 1,
    latest_status: "FAILED",
    events: [
      event(2, "FAILED", 20, "VM の起動がタイムアウトしました"),
      event(1, "QUEUED", 25),
    ],
    content_commit: contentStatus.active_commit!,
  },
  {
    team_code: 4,
    problem_code: "B02",
    revision: 1,
    latest_status: "COMPLETED",
    events: [event(0, "COMPLETED", 60)],
    content_commit: contentStatus.active_commit!,
  },
];

export const ruleMarkdown = `# 競技ルール

## 回答

- 回答は問題ごとに **20分に1回** まで提出できます
- 再展開は減点なしで2回まで行えます

| 区分 | 減点 |
| --- | --- |
| 3回目以降の再展開 | 10% |

得点は $score = marked - penalty$ で計算します。`;

export const myDiscordId = "110000000000000099";

const answerRows: [
  team: number,
  problem: number,
  number: number,
  minutesAgo: number,
  status: Schemas["AdminAnswerWorkflow"]["status"],
  assignee: string | null,
  claimed: boolean,
  score: number | null,
][] = [
  [4, 2, 3, 2, "WAITING", myDiscordId, false, null],
  [2, 0, 1, 6, "WAITING", "110000000000000001", false, null],
  [3, 1, 2, 14, "WAITING", null, false, null],
  [5, 3, 1, 33, "IN_PROGRESS", myDiscordId, true, null],
  [6, 2, 1, 48, "IN_PROGRESS", "110000000000000002", true, null],
  [2, 2, 1, 75, "COMPLETED", "110000000000000001", false, 300],
  [3, 0, 1, 120, "COMPLETED", myDiscordId, false, 180],
];

export const answers: Schemas["AdminAnswer"][] = answerRows.map(
  ([team, problem, number, minutesAgo, status, assignee, claimed, score]) => {
    const p = catalog[problem]!;
    return {
      reference: {
        team_code: team,
        problem_code: p.code,
        answer_number: number,
      },
      team: teams[team - 2]!,
      author: contestants[0]!.profile,
      problem: p,
      body: {
        type: "DESCRIPTIVE",
        body: "## 原因\n\n設定ファイルの値が誤っていました。\n\n## 対応\n\n値を修正して再起動しました。",
      },
      submitted_at: new Date(fixedNow - minutesAgo * 60_000).toISOString(),
      score:
        score === null
          ? null
          : { total: score, marked: score, penalty: 0, max: p.max_score },
      content_commit: contentStatus.active_commit!,
      workflow: {
        status,
        assignee_discord_id: assignee,
        assignment_source: claimed ? "CLAIMED" : "DEFAULT",
        revision: 1,
      },
      competition_mode: "STANDARD",
      latest: true,
      marking_blocked_reason: null,
    };
  },
);

const previousBody = `## 原因

api-gateway の設定を確認しました。

## 対応

コンテナを再起動しました。

## 確認

起動したことを確認しました。`;

const currentBody = `## 原因

api-gateway の環境変数 UPSTREAM_URL が誤っており、起動直後の疎通確認に失敗して終了していました。

## 対応

compose.yaml の UPSTREAM_URL を http://backend:8080 に修正し、docker compose up -d で再作成しました。

## 確認

10分間再起動が発生しないこと、管理画面が 200 を返すことを確認しました。`;

const detailBase = answers.find((a) => a.reference.problem_code === "B02")!;

export const detailAnswer: Schemas["AdminAnswer"] = {
  ...detailBase,
  reference: { team_code: 4, problem_code: "B02", answer_number: 3 },
  team: teams[2]!,
  body: { type: "DESCRIPTIVE", body: currentBody },
  score: { total: 240, marked: 240, penalty: 0, max: 300 },
  workflow: {
    status: "IN_PROGRESS",
    assignee_discord_id: myDiscordId,
    assignment_source: "CLAIMED",
    revision: 2,
  },
};

export const previousAnswers: Schemas["AdminAnswer"][] = [2, 1].map(
  (number) => ({
    ...detailAnswer,
    reference: { ...detailAnswer.reference, answer_number: number },
    body: { type: "DESCRIPTIVE", body: previousBody },
    submitted_at: new Date(fixedNow - (4 - number) * 1_800_000).toISOString(),
    score:
      number === 2 ? { total: 120, marked: 120, penalty: 0, max: 300 } : null,
  }),
);

export const markings: Schemas["MarkingResult"][] = [
  {
    id: "mark-2",
    answer: detailAnswer.reference,
    judge: { name: adminName },
    score: 240,
    rationale: "原因と対応は正しい。復旧確認の根拠をもう少し具体的に。",
    line_comments: [
      { line_number: 3, body: "原因の特定は十分です" },
      { line_number: 7, end_line_number: 7, body: "修正値の根拠も書くとよい" },
      {
        line_number: 11,
        body: "古いコメント",
        deleted_at: "2026-09-02T11:50:00+09:00",
        deleted_by: adminName,
      },
    ],
    created_at: "2026-09-02T11:45:00+09:00",
    visibility: "TEAM",
  },
];
