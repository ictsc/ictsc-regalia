import { createFileRoute, useRouter } from "@tanstack/react-router";
import { fetchProblems } from "@app/features/problem";
import { startTransition, use, useDeferredValue, useEffect } from "react";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { ProblemsPage } from "./problems.index/page";

export const Route = createFileRoute("/problems/")({
  component: RouteComponent,
  loader: ({ context: { transport } }) => {
    return {
      problems: fetchProblems(transport),
    };
  },
});

function RouteComponent() {
  const router = useRouter();
  const { problems: problemsPromise } = Route.useLoaderData();
  const deferredProblemsPromise = useDeferredValue(problemsPromise);
  const problems = use(deferredProblemsPromise);

  // 最も近い submittableUntil/submittableFrom に達したらリフェッチ
  useEffect(() => {
    let earliest: number | null = null;
    for (const p of problems) {
      const until = p.submissionStatus?.submittableUntil;
      const from = p.submissionStatus?.submittableFrom;
      const target = until ?? from;
      if (target == null) continue;
      const ms = timestampDate(target).getTime() - Date.now();
      if (ms <= 0) continue;
      if (earliest == null || ms < earliest) earliest = ms;
    }
    if (earliest == null) return;
    const timer = setTimeout(() => {
      startTransition(() => router.invalidate());
    }, earliest);
    return () => clearTimeout(timer);
  }, [problems, router]);

  return <ProblemsPage problems={problems} />;
}
