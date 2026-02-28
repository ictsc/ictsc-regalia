import type { Meta, StoryObj } from "@storybook/react";
import { IndexPage } from "./index.page";

export default {
  title: "pages/index",
  component: IndexPage,
} satisfies Meta<typeof IndexPage>;

type Story = StoryObj<typeof IndexPage>;

export const InContest: Story = {
  name: "競技中",
  args: {
    state: "in_contest",
    currentScheduleName: "day1-am",
    nextScheduleName: "day1-pm",
  },
};

export const Waiting: Story = {
  name: "競技時間外（待機中）",
  args: {
    state: "waiting",
    nextScheduleName: "day1-am",
  },
};

export const Ended: Story = {
  name: "競技終了",
  args: {
    state: "ended",
  },
};
