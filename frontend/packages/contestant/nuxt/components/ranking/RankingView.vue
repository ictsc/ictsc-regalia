<script setup lang="ts">
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { RankingSnapshot } from "~/features/ranking";

defineProps<{
  snapshot?: RankingSnapshot;
  /** The viewer's team, highlighted in the table. */
  ownTeamCode?: number;
  pending?: boolean;
  error?: unknown;
}>();
defineEmits<{ retry: [] }>();
</script>
<template>
  <header class="linked-page-heading">
    <h1 class="page-title">順位表</h1>
    <p v-if="snapshot?.frozen">順位表は凍結中です {{ snapshot.frozenAt }}</p>
  </header>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><p v-if="snapshot?.hidden" class="state-message">
      採点は競技終了後に行います。順位は最終結果の公開後に表示されます。
    </p>
    <template v-else
      ><section class="ranking-table">
        <div class="ranking-row ranking-heading">
          <span>順位</span><span>チーム</span><span>得点</span>
        </div>
        <div
          v-for="r in snapshot?.ranking"
          :key="r.teamCode"
          class="ranking-row"
          :class="{
            'is-team': r.teamCode === ownTeamCode,
          }"
        >
          <span class="ranking-position">{{ r.rank }}</span
          ><span class="ranking-team"
            ><strong>{{ r.teamName }}</strong
            ><small>{{ r.organization }}</small
            ><span v-if="r.teamCode === ownTeamCode" class="ranking-own-marker"
              >自チーム</span
            ></span
          ><span class="ranking-score">{{ r.score.toLocaleString() }}</span>
        </div>
      </section>
      <p v-if="!snapshot?.ranking.length" class="state-message">
        順位はまだありません。
      </p></template
    ></RequestState
  >
</template>
