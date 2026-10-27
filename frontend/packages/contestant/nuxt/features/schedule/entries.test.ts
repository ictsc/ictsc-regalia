import { describe, expect, it } from "vitest";
import { scheduleEntryState, sortedScheduleEntries } from "./entries";

const day1 = {
  name: "day1",
  startAt: "2026-10-26T10:00:00+09:00",
  endAt: "2026-10-30T17:00:00+09:00",
};
const day0 = {
  name: "day0",
  startAt: "2026-10-25T10:00:00+09:00",
  endAt: "2026-10-25T17:00:00+09:00",
};

describe("scheduleEntryState", () => {
  it("follows the entry's start and end", () => {
    expect(
      scheduleEntryState(day1, Date.parse("2026-10-26T09:59:59+09:00")),
    ).toBe("upcoming");
    expect(
      scheduleEntryState(day1, Date.parse("2026-10-26T10:00:00+09:00")),
    ).toBe("running");
    // Entries may span several days.
    expect(
      scheduleEntryState(day1, Date.parse("2026-10-28T03:00:00+09:00")),
    ).toBe("running");
    expect(
      scheduleEntryState(day1, Date.parse("2026-10-30T17:00:00+09:00")),
    ).toBe("ended");
  });
});

describe("sortedScheduleEntries", () => {
  it("orders by start without changing the input", () => {
    const entries = [day1, day0];
    expect(sortedScheduleEntries(entries).map((e) => e.name)).toEqual([
      "day0",
      "day1",
    ]);
    expect(entries.map((e) => e.name)).toEqual(["day1", "day0"]);
  });
});
