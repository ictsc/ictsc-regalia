import type { AnswerDiff, AnswerDiffHunk, AnswerDiffRow } from "./answer-diff";

export type DisplayItem =
  | { type: "line"; row: AnswerDiffRow; key: string }
  | { type: "hunk"; hunk: AnswerDiffHunk; key: string }
  | { type: "omitted"; count: number; key: string };

/**
 * Rows shown in the answer pane: every current line, or the diff hunks with
 * markers for omitted context. Returns an empty list when the diff has no
 * changes.
 */
export function buildDisplayItems(
  answerLines: readonly string[],
  comparison: AnswerDiff | null,
  hunks: readonly AnswerDiffHunk[],
  diffView: boolean,
): DisplayItem[] {
  if (!diffView || !comparison) {
    return answerLines.map((text, index) => ({
      type: "line",
      key: `current-${index + 1}`,
      row: { kind: "same", previousLine: null, currentLine: index + 1, text },
    }));
  }
  if (!hunks.length) return [];
  const items: DisplayItem[] = [];
  let shown = 0;
  for (const hunk of hunks) {
    if (hunk.omittedBefore) {
      items.push({
        type: "omitted",
        count: hunk.omittedBefore,
        key: `omitted-${shown}`,
      });
      shown += hunk.omittedBefore;
    }
    items.push({ type: "hunk", hunk, key: `hunk-${shown}` });
    for (const row of hunk.rows) {
      items.push({ type: "line", row, key: `diff-${shown++}` });
    }
  }
  const omittedAfter = comparison.rows.length - shown;
  if (omittedAfter)
    items.push({
      type: "omitted",
      count: omittedAfter,
      key: `omitted-${shown}`,
    });
  return items;
}
