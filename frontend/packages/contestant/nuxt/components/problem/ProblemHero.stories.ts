import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ProblemHero from "./ProblemHero.vue";
import { withContestantFrame } from "~/stories/decorators";

const meta = {
  component: ProblemHero,
  title: "Contestant/ProblemDetail/ProblemHero",
  decorators: [
    withContestantFrame(),
    () => ({
      template: `<div class="detail-layout is-rail-hidden"><article class="problem-content"><story /></article></div>`,
    }),
  ],
  parameters: { layout: "fullscreen", bodyClass: "detail-page" },
  args: {
    code: "B02",
    title: "コンテナが起動を繰り返す",
    category: "サーバー",
    maxScore: 300,
    cooldownMinutes: 0,
    open: true,
  },
} satisfies Meta<typeof ProblemHero>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Open: Story = {};

export const Scored: Story = {
  args: {
    score: { markedScore: 180, penalty: 0, score: 180, maxScore: 300 },
  },
};

export const Cooldown: Story = {
  args: { cooldownMinutes: 18 },
};

export const Closed: Story = {
  args: { open: false },
};

export const LongTitle: Story = {
  args: {
    title:
      "複数拠点を結ぶIPsec VPNのうち一部の拠点間でのみ通信が断続的に途切れ、経路の再計算が繰り返される",
  },
};
