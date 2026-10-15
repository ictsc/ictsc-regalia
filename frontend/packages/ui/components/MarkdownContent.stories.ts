import type { Meta, StoryObj } from "@storybook/vue3-vite";
import MarkdownContent from "./MarkdownContent.vue";

const meta = {
  component: MarkdownContent,
  title: "UI/Content/MarkdownContent",
  render: (args) => ({
    components: { MarkdownContent },
    setup: () => ({ args }),
    template: `<div class="reading-copy"><MarkdownContent v-bind="args" /></div>`,
  }),
} satisfies Meta<typeof MarkdownContent>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Basic: Story = {
  args: {
    source: `## 問題文

ルーター \`R1\` から \`R2\` へ疎通できない状態です。原因を調査し、復旧してください。

- 構成図は添付の通りです
- **設定の初期化は禁止**です
- 詳しくは [ルール](https://example.com/rule) を参照してください

> 回答には、原因と実施した作業を記載してください。`,
  },
};

export const Table: Story = {
  args: {
    source: `| ホスト | アドレス | 役割 |
| --- | --- | --- |
| R1 | 192.0.2.1/24 | ゲートウェイ |
| R2 | 198.51.100.1/24 | 対向ルーター |
| srv01 | 203.0.113.10/24 | Webサーバー |`,
  },
};

export const Code: Story = {
  args: {
    source: `実行結果は次の通りです。

\`\`\`bash
$ ip route show
default via 192.0.2.1 dev eth0
192.0.2.0/24 dev eth0 proto kernel scope link src 192.0.2.10
\`\`\`

\`\`\`yaml
interfaces:
  eth0:
    address: 192.0.2.10/24
\`\`\``,
  },
};

export const Math: Story = {
  args: {
    source: `得点は $s = \\lfloor 100 \\times r \\rfloor$ で計算します。

$$
P = \\sum_{i=1}^{n} p_i
$$`,
  },
};

export const LongText: Story = {
  args: {
    source: Array.from(
      { length: 8 },
      (_, index) =>
        `### 手順 ${index + 1}\n\n` +
        "ネットワーク機器の設定を確認し、必要に応じて修正してください。".repeat(
          4,
        ),
    ).join("\n\n"),
  },
};
