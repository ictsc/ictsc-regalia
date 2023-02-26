import { UserGroup } from "./UserGroup";

export type Answer = {
  id: string;
  body: string;
  point: number | null;
  problem_id: string;
  user_group: UserGroup;
  created_at: string;
  updated_at: string;
};
