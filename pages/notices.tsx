import CommonLayout from "@/layouts/CommonLayout";
import MarkdownPreview from "@/components/MarkdownPreview";
import LoadingPage from "@/components/LoadingPage";
import { useNotice } from "@/hooks/notice";

const Notices = () => {
  const { notices, isLoading } = useNotice();

  if (isLoading) {
    return (
      <CommonLayout title={"通知一覧"}>
        <LoadingPage />
      </CommonLayout>
    );
  }

  return (
    <CommonLayout title={"通知一覧"}>
      {notices?.map((notice) => (
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
                </div>
                <MarkdownPreview content={notice.body ?? ""} />
              </div>
            </div>
          </div>
        </div>
      ))}
    </CommonLayout>
  );
};

export default Notices;
