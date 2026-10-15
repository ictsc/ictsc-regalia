import { fileURLToPath } from "node:url";
import { defineConfig } from "vitest/config";
import vue from "@vitejs/plugin-vue";
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { "~": fileURLToPath(new URL("./nuxt", import.meta.url)) },
  },
  test: { environment: "jsdom", include: ["nuxt/**/*.test.ts"] },
});
