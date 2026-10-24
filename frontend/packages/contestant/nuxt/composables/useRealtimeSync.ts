import { spreadRefetch } from "@ictsc/api";
import { contestantRefreshPlan } from "~/features/realtime";

/** Refetches the shared screen data named by each realtime event. */
export function useRealtimeSync() {
  return useRealtime((event) => {
    const plan = contestantRefreshPlan(event);
    if (!plan) return;
    const refetch = () =>
      plan.keys === "all" ? refreshNuxtData() : refreshNuxtData(plan.keys);
    if (plan.spread) spreadRefetch(refetch);
    else void refetch();
  });
}
