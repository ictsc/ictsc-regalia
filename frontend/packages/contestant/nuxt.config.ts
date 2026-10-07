import { fileURLToPath } from "node:url";
import { themeInitScript } from "../ui/theme-init";
import { japaneseFontStylesheet, latinFontStylesheet } from "../ui/fonts";
export default defineNuxtConfig({
  compatibilityDate: "2026-09-05",
  srcDir: "nuxt/",
  ssr: false,
  devtools: { enabled: false },
  devServer: { port: 3000 },
  runtimeConfig: {
    public: {
      demoMode: process.env.ICTSC_DEMO_MODE === "true",
    },
  },
  app: {
    baseURL: "/",
    head: {
      titleTemplate: "%s / ICTSC REGALIA",
      htmlAttrs: { lang: "ja" },
      link: [
        { rel: "preconnect", href: "https://fonts.googleapis.com" },
        {
          rel: "preconnect",
          href: "https://fonts.gstatic.com",
          crossorigin: "anonymous",
        },
        {
          rel: "stylesheet",
          href: latinFontStylesheet,
        },
        { rel: "stylesheet", href: japaneseFontStylesheet },
      ],
      script: [
        { innerHTML: themeInitScript, tagPosition: "head" },
        { src: "/runtime-config.js" },
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
        "/api": {
          target: process.env.ICTSC_API_PROXY_TARGET || "http://localhost:8080",
          changeOrigin: true,
        },
      },
    },
  },
  nitro: { prerender: { crawlLinks: false } },
});
