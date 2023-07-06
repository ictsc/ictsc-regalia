import { testUserGroup, UserGroup } from "@/types/UserGroup";

export type Answer = {
  id: string;
  body: string;
  point: number | null;
  problem_id: string;
  user_group: UserGroup;
  created_at: string;
  updated_at: string;
};

export const testAnswer: Answer = {
  id: "1",
  body: "test",
  point: 100,
  problem_id: "1",
  user_group: testUserGroup,
  created_at: "2021-01-01",
  updated_at: "2021-01-01",
};
