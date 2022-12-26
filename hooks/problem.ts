import useSWR from "swr";

import {useApi} from "./api";
import {ProblemResult, Result} from "../types/_api";

export const useProblems = () => {
  const {apiClient} = useApi();

  const fetcher = (url: string) => apiClient.get(url).json<Result<ProblemResult>>();

  const {data, mutate} = useSWR('problems', fetcher);

  return {
    problems: data?.data?.problems ?? [],
    mutate,
  };
}