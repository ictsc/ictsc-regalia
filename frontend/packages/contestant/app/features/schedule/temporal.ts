import type { ScheduleEntry } from "@ictsc/proto/contestant/v1";
import { timestampDate } from "@bufbuild/protobuf/wkt";

export type ScheduleTemporalStatus = "past" | "current" | "future";

export function getTemporalStatus(
  entry: ScheduleEntry,
  now: Date,
): ScheduleTemporalStatus {
  if (entry.endAt != null && now >= timestampDate(entry.endAt)) return "past";
  if (entry.startAt != null && now >= timestampDate(entry.startAt))
    return "current";
  return "future";
}

export function startAtMs(entry: ScheduleEntry): number {
  return entry.startAt != null ? timestampDate(entry.startAt).getTime() : 0;
}

export function endAtMs(entry: ScheduleEntry): number {
  return entry.endAt != null ? timestampDate(entry.endAt).getTime() : Infinity;
}
