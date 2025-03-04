import type { Meta, StoryObj } from "@storybook/react";
import { AnnounceList } from "./page";

export default {
  title: "pages/announces",
  component: AnnounceList,
} satisfies Meta<typeof AnnounceList>;

type Story = StoryObj<typeof AnnounceList>;

export const InContest: Story = {
  render: () => (
    <div>
      <AnnounceList announce="第二報・woaの初期コンフィグの設定間違いについて" />
      <AnnounceList announce="第二報・woaの初期コンフィグの設定間違いについて" />
      <AnnounceList announce="第二報・woaの初期コンフィグの設定間違いについて" />
      <AnnounceList announce="第二報・woaの初期コンフィグの設定間違いについて" />
      <AnnounceList announce="第二報・woaの初期コンフィグの設定間違いについて" />
      <AnnounceList announce="第二報・woaの初期コンフィグの設定間違いについて" />
      <AnnounceList announce="第二報・woaの初期コンフィグの設定間違いについて" />
    </div>
  ),
};
