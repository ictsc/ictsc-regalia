import type { Meta, StoryObj } from "@storybook/vue3-vite";
import { ref } from "vue";
import GradingDetailView from "./GradingDetailView.vue";
import { defaultPaneLayout } from "~/features/grading-panes";
import { withAdminShell } from "~/stories/decorators";
import {
  adminProblem,
  contentStatus,
  detailAnswer,
  markings,
  previousAnswers,
} from "~/stories/fixtures";

const meta = {
  component: GradingDetailView,
  title: "Admin/Grading/GradingDetailView",
  decorators: [
    withAdminShell("/submissions", "workspace grading-detail-workspace"),
  ],
  parameters: { layout: "fullscreen" },
  args: {
    answer: detailAnswer,
    problem: adminProblem,
    activeCommit: contentStatus.active_commit!,
    previous: previousAnswers,
    marks: markings,
    busy: false,
    workflowBusy: false,
    collapsed: defaultPaneLayout().collapsed,
    widths: defaultPaneLayout().widths,
    submitMarking: async () => true,
  },
  // Keep pane collapse and resize interactive in the preview.
  render: (args) => ({
    components: { GradingDetailView },
    setup: () => ({
      args,
      collapsed: ref(args.collapsed),
      widths: ref(args.widths),
    }),
    template: `<GradingDetailView v-bind="args" v-model:collapsed="collapsed" v-model:widths="widths" />`,
  }),
} satisfies Meta<typeof GradingDetailView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const WithPreviousAnswers: Story = {};

export const FirstAnswer: Story = {
  args: { previous: [], marks: [] },
};

export const OutdatedProblemVersion: Story = {
  args: { activeCommit: "0000000000000000000000000000000000000000" },
};

export const CollapsedProblemPane: Story = {
  args: { collapsed: [true, false, false], widths: [23, 55, 22] },
};

export const Saving: Story = {
  args: { busy: true, workflowBusy: true },
};

export const DefaultAssignment: Story = {
  args: {
    answer: {
      ...detailAnswer,
      score: null,
      workflow: {
        ...detailAnswer.workflow,
        status: "WAITING",
        assignment_source: "DEFAULT",
        assignee_discord_id: null,
      },
    },
    marks: [],
  },
};
