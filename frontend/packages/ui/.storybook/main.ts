import type { StorybookConfig } from "@storybook/vue3-vite";
import vue from "@vitejs/plugin-vue";
const config: StorybookConfig = {
  stories: ["../**/*.stories.ts"],
  framework: "@storybook/vue3-vite",
  viteFinal: async (config) => ({
    ...config,
    plugins: [...(config.plugins ?? []), vue()],
  }),
};
export default config;
