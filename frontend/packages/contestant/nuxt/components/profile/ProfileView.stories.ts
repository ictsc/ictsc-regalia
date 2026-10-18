import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ProfileView from "./ProfileView.vue";
import { withContestantFrame } from "~/stories/decorators";
import { profile } from "~/stories/fixtures";

const meta = {
  component: ProfileView,
  title: "Contestant/Profile/ProfileView",
  decorators: [
    withContestantFrame(),
    () => ({ template: `<main class="workspace"><story /></main>` }),
  ],
  parameters: { layout: "fullscreen", route: "/profile" },
  args: {
    profile,
    saving: false,
    message: "",
    signingOut: false,
    signOutError: "",
  },
} satisfies Meta<typeof ProfileView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const Saving: Story = {
  args: { saving: true },
};

export const Saved: Story = {
  args: { message: "プロフィールを更新しました" },
};

export const SigningOut: Story = {
  args: { signingOut: true },
};

export const SignOutFailed: Story = {
  args: { signOutError: "ログアウトできませんでした" },
};

export const LoadError: Story = {
  args: { profile: undefined, error: "プロフィールを取得できませんでした" },
};
