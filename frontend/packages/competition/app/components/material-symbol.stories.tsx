import type { Meta, StoryObj } from "@storybook/react";
import { MaterialSymbol } from "./material-symbol";

export default {
  title: "MaterialSymbol",
  component: MaterialSymbol,
} satisfies Meta<typeof MaterialSymbol>;

type Story = StoryObj<typeof MaterialSymbol>;

export const Default : Story = {
  render: () => <MaterialSymbol icon="schedule" size={24} />,
}

