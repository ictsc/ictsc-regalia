import { onBeforeUnmount, ref } from "vue";

/** Spaces inserted for one Tab press in Markdown inputs. */
export const indentUnit = "    ";

/**
 * Makes Tab indent and Shift+Tab outdent the current line in a textarea, and
 * Escape leave the field so keyboard users are never trapped in it.
 */
export function createTabIndentHandler(
  options: { onIndent?: () => void } = {},
) {
  return (event: KeyboardEvent) => {
    const textarea = event.target as HTMLTextAreaElement;
    if (event.key === "Escape") {
      textarea.blur();
      return;
    }
    if (event.key !== "Tab" || event.isComposing) return;
    event.preventDefault();
    options.onIndent?.();
    const { selectionStart, selectionEnd, value } = textarea;
    if (!event.shiftKey) {
      replaceRange(textarea, selectionStart, selectionEnd, indentUnit);
      return;
    }
    const lineStart = value.lastIndexOf("\n", selectionStart - 1) + 1;
    const leading = /^ {1,4}|^\t/.exec(value.slice(lineStart));
    if (!leading) return;
    replaceRange(textarea, lineStart, lineStart + leading[0].length, "");
    const shift = leading[0].length;
    textarea.setSelectionRange(
      Math.max(lineStart, selectionStart - shift),
      Math.max(lineStart, selectionEnd - shift),
    );
  };
}

/**
 * Tab indentation for a textarea plus a short hint, shown the first time Tab
 * is used after the field gains focus, that Escape leaves the field.
 */
export function useTabIndent(hintMilliseconds = 4000) {
  const hintVisible = ref(false);
  let hintShown = false;
  let timer: ReturnType<typeof setTimeout> | undefined;
  const hideHint = () => {
    hintVisible.value = false;
    clearTimeout(timer);
  };
  const onKeydown = createTabIndentHandler({
    onIndent: () => {
      if (hintShown) return;
      hintShown = true;
      hintVisible.value = true;
      clearTimeout(timer);
      timer = setTimeout(hideHint, hintMilliseconds);
    },
  });
  const onBlur = () => {
    hintShown = false;
    hideHint();
  };
  onBeforeUnmount(() => clearTimeout(timer));
  return { onKeydown, onBlur, hintVisible };
}

// execCommand keeps the browser's undo history; setRangeText is the fallback.
function replaceRange(
  textarea: HTMLTextAreaElement,
  start: number,
  end: number,
  text: string,
) {
  textarea.setSelectionRange(start, end);
  const inserted =
    typeof document.execCommand === "function" &&
    document.execCommand(text ? "insertText" : "delete", false, text);
  if (!inserted) {
    textarea.setRangeText(text, start, end, "end");
    textarea.dispatchEvent(new Event("input", { bubbles: true }));
  }
}
