import { createContext, use } from "react";
import { type Schedule } from "@ictsc/proto/contestant/v1";

export const ScheduleContext = createContext<{
  promise: Promise<Schedule | null>;
  isPending: boolean;
} | null>(null);

export function useSchedule(): [schedule: Schedule | null, isPending: boolean] {
  const ctx = use(ScheduleContext);
  if (ctx == null) return [null, false];
  const schedule = use(ctx.promise);
  if (schedule == null) return [null, ctx.isPending];
  return [schedule, ctx.isPending];
}
