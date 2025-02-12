import { Suspense, lazy } from "react";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import type { Transport } from "@connectrpc/connect";
import { AppShell } from "@app/components/app-shell";
import { fetchViewer } from "@app/features/viewer";

interface RouterContext {
  transport: Transport;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  loader: ({ context: { transport } }) => ({
    viewer: fetchViewer(transport),
  }),
  component: Root,
});

const TanStackRouterDevtools = import.meta.env.DEV
  ? lazy(() =>
      import("@tanstack/router-devtools").then((mod) => ({
        default: mod.TanStackRouterDevtools,
      })),
    )
  : () => null;

function Root() {
  const { viewer } = Route.useLoaderData();
  return (
    <>
      <AppShell viewer={viewer}>
        <Outlet />
      </AppShell>
      <Suspense>
        <TanStackRouterDevtools />
      </Suspense>
    </>
  );
}
