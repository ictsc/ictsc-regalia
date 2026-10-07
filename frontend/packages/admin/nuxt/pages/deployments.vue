<script setup lang="ts">
import { api, expectData, subscribeAdminDeploymentEvents } from "@ictsc/api";
import {
  createAdminActions,
  mergeAdminDeployments,
} from "~/features/admin-api";
useHead({ title: "再展開" });
const { data, error, pending, refresh } = await useAsyncData(
  "admin-deployments",
  async () => {
    const [deployments, teams, problems] = await Promise.all([
      api.GET("/api/v1/admin/deployments").then(expectData),
      api.GET("/api/v1/admin/teams").then(expectData),
      api.GET("/api/v1/admin/problems").then(expectData),
    ]);
    return {
      deployments: deployments.deployments,
      teams: teams.teams,
      problems: problems.problems,
    };
  },
  { deep: true },
);
const team = ref(""),
  problem = ref(""),
  selectedKey = ref(""),
  streamError = ref(false);
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
  data.value?.deployments.find(
    (d) =>
      deploymentKey(d.team_code, d.problem_code, d.revision) ===
      selectedKey.value,
  ),
);
const filtered = computed(() =>
  data.value?.deployments.filter(
    (d) =>
      (!team.value || d.team_code === Number(team.value)) &&
      (!problem.value || d.problem_code === problem.value),
  ),
);
const { busy, message, run } = useMutation();
const actions = createAdminActions(api);
let stop: (() => void) | undefined;
onMounted(() => {
  stop = subscribeAdminDeploymentEvents(
    "/api/v1/admin/deployments/stream",
    (m) => {
      streamError.value = false;
      if (data.value)
        data.value.deployments = mergeAdminDeployments(
          data.value.deployments,
          m,
        );
    },
    () => {
      streamError.value = true;
    },
  );
});
onUnmounted(() => stop?.());
async function create() {
  await run("再展開要求", async () => {
    const r = await actions.createDeployment(Number(team.value), problem.value);
    if (data.value)
      data.value.deployments = mergeAdminDeployments(data.value.deployments, {
        type: "deployment",
        deployment:
          data.value.deployments.find(
            (d) =>
              d.team_code === r.deployment.team_code &&
              d.problem_code === r.deployment.problem_code &&
              d.revision === r.deployment.revision &&
              d.latest_status !== "QUEUED",
          ) ?? r.deployment,
      });
  });
}
</script>
<template>
  <main class="workspace">
    <h1>再展開</h1>
    <p class="admin-lead">
      チームと問題を指定して環境の再展開を要求します。履歴は状態更新に応じて反映されます。
    </p>
    <p role="status">{{ message }}</p>
    <p v-if="streamError" role="alert">状態更新に再接続しています。</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><section class="admin-panel">
        <h2>再展開を要求</h2>
        <p>対象のチームと問題を選択してください。</p>
        <form class="form-grid" @submit.prevent="create">
          <label
            >チーム<select v-model="team">
              <option value="">チームを選択</option>
              <option
                v-for="t in data?.teams"
                :key="t.code"
                :value="String(t.code)"
              >
                {{ t.code }}: {{ t.name }}
              </option>
            </select></label
          ><label
            >問題<select v-model="problem">
              <option value="">問題を選択</option>
              <option v-for="p in data?.problems" :key="p.code" :value="p.code">
                {{ p.code }}: {{ p.title }}
              </option>
            </select></label
          ><button class="button-primary" :disabled="busy || !team || !problem">
            選択した組み合わせを再展開
          </button>
        </form>
      </section>
      <p v-if="busy" role="status">再展開を要求中…</p>
      <section class="admin-panel">
        <div class="admin-panel-heading">
          <h2>再展開履歴</h2>
          <span>{{ filtered?.length ?? 0 }} 件</span>
        </div>
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
                    履歴
                  </button>
                  <button
                    class="button-secondary"
                    :disabled="busy"
                    @click="
                      run(
                        '手動同期',
                        () =>
                          actions.syncDeployment(d.team_code, d.problem_code),
                        '一度だけ状態を照会しますか？',
                      )
                    "
                  >
                    一回同期
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-if="!filtered?.length" class="admin-empty">
          該当する再展開履歴はありません。
        </p>
      </section>
      <section v-if="selected" class="admin-panel">
        <div class="admin-panel-heading">
          <h2>
            イベント履歴 · {{ selected.problem_code }} / チーム
            {{ selected.team_code }}
          </h2>
          <button class="button-secondary" @click="selectedKey = ''">
            閉じる
          </button>
        </div>
        <article
          v-for="e in selected.events"
          :key="e.event_id"
          class="content-section"
        >
          <h3>
            <span class="admin-badge">{{
              statusLabels[e.status] ?? e.status
            }}</span>
          </h3>
          <time>{{ new Date(e.occurred_at).toLocaleString("ja-JP") }}</time>
          <p>{{ e.message ?? "メッセージなし" }}</p>
        </article>
      </section></RequestState
    >
  </main>
</template>
