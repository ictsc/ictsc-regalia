import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import remarkRehype from "remark-rehype";
import rehypeKatex from "rehype-katex";
import rehypeShiki from "@shikijs/rehype/core";
import rehypeReact from "rehype-react";
import production from "react/jsx-runtime";
import { createHighlighterCore } from "shiki/core";
import { createOnigurumaEngine } from "shiki/engine/oniguruma";
import { ReactNode } from "react";

const highlighterPromise = createHighlighterCore({
  themes: [import("@shikijs/themes/material-theme-lighter")],
  langs: [
    import("@shikijs/langs/diff"),
    import("@shikijs/langs/shellscript"),
    import("@shikijs/langs/shellsession"),
    import("@shikijs/langs/hcl"),
    import("@shikijs/langs/sql"),
  ],
  engine: createOnigurumaEngine(import("shiki/wasm")),
});

export async function renderMarkdown(content: string): Promise<ReactNode> {
  const highlighter = await highlighterPromise;

  /* eslint-disable */
  const file = await unified()
    .use(remarkParse, { fragment: true })
    .use(remarkGfm)
    .use(remarkMath)
    .use(remarkRehype)
    .use(rehypeKatex)
    .use(rehypeShiki, highlighter as any, {
      theme: "material-theme-lighter",
    })
    .use(rehypeReact, production)
    .process(content);
  return file.result;
  /* eslint-enable */
}
