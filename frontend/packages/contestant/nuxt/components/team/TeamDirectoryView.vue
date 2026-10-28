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
      <header class="team-heading">
        <h2>{{ team.name }}</h2>
        <p>{{ team.organization }}</p>
      </header>
      <div class="team-members">
        <article
          v-for="member in team.members"
          :key="member.name"
          class="team-member"
        >
          <h3>
            {{ member.displayName }} <small>@{{ member.name }}</small>
          </h3>
          <MarkdownContent
            v-if="member.selfIntroduction"
            :source="member.selfIntroduction"
            block-external-images
          />
        </article>
      </div></section
  ></RequestState>
</template>
<style scoped>
.team-heading {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 4px 16px;
  margin-bottom: 16px;
}
.team-heading h2,
.team-heading p {
  margin: 0;
}
.team-heading p {
  color: var(--muted);
  font-size: 13px;
}
/* Each member is one block: name, handle, then their introduction. The
   gap alone keeps the columns apart. */
.team-members {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 24px 32px;
}
.team-member h3 {
  margin: 0;
  font-size: 15px;
}
.team-member small {
  color: var(--muted);
  font-family: var(--mono);
  font-size: 12px;
  font-weight: 400;
}
.team-member :deep(.markdown) {
  margin-top: 6px;
  font-size: 13px;
  line-height: 1.7;
}
.team-member :deep(.markdown > :first-child) {
  margin-top: 0;
}
.team-member :deep(.markdown > :last-child) {
  margin-bottom: 0;
}
</style>
