import { describe, expect, it } from "vitest";
import {
  answerSourceLines,
  buildAnswerDiffHunks,
  compareAnswerSources,
} from "./answer-diff";
import { buildDisplayItems } from "./answer-display";

describe("buildDisplayItems", () => {
  const current = "a\nb\nc\nd\ne\nf\ng\nh\ni\nj";
  const previous = "a\nb\nc\nd\nE\nf\ng\nh\ni\nj";
  const lines = answerSourceLines(current);
  const comparison = compareAnswerSources(previous, current);

  it("lists every current line outside the diff view", () => {
    const items = buildDisplayItems(lines, comparison, [], false);
    expect(items).toHaveLength(10);
    expect(items[0]).toEqual({
      type: "line",
      key: "current-1",
      row: { kind: "same", previousLine: null, currentLine: 1, text: "a" },
    });
  });

  it("adds hunk headers and omitted-context markers", () => {
    const hunks = buildAnswerDiffHunks(comparison.rows, 1);
    const items = buildDisplayItems(lines, comparison, hunks, true);
    expect(items.map((item) => item.type)).toEqual([
      "omitted",
      "hunk",
      "line",
      "line",
      "line",
      "line",
      "omitted",
    ]);
    expect(items[0]).toMatchObject({ count: 3, key: "omitted-0" });
    expect(items.at(-1)).toMatchObject({ count: 4 });
  });

  it("returns nothing when the diff has no changes", () => {
    const same = compareAnswerSources(current, current);
    expect(buildDisplayItems(lines, same, [], true)).toEqual([]);
  });
});
