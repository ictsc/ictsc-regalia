import type { Meta, StoryObj } from "@storybook/react";
import { ConfirmModal } from "./confirmModal";

export default {
  title: "pages/problem/confirmModal",
} satisfies Meta;

type Story = StoryObj;

export const Default: Story = {
  render: () => (
    <div className="grid grid-cols-2 gap-64">
      <ConfirmModal isOpen={true} onConfirm={() => {}} onCansel={() => {}} allowedDeploymentCount={1} />
    </div>
  ),
};

export const NoMoreLeft: Story = {
  render: () => (
    <div className="grid grid-cols-2 gap-64">
      <ConfirmModal isOpen={true} onConfirm={() => {}} onCansel={() => {}} allowedDeploymentCount={0} />
    </div>
  ),
};
