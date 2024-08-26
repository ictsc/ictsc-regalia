import { Suspense, lazy } from "react";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { AppShell } from "../components/app-shell";

const TanStackRouterDevtools = import.meta.env.DEV
  ? lazy(() =>
      import("@tanstack/router-devtools").then((mod) => ({
        default: mod.TanStackRouterDevtools,
      })),
    )
  : () => null;

export const Route = createRootRoute({
  component: () => (
    <>
      <AppShell>
        <Outlet />
      </AppShell>
      <Suspense>
        <TanStackRouterDevtools />
      </Suspense>
    </>
  ),
});
