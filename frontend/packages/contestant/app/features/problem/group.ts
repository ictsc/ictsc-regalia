import type { Problem, ScheduleEntry } from "@ictsc/proto/contestant/v1";
import { timestampDate } from "@bufbuild/protobuf/wkt";

export type ScheduleTemporalStatus = "past" | "current" | "future";

export type GroupScheduleInfo = {
  name: string;
  temporalStatus: ScheduleTemporalStatus;
};

export type ProblemGroup = {
  key: string;
  schedules: GroupScheduleInfo[];
  problems: Problem[];
  hasSubmittableProblem: boolean;
};

function getTemporalStatus(
  entry: ScheduleEntry,
  now: Date,
): ScheduleTemporalStatus {
  if (entry.endAt != null && now >= timestampDate(entry.endAt)) return "past";
  if (entry.startAt != null && now >= timestampDate(entry.startAt))
    return "current";
  return "future";
}

function startAtMs(entry: ScheduleEntry): number {
  return entry.startAt != null ? timestampDate(entry.startAt).getTime() : 0;
}

export function groupProblems(
  problems: Problem[],
  now: Date = new Date(),
): ProblemGroup[] {
  const groupMap = new Map<
    string,
    {
      problems: Problem[];
      entries: ScheduleEntry[];
      hasSubmittableProblem: boolean;
    }
  >();

  for (const problem of problems) {
    const schedules = problem.submissionableSchedules;
    const key = schedules
      .map((s) => s.name)
      .sort()
      .join(",");

    let group = groupMap.get(key);
    if (group == null) {
      group = {
        problems: [],
        entries: schedules.slice().sort((a, b) => startAtMs(a) - startAtMs(b)),
        hasSubmittableProblem: false,
      };
      groupMap.set(key, group);
    }
    if (problem.submissionStatus?.isSubmittable) {
      group.hasSubmittableProblem = true;
    }
    group.problems.push(problem);
  }

  return Array.from(groupMap.entries())
    .map(([key, { problems, entries, hasSubmittableProblem }]) => ({
      key,
      schedules: entries.map((e) => ({
        name: e.name,
        temporalStatus: getTemporalStatus(e, now),
      })),
      problems,
      hasSubmittableProblem,
      _sortKey: entries.length > 0 ? startAtMs(entries[0]) : Infinity,
    }))
    .sort((a, b) => {
      if (a.hasSubmittableProblem !== b.hasSubmittableProblem) {
        return a.hasSubmittableProblem ? -1 : 1;
      }
      return a._sortKey - b._sortKey;
    });
}
