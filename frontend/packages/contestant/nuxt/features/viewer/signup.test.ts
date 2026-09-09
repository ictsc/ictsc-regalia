// app/features/viewer/signup.test.ts
import { describe, expect, it, vitest } from "vitest";
import { http, HttpResponse } from "msw";
import { setupMSW } from "../../__test__/msw/node";
import { signUp } from "./signup";

const server = setupMSW();

describe("signUp", () => {
  it("signs up", async () => {
    const fn = vitest.fn();

    server.use(
      http.post(
        "http://example.test/api/v1/auth/signup",
        async ({ request }) => {
          fn(await request.json());
          return new HttpResponse(null, { status: 204 });
        },
      ),
    );

    const result = await signUp(
      {
        invitationCode: "test",
        name: "test",
        displayName: "test",
      },
      "http://example.test",
    );

    expect(fn).toHaveBeenCalledWith({
      invitation_code: "test",
      name: "test",
      display_name: "test",
    });
    expect(result).toEqual({});
  });
});
