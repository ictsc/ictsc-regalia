import type { Meta, StoryObj } from "@storybook/vue3-vite";
import AdminNav from "./AdminNav.vue";

const meta = {
  component: AdminNav,
  title: "Admin/Layout/AdminNav",
  decorators: [
    () => ({
      template: `<div class="admin-app"><div class="admin-shell"><story /><main class="workspace" /></div></div>`,
    }),
  ],
  parameters: { layout: "fullscreen" },
  args: { currentPath: "/" },
} satisfies Meta<typeof AdminNav>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Overview: Story = {};

export const NestedPage: Story = {
  args: { currentPath: "/teams/12" },
};

export const Grading: Story = {
  args: { currentPath: "/submissions" },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
