import { User } from "./User";

// {"code":200,"data":{"user":{"id":"105e93ed-6df2-43aa-9a38-2b68530afcff","created_at":"2023-02-12T11:52:28.489Z","updated_at":"2023-02-12T11:52:28.489Z","name":"admin","display_name":"admin","user_group_id":"201f0987-fd92-474e-b2d1-92eec2ee2abb","user_group":{"id":"201f0987-fd92-474e-b2d1-92eec2ee2abb","created_at":"2023-02-12T11:52:28.341Z","updated_at":"2023-02-12T11:52:28.341Z","name":"admin-group","organization":"admin-org","is_full_access":true,"bastion":{"id":"69cc8746-8fb4-4856-b865-f21a928455f3","created_at":"2023-02-12T11:52:28.494Z","updated_at":"2023-02-12T11:52:28.494Z","user_group_id":"201f0987-fd92-474e-b2d1-92eec2ee2abb","bastion_user":"bastion","bastion_password":"password","bastion_host":"bastion","bastion_port":22}},"is_read_only":false}}}

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
