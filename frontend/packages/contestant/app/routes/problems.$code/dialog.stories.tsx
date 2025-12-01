import type { Meta, StoryObj } from "@storybook/react";
import { ConfirmModal } from "./confirmModal";
import { Markdown, Typography } from "@app/components/markdown";

export default {
  title: "pages/problem/confirmModal",
} satisfies Meta;

type Story = StoryObj;

export const SubmitDeployment: Story = {
  render: () => (
    <ConfirmModal
      isOpen={true}
      onConfirm={() => {}}
      onCancel={() => {}}
      title="再展開の確認"
      confirmText="再展開する"
      cancelText="キャンセル"
      dialogClassName="w-full max-w-md transform rounded-8 bg-surface-0 p-16 text-left align-middle shadow-xl transition-all"
    >
      <span>ここに本文を書く</span>
    </ConfirmModal>
  ),
};

export const SubmitAnswer: Story = {
  render: () => (
    <ConfirmModal
      isOpen={true}
      onConfirm={() => {}}
      onCancel={() => {}}
      title="解答の確認"
      confirmText="送信する"
      cancelText="キャンセル"
      dialogClassName="w-full max-w-[1024px] transform rounded-8 bg-surface-0 p-16 text-left align-middle shadow-xl transition-all"
    >
      <div style={{ padding: "16px" }}>
        <Typography>
          <Markdown>{markdownContent}</Markdown>
        </Typography>
      </div>
    </ConfirmModal>
  ),
};

const markdownContent = `
# 見出し1
## 見出し2
### 見出し3
#### 見出し4
##### 見出し5
###### 見出し6

\`\`\`shell
echo hoge
pwd
ls
\`\`\`

| 列1 | 列2 | 列3 |
|-----|-----|-----|
| A1  | B1  | C1  |
| A2  | B2  | C2  |
| A3  | B3  | C3  |
`;
