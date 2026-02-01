import { type Transport, createClient } from "@connectrpc/connect";
import {
  ContestService,
  type Schedule,
  type ScheduleEntry,
} from "@ictsc/proto/contestant/v1";
import { timestampDate } from "@bufbuild/protobuf/wkt";

/**
 * 現在アクティブなスケジュールエントリを取得
 */
export function getCurrentScheduleEntry(schedule: Schedule | null): ScheduleEntry | null {
  if (schedule == null) return null;
  const now = new Date();
  for (const entry of schedule.schedules) {
    const start = entry.startAt != null ? timestampDate(entry.startAt) : null;
    const end = entry.endAt != null ? timestampDate(entry.endAt) : null;
    if (start != null && end != null && now >= start && now < end) {
      return entry;
    }
  }
  return null;
}

/**
 * 次のスケジュールエントリを取得
 */
export function getNextScheduleEntry(schedule: Schedule | null): ScheduleEntry | null {
  if (schedule == null) return null;
  const now = new Date();
  let nextEntry: ScheduleEntry | null = null;
  let nextStartAt: Date | null = null;
  for (const entry of schedule.schedules) {
    const start = entry.startAt != null ? timestampDate(entry.startAt) : null;
    if (start != null && start > now) {
      if (nextStartAt == null || start < nextStartAt) {
        nextEntry = entry;
        nextStartAt = start;
      }
    }
  }
  return nextEntry;
}

/**
 * 現在競技中かどうか（いずれかのスケジュール内にいるか）
 */
export function isInContest(schedule: Schedule | null): boolean {
  return getCurrentScheduleEntry(schedule) != null;
}

/**
 * 現在のスケジュールの開始時刻
 */
export function currentStartAt(schedule: Schedule | null): Date | null {
  const entry = getCurrentScheduleEntry(schedule);
  return entry?.startAt != null ? timestampDate(entry.startAt) : null;
}

/**
 * 現在のスケジュールの終了時刻
 */
export function currentEndAt(schedule: Schedule | null): Date | null {
  const entry = getCurrentScheduleEntry(schedule);
  return entry?.endAt != null ? timestampDate(entry.endAt) : null;
}

/**
 * 次のスケジュールの開始時刻
 */
export function nextStartAt(schedule: Schedule | null): Date | null {
  const entry = getNextScheduleEntry(schedule);
  return entry?.startAt != null ? timestampDate(entry.startAt) : null;
}

export async function fetchSchedule(
  transport: Transport,
): Promise<Schedule | null> {
  try {
    const client = createClient(ContestService, transport);
    const resp = await client.getSchedule({});
    return resp.schedule ?? null;
  } catch (e) {
    console.error(e);
    return null;
  }
}
