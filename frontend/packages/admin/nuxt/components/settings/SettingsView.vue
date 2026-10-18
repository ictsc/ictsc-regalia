<script setup lang="ts">
import { ref } from "vue";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";

const props = defineProps<{
  initialMarkdown: string;
  /** Freeze time as a `datetime-local` value in the browser's time zone. */
  initialFreeze: string;
  pending?: boolean;
  error?: unknown;
  busy: boolean;
  /** Result of the last operation. */
  message: string;
}>();
defineEmits<{
  saveRule: [markdown: string];
  /** Empty string clears the freeze. */
  saveFreeze: [freeze: string];
  retry: [];
}>();

// Initialized once like the page did; edits are kept across refreshes.
const markdown = ref(props.initialMarkdown);
const freeze = ref(props.initialFreeze);
</script>
<template>
  <h1>競技設定</h1>
  <p class="admin-lead">
    参加者に表示するルールと、ランキングの凍結時刻を管理します。
  </p>
  <p role="status">{{ message }}</p>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><section class="admin-panel">
      <h2>ルール編集</h2>
      <p>Markdown で編集し、下のプレビューで表示を確認できます。</p>
      <form class="form-grid" @submit.prevent="$emit('saveRule', markdown)">
        <label>Markdown<textarea v-model="markdown" rows="18" /></label
        ><button class="button-primary" :disabled="busy">ルールを更新</button>
      </form>
    </section>
    <section class="admin-panel">
      <h2>ルールのプレビュー</h2>
      <MarkdownContent :source="markdown" />
    </section>
    <section class="admin-panel">
      <h2>ランキング凍結</h2>
      <p>
        ブラウザのローカル時刻で入力します。空欄で保存すると凍結を解除します。
      </p>
      <form class="form-grid" @submit.prevent="$emit('saveFreeze', freeze)">
        <label>凍結時刻<input v-model="freeze" type="datetime-local" /></label
        ><button class="button-primary" :disabled="busy">凍結時刻を保存</button>
      </form>
    </section></RequestState
  >
</template>
