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
import { SignInPage } from "./signin.page";

type SignInSearch = {
  next?: string;
};

export const Route = createFileRoute("/signin")({
  component: RouteComponent,
  validateSearch: (search: Record<string, unknown>): SignInSearch => {
    return {
      next: typeof search.next === "string" ? search.next : undefined,
    };
  },
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
  const { next: nextPath = "/" } = Route.useSearch();

  // 既にログインしているならそのまま next に遷移する
  // ログインしているかどうかはログインの操作が可能かどうかには関係ない
  // そのためログインしているかの確認の前にページ自体は表示しておき，状態が確定したら遷移するか判断する
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

  return <SignInPage signInURL={authURL.toString()} />;
}
