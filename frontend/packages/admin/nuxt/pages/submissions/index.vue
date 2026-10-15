<script setup lang="ts">
import { api, expectData, type components } from "@ictsc/api";
import { updateAdminAnswerWorkflow } from "~/features/admin-api";

type Answer = components["schemas"]["AdminAnswer"];
type Status = Answer["workflow"]["status"];
const columns: { status: Status; label: string }[] = [
  { status: "WAITING", label: "採点待ち" },
  { status: "IN_PROGRESS", label: "採点中" },
  { status: "COMPLETED", label: "採点完了" },
];
useHead({ title: "採点" });
const { viewer } = useAdminSession();
const myDiscordID = computed(() =>
  viewer.value?.state === "ADMIN" ? viewer.value.admin.discord.id : "",
);
const { data, error, pending, refresh } = await useAsyncData(
  "admin-answers",
  async () =>
    expectData(
      await api.GET("/api/v1/admin/answers", {
        params: { query: { include_marked: true } },
      }),
    ),
);
const teams = ref<string[]>([]);
const problems = ref<string[]>([]);
const recent = ref(false);
const perfect = ref(true);
const unscored = ref(false);
const busyKey = ref("");
const actionError = ref("");
const { now, elapsedMinutes } = useSubmissionAge();

onMounted(() => {
  try {
    const saved = JSON.parse(
      localStorage.getItem("regalia/admin/answer-filters") ?? "null",
    );
    if (saved) {
      teams.value = saved.teams ?? [];
      problems.value = saved.problems ?? [];
      recent.value = !!saved.recent;
      perfect.value = saved.perfect !== false;
      unscored.value = !!saved.unscored;
    }
  } catch {
    /* Filters work without local storage. */
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
      /* Filters work without local storage. */
    }
  },
  { deep: true },
);

function key(answer: Answer) {
  const r = answer.reference;
  return `${r.team_code}:${r.problem_code}:${r.answer_number}`;
}
function url(answer: Answer) {
  const r = answer.reference;
  return `/submissions/${r.problem_code}/${r.team_code}/${r.answer_number}`;
}
const filtered = computed(() =>
  (data.value?.answers ?? []).filter(
    (a) =>
      (!teams.value.length || teams.value.includes(a.team.name)) &&
      (!problems.value.length || problems.value.includes(a.problem.code)) &&
      (!recent.value || Date.parse(a.submitted_at) >= now.value - 1200000) &&
      (perfect.value || !a.score || a.score.total !== a.score.max) &&
      (!unscored.value || !a.score),
  ),
);
function inColumn(status: Status) {
  return filtered.value
    .filter((a) => a.workflow.status === status)
    .toSorted((a, b) => {
      const mineA = Number(
        a.workflow.assignee_discord_id === myDiscordID.value &&
          !!myDiscordID.value,
      );
      const mineB = Number(
        b.workflow.assignee_discord_id === myDiscordID.value &&
          !!myDiscordID.value,
      );
      return (
        mineB - mineA ||
        Date.parse(b.submitted_at) - Date.parse(a.submitted_at) ||
        key(a).localeCompare(key(b))
      );
    });
}
async function update(
  answer: Answer,
  change: { status?: Status; assignment?: "CLAIM_SELF" | "RESET_TO_DEFAULT" },
) {
  if (change.status === "COMPLETED" && !answer.score) {
    await navigateTo(url(answer));
    return;
  }
  busyKey.value = key(answer);
  actionError.value = "";
  try {
    const updated = await updateAdminAnswerWorkflow(api, answer, change);
    if (data.value)
      data.value = {
        ...data.value,
        answers: data.value.answers.map((a) =>
          key(a) === key(answer) ? updated : a,
        ),
      };
  } catch (cause) {
    actionError.value =
      cause instanceof Error ? cause.message : "更新に失敗しました。";
    await refresh();
  } finally {
    busyKey.value = "";
  }
}
function onDragStart(event: DragEvent, answer: Answer) {
  event.dataTransfer?.setData("text/plain", key(answer));
  if (event.dataTransfer) event.dataTransfer.effectAllowed = "move";
}
function onDrop(event: DragEvent, status: Status) {
  const id = event.dataTransfer?.getData("text/plain");
  const answer = data.value?.answers.find((a) => key(a) === id);
  if (answer && answer.workflow.status !== status)
    void update(answer, { status });
}
</script>

<template>
  <main class="workspace grading-workspace">
    <h1>採点</h1>
    <p class="admin-lead">
      自分の担当回答を各列の先頭に表示します。カードを開くと採点できます。
    </p>
    <section class="admin-panel">
      <h2>回答を絞り込む</h2>
      <div class="admin-filter-grid">
        <label
          >問題（複数選択可）<select v-model="problems" multiple size="4">
            <option
              v-for="p in new Set(data?.answers.map((a) => a.problem.code))"
              :key="p"
            >
              {{ p }}
            </option>
          </select></label
        >
        <label
          >チーム（複数選択可）<select v-model="teams" multiple size="4">
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
        <label><input v-model="recent" type="checkbox" /> 直近20分のみ</label>
        <label
          ><input v-model="perfect" type="checkbox" /> 満点回答を表示</label
        >
        <label><input v-model="unscored" type="checkbox" /> 未採点のみ</label>
        <button
          class="button-secondary"
          @click="
            teams = [];
            problems = [];
            recent = false;
            perfect = true;
            unscored = false;
          "
        >
          条件をクリア
        </button>
      </div>
    </section>
    <p v-if="actionError" role="alert">{{ actionError }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh">
      <div class="grading-board">
        <section
          v-for="column in columns"
          :key="column.status"
          class="grading-column"
          :aria-label="column.label"
          @dragover.prevent
          @drop.prevent="onDrop($event, column.status)"
        >
          <header class="grading-column-header">
            <h2>{{ column.label }}</h2>
            <span>{{ inColumn(column.status).length }} 件</span>
          </header>
          <div
            v-for="answer in inColumn(column.status)"
            :key="key(answer)"
            class="grading-card"
            draggable="true"
            @dragstart="onDragStart($event, answer)"
          >
            <NuxtLink class="grading-card-link" :to="url(answer)">
              <strong
                >{{ answer.problem.code }}: {{ answer.problem.title }}</strong
              >
              <span
                >{{ answer.team.name }} / #{{
                  answer.reference.answer_number
                }}</span
              >
              <span class="grading-card-submitted">
                <time :datetime="answer.submitted_at">{{
                  new Date(answer.submitted_at).toLocaleString("ja-JP")
                }}</time>
                <span class="grading-card-age"
                  >提出から {{ elapsedMinutes(answer.submitted_at) }} 分</span
                >
              </span>
              <span v-if="answer.score"
                >{{ answer.score.total }} / {{ answer.score.max }} 点</span
              >
            </NuxtLink>
            <div class="grading-card-meta">
              <span
                v-if="
                  answer.workflow.assignee_discord_id === myDiscordID &&
                  myDiscordID
                "
                >自分の担当</span
              >
              <span v-else
                >担当:
                {{ answer.workflow.assignee_discord_id ?? "未設定" }}</span
              >
              <span v-if="answer.workflow.assignment_source === 'CLAIMED'"
                >個別担当</span
              >
            </div>
            <div class="grading-card-actions">
              <label
                >移動先<select
                  :value="answer.workflow.status"
                  :disabled="busyKey === key(answer)"
                  @change="
                    update(answer, {
                      status: ($event.target as HTMLSelectElement)
                        .value as Status,
                    })
                  "
                >
                  <option
                    v-for="target in columns"
                    :key="target.status"
                    :value="target.status"
                  >
                    {{ target.label }}
                  </option>
                </select></label
              >
              <button
                v-if="answer.workflow.assignment_source === 'DEFAULT'"
                type="button"
                class="button-secondary"
                :disabled="busyKey === key(answer)"
                @click="update(answer, { assignment: 'CLAIM_SELF' })"
              >
                自分が引き受ける
              </button>
              <button
                v-else
                type="button"
                class="button-secondary"
                :disabled="busyKey === key(answer)"
                @click="update(answer, { assignment: 'RESET_TO_DEFAULT' })"
              >
                担当を戻す
              </button>
            </div>
          </div>
          <p v-if="!inColumn(column.status).length" class="admin-empty">
            該当する回答はありません。
          </p>
        </section>
      </div>
    </RequestState>
  </main>
</template>

<style scoped>
.grading-workspace {
  max-width: none;
}
.grading-board {
  display: grid;
  width: 100%;
  min-width: 0;
  grid-template-columns: repeat(3, minmax(300px, 1fr));
  gap: 16px;
  align-items: start;
  overflow-x: auto;
}
.grading-column {
  min-width: 0;
  min-height: 340px;
  padding: 16px;
  border: 1px solid var(--line-light);
  background: var(--surface-hover);
}
.grading-column-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 16px;
}
.grading-column-header h2 {
  margin: 0;
}
.grading-card {
  margin-bottom: 12px;
  padding: 14px;
  border: 1px solid var(--line-light);
  background: var(--paper);
}
.grading-card:focus-within,
.grading-card:hover {
  border-color: var(--accent);
}
.grading-card-link {
  display: grid;
  gap: 6px;
  color: inherit;
  text-decoration: none;
}
.grading-card-link:hover strong {
  text-decoration: underline;
}
.grading-card-link span,
.grading-card-meta {
  font-size: 12px;
  color: var(--muted);
}
.grading-card-submitted {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px 10px;
}
.grading-card-link .grading-card-age {
  color: var(--accent);
  font-weight: 700;
  font-variant-numeric: tabular-nums;
}
.grading-card-meta {
  display: flex;
  gap: 10px;
  margin: 12px 0;
}
.grading-card-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  gap: 8px;
}
.grading-card-actions label {
  display: grid;
  gap: 4px;
  font-size: 12px;
}
.grading-card-actions select {
  max-width: 135px;
  padding: 6px;
}
.grading-card-actions button {
  padding: 6px 8px;
  min-height: 0;
  font-size: 12px;
}
</style>
