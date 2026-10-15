<script setup lang="ts">
import { computed, ref } from "vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { Problem } from "~/features/models";
import { groupProblems } from "~/features/problem/group";
import { sectionLabel } from "~/features/problem/section";
import { statusLabels, type ProblemStatus } from "~/features/problem/status";

export type ProblemListItem = {
  problem: Problem;
  status: ProblemStatus;
  /** Minutes until the team can submit again; 0 when not cooling down. */
  cooldownMinutes: number;
};

const props = defineProps<{
  items: ProblemListItem[];
  pending?: boolean;
  error?: unknown;
}>();
defineEmits<{ retry: [] }>();

const selectedStatuses = ref<string[]>([]);
const selectedCategories = ref<string[]>([]);
const selectedScores = ref<string[]>([]);
const scoreOptions = {
  full: "満点",
  partial: "部分点",
  zero: "0点",
  unscored: "未採点",
};
const categories = computed(() =>
  [...new Set(props.items.map((item) => item.problem.category))].sort(),
);
const matchesScore = (problem: Problem) => {
  if (!selectedScores.value.length) return true;
  return selectedScores.value.some((filter) => {
    if (filter === "unscored") return problem.score == null;
    if (filter === "full")
      return problem.score != null && problem.score.score >= problem.maxScore;
    if (filter === "partial")
      return (
        (problem.score?.score ?? 0) > 0 &&
        (problem.score?.score ?? 0) < problem.maxScore
      );
    return problem.score?.score === 0;
  });
};
const groups = computed(() => {
  const byCode = new Map(props.items.map((item) => [item.problem.code, item]));
  return groupProblems(props.items.map((item) => item.problem))
    .map((g) => ({
      slug: g.key,
      items: g.problems
        .map((p) => byCode.get(p.code)!)
        .filter(
          ({ problem: p, status }) =>
            (!selectedCategories.value.length ||
              selectedCategories.value.includes(p.category)) &&
            (!selectedStatuses.value.length ||
              selectedStatuses.value.includes(status)) &&
            matchesScore(p),
        ),
    }))
    .filter((g) => g.items.length);
});
const hasFilters = computed(
  () =>
    selectedStatuses.value.length > 0 ||
    selectedCategories.value.length > 0 ||
    selectedScores.value.length > 0,
);
function resetFilters() {
  selectedStatuses.value = [];
  selectedCategories.value = [];
  selectedScores.value = [];
}
</script>
<template>
  <RequestState :pending="pending" :error="error" @retry="$emit('retry')"
    ><h1 class="visually-hidden">問題一覧</h1>
    <div v-if="hasFilters" class="table-tools">
      <button type="button" class="filter-reset" @click="resetFilters">
        絞り込みを解除
      </button>
    </div>
    <section v-for="group in groups" :key="group.slug" class="problem-table">
      <h2 class="problem-day-heading page-title">
        {{ sectionLabel(group.slug) }}の問題
      </h2>
      <div class="problem-row problem-heading">
        <span class="problem-heading-filter"
          ><details class="column-filter">
            <summary :class="{ 'is-filtered': selectedStatuses.length }">
              <span class="filter-label-desktop">問題ID / 状態</span>
              <span class="filter-label-mobile">状態</span>
            </summary>
            <div class="column-filter-menu">
              <fieldset>
                <legend>状態で絞り込む</legend>
                <label v-for="(label, key) in statusLabels" :key="key">
                  <input
                    v-model="selectedStatuses"
                    type="checkbox"
                    :value="key"
                  />
                  <span>{{ label }}</span>
                </label>
              </fieldset>
            </div>
          </details></span
        ><span class="problem-heading-filter score-filter"
          ><details class="column-filter">
            <summary :class="{ 'is-filtered': selectedScores.length }">
              得点
            </summary>
            <div class="column-filter-menu">
              <fieldset>
                <legend>得点で絞り込む</legend>
                <label v-for="(label, key) in scoreOptions" :key="key">
                  <input
                    v-model="selectedScores"
                    type="checkbox"
                    :value="key"
                  />
                  <span>{{ label }}</span>
                </label>
              </fieldset>
            </div>
          </details></span
        ><span class="problem-heading-filter"
          ><details class="column-filter">
            <summary :class="{ 'is-filtered': selectedCategories.length }">
              カテゴリ
            </summary>
            <div class="column-filter-menu">
              <fieldset>
                <legend>カテゴリで絞り込む</legend>
                <label v-for="categoryName in categories" :key="categoryName">
                  <input
                    v-model="selectedCategories"
                    type="checkbox"
                    :value="categoryName"
                  />
                  <span>{{ categoryName }}</span>
                </label>
              </fieldset>
            </div>
          </details></span
        ><span>問題</span>
      </div>
      <NuxtLink
        v-for="{ problem: p, status, cooldownMinutes } in group.items"
        :key="p.code"
        class="problem-row"
        :to="`/problems/${p.code}`"
        :class="{ 'is-answer-closed': !p.submissionStatus?.isSubmittable }"
        ><span class="problem-id"
          ><b>{{ p.code }}</b
          ><span class="visually-hidden">{{
            cooldownMinutes
              ? `再提出可能まで${cooldownMinutes}分`
              : statusLabels[status]
          }}</span></span
        ><span class="problem-score"
          ><span class="problem-score-value"
            ><b>{{ p.score?.score ?? "—" }}</b
            ><small>/ {{ p.maxScore }}</small></span
          ><em
            v-if="cooldownMinutes > 0"
            class="problem-cooldown-badge"
            aria-hidden="true"
            >再提出まで<b>{{ cooldownMinutes }}</b
            >分</em
          ><em
            v-else-if="status === 'pending'"
            class="problem-pending"
            aria-hidden="true"
            >採点中</em
          ></span
        ><span class="problem-category">{{ p.category }}</span
        ><span class="problem-title"
          ><small
            v-if="!p.submissionStatus?.isSubmittable"
            class="answer-closed-label"
            >受付時間外</small
          ><span class="problem-title-text">{{ p.title }}</span></span
        ></NuxtLink
      >
    </section>
    <p v-if="!groups.length" class="state-message">
      該当する問題はありません。
    </p></RequestState
  >
</template>
