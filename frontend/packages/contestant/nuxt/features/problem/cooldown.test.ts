import { describe, expect, it } from "vitest";
import { remainingCooldownMinutes, remainingCooldownSeconds } from "./cooldown";

describe("problem cooldown", () => {
  const now = Date.parse("2026-09-15T12:00:00Z");

  it("returns the remaining seconds and rounds minutes up", () => {
    const next = "2026-09-15T12:03:01Z";
    expect(remainingCooldownSeconds(next, now)).toBe(181);
    expect(remainingCooldownMinutes(next, now)).toBe(4);
  });

  it("returns zero after the next submission time", () => {
    expect(remainingCooldownSeconds("2026-09-15T11:59:59Z", now)).toBe(0);
    expect(remainingCooldownMinutes(undefined, now)).toBe(0);
  });
});
