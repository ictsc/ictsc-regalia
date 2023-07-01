import useSWR from "swr";

import useApi from "@/hooks/api";
import { AuthSelfResult } from "@/types/_api";

const useAuth = () => {
  const { client } = useApi();

  const fetcher = (url: string) => client.get<AuthSelfResult>(url);

  const { data, mutate, isLoading } = useSWR("auth/self", fetcher);
  const logout = () => client.delete("auth/logout");

  return {
    user: data?.data?.user ?? null,
    logout,
    isLoading,
    mutate,
  };
};

export default useAuth;
