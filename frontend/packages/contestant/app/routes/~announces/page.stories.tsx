import type { Meta, StoryObj } from "@storybook/react";
import { AnnounceList } from "./page";

export default {
  title: "pages/announces",
  component: AnnounceList,
} satisfies Meta<typeof AnnounceList>;

type Story = StoryObj<typeof AnnounceList>;

export const InContest: Story = {
  args: {
    announces: [
      ...Array.from({ length: 10 }, () => ({
        $typeName: "contestant.v1.Notice" as const,
        title: "hogehoge",
        body: "hogehoge",
      })),
    ],
  },
};
