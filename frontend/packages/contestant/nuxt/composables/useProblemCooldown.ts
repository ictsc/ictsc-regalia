import { remainingCooldownSeconds } from "../features/problem/cooldown";

export function useProblemCooldown() {
  const remainingSeconds = (
    _problemCode: string,
    nextSubmittableAt: string | number | undefined,
    now: number,
  ) => remainingCooldownSeconds(nextSubmittableAt, now);
  const remainingMinutes = (
    problemCode: string,
    nextSubmittableAt: string | number | undefined,
    now: number,
  ) => Math.ceil(remainingSeconds(problemCode, nextSubmittableAt, now) / 60);

  return { remainingMinutes, remainingSeconds };
}
