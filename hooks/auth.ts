import useSWR from "swr";

import useApi from "@/hooks/api";
import { SignInRequest } from "@/types/SignInRequest";
import { AuthSelfResult } from "@/types/_api";

const useAuth = () => {
  const { client } = useApi();

  const fetcher = (url: string) => client.get<AuthSelfResult>(url);

  const { data, mutate, isLoading } = useSWR("auth/self", fetcher);

  const signIn = async (request: SignInRequest) =>
    client.post("auth/signin", request);
  const logout = async () => client.delete("auth/signout");

  return {
    user: data?.data?.user ?? null,
    signIn,
    logout,
    isLoading,
    mutate,
  };
};

export default useAuth;
