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
const status = ref("all"),
  day = ref("all"),
  category = ref("all");
const state = (p: NonNullable<typeof data.value>["problems"][number]) =>
  problemStatus(
    p,
    activity.value?.some((a) => a.problemCode === p.code),
    activity.value?.find((a) => a.problemCode === p.code)?.score ===
      undefined && !!activity.value?.some((a) => a.problemCode === p.code),
  );
const groups = computed(() =>
  groupProblems(data.value?.problems ?? [])
    .map((g) => ({
      slug: g.key,
      problems: g.problems.filter(
        (p) =>
          (day.value === "all" || p.sectionSlug === day.value) &&
          (category.value === "all" || p.category === category.value) &&
          (status.value === "all" || state(p) === status.value),
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
</script>
<template>
  <main>
    <RequestState :pending="pending" :error="error" @retry="refresh"
      ><h1 class="visually-hidden">問題一覧</h1>
      <div class="table-tools">
        <details class="filters-panel">
          <summary>絞り込み</summary>
          <div class="filters-panel-content">
            <label for="status-filter">採点状況</label
            ><select id="status-filter" v-model="status">
              <option value="all">すべて</option>
              <option
                v-for="(label, key) in statusLabels"
                :key="key"
                :value="key"
              >
                {{ label }}
              </option></select
            ><label for="day-filter">出題日</label
            ><select id="day-filter" v-model="day">
              <option value="all">すべて</option>
              <option
                v-for="s in new Set(data?.problems.map((p) => p.sectionSlug))"
                :key="s"
                :value="s"
              >
                {{ sectionName(s) }}
              </option></select
            ><label for="category-filter">カテゴリ</label
            ><select id="category-filter" v-model="category">
              <option value="all">すべて</option>
              <option
                v-for="c in new Set(data?.problems.map((p) => p.category))"
                :key="c"
              >
                {{ c }}
              </option>
            </select>
          </div>
        </details>
      </div>
      <section v-for="group in groups" :key="group.slug" class="problem-table">
        <h2 class="problem-day-heading">{{ sectionName(group.slug) }}の問題</h2>
        <div class="problem-row problem-heading">
          <span>問題ID / 状態</span><span>カテゴリ</span><span>問題</span
          ><span>得点</span>
        </div>
        <NuxtLink
          v-for="p in group.problems"
          :key="p.code"
          class="problem-row"
          :to="`/problems/${p.code}`"
          :class="{ 'is-answer-closed': !p.submissionStatus?.isSubmittable }"
          ><span class="problem-id" :class="`status-${state(p)}`"
            ><b>{{ p.code }}</b
            ><span class="visually-hidden">{{
              statusLabels[state(p)]
            }}</span></span
          ><span class="problem-category">{{ p.category }}</span
          ><span class="problem-title"
            ><small
              v-if="!p.submissionStatus?.isSubmittable"
              class="answer-closed-label"
              >受付時間外</small
            >{{ p.title }}</span
          ><span
            class="problem-score"
            :class="{ 'is-pending': state(p) === 'pending' }"
            ><b>{{ p.score?.score ?? "—" }}</b
            ><small
              >{{ state(p) === "pending" ? "採点中 / " : "/ "
              }}{{ p.maxScore }}</small
            ></span
          ></NuxtLink
        >
      </section>
      <p v-if="!groups.length" class="state-message">
        該当する問題はありません。
      </p></RequestState
    >
  </main>
</template>
