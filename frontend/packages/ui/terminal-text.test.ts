import { expect, it } from "vitest";
import { stripTerminalControls } from "./terminal-text";

it("removes terminal control sequences from text printed in a terminal", () => {
  expect(stripTerminalControls("\u001b[31merror\u001b[0m\tdone\n")).toBe(
    "error\tdone\n",
  );
  // Window title and clipboard writes (OSC), ended by BEL or ST.
  expect(
    stripTerminalControls("a\u001b]0;pwned\u0007b\u001b]52;c;eA==\u001b\\c"),
  ).toBe("abc");
  // A carriage return must not hide the start of a line.
  expect(stripTerminalControls("rm -rf /\rok\r\nnext")).toBe(
    "rm -rf /\nok\nnext",
  );
  expect(stripTerminalControls("x\u0008\u009b31my\u007fz")).toBe("x31myz");
  expect(stripTerminalControls("日本語のログ\n")).toBe("日本語のログ\n");
  expect(stripTerminalControls("ls \u202e# txt.exe\u202c")).toBe(
    "ls # txt.exe",
  );
});
