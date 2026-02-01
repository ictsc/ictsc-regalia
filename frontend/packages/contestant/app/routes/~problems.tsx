import { createFileRoute, Outlet, useNavigate } from "@tanstack/react-router";
import { useSchedule, isInContest } from "../features/schedule";
import { startTransition, useEffect } from "react";

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
    if (!isInContest(schedule)) {
      startTransition(() => navigate({ to: "/" }));
    }
  }, [schedule, isPending, navigate]);
  return <Outlet />;
}
