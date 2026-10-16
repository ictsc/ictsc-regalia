<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from "vue";
import type { Problem } from "~/features/models";
import { groupProblems } from "~/features/problem/group";
import { sectionLabel } from "~/features/problem/section";

export type ProblemRailItem = {
  problem: Problem;
  /** Minutes until the team can submit again; 0 when not cooling down. */
  cooldownMinutes: number;
  /** The latest submission is waiting for grading. */
  pending: boolean;
};

const props = defineProps<{ items: ProblemRailItem[]; currentCode: string }>();

const groups = computed(() => {
  const byCode = new Map(props.items.map((item) => [item.problem.code, item]));
  return groupProblems(props.items.map((item) => item.problem)).map((g) => ({
    key: g.key,
    items: g.problems.map((p) => byCode.get(p.code)!),
  }));
});

const problemRail = ref<HTMLElement | null>(null);
// On narrow screens the rail scrolls horizontally; keep the current problem in view.
async function showCurrentRailItem() {
  await nextTick();
  if (
    typeof window === "undefined" ||
    !window.matchMedia("(max-width: 760px)").matches
  )
    return;
  const container = problemRail.value;
  const active = container?.querySelector<HTMLElement>('[aria-current="page"]');
  if (!container || !active) return;
  const railBounds = container.getBoundingClientRect();
  const activeBounds = active.getBoundingClientRect();
  container.scrollLeft +=
    activeBounds.left -
    railBounds.left -
    (railBounds.width - activeBounds.width) / 2;
}
onMounted(showCurrentRailItem);
watch(groups, showCurrentRailItem, { flush: "post" });
</script>
<template>
  <nav
    id="problem-rail"
    ref="problemRail"
    class="problem-rail"
    aria-label="問題ナビゲーション"
  >
    <NuxtLink class="problem-rail-label" to="/problems">問題一覧</NuxtLink>
    <section
      v-for="group in groups"
      :key="group.key"
      class="problem-rail-group"
    >
      <h2>
        <span class="visually-hidden">ステージ: </span>
        <span class="rail-stage-label">{{ sectionLabel(group.key) }}</span>
      </h2>
      <NuxtLink
        v-for="{ problem: p, cooldownMinutes, pending } in group.items"
        :key="p.code"
        :to="`/problems/${p.code}`"
        :class="{
          'is-active': p.code === currentCode,
          'is-answer-closed': !p.submissionStatus?.isSubmittable,
          'has-status': cooldownMinutes > 0 || pending,
        }"
        :aria-current="p.code === currentCode ? 'page' : undefined"
        :aria-label="`${p.code} ${p.title} ${p.score ? `${p.score.score}点 / ${p.maxScore}点満点` : '得点未確定'}${cooldownMinutes ? ` 再提出可能まで${cooldownMinutes}分` : pending ? ' 最新の提出を採点中' : ''}`"
        ><b>{{ p.code }}</b
        ><span class="rail-score">
          <strong>{{ p.score ? p.score.score : "—" }}</strong>
          <small>/ {{ p.maxScore }}</small></span
        ><span class="rail-title">{{ p.title }}</span
        ><span v-if="cooldownMinutes" class="rail-cooldown" aria-hidden="true"
          >あと<b>{{ cooldownMinutes }}</b
          >分</span
        ><em v-else-if="pending" class="rail-pending" aria-hidden="true"
          >採点中</em
        ><span class="problem-preview"
          ><strong>{{ p.title }}</strong
          ><small>{{ p.category }}</small
          ><span class="preview-status">{{
            pending ? "最新の提出を採点中" : p.score ? "採点済み" : "未回答"
          }}</span
          ><span class="preview-score"
            >{{ p.score ? `${p.score.score}点` : "—" }} /
            {{ p.maxScore }}点満点</span
          ><span v-if="cooldownMinutes" class="preview-cooldown"
            >再提出可能まで <b>{{ cooldownMinutes }}</b
            >分</span
          ><span
            v-else-if="!p.submissionStatus?.isSubmittable"
            class="preview-closed"
            >回答終了 / 閲覧のみ</span
          ></span
        ></NuxtLink
      >
    </section>
  </nav>
</template>
