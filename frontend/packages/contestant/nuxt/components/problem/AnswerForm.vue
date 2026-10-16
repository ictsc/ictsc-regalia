<script setup lang="ts">
const body = defineModel<string>("body", { required: true });
defineProps<{
  /** Answers are being accepted now. */
  open: boolean;
  sending: boolean;
  submitted: boolean;
  /** Seconds until the team can answer again; 0 when not cooling down. */
  cooldownSeconds: number;
  /** ISO time of the last draft save in this browser. */
  savedAt: string;
  saveError: string;
  actionError?: unknown;
}>();
defineEmits<{ submit: [] }>();
</script>
<template>
  <section class="answer-section">
    <form @submit.prevent="$emit('submit')">
      <div class="answer-heading-row">
        <h2>回答</h2>
        <p v-if="saveError" role="alert">{{ saveError }}</p>
        <span v-else-if="savedAt" class="draft-save-status">
          <button
            class="draft-save-icon"
            type="button"
            :aria-label="`このブラウザに保存済み ${new Date(savedAt).toLocaleTimeString('ja-JP')}`"
          >
            <svg
              aria-hidden="true"
              viewBox="0 0 24 24"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <rect x="2.5" y="3.5" width="19" height="17" rx="2" />
              <path d="M2.5 8h19M7 5.75h.01M10 5.75h.01M8 14l2.5 2.5L16 11" />
            </svg>
          </button>
          <span class="draft-save-tooltip" role="tooltip">
            このブラウザに保存済み
            <time :datetime="savedAt">{{
              new Date(savedAt).toLocaleTimeString("ja-JP")
            }}</time>
          </span>
        </span>
      </div>
      <label class="visually-hidden" for="answer">回答</label
      ><textarea
        id="answer"
        v-model="body"
        rows="10"
        required
        :disabled="!open || sending"
        placeholder="原因 / 実施した変更 / 復旧確認の結果を記入してください"
      />
      <p v-if="submitted" role="status">回答を提出しました</p>
      <p v-if="cooldownSeconds" role="status">
        再回答可能まで {{ Math.floor(cooldownSeconds / 60) }}分{{
          cooldownSeconds % 60
        }}秒
      </p>
      <div class="answer-actions">
        <button
          class="button-primary button-submit"
          :disabled="!open || cooldownSeconds > 0 || sending || !body.trim()"
        >
          {{ sending ? "提出中…" : "回答を提出" }}
        </button>
      </div>
    </form>
    <div v-if="actionError" role="alert" class="state-message error-message">
      {{
        actionError instanceof Error ? actionError.message : String(actionError)
      }}
    </div>
  </section>
</template>
