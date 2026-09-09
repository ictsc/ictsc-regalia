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
</script>
<template>
  <main class="workspace">
    <h1>コンテンツ</h1>
    <p role="status">{{ message }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><template v-if="data"
        ><p>
          {{ data.status.state }} /
          {{ data.status.active_commit ?? "有効なコンテンツなし" }}
        </p>
        <p v-if="data.status.serving_last_known_good">最終正常版を配信中</p>
        <p v-if="data.status.last_error" role="alert">
          {{ data.status.last_error }}
        </p>
        <form
          class="form-grid"
          @submit.prevent="
            run('コンテンツ更新', () =>
              createAdminActions(api).refreshContent(commit.trim()),
            )
          "
        >
          <label>Git commit SHA<input v-model="commit" required /></label
          ><button class="button-primary" :disabled="busy">
            指定commitを取得
          </button>
        </form>
        <h2>Sections</h2>
        <RequestState :error="data.sections.error" @retry="refresh"
          ><table class="data-table">
            <tbody>
              <tr v-for="s in data.sections.data?.sections" :key="s.slug">
                <td>{{ s.slug }}</td>
                <td>{{ s.beginning }}</td>
                <td>{{ s.ending }}</td>
              </tr>
            </tbody>
          </table></RequestState
        >
        <h2>問題</h2>
        <RequestState :error="data.problems.error" @retry="refresh"
          ><table class="data-table">
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
          </table></RequestState
        >
        <h2>お知らせ</h2>
        <RequestState :error="data.announcements.error" @retry="refresh"
          ><table class="data-table">
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
          </table></RequestState
        ></template
      ></RequestState
    >
  </main>
</template>
