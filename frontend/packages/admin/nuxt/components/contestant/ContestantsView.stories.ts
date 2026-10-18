import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ContestantsView from "./ContestantsView.vue";
import { withAdminShell } from "~/stories/decorators";
import { contestants } from "~/stories/fixtures";

const meta = {
  component: ContestantsView,
  title: "Admin/Contestants/ContestantsView",
  decorators: [withAdminShell("/contestants")],
  parameters: { layout: "fullscreen" },
  args: { contestants, busy: false, message: "" },
} satisfies Meta<typeof ContestantsView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const Busy: Story = {
  args: { busy: true, message: "" },
};

export const ImpersonationFailed: Story = {
  args: { message: "代理ログインできませんでした" },
};

export const Empty: Story = {
  args: { contestants: [] },
};

export const Loading: Story = {
  args: { contestants: undefined, pending: true },
};
