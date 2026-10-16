<script setup lang="ts">
const railOpen = defineModel<boolean>("railOpen", { required: true });
defineProps<{ resetDisabled: boolean }>();
defineEmits<{ reset: [] }>();
</script>
<template>
  <div class="actions detail-toolbar">
    <button
      class="problem-rail-toggle"
      :aria-expanded="railOpen"
      :aria-label="railOpen ? '問題リンクを隠す' : '問題リンクを表示'"
      aria-controls="problem-rail"
      @click="railOpen = !railOpen"
    >
      <span class="problem-rail-toggle-icon" aria-hidden="true" />
      <span class="visually-hidden">{{
        railOpen ? "問題リンクを隠す" : "問題リンクを表示"
      }}</span>
    </button>
    <button
      class="button-secondary button-reset detail-reset"
      type="button"
      :disabled="resetDisabled"
      @click="$emit('reset')"
    >
      環境をリセット
    </button>
  </div>
</template>
