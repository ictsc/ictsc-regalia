import type { Meta, StoryObj } from "@storybook/vue3-vite";
import AsciiLogo from "./AsciiLogo.vue";

const meta = {
  component: AsciiLogo,
  title: "UI/Brand/AsciiLogo",
} satisfies Meta<typeof AsciiLogo>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const AsHomeLink: Story = {
  render: () => ({
    components: { AsciiLogo },
    template: `<a class="brand" href="#" aria-label="ICTSC 2026 REGALIA ホーム"><AsciiLogo /></a>`,
  }),
};
