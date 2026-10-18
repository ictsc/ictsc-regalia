import type { Meta, StoryObj } from "@storybook/vue3-vite";
import GradingBoard from "./GradingBoard.vue";
import { answerKey } from "~/features/grading";
import { withAdminShell } from "~/stories/decorators";
import { answers, fixedNow, myDiscordId } from "~/stories/fixtures";

const meta = {
  component: GradingBoard,
  title: "Admin/Grading/GradingBoard",
  decorators: [withAdminShell("/submissions", "workspace grading-workspace")],
  parameters: { layout: "fullscreen" },
  args: { answers, myDiscordId, busyKey: "", now: fixedNow },
} satisfies Meta<typeof GradingBoard>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const OtherAdmin: Story = {
  args: { myDiscordId: "110000000000000001" },
};

export const Updating: Story = {
  args: { busyKey: answerKey(answers[0]!) },
};

export const EmptyColumns: Story = {
  args: { answers: answers.filter((a) => a.workflow.status === "WAITING") },
};

export const LongTitles: Story = {
  args: {
    answers: answers.map((a, index) =>
      index === 0
        ? {
            ...a,
            problem: {
              ...a.problem,
              title:
                "複数拠点を結ぶIPsec VPNのうち一部の拠点間でのみ通信が断続的に途切れる",
            },
          }
        : a,
    ),
  },
};
