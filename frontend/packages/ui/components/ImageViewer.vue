<script setup lang="ts">
import { nextTick, ref } from "vue";
import UiIcon from "./UiIcon.vue";

/**
 * Shows one image or diagram over the page; click it, the backdrop or Esc to
 * close.
 */
const dialog = ref<HTMLDialogElement | null>(null);
const diagramBox = ref<HTMLElement | null>(null);
const image = ref<{ src: string; alt: string; diagram?: SVGSVGElement }>();

async function open(src: string, alt: string) {
  image.value = { src, alt };
  await nextTick();
  dialog.value?.showModal();
}
/** Shows a copy of a drawn diagram, keeping its theme colors. */
async function openDiagram(svg: SVGSVGElement, label: string) {
  image.value = { src: "", alt: label, diagram: svg };
  await nextTick();
  diagramBox.value?.replaceChildren(svg.cloneNode(true));
  dialog.value?.showModal();
}
function close() {
  dialog.value?.close();
}
defineExpose({ open, openDiagram });
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
      <div
        v-if="image.diagram"
        ref="diagramBox"
        class="image-viewer-diagram"
        role="img"
        :aria-label="image.alt"
        @click="close"
      />
      <img v-else :src="image.src" :alt="image.alt" @click="close" />
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
