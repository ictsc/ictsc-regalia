import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { fetchActivity } from "@app/features/activity";
import { fetchAnswer } from "@app/features/answer";
import { startTransition, use, useDeferredValue, useEffect } from "react";
import { useSchedule, hasContestStarted } from "../features/schedule";
import { protoScoreToProps } from "../features/score";
import { ActivityPage } from "./activity.page";

export const Route = createFileRoute("/activity")({
  component: RouteComponent,
  loader: ({ context: { transport } }) => {
    return {
      activity: fetchActivity(transport),
      fetchAnswer: (problemCode: string, num: number) =>
        fetchAnswer(transport, problemCode, num),
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

  const { activity: activityPromise, fetchAnswer: fetchAnswerFn } =
    Route.useLoaderData();
  const deferredActivityPromise = useDeferredValue(activityPromise);
  const activity = use(deferredActivityPromise);

  const formatDate = (date: Date) => {
    const pad = (num: number) => num.toString().padStart(2, "0");

    const year = date.getFullYear();
    const month = pad(date.getMonth() + 1);
    const day = pad(date.getDate());
    const hour = pad(date.getHours());
    const minute = pad(date.getMinutes());
    const second = pad(date.getSeconds());

    return `${year}-${month}-${day}-${hour}-${minute}-${second}`;
  };

  const downloadFile = (filename: string, content: string) => {
    const blob = new Blob([content], { type: "text/markdown" });
    const url = URL.createObjectURL(blob);

    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
  };

  const downloadAnswer = async (problemCode: string, num: number) => {
    const { answerBody, submittedAtString } = await fetchAnswerFn(
      problemCode,
      num,
    );
    const submittedAt = new Date(submittedAtString);
    const submittedAtFormattedString = formatDate(submittedAt);
    const filename = `${problemCode}-${submittedAtFormattedString}.md`;

    downloadFile(filename, answerBody);
  };

  return (
    <ActivityPage
      entries={activity.map((entry) => ({
        problemCode: entry.problemCode,
        problemTitle: entry.problemTitle,
        answerId: entry.answerId,
        submittedAt: entry.submittedAt,
        score: protoScoreToProps(entry.maxScore, entry.score),
        scored: entry.score != null,
        onDownload: () =>
          startTransition(() =>
            downloadAnswer(entry.problemCode, entry.answerId),
          ),
      }))}
    />
  );
}
