<script setup lang="ts">
import { api, expectData } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
useHead({ title: "コンテンツ" });
async function settled<T>(p: Promise<T>) {
  try {
    return { data: await p, error: undefined };
  } catch (error) {
    return { data: undefined, error };
  }
}
const { data, error, pending, refresh } = await useAsyncData(
  "admin-content",
  async () => {
    const [status, sections, problems, announcements] = await Promise.all([
      api.GET("/api/v1/admin/content/status").then(expectData),
      settled(api.GET("/api/v1/admin/sections").then(expectData)),
      settled(api.GET("/api/v1/admin/problems").then(expectData)),
      settled(api.GET("/api/v1/admin/announcements").then(expectData)),
    ]);
    return { status: status.content, sections, problems, announcements };
  },
);
const commit = ref(data.value?.status.active_commit ?? "");
const { busy, message, run } = useMutation(refresh);
const stateLabels: Record<string, string> = {
  NEVER: "未配信",
  REFRESHING: "更新中",
  SUCCEEDED: "配信中",
  FAILED: "更新失敗",
};
</script>
<template>
  <main class="workspace">
    <h1>コンテンツ</h1>
    <p class="admin-lead">
      配信中の版を確認し、問題・お知らせを参照します。版の更新は Git commit
      を指定して実行します。
    </p>
    <p role="status">{{ message }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><template v-if="data"
        ><section class="admin-panel">
          <h2>配信状態</h2>
          <div class="admin-meta">
            <span
              class="admin-badge"
              :class="{ 'is-active': data.status.state === 'SUCCEEDED' }"
              >{{ stateLabels[data.status.state] ?? data.status.state }}</span
            ><span
              v-if="data.status.serving_last_known_good"
              class="admin-badge"
              >最終正常版を配信中</span
            >
          </div>
          <p class="admin-code">
            有効な commit: {{ data.status.active_commit ?? "なし" }}
          </p>
          <p v-if="data.status.last_error" role="alert">
            {{ data.status.last_error }}
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
            @submit.prevent="
              run('コンテンツ更新', () =>
                createAdminActions(api).refreshContent(commit.trim()),
              )
            "
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
          <RequestState :error="data.sections.error" @retry="refresh"
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
                  <tr v-for="s in data.sections.data?.sections" :key="s.slug">
                    <td>{{ s.slug }}</td>
                    <td>{{ s.beginning }}</td>
                    <td>{{ s.ending }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p v-if="!data.sections.data?.sections.length" class="admin-empty">
              セクションはありません。
            </p></RequestState
          >
        </section>
        <section class="admin-panel">
          <h2>問題</h2>
          <RequestState :error="data.problems.error" @retry="refresh"
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
                  <tr v-for="p in data.problems.data?.problems" :key="p.code">
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
            <p v-if="!data.problems.data?.problems.length" class="admin-empty">
              問題はありません。
            </p></RequestState
          >
        </section>
        <section class="admin-panel">
          <h2>お知らせ</h2>
          <RequestState :error="data.announcements.error" @retry="refresh"
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
                    v-for="a in data.announcements.data?.announcements"
                    :key="a.slug"
                  >
                    <td>
                      <NuxtLink :to="`/content/announcements/${a.slug}`">{{
                        a.title
                      }}</NuxtLink>
                    </td>
                    <td>{{ a.effective_from }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p
              v-if="!data.announcements.data?.announcements.length"
              class="admin-empty"
            >
              お知らせはありません。
            </p></RequestState
          >
        </section></template
      ></RequestState
    >
  </main>
</template>
