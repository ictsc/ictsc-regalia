import { User } from "@/types/User";

export type UserGroup = {
  id: string;
  name: string;
  organization: string;
  created_at: Date;
  updated_at: Date;
  is_full_access: boolean;
  members: User[] | null;
  bastion: {
    bastion_user: string;
    bastion_password: string;
    bastion_host: string;
    bastion_port: number;
  } | null;
};

export const testUserGroup: UserGroup = {
  id: "1",
  name: "Test User Group",
  organization: "Test Organization",
  created_at: new Date(),
  updated_at: new Date(),
  is_full_access: false,
  members: [],
  bastion: {
    bastion_user: "test_bastion_user",
    bastion_password: "test_bastion_password",
    bastion_host: "test_bastion_host",
    bastion_port: 22,
  },
};
