import { api, expectData, expectNoContent } from "@ictsc/api";

export async function fetchWebPushConfig() {
  return expectData(await api.GET("/api/v1/contestant/web-push/config"));
}

export async function saveWebPushSubscription(subscription: PushSubscription) {
  const p256dh = subscription.getKey("p256dh");
  const auth = subscription.getKey("auth");
  if (!p256dh || !auth)
    throw new Error("Web Pushの購読鍵を取得できませんでした");
  expectNoContent(
    await api.PUT("/api/v1/contestant/web-push/subscription", {
      body: {
        endpoint: subscription.endpoint,
        keys: {
          p256dh: toBase64URL(p256dh),
          auth: toBase64URL(auth),
        },
      },
    }),
  );
}

export async function deleteWebPushSubscription(endpoint: string) {
  expectNoContent(
    await api.DELETE("/api/v1/contestant/web-push/subscription", {
      body: { endpoint },
    }),
  );
}

export function decodeVAPIDPublicKey(value: string): Uint8Array<ArrayBuffer> {
  const padding = "=".repeat((4 - (value.length % 4)) % 4);
  const binary = atob(
    (value + padding).replaceAll("-", "+").replaceAll("_", "/"),
  );
  return Uint8Array.from(binary, (character) => character.charCodeAt(0));
}

function toBase64URL(value: ArrayBuffer): string {
  const bytes = new Uint8Array(value);
  let binary = "";
  for (const byte of bytes) binary += String.fromCharCode(byte);
  return btoa(binary)
    .replaceAll("+", "-")
    .replaceAll("/", "_")
    .replaceAll("=", "");
}
