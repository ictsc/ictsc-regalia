<script setup lang="ts">
import ScheduleView from "~/components/schedule/ScheduleView.vue";
useHead({ title: "スケジュール" });
const { data, error, pending, refresh } = await useCompetition();
const now = useClock();
</script>
<template>
  <main class="workspace">
    <ScheduleView
      :entries="data?.schedule?.entries"
      :now="now"
      :pending="pending"
      :error="error"
      @retry="refresh"
    />
  </main>
</template>
