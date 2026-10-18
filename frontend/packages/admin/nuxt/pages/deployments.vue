<script setup lang="ts">
import DeploymentsView from "~/components/deployment/DeploymentsView.vue";
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
const streamError = ref(false);
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
async function create(teamCode: number, problemCode: string) {
  await run("再展開要求", async () => {
    const r = await actions.createDeployment(teamCode, problemCode);
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
const sync = (teamCode: number, problemCode: string) =>
  run(
    "手動同期",
    () => actions.syncDeployment(teamCode, problemCode),
    "一度だけ状態を照会しますか？",
  );
</script>
<template>
  <main class="workspace">
    <DeploymentsView
      :deployments="data?.deployments"
      :teams="data?.teams"
      :problems="data?.problems"
      :stream-error="streamError"
      :pending="pending"
      :error="error"
      :busy="busy"
      :message="message"
      @create="create"
      @sync="sync"
      @retry="refresh"
    />
  </main>
</template>
