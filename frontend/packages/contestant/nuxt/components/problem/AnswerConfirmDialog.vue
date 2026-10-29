<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import AttachmentList from "@ictsc/ui/components/AttachmentList.vue";
import type { AttachmentFile } from "@ictsc/ui/attachments";

const open = defineModel<boolean>("open", { default: false });
defineProps<{
  body: string;
  attachments: readonly AttachmentFile[];
  length: number;
  maxLength: number;
  sending: boolean;
}>();
defineEmits<{ confirm: [] }>();

const dialog = ref<HTMLDialogElement | null>(null);
const preview = ref<HTMLElement | null>(null);
function sync(value: boolean) {
  const element = dialog.value;
  if (!element) return;
  // jsdom lacks showModal; real browsers always have it.
  if (value && !element.open) {
    element.showModal?.();
    // Start at the top: the heading takes focus and the preview is unscrolled.
    element.scrollTop = 0;
    if (preview.value) preview.value.scrollTop = 0;
  } else if (!value && element.open) element.close();
}
watch(open, sync, { flush: "post" });
onMounted(() => sync(open.value));
function onBackdrop(event: MouseEvent) {
  if (event.target === dialog.value) open.value = false;
}
</script>
<template>
  <dialog
    ref="dialog"
    class="answer-confirm-dialog confirmation-dialog"
    aria-labelledby="answer-confirm-title"
    @click="onBackdrop"
    @close="open = false"
  >
    <p class="confirmation-dialog-eyebrow">提出内容の確認</p>
    <h2 id="answer-confirm-title" tabindex="-1" autofocus>
      この内容で提出しますか？
    </h2>
    <p class="confirmation-dialog-note">
      {{ length.toLocaleString() }} / {{ maxLength.toLocaleString() }} 文字
    </p>
    <div ref="preview" class="answer-confirm-preview">
      <MarkdownContent :source="body" block-external-images diagrams />
    </div>
    <template v-if="attachments.length">
      <h3>添付するログ（クリックで内容を確認）</h3>
      <AttachmentList :attachments="[...attachments]" />
    </template>
    <div class="confirmation-dialog-actions">
      <button type="button" class="button-secondary" @click="open = false">
        キャンセル
      </button>
      <button
        type="button"
        class="button-primary"
        :disabled="sending"
        @click="$emit('confirm')"
      >
        {{ sending ? "提出中…" : "提出する" }}
      </button>
    </div>
  </dialog>
</template>
<style scoped>
.answer-confirm-dialog {
  width: min(760px, calc(100vw - 32px));
}
.answer-confirm-preview {
  max-height: min(55vh, 520px);
  margin: 12px 0;
  padding: 16px;
  overflow: auto;
  border: 1px solid var(--line-light);
}
.answer-confirm-dialog h2:focus {
  outline: none;
}
.answer-confirm-dialog h3 {
  margin: 12px 0 6px;
  font-size: 14px;
}
.answer-confirm-preview :deep(.markdown > :first-child) {
  margin-top: 0;
}
</style>
