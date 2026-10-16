import type { Meta, StoryObj } from "@storybook/vue3-vite";
import AnswerHistory from "./AnswerHistory.vue";
import { withContestantFrame, withDetailContent } from "~/stories/decorators";
import { answerBody, answers } from "~/stories/fixtures";

const meta = {
  component: AnswerHistory,
  title: "Contestant/ProblemDetail/AnswerHistory",
  decorators: [withContestantFrame(), withDetailContent],
  parameters: { layout: "fullscreen", bodyClass: "detail-page" },
  args: { problemCode: "B02", answers, selectedAnswerNumber: null },
} satisfies Meta<typeof AnswerHistory>;
export default meta;
type Story = StoryObj<typeof meta>;

export const None: Story = {
  args: { answers: [] },
};

export const Several: Story = {};

export const AnswerSelected: Story = {
  args: { selectedAnswerNumber: 2, selectedAnswerBody: answerBody },
  parameters: { route: "/problems/B02?answer=2" },
};
