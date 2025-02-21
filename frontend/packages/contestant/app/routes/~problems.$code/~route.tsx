import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/problems/$code")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello /problems/$code!</div>;
}
