<script setup lang="ts">
import { nextTick, ref, watch } from "vue";
import type { AdminConfirmation } from "~/features/confirmation";

const props = defineProps<{
  /** The pending confirmation; null closes the dialog. */
  request: AdminConfirmation | null;
}>();
const emit = defineEmits<{ settle: [accepted: boolean] }>();

const dialog = ref<HTMLDialogElement | null>(null);
watch(
  () => props.request,
  async (request) => {
    await nextTick();
    if (request && !dialog.value?.open) dialog.value?.showModal();
    else if (!request && dialog.value?.open) dialog.value.close();
  },
  { immediate: true },
);
function onBackdrop(event: MouseEvent) {
  if (event.target === dialog.value) emit("settle", false);
}
</script>
<template>
  <dialog
    ref="dialog"
    class="confirmation-dialog"
    aria-labelledby="admin-confirm-title"
    aria-describedby="admin-confirm-message"
    @cancel.prevent="$emit('settle', false)"
    @close="$emit('settle', false)"
    @click="onBackdrop"
  >
    <template v-if="request">
      <p class="confirmation-dialog-eyebrow">
        {{ request.destructive ? "取り消せない操作" : "操作の確認" }}
      </p>
      <h2 id="admin-confirm-title">{{ request.title }}</h2>
      <p id="admin-confirm-message" class="confirmation-dialog-message">
        {{ request.message }}
      </p>
      <div class="confirmation-dialog-actions">
        <button
          type="button"
          class="button-secondary"
          autofocus
          @click="$emit('settle', false)"
        >
          キャンセル
        </button>
        <button
          type="button"
          class="button-primary"
          @click="$emit('settle', true)"
        >
          {{ request.confirmLabel ?? "実行する" }}
        </button>
      </div>
    </template>
  </dialog>
</template>
