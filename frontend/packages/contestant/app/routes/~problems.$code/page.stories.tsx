import type { Meta, StoryObj } from "@storybook/react";
import { Layout, Content, Sidebar } from "./page";

export default {
  title: "pages/problem",
} as Meta;

type Story = StoryObj<typeof Layout>;

export const Default = {
  render: () => (
    <div
      style={
        {
          "--header-height": "0",
          "--content-height": "100vh",
        } as React.CSSProperties
      }
    >
      <Layout
        content={
          <Content
            code="AAA"
            title="Title"
            body={`
## 概要

あれこれ

## 前提条件

- いろいろ
- あるよね

## 初期状態

- こうなってる

## 終了状態

- こうなる

## 接続情報

| ホスト名 | IPアドレス | ユーザ名 | パスワード |
| --- | --- | --- | --- |
| Web | 192.168.0.1 | user | password |
`}
          />
        }
        sidebar={<Sidebar />}
      />
    </div>
  ),
} as Story;
