<script setup lang="ts">
import { api } from "@ictsc/api";
import { fetchActivity } from "~/features/activity";
useHead({ title: "提出履歴" });
const { data, error, pending, refresh } = await useAsyncData("activity", () =>
  fetchActivity(api),
);
</script>
<template>
  <main>
    <header class="linked-page-heading">
      <h1 class="page-title">提出履歴</h1>
      <p>チームの回答と採点状況</p>
    </header>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><div class="table-scroll">
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
            <tr v-for="a in data" :key="`${a.problemCode}:${a.answerId}`">
              <td class="numeric">
                {{ new Date(a.submittedAt).toLocaleString("ja-JP") }}
              </td>
              <td>
                <NuxtLink :to="`/problems/${a.problemCode}`"
                  ><span class="numeric">{{ a.problemCode }}</span>
                  {{ a.problemTitle }}</NuxtLink
                >
              </td>
              <td class="numeric">#{{ a.answerId }}</td>
              <td>
                <span class="numeric">{{ a.score?.score ?? "採点中" }}</span>
                / <span class="numeric">{{ a.maxScore }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <p v-if="!data?.length">提出はありません。</p></RequestState
    >
  </main>
</template>
