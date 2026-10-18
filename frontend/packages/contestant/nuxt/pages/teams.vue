<script setup lang="ts">
import TeamDirectoryView from "~/components/team/TeamDirectoryView.vue";
import { api } from "@ictsc/api";
import { fetchTeams } from "~/features/teams";
useHead({ title: "チーム" });
const { data, error, pending, refresh } = await useAsyncData("teams", () =>
  fetchTeams(api),
);
</script>
<template>
  <main class="workspace">
    <TeamDirectoryView
      :teams="data ?? undefined"
      :pending="pending"
      :error="error"
      @retry="refresh"
    />
  </main>
</template>
