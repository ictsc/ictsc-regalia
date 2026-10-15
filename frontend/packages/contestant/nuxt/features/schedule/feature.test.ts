import { describe, expect, it } from "vitest";
import type { Schedule, ScheduleEntry } from "../models";
import { nextReloadAt } from "./feature";

function entry(name: string, startAt: string, endAt: string): ScheduleEntry {
  return { name, startAt, endAt };
}

function schedule(current?: ScheduleEntry, next?: ScheduleEntry): Schedule {
  return {
    hasStarted: current != null,
    current,
    next,
    entries: [current, next].filter(Boolean) as ScheduleEntry[],
  };
}

describe("nextReloadAt", () => {
  it("returns the current section end", () => {
    const current = entry(
      "day1",
      "2026-08-30T01:00:00Z",
      "2026-08-30T03:00:00Z",
    );
    expect(nextReloadAt(schedule(current))).toEqual(new Date(current.endAt));
  });

  it("returns the next section start while waiting", () => {
    const next = entry("day2", "2026-08-31T01:00:00Z", "2026-08-31T03:00:00Z");
    expect(nextReloadAt(schedule(undefined, next))).toEqual(
      new Date(next.startAt),
    );
  });

  it("returns null after the competition", () => {
    expect(nextReloadAt(schedule())).toBeNull();
  });
});
