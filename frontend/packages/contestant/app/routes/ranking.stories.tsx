import type { Meta, StoryObj } from "@storybook/react";
import { RankingPage } from "./ranking.page";

export default {
  title: "pages/ranking",
  component: RankingPage,
} satisfies Meta<typeof RankingPage>;

type Story = StoryObj<typeof RankingPage>;

export const Default: Story = {
  args: {
    ranking: [
      ...Array.from({ length: 2 }, () => ({
        rank: 1,
        teamName: "チーム名なまえがわからない",
        organization: "testTeam",
        score: 8888,
        lastEffectiveSubmitAt: "2025-03-04T12:00:00Z",
      })),
      ...Array.from({ length: 4 }, () => ({
        rank: 2,
        teamName: "チーム名なまえがわからない",
        organization: "testTeam",
        score: 8888,
        lastEffectiveSubmitAt: "2025-03-04T12:00:00Z",
      })),
    ],
  },
};
