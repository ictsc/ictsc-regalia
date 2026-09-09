import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import remarkBreaks from "remark-breaks";
import remarkRehype from "remark-rehype";
import rehypeSanitize, { defaultSchema } from "rehype-sanitize";
import rehypeKatex from "rehype-katex";
import rehypeShiki from "@shikijs/rehype";
import rehypeStringify from "rehype-stringify";

const processor = unified()
  .use(remarkParse)
  .use(remarkGfm)
  .use(remarkMath)
  .use(remarkBreaks)
  .use(remarkRehype)
  .use(rehypeSanitize, {
    ...defaultSchema,
    attributes: {
      ...defaultSchema.attributes,
      code: [["className", /^language-./, "math-inline", "math-display"]],
    },
  })
  .use(rehypeKatex, { trust: false })
  .use(rehypeShiki, {
    theme: "github-light",
    langs: [
      "javascript",
      "typescript",
      "bash",
      "json",
      "yaml",
      "go",
      "python",
      "sql",
      "diff",
      "html",
      "css",
    ],
    fallbackLanguage: "text",
  })
  .use(rehypeStringify);
export async function renderMarkdown(markdown: string): Promise<string> {
  return String(await processor.process(markdown));
}
