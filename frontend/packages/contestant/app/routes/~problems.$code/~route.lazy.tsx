import { createLazyFileRoute } from "@tanstack/react-router";
import { Layout, Content, Sidebar } from "./page";
import { use } from "react";

export const Route = createLazyFileRoute("/problems/$code")({
  component: RouteComponent,
});

function RouteComponent() {
  const { problem: problemPromise } = Route.useLoaderData();
  const problem = use(problemPromise);
  return <Layout content={<Content {...problem} />} sidebar={<Sidebar />} />;
}
