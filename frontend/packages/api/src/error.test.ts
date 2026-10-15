import { describe, expect, it } from "vitest";
import { ApiError, expectData } from "./error";

describe("expectData", () => {
  it("preserves RFC 9457 code and Retry-After", () => {
    const response = new Response(null, {
      status: 429,
      headers: { "Retry-After": "42" },
    });
    let thrown: unknown;
    try {
      expectData({
        response,
        error: {
          type: "about:blank",
          title: "Too Many Requests",
          status: 429,
          code: "answer_interval",
        },
      });
    } catch (error) {
      thrown = error;
    }
    expect(thrown).toBeInstanceOf(ApiError);
    expect(thrown).toMatchObject({
      status: 429,
      code: "answer_interval",
      retryAfterSeconds: 42,
    });
  });
});
