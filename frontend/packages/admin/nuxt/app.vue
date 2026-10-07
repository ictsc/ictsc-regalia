<script setup lang="ts">
import { useAdminConfirmation } from "~/features/confirmation";

const { current: confirmation, settle: closeConfirmation } =
  useAdminConfirmation();
const confirmationDialog = ref<HTMLDialogElement | null>(null);
watch(confirmation, async (request) => {
  await nextTick();
  if (request && !confirmationDialog.value?.open)
    confirmationDialog.value?.showModal();
  else if (!request && confirmationDialog.value?.open)
    confirmationDialog.value.close();
});
function onConfirmationBackdrop(event: MouseEvent) {
  if (event.target === confirmationDialog.value) closeConfirmation(false);
}

onBeforeUnmount(() => closeConfirmation(false));
</script>

<template>
  <NuxtLayout><NuxtPage /></NuxtLayout>
  <dialog
    ref="confirmationDialog"
    class="confirmation-dialog"
    aria-labelledby="admin-confirm-title"
    aria-describedby="admin-confirm-message"
    @cancel.prevent="closeConfirmation(false)"
    @close="closeConfirmation(false)"
    @click="onConfirmationBackdrop"
  >
    <template v-if="confirmation">
      <p class="confirmation-dialog-eyebrow">
        {{ confirmation.destructive ? "取り消せない操作" : "操作の確認" }}
      </p>
      <h2 id="admin-confirm-title">{{ confirmation.title }}</h2>
      <p id="admin-confirm-message" class="confirmation-dialog-message">
        {{ confirmation.message }}
      </p>
      <div class="confirmation-dialog-actions">
        <button
          type="button"
          class="button-secondary"
          autofocus
          @click="closeConfirmation(false)"
        >
          キャンセル
        </button>
        <button
          type="button"
          class="button-primary"
          @click="closeConfirmation(true)"
        >
          {{ confirmation.confirmLabel ?? "実行する" }}
        </button>
      </div>
    </template>
  </dialog>
</template>
