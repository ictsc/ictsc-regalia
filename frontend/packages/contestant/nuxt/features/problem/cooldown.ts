export function remainingCooldownSeconds(
  nextSubmittableAt: string | undefined,
  now: number,
): number {
  if (!nextSubmittableAt) return 0;
  const next = Date.parse(nextSubmittableAt);
  if (!Number.isFinite(next)) return 0;
  return Math.max(0, Math.ceil((next - now) / 1000));
}

export function remainingCooldownMinutes(
  nextSubmittableAt: string | undefined,
  now: number,
): number {
  return Math.ceil(remainingCooldownSeconds(nextSubmittableAt, now) / 60);
}
