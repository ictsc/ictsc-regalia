"use client";

import "zenn-content-css";

import React, { useState } from "react";

import Image from "next/image";
import { notFound } from "next/navigation";

import clsx from "clsx";
import { Toaster } from "react-hot-toast";

import AnswerForm from "@/app/problems/[problemId]/_components/answer-form";
import AnswerListSection from "@/app/problems/[problemId]/_components/answer-list-section";
import MultipleAnswerForm from "@/app/problems/[problemId]/_components/multiple-answer-form";
import { answerLimit, recreateRule } from "@/components/_const";
import ICTSCCard from "@/components/card";
import LoadingPage from "@/components/loading-page";
import MarkdownPreview from "@/components/markdown-preview";
import ProblemConnectionInfo from "@/components/problem-connection-info";
import ProblemMeta from "@/components/problem-meta";
import ProblemTitle from "@/components/problem-title";
import useAuth from "@/hooks/auth";
import useProblem from "@/hooks/problem";
import useReCreateInfo from "@/hooks/reCreateInfo";

function ProblemPage({ params }: { params: { problemId: string } }) {
  const { user } = useAuth();
  const { matter, problem, isLoading } = useProblem(params.problemId);
  const [isReCreateModalOpen, setIsReCreateModalOpen] = useState(false);

  const { recreateInfo, mutate, reCreate } = useReCreateInfo(
    problem?.code ?? null,
  );
  const formType = problem?.type ?? "normal";
  const isReadOnly = user?.is_read_only ?? true;

  const onReCreateSubmit = async () => {
    /* c8 ignore next 404 のところでチェックしているため必ず false になる */
    if (problem === null) {
      return;
    }

    const response = await reCreate(problem.code);

    if (response.code === 200) {
      await mutate();
    }
  };

  if (isLoading) {
    return <LoadingPage />;
  }

  if (problem === null) {
    return notFound();
  }

  return (
    <>
      <Toaster />
      <input type="checkbox" id="my-modal-5" className="modal-toggle" />
      <div className={clsx("modal", isReCreateModalOpen && "modal-open")}>
        <div className="modal-box container-ictsc">
          <h3 className="title-ictsc pt-4 pb-8">
            問題の再展開を行います。よろしいですか？
          </h3>
          <MarkdownPreview content={recreateRule} />
          <div className="modal-action">
            <button
              type="button"
              onClick={() => setIsReCreateModalOpen(false)}
              className="btn btn-link"
              data-testid="modal-close-btn"
            >
              閉じる
            </button>
            <button
              type="button"
              onClick={async () => {
                onReCreateSubmit();
                setIsReCreateModalOpen(false);
              }}
              className="btn btn-primary"
            >
              問題の再展開を行う
            </button>
          </div>
        </div>
      </div>
      <div className="container-ictsc">
        <div className="flex flex-row justify-between pt-12 justify-items-center">
          <ProblemTitle title={problem.title} />
          {!isReadOnly && formType === "normal" && (
            <button
              type="button"
              className="btn text-red-500 btn-sm"
              onClick={() => {
                setIsReCreateModalOpen(true);
              }}
              disabled={
                recreateInfo?.available != null &&
                !(recreateInfo?.available ?? false)
              }
            >
              再展開を行う
            </button>
          )}
        </div>
        <ProblemMeta problem={problem} />
        {formType === "normal" && <ProblemConnectionInfo matter={matter} />}
        {formType === "normal" &&
          recreateInfo?.available != null &&
          !(recreateInfo?.available ?? false) && (
            <div className="alert alert-info shadow-lg grow">
              <div>
                <div className="animate-spin">
                  <Image
                    src="/assets/svg/arrow-path.svg"
                    height={24}
                    width={24}
                    alt="recreate"
                  />
                </div>
                <div className="flex flex-col">
                  <span>問題を再展開中です</span>
                </div>
              </div>
            </div>
          )}
        <ICTSCCard
          className={clsx(
            formType === "normal" && "mt-8",
            formType === "multiple" && "mt-4",
          )}
        >
          <MarkdownPreview className="problem-body" content={problem.body} />
        </ICTSCCard>
        {!isReadOnly && formType === "normal" && (
          <AnswerForm code={params.problemId} />
        )}
        {!isReadOnly && formType === "multiple" && (
          <MultipleAnswerForm code={params.problemId} />
        )}
        {answerLimit && (
          <div className="text-sm pt-2">
            ※ 回答は{answerLimit}分に1度のみです
          </div>
        )}
        {formType !== "multiple" && (
          <>
            <div className="divider" />
            <AnswerListSection problem={problem} />
          </>
        )}
      </div>
    </>
  );
}

export default ProblemPage;
