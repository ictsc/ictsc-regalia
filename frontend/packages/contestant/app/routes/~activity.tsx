import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { fetchActivity } from "@app/features/activity";
import { startTransition, use, useDeferredValue, useEffect } from "react";
import { useSchedule, hasContestStarted } from "../features/schedule";
import { protoScoreToProps } from "../features/score";
import { ActivityPage } from "./activity.page";

export const Route = createFileRoute("/activity")({
  component: RouteComponent,
  loader: ({ context: { transport } }) => {
    return {
      activity: fetchActivity(transport),
    };
  },
});

function RouteComponent() {
  const [schedule, isPending] = useSchedule();
  const navigate = useNavigate();
  useEffect(() => {
    if (schedule == null || isPending) {
      return;
    }
    if (!hasContestStarted(schedule)) {
      startTransition(() => navigate({ to: "/" }));
    }
  }, [schedule, isPending, navigate]);

  const { activity: activityPromise } = Route.useLoaderData();
  const deferredActivityPromise = useDeferredValue(activityPromise);
  const activity = use(deferredActivityPromise);
  return (
    <ActivityPage
      entries={activity.map((entry) => ({
        problemCode: entry.problemCode,
        problemTitle: entry.problemTitle,
        answerId: entry.answerId,
        submittedAt: entry.submittedAt,
        score: protoScoreToProps(entry.maxScore, entry.score),
        scored: entry.score?.score != null,
      }))}
    />
  );
}
