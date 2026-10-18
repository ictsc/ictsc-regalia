import type { Meta, StoryObj } from "@storybook/vue3-vite";
import SignupView from "./SignupView.vue";
import { withContestantFrame } from "~/stories/decorators";

const meta = {
  component: SignupView,
  title: "Contestant/Signup/SignupView",
  decorators: [
    withContestantFrame(),
    () => ({ template: `<main class="workspace"><story /></main>` }),
  ],
  parameters: { layout: "fullscreen", route: "/signup" },
  args: {
    initialName: "alice",
    initialDisplayName: "Alice",
    pending: false,
    message: "",
  },
} satisfies Meta<typeof SignupView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const InvitationCode: Story = {};

export const AssignedTeam: Story = {
  args: { registrationTeamName: "KERNEL PANIC" },
};

export const Submitting: Story = {
  args: { pending: true },
};

export const Rejected: Story = {
  args: { message: "登録できませんでした: 招待コードが正しくありません" },
};
