<script setup lang="ts">
import { api } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
const route = useRoute(),
  commit = ref(String(route.query.commit ?? ""));
const { data, error, pending, refresh } = await useAsyncData(
  () => `admin-problem:${route.params.code}:${route.query.commit ?? ""}`,
  () =>
    createAdminActions(api).getProblem(
      String(route.params.code),
      String(route.query.commit ?? ""),
    ),
);
useHead({ title: computed(() => data.value?.problem.title ?? "問題詳細") });
</script>
<template>
  <main class="workspace">
    <h1>{{ data?.problem.code }}: {{ data?.problem.title }}</h1>
    <form
      class="form-grid"
      @submit.prevent="
        navigateTo({ path: route.path, query: commit ? { commit } : {} })
      "
    >
      <label>参照commit（空欄でactive）<input v-model="commit" /></label
      ><button class="button-primary">表示</button>
    </form>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><template v-if="data"
        ><p>
          {{ data.problem.category }} / {{ data.problem.section_slug }} / 満点
          {{ data.problem.max_score }}
        </p>
        <p class="mono">{{ data.problem.content_commit }}</p>
        <p>
          再展開ルール: {{ data.problem.redeploy_rule.type }} / 減点開始回数
          {{ data.problem.redeploy_rule.penalty_threshold ?? "—" }} /
          {{ data.problem.redeploy_rule.penalty_percentage ?? "—" }}%
        </p>
        <h2>問題本文</h2>
        <MarkdownContent :source="data.problem.body" />
        <h2>問題解説</h2>
        <MarkdownContent :source="data.problem.explanation" /></template
    ></RequestState>
  </main>
</template>
