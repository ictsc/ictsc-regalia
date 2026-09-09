import { api } from "@ictsc/api";
import { fetchProblems } from "~/features/problem";
import { fetchNotices } from "~/features/announce";
import { fetchRanking } from "~/features/ranking";
import { fetchSchedule } from "~/features/schedule/feature";
export function useCompetition() {
  return useAsyncData("competition", async () => {
    const [problems, notices, ranking, schedule] = await Promise.all([
      fetchProblems(api),
      fetchNotices(api),
      fetchRanking(api),
      fetchSchedule(api),
    ]);
    return { problems, notices, ranking, schedule };
  });
}
