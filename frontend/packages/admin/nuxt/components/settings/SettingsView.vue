<script setup lang="ts">
import { ref } from "vue";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { components } from "@ictsc/api";
import DevContestPhasePanel from "./DevContestPhasePanel.vue";

const props = defineProps<{
  initialMarkdown: string;
  /** Freeze time as a `datetime-local` value in the browser's time zone. */
  initialFreeze: string;
  pending?: boolean;
  error?: unknown;
  busy: boolean;
  /** Result of the last operation. */
  message: string;
  /** Secondary qualifier mode does not use the ranking freeze. */
  secondaryQualifier?: boolean;
  /** Shown only when the server enables development phase switching. */
  contestPhase?: components["schemas"]["DevContestPhase"];
}>();
defineEmits<{
  saveRule: [markdown: string];
  setPhase: [phase: components["schemas"]["ContestPhase"]];
  retry: [];
}>();

// Initialized once like the page did; edits are kept across refreshes.
const markdown = ref(props.initialMarkdown);
</script>
<template>
  <h1>競技設定</h1>
  <p class="admin-lead">
    {{
      secondaryQualifier
        ? "参加者に表示するルールを管理します。二次予選モードでは順位は最終公開まで非公開のため、凍結は使いません。"
        : "参加者に表示するルールを管理し、ランキングの凍結時刻を確認します。"
    }}
  </p>
  <p role="status">{{ message }}</p>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><DevContestPhasePanel
      v-if="contestPhase?.enabled"
      :contest-phase="contestPhase"
      :busy="busy"
      @set-phase="$emit('setPhase', $event)"
    />
    <section class="admin-panel">
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
    <section v-if="!secondaryQualifier" class="admin-panel">
      <h2>ランキング凍結</h2>
      <p>
        凍結時刻はコンテンツリポジトリの <code>schedule.yaml</code> の
        <code>ranking_freeze_at</code>
        で設定し、コンテンツ更新時に反映されます。
      </p>
      <p>
        現在の凍結時刻:
        <strong>{{
          initialFreeze
            ? new Date(initialFreeze).toLocaleString("ja-JP")
            : "未設定"
        }}</strong>
      </p>
    </section></RequestState
  >
</template>
