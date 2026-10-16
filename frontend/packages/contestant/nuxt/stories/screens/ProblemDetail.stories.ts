import type { Meta, StoryObj } from "@storybook/vue3-vite";
import { ref } from "vue";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import CompetitionHeaderView from "~/components/header/CompetitionHeaderView.vue";
import AnswerForm from "~/components/problem/AnswerForm.vue";
import AnswerHistory from "~/components/problem/AnswerHistory.vue";
import DeploymentHistory from "~/components/problem/DeploymentHistory.vue";
import ProblemDetailToolbar from "~/components/problem/ProblemDetailToolbar.vue";
import ProblemHero from "~/components/problem/ProblemHero.vue";
import ProblemRail from "~/components/problem/ProblemRail.vue";
import { headerRankingItems } from "~/features/ranking";
import { contestClock } from "~/features/schedule/clock";
import { withContestantFrame } from "~/stories/decorators";
import {
  answers,
  fixedNow,
  makeDeployment,
  ownTeamCode,
  problemBody,
  problems,
  ranking,
  sections,
} from "~/stories/fixtures";

const problem = problems.find((p) => p.code === "B02")!;

type ScreenArgs = { railOpen: boolean; cooldownMinutes: number };

const meta = {
  title: "Contestant/Screens/ProblemDetail",
  decorators: [withContestantFrame()],
  parameters: {
    layout: "fullscreen",
    bodyClass: "detail-page",
    route: "/problems/B02",
  },
  args: { railOpen: true, cooldownMinutes: 0 },
  render: (args) => ({
    components: {
      AnswerForm,
      AnswerHistory,
      CompetitionHeaderView,
      DeploymentHistory,
      MarkdownContent,
      ProblemDetailToolbar,
      ProblemHero,
      ProblemRail,
    },
    setup: () => ({
      args,
      rail: ref(args.railOpen),
      body: ref(""),
      header: {
        compact: true,
        rankingItems: headerRankingItems(ranking, ownTeamCode),
        ownTeamCode,
        maxScore: problems.reduce((sum, p) => sum + p.maxScore, 0),
        clock: contestClock(Object.values(sections), fixedNow),
        demoMode: false,
        hasUnreadNotices: false,
        teamName: "KERNEL PANIC",
        routeKey: "/problems/B02",
      },
      railItems: problems.map((p) => ({
        problem: p,
        cooldownMinutes: p.code === problem.code ? args.cooldownMinutes : 0,
        pending: p.code === "B04",
      })),
      problem,
      problemBody,
      answers,
      deployments: [makeDeployment(1, "COMPLETED")],
    }),
    template: `
<CompetitionHeaderView v-bind="header" />
<main>
  <ProblemDetailToolbar v-model:rail-open="rail" :reset-disabled="false" />
  <div class="detail-layout" :class="{ 'is-rail-hidden': !rail }">
    <ProblemRail v-if="rail" :items="railItems" :current-code="problem.code" />
    <article class="problem-content">
      <ProblemHero
        :code="problem.code"
        :title="problem.title"
        :category="problem.category"
        :score="problem.score"
        :max-score="problem.maxScore"
        :cooldown-minutes="args.cooldownMinutes"
        :open="true"
      />
      <div class="content-grid">
        <section class="reading-section">
          <div class="reading-copy"><MarkdownContent :source="problemBody" /></div>
        </section>
        <AnswerForm
          v-model:body="body"
          :open="true"
          :sending="false"
          :submitted="false"
          :cooldown-seconds="args.cooldownMinutes * 60"
          saved-at=""
          save-error=""
        />
        <DeploymentHistory
          :deployments="deployments"
          :penalty-threshold="2"
          :deploying="false"
          :stream-error="false"
        />
        <AnswerHistory
          :problem-code="problem.code"
          :answers="answers"
          :selected-answer-number="null"
        />
      </div>
    </article>
  </div>
</main>`,
  }),
} satisfies Meta<ScreenArgs>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Desktop: Story = {};

export const RailHidden: Story = {
  args: { railOpen: false },
};

export const Cooldown: Story = {
  args: { cooldownMinutes: 18 },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
