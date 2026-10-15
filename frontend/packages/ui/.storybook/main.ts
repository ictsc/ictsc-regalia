import type { StorybookConfig } from "@storybook/vue3-vite";
import vue from "@vitejs/plugin-vue";
import { japaneseFontStylesheet, latinFontStylesheet } from "../fonts";

const config: StorybookConfig = {
  stories: ["../components/**/*.stories.ts"],
  framework: "@storybook/vue3-vite",
  core: { disableTelemetry: true },
  // Load the same web fonts as the Nuxt apps' app.head links.
  previewHead: (head) => `${head}
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link rel="stylesheet" href="${latinFontStylesheet}">
<link rel="stylesheet" href="${japaneseFontStylesheet}">`,
  viteFinal: async (config) => ({
    ...config,
    plugins: [...(config.plugins ?? []), vue()],
    resolve: { ...config.resolve, dedupe: ["vue", "vue-router"] },
  }),
};
export default config;
