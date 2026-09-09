<script setup lang="ts">
import { api } from "@ictsc/api";
import {
  getAdminAnswer,
  listAdminMarkingResults,
  createAdminMarkingResult,
  createAdminActions,
} from "~/features/admin-api";
definePageMeta({ key: (route) => route.fullPath });
const route = useRoute(),
  reference = {
    problem_code: String(route.params.problem),
    team_code: Number(route.params.team),
    answer_number: Number(route.params.id),
  };
const { data, error, pending, refresh } = await useAsyncData(
  `admin-answer:${route.fullPath}`,
  async () => {
    const answer = await getAdminAnswer(api, reference);
    const [problem, active, marks] = await Promise.all([
      createAdminActions(api).getProblem(
        reference.problem_code,
        answer.content_commit,
      ),
      createAdminActions(api).getProblem(reference.problem_code),
      listAdminMarkingResults(api),
    ]);
    return {
      answer,
      problem: problem.problem,
      active: active.problem,
      marks: marks.filter(
        (m) =>
          m.answer.problem_code === reference.problem_code &&
          m.answer.team_code === reference.team_code &&
          m.answer.answer_number === reference.answer_number,
      ),
    };
  },
);
useHead({
  title: `${reference.problem_code} / Team ${reference.team_code} / #${reference.answer_number}`,
});
const score = ref(0),
  comment = ref("");
const { busy, message, run } = useMutation(refresh);
async function mark() {
  await run(
    "採点結果を保存",
    () => createAdminMarkingResult(api, reference, score.value, comment.value),
    `得点: ${score.value}\nコメント: ${comment.value}\nこの採点結果を送信しますか？`,
  );
}
</script>
<template>
  <main class="workspace">
    <h1>
      {{ reference.problem_code }} / Team {{ reference.team_code }} / #{{
        reference.answer_number
      }}
    </h1>
    <p role="status">{{ message }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><template v-if="data"
        ><p>
          チーム: {{ data.answer.team.name }} / 回答者:
          {{ data.answer.author.display_name }}
        </p>
        <p>
          提出commit: <span class="mono">{{ data.answer.content_commit }}</span>
        </p>
        <p v-if="data.answer.content_commit !== data.active.content_commit">
          現行commitと異なります
        </p>
        <p>現行commit: {{ data.active.content_commit }}</p>
        <p>
          現在の採用得点: {{ data.answer.score?.total ?? "未採点" }} /
          {{ data.answer.problem.max_score }}
        </p>
        <h2>解答</h2>
        <MarkdownContent :source="data.answer.body.body" />
        <details>
          <summary>プレーンテキスト版</summary>
          <pre>{{ data.answer.body.body }}</pre>
        </details>
        <h2>提出時点の問題</h2>
        <MarkdownContent :source="data.problem.body" />
        <h2>問題解説</h2>
        <MarkdownContent :source="data.problem.explanation" />
        <h2>採点</h2>
        <form class="form-grid" @submit.prevent="mark">
          <label
            >得点<input
              v-model.number="score"
              type="number"
              min="0"
              :max="data.answer.problem.max_score"
              required /></label
          ><label>コメント<textarea v-model="comment" /></label
          ><button class="button-primary" :disabled="busy">送信</button>
        </form>
        <h2>採点履歴</h2>
        <table class="data-table">
          <thead>
            <tr>
              <th>日時</th>
              <th>得点</th>
              <th>コメント</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="m in data.marks" :key="m.id">
              <td>{{ m.created_at }}</td>
              <td>{{ m.score }}</td>
              <td>{{ m.rationale }}</td>
            </tr>
          </tbody>
        </table></template
      ></RequestState
    >
  </main>
</template>
