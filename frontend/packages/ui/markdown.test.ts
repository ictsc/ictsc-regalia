import { expect, it } from "vitest";
import { renderMarkdown, setMarkdownImageUrl } from "./markdown";
it("renders GFM and highlighted code while rejecting executable input", async () => {
  const html = await renderMarkdown(
    "| A | B |\n|---|---|\n|1|2|\n\n- [x] Done\n\n~~old~~\n\n```js\nconst x=1;\n```\n<script>alert(1)</script>\n[x](javascript:alert(1))",
  );
  expect(html).toContain("<table>");
  expect(html).toContain('type="checkbox"');
  expect(html).toContain("<del>old</del>");
  expect(html).toContain("shiki");
  expect(html).not.toContain("<script");
  expect(html).not.toContain('href="javascript:');
});

it("keeps GFM line breaks while rendering math", async () => {
  const html = await renderMarkdown(
    "first\nsecond  \nthird\n\n$x^2$\n\n```text\n$x^2$\n```",
  );
  expect(html).toContain("first\nsecond<br>");
  expect(html).toContain("katex");
  expect(html).toContain("$x^2$");
});

it("serves repository images from the configured API path", async () => {
  setMarkdownImageUrl((name) => `/api/v1/contestant/images/${name}`);
  try {
    const html = await renderMarkdown(
      "![diagram](/images/problem1658-abc.png)\n\n![external](https://example.com/a.png)",
    );
    expect(html).toContain(
      'src="/api/v1/contestant/images/problem1658-abc.png"',
    );
    expect(html).toContain('src="https://example.com/a.png"');
    const pinned = await renderMarkdown("![d](/images/a.png)", {
      imageUrl: (name) => `/api/v1/admin/images/${name}?commit=abc`,
    });
    expect(pinned).toContain('src="/api/v1/admin/images/a.png?commit=abc"');
  } finally {
    setMarkdownImageUrl(undefined);
  }
});

it("highlights common troubleshooting languages", async () => {
  const html = await renderMarkdown(
    "```c\nint main(void) { return 0; }\n```\n\n```nginx\nserver { listen 80; }\n```",
  );
  expect(html.match(/class="shiki/g)).toHaveLength(2);
  // Light colors inline, dark ones in --shiki-dark for the dark theme.
  expect(html).toMatch(
    /<span style="color:#[0-9A-Fa-f]{6};--shiki-dark:#[0-9A-Fa-f]{6}">int<\/span>/,
  );
});

it("shows external images in contestant text as their URL", async () => {
  const html = await renderMarkdown(
    "![図](https://evil.example/t.png) ![構成](/images/net.png)",
    {
      blockExternalImages: true,
      imageUrl: (name) => `/api/v1/admin/images/${name}`,
    },
  );
  expect(html).not.toContain('src="https://evil.example');
  expect(html).toContain(
    "[外部の画像は表示しません: 図] https://evil.example/t.png",
  );
  expect(html).toContain('src="/api/v1/admin/images/net.png"');
});

it("leaves mermaid blocks as text for diagrams only when asked", async () => {
  const source = "```mermaid\ngraph LR\n  A --> B\n```";
  const diagram = await renderMarkdown(source, { mermaid: true });
  expect(diagram).toContain('<div class="mermaid-diagram">');
  expect(diagram).toContain("graph LR\n  A --> B");
  expect(diagram).not.toContain("shiki");
  const code = await renderMarkdown(source);
  expect(code).not.toContain("mermaid-diagram");
  expect(code).toContain("shiki");
});
