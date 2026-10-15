import { expect, it } from "vitest";
import { renderMarkdown } from "./markdown";
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
