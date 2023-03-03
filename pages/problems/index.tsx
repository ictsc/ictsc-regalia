import Head from "next/head";

import ICTSCNavBar from "../../components/Navbar";
import ProblemCard from "../../components/ProblemCard";
import LoadingPage from "../../components/LoadingPage";
import { useProblems } from "../../hooks/problem";
import { useNotice } from "../../hooks/notice";
import { shortRule, site } from "../../components/_const";
import ICTSCCard from "../../components/Card";
import MarkdownPreview from "../../components/MarkdownPreview";
import { dismissNoticeIdsState } from "../../hooks/state/recoil";
import { useRecoilState } from "recoil";
import Link from "next/link";

const Problems = () => {
  const [dismissNoticeIds, setDismissNoticeIds] = useRecoilState(
    dismissNoticeIdsState
  );

  const { problems, isLoading } = useProblems();
  const { notices, isLoading: isNoticeLoading } = useNotice();

  if (isLoading || isNoticeLoading) {
    return (
      <>
        <ICTSCNavBar />
        <LoadingPage />
      </>
    );
  }

  return (
    <>
      <Head>
        <title>問題一覧 - {site}</title>
      </Head>
      <ICTSCNavBar />
      <h1 className={"title-ictsc text-center py-12"}>問題一覧</h1>
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
          <div key={notice.source_id} className={"container-ictsc"}>
            <div className="bg-gray-200 p-4 rounded rounded-lg  shadow-lg grow">
              <div>
                <div className={"flex flex-col"}>
                  <div
                    className={
                      "flex flex-row justify-between justify-items-center"
                    }
                  >
                    <div className={"flex flex-row"}>
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        className="stroke-info flex-shrink-0 w-6 h-6"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        ></path>
                      </svg>
                      <span className={"pl-2 font-bold"}>{notice.title}</span>
                    </div>
                    <button
                      onClick={() => {
                        // sourceId の重複を排除しくっつける
                        setDismissNoticeIds([
                          ...new Set([...dismissNoticeIds, notice.source_id]),
                        ]);
                      }}
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke-width="1.5"
                        stroke="currentColor"
                        className="w-6 h-6"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          d="M6 18L18 6M6 6l12 12"
                        />
                      </svg>
                    </button>
                  </div>
                  <MarkdownPreview content={notice.body ?? ""} />
                </div>
              </div>
            </div>
          </div>
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
        {problems &&
          problems.map((problem, index) => (
            <li key={problem.id}>
              <ProblemCard index={index + 1} problem={problem} />
            </li>
          ))}
      </ul>
    </>
  );
};

export default Problems;
