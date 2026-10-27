<script setup lang="ts">
import { ref } from "vue";
import { clampAnswerWidth } from "~/features/problem/answer-width";

/** Answer panel width in pixels; undefined keeps the layout default. */
const width = defineModel<number | undefined>("width");
const props = defineProps<{
  /** The content grid, whose width bounds the answer panel. */
  container: HTMLElement | null;
}>();

const step = 24;
const dragging = ref(false);
function currentWidth() {
  const answer = props.container?.querySelector<HTMLElement>(".answer-section");
  return width.value ?? answer?.offsetWidth ?? 0;
}
function setWidth(next: number) {
  width.value = clampAnswerWidth(next, props.container?.clientWidth ?? next);
}
function onPointerDown(event: PointerEvent) {
  if (event.button !== 0) return;
  event.preventDefault();
  const handle = event.currentTarget as HTMLElement;
  handle.setPointerCapture(event.pointerId);
  const startX = event.clientX;
  const startWidth = currentWidth();
  const move = (moveEvent: PointerEvent) =>
    // Dragging left widens the answer panel.
    setWidth(startWidth + startX - moveEvent.clientX);
  const stop = () => {
    handle.removeEventListener("pointermove", move);
    handle.removeEventListener("pointerup", stop);
    handle.removeEventListener("pointercancel", stop);
    document.body.classList.remove("is-resizing-answer");
    dragging.value = false;
  };
  handle.addEventListener("pointermove", move);
  handle.addEventListener("pointerup", stop);
  handle.addEventListener("pointercancel", stop);
  document.body.classList.add("is-resizing-answer");
  dragging.value = true;
}
function onKeydown(event: KeyboardEvent) {
  if (event.key === "ArrowLeft") setWidth(currentWidth() + step);
  else if (event.key === "ArrowRight") setWidth(currentWidth() - step);
  else if (event.key === "Home" || event.key === "Escape")
    width.value = undefined;
  else return;
  event.preventDefault();
}
</script>
<template>
  <div
    class="answer-resizer"
    :class="{ 'is-dragging': dragging }"
    role="separator"
    aria-orientation="vertical"
    aria-label="回答欄の幅（← で広げる、→ で狭める、ダブルクリックで元に戻す）"
    :aria-valuenow="width"
    tabindex="0"
    title="ドラッグで回答欄の幅を変更（ダブルクリックで元に戻す）"
    @pointerdown="onPointerDown"
    @keydown="onKeydown"
    @dblclick="width = undefined"
  />
</template>
<style scoped>
/* A 24px hit area centered on the panel edge; both columns keep padding
   there, so it never covers text. */
.answer-resizer {
  grid-column: 2;
  grid-row: 1 / span 4;
  justify-self: start;
  position: sticky;
  z-index: 3;
  top: 20px;
  width: 24px;
  height: calc(100vh - 40px);
  margin-left: -12px;
  cursor: col-resize;
  touch-action: none;
}
/* The line shows only while the pointer is in the hit area or dragging. */
.answer-resizer::after {
  content: "";
  position: absolute;
  top: 0;
  bottom: 0;
  left: 11px;
  width: 1px;
  background: transparent;
}
.answer-resizer:hover::after,
.answer-resizer:focus-visible::after,
.answer-resizer.is-dragging::after {
  left: 10px;
  width: 3px;
  background: var(--accent);
}
.answer-resizer:focus-visible {
  outline: none;
}
@media (max-width: 1000px) {
  .answer-resizer {
    display: none;
  }
}
</style>
