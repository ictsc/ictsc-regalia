import { createFileRoute, Outlet, useNavigate } from "@tanstack/react-router";
import { useSchedule } from "../features/schedule";
import { startTransition, useEffect } from "react";
import { Phase } from "@ictsc/proto/contestant/v1";

export const Route = createFileRoute("/problems")({
  component: RouteComponent,
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
