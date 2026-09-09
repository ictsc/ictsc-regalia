<script setup lang="ts">
import { api } from "@ictsc/api";
import { fetchRule } from "~/features/rule";
useHead({ title: "ルール" });
const { data, error, pending, refresh } = await useAsyncData("rule", () =>
  fetchRule(api),
);
</script>
<template>
  <main class="workspace">
    <h1>ルール</h1>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><MarkdownContent v-if="data" :source="data.markdown"
    /></RequestState>
  </main>
</template>
