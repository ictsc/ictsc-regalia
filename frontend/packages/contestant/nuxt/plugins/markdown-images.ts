import { setMarkdownImageUrl } from "@ictsc/ui/markdown";

// Problem and notice Markdown refers to repository images as /images/<file>.
export default defineNuxtPlugin(() => {
  setMarkdownImageUrl(
    (name) => `/api/v1/contestant/images/${encodeURIComponent(name)}`,
  );
});
