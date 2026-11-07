import {
  createRefetchScheduler,
  subscribeContestantAnswerEvents,
  type ContestantAnswerChange,
} from "@ictsc/api";

/**
 * 回答が変わった問題のコードを受け取る。まとめた変更ごとに、問題 1 つに
 * つき 1 回だけ呼ぶ。undefined は接続(再接続を含む)を表す。初期取得から
 * 接続までや切断中の変更は届かないので、受け取ったら全体を再取得する。
 */
type AnswerChangeListener = (problemCode: string | undefined) => void;

const listeners = new Set<AnswerChangeListener>();

/**
 * 自チームの回答の変更を SSE で受け、問題一覧・サイドバー・提出履歴・
 * クールダウンの元データ(competition / activity)を再取得する。
 * レイアウトで一度だけ呼ぶ。チームメイトが別の端末で提出しても、
 * 再読み込みなしで「提出済み」に変わる。
 * 接続・再接続時の再取得は 0〜5 秒ばらし、全員が同時に取り直さないようにする。
 */
export function useAnswerStream() {
  const { viewer } = useSession();
  let unsubscribe: (() => void) | undefined;
  const scheduler = createRefetchScheduler<ContestantAnswerChange>(
    (changes) => {
      void refreshNuxtData(["competition", "activity"]);
      if (changes === null) {
        for (const listener of listeners) listener(undefined);
        return;
      }
      const problems = new Set(changes.map((change) => change.problem_code));
      for (const problem of problems)
        for (const listener of listeners) listener(problem);
    },
  );

  function connect() {
    // ready は初回接続でも再取得する。ページの初期取得が終わってから
    // 購読が始まるまでの間の提出を取りこぼさないため。
    unsubscribe = subscribeContestantAnswerEvents((message) => {
      if (message.type === "answer") scheduler.change(message.change);
      else scheduler.ready();
    });
  }

  function disconnect() {
    scheduler.cancel();
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
