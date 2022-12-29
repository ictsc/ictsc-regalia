import {User} from "./User";

export type UserGroup = {
  id: string;
  name: string;
  organization: string;
  created_at: Date;
  updated_at: Date;
  is_full_access: boolean;
  members: User[] | null;
}