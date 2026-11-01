<script setup lang="ts">
import RequestState from "@ictsc/ui/components/RequestState.vue";
import AdminPage from "~/components/layout/AdminPage.vue";
import AdminPanel from "~/components/layout/AdminPanel.vue";
import AdminFacts from "~/components/layout/AdminFacts.vue";
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
  <AdminPage
    title="得点・順位"
    lead="ランキングの公開状態と問題別得点を確認します。得点の再計算と最終公開は必要な場合のみ実行してください。"
    :message="message"
  >
    <template #actions
      ><button
        class="button-secondary"
        :disabled="busy"
        @click="$emit('recalculate')"
      >
        得点を再計算</button
      ><button class="button-primary" :disabled="busy" @click="$emit('reveal')">
        最終得点を公開
      </button></template
    >
    <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
      ><AdminPanel title="ランキング">
        <template #aside
          ><span
            class="admin-badge"
            :class="{ 'is-active': ranking?.frozen }"
            >{{ ranking?.frozen ? "凍結中" : "公開中" }}</span
          ></template
        >
        <AdminFacts
          v-if="ranking?.frozen_at"
          :items="[
            {
              label: '凍結時刻',
              value: new Date(ranking.frozen_at).toLocaleString('ja-JP'),
            },
          ]"
        />
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
      </AdminPanel>
      <AdminPanel title="問題別得点">
        <template #aside>{{ scores?.length ?? 0 }} 件</template>
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
        <p v-if="!scores?.length" class="admin-empty">
          得点データはありません。
        </p>
      </AdminPanel></RequestState
    >
  </AdminPage>
</template>
