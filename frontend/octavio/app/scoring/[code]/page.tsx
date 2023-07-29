"use client";

import Error from "next/error";
import { useRouter } from "next/router";

import { useForm } from "react-hook-form";

import ICTSCCard from "@/components/Card";
import LoadingPage from "@/components/LoadingPage";
import MarkdownPreview from "@/components/MarkdownPreview";
import ProblemConnectionInfo from "@/components/ProblemConnectionInfo";
import ProblemMeta from "@/components/ProblemMeta";
import ProblemTitle from "@/components/ProblemTitle";
import ScoringAnswerForm from "@/components/ScoringAnswerForm";
import useAnswers from "@/hooks/answer";
import useAuth from "@/hooks/auth";
import useProblem from "@/hooks/problem";
import BaseLayout from "@/layouts/BaseLayout";

type Input = {
  answerFilter: string;
};

function ScoringProblem() {
  const router = useRouter();
  const { code, answer_id: answerId } = router.query;

  const { register, watch } = useForm<Input>({
    defaultValues: {
      answerFilter: "2",
    },
  });

  const answerFilter = watch("answerFilter");

  const { user } = useAuth();
  const { problem, matter, isLoading } = useProblem(code as string);
  const { answers } = useAnswers(problem?.id ?? "");

  const isFullAccess = user?.user_group.is_full_access ?? false;
  const isReadOnly = user?.is_read_only ?? false;

  if (isLoading) {
    return (
      <BaseLayout title="採点">
        <LoadingPage />
      </BaseLayout>
    );
  }

  if (!isFullAccess || isReadOnly || problem === null) {
    return <Error statusCode={404} />;
  }

  return (
    <BaseLayout title={`採点(${problem.code} ${problem.title})`}>
      <div className="container-ictsc">
        <div className="flex flex-col mt-12">
          <ProblemTitle title={problem.title} />
          <ProblemMeta problem={problem} />
        </div>
        <ICTSCCard className="ml-0">
          <MarkdownPreview className="problem-body" content={problem.body} />
        </ICTSCCard>
        <div className="divider" />
        <ProblemConnectionInfo matter={matter} />
        <div className="flex flex-row justify-between mb-8 pt-2">
          <table className="table border table-compact">
            <thead>
              <tr>
                <th>未済点 ~15分</th>
                <th>15~19分</th>
                <th>20分~</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>{problem.unchecked}</td>
                <td>
                  {" "}
                  {problem.unchecked_near_overdue != null &&
                  problem.unchecked_near_overdue > 0 ? (
                    <div className="inline-block text-warning">
                      {problem.unchecked_near_overdue}
                    </div>
                  ) : (
                    <div className="inline-block">-</div>
                  )}
                </td>
                <td>
                  {problem.unchecked_overdue != null &&
                  problem.unchecked_overdue > 0 ? (
                    <div className="inline-block text-error">
                      {problem.unchecked_overdue}
                    </div>
                  ) : (
                    <div className="inline-block">-</div>
                  )}
                </td>
              </tr>
            </tbody>
          </table>
          <div className="form-control w-full max-w-[200px]">
            <div className="label">
              <span className="label-text">採点状況フィルタ</span>
            </div>
            <select
              {...register("answerFilter")}
              className="select select-sm select-bordered "
            >
              <option value={0}>すべて</option>
              <option value={1}>採点済みのみ</option>
              <option value={2}>未済点のみ</option>
            </select>
          </div>
        </div>
        {answers
          .filter((answer) => {
            if (answerFilter === "0") {
              return true;
            }
            if (answerFilter === "1") {
              return answer.point !== null;
            }
            return answer.point === null;
          })
          .sort((a, b) => {
            // date
            if (a.created_at > b.created_at) {
              return -1;
            }
            if (a.created_at < b.created_at) {
              return 1;
            }
            return 0;
          })
          .filter((answer) => {
            if (answerId == null) {
              return true;
            }
            return answer.id === answerId;
          })
          .map((answer) => (
            <ScoringAnswerForm
              key={answer.id}
              problem={problem}
              answer={answer}
            />
          ))}
      </div>
    </BaseLayout>
  );
}

export default ScoringProblem;
