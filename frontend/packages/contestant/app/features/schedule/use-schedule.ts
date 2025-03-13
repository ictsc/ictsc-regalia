import { createContext, use } from "react";
import { Phase, type Schedule } from "@ictsc/proto/contestant/v1";
import { isAfter } from "date-fns";
import { endAt } from "./feature";

export const ScheduleContext = createContext<{
  promise: Promise<Schedule | null>;
  isPending: boolean;
} | null>(null);

export function useSchedule(): [schedule: Schedule | null, isPending: boolean] {
  const ctx = use(ScheduleContext);
  if (ctx == null) return [null, false];
  const schedule = use(ctx.promise);
  if (schedule == null) return [null, ctx.isPending];
  const end = endAt(schedule);
  if (end != null && schedule.nextPhase != null && isAfter(new Date(), end)) {
    const { nextPhase, phase: _phase, startAt: _startAt, ...rest } = schedule;
    return [
      {
        ...rest,
        phase: nextPhase,
        nextPhase: Phase.UNSPECIFIED,
        startAt: schedule.endAt,
      },
      ctx.isPending,
    ];
  }

  return [schedule, ctx.isPending];
}
