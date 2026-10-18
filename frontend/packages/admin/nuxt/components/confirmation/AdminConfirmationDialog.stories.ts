import type { Meta, StoryObj } from "@storybook/vue3-vite";
import AdminConfirmationDialog from "./AdminConfirmationDialog.vue";

const meta = {
  component: AdminConfirmationDialog,
  title: "Admin/Confirmation/AdminConfirmationDialog",
  parameters: { layout: "fullscreen" },
  args: {
    request: {
      title: "得点再計算",
      message: "すべての得点を再計算しますか？",
      confirmLabel: "得点再計算",
    },
  },
} satisfies Meta<typeof AdminConfirmationDialog>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Normal: Story = {};

export const Destructive: Story = {
  args: {
    request: {
      title: "チームを削除",
      message: "KERNEL PANIC を削除しますか？この操作は取り消せません。",
      confirmLabel: "チームを削除",
      destructive: true,
    },
  },
};
