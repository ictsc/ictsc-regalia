import type { ContestantRealtimeEvent } from "@ictsc/api";

/**
 * Which shared AsyncData keys to refetch for an event. `spread` marks events
 * every contestant receives at once, whose refetch is randomly delayed.
 */
export type RefreshPlan = { keys: string[] | "all"; spread: boolean };

export function contestantRefreshPlan(
  event: ContestantRealtimeEvent,
): RefreshPlan | null {
  switch (event.type) {
    case "ready":
      return event.reconnected ? { keys: "all", spread: true } : null;
    case "answers":
    case "scores":
      return { keys: ["competition", "activity"], spread: false };
    case "ranking":
      return { keys: ["competition", "ranking"], spread: true };
    case "content":
      return { keys: ["competition", "notices"], spread: true };
    case "teams":
      return { keys: ["teams", "profile"], spread: true };
    case "competition":
      return { keys: "all", spread: true };
  }
}

/** Whether the problem detail page for `problemCode` must refetch. */
export function affectsProblem(
  event: ContestantRealtimeEvent,
  problemCode: string,
): boolean {
  return (
    ((event.type === "answers" || event.type === "scores") &&
      event.problemCode === problemCode) ||
    event.type === "content"
  );
}
