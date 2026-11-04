<script setup lang="ts">
import ActivityView from "~/components/activity/ActivityView.vue";
import { api } from "@ictsc/api";
import { fetchActivity } from "~/features/activity";
useHead({ title: "提出履歴" });
const { data, error, pending, refresh } = await useAsyncData("activity", () =>
  fetchActivity(api),
);
const { data: competition } = await useCompetition();
</script>
<template>
  <main>
    <ActivityView
      :entries="data ?? undefined"
      :results-hidden="competition?.ranking.hidden ?? false"
      :pending="pending"
      :error="error"
      @retry="refresh"
    />
  </main>
</template>
