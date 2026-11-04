// @vitest-environment jsdom
import { expect, it } from "vitest";
import { loadsResources, removeExternalReferences } from "./mermaid";

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

it("refuses diagrams that would make Mermaid load a resource", () => {
  expect(
    loadsResources('flowchart LR\n  A@{ img: "https://example.com/a.png" }'),
  ).toBe(true);
  expect(
    loadsResources("flowchart LR\n  A@{ shape: rect, label: 'x' }\n  A --> B"),
  ).toBe(false);
  expect(
    loadsResources(
      '%%{init: {"themeCSS": "@import \'https://example.com/x.css\';"}}%%\nflowchart LR\n  A',
    ),
  ).toBe(true);
  expect(
    loadsResources(
      "flowchart LR\n  A\n  style A background:url(https://example.com/a)",
    ),
  ).toBe(true);
  expect(
    loadsResources(
      'flowchart LR\n  A\n  classDef c fill:image-set("https://example.com/a" 1x)',
    ),
  ).toBe(true);
  expect(
    loadsResources("flowchart LR\n  A[https://example.com を開く] --> B"),
  ).toBe(false);
});
