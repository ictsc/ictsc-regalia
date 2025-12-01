import { use, useDeferredValue } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { AnnounceList } from "./announces.index.page";
import { Route as ParentRoute } from "./~announces";

export const Route = createFileRoute("/announces/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { announces: AnnouncePromise } = ParentRoute.useLoaderData();
  const deferredAnnouncePromise = useDeferredValue(AnnouncePromise);
  const announces = use(deferredAnnouncePromise);
  return <AnnounceList announces={announces} />;
}
