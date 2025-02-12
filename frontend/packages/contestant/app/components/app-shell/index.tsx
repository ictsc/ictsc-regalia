import {
  type ReactNode,
  startTransition,
  use,
  useDeferredValue,
  useEffect,
  useReducer,
} from "react";
import { type User } from "@app/features/viewer";
import { useNavigate, useRouterState } from "@tanstack/react-router";
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
    <>
      <Layout
        header={<Header user={viewer} />}
        navbar={
          defferedViewer?.type === "contestant" && (
            <NavbarView collapsed={collapsed} onOpenToggleClick={toggle} />
          )
        }
        navbarCollapsed={collapsed}
      >
        {children}
      </Layout>
      <Redirector viewer={defferedViewer} />
    </>
  );
}

function Redirector({ viewer }: { viewer: User | undefined }) {
  const routerState = useRouterState();
  const navigate = useNavigate();

  useEffect(() => {
    if (viewer == null) return;
    switch (viewer.type) {
      case "unauthenticated": {
        if (routerState.location.pathname !== "/signin") {
          startTransition(() =>
            navigate({
              to: "/signin",
              search: { next: routerState.location.pathname },
            }),
          );
        }
        break;
      }
      case "pre-signup": {
        if (routerState.location.pathname !== "/signup") {
          startTransition(() =>
            navigate({
              to: "/signup",
              search: { next: routerState.location.pathname },
            }),
          );
        }
        break;
      }
    }
  }, [routerState, navigate, viewer]);

  return null;
}
