import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import remarkRehype from "remark-rehype";
import rehypeSanitize, { defaultSchema } from "rehype-sanitize";
import rehypeKatex from "rehype-katex";
import rehypeShiki from "@shikijs/rehype";
import rehypeStringify from "rehype-stringify";

type HastNode = {
  type: string;
  tagName?: string;
  properties?: Record<string, unknown>;
  children?: HastNode[];
  value?: string;
};

/** Maps a content image file name to the URL that serves it. */
export type ImageUrl = (name: string) => string;

// Content repositories reference their images as /images/<file>; each app
// serves them from its own authenticated API path.
let defaultImageUrl: ImageUrl | undefined;

/** Sets how "/images/<file>" sources are rewritten when no per-render URL is given. */
export function setMarkdownImageUrl(imageUrl: ImageUrl | undefined) {
  defaultImageUrl = imageUrl;
}

export function rewriteImageSource(
  src: string,
  imageUrl = defaultImageUrl,
): string {
  return imageUrl && src.startsWith("/images/")
    ? imageUrl(decodeURIComponent(src.slice("/images/".length)))
    : src;
}

function rehypeContentImages() {
  return (tree: HastNode, file: { data: Record<string, unknown> }) => {
    const imageUrl =
      (file.data.imageUrl as ImageUrl | undefined) ?? defaultImageUrl;
    const blockExternal = file.data.blockExternalImages === true;
    const visit = (node: HastNode) => {
      if (node.tagName === "img" && typeof node.properties?.src === "string") {
        const src = node.properties.src;
        if (blockExternal && !src.startsWith("/images/")) {
          // Loading it would tell the image's host who read the text and when.
          const alt =
            typeof node.properties.alt === "string" ? node.properties.alt : "";
          node.tagName = "span";
          node.properties = { className: ["blocked-image"] };
          node.children = [
            {
              type: "text",
              value: `[外部の画像は表示しません${alt ? `: ${alt}` : ""}] ${src}`,
            },
          ];
          return;
        }
        node.properties.src = rewriteImageSource(src, imageUrl);
      }
      node.children?.forEach(visit);
    };
    visit(tree);
  };
}

const processor = unified()
  .use(remarkParse)
  .use(remarkGfm)
  .use(remarkMath)
  .use(remarkRehype)
  .use(rehypeSanitize, {
    ...defaultSchema,
    attributes: {
      ...defaultSchema.attributes,
      code: [["className", /^language-./, "math-inline", "math-display"]],
    },
  })
  .use(rehypeContentImages)
  .use(rehypeKatex, { trust: false })
  .use(rehypeShiki, {
    // Both palettes are inlined; application.css switches to the dark one.
    themes: { light: "github-light", dark: "github-dark" },
    // Languages commonly used in troubleshooting reports; others fall back to text.
    langs: [
      "javascript",
      "typescript",
      "bash",
      "shellscript",
      "sh",
      "zsh",
      "console",
      "powershell",
      "json",
      "jsonc",
      "yaml",
      "toml",
      "ini",
      "xml",
      "html",
      "css",
      "markdown",
      "go",
      "python",
      "c",
      "cpp",
      "csharp",
      "java",
      "kotlin",
      "rust",
      "ruby",
      "php",
      "perl",
      "lua",
      "sql",
      "diff",
      "dockerfile",
      "nginx",
      "apache",
      "makefile",
      "hcl",
      "ssh-config",
      "systemd",
      "log",
    ],
    fallbackLanguage: "text",
  })
  .use(rehypeStringify);
export async function renderMarkdown(
  markdown: string,
  options: { imageUrl?: ImageUrl; blockExternalImages?: boolean } = {},
): Promise<string> {
  return String(
    await processor.process({
      value: markdown,
      data: {
        imageUrl: options.imageUrl,
        blockExternalImages: options.blockExternalImages,
      },
    }),
  );
}
