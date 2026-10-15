import type { Page, Route } from "@playwright/test";

export type Deferred = {
  readonly promise: Promise<void>;
  resolve: () => void;
};

export function deferred(): Deferred {
  let resolve: () => void = () => undefined;
  const promise = new Promise<void>((done) => {
    resolve = done;
  });
  return { promise, resolve };
}

export function fulfillJson(
  route: Route,
  value: unknown,
  status = 200,
): Promise<void> {
  return route.fulfill({
    status,
    contentType: "application/json; charset=utf-8",
    body: JSON.stringify(value),
  });
}

export async function installEventSourceMock(page: Page): Promise<void> {
  await page.addInitScript(() => {
    type MockSource = EventTarget & {
      readonly url: string;
      readonly withCredentials: boolean;
      closed: boolean;
      readyState: number;
      onerror: ((event: Event) => void) | null;
      close: () => void;
    };
    const sources: MockSource[] = [];
    class MockEventSource extends EventTarget {
      static readonly CONNECTING = 0;
      static readonly OPEN = 1;
      static readonly CLOSED = 2;

      readonly CONNECTING = 0;
      readonly OPEN = 1;
      readonly CLOSED = 2;
      readonly url: string;
      readonly withCredentials: boolean;
      closed = false;
      readyState = MockEventSource.OPEN;
      onerror: ((event: Event) => void) | null = null;
      onmessage: ((event: MessageEvent) => void) | null = null;
      onopen: ((event: Event) => void) | null = null;

      constructor(url: string | URL, init?: EventSourceInit) {
        super();
        this.url = String(url);
        this.withCredentials = init?.withCredentials ?? false;
        sources.push(this);
      }

      close() {
        this.closed = true;
        this.readyState = MockEventSource.CLOSED;
      }
    }

    const target = window as unknown as {
      EventSource: typeof EventSource;
      __e2eEventSources: MockSource[];
      __emitE2eSse: (url: string, event: string, data: unknown) => number;
    };
    target.EventSource = MockEventSource as unknown as typeof EventSource;
    target.__e2eEventSources = sources;
    target.__emitE2eSse = (url, event, data) => {
      let delivered = 0;
      for (const source of sources) {
        if (source.closed || source.url !== url) continue;
        source.dispatchEvent(
          new MessageEvent(event, { data: JSON.stringify(data) }),
        );
        delivered += 1;
      }
      return delivered;
    };
  });
}

export async function emitSse(
  page: Page,
  url: string,
  event: string,
  data: unknown,
): Promise<void> {
  await page.waitForFunction((streamUrl) => {
    const target = window as unknown as {
      __e2eEventSources?: Array<{ url: string; closed: boolean }>;
    };
    return target.__e2eEventSources?.some(
      (source) => !source.closed && source.url === streamUrl,
    );
  }, url);
  const delivered = await page.evaluate(
    ({ streamUrl, eventName, payload }) => {
      const target = window as unknown as {
        __emitE2eSse: (
          sourceUrl: string,
          sourceEvent: string,
          sourceData: unknown,
        ) => number;
      };
      return target.__emitE2eSse(streamUrl, eventName, payload);
    },
    { streamUrl: url, eventName: event, payload: data },
  );
  if (delivered === 0) {
    throw new Error(`No active EventSource for ${url}`);
  }
}
