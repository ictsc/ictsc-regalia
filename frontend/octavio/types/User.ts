import { testUserGroup, UserGroup } from "@/types/UserGroup";

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

export const testUser: User = {
  id: "1",
  name: "Test User",
  display_name: "Test User",
  is_read_only: false,
  created_at: new Date(),
  updated_at: new Date(),
  user_group: {
    id: "1",
    name: "Test User Group",
    organization: "Test Organization",
    created_at: new Date(),
    updated_at: new Date(),
    is_full_access: false,
    members: [],
    bastion: null,
  },
  user_group_id: "1",
  profile: {
    id: "1",
    twitter_id: "test",
    github_id: "test",
    facebook_id: "test",
    self_introduction: "test",
    created_at: new Date(),
    updated_at: new Date(),
  },
  user_profile: {
    id: "1",
    twitter_id: "test",
    github_id: "test",
    facebook_id: "test",
    self_introduction: "test",
    created_at: new Date(),
    updated_at: new Date(),
  },
};

export const testAdminUser: User = {
  id: "1",
  name: "Test Admin User",
  display_name: "Test Admin User",
  is_read_only: false,
  created_at: new Date(),
  updated_at: new Date(),
  user_group: testUserGroup,
  user_group_id: "1",
  profile: {
    id: "1",
    twitter_id: "test",
    github_id: "test",
    facebook_id: "test",
    self_introduction: "test",
    created_at: new Date(),
    updated_at: new Date(),
  },
  user_profile: {
    id: "1",
    twitter_id: "test",
    github_id: "test",
    facebook_id: "test",
    self_introduction: "test",
    created_at: new Date(),
    updated_at: new Date(),
  },
};
