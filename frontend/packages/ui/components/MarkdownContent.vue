<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from "vue";
import { addCopyButtons, copyCode } from "../code-copy";
import { renderMarkdown, type ImageUrl } from "../markdown";
import ImageViewer from "./ImageViewer.vue";
const props = defineProps<{
  source: string;
  /** Overrides how /images/<file> is served, e.g. for a past content commit. */
  imageUrl?: ImageUrl;
  /**
   * Written by contestants (answers, profiles): images from other hosts are
   * shown as their URL instead of being loaded in the reader's browser.
   */
  blockExternalImages?: boolean;
}>();
// Two roots (the content and the image viewer): attributes go on the content.
defineOptions({ inheritAttrs: false });
const html = ref("");
const failed = ref(false);
const root = ref<HTMLElement | null>(null);
const viewer = ref<InstanceType<typeof ImageViewer> | null>(null);
const copyTimers = new Set<ReturnType<typeof setTimeout>>();

watch(
  () => [props.source, props.imageUrl, props.blockExternalImages] as const,
  async ([source, imageUrl, blockExternalImages], _, onCleanup) => {
    let active = true;
    onCleanup(() => {
      active = false;
    });
    try {
      const result = await renderMarkdown(source, {
        imageUrl,
        blockExternalImages,
      });
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
// v-html replaces the code blocks on every render, so add the buttons again.
watch(html, () => nextTick(() => root.value && addCopyButtons(root.value)), {
  flush: "post",
});

/** Copy buttons and images live in v-html, so clicks are handled here. */
async function onClick(event: MouseEvent) {
  const target = event.target as Element;
  const button = target.closest<HTMLButtonElement>(".code-copy");
  if (button) {
    const timer = await copyCode(button);
    if (timer) copyTimers.add(timer);
    return;
  }
  // Images inside links keep following the link.
  if (target instanceof HTMLImageElement && !target.closest("a"))
    await viewer.value?.open(target.currentSrc || target.src, target.alt);
}

onBeforeUnmount(() => {
  for (const timer of copyTimers) clearTimeout(timer);
});
</script>
<template>
  <pre v-if="failed" v-bind="$attrs">{{ source }}</pre>
  <!-- HTML is sanitized before math and syntax highlighting. -->
  <div
    v-else
    v-bind="$attrs"
    ref="root"
    class="markdown"
    @click="onClick"
    v-html="html"
  />
  <ImageViewer ref="viewer" />
</template>
