export type GetReCreateInfo = {
  available: boolean;
  created_time: string;
  completed_time: string | null;
};

export const testReCreateInfo: GetReCreateInfo = {
  available: true,
  created_time: "2021-01-01T00:00:00.000Z",
  completed_time: null,
};
