import { Suspense, lazy } from "react";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import type { Transport } from "@connectrpc/connect";
import { AppShell } from "@app/components/app-shell";
import { fetchMe } from "@app/features/account";

interface RouterContext {
  transport: Transport;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  loader: ({ context: { transport } }) => ({
    me: fetchMe(transport),
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
  const { me } = Route.useLoaderData();
  return (
    <>
      <AppShell me={me}>
        <Outlet />
      </AppShell>
      <Suspense>
        <TanStackRouterDevtools />
      </Suspense>
    </>
  );
}
