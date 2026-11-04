<script setup lang="ts">
import { computed, ref, watch } from "vue";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import { answerLength, type AttachmentFile } from "@ictsc/ui/attachments";
import TabHint from "@ictsc/ui/components/TabHint.vue";
import UiIcon from "@ictsc/ui/components/UiIcon.vue";
import { useTabIndent } from "@ictsc/ui/textarea";
import AnswerAttachments from "./AnswerAttachments.vue";
import AnswerConfirmDialog from "./AnswerConfirmDialog.vue";

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
// The preview keeps the editor's size and scrolls inside it.
watch(
  preview,
  (showing) => {
    if (showing && editor.value?.offsetHeight)
      previewHeight.value = `${editor.value.offsetHeight}px`;
  },
  { flush: "pre" },
);

/**
 * Files dropped outside the attachment area are ignored, so the browser does
 * not open them and leave the page.
 */
function ignoreFileDrop(event: DragEvent) {
  if (event.defaultPrevented || !event.dataTransfer?.types.includes("Files"))
    return;
  event.preventDefault();
  event.dataTransfer.dropEffect = "none";
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
          diagrams
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
      <AnswerAttachments
        v-if="allowAttachments"
        v-model="attachments"
        :open="open"
        :sending="sending"
      />
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
  font-family: "IBM Plex Mono", "Noto Sans JP", monospace;
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
</style>
