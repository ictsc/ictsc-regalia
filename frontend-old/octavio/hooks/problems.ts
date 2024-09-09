import useSWR from "swr";

import useApi from "@/hooks/api";
import { Problem } from "@/types/Problem";
import { ProblemResult } from "@/types/_api";

const useProblems = () => {
  const { client } = useApi();

  const fetcher = (url: string) => client.get<ProblemResult>(url);

  const { data, mutate, isLoading, error } = useSWR("problems", fetcher);

  let problems: Problem[];
  if (error) {
    problems = [];
  } else {
    problems = data?.data?.problems ?? [];
  }

  return {
    problems,
    mutate,
    isLoading,
  };
};

export default useProblems;
