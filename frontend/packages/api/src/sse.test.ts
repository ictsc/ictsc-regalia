import { afterEach, describe, expect, it, vi } from "vitest";
import {
  subscribeAdminAnswerEvents,
  subscribeContestantAnswerEvents,
  subscribeContestantDeploymentEvents,
  type AdminAnswerChange,
  type AnswerStreamMessage,
  type ContestantAnswerChange,
  type ContestantDeployment,
  type ContestantDeploymentStreamMessage,
} from "./sse";

class FakeEventSource {
  static latest: FakeEventSource;
  readonly listeners = new Map<string, EventListener>();
  readonly withCredentials: boolean;
  closed = false;
  onerror: ((event: Event) => void) | null = null;
  readonly url: string;
  constructor(url: string | URL, init?: EventSourceInit) {
    this.url = String(url);
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

describe("subscribeAnswerEvents", () => {
  afterEach(() => vi.unstubAllGlobals());

  it("delivers ready and answer events from the contestant stream", () => {
    vi.stubGlobal("EventSource", FakeEventSource);
    const messages: AnswerStreamMessage<ContestantAnswerChange>[] = [];
    const unsubscribe = subscribeContestantAnswerEvents((message) =>
      messages.push(message),
    );
    expect(FakeEventSource.latest.url).toBe(
      "/api/v1/contestant/answers/stream",
    );
    expect(FakeEventSource.latest.withCredentials).toBe(true);
    FakeEventSource.latest.emit("ready", {});
    FakeEventSource.latest.emit("answer", {
      problem_code: "P1",
      answer_number: 2,
    });
    expect(messages).toEqual([
      { type: "ready" },
      { type: "answer", change: { problem_code: "P1", answer_number: 2 } },
    ]);
    unsubscribe();
    expect(FakeEventSource.latest.closed).toBe(true);
  });

  it("ignores malformed data and reports stream errors", () => {
    vi.stubGlobal("EventSource", FakeEventSource);
    vi.spyOn(console, "error").mockImplementation(() => {});
    const messages: AnswerStreamMessage<AdminAnswerChange>[] = [];
    const onError = vi.fn();
    subscribeAdminAnswerEvents((message) => messages.push(message), onError);
    expect(FakeEventSource.latest.url).toBe("/api/v1/admin/answers/stream");
    FakeEventSource.latest.listeners.get("answer")?.(
      new MessageEvent("answer", { data: "{" }),
    );
    FakeEventSource.latest.onerror?.(new Event("error"));
    expect(messages).toEqual([]);
    expect(onError).toHaveBeenCalledOnce();
  });
});
