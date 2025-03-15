import { use, useMemo } from "react";
import { renderMarkdown } from "./markdown-utils";

export function Markdown({ children }: { children?: string }) {
  const nodePromise = useMemo(() => renderMarkdown(children ?? ""), [children]);
  return use(nodePromise);
}
