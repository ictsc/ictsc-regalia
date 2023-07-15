export type Notice = {
  source_id: string;
  title: string;
  body?: string;
  draft: boolean;
};

export const testNotice: Notice = {
  source_id: "1",
  title: "テスト通知タイトル",
  body: "テスト通知本文",
  draft: false,
};
