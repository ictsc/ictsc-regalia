import { use, useDeferredValue, useMemo } from "react";
import { createLazyFileRoute } from "@tanstack/react-router";
import { Markdown, Typography } from "../components/markdown";
import { Title } from "..//components/title";
import { Route as ParentRoute } from "./~announces";

export const Route = createLazyFileRoute("/announces/$slug")({
  component: RouteComponent,
});

function RouteComponent() {
  const { slug } = Route.useLoaderData();
  const { announces } = ParentRoute.useLoaderData();
  const announcePromise = useMemo(async () => {
    const list = await announces;
    return list.find((announce) => announce.slug === slug);
  }, [announces, slug]);

  const announce = use(useDeferredValue(announcePromise));

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
