<script setup lang="ts">
import TeamListView, {
  type TeamCreateForm,
} from "~/components/team/TeamListView.vue";
import { api, expectData } from "@ictsc/api";
useHead({ title: "チーム" });
const { data, error, pending, refresh } = await useAsyncData(
  "admin-teams",
  async () => expectData(await api.GET("/api/v1/admin/teams")),
);
const { busy, message, run } = useMutation(refresh);
async function create(form: TeamCreateForm) {
  await run("チームを作成", async () =>
    expectData(await api.POST("/api/v1/admin/teams", { body: form })),
  );
}
</script>
<template>
  <main class="workspace">
    <TeamListView
      :teams="data?.teams"
      :pending="pending"
      :error="error"
      :busy="busy"
      :message="message"
      @create="create"
      @retry="refresh"
    />
  </main>
</template>
