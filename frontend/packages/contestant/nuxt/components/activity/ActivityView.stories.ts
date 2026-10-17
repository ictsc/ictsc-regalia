import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ActivityView from "./ActivityView.vue";
import { withContestantFrame } from "~/stories/decorators";
import { activity } from "~/stories/fixtures";

const meta = {
  component: ActivityView,
  title: "Contestant/Activity/ActivityView",
  decorators: [
    withContestantFrame(),
    () => ({ template: "<main><story /></main>" }),
  ],
  parameters: { layout: "fullscreen", route: "/activity" },
  args: { entries: activity },
} satisfies Meta<typeof ActivityView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const Empty: Story = {
  args: { entries: [] },
};

export const Error: Story = {
  args: { entries: undefined, error: "提出履歴を取得できませんでした" },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
