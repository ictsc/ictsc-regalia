import useSWR from 'swr'
import {useApi} from "./api";
import {AuthSelfResult} from "../types/_api";

type Result<T> = {
  code: number;
  data: T;
}

export const useAuth = () => {
  const {apiClient} = useApi();

  const fetcher = (url: string) => apiClient.get(url).json<Result<AuthSelfResult>>()

  const {data, mutate} = useSWR('auth/self', fetcher)

  return {
    user: data?.data.user ?? null,
    mutate
  }
}