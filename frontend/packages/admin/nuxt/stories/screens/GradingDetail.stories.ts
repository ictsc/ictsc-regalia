import type { Meta, StoryObj } from "@storybook/vue3-vite";
import { ref } from "vue";
import GradingDetailView from "~/components/grading/GradingDetailView.vue";
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
  title: "Admin/Screens/GradingDetail",
  decorators: [
    withAdminShell("/submissions", "workspace grading-detail-workspace"),
  ],
  parameters: { layout: "fullscreen" },
  render: () => ({
    components: { GradingDetailView },
    setup: () => ({
      answer: detailAnswer,
      problem: adminProblem,
      activeCommit: contentStatus.active_commit!,
      previous: previousAnswers,
      marks: markings,
      collapsed: ref(defaultPaneLayout().collapsed),
      widths: ref(defaultPaneLayout().widths),
      submitted: new Date(detailAnswer.submitted_at).toLocaleString("ja-JP"),
      submitMarking: async () => true,
    }),
    template: `
<a class="text-link" href="#">← 採点一覧</a>
<h1>B02 / Team 4 / #3</h1>
<p class="admin-lead">{{ answer.team.name }} · {{ answer.author.display_name }} · <time :datetime="answer.submitted_at">{{ submitted }}</time></p>
<p role="status"></p>
<GradingDetailView
  v-model:collapsed="collapsed"
  v-model:widths="widths"
  :answer="answer"
  :problem="problem"
  :active-commit="activeCommit"
  :previous="previous"
  :marks="marks"
  :busy="false"
  :workflow-busy="false"
  :submit-marking="submitMarking"
/>`,
  }),
} satisfies Meta;
export default meta;
type Story = StoryObj<typeof meta>;

export const Desktop: Story = {};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
