<script setup lang="ts">
import TeamDetailView, {
  type TeamUpdateForm,
} from "~/components/team/TeamDetailView.vue";
import { api, expectData, expectNoContent } from "@ictsc/api";
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
const { busy, message, run } = useMutation(refresh);
async function save(form: TeamUpdateForm) {
  await run("チームを更新", async () =>
    expectData(
      await api.PATCH("/api/v1/admin/teams/{team_code}", {
        params: { path: { team_code: code } },
        body: form,
      }),
    ),
  );
}
async function invite(expires: string) {
  await run("招待を作成", async () =>
    expectData(
      await api.POST("/api/v1/admin/invitations", {
        body: {
          team_code: code,
          expires_at: new Date(expires).toISOString(),
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
    <TeamDetailView
      :code="code"
      :team="data?.team"
      :members="data?.members"
      :invitations="data?.invitations"
      :pending="pending"
      :error="error"
      :busy="busy"
      :message="message"
      @save="save"
      @invite="invite"
      @remove="remove"
      @retry="refresh"
    />
  </main>
</template>
