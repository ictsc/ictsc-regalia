import { type ReactNode, use, useDeferredValue, useReducer } from "react";
import { type User } from "@app/features/viewer";
import { Layout } from "./layout";
import { Header } from "./header";
import { NavbarView } from "./navbar";

type AppShellProps = {
  readonly children: ReactNode;
  readonly viewer: Promise<User>;
};

export function AppShell({ children, viewer: viewerPromise }: AppShellProps) {
  const viewer = use(useDeferredValue(viewerPromise));
  const [collapsed, toggle] = useReducer((o) => !o, false);

  return (
    <Layout
      header={<Header user={viewerPromise} />}
      navbar={
        viewer?.type === "contestant" ? (
          <NavbarView collapsed={collapsed} onOpenToggleClick={toggle} />
        ) : null
      }
      navbarCollapsed={collapsed}
    >
      {children}
    </Layout>
  );
}
