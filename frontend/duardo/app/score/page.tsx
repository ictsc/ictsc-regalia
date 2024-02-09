"use client";

import React from "react";

import clsx from "clsx";

import { Answer } from "@/proto/admin/v1/answer_pb";
import {
  MultipleChoiceProblem,
  Problem,
  ProblemType,
} from "@/proto/admin/v1/problem_pb";

function ProblemHeader({
  className,
  title,
  code,
  point,
}: {
  className?: string;
  title: string;
  code: string;
  point: number;
}) {
  return (
    <div className={className}>
      <h1 className="text-2xl md:text-xl font-bold pb-6">{title}</h1>
      <div className="flex flex-row justify-items-center pb-16">
        <p className="pr-6">
          <span className="text-gray-500 text-xs">問題コード</span>
          <span className="text-xl font-bold">{code}</span>
        </p>
        <p className="pr-6 inline-block">
          <span className="text-xl font-bold">{point}</span>
          <span className="text-gray-500 text-xs">ポイント</span>
        </p>
      </div>
    </div>
  );
}

ProblemHeader.defaultProps = {
  className: "",
};

function ProblemDescriptiveBody({ body }: { body: string }) {
  return (
    <div className="markdown-body">
      <p>{body}</p>
    </div>
  );
}

function ProblemMultipleChoiceBody({ body }: { body: MultipleChoiceProblem }) {
  return <div>Unimplemented</div>;
}

function Card({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <div
      // className={`border md:border-t px-8 pt-12 pb-8 rounded-md shadow-sm ${className}`}
      className={clsx(
        "border md:border-t px-8 pt-12 pb-8 rounded-md shadow-sm",
        className,
      )}
    >
      {children}
    </div>
  );
}

Card.defaultProps = {
  className: "",
};

function AnswerListForm() {
  const answer = new Answer({
    id: "1",
    problemId: "1",
    teamId: "1",
    problemType: ProblemType.DESCRIPTIVE,
    body: {
      case: "descriptive",
      value: {
        body: "aaaa",
      },
    },
  });

  const answers = [answer];

  return (
    <Card>
      <div className="flex flex-row justify-between pb-8">
        <div>チームICTSC（ICTSC大学）</div>
        <time>2023/12/16 18:00:05</time>
      </div>
      回答文
      <div className="divider" />
      <input type="number" className="input input-bordered input-sm" />
      <input
        type="submit"
        className="btn btn-primary btn-sm ml-2"
        value="採点"
      />
    </Card>
  );
}

function Index() {
  // TODO: バックエンド実装後コメントアウトを外す
  // const { data, error, isLoading } = useQuery(getProblem, {});
  //
  // if (isLoading) {
  //   return (
  //     <div className="flex justify-center items-center h-[100vh] w-full">
  //       <span className="loading loading-spinner loading-md" />
  //     </div>
  //   );
  // }

  // TODO: バックエンド実装後コメントアウトを外す
  const problem = new Problem({
    id: "1",
    code: "ABC",
    title: "なんか通信できない(´・ω・｀)",
    point: 100,
    body: {
      case: "descriptive",
      value: {
        body: "# 問題タイトル\n問題文",
        answer: "aaaa",
      },
    },
  });

  return (
    <main>
      <div className="flex flex-row">
        <div className="container mx-auto px-2 flex-grow">
          <ProblemHeader
            className="pt-16"
            title={problem.title}
            code={problem.code}
            point={problem.point}
          />
          <Card className="znc">
            {problem.body?.case === "descriptive" && (
              <ProblemDescriptiveBody body={problem.body.value.body} />
            )}
            {problem.body?.case === "multipleChoice" && (
              <ProblemMultipleChoiceBody body={problem.body.value} />
            )}
          </Card>
          <div className="divider py-12" />
          <AnswerListForm />
        </div>
      </div>
    </main>
  );
}

export default Index;
