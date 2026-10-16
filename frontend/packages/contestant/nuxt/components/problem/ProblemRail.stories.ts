import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ProblemRail, { type ProblemRailItem } from "./ProblemRail.vue";
import { withContestantFrame } from "~/stories/decorators";
import { problems } from "~/stories/fixtures";

const items = (
  overrides: Record<string, Partial<ProblemRailItem>> = {},
): ProblemRailItem[] =>
  problems.map((problem) => ({
    problem,
    cooldownMinutes: 0,
    pending: false,
    ...overrides[problem.code],
  }));

const meta = {
  component: ProblemRail,
  title: "Contestant/ProblemDetail/ProblemRail",
  decorators: [
    withContestantFrame(),
    () => ({ template: `<div class="detail-layout"><story /></div>` }),
  ],
  parameters: {
    layout: "fullscreen",
    bodyClass: "detail-page",
    route: "/problems/A04",
  },
  args: { items: items(), currentCode: "A04" },
} satisfies Meta<typeof ProblemRail>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const CooldownAndPending: Story = {
  args: {
    items: items({
      A04: { cooldownMinutes: 17 },
      B03: { cooldownMinutes: 2 },
      B04: { pending: true },
    }),
  },
};

export const ClosedProblemSelected: Story = {
  args: { currentCode: "A01" },
  parameters: { route: "/problems/A01" },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
