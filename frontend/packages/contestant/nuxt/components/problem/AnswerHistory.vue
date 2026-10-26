<script setup lang="ts">
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import type { Answer } from "~/features/answer";

defineProps<{
  problemCode: string;
  answers: readonly Answer[];
  selectedAnswerNumber: number | null;
  /** Body of the answer chosen with `?answer=`; omitted when none is selected. */
  selectedAnswerBody?: string;
  /** Results stay hidden until the final reveal (secondary qualifier mode). */
  resultsHidden?: boolean;
}>();
</script>
<template>
  <section
    id="answer-history"
    class="content-section"
    aria-labelledby="answer-history-heading"
  >
    <h2 id="answer-history-heading">提出履歴</h2>
    <p v-if="!answers.length">提出はありません。</p>
    <table v-else class="data-table">
      <thead>
        <tr>
          <th>回答</th>
          <th>提出日時</th>
          <th>得点</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="a in answers" :key="a.id">
          <td>
            <NuxtLink
              :to="{
                path: `/problems/${problemCode}`,
                query: { answer: String(a.id) },
                hash: '#answer-history',
              }"
              :aria-label="`回答 #${a.id} の内容を確認する`"
              :aria-current="selectedAnswerNumber === a.id ? 'true' : undefined"
            >
              <span class="numeric">#{{ a.id }}</span> を見る
            </NuxtLink>
          </td>
          <td class="numeric">
            {{ new Date(a.submittedAt).toLocaleString("ja-JP") }}
          </td>
          <td>
            <span class="numeric">{{
              a.score?.score ?? (resultsHidden ? "提出済み" : "採点中")
            }}</span>
          </td>
        </tr>
      </tbody>
    </table>
    <div
      v-if="selectedAnswerBody != null"
      class="reading-copy"
      aria-live="polite"
    >
      <MarkdownContent :source="selectedAnswerBody" />
    </div>
  </section>
</template>
