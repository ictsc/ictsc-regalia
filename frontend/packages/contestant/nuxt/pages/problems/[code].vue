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
import { problemStatus, statusLabels } from "~/features/problem/status";
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
const rail = ref(true),
  body = ref(""),
  savedAt = ref(""),
  saveError = ref(""),
  actionError = ref<unknown>(),
  sending = ref(false),
  deploying = ref(false),
  submitted = ref(false);
const now = useClock(),
  retryAt = ref(0);
const problemCooldown = useProblemCooldown();
const demoMode = useDemoMode();
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
const cooldown = computed(() => {
  const metadata = data.value?.answers.metadata;
  const intervalEnd = metadata?.lastSubmittedAt
    ? Date.parse(metadata.lastSubmittedAt) +
      metadata.submitIntervalSeconds * 1000
    : 0;
  const nextSubmittableAt = Math.max(intervalEnd, retryAt.value);
  return problemCooldown.remainingSeconds(
    code,
    nextSubmittableAt || undefined,
    now.value,
  );
});
const cooldownMinutes = computed(() => Math.ceil(cooldown.value / 60));
const cooldownSecondsFor = (
  problem: NonNullable<typeof competition.value>["problems"][number],
) => {
  const listed = problemCooldown.remainingSeconds(
    problem.code,
    problem.nextSubmittableAt,
    now.value,
  );
  return problem.code === code ? Math.max(listed, cooldown.value) : listed;
};
const cooldownMinutesFor = (
  problem: NonNullable<typeof competition.value>["problems"][number],
) =>
  problem.code === code
    ? Math.ceil(cooldownSecondsFor(problem) / 60)
    : problemCooldown.remainingMinutes(
        problem.code,
        problem.nextSubmittableAt,
        now.value,
      );
const statusOf = (
  problem: NonNullable<typeof competition.value>["problems"][number],
) => {
  const latest = activity.value?.find(
    (answer) => answer.problemCode === problem.code,
  );
  return problemStatus(
    problem,
    !!latest,
    !!latest && latest.score == null,
    demoMode.value,
  );
};
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
    await submitAnswer(api, code, body.value);
    submitted.value = true;
    await refresh();
    const lastSubmittedAt = data.value?.answers.metadata.lastSubmittedAt;
    setDemoClock(lastSubmittedAt ?? Date.now());
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
  if (
    !window.confirm(
      "環境をリセットします。再展開ルールに応じて減点される場合があります。続行しますか？",
    )
  )
    return;
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
        ><div class="actions">
          <NuxtLink to="/problems">← 問題一覧へ</NuxtLink
          ><button
            class="button-secondary"
            :aria-expanded="rail"
            aria-controls="problem-rail"
            @click="rail = !rail"
          >
            {{ rail ? "問題リンクを隠す" : "問題リンクを表示" }}
          </button>
        </div>
        <div class="detail-layout" :class="{ 'is-rail-hidden': !rail }">
          <nav
            v-if="rail"
            id="problem-rail"
            class="problem-rail"
            aria-label="問題ナビゲーション"
          >
            <template v-for="p in competition?.problems" :key="p.code"
              ><NuxtLink
                :to="`/problems/${p.code}`"
                :class="[
                  `status-${statusOf(p)}`,
                  {
                    'is-active': p.code === code,
                    'is-cooldown': cooldownSecondsFor(p) > 0,
                    'is-answer-closed': !p.submissionStatus?.isSubmittable,
                  },
                ]"
                :data-cooldown="
                  cooldownSecondsFor(p)
                    ? `${cooldownMinutesFor(p)}分`
                    : undefined
                "
                :aria-current="p.code === code ? 'page' : undefined"
                :aria-label="`${p.code} ${p.title} ${statusLabels[statusOf(p)]}${cooldownSecondsFor(p) ? ` 再提出可能まで${cooldownMinutesFor(p)}分` : ''}`"
                ><b>{{ p.code }}</b
                ><span class="problem-preview"
                  ><strong>{{ p.title }}</strong
                  ><small>{{ p.category }}</small
                  ><span class="preview-status" :class="`is-${statusOf(p)}`">{{
                    statusLabels[statusOf(p)]
                  }}</span
                  ><span class="preview-score"
                    >{{ p.score ? `${p.score.score}点` : "未確定" }} /
                    {{ p.maxScore }}点満点</span
                  ><span v-if="cooldownSecondsFor(p)" class="preview-cooldown"
                    >再提出可能まで <b>{{ cooldownMinutesFor(p) }}</b
                    >分</span
                  ><span
                    v-else-if="!p.submissionStatus?.isSubmittable"
                    class="preview-closed"
                    >回答終了 / 閲覧のみ</span
                  ></span
                ></NuxtLink
              ></template
            >
          </nav>
          <article class="problem-content">
            <header class="problem-hero">
              <div class="hero-copy">
                <h1>{{ code }}:{{ data.problem.title }}.</h1>
                <p>{{ data.problem.category }}</p>
              </div>
              <div class="problem-summary-side">
                <p class="problem-score-summary">
                  <strong>{{ score ? `${score.score}点` : "—" }}</strong
                  ><span>/ {{ data.problem.maxScore }}点満点</span>
                </p>
                <p
                  class="problem-cooldown"
                  :class="{ 'is-answer-closed': !open && !cooldown }"
                >
                  <template v-if="cooldown"
                    ><span>再回答可能まで</span
                    ><strong>{{ cooldownMinutes }}</strong
                    ><small>分</small></template
                  ><span v-else>{{ open ? "回答受付中" : "回答終了" }}</span>
                </p>
              </div>
            </header>
            <div class="content-grid">
              <section class="reading-section">
                <div class="reading-copy">
                  <MarkdownContent :source="data.problem.body" />
                </div>
              </section>
              <section class="answer-section">
                <form @submit.prevent="submit">
                  <div class="answer-heading-row">
                    <h2>回答</h2>
                    <p v-if="saveError" role="alert">{{ saveError }}</p>
                    <p v-else-if="savedAt">
                      このブラウザに保存済み
                      <time>{{
                        new Date(savedAt).toLocaleTimeString("ja-JP")
                      }}</time>
                    </p>
                  </div>
                  <label class="visually-hidden" for="answer">回答</label
                  ><textarea
                    id="answer"
                    v-model="body"
                    rows="10"
                    required
                    :disabled="!open || sending"
                    placeholder="原因 / 実施した変更 / 復旧確認の結果を記入してください"
                  />
                  <p v-if="submitted" role="status">回答を提出しました</p>
                  <p v-if="cooldown" role="status">
                    再回答可能まで {{ Math.floor(cooldown / 60) }}分{{
                      cooldown % 60
                    }}秒
                  </p>
                  <div class="answer-actions">
                    <button
                      class="button-primary"
                      :disabled="
                        !open || cooldown > 0 || sending || !body.trim()
                      "
                    >
                      {{ sending ? "提出中…" : "回答を提出" }} ↗</button
                    ><button
                      class="button-secondary"
                      type="button"
                      :disabled="!data.problem.redeployable || deploymentBusy"
                      @click="redeploy"
                    >
                      環境をリセット
                    </button>
                  </div>
                </form>
                <div
                  v-if="actionError"
                  role="alert"
                  class="state-message error-message"
                >
                  {{
                    actionError instanceof Error
                      ? actionError.message
                      : String(actionError)
                  }}
                </div>
              </section>
              <section class="content-section">
                <h2>再展開</h2>
                <p>減点なしの上限: {{ data.problem.penaltyThreashold }}回</p>
                <p v-if="deploying" role="status">QUEUED — 再展開を要求中</p>
                <p v-if="streamError" role="alert">
                  更新の接続を確認しています
                </p>
                <p v-if="!data.deployments.length">再展開履歴はありません。</p>
                <ul>
                  <li v-for="d in data.deployments" :key="d.revision">
                    #{{ d.revision }} {{ d.status }} / 減点 {{ d.penalty }} /
                    残り {{ d.allowedDeploymentCount }}回
                  </li>
                </ul>
              </section>
              <section
                id="answer-history"
                class="content-section"
                aria-labelledby="answer-history-heading"
              >
                <h2 id="answer-history-heading">提出履歴</h2>
                <p v-if="!data.answers.answers.length">提出はありません。</p>
                <table v-else class="data-table">
                  <thead>
                    <tr>
                      <th>回答</th>
                      <th>提出日時</th>
                      <th>得点</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="a in data.answers.answers" :key="a.id">
                      <td>
                        <NuxtLink
                          :to="{
                            path: `/problems/${code}`,
                            query: { answer: String(a.id) },
                            hash: '#answer-history',
                          }"
                          :aria-label="`回答 #${a.id} の内容を確認する`"
                          :aria-current="
                            selectedAnswerNumber === a.id ? 'true' : undefined
                          "
                        >
                          #{{ a.id }} を見る
                        </NuxtLink>
                      </td>
                      <td>
                        {{ new Date(a.submittedAt).toLocaleString("ja-JP") }}
                      </td>
                      <td>{{ a.score?.score ?? "採点中" }}</td>
                    </tr>
                  </tbody>
                </table>
                <div
                  v-if="data.selectedAnswer"
                  class="reading-copy"
                  aria-live="polite"
                >
                  <MarkdownContent :source="data.selectedAnswer.answerBody" />
                </div>
              </section>
            </div>
          </article></div></template
    ></RequestState>
  </main>
</template>
