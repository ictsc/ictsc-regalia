// @vitest-environment jsdom
import { expect, it } from "vitest";
import { createTabIndentHandler, indentUnit } from "./textarea";

function setup(value: string, start: number, end = start) {
  const textarea = document.createElement("textarea");
  document.body.append(textarea);
  textarea.value = value;
  textarea.setSelectionRange(start, end);
  const handler = createTabIndentHandler();
  textarea.addEventListener("keydown", handler);
  const press = (key: string, shiftKey = false) => {
    const event = new KeyboardEvent("keydown", {
      key,
      shiftKey,
      cancelable: true,
    });
    textarea.dispatchEvent(event);
    return event.defaultPrevented;
  };
  return { textarea, press };
}

it("indents with spaces and outdents the current line", () => {
  const { textarea, press } = setup("code", 0);
  expect(press("Tab")).toBe(true);
  expect(textarea.value).toBe(`${indentUnit}code`);
  expect(press("Tab", true)).toBe(true);
  expect(textarea.value).toBe("code");
});

it("leaves the field on Escape and reports indentation", () => {
  const textarea = document.createElement("textarea");
  document.body.append(textarea);
  let indents = 0;
  textarea.addEventListener(
    "keydown",
    createTabIndentHandler({ onIndent: () => indents++ }),
  );
  textarea.focus();
  textarea.dispatchEvent(
    new KeyboardEvent("keydown", { key: "Tab", cancelable: true }),
  );
  expect(indents).toBe(1);
  textarea.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape" }));
  expect(document.activeElement).not.toBe(textarea);
});
