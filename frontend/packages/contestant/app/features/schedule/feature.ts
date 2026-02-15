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
  return schedule?.current ?? null;
}

/**
 * 次のスケジュールエントリを取得
 */
export function getNextScheduleEntry(schedule: Schedule | null): ScheduleEntry | null {
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
  return entry?.startAt != null ? timestampDate(entry.startAt) : null;
}

/**
 * 現在のスケジュールの終了時刻
 */
export function currentEndAt(schedule: Schedule | null): Date | null {
  const entry = schedule?.current;
  return entry?.endAt != null ? timestampDate(entry.endAt) : null;
}

/**
 * 次のスケジュールの開始時刻
 */
export function nextStartAt(schedule: Schedule | null): Date | null {
  const entry = schedule?.next;
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
