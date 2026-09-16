/** Presentation-only samples; never stored as answers or marks. */
export function demoProblemState(problemCode: string) {
  const suffix = problemCode.match(/\d+$/)?.[0];
  const index = suffix
    ? Number(suffix.slice(-6))
    : Array.from(problemCode).reduce(
        (sum, char) => sum + char.charCodeAt(0),
        0,
      );
  const samples = [
    { status: "unanswered", cooldownSeconds: 0 },
    { status: "complete", cooldownSeconds: 0 },
    { status: "partial", cooldownSeconds: 0 },
    { status: "unanswered", cooldownSeconds: 15 * 60 },
  ] as const;
  return samples[index % samples.length]!;
}
