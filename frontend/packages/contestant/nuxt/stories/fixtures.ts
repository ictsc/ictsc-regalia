// Fake data for Storybook. It mirrors e2e/support/competition.ts and must not
// contain real participant names because the Storybook is published.
import type { Answer } from "~/features/answer";
import type { Deployment } from "~/features/deployment";
import type { Problem, Rank, ScheduleEntry } from "~/features/models";

export const fixedNow = Date.parse("2026-09-02T12:00:00+09:00");
export const teamColor = "#0083C3";

export const sections = {
  day1: {
    name: "day1",
    startAt: "2026-09-01T00:00:00+09:00",
    endAt: "2026-09-02T00:00:00+09:00",
  },
  day2: {
    name: "day2",
    startAt: "2026-09-02T00:00:00+09:00",
    endAt: "2026-09-02T14:13:48+09:00",
  },
  both: {
    name: "both",
    startAt: "2026-09-01T00:00:00+09:00",
    endAt: "2026-09-02T14:13:48+09:00",
  },
} satisfies Record<string, ScheduleEntry>;

type SectionSlug = keyof typeof sections;

export function makeProblem(
  overrides: Partial<Problem> & { code: string; sectionSlug?: SectionSlug },
): Problem {
  const section = sections[overrides.sectionSlug ?? "day2"];
  return {
    title: `問題 ${overrides.code}`,
    category: "ネットワーク",
    maxScore: 200,
    submissionableSchedules: [section],
    submissionStatus: {
      isSubmittable: section.name !== "day1",
      submittableFrom: section.startAt,
      submittableUntil: section.endAt,
    },
    ...overrides,
  };
}

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

export const problems: Problem[] = definitions.map(
  ([code, title, category, maxScore, sectionSlug, score]) =>
    makeProblem({
      code,
      title,
      category,
      maxScore,
      sectionSlug,
      score:
        score === null
          ? undefined
          : { markedScore: score, penalty: 0, score, maxScore },
    }),
);

const teamNames = [
  "NETWORK LAB",
  "PACKET",
  "KERNEL PANIC",
  "ROUTE",
  "SUBNET",
  "LOOPBACK",
  "TRACEROUTE",
  "ANYCAST",
];

export const ownTeamCode = 12;

// Ranks 11-14 surround the own team (code 12) so neighbor display is visible.
export const ranking: Rank[] = teamNames.map((teamName, index) => {
  const rank = index === 0 ? 1 : index + 10;
  return {
    rank,
    teamCode: rank,
    teamName,
    organization: "ICTSC University",
    score: index === 0 ? 2130 : 1350 - index * 60,
  };
});

export const problemBody = `## 状況

監視対象のコンテナ\`api-gateway\`が起動して数秒後に終了し、再起動を繰り返しています。利用者からは管理画面へ断続的に接続できないという報告があります。

## 環境

| ホスト | アドレス |
| --- | --- |
| router-a | 10.20.1.1 |
| app-01 | 10.20.1.10 |

\`\`\`bash
docker compose ps
docker compose logs api-gateway
\`\`\`

## 確認事項

修正後はコンテナの状態とログを確認し、原因、実施した変更、復旧確認の結果を回答欄に記載してください。`;

export function makeDeployment(
  revision: number,
  status: Deployment["status"],
  overrides: Partial<Deployment> = {},
): Deployment {
  return {
    revision,
    status,
    requestedAt: new Date(fixedNow - (10 - revision) * 600_000).toISOString(),
    allowedDeploymentCount: Math.max(0, 3 - revision),
    thresholdExceeded: revision > 3,
    penalty: revision > 3 ? 10 * (revision - 3) : 0,
    ...overrides,
  };
}

export const answers: Answer[] = [3, 2, 1].map((id) => ({
  id,
  submittedAt: new Date(fixedNow - (4 - id) * 1_800_000).toISOString(),
  contentCommit: "0f3c2a1",
  score:
    id === 3
      ? undefined
      : { markedScore: id * 60, penalty: 0, score: id * 60, maxScore: 300 },
}));

export const answerBody = `## 原因

\`api-gateway\` の環境変数 \`UPSTREAM_URL\` が誤っていたため、起動直後の疎通確認に失敗して終了していました。

## 対応

\`compose.yaml\` の値を修正し、\`docker compose up -d\` で再作成しました。

## 確認

10分間再起動が発生しないこと、管理画面が200を返すことを確認しました。`;
