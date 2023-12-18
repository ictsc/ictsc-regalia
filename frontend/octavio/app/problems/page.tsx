"use client";

import React from "react";

import Link from "next/link";

import { useRecoilState } from "recoil";

import { shortRule } from "@/components/_const";
import ICTSCCard from "@/components/card";
import LoadingPage from "@/components/loading-page";
import MarkdownPreview from "@/components/markdown-preview";
import NotificationCard from "@/components/notification-card";
import ProblemCard from "@/components/problem-card";
import ICTSCTitle from "@/components/title";
import useNotice from "@/hooks/notice";
import useProblems from "@/hooks/problems";
import { dismissNoticeIdsState } from "@/hooks/state/recoil";

function Problems() {
  const [dismissNoticeIds, setDismissNoticeIds] = useRecoilState(
    dismissNoticeIdsState,
  );

  const { problems, isLoading } = useProblems();
  const { notices, isLoading: isNoticeLoading } = useNotice();

  if (isLoading || isNoticeLoading) {
    return (
      <>
        <ICTSCTitle title="問題一覧" />
        <LoadingPage />
      </>
    );
  }

  const normalProblems = problems.filter(
    (problem) => problem.type === "normal",
  );
  const multipleProblems = problems.filter(
    (problem) => problem.type === "multiple",
  );

  return (
    <>
      <ICTSCTitle title="問題一覧" />
      <main>
        {shortRule !== "" && (
          <div className="container-ictsc">
            <ICTSCCard className="pt-4 pb-8">
              <MarkdownPreview content={shortRule} />
            </ICTSCCard>
          </div>
        )}
        {notices
          ?.filter((notice) => !dismissNoticeIds.includes(notice.source_id))
          .map((notice) => (
            <NotificationCard
              key={notice.source_id}
              notice={notice}
              onDismiss={() => {
                // sourceId の重複を排除しくっつける
                setDismissNoticeIds([
                  ...new Set([...dismissNoticeIds, notice.source_id]),
                ]);
              }}
            />
          ))}
        <div className="container-ictsc flex flex-row justify-end text-primary font-bold">
          <Link href="/notices" className="notice-link link link-hover">
            おしらせ一覧へ→
          </Link>
        </div>
        {normalProblems.length > 0 && (
          <>
            <h2 className="container-ictsc text-2xl font-bold">実技問題</h2>
            <ul className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-4 gap-8 container-ictsc">
              {normalProblems.map((problem) => (
                <li key={problem.id}>
                  <ProblemCard problem={problem} />
                </li>
              ))}
            </ul>
          </>
        )}
        {multipleProblems.length > 0 && (
          <>
            <h2 className="container-ictsc text-2xl font-bold pt-2">
              選択問題
            </h2>
            <ul className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-4 gap-8 container-ictsc">
              {multipleProblems.map((problem) => (
                <li key={problem.id}>
                  <ProblemCard problem={problem} />
                </li>
              ))}
            </ul>
          </>
        )}
      </main>
    </>
  );
}

export default Problems;
