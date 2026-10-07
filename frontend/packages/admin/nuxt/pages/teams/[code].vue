<script setup lang="ts">
import { api, expectData, expectNoContent } from "@ictsc/api";
import { teamColors, defaultTeamColor, teamStyle } from "@ictsc/ui/colors";
definePageMeta({ key: (route) => route.fullPath });
const code = Number(useRoute().params.code);
const { data, error, pending, refresh } = await useAsyncData(
  `admin-team:${code}`,
  async () => {
    const [team, invitations, contestants] = await Promise.all([
      api
        .GET("/api/v1/admin/teams/{team_code}", {
          params: { path: { team_code: code } },
        })
        .then(expectData),
      api
        .GET("/api/v1/admin/invitations", {
          params: { query: { team_code: code, include_expired: true } },
        })
        .then(expectData),
      api.GET("/api/v1/admin/contestants").then(expectData),
    ]);
    return {
      team: team.team,
      invitations: invitations.invitations,
      members: contestants.contestants.filter((c) => c.team.code === code),
    };
  },
);
useHead({ title: computed(() => data.value?.team.name ?? "チーム詳細") });
const form = reactive({
  name: data.value?.team.name ?? "",
  organization: data.value?.team.organization ?? "",
  member_limit: data.value?.team.member_limit ?? 4,
  color: data.value?.team.color ?? defaultTeamColor,
});
const expires = ref("");
const { busy, message, run } = useMutation(refresh);
async function save() {
  await run("チームを更新", async () =>
    expectData(
      await api.PATCH("/api/v1/admin/teams/{team_code}", {
        params: { path: { team_code: code } },
        body: form,
      }),
    ),
  );
}
async function invite() {
  await run("招待を作成", async () =>
    expectData(
      await api.POST("/api/v1/admin/invitations", {
        body: {
          team_code: code,
          expires_at: new Date(expires.value).toISOString(),
        },
      }),
    ),
  );
}
async function remove() {
  if (
    await run(
      "チームを削除",
      async () =>
        expectNoContent(
          await api.DELETE("/api/v1/admin/teams/{team_code}", {
            params: { path: { team_code: code } },
          }),
        ),
      {
        message: `${data.value?.team.name} を削除しますか？この操作は取り消せません。`,
        destructive: true,
      },
    )
  )
    await navigateTo("/teams");
}
</script>
<template>
  <main class="workspace">
    <NuxtLink class="text-link" to="/teams">← チーム一覧</NuxtLink>
    <h1>チーム詳細</h1>
    <p class="admin-lead">チーム情報、メンバー、招待を管理します。</p>
    <p role="status">{{ message }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><section class="admin-panel">
        <div class="admin-meta">
          <span class="team-badge" :style="teamStyle(form.color)">{{
            form.name
          }}</span
          ><span class="admin-code">チームコード {{ code }}</span>
        </div>
        <h2>基本情報</h2>
        <form class="form-grid" @submit.prevent="save">
          <label
            >チーム名<input
              v-model="form.name"
              required
              maxlength="255" /></label
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
              <option v-for="c in teamColors" :key="c" :value="c">
                {{ c }}
              </option>
            </select></label
          ><button class="button-primary" :disabled="busy">
            チーム情報を保存
          </button>
        </form>
      </section>
      <section class="admin-panel">
        <div class="admin-panel-heading">
          <h2>メンバー</h2>
          <span>{{ data?.members.length ?? 0 }} 人</span>
        </div>
        <p v-for="m in data?.members" :key="m.profile.name">
          {{ m.profile.display_name }} / {{ m.profile.name }}
        </p>
        <p v-if="!data?.members.length" class="admin-empty">
          メンバーはいません。
        </p>
      </section>
      <section class="admin-panel">
        <h2>招待</h2>
        <form class="form-grid" @submit.prevent="invite">
          <label
            >有効期限<input
              v-model="expires"
              type="datetime-local"
              required /></label
          ><button class="button-primary" :disabled="busy">招待を作成</button>
        </form>
        <div class="table-scroll">
          <table class="data-table">
            <thead>
              <tr>
                <th>招待コード</th>
                <th>有効期限</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="i in data?.invitations" :key="i.code">
                <td>{{ i.code }}</td>
                <td>{{ new Date(i.expires_at).toLocaleString("ja-JP") }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-if="!data?.invitations.length" class="admin-empty">
          招待はありません。
        </p>
      </section>
      <section class="admin-panel admin-danger">
        <h2>チームの削除</h2>
        <p>チームを削除します。対象を確認してから実行してください。</p>
        <button class="button-secondary" :disabled="busy" @click="remove">
          チームを削除
        </button>
      </section></RequestState
    >
  </main>
</template>
