"use client";

import LoadingPage from "@/components/loading-page";
import useRanking from "@/hooks/ranking";

function Page() {
  const { ranking, loading } = useRanking();

  if (loading) {
    return <LoadingPage />;
  }

  return (
    <table className="table border w-full">
      <thead>
        <tr>
          <th>#</th>
          <th>チーム名</th>
          <th>所属</th>
          <th className="text-right">得点</th>
        </tr>
      </thead>
      <tbody>
        {ranking?.map((rank) => (
          <tr key={rank.user_group.id}>
            <td className="w-[64px]">{rank.rank}</td>
            <td>{rank.user_group.name}</td>
            <td>{rank.user_group.organization}</td>
            <td className="w-[124px] text-right">{rank.point}pt</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export default Page;
