<script setup lang="ts">
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { TeamProfile } from "~/features/models";

defineProps<{
  teams?: TeamProfile[];
  pending?: boolean;
  error?: unknown;
}>();
defineEmits<{ retry: [] }>();
</script>
<template>
  <h1 class="page-title">チーム</h1>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><section v-for="team in teams" :key="team.code" class="content-section">
      <h2>{{ team.name }}</h2>
      <p>{{ team.organization }}</p>
      <article v-for="member in team.members" :key="member.name">
        <h3>
          {{ member.displayName }} <small>@{{ member.name }}</small>
        </h3>
        <MarkdownContent :source="member.selfIntroduction" />
      </article></section
  ></RequestState>
</template>
