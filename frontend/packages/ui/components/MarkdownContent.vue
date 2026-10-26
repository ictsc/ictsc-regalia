<script setup lang="ts">
import { ref, watch } from "vue";
import { renderMarkdown, type ImageUrl } from "../markdown";
const props = defineProps<{
  source: string;
  /** Overrides how /images/<file> is served, e.g. for a past content commit. */
  imageUrl?: ImageUrl;
}>();
const html = ref("");
const failed = ref(false);
watch(
  () => [props.source, props.imageUrl] as const,
  async ([source, imageUrl], _, onCleanup) => {
    let active = true;
    onCleanup(() => {
      active = false;
    });
    try {
      const result = await renderMarkdown(source, { imageUrl });
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
  <!-- HTML is sanitized before math and syntax highlighting. -->
  <div v-else class="markdown" v-html="html" />
</template>
