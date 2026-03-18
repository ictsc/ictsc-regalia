import {
  type ReactNode,
  startTransition,
  use,
  useDeferredValue,
  useReducer,
} from "react";
import { useSignOut, type User } from "../features/viewer";
import { useSchedule, hasContestStarted } from "../features/schedule";
import { Layout, Header, Navbar, AccountMenu } from "../components/app-shell";
import { MaterialSymbol } from "../components/material-symbol";

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
  const inContest = hasContestStarted(schedule);

  return (
    <>
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
              canViewActivity={inContest}
              collapsed={collapsed}
              onOpenToggleClick={toggle}
            />
          ) : null
        }
        navbarCollapsed={collapsed}
      >
        {children}
      </Layout>
      {viewer.type === "contestant" && viewer.impersonation != null ? (
        <div className="pointer-events-none fixed inset-x-0 top-0 z-20 flex h-[36px] items-center gap-12 bg-[#191970]/65 px-16">
          <MaterialSymbol
            icon="admin_panel_settings"
            size={20}
            className="text-surface-0 shrink-0 opacity-90"
          />
          <p className="text-14 text-surface-0 flex-1">
            {viewer.impersonation.adminName}として
            <span className="text-surface-0 font-bold">
              {viewer.displayName || viewer.name}
            </span>
            を操作中
          </p>
          <button
            type="button"
            className="text-surface-0 rounded-4 text-14 pointer-events-auto border border-[#ffffff]/40 px-12 py-4 hover:bg-[#ffffff]/15"
            onClick={() => {
              startTransition(() => signOutAction());
            }}
          >
            操作終了
          </button>
        </div>
      ) : null}
    </>
  );
}
