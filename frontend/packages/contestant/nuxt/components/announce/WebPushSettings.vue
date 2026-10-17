<script setup lang="ts">
defineProps<{
  /** This browser already receives notifications. */
  subscribed: boolean;
  /** The browser supports service workers and the Push API. */
  supported: boolean;
  /** Server-side Web Push is configured; undefined while loading. */
  serverEnabled?: boolean;
  busy: boolean;
  error?: unknown;
}>();
defineEmits<{ enable: []; disable: [] }>();
</script>
<template>
  <section
    class="notification-settings"
    aria-labelledby="notification-settings-title"
  >
    <div>
      <h2 id="notification-settings-title">Web Push通知</h2>
      <p v-if="subscribed">このブラウザで通知を受け取ります。</p>
      <p v-else-if="!supported">このブラウザはWeb Pushに対応していません。</p>
      <p v-else-if="serverEnabled === false">
        サーバー側のWeb Push設定が必要です。
      </p>
      <p v-else>新しいお知らせをブラウザの通知で受け取れます。</p>
    </div>
    <button
      v-if="subscribed"
      class="button-secondary"
      :disabled="busy"
      @click="$emit('disable')"
    >
      通知を停止する
    </button>
    <button
      v-else
      class="button-secondary"
      :disabled="busy || !supported || !serverEnabled"
      @click="$emit('enable')"
    >
      通知を有効にする
    </button>
    <p v-if="error" class="form-error" role="alert">
      {{ String(error) }}
    </p>
  </section>
</template>
