import type { Meta, StoryObj } from "@storybook/vue3-vite";
import RequestState from "./RequestState.vue";

const meta = {
  component: RequestState,
  title: "UI/Feedback/RequestState",
} satisfies Meta<typeof RequestState>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Loading: Story = {
  args: { pending: true },
};
export const Error: Story = {
  args: { error: "接続できませんでした" },
};
export const Loaded: Story = {
  render: (args) => ({
    components: { RequestState },
    setup: () => ({ args }),
    template: `<RequestState v-bind="args"><p>読み込んだ内容を表示します。</p></RequestState>`,
  }),
};
