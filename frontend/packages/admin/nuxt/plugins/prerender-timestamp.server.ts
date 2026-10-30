// Keep prerendered HTML byte-identical when a preview head is rebuilt.
export default defineNuxtPlugin((nuxtApp) => {
  if (!process.env.SOURCE_DATE_EPOCH) return;
  const timestamp = Number(process.env.SOURCE_DATE_EPOCH) * 1000;
  if (!Number.isSafeInteger(timestamp)) throw new Error("invalid SOURCE_DATE_EPOCH");
  nuxtApp.payload.prerenderedAt = timestamp;
});
