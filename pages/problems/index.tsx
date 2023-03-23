import Link from "next/link";

import { useRecoilState } from "recoil";

import CommonLayout from "@/layouts/CommonLayout";
import { shortRule } from "@/components/_const";
import ProblemCard from "@/components/ProblemCard";
import LoadingPage from "@/components/LoadingPage";
import ICTSCCard from "@/components/Card";
import MarkdownPreview from "@/components/MarkdownPreview";
import NotificationCard from "@/components/NotificationCard";
import { useProblems } from "@/hooks/problem";
import { useNotice } from "@/hooks/notice";
import { dismissNoticeIdsState } from "@/hooks/state/recoil";

const Problems = () => {
  const [dismissNoticeIds, setDismissNoticeIds] = useRecoilState(
    dismissNoticeIdsState
  );

  const { problems, isLoading } = useProblems();
  const { notices, isLoading: isNoticeLoading } = useNotice();

  if (isLoading || isNoticeLoading) {
    return (
      <CommonLayout title={`問題一覧`}>
        <LoadingPage />
      </CommonLayout>
    );
  }

  return (
    <CommonLayout title={`問題一覧`}>
      {shortRule != "" && (
        <div className={"container-ictsc"}>
          <ICTSCCard className={"pt-4 pb-8"}>
            <MarkdownPreview content={shortRule?.replace(/\\n/g, "\n") ?? ""} />
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
      <div
        className={
          "container-ictsc flex flex-row justify-end text-primary font-bold"
        }
      >
        <Link href={"/notices"} className={"link link-hover"}>
          おしらせ一覧へ→
        </Link>
      </div>
      <ul
        className={
          "grid grid-cols-2 md:grid-cols-4 lg:grid-cols-4 gap-8 container-ictsc"
        }
      >
        {problems.map((problem, index) => (
          <li key={problem.id}>
            <ProblemCard key={index + 1} problem={problem} />
          </li>
        ))}
      </ul>
    </CommonLayout>
  );
};

export default Problems;
