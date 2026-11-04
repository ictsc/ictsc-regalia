<script setup lang="ts">
defineProps<{ visible: boolean }>();
</script>
<template>
  <Transition name="tab-hint">
    <p v-if="visible" class="tab-hint" role="status">
      ヒント: Esc で入力欄からフォーカスが外れます
    </p>
  </Transition>
</template>
<style scoped>
.tab-hint {
  position: absolute;
  z-index: 5;
  right: 8px;
  bottom: 8px;
  margin: 0;
  padding: 6px 10px;
  background: var(--ink);
  color: var(--paper);
  font-size: 12px;
  line-height: 1.4;
  pointer-events: none;
}
.tab-hint-enter-active,
.tab-hint-leave-active {
  transition: opacity 0.2s;
}
.tab-hint-enter-from,
.tab-hint-leave-to {
  opacity: 0;
}
</style>
