<script setup lang="ts">
import { api, expectData } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
useHead({ title: "得点・順位" });
const { data, error, pending, refresh } = await useAsyncData(
  "admin-scores",
  async () => {
    const [scores, ranking] = await Promise.all([
      api.GET("/api/v1/admin/scores").then(expectData),
      api.GET("/api/v1/admin/ranking").then(expectData),
    ]);
    return { scores: scores.scores, ranking };
  },
);
const { busy, message, run } = useMutation(refresh);
const actions = createAdminActions(api);
</script>
<template>
  <main class="workspace">
    <h1>得点・順位</h1>
    <div class="actions">
      <button
        class="button-secondary"
        :disabled="busy"
        @click="
          run(
            '得点再計算',
            actions.recalculate,
            'すべての得点を再計算しますか？',
          )
        "
      >
        得点を再計算</button
      ><button
        class="button-primary"
        :disabled="busy"
        @click="
          run(
            '最終得点公開',
            actions.reveal,
            '凍結と遅延を解除して最終得点を公開します。この操作は取り消せません。続行しますか？',
          )
        "
      >
        最終得点を公開
      </button>
    </div>
    <p role="status">{{ message }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><h2>ランキング</h2>
      <p>
        {{ data?.ranking.frozen ? "凍結中" : "公開中" }}
        {{ data?.ranking.frozen_at }}
      </p>
      <table class="data-table">
        <thead>
          <tr>
            <th>順位</th>
            <th>チーム</th>
            <th>得点</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in data?.ranking.ranking" :key="r.team_code">
            <td>{{ r.rank }}</td>
            <td>{{ r.team_name }}</td>
            <td>{{ r.score }}</td>
          </tr>
        </tbody>
      </table>
      <h2>問題別得点</h2>
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
          <tr
            v-for="s in data?.scores"
            :key="`${s.team.code}:${s.problem.code}`"
          >
            <td>{{ s.team.name }}</td>
            <td>{{ s.problem.code }} {{ s.problem.title }}</td>
            <td>{{ s.score }}</td>
            <td>{{ s.marked_score }}</td>
            <td>{{ s.penalty }}</td>
          </tr>
        </tbody>
      </table></RequestState
    >
  </main>
</template>
