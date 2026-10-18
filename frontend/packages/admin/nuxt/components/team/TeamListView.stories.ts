import type { Meta, StoryObj } from "@storybook/vue3-vite";
import TeamListView from "./TeamListView.vue";
import { withAdminShell } from "~/stories/decorators";
import { teams } from "~/stories/fixtures";

const meta = {
  component: TeamListView,
  title: "Admin/Teams/TeamListView",
  decorators: [withAdminShell("/teams")],
  parameters: { layout: "fullscreen" },
  args: { teams, busy: false, message: "" },
} satisfies Meta<typeof TeamListView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const Created: Story = {
  args: { message: "チームを作成しました" },
};

export const Busy: Story = {
  args: { busy: true },
};

export const Empty: Story = {
  args: { teams: [] },
};

export const Error: Story = {
  args: { teams: undefined, error: "チーム一覧を取得できませんでした" },
};
