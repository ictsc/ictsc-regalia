<script setup lang="ts">
import {
  downloadAttachment,
  formatBytes,
  type AttachmentFile,
} from "../attachments";
import UiIcon from "./UiIcon.vue";

defineProps<{ attachments: AttachmentFile[] }>();
</script>
<template>
  <ul class="attachment-list">
    <li v-for="attachment in attachments" :key="attachment.name">
      <details>
        <summary>
          <span class="attachment-name">{{ attachment.name }}</span>
          <small>{{ formatBytes(attachment.size) }}</small>
        </summary>
        <!-- Logs are shown as text only; Vue escapes the content. -->
        <pre class="attachment-content">{{ attachment.content }}</pre>
      </details>
      <button
        type="button"
        class="attachment-download"
        :aria-label="`${attachment.name} をダウンロード`"
        title="ダウンロード"
        @click="downloadAttachment(attachment)"
      >
        <UiIcon name="download" />
      </button>
    </li>
  </ul>
</template>
<style scoped>
.attachment-list {
  display: grid;
  gap: 8px;
  margin: 0;
  padding: 0;
  list-style: none;
}
.attachment-list li {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: start;
  gap: 8px;
  padding: 8px 10px;
  border: 1px solid var(--line-light);
}
.attachment-list summary {
  cursor: pointer;
}
.attachment-name {
  margin-right: 8px;
  font-family: var(--mono);
  overflow-wrap: anywhere;
}
.attachment-content {
  max-height: 360px;
  margin: 8px 0 0;
  padding: 10px;
  overflow: auto;
  background: var(--surface-hover);
  font-size: 12px;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}
.attachment-download {
  width: 28px;
  height: 28px;
  display: inline-grid;
  place-items: center;
  padding: 0;
  border: 1px solid transparent;
  background: transparent;
  color: var(--ink);
  cursor: pointer;
}
.attachment-download:hover,
.attachment-download:focus-visible {
  border-color: var(--ink);
  outline: none;
}
</style>
