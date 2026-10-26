import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ProblemListView, { type ProblemListItem } from "./ProblemListView.vue";
import type { Problem } from "~/features/models";
import { problemStatus } from "~/features/problem/status";
import { withContestantFrame } from "~/stories/decorators";
import { makeProblem, problems } from "~/stories/fixtures";

const toItems = (list: Problem[]): ProblemListItem[] =>
  list.map((problem) => ({
    problem,
    status: problemStatus(problem),
    cooldownMinutes: 0,
  }));

const meta = {
  component: ProblemListView,
  title: "Contestant/Problems/ProblemListView",
  decorators: [
    withContestantFrame(),
    () => ({ template: "<main><story /></main>" }),
  ],
  parameters: { layout: "fullscreen", route: "/problems" },
  args: { items: toItems(problems) },
} satisfies Meta<typeof ProblemListView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const Cooldown: Story = {
  args: {
    items: toItems(problems).map((item) =>
      ["A02", "B03"].includes(item.problem.code)
        ? { ...item, cooldownMinutes: item.problem.code === "A02" ? 18 : 3 }
        : item,
    ),
  },
};

export const Pending: Story = {
  args: {
    items: toItems(problems).map((item) =>
      ["A04", "B03"].includes(item.problem.code)
        ? { ...item, status: "pending" as const }
        : item,
    ),
  },
};

/** 二次予選モードで最終公開前。得点は返らず、提出済みの問題は「提出済み」。 */
export const SecondaryQualifier: Story = {
  args: {
    resultsHidden: true,
    items: problems.map((problem) => {
      const hidden = { ...problem, score: undefined };
      return {
        problem: hidden,
        status: problemStatus(
          hidden,
          ["A01", "A04", "B03"].includes(problem.code),
          false,
          true,
        ),
        cooldownMinutes: problem.code === "B03" ? 12 : 0,
      };
    }),
  },
};

export const AllClosed: Story = {
  args: {
    items: toItems(
      problems.map((problem) => ({
        ...problem,
        submissionStatus: { isSubmittable: false },
      })),
    ),
  },
};

export const LongTitles: Story = {
  args: {
    items: toItems([
      makeProblem({
        code: "C01",
        title:
          "複数拠点を結ぶIPsec VPNのうち一部の拠点間でのみ通信が断続的に途切れ、経路の再計算が繰り返される",
        category: "ネットワーク",
      }),
      makeProblem({
        code: "C02",
        title:
          "Kubernetesクラスタ上のサービスメッシュで特定のPodからのみ外部APIへの接続がタイムアウトする",
        category: "クラウドネイティブ",
        maxScore: 1000,
      }),
    ]),
  },
};

export const Empty: Story = {
  args: { items: [] },
};

export const Loading: Story = {
  args: { items: [], pending: true },
};

export const Error: Story = {
  args: { items: [], error: "問題一覧を取得できませんでした" },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
