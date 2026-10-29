<script setup lang="ts">
import { ref } from "vue";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import { adminImageUrl } from "~/features/admin-images";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { components } from "@ictsc/api";

const props = defineProps<{
  problem?: components["schemas"]["AdminProblem"];
  /** Commit from the URL; empty means the active version. */
  initialCommit: string;
  pending?: boolean;
  error?: unknown;
}>();
defineEmits<{ viewCommit: [commit: string]; retry: [] }>();

const commit = ref(props.initialCommit);
</script>
<template>
  <NuxtLink class="text-link" to="/content">← コンテンツ一覧</NuxtLink>
  <h1>{{ problem?.code }}: {{ problem?.title }}</h1>
  <p class="admin-lead">指定した版の問題本文と再展開設定を確認します。</p>
  <section class="admin-panel">
    <h2>表示する版</h2>
    <form class="form-grid" @submit.prevent="$emit('viewCommit', commit)">
      <label>参照commit（空欄でactive）<input v-model="commit" /></label
      ><button class="button-primary">表示</button>
    </form>
  </section>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><template v-if="problem"
      ><section class="admin-panel">
        <h2>問題情報</h2>
        <p>
          {{ problem.category }} / {{ problem.section_slug }} / 満点
          {{ problem.max_score }}
        </p>
        <p class="mono">{{ problem.content_commit }}</p>
        <p>
          デフォルト採点担当者:
          {{ problem.default_grader_discord_id ?? "未設定" }}
        </p>
        <p>
          回答の文字数上限: {{ problem.answer_max_length.toLocaleString() }}
          文字（改行を除く） / ログ添付:
          {{ problem.allow_attachments ? "可" : "不可" }}
        </p>
        <p>
          再展開ルール: {{ problem.redeploy_rule.type }} / 減点開始回数
          {{ problem.redeploy_rule.penalty_threshold ?? "—" }} /
          {{ problem.redeploy_rule.penalty_percentage ?? "—" }}%
        </p>
      </section>
      <section class="admin-panel">
        <h2>問題本文</h2>
        <MarkdownContent
          :source="problem.body"
          :image-url="adminImageUrl(problem.content_commit)"
          diagrams
        />
      </section>
      <section class="admin-panel">
        <h2>問題解説</h2>
        <MarkdownContent
          :source="problem.explanation"
          :image-url="adminImageUrl(problem.content_commit)"
          diagrams
        /></section></template
  ></RequestState>
</template>
