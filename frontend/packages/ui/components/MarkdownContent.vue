<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { addCopyButtons, copyCode } from "../code-copy";
import { renderMarkdown, type ImageUrl } from "../markdown";
import { renderDiagrams } from "../mermaid";
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
  /** Draws ```mermaid blocks as diagrams (problems and answers). */
  diagrams?: boolean;
}>();
// Two roots (the content and the image viewer): attributes go on the content.
defineOptions({ inheritAttrs: false });
const html = ref("");
const failed = ref(false);
const root = ref<HTMLElement | null>(null);
const viewer = ref<InstanceType<typeof ImageViewer> | null>(null);
const copyTimers = new Set<ReturnType<typeof setTimeout>>();

watch(
  () =>
    [
      props.source,
      props.imageUrl,
      props.blockExternalImages,
      props.diagrams,
    ] as const,
  async ([source, imageUrl, blockExternalImages, diagrams], _, onCleanup) => {
    let active = true;
    onCleanup(() => {
      active = false;
    });
    try {
      const result = await renderMarkdown(source, {
        imageUrl,
        blockExternalImages,
        mermaid: diagrams,
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
function drawDiagrams() {
  // A diagram that cannot be drawn keeps showing its source.
  if (props.diagrams && root.value)
    renderDiagrams(root.value).catch(() => undefined);
}
// v-html replaces the code blocks on every render, so add the buttons and
// draw the diagrams again.
watch(
  html,
  () =>
    nextTick(() => {
      if (!root.value) return;
      addCopyButtons(root.value);
      drawDiagrams();
    }),
  { flush: "post" },
);
// Diagrams are drawn in the theme's colors, so redraw them when it changes.
let themeObserver: MutationObserver | undefined;
onMounted(() => {
  themeObserver = new MutationObserver(drawDiagrams);
  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["data-theme"],
  });
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
  const diagram = target.closest<SVGSVGElement>(".mermaid-output svg");
  if (diagram) {
    await viewer.value?.openDiagram(diagram, "図の拡大表示");
    return;
  }
  // Images inside links keep following the link.
  if (target instanceof HTMLImageElement && !target.closest("a"))
    await viewer.value?.open(target.currentSrc || target.src, target.alt);
}

onBeforeUnmount(() => {
  themeObserver?.disconnect();
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
