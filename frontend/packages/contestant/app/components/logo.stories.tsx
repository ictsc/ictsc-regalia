import type { Meta, StoryObj } from "@storybook/react";
import { Logo } from "./logo";

export default {
  title: "components/Logo",
  component: Logo,
} satisfies Meta<typeof Logo>;

type Story = StoryObj<typeof Logo>;

export const Default: Story = {
  render: () => <Logo height={100} />,
}

