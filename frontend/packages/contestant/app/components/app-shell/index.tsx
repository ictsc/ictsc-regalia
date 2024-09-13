import { type ReactNode, useReducer } from "react";
import { Layout } from "./layout";
import { Header } from "./header";
import { NavbarView } from "./navbar";
import { type User } from "@app/features/account";

type AppShellProps = {
  readonly children: ReactNode;
  readonly me: Promise<User | undefined>;
};

export function AppShell({ children, me }: AppShellProps) {
  const [collapsed, toggle] = useReducer((o) => !o, false);
  return (
    <Layout
      header={<Header user={me} />}
      navbar={<NavbarView collapsed={collapsed} onOpenToggleClick={toggle} />}
      navbarCollapsed={collapsed}
    >
      {children}
    </Layout>
  );
}
