import {
  Code,
  ConnectError,
  type Transport,
  createClient,
} from "@connectrpc/connect";
import {
  ProblemService,
  type Problem as ProtoProblem,
  type SubmissionStatus,
} from "@ictsc/proto/contestant/v1";

// TODO: アプリケーションが利用する値を定義する
export type Problem = ProtoProblem;

export async function fetchProblems(transport: Transport): Promise<Problem[]> {
  const client = createClient(ProblemService, transport);
  const problems = await client.listProblems({});
  return problems.problems;
}

export type ProblemDetail = {
  code: string;
  title: string;
  maxScore: number;
  redeployable: boolean;
  penaltyThreashold: number;
  body: string;
  submissionStatus?: SubmissionStatus;
};

export async function fetchProblem(
  transport: Transport,
  code: string,
): Promise<ProblemDetail> {
  const client = createClient(ProblemService, transport);
  const { problem } = await client.getProblem({ code });
  if (problem == null) {
    throw new ConnectError("Problem not found", Code.NotFound);
  }
  const body = problem.body?.body.value?.body;
  if (body == null) {
    throw new ConnectError("Problem body not found", Code.NotFound);
  }
  return {
    code: problem.code,
    title: problem.title,
    maxScore: problem.maxScore,
    redeployable: problem.deployment?.redeployable ?? false,
    penaltyThreashold: problem.deployment?.penaltyThreashold ?? 0,
    body,
    submissionStatus: problem.submissionStatus,
  };
}
