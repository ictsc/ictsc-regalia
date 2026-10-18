import { afterEach, describe, expect, it, vi } from "vitest";
import {
  createRealtimeChannel,
  spreadRefetch,
  subscribeAdminRealtimeEvents,
  subscribeContestantRealtimeEvents,
  type AdminRealtimeEvent,
  type ContestantRealtimeEvent,
} from "./realtime";

class FakeEventSource {
  static instances: FakeEventSource[] = [];
  readonly listeners = new Map<string, EventListener>();
  readonly withCredentials: boolean;
  closed = false;
  onerror: ((event: Event) => void) | null = null;
  constructor(
    readonly url: string | URL,
    init?: EventSourceInit,
  ) {
    this.withCredentials = init?.withCredentials ?? false;
    FakeEventSource.instances.push(this);
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

describe("realtime SSE", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
    vi.useRealTimers();
    FakeEventSource.instances = [];
  });

  it("maps contestant events and flags reconnect readiness", () => {
    vi.stubGlobal("EventSource", FakeEventSource);
    const events: ContestantRealtimeEvent[] = [];
    const stop = subscribeContestantRealtimeEvents((e) => events.push(e));
    const source = FakeEventSource.instances[0]!;
    expect(source.url).toBe("/api/v1/contestant/events/stream");
    expect(source.withCredentials).toBe(true);
    source.emit("ready", {});
    source.emit("answers", { problem_code: "ABC" });
    source.emit("scores", { problem_code: "ABC" });
    source.emit("ranking", {});
    source.emit("ready", {});
    expect(events).toEqual([
      { type: "ready", reconnected: false },
      { type: "answers", problemCode: "ABC" },
      { type: "scores", problemCode: "ABC" },
      { type: "ranking" },
      { type: "ready", reconnected: true },
    ]);
    stop();
    expect(source.closed).toBe(true);
  });

  it("maps admin answer references to camelCase", () => {
    vi.stubGlobal("EventSource", FakeEventSource);
    const events: AdminRealtimeEvent[] = [];
    subscribeAdminRealtimeEvents((e) => events.push(e));
    FakeEventSource.instances[0]!.emit("answer", {
      reason: "WORKFLOW_UPDATED",
      answer: { team_code: 2, problem_code: "ABC", answer_number: 3 },
    });
    expect(events).toEqual([
      {
        type: "answer",
        reason: "WORKFLOW_UPDATED",
        teamCode: 2,
        problemCode: "ABC",
        answerNumber: 3,
      },
    ]);
  });

  it("shares one connection and closes it with the last listener", () => {
    vi.stubGlobal("EventSource", FakeEventSource);
    const channel = createRealtimeChannel<ContestantRealtimeEvent>(
      subscribeContestantRealtimeEvents,
    );
    const first: ContestantRealtimeEvent[] = [];
    const second: ContestantRealtimeEvent[] = [];
    const failing: boolean[] = [];
    channel.onStatus((value) => failing.push(value));
    const stopFirst = channel.listen((e) => first.push(e));
    const stopSecond = channel.listen((e) => second.push(e));
    expect(FakeEventSource.instances).toHaveLength(1);
    const source = FakeEventSource.instances[0]!;
    source.onerror?.(new Event("error"));
    source.emit("teams", {});
    expect(first).toEqual([{ type: "teams" }]);
    expect(second).toEqual([{ type: "teams" }]);
    expect(failing).toEqual([true, false]);
    stopFirst();
    expect(source.closed).toBe(false);
    stopSecond();
    expect(source.closed).toBe(true);
  });

  it("spreads broadcast refetches over the delay window", () => {
    vi.useFakeTimers();
    const refetch = vi.fn();
    spreadRefetch(refetch, 3000, () => 0.5);
    vi.advanceTimersByTime(1499);
    expect(refetch).not.toHaveBeenCalled();
    vi.advanceTimersByTime(1);
    expect(refetch).toHaveBeenCalledOnce();
  });
});
