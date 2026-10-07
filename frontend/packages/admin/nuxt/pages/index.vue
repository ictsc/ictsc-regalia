<script setup lang="ts">
import { api, expectData } from "@ictsc/api";
useHead({ title: "概要" });
async function settled<T>(p: Promise<T>) {
  try {
    return { data: await p, error: undefined };
  } catch (error) {
    return { data: undefined, error };
  }
}
const { data, error, pending, refresh } = await useAsyncData(
  "admin-overview",
  async () => {
    const [teams, contestants, answers, deployments, content] =
      await Promise.all([
        settled(api.GET("/api/v1/admin/teams").then(expectData)),
        settled(api.GET("/api/v1/admin/contestants").then(expectData)),
        settled(api.GET("/api/v1/admin/answers").then(expectData)),
        settled(api.GET("/api/v1/admin/deployments").then(expectData)),
        settled(api.GET("/api/v1/admin/content/status").then(expectData)),
      ]);
    return { teams, contestants, answers, deployments, content };
  },
);
const cards = computed(() =>
  data.value
    ? [
        {
          label: "チーム",
          detail: "登録済みチームを管理",
          to: "/teams",
          count: data.value.teams.data?.teams.length,
          error: data.value.teams.error,
        },
        {
          label: "参加者",
          detail: "参加者と所属を確認",
          to: "/contestants",
          count: data.value.contestants.data?.contestants.length,
          error: data.value.contestants.error,
        },
        {
          label: "回答",
          detail: "提出された回答を確認・採点",
          to: "/submissions",
          count: data.value.answers.data?.answers.length,
          error: data.value.answers.error,
        },
        {
          label: "進行中の再展開",
          detail: "環境の準備状況を確認",
          to: "/deployments",
          count: data.value.deployments.data?.deployments.filter((d) =>
            ["QUEUED", "DEPLOYING"].includes(d.latest_status),
          ).length,
          error: data.value.deployments.error,
        },
      ]
    : [],
);
const stateLabels: Record<string, string> = {
  NEVER: "未配信",
  REFRESHING: "更新中",
  SUCCEEDED: "配信中",
  FAILED: "更新失敗",
};
</script>
<template>
  <main class="workspace">
    <h1>概要</h1>
    <p class="admin-lead">競技運営の状況を確認し、各管理画面へ移動できます。</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
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
        <RequestState :error="data?.content.error" @retry="refresh">
          <div class="admin-meta">
            <span
              class="admin-badge"
              :class="{
                'is-active': data?.content.data?.content.state === 'SUCCEEDED',
              }"
              >{{
                stateLabels[data?.content.data?.content.state ?? ""] ??
                "状態不明"
              }}</span
            >
            <span
              v-if="data?.content.data?.content.serving_last_known_good"
              class="admin-badge"
              >最終正常版を配信中</span
            >
          </div>
          <p class="admin-code">
            有効な commit:
            {{ data?.content.data?.content.active_commit ?? "未設定" }}
          </p>
        </RequestState>
      </section></RequestState
    >
  </main>
</template>
