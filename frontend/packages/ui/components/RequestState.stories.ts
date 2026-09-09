import type { Meta, StoryObj } from "@storybook/vue3-vite";
import RequestState from "./RequestState.vue";
export default {
  component: RequestState,
  title: "Feedback/RequestState",
} satisfies Meta<typeof RequestState>;
export const Loading: StoryObj<typeof RequestState> = {
  args: { pending: true },
};
export const Error: StoryObj<typeof RequestState> = {
  args: { error: "接続できませんでした" },
};
