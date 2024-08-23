import { useReducer } from "react";
import type { Meta, StoryObj } from "@storybook/react";
import { AppShellLayout as Layout } from "./layout";
import { HeaderView } from "./header";
import { NavbarView } from "./navbar";
import { ContestStateView } from "./contest-state";

function AppShell() {
  const [collapsed, toggle] = useReducer((o) => !o, false);
  return (
    <Layout
      header={
        <HeaderView
          contestState={
            <ContestStateView state="before" restDurationSeconds={73850} />
          }
        />
      }
      navbar={<NavbarView collapsed={collapsed} onOpenToggleClick={toggle} />}
      navbarCollapsed={collapsed}
    >
      <h1>Main</h1>
      <p>
        {`Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
            ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat. Duis aute irure dolor in
            reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
            pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
            culpa qui officia deserunt mollit anim id est laborum.`.repeat(100)}
      </p>
    </Layout>
  );
}

export default {
  title: "AppShell",
  component: AppShell,
} satisfies Meta<typeof AppShell>;

type Story = StoryObj<typeof AppShell>;

export const Default: Story = {
  render: () => <AppShell />,
};
