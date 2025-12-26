import { createLazyFileRoute } from "@tanstack/react-router";
import { Markdown, Typography } from "../components/markdown";
import { Title } from "..//components/title";

export const Route = createLazyFileRoute("/announces/$slug")({
  component: RouteComponent,
});

function RouteComponent() {
  const { announce } = Route.useLoaderData();

  return (
    <>
      <Title>{announce?.title}</Title>
      <div className="mt-64 flex w-full px-40">
        {announce == null ? (
          <h1 className="mx-auto font-bold">アナウンスがありません</h1>
        ) : (
          <Typography className="w-full">
            <h1>{announce.title}</h1>
            <Markdown>{announce.body}</Markdown>
          </Typography>
        )}
      </div>
    </>
  );
}
