import { describe, expect, it, vi } from "vitest";
import type { ApiClient } from "./client";
import { impersonateContestant } from "./impersonation";

describe("impersonation handoff", () => {
  for (const [name, viewer, success] of [
    ["verified", { state: "CONTESTANT", profile: { name: "alice" }, impersonated_by: "staff" }, true],
    ["missing cookie", { state: "ANONYMOUS" }, false],
    ["stale user", { state: "CONTESTANT", profile: { name: "bob" }, impersonated_by: "staff" }, false],
    ["ordinary session", { state: "CONTESTANT", profile: { name: "alice" } }, false],
  ] as const) {
    it(name, async () => {
      const calls: string[] = [];
      const POST = vi.fn(async () => { calls.push("POST"); return { response: new Response(null, {status: 204}) }; });
      const GET = vi.fn(async () => { calls.push("GET"); return { data: { viewer }, response: new Response() }; });
      const client = { POST, GET } as unknown as ApiClient;
      const result = impersonateContestant(client, "alice");
      if (success) await expect(result).resolves.toEqual(viewer);
      else await expect(result).rejects.toThrow("セッションを確認できませんでした");
      expect(calls).toEqual(["POST", "GET"]);
      expect(GET).toHaveBeenCalledWith("/api/v1/viewer", { cache: "no-store" });
    });
  }
});
