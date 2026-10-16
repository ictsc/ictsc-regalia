import type { Meta, StoryObj } from "@storybook/vue3-vite";
import CompetitionHeaderView from "~/components/header/CompetitionHeaderView.vue";
import ProblemListView from "~/components/problem/ProblemListView.vue";
import { headerRankingItems } from "~/features/ranking";
import { problemStatus } from "~/features/problem/status";
import { contestClock } from "~/features/schedule/clock";
import { withContestantFrame } from "~/stories/decorators";
import {
  fixedNow,
  ownTeamCode,
  problems,
  ranking,
  sections,
} from "~/stories/fixtures";

const header = {
  compact: false,
  rankingItems: headerRankingItems(ranking, ownTeamCode),
  ownTeamCode,
  maxScore: problems.reduce((sum, p) => sum + p.maxScore, 0),
  clock: contestClock(Object.values(sections), fixedNow),
  demoMode: false,
  hasUnreadNotices: true,
  teamName: "KERNEL PANIC",
  routeKey: "/problems",
};
const items = problems.map((problem) => ({
  problem,
  status: problemStatus(problem),
  cooldownMinutes: problem.code === "A02" ? 12 : 0,
}));

const meta = {
  title: "Contestant/Screens/ProblemList",
  decorators: [withContestantFrame()],
  parameters: { layout: "fullscreen", route: "/problems" },
  render: () => ({
    components: { CompetitionHeaderView, ProblemListView },
    setup: () => ({ header, items }),
    template: `<CompetitionHeaderView v-bind="header" /><main><ProblemListView :items="items" /></main>`,
  }),
} satisfies Meta;
export default meta;
type Story = StoryObj<typeof meta>;

export const Desktop: Story = {};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
