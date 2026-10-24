import {
  createRealtimeChannel,
  subscribeAdminRealtimeEvents,
  type AdminRealtimeEvent,
} from "@ictsc/api";

// One EventSource per tab, shared by every mounted listener.
const channel = createRealtimeChannel<AdminRealtimeEvent>(
  subscribeAdminRealtimeEvents,
);

/**
 * Receives realtime change notifications while the calling component is
 * mounted. `failing` is true while the browser is reconnecting.
 */
export function useRealtime(handler: (event: AdminRealtimeEvent) => void) {
  const failing = useState("admin-realtime-failing", () => false);
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
