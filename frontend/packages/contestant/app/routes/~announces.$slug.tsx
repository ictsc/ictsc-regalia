import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/announces/$slug")({
  async loader({ params: { slug }, parentMatchPromise }) {
    const match = await parentMatchPromise;
    const announces = await match.loaderData?.announces;
    const announce = announces?.find((a) => a.slug === slug);
    return { announce };
  },
});
