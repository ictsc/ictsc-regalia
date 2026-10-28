<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import {
  answerLength,
  formatBytes,
  maxAttachments,
  readAttachment,
  uniqueAttachmentName,
  type AttachmentFile,
} from "@ictsc/ui/attachments";
import TabHint from "@ictsc/ui/components/TabHint.vue";
import UiIcon from "@ictsc/ui/components/UiIcon.vue";
import { useTabIndent } from "@ictsc/ui/textarea";
import AnswerConfirmDialog from "./AnswerConfirmDialog.vue";
import LogComposer from "./LogComposer.vue";

const body = defineModel<string>("body", { required: true });
const attachments = defineModel<AttachmentFile[]>("attachments", {
  default: () => [],
});
const props = defineProps<{
  /** Answers are being accepted now. */
  open: boolean;
  sending: boolean;
  submitted: boolean;
  /** Seconds until the team can answer again; 0 when not cooling down. */
  cooldownSeconds: number;
  /** ISO time of the last draft save in this browser. */
  savedAt: string;
  saveError: string;
  actionError?: unknown;
  /** Characters allowed, not counting line breaks. */
  maxLength: number;
  allowAttachments: boolean;
}>();
const emit = defineEmits<{ submit: [] }>();

const preview = ref(false);
const confirming = ref(false);
const composing = ref(false);
/** Files are dragged somewhere over the page; counts nested dragenters. */
const pageDragDepth = ref(0);
/** Files are dragged over the attachment area itself. */
const dragDepth = ref(0);
const dropErrors = ref<string[]>([]);
const editor = ref<HTMLTextAreaElement | null>(null);
const previewHeight = ref<string>();
const { onKeydown, onBlur, hintVisible } = useTabIndent();
const length = computed(() => answerLength(body.value));
const remaining = computed(() => props.maxLength - length.value);
const canSubmit = computed(
  () =>
    props.open &&
    props.cooldownSeconds <= 0 &&
    !props.sending &&
    body.value.trim() !== "" &&
    remaining.value >= 0,
);
const canAttach = computed(
  () =>
    props.open && !props.sending && attachments.value.length < maxAttachments,
);
const suggestedName = computed(() => {
  for (let index = 1; ; index++) {
    const name = `log-${index}.txt`;
    if (!attachments.value.some((attachment) => attachment.name === name))
      return name;
  }
});

// The preview keeps the editor's size and scrolls inside it.
watch(
  preview,
  (showing) => {
    if (showing && editor.value?.offsetHeight)
      previewHeight.value = `${editor.value.offsetHeight}px`;
  },
  { flush: "pre" },
);

function addAttachment(attachment: AttachmentFile) {
  attachments.value = [...attachments.value, attachment];
  composing.value = false;
}
function draggingFiles(event: DragEvent) {
  return event.dataTransfer?.types.includes("Files") ?? false;
}
// While files are dragged anywhere on the page, the attachment area shows
// itself as the drop target.
function onPageDragEnter(event: DragEvent) {
  if (draggingFiles(event)) pageDragDepth.value++;
}
function onPageDragLeave(event: DragEvent) {
  if (draggingFiles(event))
    pageDragDepth.value = Math.max(0, pageDragDepth.value - 1);
}
function onPageDragEnd() {
  pageDragDepth.value = 0;
  dragDepth.value = 0;
}
onMounted(() => {
  if (!props.allowAttachments) return;
  window.addEventListener("dragenter", onPageDragEnter);
  window.addEventListener("dragleave", onPageDragLeave);
  window.addEventListener("drop", onPageDragEnd);
  window.addEventListener("dragend", onPageDragEnd);
});
onBeforeUnmount(() => {
  window.removeEventListener("dragenter", onPageDragEnter);
  window.removeEventListener("dragleave", onPageDragLeave);
  window.removeEventListener("drop", onPageDragEnd);
  window.removeEventListener("dragend", onPageDragEnd);
});

// Only the attachment area takes drops; these run on it.
function onDragEnter(event: DragEvent) {
  if (!draggingFiles(event)) return;
  event.preventDefault();
  dragDepth.value++;
}
function onDragOver(event: DragEvent) {
  if (!draggingFiles(event)) return;
  // Without this the browser opens the file instead of dropping it here.
  event.preventDefault();
  if (event.dataTransfer)
    event.dataTransfer.dropEffect = canAttach.value ? "copy" : "none";
}
function onDragLeave(event: DragEvent) {
  if (!draggingFiles(event)) return;
  dragDepth.value = Math.max(0, dragDepth.value - 1);
}
/** Attaches dropped files directly, as many as there are free slots. */
async function onDrop(event: DragEvent) {
  if (!draggingFiles(event)) return;
  event.preventDefault();
  dragDepth.value = 0;
  pageDragDepth.value = 0;
  if (!canAttach.value) return;
  const files = [...(event.dataTransfer?.files ?? [])];
  const errors: string[] = [];
  let added = [...attachments.value];
  for (const file of files) {
    if (added.length >= maxAttachments) {
      errors.push(`${file.name}: 添付は${maxAttachments}件までです`);
      continue;
    }
    const result = await readAttachment(file);
    if (typeof result === "string") {
      errors.push(result);
      continue;
    }
    added = [
      ...added,
      {
        ...result,
        name: uniqueAttachmentName(
          result.name,
          added.map((attachment) => attachment.name),
        ),
      },
    ];
  }
  attachments.value = added;
  dropErrors.value = errors;
}
/**
 * A file dropped anywhere else in the form is ignored, so the browser does
 * not open it and leave the page.
 */
function ignoreFileDrop(event: DragEvent) {
  if (event.defaultPrevented || !draggingFiles(event)) return;
  event.preventDefault();
  if (event.dataTransfer) event.dataTransfer.dropEffect = "none";
}
function removeAttachment(name: string) {
  attachments.value = attachments.value.filter((a) => a.name !== name);
}
function confirm() {
  // The form shows progress and errors, so the dialog closes right away.
  confirming.value = false;
  emit("submit");
}
</script>
<template>
  <section class="answer-section">
    <form
      class="answer-form"
      @submit.prevent="canSubmit && (confirming = true)"
      @dragover="ignoreFileDrop"
      @drop="ignoreFileDrop"
    >
      <div class="answer-heading-row">
        <h2>回答</h2>
        <div class="answer-heading-tools">
          <p v-if="saveError" role="alert">{{ saveError }}</p>
          <span v-else-if="savedAt" class="draft-save-status">
            <button
              class="draft-save-icon"
              type="button"
              :aria-label="`このブラウザに保存済み ${new Date(savedAt).toLocaleTimeString('ja-JP')}`"
            >
              <svg
                aria-hidden="true"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <!-- Floppy disk with a check mark: saved in this browser. -->
                <path
                  d="M4 5.5A1.5 1.5 0 0 1 5.5 4h11l3.5 3.5v11a1.5 1.5 0 0 1-1.5 1.5h-13A1.5 1.5 0 0 1 4 18.5z"
                />
                <path d="M8 4v4.5h7V4" />
                <path d="M9 14.5l2 2 4-4" />
              </svg>
            </button>
            <span class="draft-save-tooltip" role="tooltip">
              このブラウザに保存済み
              <time :datetime="savedAt">{{
                new Date(savedAt).toLocaleTimeString("ja-JP")
              }}</time>
            </span>
          </span>
          <button
            type="button"
            role="switch"
            class="answer-preview-switch"
            :class="{ 'is-on': preview }"
            :aria-checked="preview"
            aria-label="プレビュー"
            :title="preview ? '編集に戻る' : 'プレビューを表示'"
            @click="preview = !preview"
          >
            <span
              class="answer-preview-option"
              :class="{ 'is-active': !preview }"
              ><UiIcon name="pencil"
            /></span>
            <span
              class="answer-preview-option"
              :class="{ 'is-active': preview }"
              ><UiIcon name="eye"
            /></span>
          </button>
        </div>
      </div>
      <label class="visually-hidden" for="answer">回答</label>
      <div v-show="!preview" class="answer-editor">
        <textarea
          id="answer"
          ref="editor"
          v-model="body"
          rows="10"
          required
          :disabled="!open || sending"
          :aria-describedby="'answer-length'"
          placeholder="原因 / 実施した変更 / 復旧確認の結果を記入してください"
          @keydown="onKeydown"
          @blur="onBlur"
        />
        <TabHint :visible="hintVisible" />
      </div>
      <div
        v-if="preview"
        class="answer-preview"
        aria-label="回答のプレビュー"
        tabindex="0"
        :style="{ height: previewHeight }"
      >
        <MarkdownContent
          v-if="body.trim()"
          :source="body"
          block-external-images
        />
        <p v-else class="answer-preview-empty">回答が入力されていません</p>
      </div>
      <p
        id="answer-length"
        class="answer-length"
        :class="{ 'is-over': remaining < 0 }"
        aria-live="polite"
      >
        {{ length.toLocaleString() }} / {{ maxLength.toLocaleString() }} 文字
        <span class="visually-hidden">（改行を除く）</span>
      </p>
      <div
        v-if="allowAttachments"
        class="answer-attachments"
        :class="{
          'is-drop-ready': pageDragDepth > 0 || dragDepth > 0,
          'is-drop-over': dragDepth > 0,
        }"
        @dragenter="onDragEnter"
        @dragover="onDragOver"
        @dragleave="onDragLeave"
        @drop="onDrop"
      >
        <button
          type="button"
          class="answer-attach-button"
          title="ログファイルを作成"
          :disabled="!canAttach || composing"
          @click="composing = true"
        >
          <span v-if="pageDragDepth > 0 || dragDepth > 0">{{
            !canAttach
              ? "これ以上添付できません"
              : dragDepth > 0
                ? "離すと添付します（.txt / .log）"
                : "ここにドロップしてログを添付"
          }}</span>
          <span v-else
            >ログの添付（ここにドロップも可、.txt / .log、各1MB、{{
              maxAttachments
            }}件まで）</span
          >
          <UiIcon name="file-plus" />
        </button>
        <ul v-if="dropErrors.length" role="alert" class="answer-drop-errors">
          <li v-for="message in dropErrors" :key="message">{{ message }}</li>
        </ul>
        <ul v-if="attachments.length" class="answer-attachment-list">
          <li v-for="attachment in attachments" :key="attachment.name">
            <details>
              <summary>
                <span>{{ attachment.name }}</span>
                <small>{{ formatBytes(attachment.size) }}</small>
              </summary>
              <!-- Shown as text only; Vue escapes the content. -->
              <pre class="answer-attachment-content">{{
                attachment.content
              }}</pre>
            </details>
            <button
              type="button"
              class="icon-button answer-attachment-remove"
              :aria-label="`${attachment.name} を外す`"
              title="添付を外す"
              :disabled="sending"
              @click="removeAttachment(attachment.name)"
            >
              <UiIcon name="close" />
            </button>
          </li>
        </ul>
        <LogComposer
          v-if="composing && canAttach"
          :suggested-name="suggestedName"
          :taken-names="attachments.map((attachment) => attachment.name)"
          @attach="addAttachment"
          @cancel="composing = false"
        />
      </div>
      <p v-if="submitted" role="status">回答を提出しました</p>
      <p v-if="cooldownSeconds" role="status">
        再回答可能まで {{ Math.floor(cooldownSeconds / 60) }}分{{
          cooldownSeconds % 60
        }}秒
      </p>
      <div class="answer-actions">
        <button class="button-primary button-submit" :disabled="!canSubmit">
          {{ sending ? "提出中…" : "回答を提出" }}
        </button>
      </div>
    </form>
    <div v-if="actionError" role="alert" class="state-message error-message">
      {{
        actionError instanceof Error ? actionError.message : String(actionError)
      }}
    </div>
    <AnswerConfirmDialog
      v-model:open="confirming"
      :body="body"
      :attachments="attachments"
      :length="length"
      :max-length="maxLength"
      :sending="sending"
      @confirm="confirm"
    />
  </section>
</template>
<style scoped>
.answer-heading-tools {
  display: flex;
  align-items: center;
  gap: 12px;
}
.answer-preview-switch {
  position: relative;
  display: inline-flex;
  padding: 2px;
  border: 1px solid var(--ink);
  border-radius: 0;
  background: transparent;
  cursor: pointer;
}
/* The knob slides under the active icon. */
.answer-preview-switch::before {
  content: "";
  position: absolute;
  top: 2px;
  left: 2px;
  width: 30px;
  height: 26px;
  border-radius: 0;
  background: var(--ink);
  transition: transform 0.15s;
}
.answer-preview-switch.is-on::before {
  transform: translateX(30px);
}
.answer-preview-switch:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}
.answer-preview-option {
  position: relative;
  width: 30px;
  height: 26px;
  display: grid;
  place-items: center;
  color: var(--muted);
  transition: color 0.15s;
}
.answer-preview-option.is-active {
  color: var(--paper);
}
.answer-preview-option :deep(.ui-icon) {
  width: 16px;
  height: 16px;
}
.answer-preview {
  min-height: 250px;
  max-height: 70vh;
  padding: 20px;
  overflow: auto;
  border: 1px solid var(--ink);
  overflow-wrap: anywhere;
}
.answer-editor {
  position: relative;
}
.answer-editor textarea {
  display: block;
}
.answer-preview :deep(.markdown > :first-child) {
  margin-top: 0;
}
.answer-preview-empty {
  color: var(--muted);
}
.answer-length {
  margin: 8px 0 0;
  font-size: 13px;
  text-align: right;
  font-variant-numeric: tabular-nums;
}
.answer-length.is-over {
  color: var(--danger, #c0392b);
  font-weight: 700;
}
/* While files are dragged over the page, a dashed frame marks the
   attachment area as the drop target; over it the frame turns solid and
   thicker. The label changes too, so neither state relies on color. */
.answer-attachments.is-drop-ready::after {
  content: "";
  position: absolute;
  inset: -4px -9px;
  border: 1px dashed var(--ink);
  pointer-events: none;
}
.answer-attachments.is-drop-over::after {
  border: 2px solid var(--ink);
}
.answer-attachments.is-drop-ready .answer-attach-button {
  font-weight: 700;
}
.answer-attachments {
  position: relative;
  margin-top: 14px;
}
.answer-drop-errors {
  margin: 8px 0 0;
  padding: 0;
  list-style: none;
  color: var(--danger, #c0392b);
  font-size: 12px;
}
.answer-attach-button {
  /* Text and icon sit on the column edges like the textarea and submit
     button; the hover border reaches 9px outside them. */
  width: calc(100% + 18px);
  margin: 0 -9px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 6px 8px;
  border: 1px solid transparent;
  background: transparent;
  color: var(--ink);
  cursor: pointer;
  font: inherit;
  font-size: 12px;
  text-align: left;
}
.answer-attach-button:hover:not(:disabled),
.answer-attach-button:focus-visible {
  border-color: var(--ink);
  outline: none;
}
.answer-attach-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.icon-button {
  width: 32px;
  height: 32px;
  display: inline-grid;
  place-items: center;
  flex: none;
  padding: 0;
  border: 1px solid transparent;
  background: transparent;
  color: var(--ink);
  cursor: pointer;
}
.icon-button:hover:not(:disabled),
.icon-button:focus-visible {
  border-color: var(--ink);
  outline: none;
}
.icon-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.answer-editor textarea {
  font-family: "IBM Plex Mono", "Noto Sans JP", monospace;
}
.answer-attachment-list {
  display: grid;
  gap: 4px;
  margin: 8px 0 0;
  padding: 0;
  list-style: none;
  font-size: 13px;
}
.answer-attachment-list li {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}
.answer-attachment-list details {
  flex: 1;
  min-width: 0;
}
.answer-attachment-list summary {
  min-height: 28px;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}
/* Same states as the top bar: hover is accent with an underline. */
.answer-attachment-list summary:hover span,
.answer-attachment-list summary:focus-visible span {
  color: var(--accent);
  text-decoration: underline;
  text-underline-offset: 3px;
}
.answer-attachment-list span {
  font-family: var(--mono);
  overflow-wrap: anywhere;
}
.answer-attachment-list small {
  color: var(--muted);
}
.answer-attachment-content {
  max-height: 240px;
  margin: 4px 0 8px;
  padding: 4px 0 4px 12px;
  overflow: auto;
  border-left: 2px solid var(--line-light);
  font-family: var(--mono);
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}
.answer-attachment-remove {
  width: 28px;
  height: 28px;
  /* Puts the 18px icon on the column edge, under the attach icon. */
  margin-right: -5px;
  margin-left: auto;
}
</style>
