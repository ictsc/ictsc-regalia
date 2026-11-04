<script setup lang="ts">
import { reactive, ref } from "vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import AdminFacts from "~/components/layout/AdminFacts.vue";
import AdminPanel from "~/components/layout/AdminPanel.vue";
import AdminPage from "~/components/layout/AdminPage.vue";
import { defaultTeamColor, teamColors, teamStyle } from "@ictsc/ui/colors";
import type { components } from "@ictsc/api";

type Schemas = components["schemas"];
export type TeamUpdateForm = {
  name: string;
  organization: string;
  member_limit: number;
  color: Schemas["Team"]["color"];
};

const props = defineProps<{
  code: number;
  team?: Schemas["Team"];
  members?: Schemas["AdminContestant"][];
  invitations?: Schemas["Invitation"][];
  pending?: boolean;
  error?: unknown;
  busy: boolean;
  /** Result of the last operation. */
  message: string;
}>();
defineEmits<{
  save: [form: TeamUpdateForm];
  /** Expiry as the `datetime-local` input value. */
  invite: [expires: string];
  remove: [];
  retry: [];
}>();

// Initialized once like the page did; edits are kept across refreshes.
const form = reactive<TeamUpdateForm>({
  name: props.team?.name ?? "",
  organization: props.team?.organization ?? "",
  member_limit: props.team?.member_limit ?? 4,
  color: props.team?.color ?? defaultTeamColor,
});
const expires = ref("");
</script>
<template>
  <AdminPage
    :title="team?.name ?? `チーム ${code}`"
    lead="チーム情報、メンバー、招待を管理します。"
    :back="{ to: '/teams', label: 'チーム一覧' }"
    :message="message"
  >
    <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
      ><AdminPanel title="基本情報">
        <template #aside
          ><span class="team-badge" :style="teamStyle(form.color)">{{
            form.name || "チーム名"
          }}</span></template
        >
        <AdminFacts :items="[{ label: 'チームコード', value: code }]" />
        <form class="form-grid" @submit.prevent="$emit('save', { ...form })">
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
      </AdminPanel>
      <AdminPanel title="メンバー">
        <template #aside>{{ members?.length ?? 0 }} 人</template>
        <div v-if="members?.length" class="table-scroll">
          <table class="data-table">
            <thead>
              <tr>
                <th>表示名</th>
                <th>ユーザー名</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="m in members" :key="m.profile.name">
                <td>{{ m.profile.display_name }}</td>
                <td>{{ m.profile.name }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-else class="admin-empty">メンバーはいません。</p>
      </AdminPanel>
      <AdminPanel title="招待">
        <template #aside>{{ invitations?.length ?? 0 }} 件</template>
        <form class="form-grid" @submit.prevent="$emit('invite', expires)">
          <label
            >有効期限<input
              v-model="expires"
              type="datetime-local"
              required /></label
          ><button class="button-primary" :disabled="busy">招待を作成</button>
        </form>
        <div v-if="invitations?.length" class="table-scroll">
          <table class="data-table">
            <thead>
              <tr>
                <th>招待コード</th>
                <th>有効期限</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="i in invitations" :key="i.code">
                <td class="admin-code">{{ i.code }}</td>
                <td>{{ new Date(i.expires_at).toLocaleString("ja-JP") }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-else class="admin-empty">招待はありません。</p>
      </AdminPanel>
      <AdminPanel title="チームの削除" danger>
        <p>チームを削除します。対象を確認してから実行してください。</p>
        <button
          class="button-secondary"
          :disabled="busy"
          @click="$emit('remove')"
        >
          チームを削除
        </button>
      </AdminPanel></RequestState
    >
  </AdminPage>
</template>
