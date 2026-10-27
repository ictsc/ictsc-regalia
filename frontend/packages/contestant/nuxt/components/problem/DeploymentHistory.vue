<script setup lang="ts">
import type { Deployment } from "~/features/deployment";

defineProps<{
  deployments: Deployment[];
  penaltyThreshold: number;
  /** A redeploy request is being sent. */
  deploying: boolean;
  /** The deployment event stream is reconnecting. */
  streamError: boolean;
  resetDisabled: boolean;
}>();
defineEmits<{ reset: [] }>();
</script>
<template>
  <section class="content-section">
    <h2>再展開</h2>
    <p>減点なしの上限: {{ penaltyThreshold }}回</p>
    <!-- Kept in the reading column, away from the answer submit button; the
         dialog still asks before anything is reset. -->
    <button
      class="button-secondary button-reset detail-reset"
      type="button"
      :disabled="resetDisabled"
      @click="$emit('reset')"
    >
      環境をリセット
    </button>
    <p v-if="deploying" role="status">QUEUED — 再展開を要求中</p>
    <p v-if="streamError" role="alert">更新の接続を確認しています</p>
    <p v-if="!deployments.length">再展開履歴はありません。</p>
    <ul>
      <li v-for="d in deployments" :key="d.revision">
        #{{ d.revision }} {{ d.status }} / 減点 {{ d.penalty }} / 残り
        {{ d.allowedDeploymentCount }}回
      </li>
    </ul>
  </section>
</template>
