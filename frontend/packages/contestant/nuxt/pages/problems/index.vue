<script setup lang="ts">
import { api } from "@ictsc/api";
import ProblemListView, {
  type ProblemListItem,
} from "~/components/problem/ProblemListView.vue";
import { fetchActivity } from "~/features/activity";
import { problemStatus } from "~/features/problem/status";
useHead({ title: "問題一覧" });
const { data, pending, error, refresh } = await useCompetition();
const { data: activity } = await useAsyncData("activity", () =>
  fetchActivity(api),
);
const now = useClock();
const problemCooldown = useProblemCooldown();
const items = computed<ProblemListItem[]>(() =>
  (data.value?.problems ?? []).map((problem) => {
    const submitted = activity.value?.some(
      (a) => a.problemCode === problem.code,
    );
    return {
      problem,
      status: problemStatus(
        problem,
        submitted,
        activity.value?.find((a) => a.problemCode === problem.code)?.score ===
          undefined && !!submitted,
      ),
      cooldownMinutes: problemCooldown.remainingMinutes(
        problem.code,
        problem.nextSubmittableAt,
        now.value,
      ),
    };
  }),
);
</script>
<template>
  <main>
    <ProblemListView
      :items="items"
      :pending="pending"
      :error="error"
      @retry="refresh"
    />
  </main>
</template>
