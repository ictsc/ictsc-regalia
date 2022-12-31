import useSWR from 'swr'

import {useApi} from "./api";
import {AuthSelfResult, Result} from "../types/_api";

export const useAuth = () => {
  const {apiClient} = useApi();

  const fetcher = (url: string) => apiClient.get(url).json<Result<AuthSelfResult>>()

  const {data, mutate, error} = useSWR('auth/self', fetcher)

  const loading = !data && !error

  return {
    user: data?.data?.user ?? null,
    loading,
    mutate
  }
}