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

/** Stable demo samples: two answerable problems followed by two waiting ones. */
export function demoCooldownSeconds(problemCode: string): number {
  const suffix = problemCode.match(/\d+$/)?.[0];
  const index = suffix
    ? Number(suffix.slice(-6))
    : Array.from(problemCode).reduce(
        (sum, char) => sum + char.charCodeAt(0),
        0,
      );
  return [0, 0, 8 * 60, 15 * 60][index % 4]!;
}
