import { UserGroup } from "./UserGroup";

export type Rank = {
  rank: number;
  point: number;
  user_group: UserGroup;
  user_group_id: string;
};
