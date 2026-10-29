// Draws the ```mermaid blocks that markdown.ts leaves in rendered Markdown
// as SVG in the browser. Mermaid is large, so it is loaded only when a page
// has a diagram. Diagrams can be written by contestants: Mermaid runs with
// securityLevel "strict" (labels are text, no click handlers) and anything
// that would load from another host is removed from the SVG.
import type { Mermaid } from "mermaid";

let mermaidModule: Promise<Mermaid> | undefined;
// Mermaid's settings are global, so diagrams are drawn one batch at a time.
let queue: Promise<unknown> = Promise.resolve();
let nextId = 0;

function loadMermaid() {
  mermaidModule ??= import("mermaid").then((module) => module.default);
  return mermaidModule;
}

const outsideUrl = /url\(\s*['"]?(?!#)/i;

/**
 * Removes images, links and url() references to anything outside the SVG,
 * so a diagram never makes the reader's browser contact another host.
 */
export function removeExternalReferences(root: Element) {
  for (const element of root.querySelectorAll(
    "image, img, use, foreignObject img",
  ))
    element.remove();
  for (const link of root.querySelectorAll("a"))
    link.replaceWith(...link.childNodes);
  for (const style of root.querySelectorAll("style"))
    if (outsideUrl.test(style.textContent ?? "")) style.remove();
  for (const element of [root, ...root.querySelectorAll("*")])
    for (const attribute of [...element.attributes])
      if (
        outsideUrl.test(attribute.value) ||
        /^(xlink:)?href$|^src$/i.test(attribute.name)
      )
        element.removeAttribute(attribute.name);
}

/**
 * Draws every .mermaid-diagram under root in the current theme. A diagram
 * that cannot be parsed keeps showing its source with a note.
 */
export function renderDiagrams(root: ParentNode): Promise<void> {
  const diagrams = [...root.querySelectorAll<HTMLElement>(".mermaid-diagram")];
  if (!diagrams.length) return Promise.resolve();
  const run = queue.then(async () => {
    const mermaid = await loadMermaid();
    const dark = document.documentElement.dataset.theme === "dark";
    mermaid.initialize({
      startOnLoad: false,
      securityLevel: "strict",
      theme: dark ? "dark" : "default",
      fontFamily: getComputedStyle(diagrams[0]!).fontFamily,
    });
    for (const diagram of diagrams) {
      const source = diagram.querySelector(".mermaid-source")?.textContent;
      if (!source || !diagram.isConnected) continue;
      const id = `mermaid-${++nextId}`;
      let output = diagram.querySelector<HTMLElement>(".mermaid-output");
      if (!output) {
        output = document.createElement("div");
        output.className = "mermaid-output";
        diagram.append(output);
      }
      try {
        const { svg } = await mermaid.render(id, source);
        output.innerHTML = svg;
        removeExternalReferences(output);
        diagram.classList.add("is-rendered");
        diagram.classList.remove("is-failed");
      } catch {
        // Mermaid leaves its scratch element behind when parsing fails.
        document.getElementById(`d${id}`)?.remove();
        output.textContent =
          "図を描画できませんでした。記法を確認してください。";
        diagram.classList.add("is-failed");
        diagram.classList.remove("is-rendered");
      }
    }
  });
  queue = run.catch(() => undefined);
  return run;
}
