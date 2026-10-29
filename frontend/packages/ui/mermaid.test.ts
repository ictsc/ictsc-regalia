// @vitest-environment jsdom
import { expect, it } from "vitest";
import { removeExternalReferences } from "./mermaid";

it("removes everything in a diagram that would load from another host", () => {
  const root = document.createElement("div");
  root.innerHTML =
    '<svg><style>.a{fill:red}</style><image href="https://example.com/a.png"/>' +
    '<foreignObject><img src="https://example.com/b.png"></foreignObject>' +
    '<a href="https://example.com/"><text>link</text></a>' +
    '<path marker-end="url(#arrow)" style="fill:url(https://example.com/c)"/>' +
    "<style>.b{background:url(https://example.com/d)}</style></svg>";
  removeExternalReferences(root);
  expect(root.innerHTML).not.toContain("example.com");
  expect(root.innerHTML).toContain('marker-end="url(#arrow)"');
  expect(root.innerHTML).toContain(".a{fill:red}");
  expect(root.textContent).toContain("link");
});
