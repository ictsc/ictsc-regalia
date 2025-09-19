import { useReducer } from "react";
import type { Meta, StoryObj } from "@storybook/react";
import { action } from "storybook/actions";
import { Layout } from "./layout";
import { Header } from "./header";
import { Navbar } from "./navbar";
import { ContestStateView } from "./contest-state";
import { AccountMenu } from "./account-menu";

const signOutAction = action("sign-out");

function AppShell() {
  const [collapsed, toggle] = useReducer((o) => !o, false);
  return (
    <Layout
      header={
        <Header
          contestState={
            <ContestStateView state="before" restDurationSeconds={73850} />
          }
          accountMenu={<AccountMenu name="Alice" onSignOut={signOutAction} />}
        />
      }
      navbar={
        <Navbar
          canViewProblems
          canViewAnnounces
          collapsed={collapsed}
          onOpenToggleClick={toggle}
        />
      }
      navbarCollapsed={collapsed}
    >
      <h1>Main</h1>
      <p>
        あのイーハトーヴォのすきとおった風、夏でも底に冷たさをもつ青いそら、うつくしい森で飾られたモリーオ市、郊外のぎらぎらひかる草の波。
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
