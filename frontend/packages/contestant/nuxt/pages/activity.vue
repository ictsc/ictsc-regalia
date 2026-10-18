<script setup lang="ts">
import ActivityView from "~/components/activity/ActivityView.vue";
import { api } from "@ictsc/api";
import { fetchActivity } from "~/features/activity";
useHead({ title: "提出履歴" });
const { data, error, pending, refresh } = await useAsyncData("activity", () =>
  fetchActivity(api),
);
</script>
<template>
  <main>
    <ActivityView
      :entries="data ?? undefined"
      :pending="pending"
      :error="error"
      @retry="refresh"
    />
  </main>
</template>
