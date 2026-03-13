import { useEffect } from "react";
import { createLazyFileRoute } from "@tanstack/react-router";
import { Markdown, Typography } from "../components/markdown";
import { Title } from "../components/title";
import { useReadAnnouncements } from "../features/announce-read-status";
import { ReadToggleButton } from "./problems.index/unread-announces-banner";

export const Route = createLazyFileRoute("/announces/$slug")({
  component: RouteComponent,
});

function RouteComponent() {
  const { announce } = Route.useLoaderData();
  const { slug } = Route.useParams();
  const { markAsRead } = useReadAnnouncements();

  useEffect(() => {
    markAsRead(slug);
  }, [slug, markAsRead]);

  return (
    <>
      <Title>{announce?.title}</Title>
      <div className="mt-64 flex w-full flex-col px-40">
        {announce == null ? (
          <h1 className="mx-auto font-bold">アナウンスがありません</h1>
        ) : (
          <Typography className="w-full">
            <h1 className="flex items-center justify-between gap-8">
              {announce.title}
              <ReadToggleButton slug={slug} />
            </h1>
            <Markdown>{announce.body}</Markdown>
          </Typography>
        )}
      </div>
    </>
  );
}
