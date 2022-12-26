import {UserGroup} from "./UserGroup";

export type User = {
  id: string;
  name: string;
  display_name: string;
  is_read_only: boolean;
  created_at: Date;
  updated_at: Date;
  user_group: UserGroup;
  user_group_id: string;
}