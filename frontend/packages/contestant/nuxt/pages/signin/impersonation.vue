<script setup lang="ts">
import { api } from "@ictsc/api";
import {
  listImpersonationCandidates,
  startImpersonation,
} from "~/features/viewer/impersonation";
useHead({ title: "代理ログイン" });
const { data, error, pending, refresh } = await useAsyncData(
  "impersonations",
  () => listImpersonationCandidates(api),
);
const message = ref("");
async function start(candidate: NonNullable<typeof data.value>[number]) {
  try {
    await startImpersonation(candidate);
    clearNuxtData();
    await useSession().refresh();
    await navigateTo("/problems");
  } catch (e) {
    message.value = e instanceof Error ? e.message : String(e);
  }
}
</script>
<template>
  <main class="workspace">
    <h1>代理ログイン</h1>
    <p role="alert">{{ message }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><table class="data-table">
        <tbody>
          <tr v-for="c in data" :key="c.name">
            <td>{{ c.teamName }}</td>
            <td>{{ c.displayName }}</td>
            <td>
              <button class="button-secondary" @click="start(c)">
                代理ログイン
              </button>
            </td>
          </tr>
        </tbody>
      </table></RequestState
    >
  </main>
</template>
