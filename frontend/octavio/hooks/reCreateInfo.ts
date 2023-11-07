import useSWR from "swr";

import useApi from "@/hooks/api";
import { GetReCreateInfo } from "@/types/ReCreate";
import useProblem from "@/hooks/problem";

const useReCreateInfo = (problemCode: string | null) => {
  const { client } = useApi();

  /* c8 ignore next */
  const fetcher = (url: string) => client.get<GetReCreateInfo>(url);

  const { data, isLoading, mutate } = useSWR(
    problemCode && `recreate/${problemCode}`,
    fetcher,
    {
      refreshInterval: 30000,
    },
  );

  const reCreate = async (code: string) => client.post(`recreate/${code}`);

  return {
    recreateInfo: data?.data ?? null,
    isLoading,
    mutate,
    reCreate,
  };
};

export default useReCreateInfo;
