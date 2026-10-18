<script setup lang="ts">
import SettingsView from "~/components/settings/SettingsView.vue";
import { api, expectData } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
useHead({ title: "競技設定" });
const { data, error, pending, refresh } = await useAsyncData(
  "settings",
  async () => {
    const [rule, schedule] = await Promise.all([
      api.GET("/api/v1/admin/rule").then(expectData),
      api.GET("/api/v1/admin/dashboard-schedule").then(expectData),
    ]);
    return { rule: rule.rule, schedule: schedule.dashboard_schedule };
  },
);
const initialFreeze = data.value?.schedule.ranking_freeze_at
  ? new Date(
      Date.parse(data.value.schedule.ranking_freeze_at) -
        new Date().getTimezoneOffset() * 60000,
    )
      .toISOString()
      .slice(0, 16)
  : "";
const { busy, message, run } = useMutation(refresh);
const actions = createAdminActions(api);
const saveRule = (markdown: string) =>
  run("ルール更新", () => actions.replaceRule(markdown));
const saveFreeze = (freeze: string) =>
  run("凍結時刻更新", () =>
    actions.replaceFreeze(freeze ? new Date(freeze).toISOString() : null),
  );
</script>
<template>
  <main class="workspace">
    <SettingsView
      :initial-markdown="data?.rule.markdown ?? ''"
      :initial-freeze="initialFreeze"
      :pending="pending"
      :error="error"
      :busy="busy"
      :message="message"
      @save-rule="saveRule"
      @save-freeze="saveFreeze"
      @retry="refresh"
    />
  </main>
</template>
