import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ProblemContentView from "./ProblemContentView.vue";
import { withAdminShell } from "~/stories/decorators";
import { adminProblem } from "~/stories/fixtures";

const meta = {
  component: ProblemContentView,
  title: "Admin/Content/ProblemContentView",
  decorators: [withAdminShell("/content/problems/B02")],
  parameters: { layout: "fullscreen" },
  args: { problem: adminProblem, initialCommit: "" },
} satisfies Meta<typeof ProblemContentView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Active: Story = {};

export const SpecificCommit: Story = {
  args: {
    initialCommit: "1a2b3c4",
    problem: {
      ...adminProblem,
      default_grader_discord_id: "110000000000000001",
      redeploy_rule: {
        type: "MANUAL",
        penalty_threshold: null,
        penalty_percentage: null,
      },
    },
  },
};

export const NotFound: Story = {
  args: { problem: undefined, error: "指定した版に問題がありません" },
};
