import {
  remainingCooldownMinutes,
  remainingCooldownSeconds,
} from "../features/problem/cooldown";

const storageKey = "ictsc-demo-problem-cooldowns";

type CooldownSnapshots = Record<string, number>;

function loadSnapshots(): CooldownSnapshots {
  if (typeof sessionStorage === "undefined") return {};
  try {
    const parsed: unknown = JSON.parse(
      sessionStorage.getItem(storageKey) ?? "{}",
    );
    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed))
      return {};
    return Object.fromEntries(
      Object.entries(parsed).filter(
        ([, seconds]) =>
          typeof seconds === "number" &&
          Number.isFinite(seconds) &&
          seconds >= 0,
      ),
    );
  } catch {
    return {};
  }
}

function saveSnapshots(snapshots: CooldownSnapshots) {
  if (typeof sessionStorage === "undefined") return;
  try {
    sessionStorage.setItem(storageKey, JSON.stringify(snapshots));
  } catch {
    // sessionStorageが無効でも、画面を開いている間はuseStateで固定する。
  }
}

function targetTime(nextSubmittableAt: string | number | undefined) {
  if (typeof nextSubmittableAt === "number") return nextSubmittableAt;
  return nextSubmittableAt ? Date.parse(nextSubmittableAt) : Number.NaN;
}

export function useProblemCooldown() {
  const demoMode = useDemoMode();
  const snapshots = useState<CooldownSnapshots>(
    "demo-problem-cooldowns",
    loadSnapshots,
  );

  const remainingSeconds = (
    problemCode: string,
    nextSubmittableAt: string | number | undefined,
    now: number,
  ) => {
    if (!demoMode.value)
      return remainingCooldownSeconds(nextSubmittableAt, now);

    const target = targetTime(nextSubmittableAt);
    if (!Number.isFinite(target)) return 0;
    const key = `${problemCode}:${target}`;
    const saved = snapshots.value[key];
    if (saved != null) return saved;

    const seconds = remainingCooldownSeconds(target, now);
    snapshots.value = { ...snapshots.value, [key]: seconds };
    saveSnapshots(snapshots.value);
    return seconds;
  };

  const remainingMinutes = (
    problemCode: string,
    nextSubmittableAt: string | number | undefined,
    now: number,
  ) => {
    if (!demoMode.value)
      return remainingCooldownMinutes(nextSubmittableAt, now);
    return Math.ceil(
      remainingSeconds(problemCode, nextSubmittableAt, now) / 60,
    );
  };

  return { remainingMinutes, remainingSeconds };
}
