<script setup lang="ts">
import { computed, ref } from "vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import AdminPage from "~/components/layout/AdminPage.vue";
import AdminPanel from "~/components/layout/AdminPanel.vue";
import type { components } from "@ictsc/api";

const props = defineProps<{
  contestants?: components["schemas"]["AdminContestant"][];
  pending?: boolean;
  error?: unknown;
  busy: boolean;
  /** Result of the last operation. */
  message: string;
}>();
defineEmits<{ impersonate: [name: string]; retry: [] }>();

const query = ref("");
const filtered = computed(() =>
  props.contestants?.filter((c) =>
    [
      c.profile.name,
      c.profile.display_name,
      c.team.name,
      c.team.organization,
      c.discord_id,
    ].some((s) => s.toLowerCase().includes(query.value.toLowerCase())),
  ),
);
</script>
<template>
  <AdminPage
    title="参加者"
    lead="参加者の所属を確認します。代理ログインは対象者の表示を確認してから実行してください。"
    :message="message"
  >
    <AdminPanel title="参加者一覧">
      <template v-if="contestants" #aside
        >{{ filtered?.length ?? 0 }} / {{ contestants.length }} 人</template
      >
      <div class="form-grid">
        <label
          >検索<input
            v-model="query"
            type="search"
            placeholder="名前・チーム・Discord ID"
        /></label>
      </div>
      <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
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
                <td class="admin-code">{{ c.discord_id }}</td>
                <td>
                  <button
                    class="button-secondary"
                    :disabled="busy"
                    @click="$emit('impersonate', c.profile.name)"
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
    </AdminPanel>
  </AdminPage>
</template>
