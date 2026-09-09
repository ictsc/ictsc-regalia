<script setup lang="ts">
import { api } from "@ictsc/api";
import { fetchNotices } from "~/features/announce";
import { useReadAnnouncements } from "~/features/announce-read-status";
const { isRead, markAsRead, markAsUnread, markAllAsRead } =
  useReadAnnouncements();
useHead({ title: "お知らせ" });
const { data, error, pending, refresh } = await useAsyncData("notices", () =>
  fetchNotices(api),
);
</script>
<template>
  <main>
    <header class="linked-page-heading">
      <NuxtLink to="/problems">← 問題一覧へ</NuxtLink>
      <h1>お知らせ</h1>
      <p>競技運営からのお知らせ</p>
    </header>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><div class="actions">
        <button class="button-secondary" @click="markAllAsRead(data ?? [])">
          すべて既読にする
        </button>
      </div>
      <article v-for="n in data" :key="n.slug" class="notice-entry">
        <time>{{ new Date(n.effectiveFrom).toLocaleString("ja-JP") }}</time>
        <div>
          <h2>
            <NuxtLink :to="`/announces/${n.slug}`">{{ n.title }}</NuxtLink>
          </h2>
          <MarkdownContent :source="n.markdown" /><button
            class="button-secondary"
            @click="isRead(n.slug) ? markAsUnread(n.slug) : markAsRead(n.slug)"
          >
            {{ isRead(n.slug) ? "未読に戻す" : "既読にする" }}
          </button>
        </div>
      </article>
      <p v-if="!data?.length">現在お知らせはありません。</p></RequestState
    >
  </main>
</template>
