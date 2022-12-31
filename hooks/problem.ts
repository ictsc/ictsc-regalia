import useSWR from "swr";

import {useApi} from "./api";
import {ProblemResult, Result} from "../types/_api";
import {Problem} from "../types/Problem";

export const useProblems = () => {
  const {apiClient} = useApi();

  const fetcher = (url: string) => apiClient.get(url).json<Result<ProblemResult>>();

  const {data, mutate, isLoading} = useSWR('problems', fetcher);


  const getProblem = (code: string): Problem | null => {
    return data?.data?.problems.find(problem => problem.code === code) ?? null;
  }

  return {
    problems: data?.data?.problems ?? [],
    mutate,
    isLoading,
    getProblem
  };
}