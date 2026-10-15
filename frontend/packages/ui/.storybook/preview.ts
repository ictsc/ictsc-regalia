import type { Preview } from "@storybook/vue3-vite";
import { setup } from "@storybook/vue3-vite";
import { createMemoryHistory, createRouter, RouterLink } from "vue-router";
import "../styles.css";
import "../application.css";
import "katex/dist/katex.min.css";

const router = createRouter({
  history: createMemoryHistory(),
  routes: [{ path: "/:pathMatch(.*)*", component: { render: () => null } }],
});

setup((app) => {
  app.use(router);
  // NuxtLink renders RouterLink in the apps, so aria-current styling matches.
  app.component("NuxtLink", RouterLink);
});

const preview: Preview = {
  globalTypes: {
    theme: {
      description: "Color theme",
      toolbar: {
        title: "Theme",
        icon: "mirror",
        items: [
          { value: "light", title: "Light" },
          { value: "dark", title: "Dark" },
        ],
        dynamicTitle: true,
      },
    },
  },
  initialGlobals: { theme: "light" },
  async beforeEach({ parameters }) {
    await router.replace((parameters.route as string | undefined) ?? "/");
  },
  decorators: [
    (_story, { globals, parameters }) => {
      const root = document.documentElement;
      root.lang = "ja";
      root.dataset.theme = globals.theme as string;
      root.dataset.themePreference = globals.theme as string;
      document.body.className =
        (parameters.bodyClass as string | undefined) ?? "";
      return { template: "<story />" };
    },
  ],
};
export default preview;
