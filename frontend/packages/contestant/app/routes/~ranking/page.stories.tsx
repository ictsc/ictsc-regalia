import type { Meta, StoryObj } from "@storybook/react";
import { RankingPage } from "./page";

export default {
  title: "pages/ranking",
  component: RankingPage,
} satisfies Meta<typeof RankingPage>;

type Story = StoryObj<typeof RankingPage>;

export const Default: Story = {
  args: {
    ranking: [
      ...Array.from({ length: 2 }, () => ({
        $typeName: "contestant.v1.Rank" as const,
        rank: 1,
        teamName: "チーム名なまえがわからない",
        score: 8888,
        TimeStamp: "2025-03-04T12:00:00",
        organization: "testTeam",
      })),
      ...Array.from({ length: 4 }, () => ({
        $typeName: "contestant.v1.Rank" as const,
        rank: 2,
        teamName: "チーム名なまえがわからない",
        score: 8888,
        TimeStamp: "2025-03-04T12:00:00",
        organization: "testTeam",
      })),
    ],
  },
};
