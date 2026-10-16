<script setup lang="ts">
import { onMounted, ref, watch } from "vue";

const open = defineModel<boolean>("open", { default: false });
defineProps<{ penaltyThreshold?: number }>();
defineEmits<{ confirm: [] }>();

const dialog = ref<HTMLDialogElement | null>(null);
function sync(value: boolean) {
  const element = dialog.value;
  if (!element) return;
  if (value && !element.open) element.showModal();
  else if (!value && element.open) element.close();
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
    class="reset-dialog confirmation-dialog"
    aria-labelledby="reset-dialog-title"
    aria-describedby="reset-dialog-message"
    @click="onBackdrop"
    @close="open = false"
  >
    <p class="confirmation-dialog-eyebrow">操作の確認</p>
    <h2 id="reset-dialog-title">環境をリセットしますか？</h2>
    <p id="reset-dialog-message" class="confirmation-dialog-message">
      この問題の環境を再展開します。再展開ルールに応じて減点される場合があります。
    </p>
    <p v-if="penaltyThreshold != null" class="confirmation-dialog-note">
      減点なしの上限: {{ penaltyThreshold }}回
    </p>
    <div class="confirmation-dialog-actions">
      <button
        type="button"
        class="button-secondary"
        autofocus
        @click="open = false"
      >
        キャンセル
      </button>
      <button type="button" class="button-primary" @click="$emit('confirm')">
        環境をリセットする
      </button>
    </div>
  </dialog>
</template>
