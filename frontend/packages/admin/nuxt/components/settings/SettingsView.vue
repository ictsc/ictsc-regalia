<script setup lang="ts">
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { components } from "@ictsc/api";
import DevContestPhasePanel from "./DevContestPhasePanel.vue";

defineProps<{
  /** The active content's rule.md. */
  ruleMarkdown: string;
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
  setPhase: [phase: components["schemas"]["ContestPhase"]];
  retry: [];
}>();
</script>
<template>
  <h1>競技設定</h1>
  <p class="admin-lead">
    {{
      secondaryQualifier
        ? "参加者に表示するルールを確認します。二次予選モードでは順位は最終公開まで非公開のため、凍結は使いません。"
        : "参加者に表示するルールとランキングの凍結時刻を確認します。"
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
      <h2>ルール</h2>
      <p>
        ルールはコンテンツリポジトリの <code>rule.md</code>
        で管理し、コンテンツ更新時に反映されます。
      </p>
      <MarkdownContent v-if="ruleMarkdown.trim()" :source="ruleMarkdown" />
      <p v-else class="muted">
        有効なコンテンツに <code>rule.md</code> がありません。
      </p>
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
