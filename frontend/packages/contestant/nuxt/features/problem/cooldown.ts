import { demoProblemState } from "./demo";

export function remainingCooldownSeconds(
  nextSubmittableAt: string | number | undefined,
  now: number,
): number {
  if (nextSubmittableAt == null) return 0;
  const next =
    typeof nextSubmittableAt === "number"
      ? nextSubmittableAt
      : Date.parse(nextSubmittableAt);
  if (!Number.isFinite(next)) return 0;
  return Math.max(0, Math.ceil((next - now) / 1000));
}

export function remainingCooldownMinutes(
  nextSubmittableAt: string | number | undefined,
  now: number,
): number {
  return Math.ceil(remainingCooldownSeconds(nextSubmittableAt, now) / 60);
}

export function demoCooldownSeconds(problemCode: string): number {
  return demoProblemState(problemCode).cooldownSeconds;
}
