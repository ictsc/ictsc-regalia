import {
  createRefetchScheduler,
  subscribeAdminAnswerEvents,
  type AdminAnswerChange,
} from "@ictsc/api";
import { onMounted, onUnmounted, ref } from "vue";

/**
 * 回答の提出・ワークフロー更新・採点を SSE で受け、onChange を呼ぶ。
 * 連続した変更はまとめて 1 回にする。接続・再接続時(ready)は
 * 取りこぼしがあり得るので、0〜5 秒ばらしてから changes を空にして呼ぶ。
 * 戻り値の streamError は接続が切れて再接続を待っている間 true になる。
 */
export function useAdminAnswerStream(
  onChange: (changes: AdminAnswerChange[]) => void,
) {
  const streamError = ref(false);
  let stop: (() => void) | undefined;
  const scheduler = createRefetchScheduler<AdminAnswerChange>((changes) =>
    onChange(changes ?? []),
  );

  onMounted(() => {
    stop = subscribeAdminAnswerEvents(
      (message) => {
        streamError.value = false;
        if (message.type === "answer") scheduler.change(message.change);
        else scheduler.ready();
      },
      () => {
        streamError.value = true;
      },
    );
  });
  onUnmounted(() => {
    scheduler.cancel();
    stop?.();
  });
  return { streamError };
}
