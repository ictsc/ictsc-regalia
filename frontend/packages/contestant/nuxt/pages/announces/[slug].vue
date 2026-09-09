<script setup lang="ts">
import { api } from "@ictsc/api";
import { fetchNotices } from "~/features/announce";
import { useReadAnnouncements } from "~/features/announce-read-status";
const { markAsRead } = useReadAnnouncements();
const route = useRoute();
const { data, error, pending, refresh } = await useAsyncData("notices", () =>
  fetchNotices(api),
);
const notice = computed(() =>
  data.value?.find((n) => n.slug === route.params.slug),
);
watch(
  notice,
  (value) => {
    if (value) markAsRead(value.slug);
  },
  { immediate: true },
);
useHead({ title: computed(() => notice.value?.title ?? "お知らせ") });
</script>
<template>
  <main class="workspace">
    <NuxtLink to="/announces">← お知らせ一覧へ</NuxtLink
    ><RequestState :error="error" :pending="pending" @retry="refresh"
      ><template v-if="notice"
        ><h1>{{ notice.title }}</h1>
        <MarkdownContent :source="notice.markdown"
      /></template>
      <p v-else>お知らせが見つかりません。</p></RequestState
    >
  </main>
</template>
