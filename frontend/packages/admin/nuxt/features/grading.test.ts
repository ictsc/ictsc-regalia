import { describe, expect, it } from "vitest";
import {
  answerKey,
  answerPath,
  columnAnswers,
  defaultGradingFilters,
  filterAnswers,
  minutesSince,
  type GradingAnswer,
} from "./grading";

const now = Date.parse("2026-09-02T12:00:00Z");
const answer = (
  team: number,
  problem: string,
  number: number,
  minutesAgo: number,
  overrides: {
    status?: GradingAnswer["workflow"]["status"];
    assignee?: string | null;
    score?: { total: number; max: number } | null;
  } = {},
): GradingAnswer => ({
  reference: { team_code: team, problem_code: problem, answer_number: number },
  team: { code: team, name: `T${team}`, organization: "", member_limit: 4 },
  author: { name: "a", display_name: "A", self_introduction: "" },
  problem: { code: problem, title: problem, max_score: 100, category: "" },
  body: { type: "DESCRIPTIVE", body: "" },
  submitted_at: new Date(now - minutesAgo * 60_000).toISOString(),
  score:
    overrides.score === undefined || overrides.score === null
      ? null
      : { ...overrides.score, marked: overrides.score.total, penalty: 0 },
  content_commit: "c",
  workflow: {
    status: overrides.status ?? "WAITING",
    assignee_discord_id: overrides.assignee ?? null,
    assignment_source: "DEFAULT",
    revision: 1,
  },
});

describe("answer identity", () => {
  it("builds keys and paths from the reference", () => {
    const a = answer(12, "B02", 3, 0);
    expect(answerKey(a)).toBe("12:B02:3");
    expect(answerPath(a)).toBe("/submissions/B02/12/3");
  });
  it("counts whole minutes since submission and never goes negative", () => {
    expect(minutesSince(new Date(now - 125_000).toISOString(), now)).toBe(2);
    expect(minutesSince(new Date(now + 60_000).toISOString(), now)).toBe(0);
  });
});

describe("filterAnswers", () => {
  const answers = [
    answer(2, "A01", 1, 5, { score: { total: 100, max: 100 } }),
    answer(3, "A01", 1, 30),
    answer(3, "B02", 2, 10, { score: { total: 40, max: 100 } }),
  ];
  const filter = (change: Partial<ReturnType<typeof defaultGradingFilters>>) =>
    filterAnswers(answers, { ...defaultGradingFilters(), ...change }, now).map(
      answerKey,
    );
  it("shows everything by default", () => {
    expect(filter({})).toHaveLength(3);
  });
  it("filters by team, problem, recency, perfect score and unscored", () => {
    expect(filter({ teams: ["T3"] })).toEqual(["3:A01:1", "3:B02:2"]);
    expect(filter({ problems: ["B02"] })).toEqual(["3:B02:2"]);
    expect(filter({ recent: true })).toEqual(["2:A01:1", "3:B02:2"]);
    expect(filter({ perfect: false })).toEqual(["3:A01:1", "3:B02:2"]);
    expect(filter({ unscored: true })).toEqual(["3:A01:1"]);
  });
});

describe("columnAnswers", () => {
  it("puts the viewer's answers first, then the newest", () => {
    const answers = [
      answer(2, "A01", 1, 30),
      answer(3, "A01", 1, 5),
      answer(4, "A01", 1, 50, { assignee: "me" }),
      answer(5, "A01", 1, 1, { status: "COMPLETED" }),
    ];
    expect(columnAnswers(answers, "WAITING", "me").map(answerKey)).toEqual([
      "4:A01:1",
      "3:A01:1",
      "2:A01:1",
    ]);
    expect(columnAnswers(answers, "WAITING", "").map(answerKey)[0]).toBe(
      "3:A01:1",
    );
  });
});
