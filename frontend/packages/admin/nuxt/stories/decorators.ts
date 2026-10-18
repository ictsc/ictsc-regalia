import type { Decorator } from "@storybook/vue3-vite";
import { useRouter } from "vue-router";
import AdminHeader from "~/components/layout/AdminHeader.vue";
import AdminNav from "~/components/layout/AdminNav.vue";
import { adminName } from "./fixtures";

// Mirrors layouts/default.vue so admin.css selectors such as
// `.admin-shell > .workspace` apply to the story.
export function withAdminShell(
  currentPath: string,
  mainClass = "workspace",
): Decorator {
  return () => ({
    components: { AdminHeader, AdminNav },
    setup() {
      // Keep RouterLink's own active state in sync with the highlighted section.
      void useRouter().replace(currentPath);
      return { currentPath, mainClass, adminName };
    },
    template: `<div class="admin-app"><AdminHeader signed-in :display-name="adminName" /><div class="admin-shell"><AdminNav :current-path="currentPath" /><main :class="mainClass"><story /></main></div></div>`,
  });
}
