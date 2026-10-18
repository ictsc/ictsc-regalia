<script setup lang="ts">
import { reactive } from "vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import { defaultTeamColor, teamColors } from "@ictsc/ui/colors";
import type { components } from "@ictsc/api";

export type TeamCreateForm = components["schemas"]["CreateTeamRequest"];

defineProps<{
  teams?: components["schemas"]["Team"][];
  pending?: boolean;
  error?: unknown;
  busy: boolean;
  /** Result of the last operation. */
  message: string;
}>();
defineEmits<{ create: [form: TeamCreateForm]; retry: [] }>();

const form = reactive({
  code: 2,
  name: "",
  organization: "",
  member_limit: 4,
  color: defaultTeamColor as (typeof teamColors)[number],
});
</script>
<template>
  <h1>チーム</h1>
  <p class="admin-lead">
    競技チームの情報を確認・編集し、新しいチームを登録します。
  </p>
  <p role="status">{{ message }}</p>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><section class="admin-panel">
      <div class="admin-panel-heading">
        <h2>登録済みチーム</h2>
        <span>{{ teams?.length ?? 0 }} チーム</span>
      </div>
      <div class="table-scroll">
        <table class="data-table">
          <thead>
            <tr>
              <th>コード</th>
              <th>チーム</th>
              <th>所属</th>
              <th>定員</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in teams" :key="t.code">
              <td>{{ t.code }}</td>
              <td>
                <NuxtLink :to="`/teams/${t.code}`">{{ t.name }}</NuxtLink>
              </td>
              <td>{{ t.organization }}</td>
              <td>{{ t.member_limit }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p v-if="!teams?.length" class="admin-empty">
        登録済みのチームはありません。
      </p>
    </section></RequestState
  >
  <section class="admin-panel">
    <h2>チーム作成</h2>
    <form class="form-grid" @submit.prevent="$emit('create', { ...form })">
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
  </section>
</template>
