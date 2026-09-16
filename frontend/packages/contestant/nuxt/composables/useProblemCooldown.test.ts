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

  it("未回答のデモ問題は2問が提出可能、2問が待機中になる", () => {
    const cooldown = useProblemCooldown();
    expect(
      ["R00", "R01", "R02", "R03"].map((code) =>
        cooldown.remainingSeconds(code, undefined, 0),
      ),
    ).toEqual([0, 0, 480, 900]);
    expect(cooldown.remainingSeconds("R02", undefined, 999999)).toBe(480);
    states.clear();
    expect(useProblemCooldown().remainingSeconds("R02", undefined, 0)).toBe(
      480,
    );
  });

  it("通常モードではデモの待ち時間を使わない", () => {
    demoMode = false;
    expect(useProblemCooldown().remainingSeconds("R02", undefined, 0)).toBe(0);
    expect(useProblemCooldown().remainingSeconds("R02", 900000, 0)).toBe(900);
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
