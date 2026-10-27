<script setup lang="ts">
import { ref } from "vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
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
  <h1>コンテンツ</h1>
  <p class="admin-lead">
    配信中の版を確認し、問題・お知らせを参照します。版の更新は Git commit
    を指定して実行します。
  </p>
  <p role="status">{{ message }}</p>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><template v-if="status"
      ><section class="admin-panel">
        <h2>配信状態</h2>
        <div class="admin-meta">
          <span
            class="admin-badge"
            :class="{ 'is-active': status.state === 'SUCCEEDED' }"
            >{{ stateLabels[status.state] ?? status.state }}</span
          ><span v-if="status.serving_last_known_good" class="admin-badge"
            >最終正常版を配信中</span
          >
        </div>
        <p class="admin-code">
          有効な commit: {{ status.active_commit ?? "なし" }}
        </p>
        <p v-if="status.last_error" role="alert">
          {{ status.last_error }}
        </p>
      </section>
      <section class="admin-panel">
        <h2>配信内容を更新</h2>
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
              required
              autocomplete="off"
              spellcheck="false" /></label
          ><button class="button-primary" :disabled="busy">
            指定commitを取得
          </button>
        </form>
      </section>
      <section class="admin-panel">
        <h2>競技セクション</h2>
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
                  <td>{{ new Date(s.beginning).toLocaleString("ja-JP") }}</td>
                  <td>{{ new Date(s.ending).toLocaleString("ja-JP") }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <p v-if="!sections.data?.sections.length" class="admin-empty">
            セクションはありません。
          </p></RequestState
        >
      </section>
      <section class="admin-panel">
        <h2>問題</h2>
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
      </section>
      <section class="admin-panel">
        <h2>お知らせ</h2>
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
      </section></template
    ></RequestState
  >
</template>
