import { expect, it } from "vitest";
import { teamColors, teamStyle } from "./colors";
it("supports the reference palette and safe fallback", () => {
  expect(teamColors).toHaveLength(17);
  expect(teamStyle("invalid")["--team-color"]).toBe("#A6E35F");
  expect(teamStyle("#FFE000")["--team-text"]).toBe("#1D252D");
  expect(teamStyle("#3F43AD")["--team-text"]).toBe("#FFFFFF");
});
