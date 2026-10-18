import type { Meta, StoryObj } from "@storybook/vue3-vite";
import SettingsView from "./SettingsView.vue";
import { withAdminShell } from "~/stories/decorators";
import { ruleMarkdown } from "~/stories/fixtures";

const meta = {
  component: SettingsView,
  title: "Admin/Settings/SettingsView",
  decorators: [withAdminShell("/settings")],
  parameters: { layout: "fullscreen" },
  args: {
    initialMarkdown: ruleMarkdown,
    initialFreeze: "2026-09-02T13:13",
    busy: false,
    message: "",
  },
} satisfies Meta<typeof SettingsView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const NoFreeze: Story = {
  args: { initialFreeze: "" },
};

export const EmptyRule: Story = {
  args: { initialMarkdown: "" },
};

export const Saved: Story = {
  args: { message: "ルール更新しました" },
};
