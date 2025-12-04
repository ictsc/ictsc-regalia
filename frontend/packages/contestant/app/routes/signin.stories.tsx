import type { Meta, StoryObj } from "@storybook/react";
import { SignInPage } from "./signin.page";

export default {
  title: "pages/signin",
  component: SignInPage,
} satisfies Meta<typeof SignInPage>;

type Story = StoryObj<typeof SignInPage>;

export const Default: Story = {
  name: "デフォルト",
  args: {
    signInURL: "/",
  },
};
