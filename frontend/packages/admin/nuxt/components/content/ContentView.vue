<script setup lang="ts">
import { ref } from "vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import AdminPage from "~/components/layout/AdminPage.vue";
import AdminPanel from "~/components/layout/AdminPanel.vue";
import AdminFacts from "~/components/layout/AdminFacts.vue";
import type { components } from "@ictsc/api";

type Schemas = components["schemas"];
/** A list that loaded independently and may have failed on its own. */
type Settled<T> = { data?: T; error?: unknown };

const props = defineProps<{
  status?: Schemas["ContentStatus"];
  sections: Settled<{ sections: Schemas["Section"][] }>;
  problems: Settled<{ problems: Schemas["ProblemCatalogEntry"][] }>;
  announcements: Settled<{ announcements: Schemas["Announcement"][] }>;
  pending?: boolean;
  error?: unknown;
  busy: boolean;
  /** Result of the last operation. */
  message: string;
}>();
defineEmits<{ refreshContent: [commit: string]; retry: [] }>();

// Initialized once like the page did; edits are kept across refreshes.
const commit = ref(props.status?.active_commit ?? "");
const stateLabels: Record<string, string> = {
  NEVER: "未配信",
  REFRESHING: "更新中",
  SUCCEEDED: "配信中",
  FAILED: "更新失敗",
};
</script>
<template>
  <AdminPage
    title="コンテンツ"
    lead="配信中の版を確認し、問題・お知らせを参照します。版の更新は Git commit を指定して実行します。"
    :message="message"
  >
    <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
      ><template v-if="status"
        ><AdminPanel title="配信状態">
          <template #aside
            ><span
              class="admin-badge"
              :class="{ 'is-active': status.state === 'SUCCEEDED' }"
              >{{ stateLabels[status.state] ?? status.state }}</span
            ></template
          >
          <AdminFacts
            :items="[
              {
                label: '有効な commit',
                value: status.active_commit ?? 'なし',
                mono: true,
              },
            ]"
          />
          <p v-if="status.serving_last_known_good">
            <span class="admin-badge">最終正常版を配信中</span>
          </p>
          <p v-if="status.last_error" role="alert">
            {{ status.last_error }}
          </p>
        </AdminPanel>
        <AdminPanel title="配信内容を更新">
          <p>
            指定した commit
            のコンテンツを取得します。更新後は配信状態を確認してください。
          </p>
          <form
            class="form-grid"
            @submit.prevent="$emit('refreshContent', commit.trim())"
          >
            <label
              >Git commit SHA<input
                v-model="commit"
                class="admin-code"
                required
                autocomplete="off"
                spellcheck="false" /></label
            ><button class="button-primary" :disabled="busy">
              指定 commit を取得
            </button>
          </form>
        </AdminPanel>
        <AdminPanel title="競技セクション">
          <template #aside
            >{{ sections.data?.sections.length ?? 0 }} 件</template
          >
          <RequestState :error="sections.error" @retry="$emit('retry')"
            ><div class="table-scroll">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>識別子</th>
                    <th>開始</th>
                    <th>終了</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="s in sections.data?.sections" :key="s.slug">
                    <td>{{ s.slug }}</td>
                    <td>
                      {{ new Date(s.beginning).toLocaleString("ja-JP") }}
                    </td>
                    <td>{{ new Date(s.ending).toLocaleString("ja-JP") }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p v-if="!sections.data?.sections.length" class="admin-empty">
              セクションはありません。
            </p></RequestState
          >
        </AdminPanel>
        <AdminPanel title="問題">
          <template #aside
            >{{ problems.data?.problems.length ?? 0 }} 問</template
          >
          <RequestState :error="problems.error" @retry="$emit('retry')"
            ><div class="table-scroll">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>番号</th>
                    <th>タイトル</th>
                    <th>カテゴリ</th>
                    <th>満点</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="p in problems.data?.problems" :key="p.code">
                    <td>{{ p.code }}</td>
                    <td>
                      <NuxtLink :to="`/content/problems/${p.code}`">{{
                        p.title
                      }}</NuxtLink>
                    </td>
                    <td>{{ p.category }}</td>
                    <td>{{ p.max_score }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p v-if="!problems.data?.problems.length" class="admin-empty">
              問題はありません。
            </p></RequestState
          >
        </AdminPanel>
        <AdminPanel title="お知らせ">
          <template #aside
            >{{ announcements.data?.announcements.length ?? 0 }} 件</template
          >
          <RequestState :error="announcements.error" @retry="$emit('retry')"
            ><div class="table-scroll">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>タイトル</th>
                    <th>公開日時</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="a in announcements.data?.announcements"
                    :key="a.slug"
                  >
                    <td>
                      <NuxtLink :to="`/content/announcements/${a.slug}`">{{
                        a.title
                      }}</NuxtLink>
                    </td>
                    <td>
                      {{ new Date(a.effective_from).toLocaleString("ja-JP") }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p
              v-if="!announcements.data?.announcements.length"
              class="admin-empty"
            >
              お知らせはありません。
            </p></RequestState
          >
        </AdminPanel></template
      ></RequestState
    >
  </AdminPage>
</template>
