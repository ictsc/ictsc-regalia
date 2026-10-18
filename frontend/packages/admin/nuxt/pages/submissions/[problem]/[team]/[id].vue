<script setup lang="ts">
import { api, expectData } from "@ictsc/api";
import GradingDetailView, {
  type MarkingInput,
} from "~/components/grading/GradingDetailView.vue";
import { defaultPaneLayout, restorePaneLayout } from "~/features/grading-panes";
import {
  getAdminAnswer,
  listAdminMarkingResults,
  createAdminMarkingResult,
  deleteAdminMarkingLineComment,
  createAdminActions,
  updateAdminAnswerWorkflow,
} from "~/features/admin-api";
definePageMeta({ key: (route) => route.fullPath });
const route = useRoute();
const reference = {
  problem_code: String(route.params.problem),
  team_code: Number(route.params.team),
  answer_number: Number(route.params.id),
};
useHead({
  title: `${reference.problem_code} / Team ${reference.team_code} / #${reference.answer_number}`,
});
const { data, error, pending, refresh } = await useAsyncData(
  `admin-answer:${route.fullPath}`,
  async () => {
    const answer = await getAdminAnswer(api, reference);
    const [problem, active, previous, marks] = await Promise.all([
      createAdminActions(api).getProblem(
        reference.problem_code,
        answer.content_commit,
      ),
      createAdminActions(api).getProblem(reference.problem_code),
      api
        .GET("/api/v1/admin/answers", {
          params: {
            query: {
              include_marked: true,
              team_code: reference.team_code,
              problem_code: reference.problem_code,
            },
          },
        })
        .then(expectData),
      listAdminMarkingResults(api),
    ]);
    return {
      answer,
      problem: problem.problem,
      active: active.problem,
      previous: previous.answers
        .filter((a) => a.reference.answer_number < reference.answer_number)
        .toSorted(
          (a, b) => b.reference.answer_number - a.reference.answer_number,
        ),
      marks: marks.filter(
        (m) =>
          m.answer.problem_code === reference.problem_code &&
          m.answer.team_code === reference.team_code &&
          m.answer.answer_number === reference.answer_number,
      ),
    };
  },
);

const { elapsedMinutes } = useSubmissionAge();
const { busy, message, run } = useMutation(refresh);
async function deleteLineComment(markingId: string, commentIndex: number) {
  await run(
    "行コメントを削除",
    () => deleteAdminMarkingLineComment(api, markingId, commentIndex),
    {
      message: "この行コメントを削除しますか？削除後は回答上に表示されません。",
      destructive: true,
    },
  );
}
const workflowBusy = ref(false);
const workflowError = ref("");
function submitMarking({ score, comment, lineComments }: MarkingInput) {
  return run(
    "採点結果を保存",
    () =>
      createAdminMarkingResult(api, reference, score, comment, lineComments),
    `得点: ${score}\n全体コメント: ${comment}\n行コメント: ${lineComments.length}件\nこの採点結果を送信しますか？`,
  );
}
async function assignment(action: "CLAIM_SELF" | "RESET_TO_DEFAULT") {
  if (!data.value) return;
  workflowBusy.value = true;
  workflowError.value = "";
  try {
    const updated = await updateAdminAnswerWorkflow(api, data.value.answer, {
      assignment: action,
    });
    data.value = { ...data.value, answer: updated };
  } catch (cause) {
    workflowError.value =
      cause instanceof Error ? cause.message : "担当を更新できませんでした。";
    await refresh();
  } finally {
    workflowBusy.value = false;
  }
}

const collapsed = ref(defaultPaneLayout().collapsed);
const widths = ref(defaultPaneLayout().widths);
onMounted(() => {
  try {
    const layout = restorePaneLayout(
      localStorage.getItem("regalia/admin/grading-panes-v2"),
      localStorage.getItem("regalia/admin/grading-panes"),
      { collapsed: collapsed.value, widths: widths.value },
    );
    collapsed.value = layout.collapsed;
    widths.value = layout.widths;
  } catch {
    /* Default pane layout remains available. */
  }
});
watch(
  [collapsed, widths],
  () => {
    try {
      localStorage.setItem(
        "regalia/admin/grading-panes-v2",
        JSON.stringify({ collapsed: collapsed.value, widths: widths.value }),
      );
    } catch {
      /* Pane layout does not depend on storage. */
    }
  },
  { deep: true },
);
</script>

<template>
  <main class="workspace grading-detail-workspace">
    <NuxtLink class="text-link" to="/submissions">← 採点一覧</NuxtLink>
    <h1>
      {{ reference.problem_code }} / Team {{ reference.team_code }} / #{{
        reference.answer_number
      }}
    </h1>
    <p v-if="data" class="admin-lead">
      {{ data.answer.team.name }} · {{ data.answer.author.display_name }} ·
      <time :datetime="data.answer.submitted_at">{{
        new Date(data.answer.submitted_at).toLocaleString("ja-JP")
      }}</time>
      <span class="grading-detail-age"
        >提出から {{ elapsedMinutes(data.answer.submitted_at) }} 分</span
      >
    </p>
    <p role="status">{{ message }}</p>
    <p v-if="workflowError" role="alert">{{ workflowError }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh">
      <GradingDetailView
        v-if="data"
        v-model:collapsed="collapsed"
        v-model:widths="widths"
        :answer="data.answer"
        :problem="data.problem"
        :active-commit="data.active.content_commit"
        :previous="data.previous"
        :marks="data.marks"
        :busy="busy"
        :workflow-busy="workflowBusy"
        :submit-marking="submitMarking"
        @assignment="assignment"
        @delete-line-comment="deleteLineComment"
      />
    </RequestState>
  </main>
</template>

<style scoped>
.grading-detail-workspace {
  max-width: none;
}
.grading-detail-age {
  display: inline-block;
  margin-left: 8px;
  padding: 2px 8px;
  border: 1px solid var(--line-light);
  border-radius: 999px;
  color: var(--accent);
  font-size: 12px;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
}
</style>
