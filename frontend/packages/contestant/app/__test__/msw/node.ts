import { afterAll, afterEach, beforeAll } from "vitest";
import type { RequestHandler, SharedOptions } from "msw";
import { setupServer } from "msw/node";

export function setupMSW(
  handlers: readonly RequestHandler[] = [],
  options: Partial<SharedOptions> = { onUnhandledRequest: "error" },
) {
  const server = setupServer(...handlers);
  beforeAll(() => server.listen(options));
  afterAll(() => server.close());
  afterEach(() => server.resetHandlers());
  return server;
}
