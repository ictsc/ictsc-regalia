import { testUserGroup, UserGroup } from "@/types/UserGroup";

export type Rank = {
  rank: number;
  point: number;
  user_group: UserGroup;
  user_group_id: string;
};

export const testRank: Rank = {
  rank: 256,
  point: 100,
  user_group: testUserGroup,
  user_group_id: "1",
};
