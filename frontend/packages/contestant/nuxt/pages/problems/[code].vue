<script setup lang="ts">
import { api, ApiError } from "@ictsc/api";
import { fetchProblem } from "~/features/problem";
import { fetchAnswers, fetchAnswer, submitAnswer } from "~/features/answer";
import {
  fetchDeployments,
  deploy,
  subscribeDeployments,
  mapDeployment,
} from "~/features/deployment";
import { draftKey, loadDraft, saveDraft } from "~/features/draft";
import { problemStatus } from "~/features/problem/status";
definePageMeta({ key: (route) => route.fullPath });
const code = String(useRoute().params.code);
const { viewer } = useSession();
const { data: competition } = await useCompetition();
const { data, error, pending, refresh } = await useAsyncData(
  `problem:${code}`,
  async () => {
    const [problem, answers, deployments] = await Promise.all([
      fetchProblem(api, code),
      fetchAnswers(api, code),
      fetchDeployments(api, code),
    ]);
    return { problem, answers, deployments };
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
  submitted = ref(false),
  selectedAnswer = ref("");
const now = useClock(),
  retryAt = ref(0);
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
  return Math.max(
    0,
    Math.ceil((Math.max(intervalEnd, retryAt.value) - now.value) / 1000),
  );
});
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
      retryAt.value = Date.now() + (e.retryAfterSeconds ?? 0) * 1000;
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
async function viewAnswer(id: number) {
  try {
    selectedAnswer.value = (await fetchAnswer(api, code, id)).answerBody;
  } catch (e) {
    actionError.value = e;
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
                  `status-${problemStatus(p)}`,
                  { 'is-active': p.code === code },
                ]"
                :aria-current="p.code === code ? 'page' : undefined"
                :aria-label="`${p.code} ${p.title}`"
                ><b>{{ p.code }}</b
                ><span class="problem-preview"
                  ><strong>{{ p.title }}</strong
                  ><small>{{ p.category }}</small
                  ><span
                    >{{ p.score?.score ?? "—" }} / {{ p.maxScore }}点</span
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
              <div class="problem-score-summary">
                <strong>{{ score?.score ?? "—" }}</strong
                ><span>/ {{ data.problem.maxScore }}点満点</span>
              </div>
              <div class="problem-cooldown">
                <template v-if="cooldown"
                  ><span>再提出まで</span
                  ><strong>{{ Math.ceil(cooldown / 60) }}</strong
                  ><small>分</small></template
                ><span v-else>{{
                  open ? "回答受付中" : "受付時間外 / 閲覧のみ"
                }}</span>
              </div>
            </header>
            <div class="content-grid">
              <section class="reading-section">
                <h2>問題</h2>
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
                    再提出まで {{ Math.floor(cooldown / 60) }}分{{
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
              <section class="content-section">
                <h2>提出履歴</h2>
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
                        <button
                          class="button-secondary"
                          @click="viewAnswer(a.id)"
                        >
                          #{{ a.id }}
                        </button>
                      </td>
                      <td>
                        {{ new Date(a.submittedAt).toLocaleString("ja-JP") }}
                      </td>
                      <td>{{ a.score?.score ?? "採点中" }}</td>
                    </tr>
                  </tbody>
                </table>
                <MarkdownContent
                  v-if="selectedAnswer"
                  :source="selectedAnswer"
                />
              </section>
            </div>
          </article></div></template
    ></RequestState>
  </main>
</template>
