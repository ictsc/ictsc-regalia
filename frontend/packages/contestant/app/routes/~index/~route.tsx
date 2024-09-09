import { createFileRoute } from "@tanstack/react-router";
import { IndexPage } from "./page";

export const Route = createFileRoute("/")({
  component: IndexPage,
});
