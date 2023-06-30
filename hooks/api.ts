import axios from "axios";

import { Result } from "@/types/_api";

const useApi = () => {
  const headers = {
    Accept: "application/json, */*",
    "Content-type": "application/json",
  };

  const apiClient = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_URL,
    headers,
    withCredentials: true,
  });

  const client = {
    get: <T>(url: string) =>
      apiClient.get<Result<T>>(url).then((response) => response.data),
    post: (url: string, data?: any) =>
      apiClient.post(url, data).then((response) => response.data),
    put: (url: string, data?: any) =>
      apiClient.put(url, data).then((response) => response.data),
    patch: <T>(url: string, data?: any) =>
      apiClient.patch<Result<T>>(url, data).then((response) => response.data),
    delete: (url: string) =>
      apiClient.delete(url).then((response) => response.data),
  };

  return { client };
};

export default useApi;
