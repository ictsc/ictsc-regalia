import {
  Code,
  ConnectError,
  type Transport,
  createClient,
} from "@connectrpc/connect";
import {
  ProblemService,
  type Problem as ProtoProblem,
  DeploymentStatus,
} from "@ictsc/proto/contestant/v1";
import { Deployment } from "./deployment";
import { timestampDate } from "@bufbuild/protobuf/wkt";

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
  deployment: Deployment;
  body: string;
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
  if (problem.deployment == null) {
    throw new ConnectError("Problem deployment not found", Code.NotFound);
  }

  return {
    code: problem.code,
    title: problem.title,
    deployment: {
      events: problem.deployment.events
        .map((event) => {
          return {
            revision: event.revision != null ? Number(event.revision) : 0,
            occuredAt:
              event.occurredAt != null
                ? timestampDate(event.occurredAt).toISOString()
                : "",
            totalPenalty: event.penalty != null ? Number(event.penalty) : 0,
            type:
              event.type != null
                ? DeploymentStatus[event.type] || DeploymentStatus[event.type]
                : "-",
            isDeploying: event.type === DeploymentStatus.DEPLOYING,
          };
        })
        .sort(
          (a, b) =>
            new Date(b.occuredAt).getTime() - new Date(a.occuredAt).getTime(),
        ),
      redeployable: problem.deployment.redeployable ?? false,
      // TODO: 実装する
      maxRedeployment: 0,
    },
    body,
  };
}
