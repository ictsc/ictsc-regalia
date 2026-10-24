import { describe, expect, it } from "vitest";
import { affectsProblem, contestantRefreshPlan } from "./realtime";

describe("contestant realtime refresh", () => {
  it("refetches team data immediately and spreads broadcast refetches", () => {
    expect(
      contestantRefreshPlan({ type: "answers", problemCode: "A" }),
    ).toEqual({ keys: ["competition", "activity"], spread: false });
    expect(contestantRefreshPlan({ type: "ranking" })).toEqual({
      keys: ["competition", "ranking"],
      spread: true,
    });
    expect(contestantRefreshPlan({ type: "competition" })).toEqual({
      keys: "all",
      spread: true,
    });
  });

  it("skips the first ready but refetches everything after a reconnect", () => {
    expect(contestantRefreshPlan({ type: "ready", reconnected: false })).toBe(
      null,
    );
    expect(contestantRefreshPlan({ type: "ready", reconnected: true })).toEqual(
      { keys: "all", spread: true },
    );
  });

  it("refreshes a problem page only for its own problem or new content", () => {
    expect(affectsProblem({ type: "answers", problemCode: "A" }, "A")).toBe(
      true,
    );
    expect(affectsProblem({ type: "scores", problemCode: "B" }, "A")).toBe(
      false,
    );
    expect(affectsProblem({ type: "content" }, "A")).toBe(true);
    expect(affectsProblem({ type: "ranking" }, "A")).toBe(false);
  });
});
