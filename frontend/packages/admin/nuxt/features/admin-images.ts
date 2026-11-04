import type { ImageUrl } from "@ictsc/ui/markdown";

/**
 * Serves /images/<file> from the given content commit, or from the active
 * content when no commit is given.
 */
export function adminImageUrl(commit?: string): ImageUrl {
  const query = commit ? `?commit=${encodeURIComponent(commit)}` : "";
  return (name) => `/api/v1/admin/images/${encodeURIComponent(name)}${query}`;
}
