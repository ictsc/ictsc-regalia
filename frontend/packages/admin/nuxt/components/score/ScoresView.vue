<script setup lang="ts">
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { components } from "@ictsc/api";

type Schemas = components["schemas"];

defineProps<{
  ranking?: Schemas["RankingResponse"];
  scores?: Schemas["AdminScore"][];
  pending?: boolean;
  error?: unknown;
  busy: boolean;
  /** Result of the last operation. */
  message: string;
}>();
defineEmits<{ recalculate: []; reveal: []; retry: [] }>();
</script>
<template>
  <h1>得点・順位</h1>
  <p class="admin-lead">
    ランキングの公開状態と問題別得点を確認します。得点の再計算と最終公開は必要な場合のみ実行してください。
  </p>
  <div class="actions">
    <button
      class="button-secondary"
      :disabled="busy"
      @click="$emit('recalculate')"
    >
      得点を再計算</button
    ><button class="button-primary" :disabled="busy" @click="$emit('reveal')">
      最終得点を公開
    </button>
  </div>
  <p role="status">{{ message }}</p>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><section class="admin-panel">
      <div class="admin-panel-heading">
        <h2>ランキング</h2>
        <span class="admin-badge" :class="{ 'is-active': ranking?.frozen }">{{
          ranking?.frozen ? "凍結中" : "公開中"
        }}</span>
      </div>
      <p v-if="ranking?.frozen_at" class="muted">
        凍結時刻:
        {{ new Date(ranking.frozen_at).toLocaleString("ja-JP") }}
      </p>
      <div class="table-scroll">
        <table class="data-table">
          <thead>
            <tr>
              <th>順位</th>
              <th>チーム</th>
              <th>得点</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in ranking?.ranking" :key="r.team_code">
              <td>{{ r.rank }}</td>
              <td>{{ r.team_name }}</td>
              <td>{{ r.score }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p v-if="!ranking?.ranking.length" class="admin-empty">
        ランキングデータはありません。
      </p>
    </section>
    <section class="admin-panel">
      <h2>問題別得点</h2>
      <div class="table-scroll">
        <table class="data-table">
          <thead>
            <tr>
              <th>チーム</th>
              <th>問題</th>
              <th>得点</th>
              <th>素点</th>
              <th>減点</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in scores" :key="`${s.team.code}:${s.problem.code}`">
              <td>{{ s.team.name }}</td>
              <td>{{ s.problem.code }} {{ s.problem.title }}</td>
              <td>{{ s.score }}</td>
              <td>{{ s.marked_score }}</td>
              <td>{{ s.penalty }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p v-if="!scores?.length" class="admin-empty">得点データはありません。</p>
    </section></RequestState
  >
</template>
