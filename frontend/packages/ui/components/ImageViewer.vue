<script setup lang="ts">
import { nextTick, ref } from "vue";
import UiIcon from "./UiIcon.vue";

/** Shows one image over the page; click it, the backdrop or Esc to close. */
const dialog = ref<HTMLDialogElement | null>(null);
const image = ref<{ src: string; alt: string }>();

async function open(src: string, alt: string) {
  image.value = { src, alt };
  await nextTick();
  dialog.value?.showModal();
}
function close() {
  dialog.value?.close();
}
defineExpose({ open });
</script>
<template>
  <dialog
    ref="dialog"
    class="image-viewer"
    :aria-label="image?.alt || '画像の拡大表示'"
    @click.self="close"
    @close="image = undefined"
  >
    <template v-if="image">
      <img :src="image.src" :alt="image.alt" @click="close" />
      <button
        type="button"
        class="image-viewer-close"
        aria-label="閉じる"
        title="閉じる"
        @click="close"
      >
        <UiIcon name="close" />
      </button>
    </template>
  </dialog>
</template>
