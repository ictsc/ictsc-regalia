<script setup lang="ts">
import { api, ApiError } from "@ictsc/api";
import { fetchProblem } from "~/features/problem";
import { fetchAnswer, fetchAnswers, submitAnswer } from "~/features/answer";
import {
  fetchDeployments,
  deploy,
  subscribeDeployments,
  mapDeployment,
} from "~/features/deployment";
import { draftKey, loadDraft, saveDraft } from "~/features/draft";
import { fetchActivity } from "~/features/activity";
import { remainingCooldownSeconds } from "~/features/problem/cooldown";
import type { AttachmentFile } from "@ictsc/ui/attachments";
import AnswerForm from "~/components/problem/AnswerForm.vue";
import AnswerResizer from "~/components/problem/AnswerResizer.vue";
import {
  loadAnswerWidth,
  saveAnswerWidth,
} from "~/features/problem/answer-width";
import AnswerHistory from "~/components/problem/AnswerHistory.vue";
import DeploymentHistory from "~/components/problem/DeploymentHistory.vue";
import ProblemHero from "~/components/problem/ProblemHero.vue";
import ProblemRail, {
  type ProblemRailItem,
} from "~/components/problem/ProblemRail.vue";
import ResetDialog from "~/components/problem/ResetDialog.vue";
definePageMeta({ key: (route) => route.fullPath });
const route = useRoute();
const code = String(route.params.code);
const requestedAnswer = Array.isArray(route.query.answer)
  ? route.query.answer[0]
  : route.query.answer;
const selectedAnswerNumber =
  typeof requestedAnswer === "string" &&
  Number.isSafeInteger(Number(requestedAnswer)) &&
  Number(requestedAnswer) > 0
    ? Number(requestedAnswer)
    : null;
const { viewer } = useSession();
const { data: competition } = await useCompetition();
const { data: activity } = await useAsyncData("activity", () =>
  fetchActivity(api),
);
const { data, error, pending, refresh } = await useAsyncData(
  `problem:${code}:answer:${selectedAnswerNumber ?? "none"}`,
  async () => {
    const [problem, answers, deployments, selectedAnswer] = await Promise.all([
      fetchProblem(api, code),
      fetchAnswers(api, code),
      fetchDeployments(api, code),
      selectedAnswerNumber == null
        ? Promise.resolve(null)
        : fetchAnswer(api, code, selectedAnswerNumber),
    ]);
    return { problem, answers, deployments, selectedAnswer };
  },
  { deep: true },
);
useHead({
  title: computed(() =>
    data.value ? `${code} ${data.value.problem.title}` : code,
  ),
  bodyAttrs: { class: "detail-page" },
});
const contentGrid = ref<HTMLElement | null>(null);
const answerWidth = ref<number>();
onMounted(() => {
  try {
    answerWidth.value = loadAnswerWidth(localStorage);
  } catch {
    /* The default width works without storage. */
  }
});
watch(answerWidth, (width) => {
  try {
    saveAnswerWidth(localStorage, width);
  } catch {
    /* The width just is not remembered. */
  }
});
const body = ref(""),
  attachments = ref<AttachmentFile[]>([]),
  savedAt = ref(""),
  saveError = ref(""),
  actionError = ref<unknown>(),
  sending = ref(false),
  deploying = ref(false),
  submitted = ref(false),
  resetOpen = ref(false);
const now = useClock(false),
  retryAt = ref(0);
const displayNow = useClock();
const problemCooldown = useProblemCooldown();
const key = computed(() =>
  viewer.value?.state === "CONTESTANT"
    ? draftKey(viewer.value.profile.name, viewer.value.team.code, code)
    : null,
);
onMounted(() => {
  if (key.value) {
    try {
      const draft = loadDraft(localStorage, key.value);
      body.value = draft?.body ?? "";
      savedAt.value = draft?.savedAt ?? "";
    } catch {
      saveError.value =
        "下書きを読み込めません。このブラウザの保存設定を確認してください。";
    }
  }
});
watch(body, (value) => {
  if (key.value) {
    try {
      savedAt.value = saveDraft(localStorage, key.value, value).savedAt;
      saveError.value = "";
    } catch {
      saveError.value =
        "下書きを保存できません。入力内容をコピーして保管してください。";
    }
  }
});
const nextSubmittableAt = computed(() => {
  const metadata = data.value?.answers.metadata;
  const intervalEnd = metadata?.lastSubmittedAt
    ? Date.parse(metadata.lastSubmittedAt) +
      metadata.submitIntervalSeconds * 1000
    : 0;
  return Math.max(intervalEnd, retryAt.value) || undefined;
});
const cooldown = computed(() =>
  remainingCooldownSeconds(nextSubmittableAt.value, now.value),
);
const displayCooldown = computed(() =>
  problemCooldown.remainingSeconds(
    code,
    nextSubmittableAt.value,
    displayNow.value,
  ),
);
const cooldownMinutes = computed(() => Math.ceil(displayCooldown.value / 60));
const resultsHidden = computed(
  () => competition.value?.ranking.hidden ?? false,
);
const railItems = computed<ProblemRailItem[]>(() =>
  (competition.value?.problems ?? []).map((problem) => {
    const latest = activity.value?.find(
      (answer) => answer.problemCode === problem.code,
    );
    return {
      problem,
      cooldownMinutes:
        problem.code === code
          ? cooldownMinutes.value
          : problemCooldown.remainingMinutes(
              problem.code,
              problem.nextSubmittableAt,
              displayNow.value,
            ),
      pending: !resultsHidden.value && !!latest && latest.score == null,
      submitted: resultsHidden.value && !!latest,
    };
  }),
);
const open = computed(() => {
  const s = data.value?.problem.submissionStatus;
  return (
    !!s?.isSubmittable &&
    (!s.submittableUntil || now.value < Date.parse(s.submittableUntil))
  );
});
async function submit() {
  sending.value = true;
  actionError.value = undefined;
  submitted.value = false;
  try {
    await submitAnswer(api, code, body.value, attachments.value);
    attachments.value = [];
    submitted.value = true;
    await refresh();
    await refreshNuxtData(["competition", "activity"]);
  } catch (e) {
    actionError.value = e;
    if (e instanceof ApiError && e.status === 429)
      retryAt.value = now.value + (e.retryAfterSeconds ?? 0) * 1000;
  } finally {
    sending.value = false;
  }
}
let unsubscribe: (() => void) | undefined;
const streamError = ref(false);
function connect() {
  unsubscribe?.();
  streamError.value = false;
  unsubscribe = subscribeDeployments(
    code,
    (message) => {
      if (!data.value) return;
      streamError.value = false;
      data.value.deployments =
        message.type === "snapshot"
          ? message.deployments.map(mapDeployment)
          : [
              mapDeployment(message.deployment),
              ...data.value.deployments.filter(
                (d) => d.revision !== message.deployment.revision,
              ),
            ];
    },
    () => {
      streamError.value = true;
    },
  );
}
onMounted(connect);
onUnmounted(() => unsubscribe?.());
async function redeploy() {
  resetOpen.value = false;
  deploying.value = true;
  actionError.value = undefined;
  try {
    const deployment = await deploy(api, code);
    if (data.value)
      data.value.deployments = [
        data.value.deployments.find(
          (d) => d.revision === deployment.revision && d.status !== "QUEUED",
        ) ?? deployment,
        ...data.value.deployments.filter(
          (d) => d.revision !== deployment.revision,
        ),
      ];
  } catch (e) {
    actionError.value = e;
  } finally {
    deploying.value = false;
  }
}
const deploymentBusy = computed(
  () =>
    deploying.value ||
    data.value?.deployments.some((d) =>
      ["QUEUED", "DEPLOYING"].includes(d.status),
    ),
);
const score = computed(
  () => competition.value?.problems.find((p) => p.code === code)?.score,
);
</script>
<template>
  <main>
    <RequestState :pending="pending" :error="error" @retry="refresh"
      ><template v-if="data"
        ><div class="detail-layout">
          <ProblemRail :items="railItems" :current-code="code" />
          <article class="problem-content">
            <ProblemHero
              :code="code"
              :title="data.problem.title"
              :category="data.problem.category"
              :score="score"
              :max-score="data.problem.maxScore"
              :cooldown-minutes="cooldownMinutes"
              :open="open"
            />
            <div
              ref="contentGrid"
              class="content-grid"
              :style="
                answerWidth ? { '--answer-width': `${answerWidth}px` } : {}
              "
            >
              <AnswerResizer
                v-model:width="answerWidth"
                :container="contentGrid"
              />
              <section class="reading-section">
                <div class="reading-copy">
                  <MarkdownContent :source="data.problem.body" />
                </div>
              </section>
              <AnswerForm
                v-model:body="body"
                v-model:attachments="attachments"
                :open="open"
                :sending="sending"
                :submitted="submitted"
                :cooldown-seconds="cooldown"
                :saved-at="savedAt"
                :save-error="saveError"
                :action-error="actionError"
                :max-length="data.problem.answerMaxLength"
                :allow-attachments="data.problem.allowAttachments"
                @submit="submit"
              />
              <DeploymentHistory
                :deployments="data.deployments"
                :penalty-threshold="data.problem.penaltyThreashold"
                :deploying="deploying"
                :stream-error="streamError"
                :reset-disabled="!data.problem.redeployable || !!deploymentBusy"
                @reset="resetOpen = true"
              />
              <AnswerHistory
                :problem-code="code"
                :answers="data.answers.answers"
                :results-hidden="resultsHidden"
                :selected-answer-number="selectedAnswerNumber"
                :selected-answer-body="data.selectedAnswer?.answerBody"
                :selected-answer-attachments="data.selectedAnswer?.attachments"
              />
            </div>
          </article></div></template
    ></RequestState>
    <ResetDialog
      v-model:open="resetOpen"
      :penalty-threshold="data?.problem.penaltyThreashold"
      @confirm="redeploy"
    />
  </main>
</template>
