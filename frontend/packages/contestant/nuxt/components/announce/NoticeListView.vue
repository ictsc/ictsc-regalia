<script setup lang="ts">
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { Notice } from "~/features/models";

defineProps<{
  notices?: Notice[];
  readSlugs: readonly string[];
  pending?: boolean;
  error?: unknown;
}>();
defineEmits<{ toggleRead: [slug: string]; markAllRead: []; retry: [] }>();
</script>
<template>
  <header class="linked-page-heading">
    <h1 class="page-title">お知らせ</h1>
  </header>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><slot name="settings" />
    <div class="actions">
      <button class="button-secondary" @click="$emit('markAllRead')">
        すべて既読にする
      </button>
    </div>
    <article v-for="n in notices" :key="n.slug" class="notice-entry">
      <time>{{ new Date(n.effectiveFrom).toLocaleString("ja-JP") }}</time>
      <div>
        <h2>
          <NuxtLink :to="`/announces/${n.slug}`">{{ n.title }}</NuxtLink>
        </h2>
        <MarkdownContent :source="n.markdown" /><button
          class="button-secondary"
          @click="$emit('toggleRead', n.slug)"
        >
          {{ readSlugs.includes(n.slug) ? "未読に戻す" : "既読にする" }}
        </button>
      </div>
    </article>
    <p v-if="!notices?.length">現在お知らせはありません。</p></RequestState
  >
</template>
