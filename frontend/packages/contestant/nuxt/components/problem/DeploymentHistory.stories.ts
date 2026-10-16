import type { Meta, StoryObj } from "@storybook/vue3-vite";
import DeploymentHistory from "./DeploymentHistory.vue";
import { withContestantFrame, withDetailContent } from "~/stories/decorators";
import { makeDeployment } from "~/stories/fixtures";

const meta = {
  component: DeploymentHistory,
  title: "Contestant/ProblemDetail/DeploymentHistory",
  decorators: [withContestantFrame(), withDetailContent],
  parameters: { layout: "fullscreen", bodyClass: "detail-page" },
  args: {
    deployments: [],
    penaltyThreshold: 2,
    deploying: false,
    streamError: false,
  },
} satisfies Meta<typeof DeploymentHistory>;
export default meta;
type Story = StoryObj<typeof meta>;

export const None: Story = {};

export const Requesting: Story = {
  args: { deploying: true },
};

export const InProgress: Story = {
  args: {
    deployments: [
      makeDeployment(2, "DEPLOYING"),
      makeDeployment(1, "COMPLETED"),
    ],
  },
};

export const Failed: Story = {
  args: { deployments: [makeDeployment(1, "FAILED")] },
};

export const OverThreshold: Story = {
  args: {
    deployments: [5, 4, 3, 2, 1].map((revision) =>
      makeDeployment(revision, "COMPLETED"),
    ),
  },
};

export const Reconnecting: Story = {
  args: {
    deployments: [makeDeployment(1, "QUEUED")],
    streamError: true,
  },
};
