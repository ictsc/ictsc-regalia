import {
  startTransition,
  Suspense,
  use,
  useDeferredValue,
  useEffect,
} from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { type User } from "@app/features/viewer";
import { Route as RootRoute } from "./~__root";
import { Route as SigninRoute } from "./~signin";
import { SignInPage } from "./signin.page";

export const Route = createFileRoute("/signin/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { viewer } = RootRoute.useLoaderData();

  return (
    <Suspense>
      <Page viewerPromise={viewer} />
    </Suspense>
  );
}

function Page({ viewerPromise }: { viewerPromise: Promise<User> }) {
  const { next: nextPath = "/" } = SigninRoute.useSearch();

  const navigate = useNavigate();
  const deferredViewer = useDeferredValue(viewerPromise, null);
  const viewer = deferredViewer != null ? use(deferredViewer) : null;

  useEffect(() => {
    if (viewer == null) return;
    if (viewer.type !== "unauthenticated") {
      startTransition(() => navigate({ to: nextPath }));
    }
  }, [viewer, navigate, nextPath]);

  const authURL = new URL("/api/auth/signin", window.location.origin);
  authURL.searchParams.set("next", nextPath);
  const impersonationURL = new URL(
    "/signin/impersonation",
    window.location.origin,
  );
  impersonationURL.searchParams.set("next", nextPath);

  return (
    <SignInPage
      signInURL={authURL.toString()}
      adminTokenAvailable={viewer?.admin.canImpersonateContestants ?? false}
      impersonationURL={impersonationURL.toString()}
    />
  );
}
