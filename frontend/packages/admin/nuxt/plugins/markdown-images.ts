import { setMarkdownImageUrl } from "@ictsc/ui/markdown";
import { adminImageUrl } from "~/features/admin-images";

// Problem and notice Markdown refers to repository images as /images/<file>.
export default defineNuxtPlugin(() => {
  setMarkdownImageUrl(adminImageUrl());
});
