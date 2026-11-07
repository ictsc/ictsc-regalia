<script setup lang="ts">
import { api, expectData } from "@ictsc/api";
import AdminPage from "~/components/layout/AdminPage.vue";
import GradingBoard, {
  type GradingChange,
} from "~/components/grading/GradingBoard.vue";
import GradingFilters from "~/components/grading/GradingFilters.vue";
import { updateAdminAnswerWorkflow } from "~/features/admin-api";
import {
  answerKey,
  answerPath,
  filterAnswers,
  gradableAnswers,
  type GradingAnswer,
} from "~/features/grading";

useHead({ title: "採点" });
const { viewer } = useAdminSession();
const myDiscordID = computed(() =>
  viewer.value?.state === "ADMIN" ? viewer.value.admin.discord.id : "",
);
const { data, error, pending, refresh } = await useAsyncData(
  "admin-answers",
  async () =>
    expectData(
      await api.GET("/api/v1/admin/answers", {
        params: { query: { include_marked: true } },
      }),
    ),
);
const teams = ref<string[]>([]);
const problems = ref<string[]>([]);
const recent = ref(false);
const perfect = ref(true);
const unscored = ref(false);
const secondaryQualifier = computed(() =>
  data.value?.answers.some((a) => a.competition_mode === "SECONDARY_QUALIFIER"),
);
const busyKey = ref("");
const actionError = ref("");
const { now } = useSubmissionAge();
// 他の運営の担当変更や新しい提出を、再読み込みなしで反映する。
const { streamError } = useAdminAnswerStream(() => void refresh());

onMounted(() => {
  try {
    const saved = JSON.parse(
      localStorage.getItem("regalia/admin/answer-filters") ?? "null",
    );
    if (saved) {
      teams.value = saved.teams ?? [];
      problems.value = saved.problems ?? [];
      recent.value = !!saved.recent;
      perfect.value = saved.perfect !== false;
      unscored.value = !!saved.unscored;
    }
  } catch {
    /* Filters work without local storage. */
  }
});
watch(
  [teams, problems, recent, perfect, unscored],
  () => {
    try {
      localStorage.setItem(
        "regalia/admin/answer-filters",
        JSON.stringify({
          teams: teams.value,
          problems: problems.value,
          recent: recent.value,
          perfect: perfect.value,
          unscored: unscored.value,
        }),
      );
    } catch {
      /* Filters work without local storage. */
    }
  },
  { deep: true },
);
const gradable = computed(() => gradableAnswers(data.value?.answers ?? []));
const filtered = computed(() =>
  filterAnswers(
    gradable.value,
    {
      teams: teams.value,
      problems: problems.value,
      recent: recent.value,
      perfect: perfect.value,
      unscored: unscored.value,
    },
    now.value,
  ),
);
const problemCodes = computed(() => [
  ...new Set(gradable.value.map((a) => a.problem.code)),
]);
const teamNames = computed(() => [
  ...new Set(gradable.value.map((a) => a.team.name)),
]);
async function update(answer: GradingAnswer, change: GradingChange) {
  if (change.status === "COMPLETED" && !answer.score) {
    await navigateTo(answerPath(answer));
    return;
  }
  busyKey.value = answerKey(answer);
  actionError.value = "";
  try {
    const updated = await updateAdminAnswerWorkflow(api, answer, change);
    // 待っている間に SSE の再取得がより新しい版を読み込んでいたら、
    // 古い PATCH のレスポンスで巻き戻さない。
    if (data.value)
      data.value = {
        ...data.value,
        answers: data.value.answers.map((a) =>
          answerKey(a) === answerKey(answer) &&
          updated.workflow.revision >= a.workflow.revision
            ? updated
            : a,
        ),
      };
  } catch (cause) {
    actionError.value =
      cause instanceof Error ? cause.message : "更新に失敗しました。";
    await refresh();
  } finally {
    busyKey.value = "";
  }
}
</script>

<template>
  <main class="workspace grading-workspace">
    <AdminPage title="採点">
      <template #lead
        >自分の担当回答を各列の先頭に表示します。カードを開くと採点できます。<template
          v-if="secondaryQualifier"
          ><br />二次予選モード:
          各チーム・問題の最終答案だけを表示します。採点は競技終了後に行えます。</template
        ></template
      >
    </AdminPage>
    <GradingFilters
      v-model:teams="teams"
      v-model:problems="problems"
      v-model:recent="recent"
      v-model:perfect="perfect"
      v-model:unscored="unscored"
      :problem-codes="problemCodes"
      :team-names="teamNames"
    />
    <p v-if="actionError" role="alert">{{ actionError }}</p>
    <p v-if="streamError" role="alert" class="admin-status">
      状態更新に再接続しています。
    </p>
    <RequestState :error="error" :pending="pending && !data" @retry="refresh">
      <GradingBoard
        :answers="filtered"
        :my-discord-id="myDiscordID"
        :busy-key="busyKey"
        :now="now"
        @update="update"
      />
    </RequestState>
  </main>
</template>
