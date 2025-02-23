import {
  type ReactNode,
  startTransition,
  use,
  useDeferredValue,
  useReducer,
} from "react";
import { useSignOut, type User } from "@app/features/viewer";
import { Layout, Header, Navbar, AccountMenu } from "@app/components/app-shell";

export function AppShell({
  children,
  viewer: viewerPromise,
}: {
  readonly children?: ReactNode;
  readonly viewer: Promise<User>;
}) {
  const viewer = use(useDeferredValue(viewerPromise));
  const [collapsed, toggle] = useReducer((o) => !o, false);
  const signOutAction = useSignOut();

  return (
    <Layout
      header={
        <Header
          accountMenu={
            viewer?.type === "contestant" && (
              <AccountMenu
                name={viewer.name}
                onSignOut={() => {
                  startTransition(() => signOutAction());
                }}
              />
            )
          }
        />
      }
      navbar={
        viewer?.type === "contestant" ? (
          <Navbar collapsed={collapsed} onOpenToggleClick={toggle} />
        ) : null
      }
      navbarCollapsed={collapsed}
    >
      {children}
    </Layout>
  );
}
