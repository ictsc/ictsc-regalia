<script setup lang="ts">
import { computed, ref } from "vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import AdminPage from "~/components/layout/AdminPage.vue";
import AdminPanel from "~/components/layout/AdminPanel.vue";
import type { components } from "@ictsc/api";

type Schemas = components["schemas"];

const props = defineProps<{
  deployments?: Schemas["AdminDeployment"][];
  teams?: Schemas["Team"][];
  problems?: Schemas["ProblemCatalogEntry"][];
  /** The deployment event stream is reconnecting. */
  streamError: boolean;
  pending?: boolean;
  error?: unknown;
  busy: boolean;
  /** Result of the last operation. */
  message: string;
}>();
defineEmits<{
  create: [teamCode: number, problemCode: string];
  sync: [teamCode: number, problemCode: string];
  retry: [];
}>();

const team = ref(""),
  problem = ref(""),
  selectedKey = ref("");
const statusLabels: Record<string, string> = {
  QUEUED: "待機中",
  DEPLOYING: "展開中",
  COMPLETED: "完了",
  FAILED: "失敗",
};
const deploymentKey = (
  teamCode: number,
  problemCode: string,
  revision: number,
) => `${teamCode}:${problemCode}:${revision}`;
const selected = computed(() =>
  props.deployments?.find(
    (d) =>
      deploymentKey(d.team_code, d.problem_code, d.revision) ===
      selectedKey.value,
  ),
);
const filtered = computed(() =>
  props.deployments?.filter(
    (d) =>
      (!team.value || d.team_code === Number(team.value)) &&
      (!problem.value || d.problem_code === problem.value),
  ),
);
</script>
<template>
  <AdminPage
    title="再展開"
    lead="チームと問題を指定して環境の再展開を要求します。履歴は状態更新に応じて反映されます。"
    :message="busy ? '再展開を要求中…' : message"
  >
    <p v-if="streamError" role="alert" class="admin-status">
      状態更新に再接続しています。
    </p>
    <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
      ><AdminPanel title="再展開を要求">
        <p>対象のチームと問題を選択してください。</p>
        <form
          class="form-grid"
          @submit.prevent="$emit('create', Number(team), problem)"
        >
          <label
            >チーム<select v-model="team">
              <option value="">チームを選択</option>
              <option v-for="t in teams" :key="t.code" :value="String(t.code)">
                {{ t.code }}: {{ t.name }}
              </option>
            </select></label
          ><label
            >問題<select v-model="problem">
              <option value="">問題を選択</option>
              <option v-for="p in problems" :key="p.code" :value="p.code">
                {{ p.code }}: {{ p.title }}
              </option>
            </select></label
          ><button class="button-primary" :disabled="busy || !team || !problem">
            選択した組み合わせを再展開
          </button>
        </form>
      </AdminPanel>
      <AdminPanel title="再展開履歴">
        <template #aside>{{ filtered?.length ?? 0 }} 件</template>
        <div class="table-scroll">
          <table class="data-table">
            <thead>
              <tr>
                <th>チーム</th>
                <th>問題</th>
                <th>Revision</th>
                <th>状態</th>
                <th>Commit</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="d in filtered"
                :key="`${d.team_code}:${d.problem_code}:${d.revision}`"
              >
                <td>{{ d.team_code }}</td>
                <td>{{ d.problem_code }}</td>
                <td>{{ d.revision }}</td>
                <td>
                  <span
                    class="admin-badge"
                    :class="{
                      'is-active': ['QUEUED', 'DEPLOYING'].includes(
                        d.latest_status,
                      ),
                    }"
                    >{{
                      statusLabels[d.latest_status] ?? d.latest_status
                    }}</span
                  >
                </td>
                <td class="admin-code">{{ d.content_commit.slice(0, 12) }}</td>
                <td>
                  <div class="admin-row-actions">
                    <button
                      class="button-secondary"
                      @click="
                        selectedKey = deploymentKey(
                          d.team_code,
                          d.problem_code,
                          d.revision,
                        )
                      "
                    >
                      履歴</button
                    ><button
                      class="button-secondary"
                      :disabled="busy"
                      @click="$emit('sync', d.team_code, d.problem_code)"
                    >
                      一回同期
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-if="!filtered?.length" class="admin-empty">
          該当する再展開履歴はありません。
        </p>
      </AdminPanel>
      <AdminPanel
        v-if="selected"
        :title="`イベント履歴 · ${selected.problem_code} / チーム ${selected.team_code}`"
      >
        <template #aside
          ><button class="button-secondary" @click="selectedKey = ''">
            閉じる
          </button></template
        >
        <div class="table-scroll">
          <table class="data-table">
            <thead>
              <tr>
                <th>状態</th>
                <th>日時</th>
                <th>メッセージ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="e in selected.events" :key="e.event_id">
                <td>
                  <span class="admin-badge">{{
                    statusLabels[e.status] ?? e.status
                  }}</span>
                </td>
                <td>
                  <time :datetime="e.occurred_at">{{
                    new Date(e.occurred_at).toLocaleString("ja-JP")
                  }}</time>
                </td>
                <td>{{ e.message ?? "メッセージなし" }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </AdminPanel></RequestState
    >
  </AdminPage>
</template>
