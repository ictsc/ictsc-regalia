import type { Meta, StoryObj } from "@storybook/vue3-vite";
import CompetitionHeaderView from "./CompetitionHeaderView.vue";
import { headerRankingItems } from "~/features/ranking";
import { contestClock } from "~/features/schedule/clock";
import { withContestantFrame } from "~/stories/decorators";
import {
  fixedNow,
  ownTeamCode,
  problems,
  ranking,
  sections,
} from "~/stories/fixtures";

const entries = Object.values(sections);
const maxScore = problems.reduce((sum, p) => sum + p.maxScore, 0);

const meta = {
  component: CompetitionHeaderView,
  title: "Contestant/Header/CompetitionHeaderView",
  decorators: [withContestantFrame()],
  parameters: { layout: "fullscreen", route: "/problems" },
  args: {
    compact: false,
    rankingItems: headerRankingItems(ranking, ownTeamCode),
    ownTeamCode,
    maxScore,
    clock: contestClock(entries, fixedNow),
    demoMode: false,
    hasUnreadNotices: false,
    teamName: "KERNEL PANIC",
    routeKey: "/problems",
  },
} satisfies Meta<typeof CompetitionHeaderView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const MiddleRank: Story = {};

export const FirstPlace: Story = {
  args: {
    rankingItems: headerRankingItems(ranking, ranking[0]!.teamCode),
    ownTeamCode: ranking[0]!.teamCode,
  },
};

export const LastPlace: Story = {
  args: {
    rankingItems: headerRankingItems(ranking, ranking.at(-1)!.teamCode),
    ownTeamCode: ranking.at(-1)!.teamCode,
  },
};

const tiedRanking = ranking.map((entry) =>
  entry.teamCode === 13 ? { ...entry, rank: 12, score: 1230 } : entry,
);
export const Tie: Story = {
  args: { rankingItems: headerRankingItems(tiedRanking, ownTeamCode) },
};

export const NoOwnRank: Story = {
  args: {
    rankingItems: headerRankingItems(ranking, undefined),
    ownTeamCode: undefined,
  },
};

export const UnreadNotices: Story = {
  args: { hasUnreadNotices: true },
};

export const Compact: Story = {
  args: { compact: true, routeKey: "/problems/A04" },
  parameters: { route: "/problems/A04" },
};

export const DemoMode: Story = {
  args: { demoMode: true, clock: { label: "デモモード", text: "02:00:00" } },
};

export const BeforeStart: Story = {
  args: {
    clock: contestClock(entries, Date.parse("2026-08-31T21:00:00+09:00")),
  },
};

export const ContestFinished: Story = {
  args: {
    clock: contestClock(entries, Date.parse("2026-09-03T09:00:00+09:00")),
  },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
