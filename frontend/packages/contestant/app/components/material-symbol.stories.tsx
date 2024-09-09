import type { Meta, StoryObj } from "@storybook/react";
import { MaterialSymbol } from "./material-symbol";

export default {
  title: "components/MaterialSymbol",
  component: MaterialSymbol,
} satisfies Meta<typeof MaterialSymbol>;

type Story = StoryObj<typeof MaterialSymbol>;

export const Default: Story = {
  args: {
    icon: "schedule",
    fill: false,
    size: 24,
  },
};
