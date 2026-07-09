export default defineNuxtConfig({
  compatibilityDate: "2026-07-09",
  devtools: { enabled: true },
  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL ?? "/api",
    },
  },
  typescript: {
    strict: true,
    typeCheck: true,
  },
});
