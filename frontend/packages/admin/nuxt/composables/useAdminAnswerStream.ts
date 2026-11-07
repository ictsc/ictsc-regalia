import { subscribeAdminAnswerEvents, type AdminAnswerChange } from "@ictsc/api";
import { onMounted, onUnmounted, ref } from "vue";

/** まとめて届いた変更を 1 回の再取得にまとめる待ち時間。 */
const refreshDelayMs = 300;

/**
 * 回答の提出・ワークフロー更新・採点を SSE で受け、onChange を呼ぶ。
 * 連続した変更はまとめて 1 回にする。再接続時(2 回目以降の ready)は
 * 取りこぼしがあり得るので changes を空にして呼ぶ。
 * 戻り値の streamError は接続が切れて再接続を待っている間 true になる。
 */
export function useAdminAnswerStream(
  onChange: (changes: AdminAnswerChange[]) => void,
) {
  const streamError = ref(false);
  let stop: (() => void) | undefined;
  let timer: ReturnType<typeof setTimeout> | undefined;
  let pending: AdminAnswerChange[] = [];
  let reconnected = false;
  let connected = false;

  function schedule() {
    clearTimeout(timer);
    timer = setTimeout(() => {
      const changes = reconnected ? [] : pending;
      pending = [];
      reconnected = false;
      onChange(changes);
    }, refreshDelayMs);
  }

  onMounted(() => {
    stop = subscribeAdminAnswerEvents(
      (message) => {
        streamError.value = false;
        if (message.type === "answer") {
          pending.push(message.change);
          schedule();
          return;
        }
        // 初回の ready はページの初期取得と重なるので何もしない。
        if (connected) {
          reconnected = true;
          schedule();
        }
        connected = true;
      },
      () => {
        streamError.value = true;
      },
    );
  });
  onUnmounted(() => {
    clearTimeout(timer);
    stop?.();
  });
  return { streamError };
}
