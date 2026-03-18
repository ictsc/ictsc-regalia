import { createFileRoute, Outlet } from "@tanstack/react-router";

type SignInSearch = {
  next?: string;
};

export const Route = createFileRoute("/signin")({
  component: () => <Outlet />,
  validateSearch: (search: Record<string, unknown>): SignInSearch => {
    return {
      next: typeof search.next === "string" ? search.next : undefined,
    };
  },
});
