import {
  type Notice as ProtoNotice,
  NoticeService,
} from "@ictsc/proto/contestant/v1";
import { type Transport, createClient } from "@connectrpc/connect";

export type Notice = ProtoNotice;

export async function fetchNotices(transport: Transport): Promise<Notice[]> {
  const client = createClient(NoticeService, transport);
  const notices = await client.listNotices({});
  return notices.notices;
}
