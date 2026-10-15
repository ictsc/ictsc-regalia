import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { loadDemoClock } from "./useClock";

describe("demo clock", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-09-19T03:00:00Z"));
  });
  afterEach(() => {
    vi.restoreAllMocks();
    vi.useRealTimers();
  });

  it("restores the same reference time after real time advances", () => {
    const first = loadDemoClock();
    vi.setSystemTime(new Date("2026-09-20T03:00:00Z"));
    expect(loadDemoClock()).toBe(first);
  });

  it("replaces an invalid stored reference", () => {
    localStorage.setItem("ictsc-demo-now", "invalid");
    expect(loadDemoClock()).toBe(Date.now());
    expect(localStorage.getItem("ictsc-demo-now")).toBe(String(Date.now()));
  });

  it("falls back to the page clock when storage is blocked", () => {
    vi.spyOn(Storage.prototype, "getItem").mockImplementation(() => {
      throw new Error("storage blocked");
    });
    expect(loadDemoClock()).toBe(Date.now());
  });
});
