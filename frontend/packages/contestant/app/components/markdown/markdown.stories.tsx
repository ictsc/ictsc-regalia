import type { Meta, StoryObj } from "@storybook/react";
import { Typography, Markdown } from "./markdown";

export default {
  title: "components/Markdown",
  component: Markdown,
} satisfies Meta<typeof Markdown>;

type Story = StoryObj<typeof Markdown>;

export const HeadingsAndParagraphs: Story = {
  render: () => (
    <Typography>
      <Markdown>
        {`
# 見出し1
## 見出し2
### 見出し3
#### 見出し4
##### 見出し5
###### 見出し6

通常の段落テキストです。
複数行の段落も
このように書けます。

段落と段落の間は

このように空行を入れます。
`}
      </Markdown>
    </Typography>
  ),
};

export const ListsAndIndentation: Story = {
  render: () => (
    <Typography>
      <Markdown>
        {`
- 箇条書き1レベル目
  - 2レベル目
    - 3レベル目
      - 4レベル目
- 別の項目

1. 番号付きリスト
2. 2番目の項目
   1. ネストした番号付き
   2. 2番目のネスト
3. 3番目の項目

- [ ] チェックボックス（未チェック）
- [x] チェックボックス（チェック済み）
`}
      </Markdown>
    </Typography>
  ),
};

export const TextFormatting: Story = {
  render: () => (
    <Typography>
      <Markdown>
        {`
*イタリック* _イタリック_
**太字** __太字__
***太字イタリック*** ___太字イタリック___
~打ち消し線~
\`inline code\`

---

> 引用文
> 複数行の
> 引用文

> ネストした
>> 引用文
`}
      </Markdown>
    </Typography>
  ),
};

export const CodeAndTables: Story = {
  render: () => (
    <Typography>
      <Markdown>
        {`
問題文がここに入る  
hogehoge  
## こんな感じ
diff  
\`\`\`diff
- echo "Hello, World!"
+ echo "Hello, New World!"
\`\`\`
shellscript  
\`\`\`shellscript
echo "Hello, World!"
mkdir -p /tmp/hello
cd /tmp/hello
\`\`\`
shellsession  
\`\`\`shellsession
user@problem-nle-vm:~$ ls -la
total 32
drwxr-x--- 4 user user 4096 Mar 15 16:10 .
drwxr-xr-x 5 root root 4096 Feb 22 04:49 ..
user@problem-nle-vm:~$ sudo chmod 644 memo.txt
\`\`\`
hcl
\`\`\`hcl
io_mode = "async"

service "http" "web_proxy" {
  listen_addr = "127.0.0.1:8080"
  
  process "main" {
    command = ["/usr/local/bin/awesome-app", "server"]
  }
}

\`\`\`

| 列1 | 列2 | 列3 |
|-----|-----|-----|
| A1  | B1  | C1  |
| A2  | B2  | C2  |
| A3  | B3  | C3  |
`}
      </Markdown>
    </Typography>
  ),
};

export const LinksAndImages: Story = {
  render: () => (
    <Typography>
      <Markdown>
        {`
[リンクテキスト](https://example.com)
![画像の代替テキスト](https://example.com/image.jpg)

自動リンク: https://example.com
`}
      </Markdown>
    </Typography>
  ),
};
