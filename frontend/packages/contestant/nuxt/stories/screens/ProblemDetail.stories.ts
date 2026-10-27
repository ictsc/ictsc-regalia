import type { Meta, StoryObj } from "@storybook/vue3-vite";
import { ref } from "vue";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import CompetitionHeaderView from "~/components/header/CompetitionHeaderView.vue";
import AnswerForm from "~/components/problem/AnswerForm.vue";
import AnswerHistory from "~/components/problem/AnswerHistory.vue";
import AnswerResizer from "~/components/problem/AnswerResizer.vue";
import DeploymentHistory from "~/components/problem/DeploymentHistory.vue";
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

type ScreenArgs = { cooldownMinutes: number };

const meta = {
  title: "Contestant/Screens/ProblemDetail",
  decorators: [withContestantFrame()],
  parameters: {
    layout: "fullscreen",
    bodyClass: "detail-page",
    route: "/problems/B02",
  },
  args: { cooldownMinutes: 0 },
  render: (args) => ({
    components: {
      AnswerForm,
      AnswerHistory,
      AnswerResizer,
      CompetitionHeaderView,
      DeploymentHistory,
      MarkdownContent,
      ProblemHero,
      ProblemRail,
    },
    setup: () => ({
      args,
      body: ref(""),
      contentGrid: ref<HTMLElement | null>(null),
      answerWidth: ref<number>(),
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
  <div class="detail-layout">
    <ProblemRail :items="railItems" :current-code="problem.code" />
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
      <div
        ref="contentGrid"
        class="content-grid"
        :style="answerWidth ? { '--answer-width': answerWidth + 'px' } : {}"
      >
        <AnswerResizer v-model:width="answerWidth" :container="contentGrid" />
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
          :max-length="3000"
          :allow-attachments="true"
        />
        <DeploymentHistory
          :deployments="deployments"
          :penalty-threshold="2"
          :deploying="false"
          :stream-error="false"
          :reset-disabled="false"
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

export const Cooldown: Story = {
  args: { cooldownMinutes: 18 },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
