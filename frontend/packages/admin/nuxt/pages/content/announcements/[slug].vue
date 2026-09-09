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
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><template v-if="data"
        ><h1>{{ data.announcement.title }}</h1>
        <time>{{ data.announcement.effective_from }}</time
        ><MarkdownContent :source="data.announcement.markdown" /></template
    ></RequestState>
  </main>
</template>
