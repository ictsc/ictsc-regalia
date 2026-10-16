import type { Decorator } from "@storybook/vue3-vite";
import { teamStyle } from "@ictsc/ui/colors";
import { teamColor } from "./fixtures";

// Mirrors layouts/default.vue: team color variables and the site footer.
export function withContestantFrame(color: string = teamColor): Decorator {
  return () => ({
    setup: () => ({ style: teamStyle(color) }),
    template: `<div :style="style"><story /><footer class="site-footer"><span>ICTSC 2026 / REGALIA</span><span>競技システム</span></footer></div>`,
  });
}

// Mirrors the problem detail page: sections live inside the content grid.
export const withDetailContent: Decorator = () => ({
  template: `<div class="detail-layout is-rail-hidden"><article class="problem-content"><div class="content-grid"><story /></div></article></div>`,
});
