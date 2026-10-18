<script setup lang="ts">
import ProblemContentView from "~/components/content/ProblemContentView.vue";
import { api } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
const route = useRoute();

const { data, error, pending, refresh } = await useAsyncData(
  () => `admin-problem:${route.params.code}:${route.query.commit ?? ""}`,
  () =>
    createAdminActions(api).getProblem(
      String(route.params.code),
      String(route.query.commit ?? ""),
    ),
);
useHead({ title: computed(() => data.value?.problem.title ?? "問題詳細") });
const viewCommit = (commit: string) =>
  navigateTo({ path: route.path, query: commit ? { commit } : {} });
</script>
<template>
  <main class="workspace">
    <ProblemContentView
      :problem="data?.problem"
      :initial-commit="String(route.query.commit ?? '')"
      :pending="pending"
      :error="error"
      @view-commit="viewCommit"
      @retry="refresh"
    />
  </main>
</template>
