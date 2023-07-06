import useSWR from "swr";

import useApi from "@/hooks/api";
import { GetReCreateInfo } from "@/types/ReCreate";

const useReCreateInfo = (problemCode: string | null) => {
  const { client } = useApi();

  /* c8 ignore next */
  const fetcher = (url: string) => client.get<GetReCreateInfo>(url);

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

export default useReCreateInfo;
