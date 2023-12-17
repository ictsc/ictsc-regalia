import useSWR from "swr";

import { preRoundMode } from "@/components/_const";
import useApi from "@/hooks/api";
import { RankingResult } from "@/types/_api";

const useRanking = () => {
  const { client } = useApi();

  /* c8 ignore next */
  const fetcher = (url: string) => client.get<RankingResult>(url);

  const { data, isLoading } = useSWR(preRoundMode ? null : `ranking`, fetcher);

  return {
    ranking: data?.data?.ranking ?? null,
    loading: isLoading,
  };
};

export default useRanking;
