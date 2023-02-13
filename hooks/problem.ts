import useSWR from "swr";

import matter from 'gray-matter';

import {useApi} from "./api";
import {ProblemResult, Result} from "../types/_api";
import {Matter, Problem} from "../types/Problem";

export const useProblems = () => {
  const {apiClient} = useApi();

  const fetcher = (url: string) => apiClient.get(url).json<Result<ProblemResult>>();

  const {data, mutate, isLoading} = useSWR('problems', fetcher);


  const getProblem = (code: string): [Matter | null, Problem | null] => {
    const problem = data?.data?.problems.find(problem => problem.code === code) ?? null;

    if (problem === null) {
      return [null, null];
    }

    const matterResult = matter(problem.body ?? '');

    const matterData = matterResult.data as Matter;
    const newProblem = {...problem, body: matterResult.content};

    console.log(problem)
    console.log(matterData);

    return [matterData, newProblem];
  }


  return {
    problems: data?.data?.problems ?? [],
    mutate,
    isLoading,
    getProblem,
  };
}