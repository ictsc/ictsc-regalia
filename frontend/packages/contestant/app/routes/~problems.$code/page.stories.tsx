import type { Meta, StoryObj } from "@storybook/react";
import { action } from "@storybook/addon-actions";
import {
  Page,
  Content,
  SubmissionListItem,
  SubmissionForm,
  DeploymentListItem,
  ListContainer,
  EmptyListContainer,
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
    redeployable: true,
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
        latestPenalty={0}
      />
    ),
    submissionList: (
      <ListContainer>
        {Array.from({ length: 10 }).map((_, i) => (
          <SubmissionListItem
            key={i}
            id={i + 1}
            submittedAt="2025-02-03T00:00:00Z"
            score={{ maxScore: 100 }}
          />
        ))}
      </ListContainer>
    ),
    deploymentList: (() => {
      return (
        <ListContainer>
          <DeploymentListItem
            key={3}
            event={{
              revision: 3,
              occuredAt: "2021-01-10T12:00:00Z",
              totalPenalty: -50,
              type: "展開中",
              isDeploying: true,
            }}
            maxRedeployment={2}
            deploymentDetail={{
              revision: 3,
              remainingRedeploys: -1,
              exceededRedeployLimit: true,
              totalPenalty: -30,
            }}
            isDeploying={true}
            isLatest={true}
          />
          <DeploymentListItem
            key={2}
            event={{
              revision: 2,
              occuredAt: "2021-01-10T11:00:00Z",
              totalPenalty: 0,
              type: "展開完了",
              isDeploying: true,
            }}
            maxRedeployment={2}
            deploymentDetail={{
              revision: 2,
              remainingRedeploys: 0,
              exceededRedeployLimit: false,
              totalPenalty: 0,
            }}
            isDeploying={false}
            isLatest={false}
          />
          <DeploymentListItem
            key={1}
            event={{
              revision: 2,
              occuredAt: "2021-01-10T10:00:00Z",
              totalPenalty: 0,
              type: "展開完了",
              isDeploying: true,
            }}
            maxRedeployment={2}
            deploymentDetail={{
              revision: 1,
              remainingRedeploys: 1,
              exceededRedeployLimit: false,
              totalPenalty: 0,
            }}
            isDeploying={false}
            isLatest={false}
          />
        </ListContainer>
      );
    })(),
  },
} as Meta<typeof Page>;

type Story = StoryObj<typeof Page>;

const submitAction = action("submit");

export const Default = {} as Story;
export const Empty = {
  args: {
    submissionList: <EmptyListContainer message={""} />,
  },
} as Story;
