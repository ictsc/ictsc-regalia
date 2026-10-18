import path from "node:path";
import { fileURLToPath } from "node:url";
import type { StorybookConfig } from "@storybook/vue3-vite";
import vue from "@vitejs/plugin-vue";
import type { Plugin } from "vite";
import { japaneseFontStylesheet, latinFontStylesheet } from "../fonts";

const packagesDir = fileURLToPath(new URL("../..", import.meta.url));

// Resolve the Nuxt `~/` and `@/` aliases to the importing app's srcDir so
// extracted app components can keep their `~/features/...` imports.
function nuxtSrcAlias(): Plugin {
  return {
    name: "ictsc:nuxt-src-alias",
    enforce: "pre",
    resolveId(source, importer) {
      if (!importer || !/^[~@]\//.test(source)) return null;
      const app = /[\\/]packages[\\/](contestant|admin)[\\/]nuxt[\\/]/.exec(
        importer,
      )?.[1];
      if (!app) return null;
      return this.resolve(
        path.join(packagesDir, app, "nuxt", source.slice(2)),
        importer,
        { skipSelf: true },
      );
    },
  };
}

const config: StorybookConfig = {
  stories: [
    "../components/**/*.stories.ts",
    "../../contestant/nuxt/**/*.stories.ts",
    "../../admin/nuxt/**/*.stories.ts",
  ],
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
    plugins: [...(config.plugins ?? []), vue(), nuxtSrcAlias()],
    resolve: { ...config.resolve, dedupe: ["vue", "vue-router"] },
  }),
};
export default config;
