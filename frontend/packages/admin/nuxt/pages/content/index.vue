<script setup lang="ts">
import ContentView from "~/components/content/ContentView.vue";
import { api, expectData } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
useHead({ title: "コンテンツ" });
async function settled<T>(p: Promise<T>) {
  try {
    return { data: await p, error: undefined };
  } catch (error) {
    return { data: undefined, error };
  }
}
const { data, error, pending, refresh } = await useAsyncData(
  "admin-content",
  async () => {
    const [status, sections, problems, announcements] = await Promise.all([
      api.GET("/api/v1/admin/content/status").then(expectData),
      settled(api.GET("/api/v1/admin/sections").then(expectData)),
      settled(api.GET("/api/v1/admin/problems").then(expectData)),
      settled(api.GET("/api/v1/admin/announcements").then(expectData)),
    ]);
    return { status: status.content, sections, problems, announcements };
  },
);
const { busy, message, run } = useMutation(refresh);
const actions = createAdminActions(api);
const refreshContent = (commit: string) =>
  run("コンテンツ更新", () => actions.refreshContent(commit));
</script>
<template>
  <main class="workspace">
    <ContentView
      :status="data?.status"
      :sections="data?.sections ?? {}"
      :problems="data?.problems ?? {}"
      :announcements="data?.announcements ?? {}"
      :pending="pending"
      :error="error"
      :busy="busy"
      :message="message"
      @refresh-content="refreshContent"
      @retry="refresh"
    />
  </main>
</template>
