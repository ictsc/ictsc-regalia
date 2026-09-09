<script setup lang="ts">
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
    <header class="linked-page-heading">
      <NuxtLink to="/problems">← 問題一覧へ</NuxtLink>
      <h1>順位表</h1>
      <p v-if="data?.frozen">順位表は凍結中です {{ data.frozenAt }}</p>
    </header>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><section class="ranking-table">
        <div class="ranking-row ranking-heading">
          <span>順位</span><span>チーム</span><span>得点</span>
        </div>
        <div
          v-for="r in data?.ranking"
          :key="r.teamCode"
          class="ranking-row"
          :class="{
            'is-team':
              viewer?.state === 'CONTESTANT' && r.teamCode === viewer.team.code,
          }"
        >
          <span class="ranking-position">{{ r.rank }}</span
          ><span class="ranking-team"
            ><strong>{{ r.teamName }}</strong
            ><small>{{ r.organization }}</small></span
          ><span class="ranking-score">{{ r.score.toLocaleString() }}</span>
        </div>
      </section>
      <p v-if="!data?.ranking.length" class="state-message">
        順位はまだありません。
      </p></RequestState
    >
  </main>
</template>
