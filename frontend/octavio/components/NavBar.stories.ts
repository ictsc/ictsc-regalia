import { Meta, StoryObj } from "@storybook/react";
import { rest } from "msw";

import ICTSCNavBar from "@/components/Navbar";
import { apiUrl } from "@/components/_const";
import { testAdminUser, testUser } from "@/types/User";

const meta = {
  title: "Components/NavBar",
  component: ICTSCNavBar,
  parameters: {
    nextjs: {
      appDirectory: true,
    },
    layout: "fullscreen",
  },
} satisfies Meta<typeof ICTSCNavBar>;

export default meta;
type Story = StoryObj<typeof meta>;

const path = `${apiUrl}/auth/self`;

export const AdminLoggedIn: Story = {
  name: "管理者としてログイン",
  parameters: {
    msw: {
      handlers: [
        rest.get(path, (req, res, ctx) =>
          res(
            ctx.json({
              data: {
                user: testAdminUser,
              },
            })
          )
        ),
      ],
    },
  },
};

export const UserLoggedIn: Story = {
  name: "参加者としてログイン",
  parameters: {
    msw: {
      handlers: [
        rest.get(path, (req, res, ctx) =>
          res(
            ctx.json({
              data: {
                user: testUser,
              },
            })
          )
        ),
      ],
    },
  },
};

export const LoggedOut: Story = {
  name: "ログアウト状態",
  parameters: {
    msw: {
      handlers: [
        rest.get(path, (req, res, ctx) =>
          res(
            ctx.json({
              data: {
                user: null,
              },
            })
          )
        ),
      ],
    },
  },
};
