import type { ApiClient } from "./client";
import { expectData, expectNoContent } from "./error";

export async function impersonateContestant(client: ApiClient, name: string) {
  expectNoContent(await client.POST("/api/v1/admin/impersonations", {
    body: { contestant_name: name },
  }));
  const { viewer } = expectData(await client.GET("/api/v1/viewer", { cache: "no-store" }));
  if (viewer.state !== "CONTESTANT" || viewer.profile.name !== name || !viewer.impersonated_by) {
    throw new Error("代理ログインのセッションを確認できませんでした。ページを再読み込みして再試行してください。");
  }
  return viewer;
}
