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
