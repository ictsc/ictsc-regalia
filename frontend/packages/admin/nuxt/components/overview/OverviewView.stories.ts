import type { Meta, StoryObj } from "@storybook/vue3-vite";
import OverviewView, { type OverviewCard } from "./OverviewView.vue";
import { withAdminShell } from "~/stories/decorators";
import { contentStatus } from "~/stories/fixtures";

const cards: OverviewCard[] = [
  { label: "チーム", detail: "登録済みチームを管理", to: "/teams", count: 48 },
  {
    label: "参加者",
    detail: "参加者と所属を確認",
    to: "/contestants",
    count: 173,
  },
  {
    label: "回答",
    detail: "提出された回答を確認・採点",
    to: "/submissions",
    count: 412,
  },
  {
    label: "進行中の再展開",
    detail: "環境の準備状況を確認",
    to: "/deployments",
    count: 3,
  },
];

const meta = {
  component: OverviewView,
  title: "Admin/Overview/OverviewView",
  decorators: [withAdminShell("/")],
  parameters: { layout: "fullscreen" },
  args: { cards, content: contentStatus },
} satisfies Meta<typeof OverviewView>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Typical: Story = {};

export const PartialFailure: Story = {
  args: {
    cards: cards.map((card) =>
      card.to === "/deployments"
        ? { ...card, count: undefined, error: new Error("failed") }
        : card,
    ),
  },
};

export const LastKnownGood: Story = {
  args: {
    content: {
      ...contentStatus,
      state: "FAILED",
      last_error: "manifest.yaml: problems[3].code is duplicated",
      serving_last_known_good: true,
    },
  },
};

export const NeverDelivered: Story = {
  args: {
    content: {
      ...contentStatus,
      state: "NEVER",
      active_commit: null,
      activated_at: null,
    },
  },
};

export const ContentFailed: Story = {
  args: { content: undefined, contentError: "コンテンツ状態を取得できません" },
};

export const Loading: Story = {
  args: { cards: [], content: undefined, pending: true },
};

export const Mobile: Story = {
  globals: { viewport: { value: "mobile1", isRotated: false } },
};
