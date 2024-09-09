import { useReducer } from "react";
import type { Meta, StoryObj } from "@storybook/react";
import { Layout } from "./layout";
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
