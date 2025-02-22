import { createFileRoute } from "@tanstack/react-router";
import { fetchProblem } from "@app/features/problem";
import { lazy, useDeferredValue } from "react";

export const Route = createFileRoute("/problems/$code")({
  component: RouteComponent,
  loader: ({ context: { transport }, params: { code } }) => ({
    problem: fetchProblem(transport, code),
  }),
});

const ProblemPage = lazy(async () => {
  const mod = await import("./page");
  return { default: mod.ProblemPage };
});

function RouteComponent() {
  const { problem } = Route.useLoaderData();
  const deferredProblem = useDeferredValue(problem);
  return <ProblemPage problem={deferredProblem} />;
}
