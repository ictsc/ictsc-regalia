import type { ScheduleEntry } from "../models";

export type ScheduleEntryState = "upcoming" | "running" | "ended";

export const scheduleEntryStateLabels: Record<ScheduleEntryState, string> = {
  upcoming: "開始前",
  running: "開催中",
  ended: "終了",
};

export function scheduleEntryState(
  entry: ScheduleEntry,
  now: number,
): ScheduleEntryState {
  if (now < Date.parse(entry.startAt)) return "upcoming";
  if (now < Date.parse(entry.endAt)) return "running";
  return "ended";
}

/** Entries in start order, so the page reads as a timeline. */
export function sortedScheduleEntries(
  entries: readonly ScheduleEntry[],
): ScheduleEntry[] {
  return [...entries].sort(
    (a, b) => Date.parse(a.startAt) - Date.parse(b.startAt),
  );
}
