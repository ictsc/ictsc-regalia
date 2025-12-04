import { use } from "react";
import { createLazyFileRoute } from "@tanstack/react-router";
import { Title } from "../components/title";
import { Markdown, Typography } from "../components/markdown";

export const Route = createLazyFileRoute("/rule")({
  component: RouteComponent,
});

function RouteComponent() {
  const { rule: rulePromise } = Route.useLoaderData();
  const rule = use(rulePromise);
  return (
    <>
      <Title>ルール</Title>
      <div className="mx-40 mb-64 mt-20 max-w-screen-lg">
        <Typography>
          <Markdown>{rule.markdown}</Markdown>
        </Typography>
      </div>
    </>
  );
}
