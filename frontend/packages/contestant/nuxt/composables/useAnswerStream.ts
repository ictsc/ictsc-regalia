import {
  subscribeContestantAnswerEvents,
  type ContestantAnswerChange,
} from "@ictsc/api";

/** undefined は再接続を表し、取りこぼした変更があり得るので全体を再取得する。 */
type AnswerChangeListener = (
  change: ContestantAnswerChange | undefined,
) => void;

const listeners = new Set<AnswerChangeListener>();

/** まとめて届いた変更を 1 回の再取得にまとめる待ち時間。 */
const refreshDelayMs = 300;

/**
 * 自チームの回答の変更を SSE で受け、問題一覧・サイドバー・提出履歴・
 * クールダウンの元データ(competition / activity)を再取得する。
 * レイアウトで一度だけ呼ぶ。チームメイトが別の端末で提出しても、
 * 再読み込みなしで「提出済み」に変わる。
 */
export function useAnswerStream() {
  const { viewer } = useSession();
  let unsubscribe: (() => void) | undefined;
  let timer: ReturnType<typeof setTimeout> | undefined;
  let connected = false;

  function scheduleRefresh() {
    clearTimeout(timer);
    timer = setTimeout(() => {
      void refreshNuxtData(["competition", "activity"]);
    }, refreshDelayMs);
  }

  function notify(change: ContestantAnswerChange | undefined) {
    scheduleRefresh();
    for (const listener of listeners) listener(change);
  }

  function connect() {
    connected = false;
    unsubscribe = subscribeContestantAnswerEvents((message) => {
      if (message.type === "answer") {
        notify(message.change);
        return;
      }
      // 初回の ready はページの初期取得と重なるので何もしない。
      if (connected) notify(undefined);
      connected = true;
    });
  }

  function disconnect() {
    clearTimeout(timer);
    unsubscribe?.();
    unsubscribe = undefined;
  }

  onMounted(() => {
    watch(
      () => viewer.value?.state === "CONTESTANT",
      (contestant) => {
        disconnect();
        if (contestant) connect();
      },
      { immediate: true },
    );
  });
  onUnmounted(disconnect);
}

/** 回答の変更を購読する。アンマウント時に自動で解除する。 */
export function onAnswerChange(listener: AnswerChangeListener) {
  listeners.add(listener);
  onUnmounted(() => listeners.delete(listener));
}
