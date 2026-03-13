import { createFileRoute } from "@tanstack/react-router";
import type { ScheduleEntry } from "@ictsc/proto/contestant/v1";
import { IndexPage, Timer, type ContestState } from "./index.page";
import {
  useSchedule,
  getCurrentScheduleEntry,
  getNextScheduleEntry,
  currentEndAt,
  nextStartAt,
} from "@app/features/schedule";

export const Route = createFileRoute("/")({
  component: Page,
});

function Page() {
  const [schedule] = useSchedule();

  const currentEntry = getCurrentScheduleEntry(schedule);
  const nextEntry = getNextScheduleEntry(schedule);
  const entries: ScheduleEntry[] = schedule?.entries ?? [];

  let state: ContestState;
  let timerEndMs: number;

  if (currentEntry != null) {
    state = "in_contest";
    const endAt = currentEndAt(schedule);
    timerEndMs = endAt != null ? endAt.getTime() : 0;
  } else if (nextEntry != null) {
    state = "waiting";
    const startAt = nextStartAt(schedule);
    timerEndMs = startAt != null ? startAt.getTime() : 0;
  } else {
    state = "ended";
    timerEndMs = 0;
  }

  return (
    <IndexPage
      state={state}
      currentScheduleName={currentEntry?.name}
      nextScheduleName={nextEntry?.name}
      timer={timerEndMs > 0 ? <Timer endMs={timerEndMs} /> : undefined}
      entries={entries}
    />
  );
}
