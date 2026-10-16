import type { Meta, StoryObj } from "@storybook/vue3-vite";
import { ref } from "vue";
import ResetDialog from "./ResetDialog.vue";

const meta = {
  component: ResetDialog,
  title: "Contestant/ProblemDetail/ResetDialog",
  parameters: { layout: "fullscreen", bodyClass: "detail-page" },
  args: { penaltyThreshold: 2 },
  render: (args) => ({
    components: { ResetDialog },
    setup: () => ({ args, open: ref(true) }),
    template: `<button class="button-secondary" @click="open = true">ダイアログを開く</button><ResetDialog v-bind="args" v-model:open="open" @confirm="open = false" />`,
  }),
} satisfies Meta<typeof ResetDialog>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Open: Story = {};
