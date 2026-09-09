<script setup lang="ts">
import { api, expectData } from "@ictsc/api";
useHead({ title: "採点" });
const { data, error, pending, refresh } = await useAsyncData(
  "admin-answers",
  async () =>
    expectData(
      await api.GET("/api/v1/admin/answers", {
        params: { query: { include_marked: true } },
      }),
    ),
);
const teams = ref<string[]>([]),
  problems = ref<string[]>([]),
  recent = ref(false),
  perfect = ref(true),
  unscored = ref(false);
onMounted(() => {
  try {
    const value = JSON.parse(
      localStorage.getItem("regalia/admin/answer-filters") ?? "null",
    );
    if (value) {
      teams.value = value.teams ?? [];
      problems.value = value.problems ?? [];
      recent.value = !!value.recent;
      perfect.value = value.perfect !== false;
      unscored.value = !!value.unscored;
    }
  } catch {
    /* Defaults remain usable when storage is unavailable. */
  }
});
watch(
  [teams, problems, recent, perfect, unscored],
  () => {
    try {
      localStorage.setItem(
        "regalia/admin/answer-filters",
        JSON.stringify({
          teams: teams.value,
          problems: problems.value,
          recent: recent.value,
          perfect: perfect.value,
          unscored: unscored.value,
        }),
      );
    } catch {
      /* Filtering does not depend on persistence. */
    }
  },
  { deep: true },
);
const filtered = computed(() =>
  data.value?.answers
    .filter(
      (a) =>
        (!teams.value.length || teams.value.includes(a.team.name)) &&
        (!problems.value.length || problems.value.includes(a.problem.code)) &&
        (!recent.value || Date.parse(a.submitted_at) >= Date.now() - 1200000) &&
        (perfect.value || !a.score || a.score.total !== a.score.max) &&
        (!unscored.value || !a.score),
    )
    .toSorted((a, b) =>
      !a.score && !b.score
        ? Date.parse(a.submitted_at) - Date.parse(b.submitted_at)
        : !a.score
          ? -1
          : !b.score
            ? 1
            : Date.parse(b.submitted_at) - Date.parse(a.submitted_at),
    ),
);
</script>
<template>
  <main class="workspace">
    <h1>採点</h1>
    <div class="form-grid">
      <label
        >問題検索<select v-model="problems" multiple>
          <option
            v-for="p in new Set(data?.answers.map((a) => a.problem.code))"
            :key="p"
          >
            {{ p }}
          </option>
        </select></label
      ><label
        >チーム検索<select v-model="teams" multiple>
          <option
            v-for="t in new Set(data?.answers.map((a) => a.team.name))"
            :key="t"
          >
            {{ t }}
          </option>
        </select></label
      >
    </div>
    <div class="actions">
      <button class="button-secondary" @click="recent = !recent">
        {{ recent ? "直近20分以外も表示" : "直近20分のみ表示" }}</button
      ><button class="button-secondary" @click="perfect = !perfect">
        {{ perfect ? "満点解答を非表示" : "満点解答を表示" }}</button
      ><button class="button-secondary" @click="unscored = !unscored">
        {{ unscored ? "すべての回答を表示" : "未採点のみ表示" }}
      </button>
    </div>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><div class="table-scroll">
        <table class="data-table">
          <thead>
            <tr>
              <th>問題</th>
              <th>チーム</th>
              <th>解答ID</th>
              <th>提出時刻</th>
              <th>点数</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="a in filtered"
              :key="`${a.team.code}:${a.problem.code}:${a.reference.answer_number}`"
            >
              <td>
                <NuxtLink
                  :to="`/submissions/${a.problem.code}/${a.team.code}/${a.reference.answer_number}`"
                  >{{ a.problem.code }}: {{ a.problem.title }}</NuxtLink
                >
              </td>
              <td>{{ a.team.name }}</td>
              <td>{{ a.reference.answer_number }}</td>
              <td>{{ new Date(a.submitted_at).toLocaleString("ja-JP") }}</td>
              <td>
                {{
                  a.score
                    ? `${a.score.total} (${a.score.marked}-${a.score.penalty}) / ${a.score.max}`
                    : "未採点"
                }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <p v-if="!filtered?.length">該当する回答はありません。</p></RequestState
    >
  </main>
</template>
