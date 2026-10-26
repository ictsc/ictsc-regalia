import type { Meta, StoryObj } from "@storybook/vue3-vite";
import SettingsView from "./SettingsView.vue";
import { withAdminShell } from "~/stories/decorators";
import { ruleMarkdown } from "~/stories/fixtures";

const meta = {
  component: SettingsView,
  title: "Admin/Settings/SettingsView",
  decorators: [withAdminShell("/settings")],
  parameters: { layout: "fullscreen" },
  args: {
    initialMarkdown: ruleMarkdown,
    initialFreeze: "2026-09-02T13:13",
    busy: false,
    message: "",
  },
} satisfies Meta<typeof SettingsView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const NoFreeze: Story = {
  args: { initialFreeze: "" },
};

export const EmptyRule: Story = {
  args: { initialMarkdown: "" },
};

/** 二次予選モードは凍結を使わないので、凍結の設定を出さない。 */
export const SecondaryQualifier: Story = {
  args: { secondaryQualifier: true },
};

/** 開発モードかつ二次予選モードでだけ出る競技状態の切り替え。 */
export const DevContestPhase: Story = {
  args: {
    secondaryQualifier: true,
    contestPhase: {
      enabled: true,
      phase: "DURING",
      sections: [
        {
          slug: "day1-am",
          beginning: "2026-09-01T10:00:00+09:00",
          ending: "2026-09-01T12:00:00+09:00",
        },
        {
          slug: "day1-pm",
          beginning: "2026-09-01T13:00:00+09:00",
          ending: "2026-09-01T16:00:00+09:00",
        },
      ],
    },
  },
};

export const Saved: Story = {
  args: { message: "ルール更新しました" },
};
