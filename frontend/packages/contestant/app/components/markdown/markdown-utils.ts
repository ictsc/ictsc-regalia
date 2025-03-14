import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import remarkRehype from "remark-rehype";
import rehypeKatex from "rehype-katex";
import rehypeStringify from "rehype-stringify";
import rehypeShikiFromHighlighter from "@shikijs/rehype/core";
import { createHighlighterCore } from "shiki/core";
import { createOnigurumaEngine } from "shiki/engine/oniguruma";

const highlighterPromise = createHighlighterCore({
  themes: [import("@shikijs/themes/light-plus")],
  langs: [
    import("@shikijs/langs/diff"),
    import("@shikijs/langs/shellscript"),
    import("@shikijs/langs/shellsession"),
    import("@shikijs/langs/hcl"),
    import("@shikijs/langs/sql"),
  ],
  engine: createOnigurumaEngine(import("shiki/wasm")),
});

export async function markdownToHtml(content: string): Promise<string> {
  const highlighter = await highlighterPromise;

  const file = await unified()
    .use(remarkParse)
    .use(remarkGfm)
    .use(remarkMath)
    .use(remarkRehype, { allowDangerousHtml: true })
    .use(rehypeKatex)
    .use(rehypeShikiFromHighlighter, highlighter as never, {
      themes: {
        light: "light-plus",
        dark: "light-plus",
      },
    })
    .use(rehypeStringify, { allowDangerousHtml: true })
    .process(content);
  return file.toString();
}
