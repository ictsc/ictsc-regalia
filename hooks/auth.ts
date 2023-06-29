import useSWR from "swr";

import useApi from "@/hooks/api";
import { AuthSelfResult } from "@/types/_api";

const useAuth = () => {
  const { client } = useApi();

  const fetcher = (url: string) => client.get<AuthSelfResult>(url);

  const { data, mutate, isLoading } = useSWR("auth/self", fetcher);

  return {
    user: data?.data?.user ?? null,
    isLoading,
    mutate,
  };
};

export default useAuth;
