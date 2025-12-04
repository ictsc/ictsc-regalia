import type { Meta, StoryObj } from "@storybook/react";
import { TeamsPage } from "./teams.page";

export default {
  title: "pages/teams",
  component: TeamsPage,
} satisfies Meta;

type Story = StoryObj;

export const Default: Story = {
  args: {
    teamProfile: [
      ...Array.from({ length: 1 }, () => ({
        $typeName: "contestant.v1.TeamProfile" as const,
        name: "チーム名がわからない",
        organization: "テスト",
        members: [
          {
            name: "huge",
            displayName: "ふが",
            selfIntroduction: "よろしくお願いします！",
          },
          {
            name: "hoge",
            displayName: "ほげ",
            selfIntroduction: "がんばります！",
          },
        ],
      })),
      ...Array.from({ length: 2 }, () => ({
        $typeName: "contestant.v1.TeamProfile" as const,
        name: "チーム名",
        organization: "testTeam",
        members: [
          {
            name: "huge",
            displayName: "ふが",
            selfIntroduction: "よろしくお願いします！",
          },
          {
            name: "hoge",
            displayName: "ほげ",
            selfIntroduction: "がんばります！",
          },
        ],
      })),
    ],
  },
};
