import type { Meta, StoryObj } from "@storybook/vue3-vite";
import DeploymentsView from "./DeploymentsView.vue";
import { withAdminShell } from "~/stories/decorators";
import { catalog, deployments, teams } from "~/stories/fixtures";

const meta = {
  component: DeploymentsView,
  title: "Admin/Deployments/DeploymentsView",
  decorators: [withAdminShell("/deployments")],
  parameters: { layout: "fullscreen" },
  args: {
    deployments,
    teams,
    problems: catalog,
    streamError: false,
    busy: false,
    message: "",
  },
} satisfies Meta<typeof DeploymentsView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const Requesting: Story = {
  args: { busy: true },
};

export const Reconnecting: Story = {
  args: { streamError: true },
};

export const Empty: Story = {
  args: { deployments: [] },
};
