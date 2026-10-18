<script setup lang="ts">
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { ActivityEntry } from "~/features/activity";

defineProps<{
  entries?: ActivityEntry[];
  pending?: boolean;
  error?: unknown;
}>();
defineEmits<{ retry: [] }>();
</script>
<template>
  <header class="linked-page-heading">
    <h1 class="page-title">提出履歴</h1>
    <p>チームの回答と採点状況</p>
  </header>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><div class="table-scroll activity-table">
      <table class="data-table">
        <thead>
          <tr>
            <th>提出日時</th>
            <th>問題</th>
            <th>回答</th>
            <th>得点</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="a in entries" :key="`${a.problemCode}:${a.answerId}`">
            <td class="numeric" data-label="提出日時">
              {{ new Date(a.submittedAt).toLocaleString("ja-JP") }}
            </td>
            <td data-label="問題">
              <span class="numeric">{{ a.problemCode }}</span>
              {{ a.problemTitle }}
            </td>
            <td data-label="回答">
              <NuxtLink
                :to="{
                  path: `/problems/${a.problemCode}`,
                  query: { answer: String(a.answerId) },
                  hash: '#answer-history',
                }"
                :aria-label="`${a.problemCode}の回答 #${a.answerId} を見る`"
              >
                <span class="numeric">回答 #{{ a.answerId }}</span> を見る
              </NuxtLink>
            </td>
            <td data-label="得点">
              <span class="numeric">{{ a.score?.score ?? "採点中" }}</span>
              / <span class="numeric">{{ a.maxScore }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p v-if="!entries?.length">提出はありません。</p></RequestState
  >
</template>
