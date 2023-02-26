import ky from "ky";

export const useApi = () => {
  let headers: HeadersInit = {
    Accept: "application/json, */*",
    "Content-type": "application/json",
  };

  const apiClient = ky.create({
    prefixUrl: process.env.NEXT_PUBLIC_API_URL,
    headers: headers,
    credentials: "include",
    throwHttpErrors: false,
  });

  return { apiClient };
};
