// Fake data for Storybook. It mirrors e2e/support/competition.ts and must not
// contain real participant names because the Storybook is published.
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
