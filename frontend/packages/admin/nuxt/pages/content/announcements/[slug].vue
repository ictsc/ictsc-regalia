<script setup lang="ts">
import { api } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
const route = useRoute();
const { data, error, pending, refresh } = await useAsyncData(
  () => `admin-announcement:${route.params.slug}`,
  () => createAdminActions(api).getAnnouncement(String(route.params.slug)),
);
useHead({
  title: computed(() => data.value?.announcement.title ?? "お知らせ"),
});
</script>
<template>
  <main class="workspace">
    <NuxtLink class="text-link" to="/content">← コンテンツ一覧</NuxtLink>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><template v-if="data"
        ><h1>{{ data.announcement.title }}</h1>
        <p class="admin-lead">
          公開日時:
          <time>{{
            new Date(data.announcement.effective_from).toLocaleString("ja-JP")
          }}</time>
        </p>
        <section class="admin-panel">
          <MarkdownContent
            :source="data.announcement.markdown"
          /></section></template
    ></RequestState>
  </main>
</template>
