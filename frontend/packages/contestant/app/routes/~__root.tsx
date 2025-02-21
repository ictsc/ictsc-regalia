import { Suspense, lazy, startTransition, use, useEffect } from "react";
import {
  createRootRouteWithContext,
  Outlet,
  useNavigate,
  useRouterState,
} from "@tanstack/react-router";
import type { Transport } from "@connectrpc/connect";
import { AppShell } from "@app/components/app-shell";
import { fetchViewer, type User } from "@app/features/viewer";

interface RouterContext {
  transport: Transport;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: Root,
  loader: ({ context: { transport } }) => ({
    viewer: fetchViewer(transport),
  }),
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
        <Redirector viewer={viewer} />
      </Suspense>
      <Suspense>
        <TanStackRouterDevtools />
      </Suspense>
    </>
  );
}

function Redirector({ viewer: viewerPromise }: { viewer: Promise<User> }) {
  const routerState = useRouterState();
  const navigate = useNavigate();
  const viewer = use(viewerPromise);

  useEffect(() => {
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
