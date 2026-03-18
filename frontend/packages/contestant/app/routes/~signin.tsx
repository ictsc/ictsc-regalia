import { createFileRoute, Outlet } from "@tanstack/react-router";

export type SignInSearch = {
  next?: string;
};

export function validateSignInSearch(
  search: Record<string, unknown>,
): SignInSearch {
  return {
    next: typeof search.next === "string" ? search.next : undefined,
  };
}

export const Route = createFileRoute("/signin")({
  component: () => <Outlet />,
  validateSearch: validateSignInSearch,
});
