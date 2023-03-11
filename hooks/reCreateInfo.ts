import { useApi } from "./api";
import { Result } from "../types/_api";
import { GetReCreateInfo } from "../types/ReCreate";
import useSWR from "swr";

export const useReCreateInfo = (problemCode: string | null) => {
  const { apiClient } = useApi();

  const fetcher = (url: string) =>
    apiClient.get(url).json<Result<GetReCreateInfo>>();

  const { data, isLoading, mutate } = useSWR(
    problemCode && `recreate/${problemCode}`,
    fetcher,
    {
      refreshInterval: 30000,
    }
  );

  return {
    recreateInfo: data?.data ?? null,
    isLoading,
    mutate,
  };
};
