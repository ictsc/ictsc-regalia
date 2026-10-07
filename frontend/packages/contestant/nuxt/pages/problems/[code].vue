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
import { groupProblems } from "~/features/problem/group";
import { remainingCooldownSeconds } from "~/features/problem/cooldown";
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
const rail = useState<boolean>("problem-rail-open", () => true),
  problemRail = ref<HTMLElement | null>(null),
  body = ref(""),
  savedAt = ref(""),
  saveError = ref(""),
  actionError = ref<unknown>(),
  sending = ref(false),
  deploying = ref(false),
  submitted = ref(false),
  resetDialog = ref<HTMLDialogElement | null>(null);
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
const cooldownSecondsFor = (
  problem: NonNullable<typeof competition.value>["problems"][number],
) => {
  const listed = problemCooldown.remainingSeconds(
    problem.code,
    problem.nextSubmittableAt,
    displayNow.value,
  );
  return problem.code === code ? displayCooldown.value : listed;
};
const cooldownMinutesFor = (
  problem: NonNullable<typeof competition.value>["problems"][number],
) =>
  problem.code === code
    ? Math.ceil(cooldownSecondsFor(problem) / 60)
    : problemCooldown.remainingMinutes(
        problem.code,
        problem.nextSubmittableAt,
        displayNow.value,
      );
const latestAnswerOf = (
  problem: NonNullable<typeof competition.value>["problems"][number],
) => activity.value?.find((answer) => answer.problemCode === problem.code);
const pendingOf = (
  problem: NonNullable<typeof competition.value>["problems"][number],
) => {
  const latest = latestAnswerOf(problem);
  return !!latest && latest.score == null;
};
function openResetDialog() {
  resetDialog.value?.showModal();
}
function onResetDialogBackdrop(event: MouseEvent) {
  if (event.target === resetDialog.value) resetDialog.value?.close();
}
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
  resetDialog.value?.close();
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
const problemGroups = computed(() =>
  groupProblems(competition.value?.problems ?? []),
);
const sectionName = (slug?: string) =>
  slug === "both"
    ? "両日"
    : /^day\d+$/.test(slug ?? "")
      ? `${slug!.slice(3)}日目`
      : slug;
async function showCurrentRailItem() {
  await nextTick();
  if (!import.meta.client || !window.matchMedia("(max-width: 760px)").matches)
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
watch([problemGroups, rail], showCurrentRailItem, { flush: "post" });
</script>
<template>
  <main>
    <RequestState :pending="pending" :error="error" @retry="refresh"
      ><template v-if="data"
        ><div class="actions detail-toolbar">
          <button
            class="problem-rail-toggle"
            :aria-expanded="rail"
            :aria-label="rail ? '問題リンクを隠す' : '問題リンクを表示'"
            aria-controls="problem-rail"
            @click="rail = !rail"
          >
            <span class="problem-rail-toggle-icon" aria-hidden="true" />
            <span class="visually-hidden">{{
              rail ? "問題リンクを隠す" : "問題リンクを表示"
            }}</span>
          </button>
          <button
            class="button-secondary button-reset detail-reset"
            type="button"
            :disabled="!data.problem.redeployable || deploymentBusy"
            @click="openResetDialog"
          >
            環境をリセット
          </button>
        </div>
        <div class="detail-layout" :class="{ 'is-rail-hidden': !rail }">
          <nav
            v-if="rail"
            id="problem-rail"
            ref="problemRail"
            class="problem-rail"
            aria-label="問題ナビゲーション"
          >
            <NuxtLink class="problem-rail-label" to="/problems"
              >問題一覧</NuxtLink
            >
            <section
              v-for="group in problemGroups"
              :key="group.key"
              class="problem-rail-group"
            >
              <h2>
                <span class="visually-hidden">ステージ: </span>
                <span class="rail-stage-label">{{
                  sectionName(group.key)
                }}</span>
              </h2>
              <NuxtLink
                v-for="p in group.problems"
                :key="p.code"
                :to="`/problems/${p.code}`"
                :class="{
                  'is-active': p.code === code,
                  'is-answer-closed': !p.submissionStatus?.isSubmittable,
                }"
                :aria-current="p.code === code ? 'page' : undefined"
                :aria-label="`${p.code} ${p.title} ${p.score ? `${p.score.score}点 / ${p.maxScore}点満点` : '得点未確定'}${pendingOf(p) ? ' 最新の提出を採点中' : ''}${cooldownSecondsFor(p) ? ` 再提出可能まで${cooldownMinutesFor(p)}分` : ''}`"
                ><b>{{ p.code }}</b
                ><span class="rail-score">
                  <strong>{{ p.score ? p.score.score : "—" }}</strong>
                  <small>/ {{ p.maxScore }}</small>
                  <em v-if="pendingOf(p)" class="rail-pending">採点中</em></span
                ><span class="rail-title">{{ p.title }}</span
                ><span class="problem-preview"
                  ><strong>{{ p.title }}</strong
                  ><small>{{ p.category }}</small
                  ><span class="preview-status">{{
                    pendingOf(p)
                      ? "最新の提出を採点中"
                      : p.score
                        ? "採点済み"
                        : "未回答"
                  }}</span
                  ><span class="preview-score"
                    >{{ p.score ? `${p.score.score}点` : "—" }} /
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
              >
            </section>
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
                  :class="{ 'is-answer-closed': !open && !displayCooldown }"
                >
                  <template v-if="displayCooldown"
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
                    <span v-else-if="savedAt" class="draft-save-status">
                      <button
                        class="draft-save-icon"
                        type="button"
                        :aria-label="`このブラウザに保存済み ${new Date(savedAt).toLocaleTimeString('ja-JP')}`"
                      >
                        <svg
                          aria-hidden="true"
                          viewBox="0 0 24 24"
                          fill="none"
                          xmlns="http://www.w3.org/2000/svg"
                        >
                          <rect x="2.5" y="3.5" width="19" height="17" rx="2" />
                          <path
                            d="M2.5 8h19M7 5.75h.01M10 5.75h.01M8 14l2.5 2.5L16 11"
                          />
                        </svg>
                      </button>
                      <span class="draft-save-tooltip" role="tooltip">
                        このブラウザに保存済み
                        <time :datetime="savedAt">{{
                          new Date(savedAt).toLocaleTimeString("ja-JP")
                        }}</time>
                      </span>
                    </span>
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
                      class="button-primary button-submit"
                      :disabled="
                        !open || cooldown > 0 || sending || !body.trim()
                      "
                    >
                      {{ sending ? "提出中…" : "回答を提出" }}
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
                          <span class="numeric">#{{ a.id }}</span> を見る
                        </NuxtLink>
                      </td>
                      <td class="numeric">
                        {{ new Date(a.submittedAt).toLocaleString("ja-JP") }}
                      </td>
                      <td>
                        <span class="numeric">{{
                          a.score?.score ?? "採点中"
                        }}</span>
                      </td>
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
    <dialog
      ref="resetDialog"
      class="reset-dialog confirmation-dialog"
      aria-labelledby="reset-dialog-title"
      aria-describedby="reset-dialog-message"
      @click="onResetDialogBackdrop"
    >
      <p class="confirmation-dialog-eyebrow">操作の確認</p>
      <h2 id="reset-dialog-title">環境をリセットしますか？</h2>
      <p id="reset-dialog-message" class="confirmation-dialog-message">
        この問題の環境を再展開します。再展開ルールに応じて減点される場合があります。
      </p>
      <p v-if="data" class="confirmation-dialog-note">
        減点なしの上限: {{ data.problem.penaltyThreashold }}回
      </p>
      <div class="confirmation-dialog-actions">
        <button
          type="button"
          class="button-secondary"
          autofocus
          @click="resetDialog?.close()"
        >
          キャンセル
        </button>
        <button type="button" class="button-primary" @click="redeploy">
          環境をリセットする
        </button>
      </div>
    </dialog>
  </main>
</template>
