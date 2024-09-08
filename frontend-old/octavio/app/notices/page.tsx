"use client";

import LoadingPage from "@/components/loading-page";
import NotificationCard from "@/components/notification-card";
import useNotice from "@/hooks/notice";

function Page() {
  const { notices, isLoading } = useNotice();

  if (isLoading) {
    return <LoadingPage />;
  }

  return (
    <main>
      {notices?.map((notice) => (
        <NotificationCard key={notice.source_id} notice={notice} />
      ))}
    </main>
  );
}

export default Page;
