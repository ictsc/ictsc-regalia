import createClient, { type Client } from "openapi-fetch";
import type { paths } from "./schema";

export type ApiClient = Client<paths>;

export function createApiClient(baseUrl = ""): ApiClient {
  return createClient<paths>({
    baseUrl,
    credentials: "include",
    headers: {
      Accept: "application/json",
    },
  });
}

export const api = createApiClient();
