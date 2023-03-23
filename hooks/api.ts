import ky from "ky";

const useApi = () => {
  const headers = {
    Accept: "application/json, */*",
    "Content-type": "application/json",
  };

  const apiClient = ky.create({
    prefixUrl: process.env.NEXT_PUBLIC_API_URL,
    headers,
    credentials: "include",
    throwHttpErrors: false,
  });

  return { apiClient };
};

export default useApi;
