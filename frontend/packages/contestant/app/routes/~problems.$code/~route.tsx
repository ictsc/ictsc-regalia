import { createFileRoute } from "@tanstack/react-router";
import { fetchProblem } from "@app/features/problem";

export const Route = createFileRoute("/problems/$code")({
  loader: ({ context: { transport }, params: { code } }) => ({
    problem: fetchProblem(transport, code),
  }),
});
