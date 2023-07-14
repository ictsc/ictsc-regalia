import "zenn-content-css";

import { useState } from "react";

import Error from "next/error";
import Image from "next/image";
import { useRouter } from "next/router";

import { Toaster } from "react-hot-toast";

import AnswerForm from "@/components/AnswerForm";
import AnswerListSection from "@/components/AnswerListSection";
import ICTSCCard from "@/components/Card";
import LoadingPage from "@/components/LoadingPage";
import MarkdownPreview from "@/components/MarkdownPreview";
import ProblemConnectionInfo from "@/components/ProblemConnectionInfo";
import ProblemMeta from "@/components/ProblemMeta";
import ProblemTitle from "@/components/ProblemTitle";
import { answerLimit, recreateRule } from "@/components/_const";
import useAuth from "@/hooks/auth";
import useProblem from "@/hooks/problem";
import useReCreateInfo from "@/hooks/reCreateInfo";
import BaseLayout from "@/layouts/BaseLayout";

function ProblemPage() {
  const router = useRouter();
  const { problemId } = router.query;

  const { user } = useAuth();
  const { matter, problem, isLoading } = useProblem(problemId as string | null);
  const [isReCreateModalOpen, setIsReCreateModalOpen] = useState(false);

  const { recreateInfo, mutate, reCreate } = useReCreateInfo(
    problem?.code ?? null
  );

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
    return (
      <BaseLayout title="Loading...">
        <LoadingPage />
      </BaseLayout>
    );
  }

  if (problem === null) {
    return <Error statusCode={404} />;
  }

  return (
    <>
      <Toaster />
      <input type="checkbox" id="my-modal-5" className="modal-toggle" />
      <div className={`modal ${isReCreateModalOpen && "modal-open"}`}>
        <div className="modal-box container-ictsc">
          <h3 className="title-ictsc pt-4 pb-8">
            問題の再展開を行います。よろしいですか？
          </h3>
          <MarkdownPreview content={recreateRule.replace(/\\n/g, "\n")} />
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
      <BaseLayout title={`${problem.code} ${problem.title} 問題`}>
        <div className="container-ictsc">
          <div className="flex flex-row justify-between pt-12 justify-items-center">
            <ProblemTitle title={problem.title} />
            {!isReadOnly && (
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
          <ProblemConnectionInfo matter={matter} />
          {recreateInfo?.available != null &&
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
          <ICTSCCard className="mt-8">
            <MarkdownPreview content={problem.body ?? ""} />
          </ICTSCCard>

          {!isReadOnly && <AnswerForm />}
          {answerLimit && (
            <div className="text-sm pt-2">
              ※ 回答は{answerLimit}分に1度のみです
            </div>
          )}
          <div className="divider" />
          <AnswerListSection problem={problem} />
        </div>
      </BaseLayout>
    </>
  );
}

export default ProblemPage;
