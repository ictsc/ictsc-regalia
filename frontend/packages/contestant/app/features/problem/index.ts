import { type Transport, createClient } from "@connectrpc/connect";
import { ProblemService, type Problem } from "@ictsc/proto/contestant/v1";

export async function fetchProblems(transport: Transport): Promise<Problem[]> {
  const client = createClient(ProblemService, transport);
  const problems = await client.listProblems({});
  return problems.problems;
}
