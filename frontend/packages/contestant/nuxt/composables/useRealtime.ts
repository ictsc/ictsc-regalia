import {
  createRealtimeChannel,
  subscribeContestantRealtimeEvents,
  type ContestantRealtimeEvent,
} from "@ictsc/api";

// One EventSource per tab, shared by every mounted listener.
const channel = createRealtimeChannel<ContestantRealtimeEvent>(
  subscribeContestantRealtimeEvents,
);

/**
 * Receives realtime change notifications while the calling component is
 * mounted. `failing` is true while the browser is reconnecting.
 */
export function useRealtime(handler: (event: ContestantRealtimeEvent) => void) {
  const failing = useState("contestant-realtime-failing", () => false);
  let stop: (() => void) | undefined;
  let stopStatus: (() => void) | undefined;
  onMounted(() => {
    stopStatus = channel.onStatus((value) => {
      failing.value = value;
    });
    stop = channel.listen(handler);
  });
  onUnmounted(() => {
    stop?.();
    stopStatus?.();
  });
  return { failing };
}
