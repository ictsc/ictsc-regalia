<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import {
  appendAttachmentFiles,
  formatBytes,
  maxAttachments,
  type AttachmentFile,
} from "@ictsc/ui/attachments";
import UiIcon from "@ictsc/ui/components/UiIcon.vue";
import LogComposer from "./LogComposer.vue";

/**
 * The log files attached to an answer: drop, write or paste, and remove.
 * While mounted it listens for drags on the whole window, so it can mark
 * itself as the drop target as soon as a file is dragged over the page.
 */
const attachments = defineModel<AttachmentFile[]>({ required: true });
const props = defineProps<{
  /** Answers are being accepted now. */
  open: boolean;
  sending: boolean;
}>();

const composing = ref(false);
/** Files are dragged somewhere over the page; counts nested dragenters. */
const pageDragDepth = ref(0);
/** Files are dragged over this area itself. */
const dragDepth = ref(0);
const dropErrors = ref<string[]>([]);
const canAttach = computed(
  () =>
    props.open && !props.sending && attachments.value.length < maxAttachments,
);
const dragging = computed(() => pageDragDepth.value > 0 || dragDepth.value > 0);
const takenNames = computed(() =>
  attachments.value.map((attachment) => attachment.name),
);
const suggestedName = computed(() => {
  for (let index = 1; ; index++) {
    const name = `log-${index}.txt`;
    if (!takenNames.value.includes(name)) return name;
  }
});

function draggingFiles(event: DragEvent) {
  return event.dataTransfer?.types.includes("Files") ?? false;
}

// While files are dragged anywhere on the page, this area shows itself as
// the drop target.
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
const pageListeners = {
  dragenter: onPageDragEnter,
  dragleave: onPageDragLeave,
  drop: onPageDragEnd,
  dragend: onPageDragEnd,
} as const;
onMounted(() => {
  for (const [type, listener] of Object.entries(pageListeners))
    window.addEventListener(type, listener as EventListener);
});
onBeforeUnmount(() => {
  for (const [type, listener] of Object.entries(pageListeners))
    window.removeEventListener(type, listener as EventListener);
});

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
  onPageDragEnd();
  if (!canAttach.value) return;
  const result = await appendAttachmentFiles(
    attachments.value,
    event.dataTransfer?.files ?? [],
  );
  attachments.value = result.attachments;
  dropErrors.value = result.errors;
}

function addAttachment(attachment: AttachmentFile) {
  attachments.value = [...attachments.value, attachment];
  composing.value = false;
}
function removeAttachment(name: string) {
  attachments.value = attachments.value.filter((a) => a.name !== name);
}
</script>
<template>
  <div
    class="answer-attachments"
    :class="{ 'is-drop-ready': dragging, 'is-drop-over': dragDepth > 0 }"
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
      <span v-if="dragging">{{
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
          <pre class="answer-attachment-content">{{ attachment.content }}</pre>
        </details>
        <button
          type="button"
          class="answer-attachment-remove"
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
      :taken-names="takenNames"
      @attach="addAttachment"
      @cancel="composing = false"
    />
  </div>
</template>
<style scoped>
.answer-attachments {
  position: relative;
  margin-top: 14px;
}
/* While files are dragged over the page, a dashed frame marks this area as
   the drop target; over it the frame turns solid and thicker. The label
   changes too, so neither state relies on color. */
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
.answer-attach-button:focus-visible,
.answer-attachment-remove:hover:not(:disabled),
.answer-attachment-remove:focus-visible {
  border-color: var(--ink);
  outline: none;
}
.answer-attach-button:disabled,
.answer-attachment-remove:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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
  display: inline-grid;
  place-items: center;
  flex: none;
  /* Puts the 18px icon on the column edge, under the attach icon. */
  margin-right: -5px;
  margin-left: auto;
  padding: 0;
  border: 1px solid transparent;
  background: transparent;
  color: var(--ink);
  cursor: pointer;
}
</style>
