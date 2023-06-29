import useSWR from "swr";

import useApi from "@/hooks/api";
import { Rank } from "@/types/Rank";

type AnswerResult = {
  ranking: Rank[];
};

const useRanking = () => {
  const { client } = useApi();

  const fetcher = (url: string) => client.get<AnswerResult>(url);

  const { data, isLoading } = useSWR(`ranking`, fetcher);

  return {
    ranking: data?.data?.ranking ?? null,
    loading: isLoading,
  };
};

export default useRanking;
