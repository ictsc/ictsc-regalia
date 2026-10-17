<script setup lang="ts">
import RankingView from "~/components/ranking/RankingView.vue";
import { api } from "@ictsc/api";
import { fetchRanking } from "~/features/ranking";
useHead({ title: "順位表" });
const { viewer } = useSession();
const { data, error, pending, refresh } = await useAsyncData("ranking", () =>
  fetchRanking(api),
);
</script>
<template>
  <main>
    <RankingView
      :snapshot="data ?? undefined"
      :own-team-code="
        viewer?.state === 'CONTESTANT' ? viewer.team.code : undefined
      "
      :pending="pending"
      :error="error"
      @retry="refresh"
    />
  </main>
</template>
