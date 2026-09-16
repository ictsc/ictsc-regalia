import { describe, expect, it } from "vitest";
import type { Problem } from "../models";
import { problemStatus } from "./status";

const problem = (code: string): Problem => ({
  code,
  title: code,
  maxScore: 100,
  category: "server",
  submissionableSchedules: [],
});

describe("demo problem icons", () => {
  it("uses four stable samples only in demo mode", () => {
    expect(
      ["R00", "R01", "R02", "R03"].map((code) =>
        problemStatus(problem(code), false, false, true),
      ),
    ).toEqual(["unanswered", "complete", "partial", "unanswered"]);
    expect(problemStatus(problem("R01"))).toBe("unanswered");
  });
  it("prioritizes real submissions and scores without mutating data", () => {
    const p = problem("R01");
    expect(problemStatus(p, true, true, true)).toBe("pending");
    expect(p.score).toBeUndefined();
    p.score = { markedScore: 25, penalty: 0, score: 25, maxScore: 100 };
    expect(problemStatus(p, true, false, true)).toBe("partial");
  });
});
