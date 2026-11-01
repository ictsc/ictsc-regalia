<script setup lang="ts">
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import AdminPage from "~/components/layout/AdminPage.vue";
import AdminPanel from "~/components/layout/AdminPanel.vue";
import AdminFacts from "~/components/layout/AdminFacts.vue";
import type { components } from "@ictsc/api";

defineProps<{
  announcement?: components["schemas"]["Announcement"];
  pending?: boolean;
  error?: unknown;
}>();
defineEmits<{ retry: [] }>();
</script>
<template>
  <AdminPage
    :title="announcement?.title ?? 'お知らせ'"
    :back="{ to: '/content', label: 'コンテンツ一覧' }"
  >
    <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
      ><AdminPanel v-if="announcement" title="本文">
        <AdminFacts
          :items="[
            {
              label: '公開日時',
              value: new Date(announcement.effective_from).toLocaleString(
                'ja-JP',
              ),
            },
          ]" />
        <MarkdownContent :source="announcement.markdown" /></AdminPanel
    ></RequestState>
  </AdminPage>
</template>
