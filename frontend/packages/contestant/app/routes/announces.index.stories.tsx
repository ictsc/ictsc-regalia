import type { Meta, StoryObj } from "@storybook/react";
import { AnnounceList } from "./announces.index.page";

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
        slug: "hogehoge",
        title: "hogehoge",
        body: "hogehoge",
      })),
    ],
  },
};
