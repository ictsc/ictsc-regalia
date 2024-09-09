import type { Meta, StoryObj } from "@storybook/react";
import { IndexPage } from "./page";

export default {
  title: "pages/index",
  component: IndexPage,
} satisfies Meta<typeof IndexPage>;

type Story = StoryObj<typeof IndexPage>;

export const InContest: Story = {
  name: "競技中",
  args: { inContest: true },
};

export const OutOfContest: Story = {
  name: "競技時間外",
  args: { inContest: false },
};
