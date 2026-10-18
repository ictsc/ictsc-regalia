<script setup lang="ts">
import { defaultGradingFilters } from "~/features/grading";

const teams = defineModel<string[]>("teams", { required: true });
const problems = defineModel<string[]>("problems", { required: true });
const recent = defineModel<boolean>("recent", { required: true });
const perfect = defineModel<boolean>("perfect", { required: true });
const unscored = defineModel<boolean>("unscored", { required: true });
defineProps<{
  /** Distinct problem codes among the loaded answers. */
  problemCodes: string[];
  /** Distinct team names among the loaded answers. */
  teamNames: string[];
}>();

function clear() {
  const defaults = defaultGradingFilters();
  teams.value = defaults.teams;
  problems.value = defaults.problems;
  recent.value = defaults.recent;
  perfect.value = defaults.perfect;
  unscored.value = defaults.unscored;
}
</script>
<template>
  <section class="admin-panel">
    <h2>回答を絞り込む</h2>
    <div class="admin-filter-grid">
      <label
        >問題（複数選択可）<select v-model="problems" multiple size="4">
          <option v-for="p in problemCodes" :key="p">
            {{ p }}
          </option>
        </select></label
      >
      <label
        >チーム（複数選択可）<select v-model="teams" multiple size="4">
          <option v-for="t in teamNames" :key="t">
            {{ t }}
          </option>
        </select></label
      >
    </div>
    <div class="actions">
      <label><input v-model="recent" type="checkbox" /> 直近20分のみ</label>
      <label><input v-model="perfect" type="checkbox" /> 満点回答を表示</label>
      <label><input v-model="unscored" type="checkbox" /> 未採点のみ</label>
      <button class="button-secondary" @click="clear">条件をクリア</button>
    </div>
  </section>
</template>
