import type { Meta, StoryObj } from "@storybook/react";
import { NotificationBanner } from "./notification-banner";

export default {
  title: "components/NotificationBanner",
  component: NotificationBanner,
} satisfies Meta<typeof NotificationBanner>;

type Story = StoryObj<typeof NotificationBanner>;

export const Default: Story = {
  render: () => <NotificationBanner message="次の再展開から減点されます！" />,
};
