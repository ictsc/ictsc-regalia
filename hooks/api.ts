import axios from "axios";

import {Result} from "@/types/_api";

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
  };

  // TODO(k-shir0): apiClient は廃止し client に統一する
  return {apiClient, client};
};

export default useApi;
