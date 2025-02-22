import { fetchRule } from "@app/features/rule";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/rule")({
  loader: ({ context: { transport } }) => ({
    rule: fetchRule(transport),
  }),
});
