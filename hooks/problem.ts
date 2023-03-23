import matter from "gray-matter";
import useSWR from "swr";

import { useApi } from "@/hooks/api";
import { Matter } from "@/types/Problem";
import { ProblemResult, Result } from "@/types/_api";

export const useProblems = () => {
  const { apiClient } = useApi();

  const fetcher = (url: string) =>
    apiClient.get(url).json<Result<ProblemResult>>();

  const { data, mutate, isLoading } = useSWR("problems", fetcher);

  return {
    problems: data?.data?.problems ?? [],
    mutate,
    isLoading,
  };
};

export const useProblem = (code: string | null) => {
  const { problems } = useProblems();

  const problem = problems.find((problem) => problem.code === code) ?? null;

  if (problem === null) {
    return {
      matter: null,
      problem: null,
      isLoading: false,
    };
  }

  const matterResult = matter(problem.body ?? "");
  const matterData = matterResult.data as Matter;
  const newProblem = { ...problem, body: matterResult.content };

  return {
    matter: matterData,
    problem: newProblem,
    isLoading: false,
  };
};
