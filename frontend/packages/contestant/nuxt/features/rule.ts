import { expectData, type ApiClient } from "@ictsc/api";

export type Rule = {
  markdown: string;
};

export async function fetchRule(client: ApiClient): Promise<Rule> {
  const response = expectData(await client.GET("/api/v1/contestant/rule"));
  return { markdown: response.rule.markdown };
}
