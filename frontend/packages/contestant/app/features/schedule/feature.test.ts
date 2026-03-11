import { create } from "@bufbuild/protobuf";
import { timestampFromDate } from "@bufbuild/protobuf/wkt";
import {
  ScheduleEntrySchema,
  ScheduleSchema,
  type Schedule,
} from "@ictsc/proto/contestant/v1";
import { describe, expect, it } from "vitest";
import { nextReloadAt } from "./feature";

function scheduleEntry(
  name: string,
  startAt: Date,
  endAt: Date,
) {
  return create(ScheduleEntrySchema, {
    name,
    startAt: timestampFromDate(startAt),
    endAt: timestampFromDate(endAt),
  });
}

function schedule(data: {
  current?: ReturnType<typeof scheduleEntry>;
  next?: ReturnType<typeof scheduleEntry>;
}): Schedule {
  return create(ScheduleSchema, {
    hasStarted: data.current != null,
    current: data.current,
    next: data.next,
  });
}

describe("nextReloadAt", () => {
  it("returns the current schedule end while a slot is active", () => {
    const endAt = new Date("2026-03-11T10:30:00.000Z");

    expect(
      nextReloadAt(
        schedule({
          current: scheduleEntry(
            "day1-am",
            new Date("2026-03-11T09:00:00.000Z"),
            endAt,
          ),
          next: scheduleEntry(
            "day1-pm",
            new Date("2026-03-11T13:00:00.000Z"),
            new Date("2026-03-11T15:00:00.000Z"),
          ),
        }),
      ),
    ).toEqual(endAt);
  });

  it("returns the next schedule start while waiting for the next slot", () => {
    const startAt = new Date("2026-03-11T13:00:00.000Z");

    expect(
      nextReloadAt(
        schedule({
          next: scheduleEntry(
            "day1-pm",
            startAt,
            new Date("2026-03-11T15:00:00.000Z"),
          ),
        }),
      ),
    ).toEqual(startAt);
  });

  it("returns null after all schedule slots have ended", () => {
    expect(nextReloadAt(schedule({}))).toBeNull();
  });

  it("returns null when schedule loading failed", () => {
    expect(nextReloadAt(null)).toBeNull();
  });
});
