import type { Meta, StoryObj } from "@storybook/vue3-vite";
import WebPushSettings from "./WebPushSettings.vue";
import { withContestantFrame } from "~/stories/decorators";

const meta = {
  component: WebPushSettings,
  title: "Contestant/Announces/WebPushSettings",
  decorators: [
    withContestantFrame(),
    () => ({ template: "<main><story /></main>" }),
  ],
  parameters: { layout: "fullscreen" },
  args: {
    subscribed: false,
    supported: true,
    serverEnabled: true,
    busy: false,
  },
} satisfies Meta<typeof WebPushSettings>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Available: Story = {};

export const Subscribed: Story = {
  args: { subscribed: true },
};

export const Unsupported: Story = {
  args: { supported: false },
};

export const ServerDisabled: Story = {
  args: { serverEnabled: false },
};

export const Busy: Story = {
  args: { busy: true },
};

export const Failed: Story = {
  args: { error: new Error("ブラウザで通知が許可されていません") },
};
