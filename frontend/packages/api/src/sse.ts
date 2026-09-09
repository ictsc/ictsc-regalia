import type { components } from "./schema";

export type ContestantDeployment =
  components["schemas"]["ContestantDeployment"];
export type AdminDeployment = components["schemas"]["AdminDeployment"];
type ContestantSnapshot =
  components["schemas"]["ContestantDeploymentsResponse"];
type ContestantUpdate = components["schemas"]["ContestantDeploymentResponse"];
type AdminSnapshot = components["schemas"]["AdminDeploymentsResponse"];
type AdminUpdate = components["schemas"]["AdminDeploymentResponse"];

export type ContestantDeploymentStreamMessage =
  | { type: "snapshot"; deployments: ContestantSnapshot["deployments"] }
  | { type: "deployment"; deployment: ContestantUpdate["deployment"] };

export type AdminDeploymentStreamMessage =
  | { type: "snapshot"; deployments: AdminSnapshot["deployments"] }
  | { type: "deployment"; deployment: AdminUpdate["deployment"] };

export function subscribeContestantDeploymentEvents(
  url: string,
  onMessage: (message: ContestantDeploymentStreamMessage) => void,
  onError?: (event: Event) => void,
): () => void {
  return subscribeDeploymentEvents<ContestantSnapshot, ContestantUpdate>(
    url,
    onMessage,
    onError,
  );
}

export function subscribeAdminDeploymentEvents(
  url: string,
  onMessage: (message: AdminDeploymentStreamMessage) => void,
  onError?: (event: Event) => void,
): () => void {
  return subscribeDeploymentEvents<AdminSnapshot, AdminUpdate>(
    url,
    onMessage,
    onError,
  );
}

function subscribeDeploymentEvents<
  Snapshot extends ContestantSnapshot | AdminSnapshot,
  Update extends ContestantUpdate | AdminUpdate,
>(
  url: string,
  onMessage: (
    message:
      | { type: "snapshot"; deployments: Snapshot["deployments"] }
      | { type: "deployment"; deployment: Update["deployment"] },
  ) => void,
  onError?: (event: Event) => void,
): () => void {
  const source = new EventSource(url, { withCredentials: true });
  const parse = <Payload>(event: MessageEvent<string>): Payload | null => {
    try {
      return JSON.parse(event.data) as Payload;
    } catch (error) {
      console.error("Invalid deployment SSE event", error);
      return null;
    }
  };
  source.addEventListener("snapshot", ((event: MessageEvent<string>) => {
    const data = parse<Snapshot>(event);
    if (data != null) {
      onMessage({ type: "snapshot", deployments: data.deployments });
    }
  }) as EventListener);
  source.addEventListener("deployment", ((event: MessageEvent<string>) => {
    const data = parse<Update>(event);
    if (data != null) {
      onMessage({ type: "deployment", deployment: data.deployment });
    }
  }) as EventListener);
  source.onerror = (event) => onError?.(event);
  return () => source.close();
}
