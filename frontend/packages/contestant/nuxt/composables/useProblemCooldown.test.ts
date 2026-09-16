import { beforeEach, describe, expect, it, vi } from "vitest";
import { useProblemCooldown } from "./useProblemCooldown";

describe("useProblemCooldown", () => {
  const states = new Map<string, { value: unknown }>();
  let demoMode = true;

  beforeEach(() => {
    states.clear();
    sessionStorage.clear();
    demoMode = true;
    vi.stubGlobal("useDemoMode", () => ({ value: demoMode }));
    vi.stubGlobal("useState", (key: string, init: () => unknown) => {
      if (!states.has(key)) states.set(key, { value: init() });
      return states.get(key);
    });
  });

  it("デモモードでは時間が経過しても残り時間を固定する", () => {
    const cooldown = useProblemCooldown();
    const next = "2026-09-16T03:20:00.000Z";

    expect(
      cooldown.remainingSeconds(
        "A01",
        next,
        Date.parse("2026-09-16T03:00:00Z"),
      ),
    ).toBe(20 * 60);
    expect(
      cooldown.remainingSeconds(
        "A01",
        next,
        Date.parse("2026-09-16T03:05:00Z"),
      ),
    ).toBe(20 * 60);
  });

  it("同じsessionで画面を読み直しても固定値を復元する", () => {
    const next = "2026-09-16T03:20:00.000Z";
    useProblemCooldown().remainingSeconds(
      "A01",
      next,
      Date.parse("2026-09-16T03:00:00Z"),
    );
    states.clear();

    expect(
      useProblemCooldown().remainingSeconds(
        "A01",
        next,
        Date.parse("2026-09-16T03:10:00Z"),
      ),
    ).toBe(20 * 60);
  });

  it("通常モードでは現在時刻に合わせて減少する", () => {
    demoMode = false;
    const cooldown = useProblemCooldown();
    const next = "2026-09-16T03:20:00.000Z";

    expect(
      cooldown.remainingSeconds(
        "A01",
        next,
        Date.parse("2026-09-16T03:05:00Z"),
      ),
    ).toBe(15 * 60);
  });
});
