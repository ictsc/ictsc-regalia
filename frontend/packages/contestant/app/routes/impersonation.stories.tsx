import type { Meta, StoryObj } from "@storybook/react";
import { ImpersonationPage } from "./impersonation.page";

export default {
  title: "pages/impersonation",
  component: ImpersonationPage,
} satisfies Meta<typeof ImpersonationPage>;

type Story = StoryObj<typeof ImpersonationPage>;

export const Default: Story = {
  name: "競技者選択",
  args: {
    children: (
      <p className="text-14 text-text">
        （フォームはSuspense内でレンダリング）
      </p>
    ),
  },
};
