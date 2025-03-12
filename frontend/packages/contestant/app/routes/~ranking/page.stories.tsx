import type { Meta, StoryObj } from "@storybook/react";
import { timestampFromDate } from "@bufbuild/protobuf/wkt";
import { RankingPage } from "./page";

const date = new Date("2025-03-04T12:00:00Z");
const timestamp = timestampFromDate(date);

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
        rank: 1n,
        teamName: "チーム名なまえがわからない",
        score: 8888n,
        timestamp: timestamp,
        organization: "testTeam",
      })),
      ...Array.from({ length: 4 }, () => ({
        $typeName: "contestant.v1.Rank" as const,
        rank: 2n,
        teamName: "チーム名なまえがわからない",
        score: 8888n,
        timestamp: timestamp,
        organization: "testTeam",
      })),
    ],
  },
};
