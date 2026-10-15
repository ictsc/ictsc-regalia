import { afterEach, describe, expect, it, vi } from "vitest";
import {
  subscribeContestantDeploymentEvents,
  type ContestantDeployment,
  type ContestantDeploymentStreamMessage,
} from "./sse";

class FakeEventSource {
  static latest: FakeEventSource;
  readonly listeners = new Map<string, EventListener>();
  readonly withCredentials: boolean;
  closed = false;
  onerror: ((event: Event) => void) | null = null;
  constructor(_url: string | URL, init?: EventSourceInit) {
    this.withCredentials = init?.withCredentials ?? false;
    FakeEventSource.latest = this;
  }
  addEventListener(type: string, listener: EventListener) {
    this.listeners.set(type, listener);
  }
  close() {
    this.closed = true;
  }
  emit(type: string, data: unknown) {
    this.listeners.get(type)?.(
      new MessageEvent(type, { data: JSON.stringify(data) }),
    );
  }
}

describe("subscribeDeploymentEvents", () => {
  afterEach(() => vi.unstubAllGlobals());

  it("delivers reconnect snapshots and deployment updates with credentials", () => {
    vi.stubGlobal("EventSource", FakeEventSource);
    const messages: ContestantDeploymentStreamMessage[] = [];
    const unsubscribe = subscribeContestantDeploymentEvents(
      "/stream",
      (message) => messages.push(message),
    );
    const queued: ContestantDeployment = {
      revision: 1,
      status: "QUEUED",
      requested_at: "2026-08-30T00:00:00Z",
      penalty: 0,
      allowed_request_count: 2,
      content_commit: "1111111111111111111111111111111111111111",
    };
    const completed: ContestantDeployment = {
      ...queued,
      revision: 2,
      status: "COMPLETED",
    };
    FakeEventSource.latest.emit("snapshot", { deployments: [queued] });
    FakeEventSource.latest.emit("deployment", { deployment: completed });
    expect(FakeEventSource.latest.withCredentials).toBe(true);
    expect(messages).toEqual([
      { type: "snapshot", deployments: [queued] },
      { type: "deployment", deployment: completed },
    ]);
    unsubscribe();
    expect(FakeEventSource.latest.closed).toBe(true);
  });
});
