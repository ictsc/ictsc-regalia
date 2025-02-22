import { Markdown, Typography } from "@app/components/markdown";
import { createLazyFileRoute } from "@tanstack/react-router";
import { use } from "react";

export const Route = createLazyFileRoute("/rule")({
  component: RouteComponent,
});

function RouteComponent() {
  const { rule: rulePromise } = Route.useLoaderData();
  const rule = use(rulePromise);
  return (
    <div className="container mx-auto mb-64 mt-20">
      <Typography>
        <Markdown>{rule.markdown}</Markdown>
      </Typography>
    </div>
  );
}
