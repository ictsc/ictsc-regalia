<script setup lang="ts">
defineProps<{ pending?: boolean; error?: unknown }>();
defineEmits<{ retry: [] }>();
</script>
<template>
  <div v-if="error" role="alert" class="state-message error-message">
    <p>{{ error instanceof Error ? error.message : String(error) }}</p>
    <button class="button-secondary" @click="$emit('retry')">再試行</button>
  </div>
  <p v-else-if="pending" role="status" class="state-message">読み込み中…</p>
  <slot v-else />
</template>
