export type RefetchSchedulerOptions = {
  /** 変更通知を受けてから再取得するまでの待ち時間。連続した変更をまとめる。 */
  changeDelayMs?: number;
  /**
   * 接続・再接続(ready)時の再取得を 0 からこの時間までの間でばらす。
   * 競技開始時や API の再起動で全員が同時に取り直すのを避けるため。
   */
  readyJitterMs?: number;
  random?: () => number;
};

export type RefetchScheduler<Change> = {
  /** 個別の変更。changeDelayMs 後にまとめて再取得する。 */
  change: (change: Change) => void;
  /** 接続・再接続。取りこぼしがあり得るので、ばらした時間の後に全体を再取得する。 */
  ready: () => void;
  /** 予定している再取得を取り消す。 */
  cancel: () => void;
};

/**
 * SSE の変更通知をまとめて再取得する。refetch には、まとめた変更か、
 * 全体を取り直すべきときは null を渡す。予定が重なったときは早い方に合わせる。
 */
export function createRefetchScheduler<Change>(
  refetch: (changes: Change[] | null) => void,
  {
    changeDelayMs = 300,
    readyJitterMs = 5000,
    random = Math.random,
  }: RefetchSchedulerOptions = {},
): RefetchScheduler<Change> {
  let timer: ReturnType<typeof setTimeout> | undefined;
  let dueAt = Number.POSITIVE_INFINITY;
  let changes: Change[] = [];
  let full = false;

  function schedule(delayMs: number) {
    const due = Date.now() + delayMs;
    if (timer !== undefined && due >= dueAt) return;
    clearTimeout(timer);
    dueAt = due;
    timer = setTimeout(fire, delayMs);
  }

  function fire() {
    const batch = full ? null : changes;
    timer = undefined;
    dueAt = Number.POSITIVE_INFINITY;
    changes = [];
    full = false;
    refetch(batch);
  }

  return {
    change(change) {
      changes.push(change);
      schedule(changeDelayMs);
    },
    ready() {
      full = true;
      schedule(Math.floor(random() * readyJitterMs));
    },
    cancel() {
      clearTimeout(timer);
      timer = undefined;
      dueAt = Number.POSITIVE_INFINITY;
      changes = [];
      full = false;
    },
  };
}
