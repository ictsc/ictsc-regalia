import { type Transport, createClient } from "@connectrpc/connect";
import { ViewerService } from "@ictsc/proto/contestant/v1";

export type User = {
  name: string;
};

export async function fetchMe(transport: Transport): Promise<User | undefined> {
  const client = createClient(ViewerService, transport);
  const { viewer } = await client.getViewer({});
  if (viewer == null) {
    return;
  }
  return {
    name: viewer.name,
  };
}
