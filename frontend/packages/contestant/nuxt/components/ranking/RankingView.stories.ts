import type { Meta, StoryObj } from "@storybook/vue3-vite";
import RankingView from "./RankingView.vue";
import { withContestantFrame } from "~/stories/decorators";
import { ownTeamCode, ranking } from "~/stories/fixtures";

const meta = {
  component: RankingView,
  title: "Contestant/Ranking/RankingView",
  decorators: [
    withContestantFrame(),
    () => ({ template: "<main><story /></main>" }),
  ],
  parameters: { layout: "fullscreen", route: "/ranking" },
  args: { snapshot: { ranking, frozen: false, hidden: false }, ownTeamCode },
} satisfies Meta<typeof RankingView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const Frozen: Story = {
  args: {
    snapshot: {
      ranking,
      frozen: true,
      frozenAt: "2026-09-02T13:13:48+09:00",
      hidden: false,
    },
  },
};

export const Tie: Story = {
  args: {
    snapshot: {
      ranking: ranking.map((entry) =>
        entry.teamCode === 13 ? { ...entry, rank: 12, score: 1230 } : entry,
      ),
      frozen: false,
      hidden: false,
    },
  },
};

export const LongNames: Story = {
  args: {
    snapshot: {
      ranking: ranking.map((entry, index) =>
        index === 2
          ? {
              ...entry,
              teamName: "KERNEL PANIC AND THE INFINITE RETRANSMISSION",
              organization: "国立ネットワーク工学技術大学院大学 情報通信研究会",
            }
          : entry,
      ),
      frozen: false,
      hidden: false,
    },
  },
};

export const Empty: Story = {
  args: { snapshot: { ranking: [], frozen: false, hidden: false } },
};

export const Loading: Story = {
  args: { snapshot: undefined, pending: true },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};

/** 二次予選モードで最終公開前。順位は返らない。 */
export const Hidden: Story = {
  args: { snapshot: { ranking: [], frozen: false, hidden: true } },
};
