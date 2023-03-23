import useSWR from "swr";

import { useApi } from "@/hooks/api";
import { AuthSelfResult, Result } from "@/types/_api";

export const useAuth = () => {
  const { apiClient } = useApi();

  const fetcher = (url: string) =>
    apiClient.get(url).json<Result<AuthSelfResult>>();

  const { data, mutate, isLoading } = useSWR("auth/self", fetcher);

  return {
    user: data?.data?.user ?? null,
    isLoading,
    mutate,
  };
};
