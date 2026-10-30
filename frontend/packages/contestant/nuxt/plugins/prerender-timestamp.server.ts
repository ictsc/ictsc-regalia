// Nuxt otherwise embeds Date.now() in prerendered HTML, making identical PR
// heads produce different OCI manifests. SOURCE_DATE_EPOCH is set only by the
// trusted preview build workflow from the exact checked-out commit.
export default defineNuxtPlugin((nuxtApp) => {
  if (!import.meta.prerender || !process.env.SOURCE_DATE_EPOCH) return;
  const timestamp = Number(process.env.SOURCE_DATE_EPOCH) * 1000;
  if (!Number.isSafeInteger(timestamp)) throw new Error("invalid SOURCE_DATE_EPOCH");
  nuxtApp.hook("app:rendered", ({ ssrContext }) => {
    if (ssrContext) ssrContext.payload.prerenderedAt = timestamp;
  });
});
