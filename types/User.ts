import { UserGroup } from "./UserGroup";

export type User = {
  id: string;
  name: string;
  display_name: string;
  is_read_only: boolean;
  created_at: Date;
  updated_at: Date;
  user_group: UserGroup;
  user_group_id: string;
  // ユーザーグループ経由で取得した時のみ使われる
  profile: Profile | null;
  user_profile: Profile | null;
};

export type Profile = {
  id: string;
  twitter_id: string;
  github_id: string;
  facebook_id: string;
  self_introduction: string;
  created_at: Date;
  updated_at: Date;
};
