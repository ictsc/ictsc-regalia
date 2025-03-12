import {
  type ReactNode,
  startTransition,
  use,
  useDeferredValue,
  useReducer,
} from "react";
import { useSignOut, type User } from "../features/viewer";
import { useSchedule } from "../features/schedule";
import { Layout, Header, Navbar, AccountMenu } from "../components/app-shell";
import { Phase } from "@ictsc/proto/contestant/v1";

export function AppShell({
  children,
  viewer: viewerPromise,
}: {
  readonly children?: ReactNode;
  readonly viewer: Promise<User>;
}) {
  const viewer = use(useDeferredValue(viewerPromise));
  const [schedule] = useSchedule();
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
          <Navbar
            canViewProblems={schedule?.phase === Phase.IN_CONTEST}
            canViewAnnounces={schedule?.phase === Phase.IN_CONTEST}
            collapsed={collapsed}
            onOpenToggleClick={toggle}
          />
        ) : null
      }
      navbarCollapsed={collapsed}
    >
      {children}
    </Layout>
  );
}
