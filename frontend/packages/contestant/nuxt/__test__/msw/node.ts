import { afterAll, afterEach, beforeAll } from "vitest";
import type { RequestHandler, SharedOptions } from "msw";
import { setupServer } from "msw/node";

const server = setupServer();

beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
afterAll(() => server.close());
afterEach(() => server.resetHandlers());

export function setupMSW(
  handlers: readonly RequestHandler[] = [],
  _options: Partial<SharedOptions> = { onUnhandledRequest: "error" },
) {
  server.use(...handlers);
  return server;
}
