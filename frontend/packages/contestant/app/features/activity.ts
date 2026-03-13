import type { Transport } from "@connectrpc/connect";
import type { Score as ProtoScore } from "@ictsc/proto/contestant/v1";
import { fetchProblems } from "./problem";
import { fetchAnswers } from "./answer";

export type ActivityEntry = {
  readonly problemCode: string;
  readonly problemTitle: string;
  readonly maxScore: number;
  readonly answerId: number;
  readonly submittedAt: string;
  readonly score?: ProtoScore;
};

export async function fetchActivity(
  transport: Transport,
): Promise<ActivityEntry[]> {
  const problems = await fetchProblems(transport);
  const results = await Promise.all(
    problems.map(async (problem) => {
      const { answers } = await fetchAnswers(transport, problem.code);
      return answers.map((answer) => ({
        problemCode: problem.code,
        problemTitle: problem.title,
        maxScore: problem.maxScore,
        answerId: answer.id,
        submittedAt: answer.submittedAt,
        score: answer.score,
      }));
    }),
  );
  const entries = results.flat();
  entries.sort(
    (a, b) =>
      new Date(b.submittedAt).getTime() - new Date(a.submittedAt).getTime(),
  );
  return entries;
}
