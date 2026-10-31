<script setup lang="ts">
import { ref, watchEffect } from "vue";
const props = defineProps<{ pending?: boolean; error?: unknown }>();
defineEmits<{ retry: [] }>();
// After the first successful load, background refetches (realtime updates)
// keep showing the current content instead of flashing the loading state.
const loaded = ref(false);
watchEffect(() => {
  if (!props.pending && !props.error) loaded.value = true;
});
</script>
<template>
  <div v-if="error" role="alert" class="state-message error-message">
    <p>{{ error instanceof Error ? error.message : String(error) }}</p>
    <button class="button-secondary" @click="$emit('retry')">再試行</button>
  </div>
  <p v-else-if="pending && !loaded" role="status" class="state-message">
    読み込み中…
  </p>
  <slot v-else />
</template>
