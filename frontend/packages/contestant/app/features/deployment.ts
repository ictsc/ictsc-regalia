import { timestampDate } from "@bufbuild/protobuf/wkt";
import { type Transport, createClient } from "@connectrpc/connect";
import { ProblemService, DeploymentStatus } from "@ictsc/proto/contestant/v1";

export type Deployment = {
  revision: number;
  status: DeploymentStatus;
  requestedAt: string;
  allowedDeploymentCount: number;
  thresholdExceeded: boolean;
  penalty: number;
};

export async function fetchDeployments(
  transport: Transport,
  code: string,
): Promise<Deployment[]> {
  const client = createClient(ProblemService, transport);
  const { deployments } = await client.listDeployments({ code });
  const items = deployments.map((d) => ({
    revision: d.revision,
    status: d.status,
    requestedAt:
      d.requestedAt != null ? timestampDate(d.requestedAt).toISOString() : "",
    allowedDeploymentCount: d.allowedRequestCount,
    thresholdExceeded: d.allowedRequestCount < 0,
    penalty: d.penalty,
  }));
  items.sort((a, b) => b.revision - a.revision);
  return items;
}

export async function deploy(transport: Transport, code: string): Promise<void> {
  const client = createClient(ProblemService, transport);
  await client.deploy({ code });
}
