import {
  startTransition,
  use,
  useDeferredValue,
  useEffect,
  useState,
} from "react";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { Route as RootRoute } from "../~__root";
import { SignUpPage } from "./page";
import { signUp, type SignUpResponse } from "@app/features/viewer/signup";

export const Route = createFileRoute("/signup")({
  component: RouteComponent,
});

function RouteComponent() {
  const { viewer: viewerPromise } = RootRoute.useLoaderData();
  const viewer = use(useDeferredValue(viewerPromise));
  const router = useRouter();
  const [errState, setErrState] = useState({} as SignUpResponse);

  useEffect(() => {
    if (viewer.type !== "pre-signup") {
      startTransition(() => router.navigate({ to: "/" }));
    }
  }, [viewer, router]);

  if (viewer.type !== "pre-signup") {
    return null;
  }
  return (
    <SignUpPage
      defaultName={viewer.name}
      defaultDisplayName={viewer.displayName}
      submit={(data) => {
        setErrState({});
        startTransition(async () => {
          const resp = await signUp(data);
          if (resp.error != null) {
            setErrState(resp);
          }
          await router.invalidate();
        });
      }}
      {...errState}
    />
  );
}
