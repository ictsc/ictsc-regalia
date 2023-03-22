import { useState } from "react";

import Link from "next/link";
import Error from "next/error";

import BaseLayout from "@/layouts/BaseLayout";
import ICTSCCard from "@/components/Card";
import MarkdownPreview from "@/components/MarkdownPreview";
import LoadingPage from "@/components/LoadingPage";
import { useAuth } from "@/hooks/auth";
import { useProblem, useProblems } from "@/hooks/problem";

const Index = () => {
  const { user } = useAuth();
  const { problems, isLoading } = useProblems();
  const [selectedProblemId, setSelectedProblemId] = useState<string | null>(
    null
  );

  const { problem } = useProblem(selectedProblemId);

  const isFullAccess = user?.user_group.is_full_access ?? false;
  const isReadOnly = user?.is_read_only ?? false;

  if (isLoading) {
    return (
      <BaseLayout title={"採点"}>
        <LoadingPage />
      </BaseLayout>
    );
  }

  if (!isFullAccess || isReadOnly || problems === null) {
    return <Error statusCode={404} />;
  }

  return (
    <BaseLayout title={"採点"}>
      <div className="overflow-x-auto">
        <table className="table table-compact w-full">
          <thead>
            <tr>
              <th>採点</th>
              <th>未済点 ~15分/15~19分/20分~</th>
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
          <tbody className={"cursor-pointer"}>
            {problems.map((problem) => (
              <tr
                key={problem.id}
                onClick={() => setSelectedProblemId(problem.code)}
              >
                <td>
                  <Link
                    href={`/scoring/${problem.code}`}
                    className={"flex justify-center"}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                      strokeWidth={1.5}
                      stroke="currentColor"
                      className="w-4 h-4"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M4.5 12.75l6 6 9-13.5"
                      />
                    </svg>
                  </Link>
                </td>
                <td>
                  {problem.unchecked}/
                  {problem.unchecked_near_overdue != null &&
                  problem.unchecked_near_overdue > 0 ? (
                    <div className={"inline-block text-warning"}>
                      {problem.unchecked_near_overdue}
                    </div>
                  ) : (
                    <div className={"inline-block"}>-</div>
                  )}
                  /
                  {problem.unchecked_overdue != null &&
                  problem.unchecked_overdue > 0 ? (
                    <div className={"inline-block text-error"}>
                      {problem.unchecked_overdue}
                    </div>
                  ) : (
                    <div className={"inline-block"}>-</div>
                  )}
                </td>
                <td>{problem.id}</td>
                <td>{problem.code}</td>
                <td>{problem.title}</td>
                {/* 文は 20文字まで */}
                <td>
                  {problem.body?.length ?? 0 > 20
                    ? problem.body?.slice(0, 20) + "..."
                    : problem.body}
                </td>
                <td>{problem.point}</td>
                <td>{problem.solved_criterion}</td>
                <td>{problem.previous_problem_id}</td>
                <td>{problem.author_id === user?.id ? "自分" : ""}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className={"divider my-0"} />
      {problem != null && (
        <div className={"container-ictsc"}>
          <div className={"flex flex-row items-end py-12"}>
            <h1 className={"title-ictsc pr-4"}>{problem.title}</h1>
            満点
            {problem.point} pt 採点基準
            {problem.solved_criterion} pt
            <Link
              href={`/scoring/${problem.code}`}
              className={"link link-primary pl-2"}
            >
              採点する
            </Link>
          </div>
          <ICTSCCard>
            <MarkdownPreview content={problem.body ?? ""} />
          </ICTSCCard>
        </div>
      )}
    </BaseLayout>
  );
};

export default Index;
