import type { Meta, StoryObj } from "@storybook/vue3-vite";
import AdminHeader from "./AdminHeader.vue";
import { adminName } from "~/stories/fixtures";

const meta = {
  component: AdminHeader,
  title: "Admin/Layout/AdminHeader",
  decorators: [() => ({ template: `<div class="admin-app"><story /></div>` })],
  parameters: { layout: "fullscreen" },
  args: { signedIn: true, displayName: adminName },
} satisfies Meta<typeof AdminHeader>;
export default meta;
type Story = StoryObj<typeof meta>;

export const SignedIn: Story = {};

export const Anonymous: Story = {
  args: { signedIn: false, displayName: undefined },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
