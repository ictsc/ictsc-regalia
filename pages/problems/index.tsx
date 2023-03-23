import Link from "next/link";

import { useRecoilState } from "recoil";

import ICTSCCard from "@/components/Card";
import LoadingPage from "@/components/LoadingPage";
import MarkdownPreview from "@/components/MarkdownPreview";
import NotificationCard from "@/components/NotificationCard";
import ProblemCard from "@/components/ProblemCard";
import { shortRule } from "@/components/_const";
import useNotice from "@/hooks/notice";
import { useProblems } from "@/hooks/problem";
import { dismissNoticeIdsState } from "@/hooks/state/recoil";
import CommonLayout from "@/layouts/CommonLayout";

function Problems() {
  const [dismissNoticeIds, setDismissNoticeIds] = useRecoilState(
    dismissNoticeIdsState
  );

  const { problems, isLoading } = useProblems();
  const { notices, isLoading: isNoticeLoading } = useNotice();

  if (isLoading || isNoticeLoading) {
    return (
      <CommonLayout title="問題一覧">
        <LoadingPage />
      </CommonLayout>
    );
  }

  return (
    <CommonLayout title="問題一覧">
      {shortRule !== "" && (
        <div className="container-ictsc">
          <ICTSCCard className="pt-4 pb-8">
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
      <div className="container-ictsc flex flex-row justify-end text-primary font-bold">
        <Link href="/notices" className="link link-hover">
          おしらせ一覧へ→
        </Link>
      </div>
      <ul className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-4 gap-8 container-ictsc">
        {problems.map((problem) => (
          <li key={problem.id}>
            <ProblemCard problem={problem} />
          </li>
        ))}
      </ul>
    </CommonLayout>
  );
}

export default Problems;
