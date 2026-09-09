import { expectData, type ApiClient } from "@ictsc/api";
import type { Schedule, ScheduleEntry } from "../models";

/**
 * 現在アクティブなスケジュールエントリを取得
 */
export function getCurrentScheduleEntry(
  schedule: Schedule | null,
): ScheduleEntry | null {
  return schedule?.current ?? null;
}

/**
 * 次のスケジュールエントリを取得
 */
export function getNextScheduleEntry(
  schedule: Schedule | null,
): ScheduleEntry | null {
  return schedule?.next ?? null;
}

/**
 * 現在競技中かどうか（いずれかのスケジュール内にいるか）
 */
export function isInContest(schedule: Schedule | null): boolean {
  return schedule?.current != null;
}

/**
 * コンテストが開始済みかどうか（いずれかのスケジュールが過去に開始されたか）
 * 一度開始されたら、全スケジュール終了後もtrueを返す
 */
export function hasContestStarted(schedule: Schedule | null): boolean {
  return schedule?.hasStarted ?? false;
}

/**
 * 現在のスケジュールの開始時刻
 */
export function currentStartAt(schedule: Schedule | null): Date | null {
  const entry = schedule?.current;
  return entry?.startAt != null ? new Date(entry.startAt) : null;
}

/**
 * 現在のスケジュールの終了時刻
 */
export function currentEndAt(schedule: Schedule | null): Date | null {
  const entry = schedule?.current;
  return entry?.endAt != null ? new Date(entry.endAt) : null;
}

/**
 * 次のスケジュールの開始時刻
 */
export function nextStartAt(schedule: Schedule | null): Date | null {
  const entry = schedule?.next;
  return entry?.startAt != null ? new Date(entry.startAt) : null;
}

/**
 * schedule の再取得が必要になる次の時刻
 */
export function nextReloadAt(schedule: Schedule | null): Date | null {
  return currentEndAt(schedule) ?? nextStartAt(schedule);
}

export async function fetchSchedule(
  client: ApiClient,
): Promise<Schedule | null> {
  try {
    const response = expectData(
      await client.GET("/api/v1/contestant/schedule"),
    );
    const now = Date.now();
    const entries = response.sections
      .map((section) => ({
        name: section.slug,
        startAt: section.beginning,
        endAt: section.ending,
      }))
      .sort((a, b) => Date.parse(a.startAt) - Date.parse(b.startAt));
    return {
      hasStarted: entries.some((entry) => Date.parse(entry.startAt) <= now),
      current: entries.find(
        (entry) =>
          Date.parse(entry.startAt) <= now && now < Date.parse(entry.endAt),
      ),
      next: entries.find((entry) => Date.parse(entry.startAt) > now),
      entries,
    };
  } catch (e) {
    console.error(e);
    return null;
  }
}
