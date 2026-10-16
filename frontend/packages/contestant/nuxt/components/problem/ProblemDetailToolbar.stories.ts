import type { Meta, StoryObj } from "@storybook/vue3-vite";
import { ref } from "vue";
import ProblemDetailToolbar from "./ProblemDetailToolbar.vue";
import { withContestantFrame } from "~/stories/decorators";

const meta = {
  component: ProblemDetailToolbar,
  title: "Contestant/ProblemDetail/ProblemDetailToolbar",
  decorators: [withContestantFrame()],
  parameters: { layout: "fullscreen", bodyClass: "detail-page" },
  args: { railOpen: true, resetDisabled: false },
  render: (args) => ({
    components: { ProblemDetailToolbar },
    setup: () => ({ args, railOpen: ref(args.railOpen) }),
    template: `<ProblemDetailToolbar v-bind="args" v-model:rail-open="railOpen" />`,
  }),
} satisfies Meta<typeof ProblemDetailToolbar>;
export default meta;
type Story = StoryObj<typeof meta>;

export const RailOpen: Story = {};

export const RailHidden: Story = {
  args: { railOpen: false },
};

export const ResetUnavailable: Story = {
  args: { resetDisabled: true },
};
