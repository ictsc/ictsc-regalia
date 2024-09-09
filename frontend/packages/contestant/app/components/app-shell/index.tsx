import { type ReactNode, useReducer } from "react";
import { Layout } from "./layout";
import { HeaderView } from "./header";
import { NavbarView } from "./navbar";

// 仮置きのAPIなしで動くやつ
export function AppShell({ children }: { readonly children: ReactNode }) {
  const [collapsed, toggle] = useReducer((o) => !o, false);
  return (
    <Layout
      header={<HeaderView />}
      navbar={<NavbarView collapsed={collapsed} onOpenToggleClick={toggle} />}
      navbarCollapsed={collapsed}
    >
      {children}
    </Layout>
  );
}
