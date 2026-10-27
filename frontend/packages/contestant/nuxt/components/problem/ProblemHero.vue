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
      <!-- Category, score and answer status share one meta line under the
           title, so the answer column starts with its own heading. -->
      <p class="problem-meta">
        <span class="problem-category-label">{{ category }}</span>
        <span class="problem-score-summary">
          <strong>{{ score ? `${score.score}点` : "—" }}</strong
          ><span>/ {{ maxScore }}点満点</span>
        </span>
        <span
          class="problem-cooldown"
          :class="{ 'is-answer-closed': !open && !cooldownMinutes }"
        >
          <template v-if="cooldownMinutes"
            ><span>再回答可能まで</span><strong>{{ cooldownMinutes }}</strong
            ><small>分</small></template
          ><span v-else>{{ open ? "回答受付中" : "回答終了" }}</span>
        </span>
      </p>
    </div>
  </header>
</template>
