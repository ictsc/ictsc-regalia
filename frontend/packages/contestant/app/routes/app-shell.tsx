import {
  type ReactNode,
  startTransition,
  use,
  useDeferredValue,
  useReducer,
} from "react";
import { useSignOut, type User } from "../features/viewer";
import { useSchedule, isInContest } from "../features/schedule";
import { Layout, Header, Navbar, AccountMenu } from "../components/app-shell";

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
  const inContest = isInContest(schedule);

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
            canViewProblems={inContest}
            canViewAnnounces={inContest}
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
