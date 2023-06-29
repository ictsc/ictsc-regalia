import useSWR from "swr";

import useApi from "@/hooks/api";
import {UserGroup} from "@/types/UserGroup";

const useUserGroups = () => {
  const {client} = useApi();

  const fetcher = (url: string) => client.get<UserGroup[]>(url);

  const {data, isLoading} = useSWR("usergroups", fetcher);

  return {
    userGroups: data?.data ?? null,
    isLoading,
  };
};

export default useUserGroups;
