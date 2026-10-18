import type { components } from "./schema";

type ProblemEvent = components["schemas"]["ContestantRealtimeProblemEvent"];
type AnswerEvent = components["schemas"]["AdminRealtimeAnswerEvent"];

/** Events that name no target; the client refetches the matching screens. */
type NoticeEvent =
  | { type: "ranking" }
  | { type: "content" }
  | { type: "teams" }
  | { type: "competition" };

/**
 * `ready` arrives on every (re)connection. `reconnected` is false for the
 * first ready of a subscription, when the screen data was just loaded.
 */
type ReadyEvent = { type: "ready"; reconnected: boolean };

export type ContestantRealtimeEvent =
  | ReadyEvent
  | NoticeEvent
  | { type: "answers"; problemCode: string }
  | { type: "scores"; problemCode: string };

export type AdminRealtimeAnswerReason = AnswerEvent["reason"];

export type AdminRealtimeEvent =
  | ReadyEvent
  | NoticeEvent
  | { type: "scores" }
  | {
      type: "answer";
      reason: AdminRealtimeAnswerReason;
      teamCode: number;
      problemCode: string;
      answerNumber: number;
    };

const noticeTypes = ["ranking", "content", "teams", "competition"] as const;

export function subscribeContestantRealtimeEvents(
  onEvent: (event: ContestantRealtimeEvent) => void,
  onError?: (event: Event) => void,
  url = "/api/v1/contestant/events/stream",
): () => void {
  return subscribeRealtime(url, onEvent, onError, (listen) => {
    for (const type of noticeTypes) listen(type, () => ({ type }));
    for (const type of ["answers", "scores"] as const)
      listen(type, (data: ProblemEvent) => ({
        type,
        problemCode: data.problem_code,
      }));
  });
}

export function subscribeAdminRealtimeEvents(
  onEvent: (event: AdminRealtimeEvent) => void,
  onError?: (event: Event) => void,
  url = "/api/v1/admin/events/stream",
): () => void {
  return subscribeRealtime(url, onEvent, onError, (listen) => {
    for (const type of [...noticeTypes, "scores"] as const)
      listen(type, () => ({ type }));
    listen("answer", (data: AnswerEvent) => ({
      type: "answer",
      reason: data.reason,
      teamCode: data.answer.team_code,
      problemCode: data.answer.problem_code,
      answerNumber: data.answer.answer_number,
    }));
  });
}

function subscribeRealtime<Event extends { type: string }>(
  url: string,
  onEvent: (event: Event | ReadyEvent) => void,
  onError: ((event: globalThis.Event) => void) | undefined,
  register: (
    listen: <Data>(name: string, map: (data: Data) => Event) => void,
  ) => void,
): () => void {
  const source = new EventSource(url, { withCredentials: true });
  let connected = false;
  const listen = <Data>(name: string, map: (data: Data) => Event) => {
    source.addEventListener(name, ((message: MessageEvent<string>) => {
      try {
        onEvent(map(JSON.parse(message.data) as Data));
      } catch (error) {
        console.error("Invalid realtime SSE event", error);
      }
    }) as EventListener);
  };
  source.addEventListener("ready", () => {
    onEvent({ type: "ready", reconnected: connected });
    connected = true;
  });
  register(listen);
  source.onerror = (event) => onError?.(event);
  return () => source.close();
}

/**
 * Shares one EventSource among every listener of an app. The connection opens
 * with the first listener and closes when the last one unsubscribes.
 */
export function createRealtimeChannel<Event>(
  subscribe: (
    onEvent: (event: Event) => void,
    onError: (event: globalThis.Event) => void,
  ) => () => void,
) {
  const listeners = new Set<(event: Event) => void>();
  const statusListeners = new Set<(failing: boolean) => void>();
  let close: (() => void) | undefined;
  const setFailing = (failing: boolean) => {
    for (const listener of statusListeners) listener(failing);
  };
  return {
    listen(listener: (event: Event) => void): () => void {
      listeners.add(listener);
      close ??= subscribe(
        (event) => {
          setFailing(false);
          for (const current of [...listeners]) current(event);
        },
        () => setFailing(true),
      );
      return () => {
        listeners.delete(listener);
        if (listeners.size === 0) {
          close?.();
          close = undefined;
        }
      };
    },
    onStatus(listener: (failing: boolean) => void): () => void {
      statusListeners.add(listener);
      return () => statusListeners.delete(listener);
    },
  };
}

/**
 * Delays a refetch triggered by an event every client receives at once, so
 * hundreds of browsers do not hit the API in the same instant.
 */
export function spreadRefetch(
  refetch: () => unknown,
  maxDelayMs = 3000,
  random: () => number = Math.random,
): ReturnType<typeof setTimeout> {
  return setTimeout(() => void refetch(), Math.floor(random() * maxDelayMs));
}
