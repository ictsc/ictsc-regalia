import type { Meta, StoryObj } from "@storybook/react";
import { IndexPage } from "./index.page";
import { Phase } from "@ictsc/proto/contestant/v1";

export default {
  title: "pages/index",
  component: IndexPage,
} satisfies Meta<typeof IndexPage>;

type Story = StoryObj<typeof IndexPage>;

export const InContest: Story = {
  name: "競技中",
  args: { phase: Phase.IN_CONTEST, nextPhase: Phase.IN_CONTEST },
};

export const OutOfContest: Story = {
  name: "競技時間外",
  args: { phase: Phase.OUT_OF_CONTEST },
};
