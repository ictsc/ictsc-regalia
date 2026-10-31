import { describe, expect, it } from "vitest";
import { adminRefreshKeys } from "./realtime";

const answer = (reason: "SUBMITTED" | "WORKFLOW_UPDATED" | "MARKING_UPDATED") =>
  ({
    type: "answer",
    reason,
    teamCode: 2,
    problemCode: "A",
    answerNumber: 1,
  }) as const;

describe("admin realtime refresh", () => {
  it("maps answer changes to the overview and score pages", () => {
    expect(adminRefreshKeys(answer("SUBMITTED"))).toEqual(["admin-overview"]);
    expect(adminRefreshKeys(answer("WORKFLOW_UPDATED"))).toBe(null);
    expect(adminRefreshKeys(answer("MARKING_UPDATED"))).toEqual([
      "admin-scores",
    ]);
  });

  it("refetches everything only after a reconnect", () => {
    expect(adminRefreshKeys({ type: "ready", reconnected: false })).toBe(null);
    expect(adminRefreshKeys({ type: "ready", reconnected: true })).toBe("all");
  });
});
