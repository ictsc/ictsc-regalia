<script setup lang="ts">
import type { ScoreModel } from "~/features/models";

defineProps<{
  code: string;
  title: string;
  category: string;
  score?: ScoreModel;
  maxScore: number;
  /** Minutes until the team can answer again; 0 when not cooling down. */
  cooldownMinutes: number;
  /** Answers are being accepted now. */
  open: boolean;
}>();
</script>
<template>
  <header class="problem-hero">
    <div class="hero-copy">
      <h1>{{ code }}:{{ title }}.</h1>
      <p>{{ category }}</p>
    </div>
    <div class="problem-summary-side">
      <p class="problem-score-summary">
        <strong>{{ score ? `${score.score}点` : "—" }}</strong
        ><span>/ {{ maxScore }}点満点</span>
      </p>
      <p
        class="problem-cooldown"
        :class="{ 'is-answer-closed': !open && !cooldownMinutes }"
      >
        <template v-if="cooldownMinutes"
          ><span>再回答可能まで</span><strong>{{ cooldownMinutes }}</strong
          ><small>分</small></template
        ><span v-else>{{ open ? "回答受付中" : "回答終了" }}</span>
      </p>
    </div>
  </header>
</template>
