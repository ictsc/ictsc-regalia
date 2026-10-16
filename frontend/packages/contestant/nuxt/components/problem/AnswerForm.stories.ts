import type { Meta, StoryObj } from "@storybook/vue3-vite";
import { ref } from "vue";
import AnswerForm from "./AnswerForm.vue";
import { withContestantFrame, withDetailContent } from "~/stories/decorators";
import { answerBody, fixedNow } from "~/stories/fixtures";

const meta = {
  component: AnswerForm,
  title: "Contestant/ProblemDetail/AnswerForm",
  decorators: [withContestantFrame(), withDetailContent],
  parameters: { layout: "fullscreen", bodyClass: "detail-page" },
  args: {
    body: "",
    open: true,
    sending: false,
    submitted: false,
    cooldownSeconds: 0,
    savedAt: "",
    saveError: "",
  },
  // Keep the textarea editable in the preview.
  render: (args) => ({
    components: { AnswerForm },
    setup: () => ({ args, body: ref(args.body) }),
    template: `<AnswerForm v-bind="args" v-model:body="body" />`,
  }),
} satisfies Meta<typeof AnswerForm>;
export default meta;
type Story = StoryObj<typeof meta>;

const savedAt = new Date(fixedNow - 90_000).toISOString();

export const Empty: Story = {};

export const DraftSaved: Story = {
  args: { body: answerBody, savedAt },
};

export const Sending: Story = {
  args: { body: answerBody, savedAt, sending: true },
};

export const Submitted: Story = {
  args: { body: answerBody, savedAt, submitted: true, cooldownSeconds: 1199 },
};

export const Cooldown: Story = {
  args: { body: answerBody, savedAt, cooldownSeconds: 754 },
};

export const Closed: Story = {
  args: { open: false },
};

export const SubmitError: Story = {
  args: {
    body: answerBody,
    savedAt,
    actionError: new Error("回答を提出できませんでした。"),
  },
};

export const DraftSaveError: Story = {
  args: {
    body: answerBody,
    saveError: "下書きを保存できません。入力内容をコピーして保管してください。",
  },
};
