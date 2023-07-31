"use client";

import LoadingPage from "@/components/LoadingPage";
import NotificationCard from "@/components/NotificationCard";
import useNotice from "@/hooks/notice";

function Page() {
  const { notices, isLoading } = useNotice();

  if (isLoading) {
    return <LoadingPage />;
  }

  return notices?.map((notice) => (
    <NotificationCard key={notice.source_id} notice={notice} />
  ));
}

export default Page;
