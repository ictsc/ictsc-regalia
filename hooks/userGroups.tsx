import useSWR from "swr";

import {useApi} from "./api";
import {Result} from "../types/_api";
import {UserGroup} from "../types/UserGroup";

export const useUserGroups = () => {
  const {apiClient} = useApi();

  const fetcher = (url: string) => apiClient.get(url).json<Result<UserGroup[]>>();

  const {data, error} = useSWR('usergroups', fetcher)

  const loading = !data && !error;

  return {
    userGroups: data?.data ?? [],
    loading,
  }
}