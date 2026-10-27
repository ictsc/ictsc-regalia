<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from "vue";
import { renderMarkdown, type ImageUrl } from "../markdown";
import UiIcon from "./UiIcon.vue";
const props = defineProps<{
  source: string;
  /** Overrides how /images/<file> is served, e.g. for a past content commit. */
  imageUrl?: ImageUrl;
}>();
// Two roots (the content and the image viewer): attributes go on the content.
defineOptions({ inheritAttrs: false });
const html = ref("");
const failed = ref(false);
const root = ref<HTMLElement | null>(null);
const viewer = ref<HTMLDialogElement | null>(null);
const viewing = ref<{ src: string; alt: string }>();
const copyResetTimers = new Set<ReturnType<typeof setTimeout>>();

// Static markup for the buttons added inside v-html, drawn like UiIcon.
const iconSvg = (paths: string) =>
  `<svg class="ui-icon" aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">${paths}</svg>`;
const copyIcons = {
  idle: iconSvg(
    '<rect x="8.5" y="8.5" width="11" height="11" rx="1.5"/><path d="M15.5 8.5V5.5A1.5 1.5 0 0 0 14 4H5.5A1.5 1.5 0 0 0 4 5.5V14a1.5 1.5 0 0 0 1.5 1.5h3"/>',
  ),
  done: iconSvg('<path d="M5 12.5l4.5 4.5L19 7.5"/>'),
  failed: iconSvg('<path d="M6 6l12 12M18 6L6 18"/>'),
};
function setCopyState(
  button: HTMLButtonElement,
  state: keyof typeof copyIcons,
) {
  const label = {
    idle: "コードをコピー",
    done: "コピーしました",
    failed: "コピーできませんでした",
  }[state];
  button.innerHTML = copyIcons[state];
  button.setAttribute("aria-label", label);
  button.title = label;
}

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

/**
 * Adds a copy button to each code block of the rendered (sanitized) HTML.
 * The block is wrapped so the button stays put while the code scrolls.
 */
function addCopyButtons() {
  for (const pre of root.value?.querySelectorAll("pre") ?? []) {
    if (!pre.querySelector(":scope > code")) continue;
    if (pre.parentElement?.classList.contains("code-block")) continue;
    const wrapper = document.createElement("div");
    wrapper.className = "code-block";
    pre.replaceWith(wrapper);
    wrapper.append(pre);
    const button = document.createElement("button");
    button.type = "button";
    button.className = "code-copy";
    setCopyState(button, "idle");
    wrapper.append(button);
  }
}
watch(html, () => nextTick(addCopyButtons), { flush: "post" });

async function copyCode(button: HTMLButtonElement) {
  const code = button.parentElement?.querySelector("pre > code");
  if (!code) return;
  try {
    await navigator.clipboard.writeText(code.textContent ?? "");
    setCopyState(button, "done");
  } catch {
    setCopyState(button, "failed");
  }
  const timer = setTimeout(() => {
    setCopyState(button, "idle");
    copyResetTimers.delete(timer);
  }, 2000);
  copyResetTimers.add(timer);
}

/** Copy buttons and images are handled here, since they live in v-html. */
function onClick(event: MouseEvent) {
  const target = event.target as HTMLElement;
  const button = target.closest<HTMLButtonElement>(".code-copy");
  if (button) {
    void copyCode(button);
    return;
  }
  // Images inside links keep following the link.
  if (target instanceof HTMLImageElement && !target.closest("a")) {
    viewing.value = { src: target.currentSrc || target.src, alt: target.alt };
    void nextTick(() => viewer.value?.showModal());
  }
}
function closeViewer() {
  viewer.value?.close();
}
function onViewerClose() {
  viewing.value = undefined;
}

onBeforeUnmount(() => {
  for (const timer of copyResetTimers) clearTimeout(timer);
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
  <dialog
    ref="viewer"
    class="image-viewer"
    :aria-label="viewing?.alt || '画像の拡大表示'"
    @click.self="closeViewer"
    @close="onViewerClose"
  >
    <template v-if="viewing">
      <img :src="viewing.src" :alt="viewing.alt" @click="closeViewer" />
      <button
        type="button"
        class="image-viewer-close"
        aria-label="閉じる"
        title="閉じる"
        @click="closeViewer"
      >
        <UiIcon name="close" />
      </button>
    </template>
  </dialog>
</template>
