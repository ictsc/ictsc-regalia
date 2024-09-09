import type { Meta, StoryObj } from "@storybook/react";
import { ContestStateView, DAY, HOUR, MINUTE } from "./contest-state";

export default {
  title: "ContentState",
  component: ContestStateView,
} satisfies Meta<typeof ContestStateView>;

type Story = StoryObj<typeof ContestStateView>;

export const All: Story = {
  render: () => (
    <div className="flex flex-col items-start gap-[20px]">
      <ContestStateView
        state="break"
        restDurationSeconds={15 * HOUR + 30 * MINUTE + 50}
      />
      <ContestStateView
        state="running"
        restDurationSeconds={3 * HOUR + 56 * MINUTE + 24}
      />
      <ContestStateView
        state="before"
        restDurationSeconds={20 * HOUR + 30 * MINUTE + 50}
      />
      <ContestStateView state="before" restDurationSeconds={6 * DAY + 50} />
      <ContestStateView state="finished" restDurationSeconds={0} />
    </div>
  ),
};
