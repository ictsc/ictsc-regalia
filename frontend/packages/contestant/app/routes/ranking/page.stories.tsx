import type { Meta, StoryObj } from "@storybook/react";
import { RankingItem } from "./page";

export default {
  title: "pages/ranking",
} satisfies Meta;

type Story = StoryObj;

export const Default: Story = {
  render: () => (
    <div className="gap-64 w-full">
      <RankingItem
        rank={1}
        teamName="チーム名なまえがわからない"
        affiliation="testTeam"
        score={8888}
        timeStamp={new Date('2025-03-04T12:00:00')}
        penalty={0}
        fullScore
        rawFullScore
      />
       <RankingItem
        rank={2}
        teamName="チーム名がわからないいいいいいいいいいいいい"
        affiliation="testTeam"
        score={8888}
        timeStamp={new Date('2025-03-04T12:00:00')}
        penalty={0}
        fullScore
        rawFullScore
      />
       <RankingItem
        rank={3}
        teamName="チーム名がわからないいいいいいいいいいいいい"
        affiliation="testTeam"
        score={8888}
        timeStamp={new Date('2025-03-04T12:00:00')}
        penalty={0}
        fullScore
        rawFullScore
      />
       <RankingItem
        rank={4}
        teamName="チーム名がわからないいいいいいいいいいいいい"
       affiliation="testTeam"
        score={8888}
        timeStamp={new Date('2025-03-04T12:00:00')}
        penalty={0}
        fullScore
        rawFullScore
      />
       <RankingItem
        rank={5}
        teamName="チーム名がわからないいいいいいいいいいいいい"
         affiliation="testTeam"
        score={500}
        timeStamp={new Date('2025-03-04T12:00:00')}
        penalty={0}
        fullScore
        rawFullScore
      />
       <RankingItem
        rank={6}
        teamName="チーム名がわからないいいいいいいいいいいいい"
        affiliation="testTeam"
        score={0}
        timeStamp={new Date('2025-03-04T12:00:00')}
        penalty={0}
        fullScore
        rawFullScore
      />
    </div>
  ),
};