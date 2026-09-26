<script setup lang="ts">
import { api } from "@ictsc/api";
import { fetchNotices } from "~/features/announce";
import { useReadAnnouncements } from "~/features/announce-read-status";
import {
  decodeVAPIDPublicKey,
  deleteWebPushSubscription,
  fetchWebPushConfig,
  saveWebPushSubscription,
} from "~/features/web-push";
const { isRead, markAsRead, markAsUnread, markAllAsRead } =
  useReadAnnouncements();
useHead({ title: "お知らせ" });
const { data, error, pending, refresh } = await useAsyncData("notices", () =>
  fetchNotices(api),
);
const { data: pushConfig, error: pushConfigError } = await useAsyncData(
  "web-push-config",
  fetchWebPushConfig,
);
const pushSupported = ref(false);
const pushSubscribed = ref(false);
const pushBusy = ref(false);
const pushError = ref<unknown>();

onMounted(async () => {
  pushSupported.value =
    "serviceWorker" in navigator &&
    "PushManager" in window &&
    "Notification" in window;
  if (!pushSupported.value) return;
  const registration = await navigator.serviceWorker.getRegistration("/");
  pushSubscribed.value = Boolean(
    await registration?.pushManager.getSubscription(),
  );
});

async function enablePush() {
  if (!pushConfig.value?.enabled || !pushConfig.value.vapid_public_key) return;
  pushBusy.value = true;
  pushError.value = undefined;
  try {
    const permission = await Notification.requestPermission();
    if (permission !== "granted") {
      throw new Error("ブラウザで通知が許可されていません");
    }
    const registration = await navigator.serviceWorker.register(
      "/web-push-sw.js",
      { scope: "/" },
    );
    const existing = await registration.pushManager.getSubscription();
    const subscription =
      existing ??
      (await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: decodeVAPIDPublicKey(
          pushConfig.value.vapid_public_key,
        ),
      }));
    await saveWebPushSubscription(subscription);
    pushSubscribed.value = true;
  } catch (error) {
    pushError.value = error;
  } finally {
    pushBusy.value = false;
  }
}

async function disablePush() {
  pushBusy.value = true;
  pushError.value = undefined;
  try {
    const registration = await navigator.serviceWorker.getRegistration("/");
    const subscription = await registration?.pushManager.getSubscription();
    if (subscription) {
      await deleteWebPushSubscription(subscription.endpoint);
      await subscription.unsubscribe();
    }
    pushSubscribed.value = false;
  } catch (error) {
    pushError.value = error;
  } finally {
    pushBusy.value = false;
  }
}
</script>
<template>
  <main>
    <header class="linked-page-heading">
      <h1 class="page-title">お知らせ</h1>
    </header>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><section
        class="notification-settings"
        aria-labelledby="notification-settings-title"
      >
        <div>
          <h2 id="notification-settings-title">Web Push通知</h2>
          <p v-if="pushSubscribed">このブラウザで通知を受け取ります。</p>
          <p v-else-if="!pushSupported">
            このブラウザはWeb Pushに対応していません。
          </p>
          <p v-else-if="pushConfig && !pushConfig.enabled">
            サーバー側のWeb Push設定が必要です。
          </p>
          <p v-else>新しいお知らせをブラウザの通知で受け取れます。</p>
        </div>
        <button
          v-if="pushSubscribed"
          class="button-secondary"
          :disabled="pushBusy"
          @click="disablePush"
        >
          通知を停止する
        </button>
        <button
          v-else
          class="button-secondary"
          :disabled="pushBusy || !pushSupported || !pushConfig?.enabled"
          @click="enablePush"
        >
          通知を有効にする
        </button>
        <p v-if="pushError || pushConfigError" class="form-error" role="alert">
          {{ String(pushError ?? pushConfigError) }}
        </p>
      </section>
      <div class="actions">
        <button class="button-secondary" @click="markAllAsRead(data ?? [])">
          すべて既読にする
        </button>
      </div>
      <article v-for="n in data" :key="n.slug" class="notice-entry">
        <time>{{ new Date(n.effectiveFrom).toLocaleString("ja-JP") }}</time>
        <div>
          <h2>
            <NuxtLink :to="`/announces/${n.slug}`">{{ n.title }}</NuxtLink>
          </h2>
          <MarkdownContent :source="n.markdown" /><button
            class="button-secondary"
            @click="isRead(n.slug) ? markAsUnread(n.slug) : markAsRead(n.slug)"
          >
            {{ isRead(n.slug) ? "未読に戻す" : "既読にする" }}
          </button>
        </div>
      </article>
      <p v-if="!data?.length">現在お知らせはありません。</p></RequestState
    >
  </main>
</template>
