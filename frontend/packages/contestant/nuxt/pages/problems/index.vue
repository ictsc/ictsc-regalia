<script setup lang="ts">
import { api } from "@ictsc/api";
import { groupProblems } from "~/features/problem/group";
import { fetchActivity } from "~/features/activity";
import { problemStatus, statusLabels } from "~/features/problem/status";
useHead({ title: "問題一覧" });
const { data, pending, error, refresh } = await useCompetition();
const { data: activity } = await useAsyncData("activity", () =>
  fetchActivity(api),
);
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
  [...new Set((data.value?.problems ?? []).map((p) => p.category))].sort(),
);
const now = useClock();
const problemCooldown = useProblemCooldown();
const cooldownMinutes = (
  problem: NonNullable<typeof data.value>["problems"][number],
) =>
  problemCooldown.remainingMinutes(
    problem.code,
    problem.nextSubmittableAt,
    now.value,
  );
const state = (p: NonNullable<typeof data.value>["problems"][number]) =>
  problemStatus(
    p,
    activity.value?.some((a) => a.problemCode === p.code),
    activity.value?.find((a) => a.problemCode === p.code)?.score ===
      undefined && !!activity.value?.some((a) => a.problemCode === p.code),
  );
const matchesScore = (
  problem: NonNullable<typeof data.value>["problems"][number],
) => {
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
const groups = computed(() =>
  groupProblems(data.value?.problems ?? [])
    .map((g) => ({
      slug: g.key,
      problems: g.problems.filter(
        (p) =>
          (!selectedCategories.value.length ||
            selectedCategories.value.includes(p.category)) &&
          (!selectedStatuses.value.length ||
            selectedStatuses.value.includes(state(p))) &&
          matchesScore(p),
      ),
    }))
    .filter((g) => g.problems.length),
);
const sectionName = (slug?: string) =>
  slug === "both"
    ? "両日"
    : /^day\d+$/.test(slug ?? "")
      ? `${slug!.slice(3)}日目`
      : slug;
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
  <main>
    <RequestState :pending="pending" :error="error" @retry="refresh"
      ><h1 class="visually-hidden">問題一覧</h1>
      <div v-if="hasFilters" class="table-tools">
        <button type="button" class="filter-reset" @click="resetFilters">
          絞り込みを解除
        </button>
      </div>
      <section v-for="group in groups" :key="group.slug" class="problem-table">
        <h2 class="problem-day-heading page-title">
          {{ sectionName(group.slug) }}の問題
        </h2>
        <div class="problem-row problem-heading">
          <span class="problem-heading-filter"
            ><details class="column-filter">
              <summary :class="{ 'is-filtered': selectedStatuses.length }">
                問題ID / 状態
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
          v-for="p in group.problems"
          :key="p.code"
          class="problem-row"
          :to="`/problems/${p.code}`"
          :class="{ 'is-answer-closed': !p.submissionStatus?.isSubmittable }"
          ><span
            class="problem-id"
            :class="[
              `status-${state(p)}`,
              { 'is-cooldown': cooldownMinutes(p) > 0 },
            ]"
            :data-cooldown="
              cooldownMinutes(p) ? `${cooldownMinutes(p)}分` : undefined
            "
            ><b>{{ p.code }}</b
            ><span class="visually-hidden">{{
              cooldownMinutes(p)
                ? `${statusLabels[state(p)]} / 再提出可能まで${cooldownMinutes(p)}分`
                : statusLabels[state(p)]
            }}</span></span
          ><span class="problem-score"
            ><b>{{ p.score?.score ?? "—" }}</b
            ><small>/ {{ p.maxScore }}</small></span
          ><span class="problem-category">{{ p.category }}</span
          ><span class="problem-title"
            ><small
              v-if="!p.submissionStatus?.isSubmittable"
              class="answer-closed-label"
              >受付時間外</small
            >{{ p.title }}</span
          ></NuxtLink
        >
      </section>
      <p v-if="!groups.length" class="state-message">
        該当する問題はありません。
      </p></RequestState
    >
  </main>
</template>
