<script setup lang="ts">
import { api, expectData } from "@ictsc/api";
import { teamColors, defaultTeamColor } from "@ictsc/ui/colors";
useHead({ title: "チーム" });
const { data, error, pending, refresh } = await useAsyncData(
  "admin-teams",
  async () => expectData(await api.GET("/api/v1/admin/teams")),
);
const { busy, message, run } = useMutation(refresh);
const form = reactive({
  code: 2,
  name: "",
  organization: "",
  member_limit: 4,
  color: defaultTeamColor as (typeof teamColors)[number],
});
async function create() {
  await run("チームを作成", async () =>
    expectData(await api.POST("/api/v1/admin/teams", { body: form })),
  );
}
</script>
<template>
  <main class="workspace">
    <h1>チーム</h1>
    <p role="status">{{ message }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><table class="data-table">
        <thead>
          <tr>
            <th>コード</th>
            <th>チーム</th>
            <th>所属</th>
            <th>定員</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in data?.teams" :key="t.code">
            <td>{{ t.code }}</td>
            <td>
              <NuxtLink :to="`/teams/${t.code}`">{{ t.name }}</NuxtLink>
            </td>
            <td>{{ t.organization }}</td>
            <td>{{ t.member_limit }}</td>
          </tr>
        </tbody>
      </table></RequestState
    >
    <h2>チーム作成</h2>
    <form class="form-grid" @submit.prevent="create">
      <label
        >チームコード<input
          v-model.number="form.code"
          type="number"
          min="2"
          max="99"
          required /></label
      ><label
        >チーム名<input v-model="form.name" required maxlength="255" /></label
      ><label
        >所属<input
          v-model="form.organization"
          required
          maxlength="255" /></label
      ><label
        >定員<input
          v-model.number="form.member_limit"
          type="number"
          min="1"
          required /></label
      ><label
        >チームカラー<select v-model="form.color">
          <option v-for="c in teamColors" :key="c" :value="c">{{ c }}</option>
        </select></label
      ><button class="button-primary" :disabled="busy">チームを作成</button>
    </form>
  </main>
</template>
