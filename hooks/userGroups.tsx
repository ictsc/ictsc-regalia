import useSWR from "swr";

import { useApi } from "@/hooks/api";
import { UserGroup } from "@/types/UserGroup";
import { Result } from "@/types/_api";

export const useUserGroups = () => {
  const { apiClient } = useApi();

  const fetcher = (url: string) =>
    apiClient.get(url).json<Result<UserGroup[]>>();

  const { data, isLoading } = useSWR("usergroups", fetcher);

  return {
    userGroups: data?.data ?? null,
    isLoading,
  };
};
