<script setup lang="ts">
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import AdminPage from "~/components/layout/AdminPage.vue";
import AdminPanel from "~/components/layout/AdminPanel.vue";
import AdminFacts from "~/components/layout/AdminFacts.vue";
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
  <AdminPage
    title="競技設定"
    :lead="
      secondaryQualifier
        ? '参加者に表示するルールを確認します。二次予選モードでは順位は最終公開まで非公開のため、凍結は使いません。'
        : '参加者に表示するルールとランキングの凍結時刻を確認します。'
    "
    :message="message"
  >
    <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
      ><DevContestPhasePanel
        v-if="contestPhase?.enabled"
        :contest-phase="contestPhase"
        :busy="busy"
        @set-phase="$emit('setPhase', $event)" />
      <AdminPanel title="ルール">
        <p>
          ルールはコンテンツリポジトリの <code>rule.md</code>
          で管理し、コンテンツ更新時に反映されます。
        </p>
        <MarkdownContent v-if="ruleMarkdown.trim()" :source="ruleMarkdown" />
        <p v-else class="admin-empty">
          有効なコンテンツに <code>rule.md</code> がありません。
        </p>
      </AdminPanel>
      <AdminPanel v-if="!secondaryQualifier" title="ランキング凍結">
        <p>
          凍結時刻はコンテンツリポジトリの <code>schedule.yaml</code> の
          <code>ranking_freeze_at</code>
          で設定し、コンテンツ更新時に反映されます。
        </p>
        <AdminFacts
          :items="[
            {
              label: '現在の凍結時刻',
              value: initialFreeze
                ? new Date(initialFreeze).toLocaleString('ja-JP')
                : '未設定',
            },
          ]"
        /> </AdminPanel
    ></RequestState>
  </AdminPage>
</template>
