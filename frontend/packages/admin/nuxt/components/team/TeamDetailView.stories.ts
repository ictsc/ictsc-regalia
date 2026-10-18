import type { Meta, StoryObj } from "@storybook/vue3-vite";
import TeamDetailView from "./TeamDetailView.vue";
import { withAdminShell } from "~/stories/decorators";
import { contestants, invitations, teams } from "~/stories/fixtures";

const team = teams[2]!;

const meta = {
  component: TeamDetailView,
  title: "Admin/Teams/TeamDetailView",
  decorators: [withAdminShell(`/teams/${team.code}`)],
  parameters: { layout: "fullscreen" },
  args: {
    code: team.code,
    team,
    members: contestants.filter((c) => c.team.code === team.code),
    invitations,
    busy: false,
    message: "",
  },
} satisfies Meta<typeof TeamDetailView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const NoMembersOrInvitations: Story = {
  args: { members: [], invitations: [] },
};

export const Saved: Story = {
  args: { message: "チームを更新しました" },
};

export const Busy: Story = {
  args: { busy: true },
};
