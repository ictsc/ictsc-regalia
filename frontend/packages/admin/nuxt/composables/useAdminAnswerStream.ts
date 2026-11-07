import { subscribeAdminAnswerEvents, type AdminAnswerChange } from "@ictsc/api";
import { onMounted, onUnmounted, ref } from "vue";

/** まとめて届いた変更を 1 回の再取得にまとめる待ち時間。 */
const refreshDelayMs = 300;

/**
 * 回答の提出・ワークフロー更新・採点を SSE で受け、onChange を呼ぶ。
 * 連続した変更はまとめて 1 回にする。接続・再接続時(ready)は
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
        // ready は初回接続でも再取得する。ページの初期取得が終わってから
        // 購読が始まるまでの間や、切断中の変更を取りこぼさないため。
        reconnected = true;
        schedule();
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
