import type { Meta, StoryObj } from "@storybook/vue3-vite";
import { computed, reactive } from "vue";
import GradingBoard from "~/components/grading/GradingBoard.vue";
import GradingFilters from "~/components/grading/GradingFilters.vue";
import { defaultGradingFilters, filterAnswers } from "~/features/grading";
import { withAdminShell } from "~/stories/decorators";
import { answers, fixedNow, myDiscordId } from "~/stories/fixtures";

const meta = {
  title: "Admin/Screens/GradingBoard",
  decorators: [withAdminShell("/submissions", "workspace grading-workspace")],
  parameters: { layout: "fullscreen" },
  render: () => ({
    components: { GradingBoard, GradingFilters },
    setup() {
      const filters = reactive(defaultGradingFilters());
      const filtered = computed(() =>
        filterAnswers(answers, filters, fixedNow),
      );
      return {
        filters,
        filtered,
        myDiscordId,
        fixedNow,
        problemCodes: [...new Set(answers.map((a) => a.problem.code))],
        teamNames: [...new Set(answers.map((a) => a.team.name))],
      };
    },
    template: `
<h1>採点</h1>
<p class="admin-lead">自分の担当回答を各列の先頭に表示します。カードを開くと採点できます。</p>
<GradingFilters
  v-model:teams="filters.teams"
  v-model:problems="filters.problems"
  v-model:recent="filters.recent"
  v-model:perfect="filters.perfect"
  v-model:unscored="filters.unscored"
  :problem-codes="problemCodes"
  :team-names="teamNames"
/>
<GradingBoard :answers="filtered" :my-discord-id="myDiscordId" busy-key="" :now="fixedNow" />`,
  }),
} satisfies Meta;
export default meta;
type Story = StoryObj<typeof meta>;

export const Desktop: Story = {};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
