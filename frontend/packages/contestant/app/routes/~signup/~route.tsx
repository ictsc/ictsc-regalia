import { startTransition, use, useState } from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Route as RootRoute } from "../~__root";
import { SignUpPage } from "./page";
import { signUp, SignUpResponse } from "@app/features/viewer/signup";

export const Route = createFileRoute("/signup")({
  component: RouteComponent,
});

function RouteComponent() {
  const { viewer: viewerPromise } = RootRoute.useLoaderData();
  const viewer = use(viewerPromise);
  const navigate = useNavigate();
  const [errState, setErrState] = useState({} as SignUpResponse);

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
          if (resp.error == null) {
            await navigate({ to: "/" });
          }
          setErrState(resp);
        });
      }}
      {...errState}
    />
  );
}
