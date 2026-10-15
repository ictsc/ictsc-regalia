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
    <p class="admin-lead">
      ランキングの公開状態と問題別得点を確認します。得点の再計算と最終公開は必要な場合のみ実行してください。
    </p>
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
          run('最終得点公開', actions.reveal, {
            message:
              '凍結と遅延を解除して最終得点を公開します。この操作は取り消せません。続行しますか？',
            destructive: true,
          })
        "
      >
        最終得点を公開
      </button>
    </div>
    <p role="status">{{ message }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><section class="admin-panel">
        <div class="admin-panel-heading">
          <h2>ランキング</h2>
          <span
            class="admin-badge"
            :class="{ 'is-active': data?.ranking.frozen }"
            >{{ data?.ranking.frozen ? "凍結中" : "公開中" }}</span
          >
        </div>
        <p v-if="data?.ranking.frozen_at" class="muted">
          凍結時刻:
          {{ new Date(data.ranking.frozen_at).toLocaleString("ja-JP") }}
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
              <tr v-for="r in data?.ranking.ranking" :key="r.team_code">
                <td>{{ r.rank }}</td>
                <td>{{ r.team_name }}</td>
                <td>{{ r.score }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-if="!data?.ranking.ranking.length" class="admin-empty">
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
          </table>
        </div>
        <p v-if="!data?.scores.length" class="admin-empty">
          得点データはありません。
        </p>
      </section></RequestState
    >
  </main>
</template>
