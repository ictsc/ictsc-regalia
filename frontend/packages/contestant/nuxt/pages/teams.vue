<script setup lang="ts">
import { api } from "@ictsc/api";
import { fetchTeams } from "~/features/teams";
useHead({ title: "チーム" });
const { data, error, pending, refresh } = await useAsyncData("teams", () =>
  fetchTeams(api),
);
</script>
<template>
  <main class="workspace">
    <h1>チーム</h1>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><section v-for="team in data" :key="team.code" class="content-section">
        <h2>{{ team.name }}</h2>
        <p>{{ team.organization }}</p>
        <article v-for="member in team.members" :key="member.name">
          <h3>
            {{ member.displayName }} <small>@{{ member.name }}</small>
          </h3>
          <MarkdownContent :source="member.selfIntroduction" />
        </article></section
    ></RequestState>
  </main>
</template>
