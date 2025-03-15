import type { Meta, StoryObj } from "@storybook/react";
import { ConfirmModal } from "./confirmModal";

export default {
  title: "pages/problem/confirmModal",
} satisfies Meta;

type Story = StoryObj;

export const SubmitDeployment: Story = {
  render: () => (
    <ConfirmModal
      isOpen={true}
      onConfirm={() => {}}
      onCancel={() => {}}
      title="再展開の確認"
      confirmText="再展開する"
      cancelText="キャンセル"
    >
      <span>ここに本文を書く</span>
    </ConfirmModal>
  ),
};

export const SubmitAnswer: Story = {
  render: () => (
    <ConfirmModal
      isOpen={true}
      onConfirm={() => {}}
      onCancel={() => {}}
      title="解答の確認"
      confirmText="送信する"
      cancelText="キャンセル"
    >
      <div className="my-12">
        <p className="text-16 text-text">本当にこの問題を提出しますか？</p>
      </div>
    </ConfirmModal>
  ),
};
