import useSWR from "swr";

import { useApi } from "@/hooks/api";
import { Notice } from "@/types/Notice";
import { Result } from "@/types/_api";

export const useNotice = () => {
  const { apiClient } = useApi();

  const fetcher = (url: string) => apiClient.get(url).json<Result<Notice[]>>();

  const { data, mutate, isLoading } = useSWR("notices", fetcher);

  console.log(data?.data);

  return {
    notices: data?.data ?? null,
    mutate,
    isLoading,
  };
};
