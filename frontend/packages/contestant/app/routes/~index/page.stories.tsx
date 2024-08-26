import type { Meta, StoryObj } from "@storybook/react";
import { IndexPage } from "./page";

export default {
  title: "pages/index",
  component: IndexPage,
} satisfies Meta<typeof IndexPage>;

type Story = StoryObj<typeof IndexPage>;

export const Default: Story = {
  args: {},
};
