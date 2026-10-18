<script setup lang="ts">
import ScoresView from "~/components/score/ScoresView.vue";
import { api, expectData } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
useHead({ title: "得点・順位" });
const { data, error, pending, refresh } = await useAsyncData(
  "admin-scores",
  async () => {
    const [scores, ranking] = await Promise.all([
      api.GET("/api/v1/admin/scores").then(expectData),
      api.GET("/api/v1/admin/ranking").then(expectData),
    ]);
    return { scores: scores.scores, ranking };
  },
);
const { busy, message, run } = useMutation(refresh);
const actions = createAdminActions(api);
const recalculate = () =>
  run("得点再計算", actions.recalculate, "すべての得点を再計算しますか？");
const reveal = () =>
  run("最終得点公開", actions.reveal, {
    message:
      "凍結と遅延を解除して最終得点を公開します。この操作は取り消せません。続行しますか？",
    destructive: true,
  });
</script>
<template>
  <main class="workspace">
    <ScoresView
      :ranking="data?.ranking"
      :scores="data?.scores"
      :pending="pending"
      :error="error"
      :busy="busy"
      :message="message"
      @recalculate="recalculate"
      @reveal="reveal"
      @retry="refresh"
    />
  </main>
</template>
