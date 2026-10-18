import type { Meta, StoryObj } from "@storybook/vue3-vite";
import { ref } from "vue";
import NoticeListView from "./NoticeListView.vue";
import WebPushSettings from "./WebPushSettings.vue";
import { withContestantFrame } from "~/stories/decorators";
import { notices } from "~/stories/fixtures";

const meta = {
  component: NoticeListView,
  title: "Contestant/Announces/NoticeListView",
  decorators: [
    withContestantFrame(),
    () => ({ template: "<main><story /></main>" }),
  ],
  parameters: { layout: "fullscreen", route: "/announces" },
  args: { notices, readSlugs: ["start"] },
  // Toggle read state locally so the buttons work in the preview.
  render: (args) => ({
    components: { NoticeListView, WebPushSettings },
    setup() {
      const read = ref([...args.readSlugs]);
      const toggleRead = (slug: string) => {
        read.value = read.value.includes(slug)
          ? read.value.filter((s) => s !== slug)
          : [...read.value, slug];
      };
      const markAllRead = () => {
        read.value = (args.notices ?? []).map((n) => n.slug);
      };
      return { args, read, toggleRead, markAllRead };
    },
    template: `
<NoticeListView v-bind="args" :read-slugs="read" @toggle-read="toggleRead" @mark-all-read="markAllRead">
  <template #settings>
    <WebPushSettings :subscribed="false" :supported="true" :server-enabled="true" :busy="false" />
  </template>
</NoticeListView>`,
  }),
} satisfies Meta<typeof NoticeListView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const AllRead: Story = {
  args: { readSlugs: notices.map((n) => n.slug) },
};

export const Empty: Story = {
  args: { notices: [] },
};

export const Loading: Story = {
  args: { notices: undefined, pending: true },
};

export const Error: Story = {
  args: { notices: undefined, error: "お知らせを取得できませんでした" },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
