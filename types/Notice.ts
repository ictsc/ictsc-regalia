export type Notice = {
  source_id: string;
  title: string;
  body?: string;
  draft: boolean;
};

export const testNotice: Notice = {
  source_id: "1",
  title: "Test Notice",
  body: "Test Notice Body",
  draft: false,
};
