import { beforeEach, describe, expect, it, vi } from "vitest";
import { remainingCooldownSeconds } from "../features/problem/cooldown";
import { setDemoClock } from "./useClock";

describe("demo clock", () => {
  const states = new Map<string, { value: unknown }>();

  beforeEach(() => {
    states.clear();
    states.set("demo-mode", { value: true });
    vi.stubGlobal("useState", (key: string, init: () => unknown) => {
      if (!states.has(key)) states.set(key, { value: init() });
      return states.get(key);
    });
  });

  it("提出時刻を基準に再回答までの20分を固定する", () => {
    const submittedAt = "2026-09-16T03:00:00.000Z";
    const nextSubmittableAt = "2026-09-16T03:20:00.000Z";

    setDemoClock(submittedAt);

    const demoNow = states.get("demo-now")?.value;
    expect(demoNow).toBe(Date.parse(submittedAt));
    expect(remainingCooldownSeconds(nextSubmittableAt, demoNow as number)).toBe(
      20 * 60,
    );
  });
});
