import { createFileRoute } from "@tanstack/react-router";
import { fetchProblems } from "@app/features/problem";
import { use, useDeferredValue } from "react";
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
  const { problems: problemsPromise } = Route.useLoaderData();
  const deferredProblemsPromise = useDeferredValue(problemsPromise);
  const problems = use(deferredProblemsPromise);
  return <ProblemsPage problems={problems} />;
}
