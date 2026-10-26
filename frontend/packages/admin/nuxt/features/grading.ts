import type { components } from "@ictsc/api";

export type GradingAnswer = components["schemas"]["AdminAnswer"];
export type GradingStatus = GradingAnswer["workflow"]["status"];
export type GradingFilters = {
  /** Team names to show; empty shows every team. */
  teams: string[];
  /** Problem codes to show; empty shows every problem. */
  problems: string[];
  /** Only answers submitted in the last 20 minutes. */
  recent: boolean;
  /** Include answers that already have the full score. */
  perfect: boolean;
  /** Only answers without a score. */
  unscored: boolean;
};

export const gradingColumns: { status: GradingStatus; label: string }[] = [
  { status: "WAITING", label: "採点待ち" },
  { status: "IN_PROGRESS", label: "採点中" },
  { status: "COMPLETED", label: "採点完了" },
];

export const defaultGradingFilters = (): GradingFilters => ({
  teams: [],
  problems: [],
  recent: false,
  perfect: true,
  unscored: false,
});

export const markingBlockedLabels: Record<
  NonNullable<GradingAnswer["marking_blocked_reason"]>,
  string
> = {
  grading_not_started: "競技終了後に採点",
  answer_not_latest: "採点対象外（旧回答）",
};

const recentWindowMs = 20 * 60_000;

export function answerKey(answer: GradingAnswer) {
  const r = answer.reference;
  return `${r.team_code}:${r.problem_code}:${r.answer_number}`;
}

export function answerPath(answer: GradingAnswer) {
  const r = answer.reference;
  return `/submissions/${r.problem_code}/${r.team_code}/${r.answer_number}`;
}

export function minutesSince(submittedAt: string, now: number) {
  return Math.max(0, Math.floor((now - Date.parse(submittedAt)) / 60_000));
}

/**
 * Answers that belong on the grading board. Secondary qualifier mode grades
 * only each team's final answer per problem, so earlier ones are left out.
 */
export function gradableAnswers(answers: readonly GradingAnswer[]) {
  return answers.filter(
    (a) => a.competition_mode !== "SECONDARY_QUALIFIER" || a.latest,
  );
}

export function filterAnswers(
  answers: readonly GradingAnswer[],
  filters: GradingFilters,
  now: number,
) {
  return answers.filter(
    (a) =>
      (!filters.teams.length || filters.teams.includes(a.team.name)) &&
      (!filters.problems.length || filters.problems.includes(a.problem.code)) &&
      (!filters.recent || Date.parse(a.submitted_at) >= now - recentWindowMs) &&
      (filters.perfect || !a.score || a.score.total !== a.score.max) &&
      (!filters.unscored || !a.score),
  );
}

/** Answers in one column: the viewer's own first, then newest first. */
export function columnAnswers(
  answers: readonly GradingAnswer[],
  status: GradingStatus,
  myDiscordId: string,
) {
  const mine = (a: GradingAnswer) =>
    Number(a.workflow.assignee_discord_id === myDiscordId && !!myDiscordId);
  return answers
    .filter((a) => a.workflow.status === status)
    .toSorted(
      (a, b) =>
        mine(b) - mine(a) ||
        Date.parse(b.submitted_at) - Date.parse(a.submitted_at) ||
        answerKey(a).localeCompare(answerKey(b)),
    );
}
