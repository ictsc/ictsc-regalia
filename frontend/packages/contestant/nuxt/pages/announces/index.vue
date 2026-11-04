<script setup lang="ts">
import NoticeListView from "~/components/announce/NoticeListView.vue";
import WebPushSettings from "~/components/announce/WebPushSettings.vue";
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
const readSlugs = computed(() =>
  (data.value ?? []).filter((n) => isRead(n.slug)).map((n) => n.slug),
);
function toggleRead(slug: string) {
  if (isRead(slug)) markAsUnread(slug);
  else markAsRead(slug);
}
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
    await navigator.serviceWorker.register("/web-push-sw.js", { scope: "/" });
    // subscribe() は有効化済みの Service Worker が必要。登録直後は installing のため待つ。
    const registration = await navigator.serviceWorker.ready;
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
    <NoticeListView
      :notices="data ?? undefined"
      :read-slugs="readSlugs"
      :pending="pending"
      :error="error"
      @toggle-read="toggleRead"
      @mark-all-read="markAllAsRead(data ?? [])"
      @retry="refresh"
    >
      <template #settings>
        <WebPushSettings
          :subscribed="pushSubscribed"
          :supported="pushSupported"
          :server-enabled="pushConfig?.enabled"
          :busy="pushBusy"
          :error="pushError ?? pushConfigError ?? undefined"
          @enable="enablePush"
          @disable="disablePush"
        />
      </template>
    </NoticeListView>
  </main>
</template>
