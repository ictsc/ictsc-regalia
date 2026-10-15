import type { ScheduleEntry } from "../models";

export type ScheduleTemporalStatus = "past" | "current" | "future";

export function getTemporalStatus(
  entry: ScheduleEntry,
  now: Date,
): ScheduleTemporalStatus {
  if (now >= new Date(entry.endAt)) return "past";
  if (now >= new Date(entry.startAt)) return "current";
  return "future";
}

export function startAtMs(entry: ScheduleEntry): number {
  return new Date(entry.startAt).getTime();
}

export function endAtMs(entry: ScheduleEntry): number {
  return new Date(entry.endAt).getTime();
}
