import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ScoresView from "./ScoresView.vue";
import { withAdminShell } from "~/stories/decorators";
import { ranking, scores } from "~/stories/fixtures";

const meta = {
  component: ScoresView,
  title: "Admin/Scores/ScoresView",
  decorators: [withAdminShell("/scores")],
  parameters: { layout: "fullscreen" },
  args: { ranking, scores, busy: false, message: "" },
} satisfies Meta<typeof ScoresView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Public: Story = {};

export const Frozen: Story = {
  args: {
    ranking: {
      ...ranking,
      frozen: true,
      frozen_at: "2026-09-02T13:13:48+09:00",
    },
  },
};

export const Recalculated: Story = {
  args: { message: "得点再計算しました" },
};

export const Empty: Story = {
  args: {
    ranking: { ranking: [], frozen: false, frozen_at: null, hidden: false },
    scores: [],
  },
};
