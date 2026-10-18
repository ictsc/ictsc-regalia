import type { Meta, StoryObj } from "@storybook/vue3-vite";
import { reactive } from "vue";
import GradingFilters from "./GradingFilters.vue";
import { defaultGradingFilters } from "~/features/grading";
import { withAdminShell } from "~/stories/decorators";
import { catalog, teams } from "~/stories/fixtures";

const meta = {
  component: GradingFilters,
  title: "Admin/Grading/GradingFilters",
  decorators: [withAdminShell("/submissions", "workspace grading-workspace")],
  parameters: { layout: "fullscreen" },
  args: {
    ...defaultGradingFilters(),
    problemCodes: catalog.map((p) => p.code),
    teamNames: teams.map((t) => t.name),
  },
  // Keep the filters interactive in the preview.
  render: (args) => ({
    components: { GradingFilters },
    setup: () => ({ args, filters: reactive({ ...args }) }),
    template: `<GradingFilters v-bind="args" v-model:teams="filters.teams" v-model:problems="filters.problems" v-model:recent="filters.recent" v-model:perfect="filters.perfect" v-model:unscored="filters.unscored" />`,
  }),
} satisfies Meta<typeof GradingFilters>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Filtered: Story = {
  args: {
    problems: ["B02"],
    teams: ["KERNEL PANIC"],
    recent: true,
    perfect: false,
    unscored: true,
  },
};
