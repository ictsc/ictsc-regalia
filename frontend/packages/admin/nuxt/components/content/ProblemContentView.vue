<script setup lang="ts">
import { ref } from "vue";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";
import { adminImageUrl } from "~/features/admin-images";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import AdminPage from "~/components/layout/AdminPage.vue";
import AdminPanel from "~/components/layout/AdminPanel.vue";
import AdminFacts from "~/components/layout/AdminFacts.vue";
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
  <AdminPage
    :title="problem ? `${problem.code}: ${problem.title}` : '問題'"
    lead="指定した版の問題本文と再展開設定を確認します。"
    :back="{ to: '/content', label: 'コンテンツ一覧' }"
  >
    <AdminPanel title="表示する版">
      <form class="form-grid" @submit.prevent="$emit('viewCommit', commit)">
        <label
          >参照 commit（空欄で配信中の版）<input
            v-model="commit"
            class="admin-code"
            autocomplete="off"
            spellcheck="false" /></label
        ><button class="button-primary">表示</button>
      </form>
    </AdminPanel>
    <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
      ><template v-if="problem"
        ><AdminPanel title="問題情報">
          <AdminFacts
            :items="[
              { label: 'カテゴリ', value: problem.category },
              { label: 'セクション', value: problem.section_slug },
              { label: '満点', value: problem.max_score },
              { label: 'commit', value: problem.content_commit, mono: true },
              {
                label: 'デフォルト採点担当者',
                value: problem.default_grader_discord_id ?? '未設定',
                mono: !!problem.default_grader_discord_id,
              },
              {
                label: '回答の文字数上限',
                value: `${problem.answer_max_length.toLocaleString()} 文字（改行を除く）`,
              },
              {
                label: 'ログ添付',
                value: problem.allow_attachments ? '可' : '不可',
              },
              {
                label: '再展開ルール',
                value: `${problem.redeploy_rule.type} / 減点開始回数 ${problem.redeploy_rule.penalty_threshold ?? '—'} / ${problem.redeploy_rule.penalty_percentage ?? '—'}%`,
              },
            ]"
          />
        </AdminPanel>
        <AdminPanel title="問題本文">
          <MarkdownContent
            :source="problem.body"
            :image-url="adminImageUrl(problem.content_commit)"
            diagrams
          />
        </AdminPanel>
        <AdminPanel title="問題解説">
          <MarkdownContent
            :source="problem.explanation"
            :image-url="adminImageUrl(problem.content_commit)"
            diagrams /></AdminPanel></template
    ></RequestState>
  </AdminPage>
</template>
