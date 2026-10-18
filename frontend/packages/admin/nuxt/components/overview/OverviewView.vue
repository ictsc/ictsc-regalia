<script setup lang="ts">
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { components } from "@ictsc/api";

export type OverviewCard = {
  label: string;
  detail: string;
  to: string;
  /** Item count; undefined while unknown. */
  count?: number;
  /** Set when this card's request failed. */
  error?: unknown;
};

defineProps<{
  cards: OverviewCard[];
  content?: components["schemas"]["ContentStatus"];
  contentError?: unknown;
  pending?: boolean;
  error?: unknown;
}>();
defineEmits<{ retry: [] }>();

const stateLabels: Record<string, string> = {
  NEVER: "未配信",
  REFRESHING: "更新中",
  SUCCEEDED: "配信中",
  FAILED: "更新失敗",
};
</script>
<template>
  <h1>概要</h1>
  <p class="admin-lead">競技運営の状況を確認し、各管理画面へ移動できます。</p>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><div class="admin-summary-grid">
      <NuxtLink
        v-for="card in cards"
        :key="card.to"
        :to="card.to"
        class="admin-summary-card"
      >
        <span class="admin-summary-label">{{ card.label }}</span>
        <span class="admin-summary-number">{{
          card.error ? "—" : (card.count ?? 0)
        }}</span>
        <span class="admin-summary-footer"
          >{{ card.error ? "取得できませんでした" : card.detail }} →</span
        >
      </NuxtLink>
    </div>
    <section class="admin-panel">
      <div class="admin-panel-heading">
        <h2>配信コンテンツ</h2>
        <NuxtLink class="text-link" to="/content">管理画面へ →</NuxtLink>
      </div>
      <RequestState :error="contentError" @retry="$emit('retry')">
        <div class="admin-meta">
          <span
            class="admin-badge"
            :class="{
              'is-active': content?.state === 'SUCCEEDED',
            }"
            >{{ stateLabels[content?.state ?? ""] ?? "状態不明" }}</span
          >
          <span v-if="content?.serving_last_known_good" class="admin-badge"
            >最終正常版を配信中</span
          >
        </div>
        <p class="admin-code">
          有効な commit:
          {{ content?.active_commit ?? "未設定" }}
        </p>
      </RequestState>
    </section></RequestState
  >
</template>
