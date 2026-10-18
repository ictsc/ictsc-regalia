import { adminRefreshKeys } from "~/features/realtime";

/** Refetches the shared screen data named by each realtime event. */
export function useRealtimeSync() {
  return useRealtime((event) => {
    const keys = adminRefreshKeys(event);
    if (keys === "all") void refreshNuxtData();
    else if (keys) void refreshNuxtData(keys);
  });
}
