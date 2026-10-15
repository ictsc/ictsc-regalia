import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ThemeToggle from "./ThemeToggle.vue";

const meta = {
  component: ThemeToggle,
  title: "UI/Controls/ThemeToggle",
  parameters: {
    docs: {
      description: {
        component:
          "クリックで自動・明・暗を切り替えます。各ストーリーのテーマはツールバーの Theme で切り替えてください。",
      },
    },
  },
} satisfies Meta<typeof ThemeToggle>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
