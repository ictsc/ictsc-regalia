import CommonLayout from "@/layouts/CommonLayout";
import LoadingPage from "@/components/LoadingPage";
import NotificationCard from "@/components/NotificationCard";
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
        <NotificationCard key={notice.source_id} notice={notice} />
      ))}
    </CommonLayout>
  );
};

export default Notices;
