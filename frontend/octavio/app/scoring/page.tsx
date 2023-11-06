"use client";

import React, { useState } from "react";

import Image from "next/image";
import Link from "next/link";
import { notFound } from "next/navigation";

import ICTSCCard from "@/components/card";
import ConditionalText from "@/components/conditional-text";
import LoadingPage from "@/components/loading-page";
import MarkdownPreview from "@/components/markdown-preview";
import useAuth from "@/hooks/auth";
import useProblem from "@/hooks/problem";
import useProblems from "@/hooks/problems";

function Index() {
  const { user } = useAuth();
  const { problems, isLoading } = useProblems();
  const [selectedProblemId, setSelectedProblemId] = useState<string | null>(
    null,
  );

  const { problem } = useProblem(selectedProblemId);

  const isFullAccess = user?.user_group.is_full_access ?? false;
  const isReadOnly = user?.is_read_only ?? false;

  if (isLoading) {
    return <LoadingPage />;
  }

  if (!isFullAccess || isReadOnly || problems.length === 0) {
    return notFound();
  }

  return (
    <>
      <div className="overflow-x-auto">
        <table className="table table-compact w-full">
          <thead>
            <tr>
              <th>採点</th>
              <th>未採点 ~15分/15~19分/20分~</th>
              <th>ID</th>
              <th>問題コード</th>
              <th>タイトル</th>
              <th>コンテンツ</th>
              <th>ポイント</th>
              <th>採点基準ポイント</th>
              <th>前提問題</th>
              <th>著者</th>
            </tr>
          </thead>
          <tbody className="cursor-pointer">
            {problems.map((prob) => (
              <tr key={prob.id} onClick={() => setSelectedProblemId(prob.code)}>
                <td>
                  <Link
                    href={`/scoring/${prob.code}`}
                    className="flex justify-center"
                  >
                    <Image
                      src="/assets/svg/check.svg"
                      width={20}
                      height={20}
                      alt="採点へ"
                    />
                  </Link>
                </td>
                <td>
                  <ConditionalText value={prob.unchecked ?? 0} />
                  /
                  <ConditionalText
                    value={prob.unchecked_near_overdue ?? 0}
                    className="text-warning"
                  />
                  /
                  <ConditionalText
                    value={prob.unchecked_overdue ?? 0}
                    className="text-error"
                  />
                </td>
                <td>{prob.id}</td>
                <td>{prob.code}</td>
                <td>{prob.title}</td>
                {/* 文は 20文字まで */}
                <td>
                  {prob.body?.length ?? 0 > 20
                    ? `${prob.body?.slice(0, 20)}...`
                    : prob.body}
                </td>
                <td>{prob.point}</td>
                <td>{prob.solved_criterion}</td>
                <td>{prob.previous_problem_id}</td>
                <td>{prob.author_id === user?.id ? "自分" : ""}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="divider my-0" />
      {problem != null && (
        <div className="container-ictsc">
          <div className="problem-info flex flex-row items-end py-12">
            <h1 className="title-ictsc pr-4">{problem.title}</h1>
            満点
            {problem.point} pt 採点基準
            {problem.solved_criterion} pt
            <Link
              href={`/scoring/${problem.code}`}
              className="link link-primary pl-2"
            >
              採点する
            </Link>
          </div>
          <ICTSCCard>
            <MarkdownPreview
              className="problem-preview"
              content={problem.body}
            />
          </ICTSCCard>
        </div>
      )}
    </>
  );
}

export default Index;
