<script setup lang="ts">
import type { components } from "@ictsc/api";

type ContestPhase = components["schemas"]["ContestPhase"];

defineProps<{
  contestPhase: components["schemas"]["DevContestPhase"];
  busy: boolean;
}>();
defineEmits<{ setPhase: [phase: ContestPhase] }>();

const phases: { phase: ContestPhase; label: string; description: string }[] = [
  {
    phase: "BEFORE",
    label: "競技前",
    description: "最初のセクションが10分後に始まる",
  },
  {
    phase: "DURING",
    label: "競技中",
    description: "最初のセクションが10分前に始まった",
  },
  {
    phase: "AFTER",
    label: "競技後",
    description: "最後のセクションが10分前に終わった",
  },
];
const labelOf = (phase: ContestPhase) =>
  phases.find((p) => p.phase === phase)?.label ?? phase;
const format = (value: string) => new Date(value).toLocaleString("ja-JP");
</script>
<template>
  <section class="admin-panel dev-contest-phase">
    <h2>開発用: 競技状態の切り替え</h2>
    <p>
      開発モードかつ二次予選モードでのみ表示されます。セクションの数・長さ・間隔は保ったまま、現在時刻を基準に全体を前後へずらします。
    </p>
    <p>
      現在の状態: <strong>{{ labelOf(contestPhase.phase) }}</strong>
    </p>
    <ul class="dev-contest-phase-sections">
      <li v-for="section in contestPhase.sections" :key="section.slug">
        <code>{{ section.slug }}</code>
        {{ format(section.beginning) }} 〜 {{ format(section.ending) }}
      </li>
    </ul>
    <div class="actions">
      <button
        v-for="option in phases"
        :key="option.phase"
        type="button"
        :class="
          option.phase === contestPhase.phase
            ? 'button-primary'
            : 'button-secondary'
        "
        :disabled="busy"
        :title="option.description"
        @click="$emit('setPhase', option.phase)"
      >
        {{ option.label }}にする
      </button>
    </div>
  </section>
</template>

<style scoped>
.dev-contest-phase-sections {
  margin: 0 0 16px;
  padding-left: 1.2em;
  font-variant-numeric: tabular-nums;
}
</style>
