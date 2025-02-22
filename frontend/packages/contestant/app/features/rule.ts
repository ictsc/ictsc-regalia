import {
  Code,
  ConnectError,
  type Transport,
  createClient,
} from "@connectrpc/connect";
import { ContestService } from "@ictsc/proto/contestant/v1";

export type Rule = {
  markdown: string;
};

export async function fetchRule(transport: Transport): Promise<Rule> {
  const client = createClient(ContestService, transport);
  const { rule } = await client.getRule({});
  if (rule == null) {
    throw new ConnectError("rule not found", Code.NotFound);
  }
  return { markdown: rule.markdown };
}
