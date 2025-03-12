import { startTransition, useEffect } from "react";
import { createFileRoute, Outlet, useNavigate } from "@tanstack/react-router";
import { Phase } from "@ictsc/proto/contestant/v1";
import { useSchedule } from "../features/schedule";
import { fetchNotices } from "../features/announce";

export const Route = createFileRoute("/announces")({
  component: RouteComponent,
  loader: ({ context: { transport } }) => ({
    announces: fetchNotices(transport),
  }),
});

function RouteComponent() {
  const [schedule, isPending] = useSchedule();
  const navigate = useNavigate();
  useEffect(() => {
    if (schedule == null || isPending) {
      return;
    }
    if (schedule.phase !== Phase.IN_CONTEST) {
      startTransition(() => navigate({ to: "/" }));
    }
  }, [schedule, isPending, navigate]);
  return <Outlet />;
}
