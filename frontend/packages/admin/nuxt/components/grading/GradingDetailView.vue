<script setup lang="ts">
import { computed, ref, watch } from "vue";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import AttachmentList from "@ictsc/ui/components/AttachmentList.vue";
import type { AttachmentFile } from "@ictsc/ui/attachments";
import type { components } from "@ictsc/api";
import {
  answerSourceLines,
  buildAnswerDiffHunks,
  compareAnswerSources,
} from "~/features/answer-diff";
import { buildDisplayItems } from "~/features/answer-display";
import { adminImageUrl } from "~/features/admin-images";
import { resizePaneWidths } from "~/features/grading-panes";
import {
  activeCommentCount,
  clearLineDrafts,
  draftLineComments as toLineComments,
  emptyLineDraftState,
  lineHistory as groupLineHistory,
  rangeLabel,
  removeDraft as withoutDraft,
  selectLine as withSelectedLine,
} from "~/features/line-drafts";

type Schemas = components["schemas"];
export type MarkingInput = {
  score: number;
  comment: string;
  lineComments: Schemas["AnswerLineComment"][];
};

const props = defineProps<{
  answer: Schemas["AdminAnswer"];
  /** The problem as of the answer's content commit. */
  problem: Schemas["AdminProblem"];
  /** Commit currently served to contestants. */
  activeCommit: string;
  /** Earlier answers to the same problem, newest first. */
  previous: Schemas["AdminAnswer"][];
  marks: Schemas["MarkingResult"][];
  /** Text logs submitted with the answer. */
  attachments?: AttachmentFile[];
  /** A marking or comment deletion is running. */
  busy: boolean;
  /** The assignment is being updated. */
  workflowBusy: boolean;
  /** Saves a marking; resolves true when it was saved. */
  submitMarking: (input: MarkingInput) => Promise<boolean>;
}>();
defineEmits<{
  assignment: [action: "CLAIM_SELF" | "RESET_TO_DEFAULT"];
  deleteLineComment: [markingId: string, commentIndex: number];
}>();
const collapsed = defineModel<boolean[]>("collapsed", { required: true });
const widths = defineModel<number[]>("widths", { required: true });

const selectedPrevious = ref<number | null>(null);
watch(
  () => props.previous,
  (answers) => {
    selectedPrevious.value = answers[0]?.reference.answer_number ?? null;
  },
  { immediate: true },
);
const previousAnswer = computed(() =>
  props.previous.find(
    (a) => a.reference.answer_number === selectedPrevious.value,
  ),
);
const score = ref<number | null>(null);
const comment = ref("");
const draftState = ref(emptyLineDraftState());
const lineDrafts = computed(() => draftState.value.drafts);
const activeDraftId = computed({
  get: () => draftState.value.activeId,
  set: (activeId) => {
    draftState.value = { ...draftState.value, activeId };
  },
});
const answerLines = computed(() => answerSourceLines(props.answer.body.body));
const comparison = computed(() =>
  previousAnswer.value
    ? compareAnswerSources(
        previousAnswer.value.body.body,
        props.answer.body.body,
      )
    : null,
);
const answerView = ref<"diff" | "current">("diff");
const showAllContext = ref(false);
const isDiffView = computed(
  () => Boolean(previousAnswer.value) && answerView.value === "diff",
);
const diffHunks = computed(() =>
  comparison.value
    ? buildAnswerDiffHunks(
        comparison.value.rows,
        showAllContext.value ? Infinity : 3,
      )
    : [],
);
const displayedItems = computed(() =>
  buildDisplayItems(
    answerLines.value,
    comparison.value,
    diffHunks.value,
    isDiffView.value,
  ),
);
watch(selectedPrevious, () => {
  showAllContext.value = false;
});
const lineHistory = computed(() => groupLineHistory(props.marks));
const savedRanges = computed(() => Object.values(lineHistory.value).flat());
const draftLineComments = computed(() => toLineComments(lineDrafts.value));
function draftsAtLine(line: number) {
  return lineDrafts.value.filter((draft) => draft.endLine === line);
}
function selectLine(line: number, extend: boolean) {
  draftState.value = withSelectedLine(draftState.value, line, extend);
}
function removeDraft(id: number) {
  draftState.value = withoutDraft(draftState.value, id);
}
// Images follow the problem version the answer was submitted against.
const problemImageUrl = computed(() =>
  adminImageUrl(props.problem.content_commit),
);
// 二次予選は最終答案だけを採点し比較もしないので、回答は Markdown で読むだけにする。
const markdownOnly = computed(
  () => props.answer.competition_mode === "SECONDARY_QUALIFIER",
);
const blockedMessages: Record<Schemas["MarkingBlockedReason"], string> = {
  grading_not_started: "二次予選モードでは競技終了後に採点できます。",
  answer_not_latest:
    "二次予選モードでは最新の回答だけを採点します。この回答は採点対象外です。",
};
const blockedMessage = computed(() =>
  props.answer.marking_blocked_reason
    ? blockedMessages[props.answer.marking_blocked_reason]
    : "",
);
async function mark() {
  if (score.value == null || blockedMessage.value) return;
  const saved = await props.submitMarking({
    score: score.value,
    comment: comment.value,
    lineComments: draftLineComments.value,
  });
  if (saved) draftState.value = clearLineDrafts(draftState.value);
}

const paneLabels = ["提出時点の問題", "回答", "採点"];
const visiblePaneIndices = computed(() => [0, 1, 2]);
const paneContainer = ref<HTMLElement | null>(null);
function toggle(index: number) {
  collapsed.value = collapsed.value.map((value, i) =>
    i === index ? !value : value,
  );
}
function startResize(
  event: PointerEvent,
  leftIndex: number,
  rightIndex: number,
) {
  if (
    !paneContainer.value ||
    collapsed.value[leftIndex] ||
    collapsed.value[rightIndex]
  )
    return;
  event.preventDefault();
  const startX = event.clientX;
  const start = [...widths.value];
  const openWeight = visiblePaneIndices.value
    .filter((i) => !collapsed.value[i])
    .reduce((sum, i) => sum + start[i]!, 0);
  const available =
    paneContainer.value.clientWidth -
    (visiblePaneIndices.value.length - 1) * 6 -
    visiblePaneIndices.value.filter((i) => collapsed.value[i]).length * 46;
  const move = (moveEvent: PointerEvent) => {
    const delta = ((moveEvent.clientX - startX) / available) * openWeight;
    widths.value = resizePaneWidths(start, leftIndex, rightIndex, delta);
  };
  const stop = () => {
    window.removeEventListener("pointermove", move);
    window.removeEventListener("pointerup", stop);
  };
  window.addEventListener("pointermove", move);
  window.addEventListener("pointerup", stop, { once: true });
}
</script>
<template>
  <div class="grading-detail-meta">
    <span
      >現在の得点:
      <strong
        >{{ answer.score?.total ?? "未採点" }} /
        {{ answer.problem.max_score }}</strong
      ></span
    >
    <span>担当: {{ answer.workflow.assignee_discord_id ?? "未設定" }}</span>
    <button
      v-if="answer.workflow.assignment_source === 'DEFAULT'"
      type="button"
      class="button-secondary"
      :disabled="workflowBusy"
      @click="$emit('assignment', 'CLAIM_SELF')"
    >
      自分が引き受ける
    </button>
    <button
      v-else
      type="button"
      class="button-secondary"
      :disabled="workflowBusy"
      @click="$emit('assignment', 'RESET_TO_DEFAULT')"
    >
      担当を戻す
    </button>
  </div>
  <div ref="paneContainer" class="grading-panes">
    <template v-for="(index, position) in visiblePaneIndices" :key="index">
      <div
        v-if="position > 0"
        class="grading-resizer"
        role="separator"
        aria-orientation="vertical"
        :aria-label="`${paneLabels[visiblePaneIndices[position - 1]!]}と${paneLabels[index]}の幅を変更`"
        @pointerdown="
          startResize($event, visiblePaneIndices[position - 1]!, index)
        "
      />
      <section
        class="grading-pane"
        :class="{ 'is-collapsed': collapsed[index] }"
        :style="{
          flex: collapsed[index] ? '0 0 46px' : `${widths[index]} 1 0px`,
        }"
      >
        <header class="grading-pane-header">
          <h2>{{ paneLabels[index] }}</h2>
          <button
            type="button"
            class="button-secondary"
            :aria-label="`${paneLabels[index]}を${collapsed[index] ? '展開' : '折りたたむ'}`"
            @click="toggle(index)"
          >
            {{ collapsed[index] ? "開く" : "閉じる" }}
          </button>
        </header>
        <div v-if="!collapsed[index]" class="grading-pane-content">
          <template v-if="index === 0">
            <p
              v-if="answer.content_commit !== activeCommit"
              class="admin-badge"
            >
              提出時点の版を表示中
            </p>
            <MarkdownContent
              :source="problem.body"
              :image-url="problemImageUrl"
              diagrams
            />
            <details v-if="problem.explanation">
              <summary>問題解説</summary>
              <MarkdownContent
                :source="problem.explanation"
                :image-url="problemImageUrl"
                diagrams
              />
            </details>
            <details>
              <summary>提出時の commit</summary>
              <code>{{ answer.content_commit }}</code>
            </details>
          </template>
          <template v-else-if="index === 1">
            <div
              class="grading-answer-toolbar"
              :class="{ 'is-plain': markdownOnly }"
            >
              <div class="grading-answer-toolbar-top">
                <div>
                  <small class="grading-answer-eyebrow">提出された回答</small>
                  <h3>今回の回答 #{{ answer.reference.answer_number }}</h3>
                </div>
                <div
                  v-if="previous.length"
                  class="grading-view-switch"
                  role="group"
                  aria-label="回答の表示方法"
                >
                  <button
                    type="button"
                    :aria-pressed="isDiffView"
                    @click="answerView = 'diff'"
                  >
                    差分
                  </button>
                  <button
                    type="button"
                    :aria-pressed="!isDiffView"
                    @click="answerView = 'current'"
                  >
                    今回の回答のみ
                  </button>
                </div>
              </div>
              <div v-if="previous.length" class="grading-compare-controls">
                <label class="grading-previous-picker">
                  比較する回答
                  <select v-model.number="selectedPrevious">
                    <option
                      v-for="old in previous"
                      :key="old.reference.answer_number"
                      :value="old.reference.answer_number"
                    >
                      #{{ old.reference.answer_number }} ·
                      {{ old.score?.total ?? "未採点" }} 点
                    </option>
                  </select>
                </label>
                <p v-if="previousAnswer" class="grading-previous-meta">
                  <span
                    >過去 #{{ previousAnswer.reference.answer_number }} ·
                    {{ previousAnswer.score?.total ?? "未採点"
                    }}<template v-if="previousAnswer.score">
                      / {{ previousAnswer.problem.max_score }}</template
                    >
                    点</span
                  >
                  <time :datetime="previousAnswer.submitted_at">{{
                    new Date(previousAnswer.submitted_at).toLocaleString(
                      "ja-JP",
                    )
                  }}</time>
                </p>
              </div>
            </div>
            <MarkdownContent
              v-if="markdownOnly"
              class="grading-markdown-answer"
              :source="answer.body.body"
              block-external-images
              diagrams
            />
            <template v-else>
              <div v-if="isDiffView && comparison" class="grading-diff-filebar">
                <span
                  >過去 #{{ previousAnswer?.reference.answer_number }} → 今回
                  #{{ answer.reference.answer_number }}</span
                >
                <span class="grading-diff-stats"
                  ><b>+{{ comparison.addedLines }}</b
                  ><b>−{{ comparison.removedLines }}</b></span
                >
              </div>
              <p
                v-if="isDiffView && comparison && !diffHunks.length"
                class="grading-no-diff"
              >
                選択した過去回答から変更はありません。「今回の回答のみ」で全文を確認できます。
              </p>
              <div
                v-else
                class="grading-answer-lines"
                :class="{ 'is-unified-diff': isDiffView }"
              >
                <template v-for="item in displayedItems" :key="item.key">
                  <button
                    v-if="item.type === 'omitted'"
                    type="button"
                    class="grading-diff-omitted"
                    @click="showAllContext = true"
                  >
                    … {{ item.count }} 行を展開
                  </button>
                  <div
                    v-else-if="item.type === 'hunk'"
                    class="grading-diff-hunk"
                  >
                    @@ -{{ item.hunk.previousStart }},{{
                      item.hunk.previousCount
                    }}
                    +{{ item.hunk.currentStart }},{{ item.hunk.currentCount }}
                    @@
                  </div>
                  <div
                    v-else
                    class="grading-answer-line"
                    :class="`is-${item.row.kind}`"
                  >
                    <div
                      v-if="item.row.currentLine === null"
                      class="grading-line-source is-removed"
                    >
                      <span class="grading-old-line-number">{{
                        item.row.previousLine
                      }}</span>
                      <span class="grading-empty-line-number" />
                      <span class="grading-current-marker" aria-hidden="true"
                        >−</span
                      >
                      <span class="grading-line-text">{{
                        item.row.text || " "
                      }}</span>
                    </div>
                    <div
                      v-else
                      class="grading-line-source"
                      :class="{
                        'is-selected': lineDrafts.some(
                          (draft) =>
                            draft.id === activeDraftId &&
                            item.row.currentLine! >= draft.startLine &&
                            item.row.currentLine! <= draft.endLine,
                        ),
                        'has-draft': lineDrafts.some(
                          (draft) =>
                            draft.body.trim() &&
                            item.row.currentLine! >= draft.startLine &&
                            item.row.currentLine! <= draft.endLine,
                        ),
                        'has-history': savedRanges.some(
                          (note) =>
                            item.row.currentLine! >= note.startLine &&
                            item.row.currentLine! <= note.endLine,
                        ),
                      }"
                      @click="
                        selectLine(item.row.currentLine!, $event.shiftKey)
                      "
                    >
                      <span v-if="isDiffView" class="grading-old-line-number">{{
                        item.row.previousLine ?? ""
                      }}</span>
                      <button
                        type="button"
                        class="grading-line-number"
                        :aria-label="`行 ${item.row.currentLine} を選択`"
                        @click.stop="
                          selectLine(item.row.currentLine!, $event.shiftKey)
                        "
                      >
                        {{ item.row.currentLine }}
                      </button>
                      <span
                        v-if="isDiffView"
                        class="grading-current-marker"
                        aria-hidden="true"
                        >{{ item.row.kind === "added" ? "+" : " " }}</span
                      >
                      <span class="grading-line-text">{{
                        item.row.text || " "
                      }}</span>
                    </div>
                    <div
                      v-for="draft in item.row.currentLine === null
                        ? []
                        : draftsAtLine(item.row.currentLine)"
                      :key="draft.id"
                      class="grading-line-editor"
                    >
                      <template v-if="draft.id === activeDraftId">
                        <label :for="`line-comment-${draft.id}`"
                          >{{
                            rangeLabel(draft.startLine, draft.endLine)
                          }}へのコメント</label
                        >
                        <textarea
                          :id="`line-comment-${draft.id}`"
                          v-model="draft.body"
                          maxlength="2000"
                          rows="3"
                          placeholder="この行についてコメント"
                        />
                        <button
                          type="button"
                          class="button-secondary"
                          @click="removeDraft(draft.id)"
                        >
                          コメントを破棄
                        </button>
                      </template>
                      <button
                        v-else
                        type="button"
                        class="grading-line-draft"
                        @click="activeDraftId = draft.id"
                      >
                        <strong
                          >{{
                            rangeLabel(draft.startLine, draft.endLine)
                          }}の下書き</strong
                        >
                        <span>{{ draft.body }}</span>
                      </button>
                    </div>
                    <div
                      v-for="(note, noteIndex) in item.row.currentLine === null
                        ? []
                        : (lineHistory[item.row.currentLine] ?? [])"
                      :key="`${note.id}-${noteIndex}`"
                      class="grading-line-history"
                    >
                      <small>{{
                        rangeLabel(note.startLine, note.endLine)
                      }}</small
                      ><br /><strong>{{ note.judge }}</strong
                      >: {{ note.body }}
                    </div>
                  </div>
                </template>
              </div>
              <details class="grading-markdown-preview">
                <summary>今回の回答を Markdown で表示</summary>
                <MarkdownContent
                  :source="answer.body.body"
                  block-external-images
                  diagrams
                />
              </details>
            </template>
            <section
              v-if="attachments?.length"
              class="grading-attachments"
              aria-label="添付ログ"
            >
              <h3>添付ログ</h3>
              <AttachmentList :attachments="attachments" />
            </section>
          </template>
          <template v-else>
            <p v-if="blockedMessage" class="grading-blocked" role="note">
              {{ blockedMessage }}
            </p>
            <form v-else class="form-grid" @submit.prevent="mark">
              <label
                >得点<input
                  v-model.number="score"
                  type="number"
                  min="0"
                  :max="answer.problem.max_score"
                  required
              /></label>
              <label>全体コメント<textarea v-model="comment" /></label>
              <p v-if="draftLineComments.length">
                行コメント
                {{ draftLineComments.length }} 件を含めて保存します。
              </p>
              <button class="button-primary" :disabled="busy">
                採点結果を送信
              </button>
            </form>
            <h3>採点履歴</h3>
            <ol class="grading-history">
              <li v-for="marking in marks" :key="marking.id">
                <strong>{{ marking.score }} 点</strong> · {{ marking.judge.name
                }}<br />
                <small>{{
                  new Date(marking.created_at).toLocaleString("ja-JP")
                }}</small>
                <p>{{ marking.rationale }}</p>
                <details
                  v-if="marking.line_comments.length"
                  class="grading-history-comments"
                >
                  <summary>
                    行コメント {{ activeCommentCount(marking) }} 件
                    <span
                      v-if="
                        activeCommentCount(marking) !==
                        marking.line_comments.length
                      "
                      >（削除済み
                      {{
                        marking.line_comments.length -
                        activeCommentCount(marking)
                      }}
                      件）</span
                    >
                  </summary>
                  <ol>
                    <li
                      v-for="(note, noteIndex) in marking.line_comments"
                      :key="noteIndex"
                      :class="{ 'is-deleted': note.deleted_at }"
                    >
                      <div class="grading-history-comment-head">
                        <strong>{{
                          rangeLabel(
                            note.line_number,
                            note.end_line_number ?? note.line_number,
                          )
                        }}</strong>
                        <button
                          v-if="!note.deleted_at"
                          type="button"
                          class="button-secondary"
                          :disabled="busy"
                          :aria-label="`${rangeLabel(note.line_number, note.end_line_number ?? note.line_number)}のコメントを削除`"
                          @click="
                            $emit('deleteLineComment', marking.id, noteIndex)
                          "
                        >
                          削除
                        </button>
                      </div>
                      <p
                        v-if="note.deleted_at"
                        class="grading-history-comment-deleted"
                      >
                        削除済み · {{ note.deleted_by }} ·
                        {{ new Date(note.deleted_at).toLocaleString("ja-JP") }}
                      </p>
                      <p v-else>{{ note.body }}</p>
                    </li>
                  </ol>
                </details>
              </li>
            </ol>
            <p v-if="!marks.length" class="admin-empty">
              採点履歴はありません。
            </p>
          </template>
        </div>
      </section>
    </template>
  </div>
</template>

<style scoped>
.grading-detail-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px 22px;
  margin-bottom: 20px;
}
.grading-panes {
  display: flex;
  width: 100%;
  min-height: min(70vh, 800px);
  overflow-x: auto;
  border: 1px solid var(--line-light);
}
.grading-pane {
  min-width: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.grading-pane:not(.is-collapsed) {
  min-width: 220px;
}
.grading-pane + .grading-resizer,
.grading-resizer + .grading-pane {
  border-left: 1px solid var(--line-light);
}
.grading-pane-header {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid var(--line-light);
}
.grading-pane-header h2 {
  margin: 0;
  font-size: 16px;
}
.grading-pane-header button {
  min-height: 0;
  padding: 5px;
  font-size: 12px;
}
.grading-pane-content {
  min-width: 0;
  overflow: auto;
  padding: 16px;
  flex: 1;
  overflow-wrap: anywhere;
}
.grading-pane-content pre {
  white-space: pre-wrap;
}
.grading-pane-content select {
  display: block;
  max-width: 100%;
  margin: 8px 0 16px;
}
.grading-answer-toolbar {
  display: grid;
  gap: 14px;
  padding: 14px 16px;
  border: 1px solid var(--line-light);
  border-bottom: 0;
  background: var(--surface-hover);
}
.grading-answer-toolbar-top,
.grading-compare-controls,
.grading-previous-meta,
.grading-diff-filebar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 10px 18px;
}
.grading-answer-eyebrow {
  color: var(--muted);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.06em;
}
.grading-answer-toolbar h3 {
  margin: 2px 0 0;
  font-size: 17px;
}
.grading-view-switch {
  display: inline-flex;
  gap: 2px;
  padding: 3px;
  border: 1px solid var(--line-light);
  border-radius: 8px;
  background: var(--paper);
}
.grading-view-switch button {
  min-height: 32px;
  padding: 5px 10px;
  border: 0;
  border-radius: 5px;
  background: transparent;
  color: var(--muted);
  font: inherit;
  font-size: 12px;
  cursor: pointer;
}
.grading-view-switch button[aria-pressed="true"] {
  background: var(--accent);
  color: var(--accent-interactive-text);
  font-weight: 700;
}
.grading-view-switch button:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}
.grading-compare-controls {
  padding-top: 12px;
  border-top: 1px solid var(--line-light);
}
.grading-previous-picker {
  display: grid;
  gap: 4px;
  min-width: min(100%, 220px);
  font-size: 12px;
  font-weight: 700;
}
.grading-previous-picker select {
  width: 100%;
  max-width: 100%;
  margin: 0;
  font-size: 13px;
}
.grading-previous-meta {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
}
.grading-diff-filebar {
  padding: 9px 14px;
  border: 1px solid var(--line-light);
  background: var(--paper);
  font-size: 12px;
  font-weight: 700;
}
.grading-diff-stats {
  display: inline-flex;
  gap: 12px;
  font-family: var(--mono);
}
.grading-diff-stats b:first-child {
  color: #208454;
}
.grading-diff-stats b:last-child {
  color: var(--accent);
}
.grading-no-diff {
  margin: 0;
  padding: 24px 16px;
  border: 1px solid var(--line-light);
  background: var(--paper);
  color: var(--muted);
}
.grading-answer-lines {
  border: 1px solid var(--line-light);
  border-top: 0;
  font-family: var(--sans-read);
  font-size: 13px;
  line-height: 1.65;
}
.grading-diff-hunk,
.grading-diff-omitted {
  display: block;
  width: 100%;
  padding: 7px 12px;
  border: 0;
  border-bottom: 1px solid var(--line-light);
  background: color-mix(in srgb, var(--accent) 7%, var(--paper));
  color: var(--muted);
  font-family: var(--mono);
  font-size: 12px;
  text-align: left;
}
.grading-diff-omitted {
  cursor: pointer;
}
.grading-diff-omitted:hover,
.grading-diff-omitted:focus-visible {
  color: var(--accent);
  text-decoration: underline;
  text-underline-offset: 3px;
}
.grading-blocked {
  margin: 0 0 16px;
  padding: 12px 14px;
  border-left: 3px solid var(--accent);
  background: var(--surface-hover);
  font-size: 13px;
  line-height: 1.6;
}
.grading-answer-line {
  border-bottom: 1px solid var(--line-light);
}
.grading-answer-line.is-added .grading-line-source {
  background: color-mix(in srgb, #28a36a 13%, var(--paper));
}
.grading-answer-line.is-removed .grading-line-source {
  background: color-mix(in srgb, var(--accent) 11%, var(--paper));
}
.grading-answer-line.is-added .grading-current-marker {
  color: #208454;
  font-weight: 700;
}
.grading-answer-line.is-removed .grading-current-marker {
  color: var(--accent);
  font-weight: 700;
}
.grading-line-source {
  display: grid;
  grid-template-columns: 42px minmax(0, 1fr);
  min-height: 31px;
  cursor: pointer;
}
.is-unified-diff .grading-line-source {
  grid-template-columns: 42px 42px 22px minmax(0, 1fr);
}
.grading-line-source.is-removed {
  cursor: default;
}
.grading-old-line-number,
.grading-empty-line-number {
  padding: 4px;
  border-right: 1px solid var(--line-light);
  color: var(--muted);
  font-family: var(--mono);
  text-align: right;
  font-variant-numeric: tabular-nums;
}
.grading-line-source.has-draft,
.grading-line-source.has-history,
.grading-line-source:focus-within {
  background: var(--surface-hover);
}
/* Lines with a comment are underlined, not only shaded. */
.grading-line-source.has-draft .grading-line-text,
.grading-line-source.has-history .grading-line-text {
  text-decoration: underline;
  text-decoration-style: dotted;
  text-underline-offset: 4px;
}
.grading-line-source.is-selected {
  background: color-mix(in srgb, var(--accent) 16%, var(--paper));
  box-shadow: inset 3px 0 var(--accent);
}
.grading-line-number {
  border: 0;
  border-right: 1px solid var(--line-light);
  padding: 4px;
  background: transparent;
  color: var(--accent);
  text-align: right;
  font: inherit;
  font-family: var(--mono);
  cursor: pointer;
}
.grading-current-marker {
  padding: 4px 0;
  text-align: center;
}
.grading-line-text {
  padding: 4px 8px;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}
.grading-line-history,
.grading-line-editor {
  margin: 0 8px 8px 50px;
  padding: 8px;
  background: var(--surface-hover);
  overflow-wrap: anywhere;
}
.is-unified-diff .grading-line-history,
.is-unified-diff .grading-line-editor {
  margin-left: 114px;
}
.grading-line-draft {
  display: grid;
  width: 100%;
  gap: 4px;
  border: 0;
  padding: 0;
  background: transparent;
  color: inherit;
  text-align: left;
}
.grading-line-editor label,
.grading-line-editor textarea {
  display: block;
  width: 100%;
  margin-bottom: 8px;
}
.grading-line-editor textarea {
  border: 1px solid var(--line-light);
  padding: 8px;
  color: var(--ink);
  background: var(--paper);
  font: inherit;
}
.grading-line-editor button {
  min-height: 0;
  padding: 5px 8px;
}
.grading-markdown-preview {
  margin-top: 20px;
}
.grading-attachments {
  margin-top: 20px;
}
.grading-attachments h3 {
  margin: 0 0 8px;
  font-size: 14px;
}
/* Markdown-only answers sit directly in the pane without an inner box. */
.grading-answer-toolbar.is-plain {
  padding: 0;
  border: 0;
  background: none;
}
.grading-markdown-answer {
  margin-top: 16px;
}
.grading-pane-content .form-grid {
  margin: 0 0 24px;
}
.grading-pane-content .form-grid textarea {
  min-height: 110px;
}
.grading-pane.is-collapsed .grading-pane-header {
  flex-direction: column;
}
.grading-pane.is-collapsed h2 {
  writing-mode: vertical-rl;
}
.grading-resizer {
  flex: 0 0 6px;
  background: var(--line-light);
  cursor: col-resize;
  touch-action: none;
}
.grading-history {
  padding-left: 20px;
}
.grading-history li {
  margin-bottom: 16px;
}
.grading-history-comments {
  margin-top: 10px;
  border: 1px solid var(--line-light);
  border-radius: 8px;
}
.grading-history-comments > summary {
  padding: 9px 12px;
  color: var(--accent);
  font-weight: 600;
  cursor: pointer;
}
.grading-history-comments > summary:hover,
.grading-history-comments > summary:focus-visible {
  text-decoration: underline;
  text-underline-offset: 3px;
}
.grading-history-comments > ol {
  margin: 0;
  padding: 4px 12px 6px 34px;
  border-top: 1px solid var(--line-light);
}
.grading-history-comments li {
  margin: 8px 0;
  padding-bottom: 8px;
  overflow-wrap: anywhere;
}
.grading-history-comments li + li {
  border-top: 1px solid var(--line-light);
  padding-top: 10px;
}
.grading-history-comment-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.grading-history-comment-head button {
  min-height: 30px;
  padding: 4px 10px;
}
.grading-history-comments li p {
  margin: 6px 0 0;
  white-space: pre-wrap;
}
.grading-history-comment-deleted {
  color: var(--muted);
  font-size: 12px;
}
@media (max-width: 760px) {
  .grading-panes {
    display: grid;
    min-height: 0;
    overflow: visible;
  }
  .grading-pane {
    flex-basis: auto !important;
    min-height: 48px;
    max-height: 65vh;
    border-bottom: 1px solid var(--line-light);
  }
  .grading-pane.is-collapsed h2 {
    writing-mode: initial;
  }
  .grading-pane.is-collapsed .grading-pane-header {
    flex-direction: row;
  }
  .grading-resizer {
    display: none;
  }
}
</style>
