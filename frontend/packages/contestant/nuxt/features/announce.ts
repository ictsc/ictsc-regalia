import { expectData, type ApiClient } from "@ictsc/api";
import type { Notice } from "./models";

export type { Notice } from "./models";

export async function fetchNotices(client: ApiClient): Promise<Notice[]> {
  const response = expectData(
    await client.GET("/api/v1/contestant/announcements"),
  );
  return response.announcements.map((announcement) => ({
    slug: announcement.slug,
    title: announcement.title,
    markdown: announcement.markdown,
    effectiveFrom: announcement.effective_from,
    body: announcement.markdown,
  }));
}
