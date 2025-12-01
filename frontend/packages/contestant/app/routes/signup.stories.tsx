import type { Meta, StoryObj } from "@storybook/react";
import { action } from "storybook/actions";
import { SignUpPage } from "./signup.page";
import { startTransition } from "react";

const submitAction = action("submit");

export default {
  title: "pages/signup",
  component: SignUpPage,
  args: {
    submit: (data: unknown) => {
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
