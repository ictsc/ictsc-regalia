<script setup lang="ts">
import { reactive, ref } from "vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
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
  <NuxtLink class="text-link" to="/teams">← チーム一覧</NuxtLink>
  <h1>チーム詳細</h1>
  <p class="admin-lead">チーム情報、メンバー、招待を管理します。</p>
  <p role="status">{{ message }}</p>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><section class="admin-panel">
      <div class="admin-meta">
        <span class="team-badge" :style="teamStyle(form.color)">{{
          form.name
        }}</span
        ><span class="admin-code">チームコード {{ code }}</span>
      </div>
      <h2>基本情報</h2>
      <form class="form-grid" @submit.prevent="$emit('save', { ...form })">
        <label
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
        <span>{{ members?.length ?? 0 }} 人</span>
      </div>
      <p v-for="m in members" :key="m.profile.name">
        {{ m.profile.display_name }} / {{ m.profile.name }}
      </p>
      <p v-if="!members?.length" class="admin-empty">メンバーはいません。</p>
    </section>
    <section class="admin-panel">
      <h2>招待</h2>
      <form class="form-grid" @submit.prevent="$emit('invite', expires)">
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
            <tr v-for="i in invitations" :key="i.code">
              <td>{{ i.code }}</td>
              <td>{{ new Date(i.expires_at).toLocaleString("ja-JP") }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p v-if="!invitations?.length" class="admin-empty">招待はありません。</p>
    </section>
    <section class="admin-panel admin-danger">
      <h2>チームの削除</h2>
      <p>チームを削除します。対象を確認してから実行してください。</p>
      <button
        class="button-secondary"
        :disabled="busy"
        @click="$emit('remove')"
      >
        チームを削除
      </button>
    </section></RequestState
  >
</template>
