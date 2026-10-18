import type { AdminRealtimeEvent } from "@ictsc/api";

/**
 * Which shared AsyncData keys to refetch for an event. Pages with their own
 * listener (the grading board and detail) update themselves in place.
 */
export function adminRefreshKeys(
  event: AdminRealtimeEvent,
): string[] | "all" | null {
  switch (event.type) {
    case "ready":
      return event.reconnected ? "all" : null;
    case "answer":
      return event.reason === "SUBMITTED"
        ? ["admin-overview"]
        : event.reason === "WORKFLOW_UPDATED"
          ? null
          : ["admin-scores"];
    case "scores":
    case "ranking":
      return ["admin-scores"];
    case "content":
      return ["admin-content", "admin-overview"];
    case "teams":
      return ["admin-teams", "admin-contestants", "admin-overview"];
    case "competition":
      return ["settings", "admin-scores"];
  }
}
