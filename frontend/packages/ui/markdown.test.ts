import { expect, it } from "vitest";
import { renderMarkdown } from "./markdown";
it("renders tables, math, and highlighted code while rejecting executable input", async () => {
  const html = await renderMarkdown(
    "| A | B |\n|---|---|\n|1|2|\n\n$x^2$\n\n```js\nconst x=1;\n```\n<script>alert(1)</script>\n[x](javascript:alert(1))",
  );
  expect(html).toContain("<table>");
  expect(html).toContain("katex");
  expect(html).toContain("shiki");
  expect(html).not.toContain("<script");
  expect(html).not.toContain('href="javascript:');
});
