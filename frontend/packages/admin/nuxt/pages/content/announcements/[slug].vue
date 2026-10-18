<script setup lang="ts">
import AnnouncementContentView from "~/components/content/AnnouncementContentView.vue";
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
    <AnnouncementContentView
      :announcement="data?.announcement"
      :pending="pending"
      :error="error"
      @retry="refresh"
    />
  </main>
</template>
