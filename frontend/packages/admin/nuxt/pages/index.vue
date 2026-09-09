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
          to: "/teams",
          count: data.value.teams.data?.teams.length,
          error: data.value.teams.error,
        },
        {
          label: "参加者",
          to: "/contestants",
          count: data.value.contestants.data?.contestants.length,
          error: data.value.contestants.error,
        },
        {
          label: "未採点回答",
          to: "/submissions",
          count: data.value.answers.data?.answers.length,
          error: data.value.answers.error,
        },
        {
          label: "進行中の再展開",
          to: "/deployments",
          count: data.value.deployments.data?.deployments.filter((d) =>
            ["QUEUED", "DEPLOYING"].includes(d.latest_status),
          ).length,
          error: data.value.deployments.error,
        },
      ]
    : [],
);
</script>
<template>
  <main class="workspace">
    <h1>概要</h1>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><section v-for="card in cards" :key="card.to" class="content-section">
        <h2>{{ card.label }}</h2>
        <RequestState :error="card.error" @retry="refresh"
          ><p class="mono">{{ card.count ?? 0 }}</p>
          <NuxtLink class="text-link" :to="card.to"
            >詳細 →</NuxtLink
          ></RequestState
        >
      </section>
      <h2>コンテンツ</h2>
      <RequestState :error="data?.content.error" @retry="refresh"
        ><p>{{ data?.content.data?.content.state }}</p>
        <p>{{ data?.content.data?.content.active_commit ?? "未取得" }}</p>
        <p v-if="data?.content.data?.content.serving_last_known_good">
          最終正常版を配信中
        </p>
        <NuxtLink class="text-link" to="/content"
          >コンテンツ管理</NuxtLink
        ></RequestState
      ></RequestState
    >
  </main>
</template>
