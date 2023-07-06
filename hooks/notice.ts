import useSWR from "swr";

import useApi from "@/hooks/api";
import { Notice } from "@/types/Notice";

const useNotice = () => {
  const { client } = useApi();

  /* c8 ignore next */
  const fetcher = (url: string) => client.get<Notice[]>(url);

  const { data, mutate, isLoading } = useSWR("notices", fetcher);

  return {
    notices: data?.data ?? null,
    mutate,
    isLoading,
  };
};

export default useNotice;
