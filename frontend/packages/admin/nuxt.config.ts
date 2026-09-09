import { fileURLToPath } from "node:url";
export default defineNuxtConfig({
  compatibilityDate: "2026-09-05",
  srcDir: "nuxt/",
  ssr: false,
  devtools: { enabled: false },
  devServer: { port: 3001 },
  app: {
    baseURL: "/admin/",
    head: {
      titleTemplate: "%s / ICTSC REGALIA",
      htmlAttrs: { lang: "ja" },
      link: [
        {
          rel: "stylesheet",
          href: "https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500&family=Zen+Kaku+Gothic+New:wght@400;500;700&display=swap",
        },
      ],
    },
  },
  css: [
    "@ictsc/ui/styles.css",
    "@ictsc/ui/application.css",
    "katex/dist/katex.min.css",
  ],
  components: [
    { path: "~/components" },
    {
      path: fileURLToPath(new URL("../ui/components", import.meta.url)),
      pathPrefix: false,
      extensions: ["vue"],
    },
  ],
  modules: ["@nuxt/eslint"],
  eslint: { config: { stylistic: false } },
  vite: {
    server: {
      proxy: {
        "/api": { target: "http://localhost:8080", changeOrigin: true },
      },
    },
  },
  nitro: { prerender: { crawlLinks: false } },
});
