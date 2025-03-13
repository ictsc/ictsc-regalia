import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/announces/$slug")({
  loader({ params: { slug } }) {
    return { slug };
  },
});
