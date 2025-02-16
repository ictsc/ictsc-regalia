// app/features/viewer/signup.test.ts
import { describe, expect, it, vitest } from "vitest";
import { http, HttpResponse } from "msw";
import { setupMSW } from "@app/__test__/msw/node";
import { signUp } from "./signup";

const server = setupMSW();

describe("signUp", () => {
  it("signs up", async () => {
    const fn = vitest.fn();

    server.use(
      http.post("http://example.test/api/auth/signup", async ({ request }) => {
        fn(await request.json());
        return HttpResponse.json({}, { status: 201 });
      }),
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
