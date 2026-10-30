// Keep prerendered HTML byte-identical when a preview head is rebuilt.
export default defineNuxtPlugin((nuxtApp) => {
  if (!import.meta.prerender || !process.env.SOURCE_DATE_EPOCH) return;
  const timestamp = Number(process.env.SOURCE_DATE_EPOCH) * 1000;
  if (!Number.isSafeInteger(timestamp)) throw new Error("invalid SOURCE_DATE_EPOCH");
  nuxtApp.hook("app:rendered", ({ ssrContext }) => {
    if (ssrContext) ssrContext.payload.prerenderedAt = timestamp;
  });
});
