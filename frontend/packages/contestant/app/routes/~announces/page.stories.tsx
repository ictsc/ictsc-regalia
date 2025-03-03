import type { Meta, StoryObj } from "@storybook/react";
import { NotificationList } from "./page";

export default {
  title: "pages/announces",
  component: NotificationList,
} satisfies Meta<typeof NotificationList>;

type Story = StoryObj<typeof NotificationList>;

export const InContest: Story = {
  name: "Default",
  args: {
    notifications: [
      { text: '第二報・woaの初期コンフィグの設定間違いについて'},
      { text: '第二報・woaの初期コンフィグの設定間違いについて'},
      { text: '第二報・woaの初期コンフィグの設定間違いについて'},
      { text: '第二報・woaの初期コンフィグの設定間違いについて'},
      { text: '第二報・woaの初期コンフィグの設定間違いについて'},
      { text: '第二報・woaの初期コンフィグの設定間違いについて'},
      { text: '第二報・woaの初期コンフィグの設定間違いについて'},
      { text: '第二報・woaの初期コンフィグの設定間違いについて'},
    ],
  },
};