import { type Transport, createClient } from "@connectrpc/connect";
import {
  ContestService,
  type Schedule,
  Phase,
} from "@ictsc/proto/contestant/v1";
import { timestampDate } from "@bufbuild/protobuf/wkt";

export function startAt(schedule: Schedule): Date | null {
  return schedule.startAt != null ? timestampDate(schedule.startAt) : null;
}

export function endAt(schedule: Schedule): Date | null {
  return schedule.endAt != null ? timestampDate(schedule.endAt) : null;
}

export function isInContest(schedule: Schedule | null): boolean {
  const phase = schedule?.phase ?? Phase.UNSPECIFIED;
  return phase === Phase.IN_CONTEST
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
