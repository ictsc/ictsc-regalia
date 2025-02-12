import { createFileRoute } from "@tanstack/react-router";

type SignInSearch = {
  next?: string;
};

export const Route = createFileRoute("/signin")({
  component: RouteComponent,
  validateSearch: (search: Record<string, unknown>): SignInSearch => {
    return {
      next: typeof search.next === "string" ? search.next : undefined,
    };
  },
});

function RouteComponent() {
  const { next: nextPath = "/" } = Route.useSearch();
  const authURL = new URL("/api/auth/signin", window.location.origin);
  authURL.searchParams.set("next", nextPath);
  return (
    <div>
      <a href={authURL.toString()}>Sign in</a>
    </div>
  );
}
