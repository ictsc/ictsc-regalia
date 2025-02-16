import type { Meta, StoryObj } from "@storybook/react";
import { action } from "@storybook/addon-actions";
import { SignUpPage } from "./page";
import { startTransition } from "react";

const submitAction = action("submit");

export default {
  title: "pages/signup",
  component: SignUpPage,
  args: {
    submit: (data) => {
      startTransition(async () => {
        submitAction(data);
        await new Promise((resolve) => setTimeout(resolve, 2000));
      });
    },
  },
} satisfies Meta<typeof SignUpPage>;

type Story = StoryObj<typeof SignUpPage>;

export const Default: Story = {
  args: {},
};

export const Error: Story = {
  args: {
    error: "invalid",
    invitationCodeError: "required",
    nameError: "duplicate",
    displayNameError: "invalid",
  },
};
