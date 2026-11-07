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
  return subscribeSSE<{ snapshot: Snapshot; deployment: Update }>(
    url,
    {
      snapshot: (data) =>
        onMessage({ type: "snapshot", deployments: data.deployments }),
      deployment: (data) =>
        onMessage({ type: "deployment", deployment: data.deployment }),
    },
    onError,
  );
}

export type ContestantAnswerChange =
  components["schemas"]["ContestantAnswerChange"];
export type AdminAnswerChange = components["schemas"]["AdminAnswerChange"];

/**
 * ready は接続・再接続のたびに届く。取りこぼしを回復するため、受け取ったら
 * 回答の一覧を再取得する。answer は変更された回答の識別子だけを運ぶ。
 */
export type AnswerStreamMessage<Change> =
  | { type: "ready" }
  | { type: "answer"; change: Change };

export function subscribeContestantAnswerEvents(
  onMessage: (message: AnswerStreamMessage<ContestantAnswerChange>) => void,
  onError?: (event: Event) => void,
): () => void {
  return subscribeAnswerEvents(
    "/api/v1/contestant/answers/stream",
    onMessage,
    onError,
  );
}

export function subscribeAdminAnswerEvents(
  onMessage: (message: AnswerStreamMessage<AdminAnswerChange>) => void,
  onError?: (event: Event) => void,
): () => void {
  return subscribeAnswerEvents(
    "/api/v1/admin/answers/stream",
    onMessage,
    onError,
  );
}

function subscribeAnswerEvents<Change>(
  url: string,
  onMessage: (message: AnswerStreamMessage<Change>) => void,
  onError?: (event: Event) => void,
): () => void {
  return subscribeSSE<{ ready: unknown; answer: Change }>(
    url,
    {
      ready: () => onMessage({ type: "ready" }),
      answer: (change) => onMessage({ type: "answer", change }),
    },
    onError,
  );
}

/** 名前付き SSE イベントごとに JSON を解釈して渡す。戻り値で接続を閉じる。 */
function subscribeSSE<Events extends Record<string, unknown>>(
  url: string,
  handlers: { [Name in keyof Events]: (data: Events[Name]) => void },
  onError?: (event: Event) => void,
): () => void {
  const source = new EventSource(url, { withCredentials: true });
  for (const name of Object.keys(handlers) as (keyof Events & string)[]) {
    source.addEventListener(name, ((event: MessageEvent<string>) => {
      let data: Events[typeof name];
      try {
        data = JSON.parse(event.data) as Events[typeof name];
      } catch (error) {
        console.error(`Invalid SSE ${name} event`, error);
        return;
      }
      handlers[name](data);
    }) as EventListener);
  }
  source.onerror = (event) => onError?.(event);
  return () => source.close();
}
