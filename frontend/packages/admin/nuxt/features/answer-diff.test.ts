import { describe, expect, it } from "vitest";
import {
  answerSourceLines,
  buildAnswerDiffHunks,
  compareAnswerSources,
} from "./answer-diff";

describe("unified answer diff", () => {
  it("puts removed lines before added lines and keeps both original line numbers", () => {
    const comparison = compareAnswerSources(
      "same\nold\nkeep\nremoved\nend",
      "same\nnew\nkeep\nend\nadded",
    );
    expect(
      comparison.rows.map(({ kind, previousLine, currentLine }) => [
        kind,
        previousLine,
        currentLine,
      ]),
    ).toEqual([
      ["same", 1, 1],
      ["removed", 2, null],
      ["added", null, 2],
      ["same", 3, 3],
      ["removed", 4, null],
      ["same", 5, 4],
      ["added", null, 5],
    ]);
    expect(comparison).toMatchObject({ addedLines: 2, removedLines: 2 });
  });

  it("groups distant changes into hunks with correct line ranges and omitted context", () => {
    const previous = Array.from(
      { length: 20 },
      (_, index) => `line ${index + 1}`,
    );
    const current = [...previous];
    current[6] = "replacement";
    current[17] = "another replacement";
    const comparison = compareAnswerSources(
      previous.join("\n"),
      current.join("\n"),
    );
    const hunks = buildAnswerDiffHunks(comparison.rows, 2);
    expect(hunks).toHaveLength(2);
    expect(hunks[0]).toMatchObject({
      previousStart: 5,
      previousCount: 5,
      currentStart: 5,
      currentCount: 5,
      omittedBefore: 4,
    });
    expect(hunks[1]).toMatchObject({
      previousStart: 16,
      previousCount: 5,
      currentStart: 16,
      currentCount: 5,
      omittedBefore: 6,
    });
    expect(buildAnswerDiffHunks(comparison.rows, Infinity)).toHaveLength(1);
  });

  it("handles repeated lines, blank lines and CRLF", () => {
    expect(answerSourceLines("a\r\n\r\nb\r")).toEqual(["a", "", "b", ""]);
    const comparison = compareAnswerSources(
      "same\nsame\n",
      "same\nnew\nsame\n",
    );
    expect(comparison.rows.filter((row) => row.text === "new")).toEqual([
      { kind: "added", previousLine: null, currentLine: 2, text: "new" },
    ]);
    expect(
      buildAnswerDiffHunks(compareAnswerSources("same", "same").rows),
    ).toEqual([]);
  });
});
