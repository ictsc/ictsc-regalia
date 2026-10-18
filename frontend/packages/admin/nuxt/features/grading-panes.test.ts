import { describe, expect, it } from "vitest";
import {
  defaultPaneLayout,
  resizePaneWidths,
  restorePaneLayout,
} from "./grading-panes";

describe("restorePaneLayout", () => {
  const fallback = defaultPaneLayout();
  it("keeps the defaults when nothing is saved", () => {
    expect(restorePaneLayout(null, null, fallback)).toEqual(fallback);
  });
  it("restores a valid saved layout and ignores invalid widths", () => {
    expect(
      restorePaneLayout(
        JSON.stringify({ collapsed: [1, 0, 0], widths: [30, 40, 30] }),
        null,
        fallback,
      ),
    ).toEqual({ collapsed: [true, false, false], widths: [30, 40, 30] });
    expect(
      restorePaneLayout(
        JSON.stringify({ collapsed: [0, 0, 0], widths: [5, 80, 15] }),
        null,
        fallback,
      ).widths,
    ).toEqual(fallback.widths);
  });
  it("migrates the legacy four-pane layout", () => {
    expect(
      restorePaneLayout(
        null,
        JSON.stringify({
          collapsed: [false, true, false, true],
          widths: [20, 15, 45, 20],
        }),
        fallback,
      ),
    ).toEqual({ collapsed: [false, false, true], widths: [20, 60, 20] });
  });
  it("throws on malformed JSON", () => {
    expect(() => restorePaneLayout("{", null, fallback)).toThrow();
  });
});

describe("resizePaneWidths", () => {
  it("moves the boundary but keeps both panes at least 12 wide", () => {
    expect(resizePaneWidths([23, 55, 22], 0, 1, 5)).toEqual([28, 50, 22]);
    expect(resizePaneWidths([23, 55, 22], 1, 2, 20)).toEqual([23, 65, 12]);
    expect(resizePaneWidths([23, 55, 22], 0, 1, -30)).toEqual([12, 66, 22]);
  });
});
