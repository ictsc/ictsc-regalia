<script setup lang="ts">
import { ref } from "vue";
import {
  readAttachment,
  textAttachment,
  type AttachmentFile,
} from "@ictsc/ui/attachments";
import TabHint from "@ictsc/ui/components/TabHint.vue";
import { useTabIndent } from "@ictsc/ui/textarea";

const props = defineProps<{
  /** Default file name, e.g. log-2.txt. */
  suggestedName: string;
  /** Names already attached; a new log may not reuse them. */
  takenNames: string[];
}>();
const emit = defineEmits<{
  attach: [attachment: AttachmentFile];
  cancel: [];
}>();

const name = ref(props.suggestedName);
const content = ref("");
const error = ref("");
const { onKeydown, onBlur, hintVisible } = useTabIndent();

async function upload(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  input.value = "";
  if (!file) return;
  const result = await readAttachment(file);
  if (typeof result === "string") {
    error.value = result;
    return;
  }
  error.value = "";
  name.value = result.name;
  content.value = result.content;
}
function attach() {
  const result = textAttachment(name.value, content.value);
  if (typeof result === "string") error.value = result;
  else if (props.takenNames.includes(result.name))
    error.value = `${result.name} は添付済みです。別のファイル名にしてください`;
  else emit("attach", result);
}
</script>
<template>
  <div class="log-composer" role="group" aria-label="ログファイルの作成">
    <label class="log-composer-name"
      >ファイル名<input
        v-model="name"
        type="text"
        maxlength="100"
        spellcheck="false"
        placeholder="log-1.txt"
    /></label>
    <label class="visually-hidden" for="log-composer-content">ログの内容</label>
    <div class="log-composer-editor">
      <textarea
        id="log-composer-content"
        v-model="content"
        rows="8"
        spellcheck="false"
        placeholder="ログを貼り付けるか、ファイルから読み込んでください"
        @keydown="onKeydown"
        @blur="onBlur"
      />
      <TabHint :visible="hintVisible" />
    </div>
    <p v-if="error" role="alert" class="log-composer-error">{{ error }}</p>
    <div class="log-composer-actions">
      <label class="button-secondary log-composer-upload">
        ファイルからアップロード
        <input
          type="file"
          class="visually-hidden"
          accept=".txt,.log,text/plain"
          @change="upload"
        />
      </label>
      <span class="log-composer-spacer" />
      <button type="button" class="button-secondary" @click="emit('cancel')">
        キャンセル
      </button>
      <button type="button" class="button-primary" @click="attach">添付</button>
    </div>
  </div>
</template>
<style scoped>
.log-composer {
  display: grid;
  gap: 8px;
  margin-top: 8px;
  padding: 12px;
  border: 1px solid var(--line-light);
}
.log-composer-name {
  display: grid;
  gap: 4px;
  font-size: 12px;
}
.log-composer-name input {
  padding: 6px 8px;
  border: 1px solid var(--ink);
  background: var(--paper);
  color: var(--ink);
  font: inherit;
  font-size: 14px;
}
.log-composer-editor {
  position: relative;
}
.log-composer textarea {
  display: block;
  font-family: "IBM Plex Mono", "Noto Sans JP", monospace;
  min-height: 160px;
  max-height: 50vh;
  padding: 12px;
  font-size: 13px;
  line-height: 1.6;
}
.log-composer-error {
  margin: 0;
  color: var(--danger, #c0392b);
  font-size: 12px;
}
.log-composer-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}
.log-composer-actions button,
.log-composer-upload {
  flex: none;
  width: auto;
  min-height: 0;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 12px;
}
.log-composer-spacer {
  flex: 1;
}
</style>
