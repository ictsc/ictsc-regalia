import { Response } from "next/dist/compiled/@edge-runtime/primitives";

import { Meta, StoryObj } from "@storybook/react";
import { http } from "msw";

import { apiUrl } from "@/components/_const";
import ICTSCNavBar from "@/components/navbar";
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
        http.get(
          path,
          (info) =>
            new Response(
              JSON.stringify({
                data: {
                  user: testAdminUser,
                },
              }),
              // {
              //   status: 200,
              //   headers: {
              //     "Content-Type": "application/json",
              //   },
              // },
            ),
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
        // http.get(path, (req, res, ctx) =>
        //   res(
        //     ctx.json({
        //       data: {
        //         user: testUser,
        //       },
        //     }),
        //   ),
        // ),
        http.get(
          path,
          (info) =>
            new Response(
              JSON.stringify({
                data: {
                  user: testUser,
                },
              }),
              // {
              //   status: 200,
              //   headers: {
              //     "Content-Type": "application/json",
              //   },
              // },
            ),
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
        // rest.get(path, (req, res, ctx) =>
        //   res(
        //     ctx.json({
        //       data: {
        //         user: null,
        //       },
        //     }),
        //   ),
        // ),
        http.get(
          path,
          (info) =>
            new Response(
              JSON.stringify({
                data: {
                  user: null,
                },
              }),
              // {
              //   status: 200,
              //   headers: {
              //     "Content-Type": "application/json",
              //   },
              // },
            ),
        ),
      ],
    },
  },
};
