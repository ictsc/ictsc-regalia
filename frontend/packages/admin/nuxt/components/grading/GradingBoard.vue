<script setup lang="ts">
import {
  answerKey,
  answerPath,
  columnAnswers,
  gradingColumns,
  markingBlockedLabels,
  minutesSince,
  type GradingAnswer,
  type GradingStatus,
} from "~/features/grading";

export type GradingChange = {
  status?: GradingStatus;
  assignment?: "CLAIM_SELF" | "RESET_TO_DEFAULT";
};

const props = defineProps<{
  /** Answers after filtering. */
  answers: GradingAnswer[];
  /** The signed-in admin; their answers are listed first. */
  myDiscordId: string;
  /** Key of the answer being updated, which disables its controls. */
  busyKey: string;
  /** Current time used for the "minutes since submission" label. */
  now: number;
}>();
const emit = defineEmits<{
  update: [answer: GradingAnswer, change: GradingChange];
}>();

const inColumn = (status: GradingStatus) =>
  columnAnswers(props.answers, status, props.myDiscordId);
function onDragStart(event: DragEvent, answer: GradingAnswer) {
  event.dataTransfer?.setData("text/plain", answerKey(answer));
  if (event.dataTransfer) event.dataTransfer.effectAllowed = "move";
}
function onDrop(event: DragEvent, status: GradingStatus) {
  const id = event.dataTransfer?.getData("text/plain");
  const answer = props.answers.find((a) => answerKey(a) === id);
  if (answer && answer.workflow.status !== status)
    emit("update", answer, { status });
}
</script>
<template>
  <div class="grading-board">
    <section
      v-for="column in gradingColumns"
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
        :key="answerKey(answer)"
        class="grading-card"
        draggable="true"
        @dragstart="onDragStart($event, answer)"
      >
        <NuxtLink class="grading-card-link" :to="answerPath(answer)">
          <strong>{{ answer.problem.code }}: {{ answer.problem.title }}</strong>
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
              >提出から {{ minutesSince(answer.submitted_at, now) }} 分</span
            >
          </span>
          <span v-if="answer.score"
            >{{ answer.score.total }} / {{ answer.score.max }} 点</span
          >
          <span
            v-if="answer.marking_blocked_reason"
            class="grading-card-blocked"
            >{{ markingBlockedLabels[answer.marking_blocked_reason] }}</span
          >
        </NuxtLink>
        <div class="grading-card-meta">
          <span
            v-if="
              answer.workflow.assignee_discord_id === myDiscordId && myDiscordId
            "
            >自分の担当</span
          >
          <span v-else
            >担当: {{ answer.workflow.assignee_discord_id ?? "未設定" }}</span
          >
          <span v-if="answer.workflow.assignment_source === 'CLAIMED'"
            >個別担当</span
          >
        </div>
        <div class="grading-card-actions">
          <label
            >移動先<select
              :value="answer.workflow.status"
              :disabled="busyKey === answerKey(answer)"
              @change="
                $emit('update', answer, {
                  status: ($event.target as HTMLSelectElement)
                    .value as GradingStatus,
                })
              "
            >
              <option
                v-for="target in gradingColumns"
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
            :disabled="busyKey === answerKey(answer)"
            @click="$emit('update', answer, { assignment: 'CLAIM_SELF' })"
          >
            自分が引き受ける
          </button>
          <button
            v-else
            type="button"
            class="button-secondary"
            :disabled="busyKey === answerKey(answer)"
            @click="$emit('update', answer, { assignment: 'RESET_TO_DEFAULT' })"
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
</template>

<style scoped>
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
.grading-card-link .grading-card-blocked {
  justify-self: start;
  padding: 1px 8px;
  border: 1px solid var(--line-light);
  border-radius: 999px;
  color: var(--muted);
  font-size: 12px;
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
