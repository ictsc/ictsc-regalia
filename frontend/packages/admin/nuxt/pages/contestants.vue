<script setup lang="ts">
import ContestantsView from "~/components/contestant/ContestantsView.vue";
import { api, expectData } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
useHead({ title: "参加者" });
const { data, error, pending, refresh } = await useAsyncData(
  "admin-contestants",
  async () => expectData(await api.GET("/api/v1/admin/contestants")),
);
const { busy, message, run } = useMutation();
async function impersonate(name: string) {
  if (
    await run(
      "代理ログイン",
      () => createAdminActions(api).impersonate(name),
      `${name} として代理ログインしますか？`,
    )
  )
    window.location.assign("/");
}
</script>
<template>
  <main class="workspace">
    <ContestantsView
      :contestants="data?.contestants"
      :pending="pending"
      :error="error"
      :busy="busy"
      :message="message"
      @impersonate="impersonate"
      @retry="refresh"
    />
  </main>
</template>
