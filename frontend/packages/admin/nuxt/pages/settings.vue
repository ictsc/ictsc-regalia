<script setup lang="ts">
import SettingsView from "~/components/settings/SettingsView.vue";
import { api, expectData } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
useHead({ title: "競技設定" });
const { data, error, pending, refresh } = await useAsyncData(
  "settings",
  async () => {
    const [rule, schedule, phase] = await Promise.all([
      api.GET("/api/v1/admin/rule").then(expectData),
      api.GET("/api/v1/admin/dashboard-schedule").then(expectData),
      api.GET("/api/v1/admin/dev/contest-phase").then(expectData),
    ]);
    return {
      rule: rule.rule,
      schedule: schedule.dashboard_schedule,
      secondaryQualifier: schedule.competition_mode === "SECONDARY_QUALIFIER",
      contestPhase: phase.contest_phase,
    };
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
const phaseLabels = { BEFORE: "競技前", DURING: "競技中", AFTER: "競技後" };
const setPhase = (phase: keyof typeof phaseLabels) =>
  run(
    `競技状態を${phaseLabels[phase]}に切り替え`,
    () => actions.setDevContestPhase(phase),
    "有効なコンテンツのセクション時刻を、現在時刻を基準にずらして書き換えます。",
  );
</script>
<template>
  <main class="workspace">
    <SettingsView
      :rule-markdown="data?.rule.markdown ?? ''"
      :initial-freeze="initialFreeze"
      :pending="pending"
      :error="error"
      :busy="busy"
      :message="message"
      :contest-phase="data?.contestPhase"
      :secondary-qualifier="data?.secondaryQualifier"
      @set-phase="setPhase"
      @retry="refresh"
    />
  </main>
</template>
