import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ScheduleView from "./ScheduleView.vue";
import { withContestantFrame } from "~/stories/decorators";
import { fixedNow, sections } from "~/stories/fixtures";

const meta = {
  component: ScheduleView,
  title: "Contestant/Schedule/ScheduleView",
  decorators: [
    withContestantFrame(),
    () => ({ template: "<main><story /></main>" }),
  ],
  parameters: { layout: "fullscreen", route: "/schedule" },
  args: {
    entries: [sections.day2, sections.day1],
    now: fixedNow,
  },
} satisfies Meta<typeof ScheduleView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const Empty: Story = {
  args: { entries: [] },
};

export const Error: Story = {
  args: { entries: undefined, error: "スケジュールを取得できませんでした" },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
