// Copy buttons for code blocks in rendered Markdown. They are added to the
// DOM that MarkdownContent creates with v-html, so they are plain elements,
// not Vue components, and must be added again after every render.
import { stripTerminalControls } from "./terminal-text";

// Static markup for buttons added inside rendered HTML, drawn like UiIcon.
const iconSvg = (paths: string) =>
  `<svg class="ui-icon" aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">${paths}</svg>`;

const states = {
  idle: {
    label: "コードをコピー",
    icon: iconSvg(
      '<rect x="8.5" y="8.5" width="11" height="11" rx="1.5"/><path d="M15.5 8.5V5.5A1.5 1.5 0 0 0 14 4H5.5A1.5 1.5 0 0 0 4 5.5V14a1.5 1.5 0 0 0 1.5 1.5h3"/>',
    ),
  },
  done: {
    label: "コピーしました",
    icon: iconSvg('<path d="M5 12.5l4.5 4.5L19 7.5"/>'),
  },
  failed: {
    label: "コピーできませんでした",
    icon: iconSvg('<path d="M6 6l12 12M18 6L6 18"/>'),
  },
};

function setState(button: HTMLButtonElement, state: keyof typeof states) {
  button.innerHTML = states[state].icon;
  button.setAttribute("aria-label", states[state].label);
  button.title = states[state].label;
}

/**
 * Wraps each code block of rendered (sanitized) HTML and adds a copy button
 * in its corner; the wrapper keeps the button put while the code scrolls.
 * Blocks that already have one are skipped.
 */
export function addCopyButtons(root: ParentNode) {
  for (const pre of root.querySelectorAll("pre")) {
    if (!pre.querySelector(":scope > code")) continue;
    if (pre.parentElement?.classList.contains("code-block")) continue;
    const wrapper = document.createElement("div");
    wrapper.className = "code-block";
    pre.replaceWith(wrapper);
    wrapper.append(pre);
    const button = document.createElement("button");
    button.type = "button";
    button.className = "code-copy";
    setState(button, "idle");
    wrapper.append(button);
  }
}

/**
 * Copies the code next to a copy button and shows the result on the button
 * for two seconds. Returns the timer so the caller can cancel it.
 */
export async function copyCode(
  button: HTMLButtonElement,
): Promise<ReturnType<typeof setTimeout> | undefined> {
  const code = button.parentElement?.querySelector("pre > code");
  if (!code) return undefined;
  try {
    // What is pasted into a terminal must be what is on screen.
    await navigator.clipboard.writeText(
      stripTerminalControls(code.textContent ?? ""),
    );
    setState(button, "done");
  } catch {
    setState(button, "failed");
  }
  return setTimeout(() => setState(button, "idle"), 2000);
}
