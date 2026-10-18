<script setup lang="ts">
import { api, expectData } from "@ictsc/api";
import OverviewView, {
  type OverviewCard,
} from "~/components/overview/OverviewView.vue";
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
const cards = computed<OverviewCard[]>(() =>
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
</script>
<template>
  <main class="workspace">
    <OverviewView
      :cards="cards"
      :content="data?.content.data?.content"
      :content-error="data?.content.error"
      :pending="pending"
      :error="error"
      @retry="refresh"
    />
  </main>
</template>
