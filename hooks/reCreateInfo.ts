import useSWR from "swr";

import { useApi } from "@/hooks/api";
import { GetReCreateInfo } from "@/types/ReCreate";
import { Result } from "@/types/_api";

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
