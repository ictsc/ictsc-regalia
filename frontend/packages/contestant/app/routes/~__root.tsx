import { Suspense, lazy, startTransition, use, useEffect } from "react";
import { ErrorBoundary, type FallbackProps } from "react-error-boundary";
import {
  createRootRouteWithContext,
  Outlet,
  useNavigate,
  useRouterState,
} from "@tanstack/react-router";
import { ConnectError, Code, type Transport } from "@connectrpc/connect";
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
        <ErrorBoundary FallbackComponent={ErrorFallback}>
          <Suspense fallback={null}>
            <Outlet />
          </Suspense>
        </ErrorBoundary>
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

function ErrorFallback(props: FallbackProps) {
  const error: unknown = props.error;
  if (error instanceof ConnectError && error.code === Code.Unauthenticated) {
    return <UnauthorizedFallback {...props} />;
  }
  // TODO: エラー画面を表示する
  return null;
}

function UnauthorizedFallback({ resetErrorBoundary }: FallbackProps) {
  // viewer の取得よりも先に認証エラーが発生した場合 ErrorBoundary が反応するが，
  // 認証エラーは Redirector により処理されるため，適当に再開しつつエラーを無視する
  useEffect(() => {
    resetErrorBoundary();
  }, [resetErrorBoundary]);
  return null;
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
