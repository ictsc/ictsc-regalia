// Draws the ```mermaid blocks that markdown.ts leaves in rendered Markdown
// as SVG in the browser. Mermaid is large, so it is loaded only when a page
// has a diagram. Diagrams can be written by contestants, and Mermaid loads
// some resources while it draws, before its SVG reaches us. So a diagram
// that names an image or a stylesheet is not drawn at all, directives cannot
// change the styling or sanitizer settings, HTML in labels is limited to
// text formatting, and the finished SVG is stripped once more.
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

// Mermaid fetches node images (A@{ img: "..." }) and CSS resources itself.
const imageShape = /@\{[^}]*\bimg\b/i;
const cssResource = /\burl\s*\(|image-set\s*\(|@import/i;

/** Reports whether drawing source would make Mermaid load a resource. */
export function loadsResources(source: string): boolean {
  return imageShape.test(source) || cssResource.test(source);
}

// Mermaid's own secure keys plus everything that changes styling or the
// sanitizer, so %%{init}%% directives in a diagram cannot set them.
const secure = [
  "secure",
  "securityLevel",
  "startOnLoad",
  "maxTextSize",
  "suppressErrorRendering",
  "maxEdges",
  "dompurifyConfig",
  "theme",
  "themeCSS",
  "themeVariables",
  "darkMode",
  "fontFamily",
  "altFontFamily",
  "htmlLabels",
];
// Labels may use text formatting only: no images, media or styles.
const dompurifyConfig = {
  ALLOWED_TAGS: [
    "b",
    "strong",
    "i",
    "em",
    "u",
    "s",
    "sub",
    "sup",
    "small",
    "code",
    "br",
    "span",
    "div",
    "p",
  ],
  ALLOWED_ATTR: ["class"],
};

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
      secure,
      dompurifyConfig,
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
      if (loadsResources(source)) {
        output.textContent = "外部の画像やスタイルを読み込む図は表示しません。";
        diagram.classList.add("is-failed");
        diagram.classList.remove("is-rendered");
        continue;
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
