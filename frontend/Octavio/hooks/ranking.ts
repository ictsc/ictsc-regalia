import useSWR from "swr";

import useApi from "@/hooks/api";
import { RankingResult } from "@/types/_api";

const useRanking = () => {
  const { client } = useApi();

  /* c8 ignore next */
  const fetcher = (url: string) => client.get<RankingResult>(url);

  const { data, isLoading } = useSWR(`ranking`, fetcher);

  return {
    ranking: data?.data?.ranking ?? null,
    loading: isLoading,
  };
};

export default useRanking;
