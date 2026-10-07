<script setup lang="ts">
import { api, expectData } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
useHead({ title: "参加者" });
const { data, error, pending, refresh } = await useAsyncData(
  "admin-contestants",
  async () => expectData(await api.GET("/api/v1/admin/contestants")),
);
const query = ref("");
const filtered = computed(() =>
  data.value?.contestants.filter((c) =>
    [
      c.profile.name,
      c.profile.display_name,
      c.team.name,
      c.team.organization,
      c.discord_id,
    ].some((s) => s.toLowerCase().includes(query.value.toLowerCase())),
  ),
);
const { busy, message, run } = useMutation();
async function impersonate(name: string) {
  if (
    await run(
      "代理ログイン",
      () => createAdminActions(api).impersonate(name),
      `${name} として代理ログインしますか？`,
    )
  )
    window.location.assign("/");
}
</script>
<template>
  <main class="workspace">
    <h1>参加者</h1>
    <p class="admin-lead">
      参加者の所属を確認します。代理ログインは対象者の表示を確認してから実行してください。
    </p>
    <section class="admin-panel">
      <label
        >検索
        <input
          v-model="query"
          class="field"
          placeholder="名前・チーム・Discord ID"
      /></label>
      <p v-if="data" class="muted">
        {{ filtered?.length ?? 0 }} / {{ data.contestants.length }} 人を表示
      </p>
      <p role="status">{{ message }}</p>
      <RequestState :error="error" :pending="pending" @retry="refresh"
        ><div class="table-scroll">
          <table class="data-table">
            <thead>
              <tr>
                <th>参加者</th>
                <th>チーム</th>
                <th>Discord ID</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="c in filtered" :key="c.profile.name">
                <td>
                  <strong>{{ c.profile.display_name }}</strong
                  ><br /><small>{{ c.profile.name }}</small>
                </td>
                <td>{{ c.team.name }}</td>
                <td>{{ c.discord_id }}</td>
                <td>
                  <button
                    class="button-secondary"
                    :disabled="busy"
                    @click="impersonate(c.profile.name)"
                  >
                    代理ログイン
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-if="!filtered?.length" class="admin-empty">
          該当する参加者はいません。
        </p></RequestState
      >
    </section>
  </main>
</template>
