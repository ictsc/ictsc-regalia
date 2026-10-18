import type { Meta, StoryObj } from "@storybook/vue3-vite";
import ContentView from "./ContentView.vue";
import { withAdminShell } from "~/stories/decorators";
import {
  announcements,
  catalog,
  contentStatus,
  sections,
} from "~/stories/fixtures";

const meta = {
  component: ContentView,
  title: "Admin/Content/ContentView",
  decorators: [withAdminShell("/content")],
  parameters: { layout: "fullscreen" },
  args: {
    status: contentStatus,
    sections: { data: { sections } },
    problems: { data: { problems: catalog } },
    announcements: { data: { announcements } },
    busy: false,
    message: "",
  },
} satisfies Meta<typeof ContentView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const RefreshFailed: Story = {
  args: {
    status: {
      ...contentStatus,
      state: "FAILED",
      last_error: "manifest.yaml: problems[3].code is duplicated",
      serving_last_known_good: true,
    },
  },
};

export const PartialError: Story = {
  args: { problems: { error: "問題一覧を取得できませんでした" } },
};

export const Empty: Story = {
  args: {
    sections: { data: { sections: [] } },
    problems: { data: { problems: [] } },
    announcements: { data: { announcements: [] } },
  },
};
