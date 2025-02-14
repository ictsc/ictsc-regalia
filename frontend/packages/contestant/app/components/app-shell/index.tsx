import { type ReactNode, use, useDeferredValue, useReducer } from "react";
import { type User } from "@app/features/viewer";
import { Layout } from "./layout";
import { Header } from "./header";
import { NavbarView } from "./navbar";

type AppShellProps = {
  readonly children: ReactNode;
  readonly viewer: Promise<User>;
};

const initialViewer = Promise.resolve<User | undefined>(undefined);

export function AppShell({ children, viewer }: AppShellProps) {
  const defferedViewer = use(useDeferredValue(viewer, initialViewer));
  const [collapsed, toggle] = useReducer((o) => !o, false);

  return (
    <Layout
      header={<Header user={viewer} />}
      navbar={
        defferedViewer?.type === "contestant" ? (
          <NavbarView collapsed={collapsed} onOpenToggleClick={toggle} />
        ) : null
      }
      navbarCollapsed={collapsed}
    >
      {children}
    </Layout>
  );
}
