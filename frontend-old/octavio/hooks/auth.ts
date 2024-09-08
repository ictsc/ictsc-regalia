import useSWR from "swr";

import useApi from "@/hooks/api";
import { PutProfileRequest } from "@/types/PutProfileRequest";
import { SignInRequest } from "@/types/SignInRequest";
import { User } from "@/types/User";
import { AuthSelfResult, SignUpRequest } from "@/types/_api";

const useAuth = () => {
  const { client } = useApi();

  const fetcher = (url: string) => client.get<AuthSelfResult>(url);

  const { data, mutate, isLoading, error } = useSWR("auth/self", fetcher);

  let user: User | null;
  if (error) {
    user = null;
  } else {
    user = data?.data?.user ?? null;
  }

  const signUp = async (request: SignUpRequest) =>
    client.post("users", request);
  const signIn = async (request: SignInRequest) =>
    client.post("auth/signin", request);
  const logout = async () => client.delete("auth/signout");
  const putProfile = async (userId: string, request: PutProfileRequest) =>
    client.put(`users/${userId}`, request);

  return {
    user,
    signUp,
    signIn,
    logout,
    putProfile,
    isLoading,
    mutate,
  };
};

export default useAuth;
