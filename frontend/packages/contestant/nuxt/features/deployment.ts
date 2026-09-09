import {
  expectData,
  subscribeContestantDeploymentEvents,
  type ApiClient,
  type ContestantDeploymentStreamMessage,
} from "@ictsc/api";
import { DeploymentStatus, type DeploymentStatus as Status } from "./models";

export { DeploymentStatus } from "./models";

export type Deployment = {
  revision: number;
  status: Status;
  requestedAt: string;
  allowedDeploymentCount: number;
  thresholdExceeded: boolean;
  penalty: number;
  contentCommit?: string;
};

export async function fetchDeployments(
  client: ApiClient,
  code: string,
): Promise<Deployment[]> {
  const response = expectData(
    await client.GET("/api/v1/contestant/problems/{problem_code}/deployments", {
      params: { path: { problem_code: code } },
    }),
  );
  const items = response.deployments.map(mapDeployment);
  items.sort((a, b) => b.revision - a.revision);
  return items;
}

export async function deploy(
  client: ApiClient,
  code: string,
): Promise<Deployment> {
  const response = expectData(
    await client.POST(
      "/api/v1/contestant/problems/{problem_code}/deployments",
      { params: { path: { problem_code: code } } },
    ),
  );
  return mapDeployment(response.deployment);
}

export function subscribeDeployments(
  problemCode: string,
  onMessage: (message: ContestantDeploymentStreamMessage) => void,
  onError?: (event: Event) => void,
): () => void {
  return subscribeContestantDeploymentEvents(
    `/api/v1/contestant/problems/${encodeURIComponent(problemCode)}/deployments/stream`,
    onMessage,
    onError,
  );
}

export function mapDeployment(deployment: {
  revision: number;
  status: string;
  requested_at: string;
  penalty: number;
  allowed_request_count: number;
  content_commit?: string;
}): Deployment {
  const allowed = deployment.allowed_request_count;
  return {
    revision: deployment.revision,
    status: mapStatus(deployment.status),
    requestedAt: deployment.requested_at,
    allowedDeploymentCount: allowed,
    thresholdExceeded: allowed === 0,
    penalty: deployment.penalty,
    contentCommit: deployment.content_commit,
  };
}

function mapStatus(status: string): Status {
  if (status === "COMPLETED") return DeploymentStatus.COMPLETED;
  if (status === "FAILED") return DeploymentStatus.FAILED;
  if (status === "DEPLOYING") return DeploymentStatus.DEPLOYING;
  return DeploymentStatus.QUEUED;
}
