import { fetchNotices } from "@app/features/announce";
import { createFileRoute } from "@tanstack/react-router";
import { AnnounceList } from "./page";
import { use, useDeferredValue } from "react";

export const Route = createFileRoute("/announces/")({
  component: RouteComponent,
  loader: ({ context: { transport } }) => {
    return {
      announces: fetchNotices(transport),
    };
  },
});

function RouteComponent() {
  const { announces: AnnouncePromise } = Route.useLoaderData();
  const deferredAnnouncePromise = useDeferredValue(AnnouncePromise);
  const announces = use(deferredAnnouncePromise);
  return <AnnounceList announces={announces} />;
}
