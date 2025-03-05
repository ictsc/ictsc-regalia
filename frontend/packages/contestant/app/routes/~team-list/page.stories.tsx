import type { Meta, StoryObj } from "@storybook/react";
import {  TeamList } from "./page";

export default {
  title: "pages/team-list",
} satisfies Meta;

type Story = StoryObj;

export const Default: Story = {
  render: () => (
    <div className="w-full gap-64">
      <TeamList
        teamName="チーム名"
        affiliation="testTeam"
        userNames={["hoge","huga"]}
      />
      <TeamList
        teamName="チーム名がわからないいいいいいいいいいいいい"
        affiliation="testTeam"
        userNames={["hoge","huga"]}
      />
      <TeamList

        teamName="チーム名がわからないいいいいいいいいいいいい"
        affiliation="testTeam"
        userNames={["hoge","huga"]}
      />
      <TeamList
        teamName="チーム名がわからないいいいいいいいいいいいい"
        affiliation="testTeam"
        userNames={["hoge","huga"]}
      />
      <TeamList
        teamName="チーム名がわからないいいいいいいいいいいいい"
        affiliation="testTeam"
        userNames={["hoge","huga"]}
      />
      <TeamList
        teamName="チーム名がわからないいいいいいいいいいいいい"
        affiliation="testTeam"
        userNames={["hoge","huga"]}
      />
    </div>
  ),
};
