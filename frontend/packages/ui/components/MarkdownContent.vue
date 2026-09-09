<script setup lang="ts">
import { ref, watch } from "vue";
import { renderMarkdown } from "../markdown";
const props = defineProps<{ source: string }>();
const html = ref("");
const failed = ref(false);
watch(
  () => props.source,
  async (source, _, onCleanup) => {
    let active = true;
    onCleanup(() => {
      active = false;
    });
    try {
      const result = await renderMarkdown(source);
      if (active) {
        html.value = result;
        failed.value = false;
      }
    } catch {
      if (active) failed.value = true;
    }
  },
  { immediate: true },
);
</script>
<template>
  <pre v-if="failed">{{ source }}</pre>
  <!-- HTML is sanitized before trusted math/highlight transforms. -->
  <div v-else class="markdown" v-html="html" />
</template>
