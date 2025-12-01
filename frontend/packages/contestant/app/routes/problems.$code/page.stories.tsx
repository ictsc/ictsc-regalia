import type { Meta, StoryObj } from "@storybook/react";
import { action } from "storybook/actions";
import {
  Page,
  Content,
  SubmissionList,
  SubmissionListItem,
  EmptySubmissionList,
  SubmissionForm,
} from "./page";

export default {
  title: "pages/problem",
  render: (props) => (
    <div
      style={
        {
          "--header-height": "0",
          "--content-height": "100vh",
          "--content-width": "100%",
        } as React.CSSProperties
      }
    >
      <Page {...props} />
    </div>
  ),
  args: {
    redeployable: false,
    content: (
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
    ),
    submissionForm: (
      <SubmissionForm
        action={() => {
          submitAction("submit");
          return Promise.resolve("success");
        }}
      />
    ),
    submissionList: (
      <SubmissionList>
        {Array.from({ length: 10 }).map((_, i) => (
          <SubmissionListItem
            key={i}
            id={i + 1}
            submittedAt="2025-02-03T00:00:00Z"
            score={{ maxScore: 100 }}
            downloadAnswer={() => {}}
          />
        ))}
      </SubmissionList>
    ),
    deploymentList: null,
  },
} as Meta<typeof Page>;

type Story = StoryObj<typeof Page>;

const submitAction = action("submit");

export const Default = {} as Story;
export const Empty = {
  args: {
    submissionList: <EmptySubmissionList />,
  },
} as Story;
