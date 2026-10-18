<script setup lang="ts">
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { components } from "@ictsc/api";

defineProps<{
  announcement?: components["schemas"]["Announcement"];
  pending?: boolean;
  error?: unknown;
}>();
defineEmits<{ retry: [] }>();
</script>
<template>
  <NuxtLink class="text-link" to="/content">← コンテンツ一覧</NuxtLink>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><template v-if="announcement"
      ><h1>{{ announcement.title }}</h1>
      <p class="admin-lead">
        公開日時:
        <time>{{
          new Date(announcement.effective_from).toLocaleString("ja-JP")
        }}</time>
      </p>
      <section class="admin-panel">
        <MarkdownContent :source="announcement.markdown" /></section></template
  ></RequestState>
</template>
