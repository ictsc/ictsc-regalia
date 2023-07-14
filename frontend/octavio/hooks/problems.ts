import useSWR from "swr";

import useApi from "@/hooks/api";
import { ProblemResult } from "@/types/_api";

const useProblems = () => {
  const { client } = useApi();

  const fetcher = (url: string) => client.get<ProblemResult>(url);

  const { data, mutate, isLoading } = useSWR("problems", fetcher);

  return {
    problems: data?.data?.problems ?? [],
    mutate,
    isLoading,
  };
};

export default useProblems;
