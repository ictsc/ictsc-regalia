import type { Meta, StoryObj } from "@storybook/vue3-vite";
import TeamDirectoryView from "./TeamDirectoryView.vue";
import { withContestantFrame } from "~/stories/decorators";
import { teams } from "~/stories/fixtures";

const meta = {
  component: TeamDirectoryView,
  title: "Contestant/Teams/TeamDirectoryView",
  decorators: [
    withContestantFrame(),
    () => ({ template: `<main class="workspace"><story /></main>` }),
  ],
  parameters: { layout: "fullscreen", route: "/teams" },
  args: { teams },
} satisfies Meta<typeof TeamDirectoryView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const LongIntroduction: Story = {
  args: {
    teams: [
      {
        ...teams[0]!,
        members: [
          {
            name: "alice",
            displayName: "Alice",
            selfIntroduction: Array.from(
              { length: 4 },
              () =>
                "学部ではネットワークの研究をしていて、普段はルーターの設定やパケットキャプチャをしています。",
            ).join("\n\n"),
          },
        ],
      },
    ],
  },
};

export const Loading: Story = {
  args: { teams: undefined, pending: true },
};
