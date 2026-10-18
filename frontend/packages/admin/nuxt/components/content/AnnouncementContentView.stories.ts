import type { Meta, StoryObj } from "@storybook/vue3-vite";
import AnnouncementContentView from "./AnnouncementContentView.vue";
import { withAdminShell } from "~/stories/decorators";
import { announcements } from "~/stories/fixtures";

const meta = {
  component: AnnouncementContentView,
  title: "Admin/Content/AnnouncementContentView",
  decorators: [withAdminShell("/content/announcements/day2-start")],
  parameters: { layout: "fullscreen" },
  args: { announcement: announcements[0] },
} satisfies Meta<typeof AnnouncementContentView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const Loading: Story = {
  args: { announcement: undefined, pending: true },
};
